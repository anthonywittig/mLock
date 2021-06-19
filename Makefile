deploy-backend: gen-protos
	./backend/bin/deploy.sh

deploy-frontend:
	./frontend/bin/deploy.sh

gen-protos:
	rm -rf backend/shared/protos/*
	protoc  -I=backend --go_out=backend/shared/protos backend/protos/*
	cd backend && go mod tidy

build-onprem-mac: gen-protos
	cd backend/onprem && go test ./...
	cd backend/onprem && go build -o mlock-onprem
	mkdir -p backend/onprem/dist
	mv backend/onprem/mlock-onprem backend/onprem/dist
	cp backend/onprem/.env backend/onprem/dist

build-onprem-rpi: gen-protos
	cd backend/onprem && go test ./...
	cd backend/onprem && env GOOS=linux GOARCH=arm GOARM=5 go build -o mlock-onprem
	mkdir -p backend/onprem/dist
	mv backend/onprem/mlock-onprem backend/onprem/dist
	cp backend/onprem/.env backend/onprem/dist

run-onprem: build-onprem-mac
	./backend/onprem/dist/mlock-onprem

run-onprem-tests-integ:
	cd backend && export $(cat onprem/.env | sed 's/#.*//g' | xargs) && go test mlock/onprem/hab