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

**The stop and the restart.** Neither happened on S-0558: the run reached review 5 minutes after the operator was asked for the stop, and its agent asked nothing on a thread. Both moved to S-0556's run (TH-0836).

### The second run: S-0556, Go, on this repository's layer

S-0556 was worked from `story_start` to review in three container runs on `flai-agent-sf-claude-code:local`, the layer `agent.image.layer` names. It went to ready at 06:17:39 and to review at 06:57:43. The operator moved it to ready and stopped it once from the dashboard, and its agent asked once on a thread.

| Run | Session | Started | Ended | How it ended |
|-----|---------|---------|-------|--------------|
| 1 | `fc3d3d78-…` (new) | 06:17:41 | 06:25:00 | `stopped` by the operator from the dashboard, exit 143 |
| 2 | `c8242e4c-…` (new) | 06:25:30, Retry | 06:51:36 | ended asking on TH-0849 |
| 3 | `c8242e4c-…` (resumed) | 06:51:55, TH-0849 answered | after 06:57:43 | `worked` |

The stop and the two restarts are written up in S-0565's Notes, under "The stop and the restart (T-2450)". In short:

- The stop removed the container within 3 s.
- Retry started a new agent in a new session, as Retry does on a host.
- The answer on TH-0849 started the container again in the same session, from the story's writable Claude Code folder, and the agent went on from where it stopped.
- The one limitation the operator met was on the dashboard (I-0278).

**The close-out on the host.** The `flai serve` job passed at `45b5126`: rebase, sync, narrative, check, gofmt, vet, golangci-lint, markdown, template, integration (3m52s), and smoke (44s). `go-test` was skipped, covered by integration (`flai verify S-0556 --last`). The orchestrator accepted S-0556 at 06:59:15.

**The close-out inside the container.** TH-0821's option (a) runs it once by hand inside the layered image on S-0556's worktree. Four things stood in the way:

1. `scripts/agent-container.sh` always runs `flai-agent-claude-code:local` and ignores `agent.image.layer` (I-0279). The run used `docker run` on `flai-agent-sf-claude-code:local` directly, with the mounts and environment `flai serve` gave S-0556's agent.
2. At a commit whose pass the verify record holds, `flai verify` answers from the record and runs nothing (ADR-0229 §5). `scripts/close-out.sh` cannot ask it for `--fresh` (I-0280). The run used the close-out's own two verify steps, `scripts/flai.sh verify S-0556 --fresh --record-issues` and `--sync-only`. The worktree was clean, so its commit step had nothing to do.
3. The orchestrator's acceptance run had started at about 06:58:21. It accepted S-0556 at 06:59:15 and removed the worktree under the run, 36 s after a hold was asked for on TH-0836. Before the removal, rebase, sync, narrative, check, gofmt, vet, and golangci-lint had passed inside. The markdown tier then failed for a working folder that was gone. The run also left its result in S-0556's verify record at `45b5126`, which `.flai-cache/verify` mounted read-write lets a container write.
4. The rest of the tiers ran inside on a detached scratch worktree at `45b5126` under `/tmp`, removed afterwards. They used the same image and mounts, with `CLOSE_OUT_STORY=S-0556` for the smoke tier as a close-out sets it.

| Tier | Host close-out | Inside the container |
|------|----------------|---------------------|
| gofmt | passed | passed, 14 ms |
| vet | passed | passed, 16.5 s |
| golangci-lint | passed | passed, 7.5 s |
| go-test (`-race -short`) | covered by integration | passed, 3m38s |
| markdown | passed | passed, 6.7 s |
| template | passed | passed, 177 ms |
| integration | passed, 3m52s | passed, 4m10s |
| smoke | passed, 44 s | **failed** at its installer test |

Smoke inside passed its template render, its repository check (3152 items, 0 findings), and its markdown lint, then failed in `scripts/install-test.sh`. The layer has no curl, so `install.sh` stops with "curl is required". On a throwaway image of the layer plus curl, the test failed next at `scripts/flai-snapshot.sh`: "goreleaser: not found". GoReleaser lives in the host's `bin/`, not in the layer. Both are I-0281.

**I-0259.** No container run installed `flaiover/node_modules`: the markdown lint ran from the container's own cache, and S-0556 selected no flaiover tier. The host's flaiover tiers, run afterwards with `flai test flaiover` on the main checkout, passed: flaiover-lint, vitest, and svelte-check.

**The usage.** S-0556 records 2291 agent seconds over its three runs and 9.40 USD on `claude-opus-5-5`, with one empty wake. It is `priced_by: estimate`, not `harness`, because the stopped run's log has no result line to price it by. A run stopped on the host would read the same.

**Limitations met in this run.**

- I-0272, bumped by S-0556's agent: its in-container `flai stream sync` failed on the read-only `.git` again, this time with HEAD left attached.
- I-0228, bumped by S-0556's agent.
- I-0277, filed by S-0556's agent: an agent in a container cannot write its narrative's Decisions, since `wip/` is read-only there and `stream_state` writes only Current state and Next steps.
- I-0278: no Retry on the agent page after a stop.
- I-0279: the hand-start script ignores the layer.
- I-0280: `close-out.sh` cannot force a fresh run.
- I-0281: the layer lacks curl and GoReleaser, so smoke fails inside.

### Against the success measure

| Criterion | Result |
|-----------|--------|
| 1. A story worked end to end in its container, usage measured, card's status as for a host run | Met by S-0558, and again by S-0556. Usage was measured from the run's log as for a host run. The card reads the same `agent.status` and `agent.stream` records as a host run's, with `harness: container` |
| 2. A stop from the dashboard and a restart after an answered thread, each tried once | Met on S-0556's container (S-0565's Notes, "The stop and the restart (T-2450)"). The stop removed the container. The answered thread resumed the same session. Retry after a stop gave a new session, as on a host |
| 3. A Go story worked end to end on the layer, its close-out passing every selected tier inside the container | S-0556 was worked end to end on the layer, and its close-out passed as `flai serve`'s host job. Inside the container, seven of its eight tiers passed; smoke failed at its installer test for want of curl and GoReleaser (I-0281). The operator ticked this criterion on 2026-11-07 after the two runs (TH-0836); I-0281 is what stands between it and a close-out that passes inside |

Issues met, and what carries each on:

| Issue | What | Carried on by |
|-------|------|---------------|
| I-0228 | `story_start` and `task_done` over HTTP within 18 s of Claude Code's 60 s limit | S-0583, in progress |
| I-0272 | `flai stream sync` inside the container cannot rebase on the read-only `.git` | none yet |
| I-0277 | an agent in a container cannot write its narrative's Decisions | none yet |
| I-0278 | no Retry on the story's agent page after a stop; only on the card's menu | none yet; the dashboard's, not the container's |
| I-0279 | `scripts/agent-container.sh` ignores `agent.image.layer` | none yet |
| I-0280 | `scripts/close-out.sh` cannot force a fresh run | none yet |
| I-0281 | the layer lacks curl and GoReleaser, so smoke fails inside | none yet |

## Recommendation

**Adapt.** Keep the `container` harness, and keep it opt-in per story for now. Do not make it this repository's default yet.

What works is the part that matters most:

- `flai serve` starts the container with exactly the mounts the design names.
- The agent works the story through `story_start`, `task_done`, and the close-out job, all over HTTP, and every commit stays the host's.
- Usage and the card read as for a host run.
- A stop removes the container, and an answered thread resumes the same session.

A documentation story and a Go story each reached review this way. Neither needed help beyond the operator's host setup.

What must change first, in this order:

1. **I-0281.** Add curl and GoReleaser to `flai/images/system-flow/Dockerfile`, and check the built layer with `flai test smoke` in it. Until then no story touching `flai/` or `template/` can close out inside the container.
2. **I-0272.** Refuse `flai stream sync` inside the container, saying the host syncs through `task_done`, as ADR-0246 meant S-0521 to. Both proof agents tried it, and one was left with a detached HEAD.
3. **I-0228** (S-0583). The HTTP calls that do the most work came within 18 s of Claude Code's limit.
4. **I-0277.** An agent in a container needs a way to write its narrative's Decisions through flai.

I-0279 and I-0280 are for running things by hand in the container. I-0278 belongs to the dashboard. None of the three stands in the way of `flai serve`'s own runs.
