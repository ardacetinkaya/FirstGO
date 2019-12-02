# Start from the latest golang base image
FROM golang:latest

LABEL maintainer="Arda Cetinkaya"

WORKDIR /app

# Copy go mod and sum files
COPY go.mod  go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8090

# Command to run the executable
CMD ["./main"]