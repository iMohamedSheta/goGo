# Generate Swagger Specification documentation in YAML format
swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
		
swagger-serve:
	GO111MODULE=off swagger serve -F=swagger ./swagger.yaml


# Migrations 
migrate:
	goose up

migrate-rollback:
	goose down

migrate-reset:
	goose reset

migrate-status:
	goose status

# Create a new migration file 
type ?= sql
migration:
	goose create $(name) $(type)

# Define directories and proto files
PROTO_FILES := $(wildcard proto/*.proto)
OUT := grpc

generate-grpc:
	protoc --go_out=$(OUT) --go-grpc_out=$(OUT) $(PROTO_FILES)


# install dependencies linux
install:
# install protobuf compiler
	apt install -y protobuf-compiler
# install go modules
	go mod tidy
	go mod vendor

# install dependencies windows
install-windows:
# install protobuf compiler
	winget install protobuf
# install go modules
	go mod tidy






