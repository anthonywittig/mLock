runonprem:
	cd backend/onprem && go test ./...
	cd backend/onprem && go build
	mkdir -p backend/onprem/dist
	mv backend/onprem/onprem backend/onprem/dist
	cp backend/onprem/.env backend/onprem/dist
	./backend/onprem/dist/onprem

	# env GOOS=linux GOARCH=arm GOARM=5 go build
