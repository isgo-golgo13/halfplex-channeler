# HalfPlexChanneler (Go)
Go 1.22 Half Duplex (Fifo) HalfPlexChannel interface using io.Write, io.Read, io.Close for HalfPlexChanneler structs for client and server for data exchange.


## Source ToC (Table of Contents)

- server/server.go - The `server` imports `svckit`package to use HalfPlexChannel interface and HalfPlexChanneler struct.
- svckit/svckit.go - The `svckit`package provides SoC-style HalfPlexChannel interface and HalfPlexChanneler struct.
- client/client.go - The `client` imports `svckit`package to use HalfPlexChannel interface and HalfPlexChanneler struct.

The `svckit` package provides a Separation of Concerns-style API with its HalfPlexChannel interface to allow sending and receving
of any data source that fulfills the `io.Reader`interface and `io.Writer` interface and defines a `Close` API. The APIs additionally
include a Go `context.Context` parameter to disallow deadline exceeding and cancellation control. The `svckit` holds all the
FIFO half-plex connection, read, write and device closing logic and prevents code clutter of this logic into the client and server sources.

The `client` uses the `HalfPlexChannel` interface as the handle to send and receive data exchanges to the server and constructs
the `HalfPlexChanneler` struct instance assigning it to the interface handle.

The `server` uses the `HalfPlexChannel` interface as the handle to send and receive data exchanges to the client and constructs
the `HalfPlexChanneler` struct instance assigning it to the interface handle.

## Instructions
The included `Makefile` provides targets to do the following.

- Compile the `client` client.go binary individually.
- Compile the `server` server.go binary individually.
- Compile both the `client` and `server` go sources into binaries.
- Build Docker image for the client (see client/Dockerfile) with latest tag.
- Build Docker image for the server (see server/Dockerfile) with latest tag.
- Build both Docker images for client and server simulateneously with latest tags.

For reference here is the `Makefile`.

```Makefile
# Variables
PROJECT_NAME := halfplex
CLIENT_DIR := ./client
SERVER_DIR := ./server
DOCKER_COMPOSE_FILE := ./docker-compose.yml

# Build the client binary
build-client:
	@echo "Building client..."
	cd $(CLIENT_DIR) && go build -o client client.go

# Build the server binary
build-server:
	@echo "Building server..."
	cd $(SERVER_DIR) && go build -o server server.go

# Build client Docker image
docker-build-client:
	@echo "Building client Docker image..."
	docker build -t $(PROJECT_NAME)-client $(CLIENT_DIR)

# Build server Docker image
docker-build-server:
	@echo "Building server Docker image..."
	docker build -t $(PROJECT_NAME)-server $(SERVER_DIR)

# Build all Docker images
docker-build-all: docker-build-client docker-build-server

# Run Docker Compose
docker-up:
	@echo "Starting Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

# Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f $(CLIENT_DIR)/client $(SERVER_DIR)/server
	docker system prune -f

# Full build and deploy
deploy: docker-build-all docker-up

# Full clean up
clean-all: docker-down clean
```



To compile the `client` and `server`.

```sh
make build-client && make build-server
```


To build the Docker images for the `client` and `server` do the following.

```sh
make docker-build-all
```

To build the Docker image individually for the `client` do the following.

```sh
make docker-build-client
```

To build the Docker image individually for the `server` do the following.

```sh
make docker-build-server
```
