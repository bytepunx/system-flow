---
title: Operators guide
updated: 2026-09-20
status: draft
---

# system-flow for operators

## Running the dashboard

`flai dashboard` runs `ghcr.io/bytepunx/flaiover` detached as `flaiover-<project>`, published on every interface at port `4242` by default (`--bind 127.0.0.1`, or `dashboard.bind`, restricts it to this host), and prints a login link; see Authentication below. The container is given that port and two secrets, each one file mounted read-only: the login token and the credential the host's flai proves itself with. **No file of the project is mounted into it** ([ADR-0031](../../design/adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). Everything the dashboard shows and everything it changes it asks of `flai serve` on this host, which `flai dashboard` registers the project with and starts; see the next section. `flai dashboard status`, `logs`, and `stop` manage it. Change the image, tag, port, or bind address in `~/.flai/config.json` or per project in `system-flow.yaml` under `dashboard`. Inside the monorepo `flai dashboard --build` builds the image from `flaiover/` as `flaiover:local` instead of pulling.

The container runs as you (`--user`), for one reason: the two secrets are files only you can read. There is nothing of yours within its reach to write.

Without `flai`, the container alone is:

```bash
docker run --detach --rm --name flaiover-myproject \
  --publish 0.0.0.0:4242:3000 \
  --mount type=bind,source="$PWD/.flai-cache/dashboard.token",target=/run/secrets/flaiover_token,readonly \
  --env FLAIOVER_TOKEN_FILE=/run/secrets/flaiover_token \
  --mount type=bind,source="$PWD/.flai-cache/dashboard.agent-key",target=/run/secrets/flaiover_agent_key,readonly \
  --env FLAIOVER_AGENT_KEY_FILE=/run/secrets/flaiover_agent_key \
  --user "$(id -u):$(id -g)" ghcr.io/bytepunx/flaiover:latest
```

It shows nothing until a `flai serve` that knows the project and holds the same credential connects to it, which is what `flai dashboard` arranges; started by hand, every page says that no flai is connected.

### Upgrading from a release that mounted the repository

Earlier releases mounted the repository read-write into the container, with `.git/hooks`, `.git/config`, and `.git/info` read-only over it, your git identity and global excludes passed in, and optionally an SSH key to push with. All of that is gone.

- **Restart the dashboard** with the new flai: `flai dashboard stop`, then `flai dashboard`. A container started by an older flai keeps its mounts until then, and `flai dashboard status` says so.
- **If you configured a push key** (`dashboard.push_key`, `--push-key`, `dashboard.push_known_hosts`): it is ignored, and `flai dashboard` says so at every start until you clear it with `flai config set dashboard.push_key ""` (and `dashboard.push_known_hosts` likewise). The container has no git or ssh to push with and no repository to push from. An acceptance from the board is committed and tagged on the host and pushed from the host with your own credentials: `flai push --pending`, below. If the key was a deploy key made for this purpose, delete it from the repository's settings and from `~/.ssh`; nothing uses it any more. `.flai-cache/dashboard.known_hosts` and `.flai-cache/dashboard.passwd` can be deleted.
- **Nothing to do for git settings.** Commits from the board are made on the host with your git configuration of the moment; the advice to restart the dashboard after changing git settings no longer applies.
- **A Windows or otherwise unusual host path** no longer matters: git runs where `flai serve` runs, so a story with a branch can be accepted from the board anywhere.

### Who acts when the dashboard writes

The dashboard has one project token and so one holder. What it does for them is recorded as the manifest's `owner` (`system-flow.yaml`), or `designer` when there is none: thread entries, moves made on the board, and acceptances, which run `flai accept <id> --by <owner>`. All of it is done by `flai serve` on the host, as you: commits carry your git identity and the dashboard's `Co-Authored-By` trailer.

### When an acceptance has not been pushed

A story accepted from the board is merged, committed, and tagged in your clone by `flai serve`, and waits there: nothing pushes it unasked, because whoever can push a release tag can publish a release. The board, the story's page, `flai board`, and the agents' MCP `inbox` all keep saying so until it is pushed. On the host, `flai push --pending` pushes the branch and the release tags of those acceptances with your own credentials; it never forces, and refuses when the remote has moved until you fetch and merge. Tags go three to a push, the branch last, because GitHub starts no tag-triggered workflow when one push carries more than three. An agent session that is running does this itself when `inbox` reports it. To make it unattended, run it from a timer of your own (a systemd user timer or cron entry calling `flai push --pending` in the repository); that is a push nobody approved, with your full credentials, and is your decision to make on your host.

### What the container can and cannot reach

It can reach its port, the network, and the two secrets. It cannot read or write any file of the project or of the host: no work tree, no `.git`, no `.flai-cache` beyond those two files, no flai configuration. So it cannot leave a git hook or setting that would run on your machine, change a tracked file or a branch, or plant an ignored file your tools execute, which were the routes open or guarded while the repository was mounted (I-0022, ADR-0027, superseded).

What a compromised container could still do is what the dashboard itself does: ask `flai serve` for the named methods it offers (reads of the three folders, moves, saves of Markdown under them, an acceptance of a story in review), each of which flai checks and performs itself, and present the login token. Treat a dashboard you expose beyond your own network accordingly, and see the next section for how to turn the connection off.

| Setting | Where | Default | Effect |
|---------|-------|---------|--------|
| `dashboard.notify_url` | `system-flow.yaml`, per project | unset | A webhook for the designer's inbox. When it is an `http` or `https` URL, the dashboard's server POSTs `{ "project": "<name>", "entry": { "key", "kind", "title", "href", "at" } }` as JSON for each inbox entry that appears after the server started: `kind` is `thread`, `question`, `review`, `blocked`, or `overlap`, and `href` is a path in the dashboard. One attempt per entry, a five second timeout, no retry; a failure is a warning in the log naming the host only, since the URL may carry a secret. The project token, file contents, and anything else are never sent. Entries that existed at start are not posted, so a restart does not replay the inbox. Read once at start: restart the dashboard after changing it |
| `dashboard.autocommit` | `system-flow.yaml`, per project | `true` | Documents saved from the dashboard's editor are committed on the main checkout by `flai serve`, one path per commit, authored by your git identity, with a `Co-Authored-By: flaiover` trailer. `false` leaves them uncommitted: agents on story branches then do not see the edit until someone commits it, and an acceptance from the board lists it as an uncommitted change. Nothing is pushed ([ADR-0023](../../design/adrs/0023-documents-are-saved-through-flai.md)) |
| `worktrees.relative_paths` | `~/.flai/config.json`, per user (`flai config set`) | `false` | With git 2.48 or newer, `flai stream open` links new worktrees with relative paths, for clones that are moved or shared. The dashboard no longer has any use for it. Sets `extensions.relativeWorktrees` on the clone, after which git older than 2.48 refuses the repository. Never enabled automatically. How to turn it back off is in [the flai guide](../users/flai.md) |

## The connection from flai on the host

`flai dashboard` starts a second thing beside the container: `flai serve`, one process per user on the host, which opens a WebSocket to the dashboard's `/agent` endpoint and keeps it open. Over it the dashboard asks flai for named things and flai answers. The direction is deliberate: the container is given no socket, pipe, or address of the host's, and cannot start a conversation with it. It is the only way the dashboard reaches the project: everything it reads and everything it changes goes over it, and the container holds no file of the project ([ADR-0031](../../design/adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). The image holds no `flai`, no `git`, and no `ssh`; commits made from the board are made by `flai serve` on the host as you, with your git configuration, and carry the dashboard's trailer.

- **The credential.** `flai dashboard` creates `.flai-cache/dashboard.agent-key` (mode 0600) and mounts it read-only at `/run/secrets/flaiover_agent_key`. It is not the login token and opens nothing but `/agent`. It never travels: `flai serve` sends a nonce, the dashboard answers with its own nonce and an HMAC of both under the credential, and `flai serve` answers with the mirror HMAC. A flai that dials a port something else is listening on gives nothing away and says so in `flai serve status`; a client that cannot prove the credential is dropped before it is served anything.
- **Who can reach `/agent`.** Anyone who can reach the dashboard's port can open the socket, and gets nothing without the credential. A request with an `Origin` header is refused outright, so no page in a browser can try. One connection per credential: a newer proven connection replaces the older.
- **What can be asked.** Only the methods `flai serve` offers (`flai hostapi` lists them), each checked against the project's key. There is no method that takes a command line or a file path of the caller's choosing, and the dashboard cannot add one. What a holder of the dashboard token can do is what the dashboard's pages can do, as before; what a compromised container can do is call those same methods and nothing more: it has no file of the project, no `git`, and no credential. Whatever it sends is validated as data (an item ID, a state, a Markdown path under the project's three folders, text) and reaches a command only as `--flag=value` or after `--`. It cannot choose who a change is recorded as.
- **When it is not there.** Pings every 4 seconds in both directions; a `flai serve` that is frozen or gone is marked so within ten seconds, a restarted container is reconnected within a second, and requests in flight fail with a clear error. Tried on Linux under WSL2 with Docker Engine. On macOS and Windows with Docker Desktop nothing differs in principle, since flai dials the published port on the host as a browser does, but it has **not been tried** there; on Windows `flai serve` runs in its own process group and `flai serve stop` ends it without a graceful signal.
- **Where its files are.** A folder named `serve` beside flai's config file: `projects.json` (what it serves), `state.json` (rewritten every second while it runs), `serve.log`. Nothing in it is secret; the credential stays in the project's `.flai-cache`.
- **Turning it off.** `flai dashboard --no-serve` starts the container without registering the project. The dashboard then has no flai to ask: the board, items, and threads answer 503, every page says why, and `/_ready` reports not ready (its `host_flai` check), which matters if you probe it. Use it when you run `flai serve` yourself in a terminal.
- **What flai serve reads.** Each registered project's design, docs, and wip folders and its manifest, looked at every 300 ms for changes (a stat of each file; a few hundred files cost well under a millisecond of CPU a look), and, when the dashboard asks, Markdown files under those three folders and nothing else: a path that climbs out, a link that leads out, another folder, or another file type is refused by flai, whatever the dashboard sends. The search index is built from the same files and kept in `flai serve`'s memory.
- **Stopping the container.** The image's server now closes every connection five seconds after SIGTERM (`SHUTDOWN_TIMEOUT`, in seconds), so a container asked to stop does stop, even with event streams open; before, a server that had stopped listening could linger while its container showed as up.

## MCP over HTTP

MCP is served by flai on the host, not by the dashboard ([ADR-0030](../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)). An agent on the host needs nothing from you: `.mcp.json` starts `flai mcp` on stdio. For an agent that cannot start a process there, run the same server over HTTP, one per project:

```bash
flai mcp start      # in the background; flai mcp http runs it in the foreground
flai mcp status     # where it listens, and what an agent's configuration looks like
flai mcp token      # its bearer token; --rotate replaces it and restarts the server
flai mcp stop
```

| | |
|-|-|
| Address | `http://127.0.0.1:4243/mcp` unless `--addr` says otherwise. The address is remembered per project (`.flai-cache/mcp-http.addr`), so an agent's configuration survives a restart; a second project on the same machine needs another port |
| Transport | MCP Streamable HTTP, answers as `application/json`. Clients on revisions up to 2025-11-25 get a session (`Mcp-Session-Id`), ended by `DELETE` or after `--idle` without a request (default 30 minutes), at most `--max-sessions` at once (default 16), after which `initialize` answers 503. Clients on 2026-07-28, which has no sessions, are served request by request at the same address |
| Authentication | `Authorization: Bearer <token>` only, from `.flai-cache/mcp.token` (mode 0600, git-ignored, created when first needed). It is not the dashboard's token, and the dashboard's token does not open it. A request with an `Origin` header, which is what a browser sends, is refused with 403 |
| The agent's name | The `X-Flai-Agent` header, else the client's own name, made safe for a file name. It is who thread entries and transitions are attributed to, and whose cursor `inbox` keeps |
| Long requests | `wait_for_events` holds its request open until something changes, for up to five minutes. A proxy or tunnel in front must allow an idle response that long, or agents see their wait cut short. At most 64 requests are in flight at once; more answer 503 |
| State and log | `.flai-cache/mcp-http.json` while it runs, `.flai-cache/mcp-http.log` for its events (`mcp server started`, `mcp session requested` with the agent and how many sessions are open, `mcp server stopped`) |
| What it serves | The main checkout, read when it starts: restart it after changing `system-flow.yaml` |
| Stopping | `flai mcp stop` ends held `wait_for_events` calls as waits that ran out, so a waiting agent gets an answer, not a broken connection, and simply starts a session again when the server is back |

It listens on this machine only by default. `--addr 0.0.0.0:4243` or another interface is allowed and logged as a warning: the token travels in every request and flai does not encrypt it, so beyond the machine put an SSH forward, a tunnel, or a proxy that terminates TLS in front and give agents that address.

**The dashboard's `/mcp` is gone.** Until flaiover 0.22 the dashboard served MCP itself (ADR-0024, superseded). Any request to `/mcp` now answers 410 with a message naming `flai mcp start`, so an agent configured for the old address is told what to do. `FLAIOVER_MCP_MAX_SESSIONS` and `FLAIOVER_MCP_IDLE_MINUTES` no longer mean anything.

### The tunnel expectation

The dashboard speaks plain HTTP and the token travels in every request. On the machine itself or a network you trust that is acceptable. Anywhere else, put a tunnel or a reverse proxy that terminates TLS in front of it and use the `https` address; never publish the dashboard's port to the internet as it is. To keep it to the host and let only the tunnel reach it, bind it to loopback: `dashboard.bind: 127.0.0.1` in `system-flow.yaml`, or `flai dashboard --bind 127.0.0.1`. The same holds for `flai mcp` over HTTP, which has its own port and token.

### Project identity

Every `/api/*` response of the dashboard, and every answer of `flai mcp` over HTTP, names the project it came from: the headers `X-Flai-Project-Key` and `X-Flai-Project-Name` (URI-encoded), and `project: { name, key }` in JSON object bodies. Both come from `name` and `key` in `system-flow.yaml`; `flai check` warns when `key` is missing. The dashboard's refusals for lack of a token carry neither; flai's MCP server sends the headers on refusals too, since its caller chose the project by choosing the address.

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

`ghcr.io/bytepunx/flaiover` is built from `flaiover/Dockerfile` at the repository root: the SvelteKit build on `node:24-alpine` and nothing else (no `flai`, no `git`, no `ssh`; the label `dev.system-flow.flai` names the flai of the commit it was built and tested with), listening on `3000`, running as an unprivileged user by default and working under any `--user`. Tags: `latest` (main), `X.Y.Z` and `X` from `flaiover/vX.Y.Z` release tags, and `sha-<commit>`. While the repository is private the package is too: `docker login ghcr.io` with a token that has `read:packages` before `flai dashboard` can pull it. The container writes nothing that outlives it: it has no volume, and its only home is `/tmp`.

Build locally with `make flaiover-image` (tag `flaiover:local`) and run it with `flai dashboard --image flaiover --tag local`.

## Telemetry

flaiover follows the logging and telemetry conventions.

| Signal | Where | Notes |
|--------|-------|-------|
| Logs | stdout, one JSON event per line | `ts`, `level`, `service`, `component`, `msg`, and fields; one line per request with `trace_id` (or `request_id` when tracing is off), `method`, `route`, `path`, `status`, `duration_ms`; `LOG_LEVEL` (`debug`, `info`, `warn`, `error`), `LOG_FORMAT=text` for key-value text outside production |
| Liveness | `GET /_health` | Always 200 while the process runs; touches nothing |
| Readiness | `GET /_ready` | Checks that flai on the host is connected (`host_flai`) and, through it, the project's manifest and item listing, with a two second timeout each; 503 names the failing check |
| Metrics | `GET /metrics` | Prometheus format: `flaiover_http_requests_total{method,route,status}`, `flaiover_http_request_duration_seconds` (histogram), `flaiover_http_requests_in_flight`, `flaiover_build_info{version,commit}` (the flaiover release tag and commit baked into the image), plus Node process metrics prefixed `flaiover_`. Route labels are SvelteKit route ids, never paths with IDs |
| Traces | OTLP/HTTP | Exported only when `OTEL_EXPORTER_OTLP_ENDPOINT` is set; one server span per request named `<method> <path>` with W3C context taken from the incoming headers; standard `OTEL_*` variables apply (`OTEL_SERVICE_NAME`, `OTEL_TRACES_SAMPLER`, `OTEL_TRACES_SAMPLER_ARG`) |

Local stack: `PROJECT=$PWD docker compose -f flaiover/compose.yaml up --build` runs the dashboard on <http://localhost:4242> next to `grafana/otel-lgtm` (collector, Tempo, Prometheus, Loki, Grafana on <http://localhost:3001>), with traces exported to it by default.

## Security posture

The dashboard authenticates every request with the project token (above), and through `flai serve` it can change the project: moves, saves, acceptances. By default it is published on every interface of the host so a team can reach it over a private network or VPN; beyond a trusted LAN put a TLS-terminating tunnel or proxy in front of it, because the token travels in clear over plain HTTP. To keep it to the machine it runs on, set `dashboard.bind: 127.0.0.1` in `system-flow.yaml` or config, or pass `--bind 127.0.0.1`. The container holds no file of the project and no git credential, so the token is not the power to publish: an acceptance waits on the host until you push it. Keep the dashboard off public addresses or behind a tunnel you trust all the same, and rotate the token (`flai dashboard token --rotate`) when in doubt.

## Requirements

- Docker Engine 24 or newer on `PATH`
- `git` on `PATH` for template cloning
