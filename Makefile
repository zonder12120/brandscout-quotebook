IMAGE_NAME=brandscout-quotebook
CONTAINER_NAME=quotebook-container
ENV_FILE=config\.env

build:
	@echo "Building the Docker image..."
	docker build -t $(IMAGE_NAME) .

run:
	@echo "Running the Docker container..."
	docker run -d --name $(CONTAINER_NAME) -p 8080:8080 --restart unless-stopped --env-file $(ENV_FILE) $(IMAGE_NAME)

stop:
	@echo "Stopping the container..."
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

restart: stop run

logs:
	docker logs -f $(CONTAINER_NAME)

clean:
	docker rmi $(IMAGE_NAME)

.PHONY: build run stop restart logs clean