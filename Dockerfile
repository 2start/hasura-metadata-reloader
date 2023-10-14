# Name the build stage as "build"
# Using golang:1.21-alpine image as the base image
FROM golang:1.21-alpine AS build

# Add a new user "gouser" with user id 1001, -D flag is for disabling password for this user
RUN adduser -D -u 1001 gouser

# Set the Working Directory inside the Docker container, each command that follows will be executed in this directory
WORKDIR /app

# Copy Go modules manifests (go.mod, go.sum file) present in the root to the image
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy all relevant Go source code files (main.go), cmd and internal directory to the Docker image
COPY main.go ./
COPY ./cmd ./cmd
COPY ./internal ./internal

# Build the Go application, -s -w reduces the size of the binary by omitting the symbol table and debug information
# The binary is output to ./bin/hasura-metadata-reloader
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/hasura-metadata-reloader ./main.go


# Scratch is a minimalist Docker image with no operating system.
# It's used when you want to distribute a self-contained application binary.
# The final image will be just a few megabytes in size (last checked < 5MiB).
FROM scratch

# Add metadata to the image
LABEL org.opencontainers.image.source=https://github.com/2start/hasura-metadata-reloader
LABEL org.opencontainers.image.description="Hasura Metadata Reloader"

# Copy the binary from the "build" stage (/app/bin/hasura-metadata-reloader) to the "deploy" stage (/bin/hasura-metadata-reloader)
COPY --from=build /app/bin/hasura-metadata-reloader /bin/hasura-metadata-reloader 

# Copy the /etc/passwd file created in the build stage to have the same users in the scratch image
COPY --from=build /etc/passwd /etc/passwd

# Copy CA certificates to support SSL calls to external services
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Change to user 1001 that we created at the first stage
USER 1001

# Expose port 8080 to allow communication to/from server 
EXPOSE 8080

# Command to run the server
CMD ["/bin/hasura-metadata-reloader"]
