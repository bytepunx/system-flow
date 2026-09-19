---
title: Operators guide
updated: 2026-09-19
status: draft
---

# system-flow for operators

## Running the dashboard

`flai dashboard` runs `ghcr.io/bytepunx/flaiover` detached as `flaiover-<project>` with the repository mounted read-write at the same absolute path it has on the host and `PROJECT_DIR` set to that path, published on every interface at port `4242` by default (`--bind 127.0.0.1`, or `dashboard.bind`, restricts it to this host), as the invoking user. It prints a login link; see Authentication below. It also passes your git `user.name` and `user.email` into the container, so an acceptance made from the dashboard commits as you, and your global git excludes file, so git there ignores what it ignores for you; the container has no credentials, so such an acceptance is committed and tagged locally and you push it from a shell. `flai dashboard status`, `logs`, and `stop` manage it. Change the image, tag, port, or bind address in `~/.flai/config.json` or per project in `system-flow.yaml` under `dashboard`. Inside the monorepo `flai dashboard --build` builds the image from `flaiover/` as `flaiover:local` instead of pulling.

Without `flai`, the equivalent is:

```bash
docker run --detach --rm --name flaiover-myproject \
  --publish 0.0.0.0:4242:3000 --volume "$PWD:$PWD" --env PROJECT_DIR="$PWD" \
  --mount type=bind,source="$PWD/.flai-cache/dashboard.token",target=/run/secrets/flaiover_token,readonly \
  --env FLAIOVER_TOKEN_FILE=/run/secrets/flaiover_token \
  --user "$(id -u):$(id -g)" ghcr.io/bytepunx/flaiover:latest
```

### Why the mount path matters

Git links a story worktree (`.flai-cache/worktrees/S-nnnn`) to the repository with absolute paths in both directions. The container runs git for acceptance, so those paths must exist inside it, which they do when the repository is mounted at its host path ([ADR-0022](../../design/adrs/0022-repository-mounted-at-its-host-path.md)). A container that sees the repository anywhere else, including one started by an older flai at `/project`, cannot accept a story that has a branch; the confirmation says so and points at `flai accept`. The image's default is still `PROJECT_DIR=/project` for mounts made by hand, which is fine for reading and for projects without story branches.

| Setting | Where | Default | Effect |
|---------|-------|---------|--------|
| `GIT_CONFIG_COUNT`, `GIT_CONFIG_KEY_0`, `GIT_CONFIG_VALUE_0` | container environment, set by `flai dashboard` | `core.excludesFile` pointing at `/run/flaiover/gitignore` when the host has a global excludes file, unset otherwise | The host's global git excludes (`core.excludesFile`, else `$XDG_CONFIG_HOME/git/ignore`, else `~/.config/git/ignore`) is bind-mounted read-only there, so a file ignored only on the host is not reported as uncommitted in the dashboard and does not stop an acceptance (S-0051) |
| `dashboard.autocommit` | `system-flow.yaml`, per project | `true` | Documents saved from the dashboard's editor are committed on the main checkout, one path per commit, authored by the git identity `flai dashboard` passes in, with a `Co-Authored-By: flaiover` trailer. `false` leaves them uncommitted: agents on story branches then do not see the edit until someone commits it, and an acceptance from the board lists it as an uncommitted change. Commits are never pushed from the container ([ADR-0023](../../design/adrs/0023-documents-are-saved-through-flai.md)) |
| `PROJECT_DIR` | container environment, set by `flai dashboard` | the repository's host path (`/project` in the image) | The repository flaiover serves and the folder flai's config and cache are read from (`.flai-cache`) |
| `worktrees.relative_paths` | `~/.flai/config.json`, per user (`flai config set`) | `false` | With git 2.48 or newer, `flai stream open` links new worktrees with relative paths so they work at any mount path. Sets `extensions.relativeWorktrees` on the clone, after which git older than 2.48 refuses the repository. Never enabled automatically. How to turn it back off is in [the flai guide](../users/flai.md) |

When the host path cannot be used in a Linux container (a Windows drive path, or a path containing a colon), `flai dashboard` mounts at `/project`, logs a warning, and stories with a branch are accepted from a shell unless worktrees are relative.

## Authentication

Every request except `/_health` and `/_ready` needs the project's token (ADR-0018).

| Step | How |
|------|-----|
| Set | `flai dashboard` generates it on first run: 32 random bytes, base64url. `flai dashboard token` prints it; `--rotate` replaces it and restarts a running dashboard, ending every session |
| Stored | `.flai-cache/dashboard.token`, mode 0600, one per project. `.flai-cache/` must be git-ignored; `flai check` refuses a repository where it is not |
| Handed to the container | Bind-mounted read-only at `/run/secrets/flaiover_token` with `FLAIOVER_TOKEN_FILE` pointing at it. Never an environment variable, so `docker inspect` does not show it |
| Browser | Open the login link (`http://host:4242/login#token=...`). The fragment never leaves the browser; the page exchanges it for an HttpOnly, SameSite=Lax cookie (Secure over HTTPS) and rewrites history. Pasting the token on `/login` also works |
| Agents and tools | `Authorization: Bearer <token>` on every request |
| Metrics | `/metrics` needs the token unless `FLAIOVER_METRICS_PUBLIC=true` |
| Development | `scripts/flaiover-dev.sh` points the dev server at the same file; `FLAIOVER_AUTH=off` disables authentication outside production only |

Exposure, plainly:

| Risk | Answer |
|------|--------|
| Sniffed on the network | flaiover speaks plain HTTP. Beyond a trusted LAN, terminate TLS in a tunnel or proxy; the cookie is then marked Secure |
| Token in logs | Request logs carry method, route, status, and duration only; never header or cookie values |
| Token in the image or repository | Never baked in, never committed; verified by `flai check` |
| Guessing | 256 bits of randomness |
| Stolen cookie | HttpOnly blocks script access; rotation invalidates every session |
| Docker socket access on the host | Can read the mounted file; that is host ownership, not something the token scheme addresses |

## Access to the image

`ghcr.io/bytepunx/flaiover` is a private package while the repository is private, so a plain `docker pull` is refused. Three ways in:

| Path | What to do |
|------|------------|
| Let `flai` log in | Have `GITHUB_TOKEN` or `GH_TOKEN` set, or be logged in with `gh`, and give the token the `read:packages` scope: `gh auth refresh -h github.com -s read:packages`. On a refused pull `flai dashboard` runs `docker login ghcr.io` with that token and retries once. The token never appears in arguments or logs |
| Log in yourself | `gh auth token \| docker login ghcr.io --username <github-user> --password-stdin`, then `flai dashboard` or the `docker run` above |
| Build locally | In the monorepo, `flai dashboard --build` (or `make flaiover-image` then `flai dashboard --image flaiover --tag local`) builds `flaiover:local` from `flaiover/Dockerfile` and runs it; no registry access needed |

Making the package public on GHCR removes the need for a token entirely; that is an organisation setting, not something `flai` changes.

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

The dashboard authenticates every request with the project token (above) and can write to the mounted repository. By default it is published on every interface of the host so a team can reach it over a private network or VPN; beyond a trusted LAN put a TLS-terminating tunnel or proxy in front of it, because the token travels in clear over plain HTTP. To keep it to the machine it runs on, set `dashboard.bind: 127.0.0.1` in `system-flow.yaml` or config, or pass `--bind 127.0.0.1`. Treat the mount as you would a shared working copy.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
