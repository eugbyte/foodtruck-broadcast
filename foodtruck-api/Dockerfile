# syntax=docker/dockerfile:1

FROM golang:1.20
WORKDIR /app

# Update the linux OS
RUN apt-get update && apt-get install make && apt-get install sudo

# Download Go modules
COPY . .
RUN ls
RUN go mod download

# Build the executable
ARG SRC_PATH
RUN env GO111MODULE=on STAGE=${STAGE} GOOS=linux go build -o main -ldflags="-s -w" ${SRC_PATH}
RUN ls

# The command to be executed when a Docker container starts.
EXPOSE 8080
CMD ["./main"]
