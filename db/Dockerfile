from redis:alpine

EXPOSE 6379FROM golang:1.22-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/main .

# Expose the port  application runs on
EXPOSE 3000

CMD ["./main"]