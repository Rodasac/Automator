# Automator - Automate your complex tasks

## What is Automator?

Automator is a project that aims to automate the monitoring and automation of tasks on every website (apps may be
supported in the future). It is a project that is still in development and is not ready for production use.

## How does it work?

Automator uses a simple and easy to use API to create tasks and monitor them. It's composed of 4 main components:

- The **Automator Robot** - The robot is the main component of the project. It is responsible for executing the tasks
  and monitoring them. It is also responsible for saving resources and data to the correct storage. 
- The **Automator Queue** - The queue orchestrate the tasks execution, it is responsible for sending the tasks to the
  robot through a queue system.
- The **Automator Authenticator** - The authenticator is responsible for confirming the tasks execution. It is
  responsible for checking if the tasks were executed correctly and if not and validates the resulting resources.
- The **Automator API** - The API is responsible for serving the tasks to the users, to managing all the system
  data/configuration and to manage the users and their permissions.

## Development status

| Component     | Status         |
|---------------|----------------|
| Robot         | ⚠️ In Progress |
| Queue         | ⭕️ Not started |
| Authenticator | ⭕️ Not started |
| API           | ⭕️ Not started |

## Requirements

- Docker (optional)
- Docker Compose (optional)
- Go 1.21+
- PostgreSQL 15+ (older versions may work)
- Redis 6+ (older versions may work)
- RabbitMQ 3.12+ (older versions may work)
- A CDP compatible web browser (Chrome, Firefox, Edge, etc), if not present rod (the base library used by the robot)
  will download a compatible version of Chromium.

## How to run (meant for development)

1. Clone the repository
2. Run `docker-compose up -d` to start the database and queue services.
3. Copy the .env.template file to .env inside every service and fill the variables.
4. Run robot migrations `cd robot && go run main/db/cli.go init && go run main/db/cli.go migrate`
5. Start the robot `go run main/file_automator/main.go`
6. Start the robot grpc server `go run main/grpc_server/main.go` (starts on port 50051, you can see grpc/media.proto for
   the available methods)
