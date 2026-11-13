---
title: "Runbook: run a story's agent in a container"
updated: 2026-11-13
status: active
---

# Run a story's agent in a container

Building the agent images and running one story's Claude Code agent in a container by hand, or a probe that tries an image without working the story. This is the first step of E-0029, flai running agents in Docker containers. Since S-0508 `flai serve` starts a story's agent this way on `harness: container`, as [docs/users/flai.md](../../users/flai.md#in-a-container-the-container-harness) describes; this runbook is the hand-run form. The scripts came from the experiment S-0505, and its results document records what was measured.

## The images

Two images, and a third when the project names a layer. This runbook runs the ones built on this machine from this repository's tree:

| Image | What it holds |
|-------|---------------|
| `flai-agent-bedrock:local` | Alpine, `git`, `ca-certificates`, `tini`, and `flai` built from the tree. The base of every agent image. Dockerfile: `flai/images/bedrock/Dockerfile`. |
| `flai-agent-claude-code:local` | Claude Code from npm, at the version the Dockerfile pins, with Node and ripgrep, on the bedrock image. Auto-update is off: a new Claude Code is a rebuild. Dockerfile: `flai/images/claude-code/Dockerfile`. |
| `flai-agent-<key>-claude-code:local` | The project's layer, on the claude-code image, when `system-flow.yaml` names its Dockerfile under `agent.image.layer`. This repository's is `flai-agent-sf-claude-code:local`, from `flai/images/system-flow/Dockerfile`: Go, golangci-lint, ShellCheck, and pnpm, the tools the close-out's tiers need ([docker.md](../../../design/tech/docker.md)). `go build` in it stamps VCS information, as on the host: for a story's worktree Go runs git in the main checkout, which git accepts only through the `safe.directory` the container is given, so a container started without it fails with `error obtaining VCS status: exit status 128` (I-0255). `scripts/agent-vcs-test.sh` checks it. |

These are local: nothing pushes them. They are not stamped and not signed, unlike flai's releases ([release-key.md](release-key.md)). None holds a login, a key, or a token: the agent's credentials reach the container when it runs. The first two hold no Go, golangci-lint, pnpm, or markdownlint, so an agent in them cannot run this repository's close-out tiers that need them; the layer image holds the tiers' tools.

### The published images

Each flai release, from the tag `flai/vX.Y.Z`, also publishes the first two images to GHCR, for linux/amd64 and linux/arm64, tagged with the version without its `v` ([ADR-0295](../../../design/adrs/0295-an-installed-flai-pulls-its-agent-images-by-the-digests-a-signed-list-on-its.md)):

| Image | What it holds |
|-------|---------------|
| `ghcr.io/bytepunx/flai-agent-bedrock:X.Y.Z` | The bedrock image, with the release's own `flai` and `pilot` for Linux in place of a tree build. |
| `ghcr.io/bytepunx/flai-agent-claude-code:X.Y.Z` | The claude-code image, built over that bedrock image by its digest. |

Both carry the labels `dev.system-flow.flai`, the release's version, and `dev.system-flow.commit`, its commit. Each is signed with the release key by its index digest. The GitHub release `flai vX.Y.Z` holds one digest list per image, `flai-agent-bedrock_X.Y.Z.digests` and `flai-agent-claude-code_X.Y.Z.digests`, each with its signature `.digests.sig`, in the format of [release-signing.md](../../../design/system/release-signing.md#the-digest-list), beside `cosign.pub`. Compare the key's fingerprint with [release-key.md](release-key.md#current-key) first.

Verify a list with cosign:

```bash
gh release download flai/vX.Y.Z --repo bytepunx/system-flow \
  --pattern 'flai-agent-*.digests*' --pattern cosign.pub
cosign verify-blob --key cosign.pub --signature flai-agent-bedrock_X.Y.Z.digests.sig \
  --insecure-ignore-tlog flai-agent-bedrock_X.Y.Z.digests
```

or with openssl:

```bash
openssl base64 -d -A -in flai-agent-bedrock_X.Y.Z.digests.sig -out flai-agent-bedrock_X.Y.Z.digests.der
openssl dgst -sha256 -verify cosign.pub -signature flai-agent-bedrock_X.Y.Z.digests.der flai-agent-bedrock_X.Y.Z.digests
```

Verify an image by the `index` digest its list names:

```bash
cosign verify --key cosign.pub --insecure-ignore-tlog=true ghcr.io/bytepunx/flai-agent-bedrock@sha256:<index digest>
```

The same commands with `claude-code` in place of `bedrock` check the other image.

GHCR makes a package created from a private repository private. A pull with no login needs both packages public, so after the first flai release that publishes them, set each one's visibility to public: on GitHub, Packages, `flai-agent-bedrock`, then Package settings, Change visibility; then the same for `flai-agent-claude-code`. Once a package is public it stays so for every later release.

`flai image ensure` pulls these images by the digests a verified list names (S-0740). The container harness still runs the image by the name the start uses, `flai-agent-claude-code:local`, and `flai image ensure` tags the image it pulls or builds with that name and with `flai-agent-claude-code:X.Y.Z`. The scripts below run the `:local` images too.

## Build

`flai serve` brings an agent image back by itself (S-0751). A start that finds a harness image or the project's layer missing or stale is held while `flai serve`'s image job makes it on the host. That holds for a story's agent, a strategic agent, the close-out's container, and the checks' container. One job runs at a time, and it does what `flai image ensure` does, below. Each job keeps its record and its log, `<id>.job.json` and `<id>.log`, in `~/.flai/serve/image-jobs/`. The host journal records its start and its end as action `image`, and a held run's `image_job` in `~/.flai/serve/agents.json` names it. When the job fails, each held start fails once, naming the cause, the failing step, and the log; the next start asks for a new job. The user guide has the rest ([In a container](../../users/flai.md#in-a-container-the-container-harness)).

Run `flai image ensure` yourself, with Docker running, when:

- no `flai serve` runs the start, as for this runbook's scripts or a start by hand, which refuses a missing image rather than make it;
- a job failed: mend what its log names, then run it to see the failure again or to make the images;
- you want the images before the first start needs them.

```bash
flai image ensure
```

It makes each harness image the project's agents need, as `flai image status` lists them, that is missing or stale. A released flai pulls it by the index digest its release's signed list names, and refuses a list or an image that does not verify, building nothing in its place. When the release has no list, the release or the registry cannot be reached and no list is kept, or `docker pull` fails, it builds the image on this host from the Dockerfiles flai carries and the release's Linux `flai` and `pilot`, checked against the release's signed `checksums.txt`. A development flai, such as the one `scripts/flai.sh` builds from this tree, is never pulled for: in this repository it builds from the tree's own Dockerfiles and source, and elsewhere it refuses. Then, when `system-flow.yaml` names `agent.image.layer`, it builds the project's layer over the harness image, in any project, when the layer is missing or stale or its harness image was just pulled or built (S-0745). It says which way each image took, then prints the status lines, and exits non-zero when a harness image or the layer is still missing or stale.

The layer's build runs from the project's root:

```bash
docker build -f <layer> --build-arg BASE=<harness image> --build-arg <name>=<value>… -t flai-agent-<key>-claude-code:local .
```

flai sets `BASE` to the harness image. Each `ARG` of the Dockerfile keeps its default unless `agent.image.args` in `system-flow.yaml` names it with a literal value ([settings](../settings.md#project-manifest)). The image is labelled with its recipe, over the Dockerfile and the arguments, and with the harness image's ID, so `flai image status` says it is stale after a change to the Dockerfile, to an argument, or to the harness image under it.

Each docker step's output goes to the terminal and to `.flai-cache/agent-images/<UTC stamp>.log` in the main checkout, which the first and last lines name. A step that writes nothing for 10 minutes is stopped and reported. A step that fails on the network, with a connection reset, a timeout, a DNS or TLS handshake failure, or a registry 5xx, is tried twice more, after 30 seconds and then 2 minutes. `flai image build` builds the images on this host without the pull, current or not, the layer last. Neither touches an image you named with `--image`. After a change to the layer alone, rebuild just it, over the harness image already on this host:

```bash
flai image build --layer
```

It refuses when the harness image is not on this host, naming it and `flai image ensure`, and when the manifest names no layer.

`make agent-images` is this repository's own build of all three images, from its tree. It is not how a missing image comes back: the image job or `flai image ensure` does that. From the repository root:

```bash
make agent-images
```

It runs `scripts/agent-image.sh`, which builds `flai-agent-bedrock:local` and then `flai-agent-claude-code:local` on top of it, and prints `built flai-agent-bedrock:local and flai-agent-claude-code:local`. When the manifest names a layer, it then builds the layer image over them and prints `built flai-agent-<key>-claude-code:local`; `scripts/agent-image.sh --layer` builds the layer alone. The script passes the layer `GO_VERSION` from `flai/go.mod` and the tool pins from `scripts/install-tools.sh`, where flai passes only `agent.image.args`, which this repository does not name yet; S-0760 moves them onto `agent.image.args`, and until then the script stays the way to build this repository's layer with those pins. `scripts/agent-image.sh --tag` prints the image a story's agent runs in and builds nothing. The flai inside is built from the checked-out tree, labelled with the last `flai/v*` tag and the commit. An unchanged tree rebuilds from Docker's cache. Build again after a change to flai that the agent should have, or to move Claude Code to another version.

Before it builds, the script asks the credential helper your Docker configuration names in `credsStore` (`$DOCKER_CONFIG/config.json`, else `~/.docker/config.json`) to list its logins, and waits ten seconds at most. Docker Desktop's helper, `desktop.exe`, cannot start from WSL while Windows interop is not answering, and a build would then stop at the first base image with `error getting credentials` ([I-0339](../../../design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md)). When the helper is not on `PATH`, fails, or does not answer, the script prints one line naming it and builds with no login, from a scratch Docker configuration that keeps your current context and links the other entries of your configuration's folder, such as `cli-plugins`, `contexts`, and `buildx`. The base images are public, so they pull anonymously. Your `config.json` is not changed, and the scratch configuration is removed when the script ends. A working helper, or none named, leaves the build as it was.

## What the run needs on the host

- Docker. S-0505 ran on Docker Desktop with WSL2 integration; a native Docker Engine host was not tried. There the MCP server must also listen on the bridge's gateway, since the engine routes `host.docker.internal` to it. `mcp.bind`'s default, `auto`, does so at each server start, and `flai serve` refuses its own container starts while the server does not ([ADR-0293](../../../design/adrs/0293-mcp-bind-holds-one-address-or-several-or-auto-the-default-which-flai-host.md)), but this script makes no such check.
- The images, made as above with `flai image ensure`, or `make agent-images` in this repository: the scripts run outside `flai serve`, so no image job makes them.
- The story started, so that its worktree exists: `flai story start S-nnnn`.
- The project's HTTP MCP server running: `flai host` runs one for every project `flai serve` serves, or start it with `flai mcp start` ([MCP over HTTP](../index.md#mcp-over-http)). `flai mcp status` says where it listens.
- `flai mcp token` answering the server's bearer token. The script reads it with the installed `flai` on `PATH`, not this tree's `bin/flai`; set `FLAI_BIN` to use another.
- A Claude Code login on the host, `~/.claude/.credentials.json`: sign in with `claude` on the host first.

## Run

```bash
scripts/agent-container.sh S-nnnn                   # the story's agent, working the live story
scripts/agent-container.sh --probe S-nnnn           # a trial of the image, which works no item
scripts/agent-container.sh [--probe] S-nnnn -- cmd  # cmd in place of pilot, same mounts
```

The script takes no prompt. pilot, the image's entrypoint, builds it from the start request the script feeds it, as `flai serve` does ([flai-cli.md](../../../design/system/flai-cli.md#pilot-the-agent-images-entrypoint-s-0507)).

| Form | Agent name | Session | `FLAI_STORY` | Prompt |
|------|------------|---------|--------------|--------|
| `S-nnnn` | `agent-S-nnnn` | `FLAI_AGENT_SESSION`, or a new one | `S-nnnn` | the story's agent's: it works the live story |
| `--probe S-nnnn` | `probe-S-nnnn` | always a new one | none | read and report: it calls flai's read tools, runs `flai version`, `git status`, and `git log -1`, reports one line per check, and ends |

A run without `--probe` is the story's agent. It moves the story's tasks, commits through `task_done`, and writes the narrative, as `flai serve` would. Use `--probe` to try an image: its agent works no item and never calls a writing tool, so the live story is left as it was ([I-0224](../../../design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md), S-0682). The script says on standard error that the run is a probe. The mounts are the story's in both forms, so a probe sees what the story's agent would.

A capped trial, on the smallest model and a budget of two cents:

```bash
FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-nnnn
```

With `--` the script runs the command in place of pilot, with `FLAI_AGENT` and `FLAI_SESSION` set as in the table, and `FLAI_STORY` only without `--probe`.

Without `--`, pilot runs `claude -p` headless in the agent image ([Settings](#settings)) with the arguments `flai serve` gives a story's agent, as the operator's uid, and the container is removed when it ends (`--rm`). The output is Claude Code's stream JSON on standard output.

The script and pilot write these to `.flai-cache/agents/S-nnnn/`, which is git-ignored:

- `mcp.json`, mode 600, which pilot writes. It names the HTTP MCP server, with its URL, token, and Host as variables that Claude Code expands as it reads the file, so the token is not written. Its `timeout` of 1800000 milliseconds keeps a slow `task_done`, `story_start`, or `test` from being dropped at 60 seconds ([ADR-0258](../../../design/adrs/0258-the-mcp-file-flai-writes-for-its-http-server-names-a-per-server-timeout-of.md)). An image built before that change writes it without the `timeout` and keeps the 60-second limit until `make agent-images` builds it again.
- `claude/.credentials.json`, a symbolic link to `/run/claude-login/.credentials.json`, the login inside the container. It does not resolve on the host, and nothing from the login is copied.
- `flai/`, an empty folder for the container's flai. On its first run that flai writes its config there, `flai/config.json`, and keeps its cache in `flai/cache`.

### What it mounts

| Host path | In the container | Mode | For |
|-----------|------------------|------|-----|
| `.flai-cache/worktrees/S-nnnn` | the same path, the working directory | read-write | the story's worktree |
| `.git` of the main checkout | the same path | read-only | git in the worktree reads the objects, refs, and `config` from it |
| `.git/worktrees/S-nnnn` of the main checkout, the story's registration | the same path, over the read-only `.git` | read-write | the worktree's own `index`, `HEAD`, logs, and locks, which its `.git` file names |
| `wip/` of the main checkout, as the manifest's `layout.wip` names it | the same path | read-only | the narrative, the board, the items, the threads, and the messages that `flai verify` checks ([ADR-0253](../../../design/adrs/0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md)) |
| `system-flow.yaml` of the main checkout | the same path | read-only | by it flai knows the main checkout is there, so `flai verify` runs rather than refusing |
| `.flai-cache/verify` of the main checkout, made first so it is yours | the same path | read-write | the report `flai verify` writes, which `flai verify --last` reads on the host |
| `.flai-cache/items-checked.lock` of the main checkout, made first as an empty 0600 file so Docker does not make it a folder | the same path | read-write | the checked-write lock, so `--record-issues` in the container waits for a checked write of a work item on the host ([ADR-0283](../../../design/adrs/0283-a-close-out-s-record-step-confirms-each-finding-outside-the-story-on-a-work.md)) |
| `.flai-cache/agents/S-nnnn` | the same path | read-write | in `claude/`, the `CLAUDE_CONFIG_DIR`, Claude Code's state: settings, transcript, caches; in `flai/`, the container's flai config and cache |
| `~/.claude` (`FLAI_CLAUDE_LOGIN`) | `/run/claude-login` | read-only | the login; only `.credentials.json` is read |

Nothing else is mounted. The Docker socket is never mounted, so the agent cannot start containers or reach the host's Docker. The two git mounts are [ADR-0246](../../../design/adrs/0246-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-git-read-only.md)'s: git in the worktree reads its history and refreshes its index, so `git status`, `diff`, `log`, and `show` work, and a `git commit`, `stash`, or new branch is refused by git on the read-only objects and refs. Commits are the host's, through `task_done` and the close-out. So are syncs: with `FLAI_STORY` set and the container's flai config in the story's `.flai-cache/agents/S-nnnn`, as here, `flai stream sync` and the shell's `flai task done` refuse before any git call, naming the MCP tools `task_done` and `verify`, which sync on the host. The script refuses to start when the story's registration, `.git/worktrees/S-nnnn`, is missing, and names its path.

#### Outside folders

The script mounts no folder outside the project. Under `flai serve`, the container harness mounts each folder you grant with `flai serve agent harness container --outside <dir>` ([ADR-0289](../../../design/adrs/0289-flai-serve-gives-a-story-s-agent-and-its-sub-agents-the-outside-folders-the.md), [the operators guide](../index.md)). It mounts each read-write at its own host path, after the agent folder and before the login. A story's agent's container and a strategic agent's container both get them. pilot passes each to the inner `claude` as `--add-dir`, so the agent and its sub-agents may read and write there. The setting itself refuses your home folder, `~/.flai` and any folder under or holding it, and any folder that holds the project or that the project holds. At the start, the container harness also refuses a folder that:

| Refused | Why |
|---------|-----|
| is not an absolute path | a mount needs one |
| holds a colon | a `-v` mount cannot carry it |
| is or holds the main checkout, your home, `~/.flai`, `/var/run/docker.sock`, or `/run/docker.sock` | a container never mounts them; `/` holds them all |
| lies under the main checkout or `~/.flai` | the checkout's mounts are the ones flai fixes, and `~/.flai` is never mounted |

A refused folder refuses the start, with an error that names the folder and says to grant another or clear them with `--outside ""`. A folder given twice is mounted once.

The worktree's `flaiover/node_modules` is inside the worktree mount, so the container and the host share it: a flaiover tier run in either installs it for both. It holds the native bindings of both libcs, glibc for the host and musl for the Alpine image, since `flaiover/pnpm-workspace.yaml` sets pnpm's `supportedArchitectures` ([I-0259](../../../design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md), S-0706). An install made before that setting holds only the libc it was made on. The next `scripts/flaiover-install.sh --if-needed` replaces it, since that script installs again whenever `pnpm-workspace.yaml` is newer than the install, and so does `make flaiover-install`.

#### A mount you type by hand in zsh

A `docker run` you type in zsh can mount the main checkout as an empty folder. zsh reads `$ROOT:r` as a modifier, even inside double quotes, so `-v "$ROOT:$ROOT:ro"` reaches docker as `<ROOT>:<ROOT>o`. Docker then mounts the checkout read-write at `<ROOT>o`, mounts nothing at `<ROOT>`, and `-w "$ROOT"` makes it an empty working directory. `docker inspect --format '{{json .HostConfig.Binds}}' <container>` shows the mangled target. Brace the name before the colon, `-v "${ROOT}:${ROOT}:ro"`, or use `--mount type=bind,src="${ROOT}",dst="${ROOT}",readonly`. This script and `flai serve` build docker's arguments in `sh` and in Go, so they are not affected ([I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md), S-0667).

### The close-out's and the checks' containers

Under `flai serve`, a story whose agent runs on `harness: container` is closed out, and checked, in containers of their own from the same image, so nothing from its worktree runs on the host ([ADR-0274](../../../design/adrs/0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)). This script does not start them; `flai serve` does. Their mounts beside the agent's:

| Mount | The agent's container | The close-out's and the checks' |
|-------|-----------------------|---------------------------------|
| the story's worktree, the working directory | read-write | read-write |
| `.git` of the main checkout | read-only | read-only |
| the story's registration, `.git/worktrees/S-nnnn` | read-write | read-write |
| `wip/` and `system-flow.yaml` of the main checkout | read-only | read-only |
| `.flai-cache/verify` of the main checkout | read-write | read-write |
| `.flai-cache/items-checked.lock` of the main checkout | read-write | read-write |
| `.flai-cache/agents/S-nnnn` | read-write | not mounted |
| each outside folder you grant, at its own path | read-write | not mounted |
| the login | read-only | not mounted |

The close-out's container is `flai-closeout-S-nnnn-<stamp>`, with `-r2` and `-r3` when main moved and it ran again, and runs `flai verify S-nnnn --record-issues`. Each check's is `flai-check-S-nnnn-<stamp>-<n>`. Neither gets the MCP server's URL or token, a provider's key, or `host.docker.internal`. The close-out's commit, its sync check, and its sync run on the host through git and `flai serve`'s own flai, with the main checkout's hooks only. `docker ps --filter name=flai-closeout-` and `--filter name=flai-check-` list them while they run.

The script sets two variables in the container so that its flai writes nothing under the container's home, `/home/agent`:

| Variable | Points at | For |
|----------|-----------|-----|
| `FLAI_CONFIG` | `.flai-cache/agents/S-nnnn/flai/config.json` | the config file the container's flai reads, and writes with defaults on its first run |
| `FLAI_CACHE_DIR` | `.flai-cache/agents/S-nnnn/flai/cache` | the `cache_dir` that first run writes into that config |

Both are host paths, the same inside the container, in the mounted agent folder. They do not touch the host's flai config: the script's own `flai mcp status` and `flai mcp token` run on the host, with the `FLAI_CONFIG` that `scripts/env.sh` sets there.

### Settings

Environment variables the script reads:

| Variable | Default | Effect |
|----------|---------|--------|
| `FLAI_BIN` | the `flai` on `PATH` | the flai that answers `flai mcp status` and `flai mcp token` |
| `FLAI_AGENT_IMAGE` | what `scripts/agent-image.sh --tag` prints: the layer image when the manifest names `agent.image.layer`, else `flai-agent-claude-code:local` | the image the run uses. The default is the one `flai serve`'s container harness runs (S-0670). The script refuses an image not built, naming `scripts/agent-image.sh`; `scripts/agent-container.sh --image` prints its choice and runs nothing |
| `FLAI_CLAUDE_LOGIN` | `~/.claude` | the folder holding the login, mounted read-only |
| `FLAI_AGENT_MODEL` | Claude Code's own | `claude --model`, such as `haiku` |
| `FLAI_AGENT_NETWORK` | `bridge` | `bridge` reaches the host as `host.docker.internal`; `host` uses Docker's host network |
| `FLAI_AGENT_MAX_BUDGET_USD` | none | the agent's `max_budget_usd`, such as `0.02`, for a trial run |
| `FLAI_AGENT_SESSION` | a new UUID | the session the agent logs under; a second run with the same one resumes it. Ignored with `--probe`, which always starts a new one |

On the bridge network the agent sends the server's own `Host` header, `127.0.0.1:<port>`, whatever address it dials, because the server refuses a request on its loopback listener that names another host. On Docker Desktop the `host` network is the VM's, not the machine's, and the server is not reachable from it.

## Stop

The agent stops when `claude -p` finishes. To stop it sooner, interrupt the script, or stop the container:

```bash
docker ps --filter ancestor=flai-agent-claude-code:local
docker stop <container id>
```

The container is removed when it stops. Whatever the agent left in the worktree stays there: read `git status` in the worktree before you go on with the story.

## Clean up

After the story, remove its cache folder. It holds the agent's transcripts, a probe's among them, and `mcp.json`:

```bash
rm -rf .flai-cache/agents/S-nnnn
```

Remove the images when you no longer run agents in containers:

```bash
docker image rm flai-agent-sf-claude-code:local flai-agent-claude-code:local flai-agent-bedrock:local
```

If the MCP token may have been copied elsewhere, rotate it with `flai mcp token --rotate`; every agent connected over HTTP then needs the new one.

## Limitations

S-0505's results document records what was tried, what worked, and what is still open, among them the network attempts, the login's OAuth refresh, which was not exercised, and git inside the container, which ADR-0246's two git mounts have since settled. The close-out runs inside the container past `flai verify`'s own steps, since the main checkout's `wip/`, `system-flow.yaml`, and `.flai-cache/verify` are mounted (ADR-0253, I-0218). Its tiers need the project's layer image, `flai-agent-sf-claude-code:local` here (S-0510), which this script runs since S-0670: in `flai-agent-claude-code:local` the close-out stops at the first tier that needs `go`. Its commit step stops inside the container whenever anything is uncommitted, such as an issue `--record-issues` recorded, since commits are the host's: commit through `task_done` and ask for the close-out from the host, whose job verifies in a container and commits on the host ([above](#the-close-outs-and-the-checks-containers)). Read S-0505's results before relying on this runbook for long runs: [S-0505's results](../../../design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md).
