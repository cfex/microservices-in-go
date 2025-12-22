PROTO_DIR := protobuf
OUT_DIR := ./services/common/pb
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

GOBIN := $(shell go env GOPATH)/bin

.PHONY: all tools gen proto clean

all: gen

tools:
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Done."

gen: proto

proto:
	@PATH="$(GOBIN):$$PATH" protoc \
	-I $(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

clean:
	@find $(OUT_DIR) -name "*.pb.go" -delete
