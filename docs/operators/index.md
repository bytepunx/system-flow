---
title: Operators guide
updated: 2026-09-15
status: draft
---

# system-flow for operators

## Running the dashboard

`flai dashboard` runs `ghcr.io/bytepunx/flaiover` detached as `flaiover-<project>` with the repository mounted read-write at `/project`, bound to `127.0.0.1:4242` by default, as the invoking user. `flai dashboard status`, `logs`, and `stop` manage it. Change the image, tag, or port in `~/.flai/config.json` or per project in `system-flow.yaml` under `dashboard`.

Without `flai`, the equivalent is:

```bash
docker run --detach --rm --name flaiover-myproject \
  --publish 127.0.0.1:4242:3000 --volume "$PWD:/project" --env PROJECT_DIR=/project \
  --user "$(id -u):$(id -g)" ghcr.io/bytepunx/flaiover:latest
```

## The image

`ghcr.io/bytepunx/flaiover` is built from `flaiover/Dockerfile` at the repository root: a `flai` binary from the same commit at `/usr/local/bin/flai`, the SvelteKit build on `node:24-alpine`, listening on `3000`, running as an unprivileged user by default and working under any `--user`. Tags: `latest` (main), `X.Y.Z` and `X` from `flaiover/vX.Y.Z` release tags, and `sha-<commit>`. While the repository is private the package is too: `docker login ghcr.io` with a token that has `read:packages` before `flai dashboard` can pull it. The container writes flai's config and cache under the mounted project's `.flai-cache/`, which is git-ignored.

Build locally with `make flaiover-image` (tag `flaiover:local`) and run it with `flai dashboard --image flaiover --tag local`.

## Telemetry

flaiover follows the logging and telemetry conventions.

| Signal | Where | Notes |
|--------|-------|-------|
| Logs | stdout, one JSON event per line | `ts`, `level`, `service`, `component`, `msg`, and fields; one line per request with `trace_id` (or `request_id` when tracing is off), `method`, `route`, `path`, `status`, `duration_ms`; `LOG_LEVEL` (`debug`, `info`, `warn`, `error`), `LOG_FORMAT=text` for key-value text outside production |
| Liveness | `GET /_health` | Always 200 while the process runs; touches nothing |
| Readiness | `GET /_ready` | Checks the mounted project's manifest and item listing with a two second timeout each and reports whether the bundled flai is available; 503 names the failing check |
| Metrics | `GET /metrics` | Prometheus format: `flaiover_http_requests_total{method,route,status}`, `flaiover_http_request_duration_seconds` (histogram), `flaiover_http_requests_in_flight`, `flaiover_build_info{version,commit}` (the flaiover release tag and commit baked into the image), plus Node process metrics prefixed `flaiover_`. Route labels are SvelteKit route ids, never paths with IDs |
| Traces | OTLP/HTTP | Exported only when `OTEL_EXPORTER_OTLP_ENDPOINT` is set; one server span per request named `<method> <path>` with W3C context taken from the incoming headers; standard `OTEL_*` variables apply (`OTEL_SERVICE_NAME`, `OTEL_TRACES_SAMPLER`, `OTEL_TRACES_SAMPLER_ARG`) |

Local stack: `PROJECT=$PWD docker compose -f flaiover/compose.yaml up --build` runs the dashboard on <http://localhost:4242> next to `grafana/otel-lgtm` (collector, Tempo, Prometheus, Loki, Grafana on <http://localhost:3001>), with traces exported to it by default.

## Security posture

The dashboard has no authentication and can write to the mounted repository. Keep it bound to localhost. Do not put it behind a public reverse proxy. If a team needs shared access, run it on a host they reach over a private network or VPN, and treat the mount as you would a shared working copy.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
