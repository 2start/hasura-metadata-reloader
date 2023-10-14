![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/2start/hasura-metadata-reloader/build-and-push-with-semver.yaml)

# Hasura Metadata Reloader

## About

At canida.io, we use Hasura to generate GraphQL APIs for most of our projects.
The following project provides a CLI tool and an API to call the Hasura Metadata API and reload the Hasura metadata.
If there are inconsistencies discovered in the metadata despite a reload, the inconsistencies are reported to Sentry.

We regularly faced problems with inconsistent state in Hasura for the following reasons:

- **Remote Schema temporarily unavailable**: Hasura just gives up to contact remote schemas after a short time. I.e. the remote schema may come back up but Hasura will not try to contact it again. Then, the Hasura API can not serve the remote schema.
- **Remote Schema error**: The remote schema may go out of sync because there is a bug in the service and Hasura cannot introspect the schema properly. The service may be fixed but Hasura already gave up on our service.
- The **underlying database is modified** by forces outside of our control. Let's say a customer runs a job to recreate certain database tables on a regular basis. This may lead to Hasura throwing inconsistency errors because it cannot find the table for a period of time and then it will not recover.

We use this tool as follows:
- After a Hasura remote schema starts up, we reload the metadata to ensure that Hasura picks up on the service again. This ensures auto-recovery after a service downtime.
- A Kubernetes CronJob regularly triggers reload metadata. If there are inconsistencies found, they will be reported to Sentry. Then, we receive a Slack notification to quickly fix the problem.

## Run

The application can be configured via env variables (see sample.env) or cli flags.
It can be run by building the go binary first or by just using the provided docker container.

There are two ways to trigger a reload of the Hasura metadata:

1. Use the CLI command `reload`.
2. Use the CLI command `server`. Then, call the endpoint `/reload_metadata`. The request method does not matter.

### Build & Run

```shell
go build -o hasura-metadata-reloader ./cmd
./hasura-metadata-reloader reload --hasura-endpoint=$HASURA_ENDPOINT --hasura-admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
./hasura-metadata-reloader server --hasura-endpoint=$HASURA_ENDPOINT --hasura-admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
```

### Docker
```shell
docker run ghcr.io/2start/hasura-metadata-reloader:latest hasura-metadata-reloader reload --hasura-endpoint=$HASURA_ENDPOINT --admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
docker run ghcr.io/2start/hasura-metadata-reloader:latest hasura-metadata-reloader server --hasura-endpoint=$HASURA_ENDPOINT --admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
```

### Sample Output without Metadata Inconsistencies

```shell
{"level":"info","time":"2023-10-12T19:32:56Z","message":"Sentry initialized."}
{"level":"info","is_consistent":true,"time":"2023-10-12T19:32:58Z","message":"Metadata is consistent."}
```

### Sample Output with Metadata Inconsistency

```shell
{"level":"info","time":"2023-10-12T18:13:23+02:00","message":"Sentry initialized."}
{"level":"error","error":"metadata is inconsistent","inconsistent_objects":[{"definition":{"comment":null,"definition":{"customization":{"root_fields_namespace":"supabase"},"forward_client_headers":true,"timeout_seconds":60,"url_from_env":"SUPABASE_AUTH_CONNECTOR_URL"},"name":"supabase-auth-connector","permissions":[],"remote_relationships":[]},"message":{"message":"Connection failure: Network.Socket.connect: <socket: 31>: does not exist (Connection refused)","request":{"host":"supabase-auth-connector.staging","method":"POST","path":"/graphql","port":80,"queryString":"","requestHeaders":{"Content-Type":"application/json","User-Agent":"hasura-graphql-engine/v2.28.1"},"responseTimeout":"ResponseTimeoutMicro 60000000","secure":false},"type":"http_exception"},"name":"remote_schema supabase-auth-connector","reason":"Inconsistent object: HTTP exception occurred while sending the request to \"SUPABASE_AUTH_CONNECTOR_URL\"","type":"remote_schema"}],"time":"2023-10-12T18:13:25+02:00","message":"Metadata is inconsistent."}
```

## Contributing

### pre commit hooks

I use pre-commit hooks to check several conditions before commiting e.g. code is formatted properly with gofmt. The pre-commit hooks are configured in `.pre-commit-config.yaml`.

Install [pre-commit](https://pre-commit.com/#install) first and then install the pre-commit hooks via `pre-commit install`.
