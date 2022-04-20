SERVER_NAME=stratum
CLIENT_NAME=strat_xp
# build and install client and server
all: 
	go build -o $(CLIENT_NAME) ./cmd/client/client.go
	go build -o $(SERVER_NAME) ./cmd/server/main.go
	@echo "Build Successful, now run ./$(SERVER_NAME) to start server and ./$(CLIENT_NAME) to start client"

# test server methods
test:
	go test -v ./...


.PHONY: all