.PHONY: compile
compile:
	protoc proto/stat/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.
<<<<<<< Updated upstream
=======

build:
	goreleaser --snapshot --skip-publish --rm-dist
>>>>>>> Stashed changes
