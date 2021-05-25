deploy-backend: gen-protos
	./backend/bin/deploy.sh

deploy-frontend:
	./frontend/bin/deploy.sh

gen-protos:
	rm -rf backend/shared/protos/*
	protoc  -I=backend --go_out=backend/shared/protos backend/protos/*
	cd backend && go mod tidy

run-onprem: gen-protos
	cd backend/onprem && go test ./...
	cd backend/onprem && go build
	mkdir -p backend/onprem/dist
	mv backend/onprem/onprem backend/onprem/dist
	cp backend/onprem/.env backend/onprem/dist
	./backend/onprem/dist/onprem

	# env GOOS=linux GOARCH=arm GOARM=5 go build
