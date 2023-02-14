#protoc -I=$HOME/go/pkg/mod/github.com/gogo/protobuf@v1.3.1:. --gogofaster_out=plugins=grpc,paths=source_relative:. \
api/proto/v1/*.proto

#protoc -I=$HOME/go/pkg/mod/github.com/gogo/protobuf@v1.3.1:. --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,paths=source_relative:. api/proto/v1/*-dto.proto
WorkDir=$(shell pwd)
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
 		--proto_path=.:$(WorkDir) \
		proto/*/*/*.proto



