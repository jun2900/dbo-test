# Project dbo-test

Project for DBO Interview

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Environment Setup

Before you start, rename the `.env.example` file to `.env` and update the environment variables as needed.

### Swagger Documentation

After running the application, you can access the Swagger documentation by navigating to the following URL in your browser (the port is in the .env file):
[http://localhost:[your_port]/swagger/index.html](http://localhost:[your_port]/swagger/index.html)

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```