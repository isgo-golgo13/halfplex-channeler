# Makefile for compiling client and server

# Module path
MODULE_PATH=github.com/isgo-golgo13/fifochannel

# Directories
CLIENT_DIR=cmd/client
SERVER_DIR=cmd/server

# Output binaries
CLIENT_BIN=client
SERVER_BIN=server

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

# Phony targets
.PHONY: compile-client compile-server compile-all clean
