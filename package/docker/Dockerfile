FROM docker.dyn.ntnxdpro.com/ntnx-general-docker/golang:1.19 as build

LABEL maintainer="acropolis-catalog@nutanix.com"

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV PROTOC_ZIP=protoc-3.15.8-linux-x86_64.zip
ENV GOPRIVATE=github.com/nutanix-core,github.com/nutanix

RUN apt-get update && apt-get install -y unzip

# Install protoc compiler and required bindings for protocol buffers.
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/$PROTOC_ZIP \
    && unzip protoc-3.15.8-linux-x86_64.zip -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \
    && rm -f $PROTOC_ZIP

RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.4.0

# Download public key for github.com
# Setup git and ssh to be able to clone from Nutanix private github repos
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN git config --global url."git@github.com:".insteadOf "https://github.com"

WORKDIR /go/src/github/marina

COPY . .

# The mount prefix here will allow the command to use ssh keys to authenticate to Nutanix github.
RUN --mount=type=secret,id=ssh,dst=/root/.ssh/id_rsa make build-protos
RUN --mount=type=secret,id=ssh,dst=/root/.ssh/id_rsa go get -d -v ./...
RUN --mount=type=secret,id=ssh,dst=/root/.ssh/id_rsa make server

# Start a new image and copy the binary file.
FROM docker.dyn.ntnxdpro.com/ntnx-general-docker/alpine:3.11.6

COPY --from=build /go/src/github/marina/marina_server /usr/local/bin/marina
COPY --from=build /go/src/github/marina/marina_client /usr/local/bin/marina_client
WORKDIR /

EXPOSE 9200 9201 9202

CMD ["marina"]
