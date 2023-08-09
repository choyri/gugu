GIT_TAG	:= $(shell git describe --tags || echo dev)

GO_LDFLAGS 	:= -s -w -X main.version=$(GIT_TAG)

init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go mod download

proto:
	protoc \
        --proto_path ./api \
        --go_out ./api \
            --go_opt paths=source_relative \
        --go-grpc_out ./api \
            --go-grpc_opt paths=source_relative \
        --grpc-gateway_out ./api \
            --grpc-gateway_opt paths=source_relative \
        api/*/v*/*.proto

go-build:
	go build -o bin/ -trimpath -ldflags="$(GO_LDFLAGS)"
