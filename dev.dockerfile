# Go Version
FROM golang:1.20-alpine

SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Environment variables which CompileDaemon requires to runs
ENV PROJECT_DIR=/app \
  GO111MODULE=on \
  CGO_ENABLED=0

# Basic setup of the container
RUN mkdir /app
COPY .. /app
WORKDIR /app

# Install dependencies
RUN apk --no-cache add curl=8.4.0-r0 && \
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
  mv migrate /usr/bin/ && \
  go get github.com/githubnemo/CompileDaemon && \
  go install github.com/githubnemo/CompileDaemon && \
  curl -L https://github.com/go-task/task/releases/download/v3.31.0/task_linux_amd64.tar.gz | tar xvz && \
  mv task /usr/bin/

# Build and start listening for file changes
ENTRYPOINT ["task", "dev"]
