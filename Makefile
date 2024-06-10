.PHONY: compile
compile:
	protoc proto/stat/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

build:
	GOOS=linux GOARCH=amd64 go build -o bin/server server/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/client client/*.go
