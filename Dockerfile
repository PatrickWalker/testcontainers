# Start from golang base image
FROM golang:alpine as builder

# Add Maintainer info
LABEL maintainer="pwok"

# Install git.
# Git is required for fetching the dependencies
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependecies. Dependencies will be cached if the go.mod and the go.sum files are not changes
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage. Observe we also copied the app.env file
COPY --from=builder /app/main .

# Exposet post 8080 to the outside world
EXPOSE 8080 8080

# Command to run the executable
CMD ["./main"]