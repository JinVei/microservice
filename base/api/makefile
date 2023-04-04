
WorkDir=$(shell pwd)
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
 		--proto_path=.:$(WorkDir) \
		--go-grpc_opt=require_unimplemented_servers=false \
		proto/*/*/*.proto



