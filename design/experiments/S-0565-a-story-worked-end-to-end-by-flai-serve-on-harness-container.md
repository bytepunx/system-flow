---
title: A story worked end to end by flai serve on harness container
updated: 2026-11-07
status: active
story: S-0565
---

# A story worked end to end by flai serve on harness container

S-0505 ran a story's agent in a container started by hand. This experiment runs it where it is meant to run: `flai serve` on this host, the installed release, starts a story of this repository on `harness: container`, and its agent works the story inside the container.

## Hypothesis

The `container` harness of S-0508, with S-0510's image layer and the mounts of S-0517 (ADR-0246) and S-0520 (ADR-0253), runs a story's agent under `flai serve` from `story_start` to review as a host run does. Its usage is measured and its card shows its status as a host run's does. A stop from the dashboard stops the container, and an answered thread starts it again in the same session. A story that changes Go closes out inside the layered container with every tier the manifest's `tests` select.

## Success measure

S-0565's three acceptance criteria, each met on this host:

1. One story of this repository is worked end to end by `flai serve` with `harness: container`, from `story_start` to review, its usage measured and its card's status shown as for a host run. Each limitation met is an issue under `design/issues`.
2. A stop from the dashboard and a restart after an answered thread are each tried once on that story's container, and what happened is recorded in S-0565's Notes.
3. One story that changes Go is worked end to end by `flai serve` on `harness: container` with the project's image layer (`flai/images/system-flow/Dockerfile`), and `scripts/close-out.sh` runs every tier the manifest's `tests` select and passes inside the container. Each limitation met is an issue under `design/issues`.

## What was done

### Proof stories

TH-0821 settled the two proof stories. Both carry `harness: container` on their agent since 2026-11-07, set with `flai edit S-nnnn --harness container`.

| Story | Changes | For criteria |
|-------|---------|--------------|
| S-0558 | `safety.md` against `git.md` on force pushes, documentation only | 1 and 2 |
| S-0556 | what `flai guard` lists among a sub-agent's reads, Go | 3 |

Under `flai serve` the close-out is a job on the host ([ADR-0229](../adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)), so TH-0821 settled option (a): S-0556's close-out also runs once by hand inside the layered image, with `scripts/agent-container.sh S-0556 -- scripts/close-out.sh S-0556`.

### Host preconditions

Checked on 2026-11-07 at about 03:43Z. The host runs Docker Desktop 29.7.2 on WSL2.

| Precondition | Command | Answer | Holds |
|--------------|---------|--------|-------|
| The installed flai holds S-0508, S-0510, S-0517, and S-0520 | `flai version`; `flai serve status`; `git tag --contains` on each story's acceptance commit | flai 1.83.0 is installed, and `flai serve` 1.83.0 has run since 2026-11-03T18:12:07Z. S-0508 and S-0517 are first in 1.81.0, S-0510 in 1.83.0, and S-0520 only in 1.84.0, published but not installed | No: S-0520's mounts (ADR-0253) are missing |
| The images are built | `docker image inspect --format '{{.RepoTags}} {{.Created}}'` | `flai-agent-bedrock:local` 2026-11-02T18:51:24Z, `flai-agent-claude-code:local` 2026-11-02T18:51:49Z, `flai-agent-sf-claude-code:local` 2026-11-02T18:52:42Z | No: the layer image predates S-0510's layer commits of 2026-11-03T01:19Z (`d89ce9278`, `7307354e7`), and every image predates 1.84.0 |
| The container harness is set on the host | `flai serve agent harness` | program `docker`; image `flai-agent-claude-code:local` by default; network Docker's default bridge; login not set | No: `--login` is not set |
| A container reaches the host's HTTP MCP listener (I-0254) | `wget -S http://host.docker.internal:4243/mcp` from `flai-agent-claude-code:local` on the default bridge | `HTTP/1.1 401 Unauthorized` from 192.168.65.254, though the listener binds 127.0.0.1:4243 | Yes: Docker Desktop's host gateway reaches the loopback bind, and the 401 is the bearer-token gate |
| The proof stories can be started | `flai show S-0558`, `flai show S-0556` | both in `backlog` | No: moving them to ready is the operator's or the orchestrator's |

Each precondition that did not hold was the operator's: `flai guard` refuses an agent `flai serve` started every write to the host configuration ([ADR-0220](../adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md)), and upgrading the host restarts `flai serve`. TH-0836 asked for them. The operator confirmed the orchestrator's recommendation there at 2026-11-07T04:19:06Z: the operator upgrades the host, rebuilds the images, and sets the login, and S-0558 goes to ready only after that, moved by the operator or the orchestrator.

Checked again at 2026-11-07T04:19Z:

| Precondition | Command | Answer | Holds |
|--------------|---------|--------|-------|
| The installed flai holds all four stories | `flai version`; `flai serve status` | flai 1.84.0; `flai serve` 1.84.0 since 2026-11-07T03:56:31Z | Yes |
| The images are built | `docker image inspect`; `flai version` and `claude --version` inside `flai-agent-sf-claude-code:local` | bedrock 03:46:02Z, claude-code 03:46:25Z, the layer 03:47:29Z; flai 1.84.0 (`e41d6124e`) and Claude Code 2.1.312 inside | Yes |
| The container harness is set | `flai serve agent harness` | login `/home/alex/.claude-me`, mounted read-only; image and network as before | Yes |
| The proof stories can be started | `flai show S-0558` | S-0558 in `backlog`; the in-progress limit is full, 5 of 5 | Waits for the move to ready and for room |

## Results

Not yet run.

## Recommendation

Pending: whether to adopt, adapt, or drop running this repository's stories on `harness: container` is decided once both proof runs are recorded (T-2453).
