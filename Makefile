# Makefile for compiling client and server

# Module path
MODULE_PATH=github.com/isgo-golgo13/fifochannel

# Directories
CLIENT_DIR=cmd/client
SERVER_DIR=cmd/server

# Output binaries
CLIENT_BIN=client
SERVER_BIN=server

# Docker image names and tags
CLIENT_IMAGE=isgo-golgo13/fifochannel-client
SERVER_IMAGE=isgo-golgo13/fifochannel-server
IMAGE_TAG=latest

# Compile the client
compile-client:
	@echo "Compiling client..."
	@go build -o $(CLIENT_DIR)/$(CLIENT_BIN) $(MODULE_PATH)/$(CLIENT_DIR)

# Compile the server
compile-server:
	@echo "Compiling server..."
	@go build -o $(SERVER_DIR)/$(SERVER_BIN) $(MODULE_PATH)/$(SERVER_DIR)

# Compile both client and server
compile-all: compile-client compile-server
	@echo "Compilation complete."

# Clean up binaries
clean:
	@echo "Cleaning up binaries..."
	@rm -f $(CLIENT_DIR)/$(CLIENT_BIN) $(SERVER_DIR)/$(SERVER_BIN)
	@echo "Clean up complete."

# Build Docker image for client
docker-image-client:
	@echo "Building Docker image for client..."
	@docker build -t $(CLIENT_IMAGE):$(IMAGE_TAG) -f Dockerfile.client .

# Build Docker image for server
docker-image-server:
	@echo "Building Docker image for server..."
	@docker build -t $(SERVER_IMAGE):$(IMAGE_TAG) -f Dockerfile.server .

# Build Docker images for both client and server
docker-image-all: docker-image-client docker-image-server
	@echo "Docker images built and tagged."

# Phony targets
.PHONY: compile-client compile-server compile-all clean docker-image-client docker-image-server docker-image-all
