VERSION := 0.1.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: help robot-init-db robot-migrate-db robot-rollback-db robot-start-db unit-tests integration-tests test-coverage vet static-checks fmt-check check clean build-robot-stream-automator build-robot-grpc-server build-robot-file-automator build-robot start-robot-stream-automator start-robot-grpc-server start-robot robot-init

help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

robot-init-db: # Initialize the database.
	go run cmd/db/cli.go db init

robot-migrate-db: # Migrate the database.
	go run cmd/db/cli.go db migrate

robot-rollback-db: # Rollback the database.
	go run cmd/db/cli.go db rollback

robot-start-db: robot-init-db robot-migrate-db # Initialize and migrate the database.

unit-tests: # Run unit tests.
	go test -v ./entities/... ./usecases/...

integration-tests: # Run integration tests.
	go test -v ./adapters/...

test-coverage: # Run tests with coverage.
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

vet: # Run go vet.
	go vet ./...

static-checks: # Run static checks.
	staticcheck ./...

fmt-check: # Run gofmt check.
	gofmt -l -e .

check: fmt-check vet static-checks unit-tests integration-tests # Run all checks.

clean: # Clean the project.
	go clean
	rm -rf bin

build-robot-stream-automator: # Build the robot stream automator.
	go build $(LDFLAGS) -o bin/robot-stream-automator cmd/stream_automator/main.go

build-robot-grpc-server: # Build the robot gRPC server.
	go build $(LDFLAGS) -o bin/robot-grpc-server cmd/grpc_server/main.go

build-robot-file-automator: # Build the robot file automator.
	go build $(LDFLAGS) -o bin/robot-file-automator cmd/file_automator/main.go

build-robot: clean build-robot-stream-automator build-robot-grpc-server build-robot-file-automator # Build the robot.

start-robot-stream-automator: # Start the robot stream automator consumer.
	bin/robot-stream-automator

start-robot-grpc-server: # Start the robot gRPC server.
	bin/robot-grpc-server

start-robot: # Start the robot stream automator consumer and the gRPC server.
	bin/robot-stream-automator &
	bin/robot-grpc-server

robot-init: build-robot start-robot # Build and start the robot stream automator consumer and the gRPC server.
