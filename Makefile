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
	rm -rf build/
	rm -rf mocks/

server:
	$(info Building binary at the project root)
	$(info ===================================)
	$(info )
	env GOOS=linux GOARCH=amd64 go build -o build/marina_server marina/*.go
	env GOOS=linux GOARCH=amd64 go build -o build/marina_client internal/*.go

build-debug-server:
	$(info Building binary at the project root with gcflags "all=-N -l")
	env GOOS=linux GOARCH=amd64 go build -o build/marina_server_debug -gcflags "all=-N -l" marina/*.go

unit-tests:
	rm -rf ./mocks/
	bash hooks/test-services-finalize.sh

echo:
	$(info Building binary at the project root)
	$(info ===================================)
	$(info )
	go build -o build/echo_server services/echo_server.go
	go build -o build/echo_client client/echo_client.go

deps:
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/fiorix/protoc-gen-cobra
