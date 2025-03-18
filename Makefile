.PHONY: go docker clean deepClean

# Run the application locally
go:
	go run main.go

# Build and run the Docker container
docker:
	PORT=$$(go run main.go -print-port) && \
	docker build --build-arg PORT=$$PORT -t forum . && \
	docker run -e PORT=$$PORT -p $$PORT:$$PORT --name forum \
		-v $(PWD)/database:/app/database \
		--rm forum

# Stop and clean up Docker resources
clean:
	-docker stop forum || true
	-docker rmi forum || true
	-docker system prune -f --volumes

# Stop and clean up All Docker resources
deepClean:
	-docker stop $$(docker ps -aq) 
	-docker rm $$(docker ps -aq) 
	-docker rmi $$(docker images -q)
	-docker system prune -a -f --volumes
