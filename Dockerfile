# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##
FROM golang:1.21-alpine AS build
RUN adduser -D -u 1001 gouser

WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod go.sum ./

# download Go modules and dependencies
RUN go mod download

# copy all files ending with .go
COPY *.go ./
COPY ./cmd ./cmd
COPY ./internal ./internal

# compile application -s -w flags to reduce binary size
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/hasura-metadata-reloader ./main.go

##
## STEP 2 - DEPLOY
##
FROM scratch
LABEL org.opencontainers.image.description "A tool to reload hasura metadata and report inconsistencies."

# Copy the binary.
COPY --from=build /app/bin/hasura-metadata-reloader /bin/hasura-metadata-reloader
# Copy the /etc/passwd file created in the build stage to use the gouser in the scratch image
COPY --from=build /etc/passwd /etc/passwd
# Copy CA certificates o.w. SSL calls to external services will fail. 
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 1001

EXPOSE 8080

CMD ["/hasura-metadata-reloader"]
