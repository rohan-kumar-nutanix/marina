build-protos:
	$(info Building protos)
	protoc \
  --go_out . --go_opt paths=source_relative \
  --go-grpc_out . --go-grpc_opt paths=source_relative \
	protos/marina/*.proto

build-protos-with-cobra:
	$(info Building protos)
	protoc \
  --go_out . --go_opt paths=source_relative \
  --go-grpc_out . --go-grpc_opt paths=source_relative \
  --cobra_out=plugins=client:. \
  protos/marina/*.proto

clean:
	find ./protos -name \*.pb.go -type f -exec rm -f {} +
	rm -f echo_server echo_client
	rm -f marina_server client

server:
	$(info Building binary at the project root)
	$(info ===================================)
	$(info )
	env GOOS=linux GOARCH=amd64
	go build -o marina_server marina/*.go
	go build -o marina_client client/*.go
#	go build -o marina_server services/marina/server.go
#	go build -o marina_client client/marina_grpc_client.go

echo:
	$(info Building binary at the project root)
	$(info ===================================)
	$(info )
	go build -o echo_server services/echo_server.go
	go build -o echo_client client/echo_client.go

deps:
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/fiorix/protoc-gen-cobra