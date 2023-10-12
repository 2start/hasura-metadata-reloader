FROM golang:alpine AS build

RUN apk --no-cache add \
    gcc \
    g++ \
    make \
    git


WORKDIR /go/src/app

COPY ./cmd ./cmd
COPY ./internal ./internal


COPY go.mod .
COPY go.sum .
COPY *.go ./
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/hasura-metadata-reloader ./main.go


FROM alpine:3.17
LABEL org.opencontainers.image.description "A tool to reload hasura metadata and report inconsistencies."

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /go/src/app/bin /app

ENV PATH="/app:${PATH}"
CMD ["/app/hasura-metadata-reloader"]

