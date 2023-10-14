# Hasura Metadata Reloader

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/2start/hasura-metadata-reloader/build-and-push-with-semver.yaml)

## Overview

This project offers a CLI tool and an API for invoking Hasura Metadata API, facilitating metadata reload. Any discrepancies in the metadata, persisting post-reload, are reported to Sentry.

Issues provoking an inconsistent state in Hasura include:

- **Transient non-availability of Remote Schema**: Hasura's temporary inability to connect to remote schemas that eventually recover
- **Errors in Remote Schema**: Bugs in the service causing Hasura to have issues introspecting the schema
- **Modifications in the underlying database by external catalysts**: Jobs run by customers to recreate certain database tables periodically can cause inconsistency issues in Hasura due to temporary unavailability of tables

We utilize this tool to:

- Ensure auto-recovery after service downtime by reloading metadata once a Hasura remote schema reboots
- Trigger a metadata reload with Kubernetes CronJob, sending inconsistency reports to Sentry, followed by a prompt Slack notification for quick troubleshooting

## Usage

The application runs through the configuration of environment variables (refer sample.env) or CLI flags. You can either run it by creating the Go binary first or by simply running the given Docker container.

There are two methods to prompt a reload of the Hasura metadata:

1. The CLI command `reload`.
2. The CLI command `server`, paired with the endpoint call `/reload_metadata` (Request method is irrelevant).

### Compilation & Execution

```shell
go build -o hasura-metadata-reloader ./cmd
./hasura-metadata-reloader reload --hasura-endpoint=$HASURA_ENDPOINT --hasura-admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
./hasura-metadata-reloader server --hasura-endpoint=$HASURA_ENDPOINT --hasura-admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
```

### Docker Execution

```shell
docker run ghcr.io/2start/hasura-metadata-reloader:latest hasura-metadata-reloader reload --hasura-endpoint=$HASURA_ENDPOINT --admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
docker run ghcr.io/2start/hasura-metadata-reloader:latest hasura-metadata-reloader server --hasura-endpoint=$HASURA_ENDPOINT --admin-secret $HASURA_ADMIN_SECRET --sentry-dsn $SENTRY_DSN
```

### Sample Outputs

#### Without Metadata Inconsistencies

```shell
{"level":"info","time":"2023-10-12T19:32:56Z","message":"Sentry initialized."}
{"level":"info","is_consistent":true,"time":"2023-10-12T19:32:58Z","message":"Metadata is consistent."}
```

#### With Metadata Inconsistency

```shell
{"level":"info","time":"2023-10-12T18:13:23+02:00","message":"Sentry initialized."}
{"level":"error","error":"metadata is inconsistent","inconsistent_objects":[{"definition":{"comment":null,"definition":{"customization":{"root_fields_namespace":"supabase"},"forward_client_headers":true,"timeout_seconds":60,"url_from_env":"SUPABASE_AUTH_CONNECTOR_URL"},"name":"supabase-auth-connector","permissions":[],"remote_relationships":[]},"message":{"message":"Connection failure: Network.Socket.connect: <socket: 31>: does not exist (Connection refused)","request":{"host":"supabase-auth-connector.staging","method":"POST","path":"/graphql","port":80,"queryString":"","requestHeaders":{"Content-Type":"application/json","User-Agent":"hasura-graphql-engine/v2.28.1"},"responseTimeout":"ResponseTimeoutMicro 60000000","secure":false},"type":"http_exception"},"name":"remote_schema supabase-auth-connector","reason":"Inconsistent object: HTTP exception occurred while sending the request to \"SUPABASE_AUTH_CONNECTOR_URL\"","type":"remote_schema"}],"time":"2023-10-12T18:13:25+02:00","message":"Metadata is inconsistent."}
```

## Contributions

### Pre-commit Hooks

Pre-commit hooks are implemented to ensure the fulfillment of several conditions before committing, like properly formatting the code with "gofmt". The hooks' configuration can be found in `.pre-commit-config.yaml`.

The first step is to install [pre-commit](https://pre-commit.com/#install), followed by the installation of the pre-commit hooks, which can be achieved through `pre-commit install`.
