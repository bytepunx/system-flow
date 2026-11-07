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

### The first run: S-0558, documentation only

S-0558 was worked from `story_start` to review by its agent in its container, under `flai serve` 1.84.0, in about six minutes.

| When (UTC) | What | Seen in |
|------------|------|---------|
| 04:21:54 | The orchestrator moves S-0558 to ready; the in-progress limit is full, 5 of 5 | `moved` event |
| 04:36:23 | `flai serve` starts `docker` (container) for S-0558 as `agent-S-0558`, once S-0555 moves to review | `agent` event `started`, why "it entered ready"; `flai serve journal` |
| 04:36:27 | Container `flai-agent-S-0558-47f3eec4-8dd1-41ac-adfb-cbc7c6aff228` is up, image `flai-agent-sf-claude-code:local` | `docker ps` |
| 04:36:31 | pilot starts Claude Code 2.1.312 in a new session, `47f3eec4-8dd1-41ac-adfb-cbc7c6aff228` | the run's log |
| 04:36:34 to 04:37:16 | `story_start` over HTTP; S-0558 moves to in-progress at 04:37:10 | the log; `moved` event |
| 04:38:32, 04:39:21 | T-2457 and T-2458 closed with `task_done` | `moved` events |
| 04:40:13 to 04:40:21, 04:42:15 to 04:42:27 | Two close-out jobs on the host, each `passed` | `flai serve journal`; `flai verify S-0558 --last` at `6e4f917`: rebase, sync, narrative, check, markdown, template |
| 04:42:37 | S-0558 moves to review | `moved` event |
| 04:42:48 | Claude Code exits 0; the run ends `worked`; `docker ps -a` no longer lists the container (`--rm`) | `agent` event `ended`, why "ended with its work done" |

**The start.** From `flai serve`'s start to the agent's first event took 8 s in the container, against 4.5 s for a host run started at 04:14:38 (S-0570). The container was up 4 s after the start.

**The mounts.** `docker inspect` listed the eight mounts `design/tech/docker.md` names, and nothing else: the worktree read-write; the main checkout's `.git` read-only and `.git/worktrees/S-0558` read-write; `wip/` and `system-flow.yaml` read-only and `.flai-cache/verify` read-write; the story's cache folder read-write; the login read-only at `/run/claude-login`. The environment carried the names ADR-0246 and ADR-0253 give, the token by name only.

**Every step that needed the host.** Each commit was the host's, through `task_done` over HTTP, as ADR-0246 has it; the branch's commits carry the host's git identity. Both close-outs ran as `flai serve` jobs on the host (ADR-0229). One step went wrong inside: the agent ran `flai stream sync` in the container, whose rebase could not lock the branch's ref on the read-only `.git` and left HEAD detached; the agent reattached it with `git symbolic-ref` and synced through `task_done` on the host. S-0558's agent filed that as I-0272.

**The MCP calls over HTTP.** None failed, but the slowest came near Claude Code's 60 s limit on an HTTP tool call (I-0228): `story_start` took 41.8 s, against 4.9 s over stdio for S-0570 on the host, and the three `task_done` calls took 26.2 s, 21.6 s, and 14.2 s. The MCP file the container adapter writes sets no per-server timeout.

**The usage.** Measured from the run's log as for a host run, once S-0558 was in review:

| | S-0558, container | Host runs of 2026-11-07 |
|---|---|---|
| `source` | `log` | `log` |
| `priced_by` | `harness` | `harness` |
| Agent seconds | 376 (forecast 12m) | S-0557 1554, S-0571 1575, S-0555 2417 |
| Cost | 4.12 USD, `claude-opus-5-5` | 5.35, 6.85, 10.80 USD |
| Turns | 78 (`work` 38, `test_runs` 1 on the day) | |

S-0558 had two tasks, the host runs three or four, so the figures say the measure works, not that a container run is cheaper.

**The card.** The dashboard's card and story page read `agent.status` from `flai serve` (`flai/internal/hostapi/hostapi.go`) and `agent.stream` from the run's log. For S-0558, `~/.flai/serve/agents.json` keeps the same fields as for a host run (story, agent, session, started, ended, exit 0, outcome `worked`, `harness_version` 2.1.312), with `harness: container` and `command: docker`. `flai serve agent stream S-0558` reads the log as it reads a host run's, down to its `end` line "session ended (success, 78 turns, 378s, $4.12)".

**Limitations met.** I-0272: `flai stream sync` inside the container cannot rebase and leaves HEAD detached (filed by S-0558's agent, on its branch until it is accepted). [I-0228](../issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md), bumped: `story_start` and `task_done` over HTTP within 18 s of the 60 s limit. I-0270, also filed by S-0558's agent, is not the container's: the integration and smoke tiers do not select a change only under `template/` on the host either.

**The stop and the restart.** Neither happened on S-0558: the run reached review 5 minutes after the operator was asked for the stop, and its agent asked nothing on a thread. Both move to S-0556's run (TH-0836).

## Recommendation

Pending: whether to adopt, adapt, or drop running this repository's stories on `harness: container` is decided once both proof runs are recorded (T-2453).
