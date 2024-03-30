FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app

# copy necessary files into build container
COPY commonlib ./commonlib
COPY microsvc ./microsvc
COPY vendor ./vendor

# Copy go mod and sum files
COPY go.mod go.sum ./

WORKDIR /app/microsvc

# Build the Go app
RUN go build -buildvcs=false -o ./microsvc ./

# Start a new stage from scratch
FROM alpine:3.19
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/microsvc/microsvc .

#Command to run the executable
CMD ["./microsvc"]