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

### Who acts when the dashboard writes

The dashboard has one project token and so one holder. What it does for them is recorded as the manifest's `owner` (`system-flow.yaml`), or `designer` when there is none: thread entries, moves made on the board, and acceptances, which run `flai accept <id> --by <owner>`. Git commits made in the container are authored by the git identity `flai dashboard` passes in from the host. An acceptance from the dashboard needs that identity, the repository mounted at its host path, and the host's global git excludes, all of which `flai dashboard` sets up; it is committed and tagged locally and never pushed, because the container holds no credentials.

### Pushing what the board accepts

By default the container holds no git credential. A story accepted from the board is merged, committed, and tagged in your clone, the dashboard says it was accepted locally, and someone pushes from a shell. That is deliberate: whoever can push a release tag can publish a release.

If you want an acceptance from the board to reach the remote at once, release tags included, name an SSH private key on the host and `flai dashboard` mounts it into the container read-only ([ADR-0026](../../design/adrs/0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md)). It is a setting of this machine, never of the repository:

```bash
flai config set dashboard.push_key ~/.ssh/flaiover-push   # or: flai dashboard --push-key <path>
flai dashboard stop && flai dashboard
flai dashboard status                                     # says which key the running container holds
flai config set dashboard.push_key ""                     # back to holding nothing, after a restart
```

Understand what this changes before turning it on. The dashboard token becomes the power to publish: whoever holds it can release any story that is in review by accepting it. And if the container were ever compromised, the key is a file it can read.

**Which key.** In order of how little a stolen copy opens:

1. **A key made for this repository alone** (recommended). Make one without a passphrase and add its public half to the repository as a deploy key with write access (on GitHub: Settings, Deploy keys, Add deploy key, Allow write access). A stolen copy opens this one repository and nothing else.

   ```bash
   ssh-keygen -t ed25519 -N "" -C "flaiover push $(hostname)" -f ~/.ssh/flaiover-push
   cat ~/.ssh/flaiover-push.pub      # add this as the deploy key
   ```

   If the same host already uses another key for that provider, the remote URL stays as it is: the container is told to use only the key you named.
2. **Your own key.** It works the same way and is the quickest to set up, but it opens every repository of your account and every server that trusts it, and the container can read the file. If it has a passphrase, as your own key should, `flai dashboard` refuses it: nobody is there to type the passphrase when an acceptance is pushed.

Your SSH agent is never forwarded into the container. It would lend the container every key in the agent for as long as it runs, and it stops working when you log out.

**What `flai dashboard` checks before it starts anything**, each with the reason when it refuses: the file exists, is no looser than mode 0600, is a private key, and has no passphrase; the `origin` remote is an SSH URL; and this machine already has a `known_hosts` entry for the remote's host. The container never accepts a host key on first use: connect once from a shell (`ssh -T git@github.com`), check the fingerprint against the ones your provider publishes, and accept it there. `--push-known-hosts <file>` (or `dashboard.push_known_hosts`) takes the host keys from a file you curate instead of your `known_hosts` and the system's. The pinned host keys and a passwd entry for your user ID, which OpenSSH needs, are written under `.flai-cache/` and mounted read-only. At start `flai dashboard` prints the key's fingerprint and comment, never the key.

**Limit what a stolen key can do** with rules on the repository: block force pushes to the default branch, and restrict deletion of the default branch and of release tags (`flai/v*`, `flaiover/v*`). No rule stops a key that may push from publishing a release; that is what the key is for.

**To revoke**: delete the deploy key from the repository (or remove your own key from your account), unset `dashboard.push_key`, and restart the dashboard. A dedicated key never expires on its own.

A push that fails (no network, a revoked key) leaves the acceptance standing, as without a key: the dashboard says it was accepted locally and shows the command. Not supported on a Windows host yet; `flai dashboard` says so and refuses the key.

### When an acceptance has not been pushed

The dashboard's container holds no git credential unless you give it a push key (above), so an acceptance from the board is committed and tagged in your clone and waits there. The board, the story's page, `flai board`, and the agents' MCP `inbox` all keep saying so until it is pushed. On the host, `flai push --pending` pushes the branch and the release tags of those acceptances with your own credentials; it never forces, and refuses when the remote has moved until you fetch and merge. An agent session that is running does this itself when `inbox` reports it. To make it unattended without giving the container anything, run it from a timer of your own (a systemd user timer or cron entry calling `flai push --pending` in the repository); that is a push nobody approved, with your full credentials, and is your decision to make on your host.

### What the container can and cannot write

The repository is mounted read-write, `.git` included, because accepting a story from the board merges, commits, tags, and removes a worktree. A few paths inside it are mounted a second time, read-only, so that a compromised container cannot leave something that git would later run on your machine as you ([ADR-0027](../../design/adrs/0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md)):

- `.git/hooks`: no hook can be added or changed from the container.
- `.git/config`: no git setting that runs a command or redirects a remote (`core.sshCommand`, `core.hooksPath`, a shell alias, a filter, `url.*.insteadOf`, a `pushurl`, and the like).
- `.git/info`: `exclude` cannot be used to hide a planted file from `git status`, and `attributes` cannot name a filter.
- `.flai-cache/dashboard.token`, and your flai config when it lives inside the repository (it names the image and the push key the next `flai dashboard` uses).

`flai dashboard` says so when it starts, and also lists what the clone already holds of those kinds (an enabled hook, such a setting) so you can check that they are yours; `flai dashboard status` repeats it, and tells you when the running container was started by an older flai without these mounts: restart it.

Two things follow. The container reads the git config as it was when the dashboard started, so after changing git settings on the host (a new remote, say), restart the dashboard. And this closes the routes git itself offers, not every route a writable working copy offers: the container can still change a tracked file, which you would see in `git status` and `git diff` and which acceptance refuses to include unasked, and it can write files that git ignores. If your project runs ignored files on the host (a built binary, tool caches), treat a dashboard you expose beyond your own network accordingly.

### Why the mount path matters

Git links a story worktree (`.flai-cache/worktrees/S-nnnn`) to the repository with absolute paths in both directions. The container runs git for acceptance, so those paths must exist inside it, which they do when the repository is mounted at its host path ([ADR-0022](../../design/adrs/0022-repository-mounted-at-its-host-path.md)). A container that sees the repository anywhere else, including one started by an older flai at `/project`, cannot accept a story that has a branch; the confirmation says so and points at `flai accept`. The image's default is still `PROJECT_DIR=/project` for mounts made by hand, which is fine for reading and for projects without story branches.

| Setting | Where | Default | Effect |
|---------|-------|---------|--------|
| `GIT_CONFIG_COUNT`, `GIT_CONFIG_KEY_0`, `GIT_CONFIG_VALUE_0` | container environment, set by `flai dashboard` | `core.excludesFile` pointing at `/run/flaiover/gitignore` when the host has a global excludes file, unset otherwise | The host's global git excludes (`core.excludesFile`, else `$XDG_CONFIG_HOME/git/ignore`, else `~/.config/git/ignore`) is bind-mounted read-only there, so a file ignored only on the host is not reported as uncommitted in the dashboard and does not stop an acceptance (S-0051) |
| `dashboard.notify_url` | `system-flow.yaml`, per project | unset | A webhook for the designer's inbox. When it is an `http` or `https` URL, the dashboard's server POSTs `{ "project": "<name>", "entry": { "key", "kind", "title", "href", "at" } }` as JSON for each inbox entry that appears after the server started: `kind` is `thread`, `question`, `review`, `blocked`, or `overlap`, and `href` is a path in the dashboard. One attempt per entry, a five second timeout, no retry; a failure is a warning in the log naming the host only, since the URL may carry a secret. The project token, file contents, and anything else are never sent. Entries that existed at start are not posted, so a restart does not replay the inbox. Read once at start: restart the dashboard after changing it |
| `dashboard.autocommit` | `system-flow.yaml`, per project | `true` | Documents saved from the dashboard's editor are committed on the main checkout, one path per commit, authored by the git identity `flai dashboard` passes in, with a `Co-Authored-By: flaiover` trailer. `false` leaves them uncommitted: agents on story branches then do not see the edit until someone commits it, and an acceptance from the board lists it as an uncommitted change. Commits are never pushed from the container ([ADR-0023](../../design/adrs/0023-documents-are-saved-through-flai.md)) |
| `PROJECT_DIR` | container environment, set by `flai dashboard` | the repository's host path (`/project` in the image) | The repository flaiover serves and the folder flai's config and cache are read from (`.flai-cache`) |
| `worktrees.relative_paths` | `~/.flai/config.json`, per user (`flai config set`) | `false` | With git 2.48 or newer, `flai stream open` links new worktrees with relative paths so they work at any mount path. Sets `extensions.relativeWorktrees` on the clone, after which git older than 2.48 refuses the repository. Never enabled automatically. How to turn it back off is in [the flai guide](../users/flai.md) |

When the host path cannot be used in a Linux container (a Windows drive path, or a path containing a colon), `flai dashboard` mounts at `/project`, logs a warning, and stories with a branch are accepted from a shell unless worktrees are relative.

## MCP over HTTP

The dashboard serves the project's MCP server at `/mcp` ([ADR-0024](../../design/adrs/0024-mcp-over-http-and-project-identity.md)), so an agent on another machine, and later a hub, reach the same tools an agent on the host has through `flai mcp`.

| | |
|-|-|
| Transport | MCP Streamable HTTP: `POST /mcp` for messages, `DELETE /mcp` to end a session, `GET /mcp` answers 405 (no server-initiated stream is offered) |
| Authentication | `Authorization: Bearer <project token>` only. The browser session cookie does not open `/mcp`, and a request whose `Origin` is another site is refused |
| Sessions | Each session is its own `flai mcp` process on the host, running as the agent named by the `X-Flai-Agent` header on `initialize`, else the client's name. A session ends on `DELETE`, or after `FLAIOVER_MCP_IDLE_MINUTES` without a request (default 30); at most `FLAIOVER_MCP_MAX_SESSIONS` exist at once (default 16), after which `initialize` answers 503 |
| Long requests | `wait_for_events` holds its request open until something changes, for up to five minutes. A proxy or tunnel in front of the dashboard must allow an idle response that long, or agents see their wait cut short |
| Log | `mcp session started` and `mcp session ended` events, with the session, the agent, and how many sessions are open |

### The tunnel expectation

The dashboard speaks plain HTTP and the token travels in every request. On the machine itself or a network you trust that is acceptable. Anywhere else, put a tunnel or a reverse proxy that terminates TLS in front of it and give agents the `https` address; never publish the dashboard's port to the internet as it is. To keep it to the host and let only the tunnel reach it, bind it to loopback: `dashboard.bind: 127.0.0.1` in `system-flow.yaml`, or `flai dashboard --bind 127.0.0.1`. A hub, when there is one, is reached the other way round: the dashboard dials out to it, so no inbound port is opened at all.

### Project identity

Every `/api/*` and `/mcp` response names the project it came from: the headers `X-Flai-Project-Key` and `X-Flai-Project-Name` (URI-encoded), and `project: { name, key }` in JSON object bodies. Both come from `name` and `key` in `system-flow.yaml`; `flai check` warns when `key` is missing. Responses that refuse a request for lack of a token carry neither.

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

The dashboard authenticates every request with the project token (above) and can write to the mounted repository. By default it is published on every interface of the host so a team can reach it over a private network or VPN; beyond a trusted LAN put a TLS-terminating tunnel or proxy in front of it, because the token travels in clear over plain HTTP. To keep it to the machine it runs on, set `dashboard.bind: 127.0.0.1` in `system-flow.yaml` or config, or pass `--bind 127.0.0.1`. Treat the mount as you would a shared working copy; the git hooks, config, and info inside it are read-only to the container (above). The container holds no git credential unless you give it a push key (above); with one, the token is also the power to publish a release, so keep the dashboard off public addresses or behind a tunnel you trust, and rotate the token (`flai dashboard token --rotate`) when in doubt.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
