# This file should contain our go build and go test stuff.
curl -o mockery.tar.gz -L'#' https://github.com/vektra/mockery/releases/download/v2.15.0/mockery_2.15.0_Linux_x86_64.tar.gz
tar -xf mockery.tar.gz mockery && sudo mv mockery /usr/local/bin/mockery
sudo chmod +x /usr/local/bin/mockery

curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
unzip protoc-3.15.8-linux-x86_64.zip -d /home/circleci/protoc
export PATH=/home/circleci/protoc/bin:$PATH

export PATH=$PATH:/home/circleci/project/.workspace/build-tools/golang-go1.19/go/bin
go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.4.0
go install gotest.tools/gotestsum@latest
export PATH=$PATH:/home/circleci/go/bin

make build-protos