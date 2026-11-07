---
id: S-0567
type: story
nature: feature
title: flai serve starts the planner in a container on harness container, and it plans one backlog story of this repository with its activity logged
status: done
parent: E-0029
owner: alex
created: 2026-11-02T19:02:15Z
updated: 2026-11-07T06:08:54Z
transitions:
  - to: ready
    at: 2026-11-07T03:43:06Z
    by: alex
  - to: in-progress
    at: 2026-11-07T03:43:19Z
    by: agent-S-0567
  - to: review
    at: 2026-11-07T05:45:58Z
    by: agent-S-0567
  - to: done
    at: 2026-11-07T06:08:54Z
    by: alex
blocked:
  - from: 2026-11-07T03:50:10Z
    until: 2026-11-07T03:54:00Z
    reason: awaiting the operator on TH-0838
  - from: 2026-11-07T03:54:43Z
    until: 2026-11-07T04:06:50Z
    reason: awaiting the operator on TH-0843
  - from: 2026-11-07T04:07:38Z
    until: 2026-11-07T05:38:11Z
    reason: awaiting the operator on TH-0843
tags: []
topics: [harness, docker, planning]
touches: [wip/agents/planner.md, design/system/strategic-agents.md, design/issues, design/issues/summary.md, design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md, design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md]
after: [S-0509]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 959
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 5
      hand_edits: 6
      work: 87
  models:
    - model: claude-opus-5-5
      input: 212
      output: 48692
      cache_read: 25232215
      cache_write: 554204
      cost: 10.4548
  strategic:
    - kind: orchestrator
      seconds: 183
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 327
          cache_read: 2238005
          cache_write: 193891
          cost: 0.5919
        - model: claude-sonnet-5-5
          input: 6
          output: 61
          cache_read: 43729
          cache_write: 36959
          cost: 0.0579
forecast:
  duration: 1h
  basis: "The planner's figure on TH-0838: flai's 19m size estimate plus the operator's setup and watching one planner session, as S-0565 is forecast at 1h30m"
  by: agent-S-0567
  at: 2026-11-07T03:49:29Z
finalized:
  by: alex
  at: 2026-11-07T03:41:47Z
---
# S-0567 flai serve starts the planner in a container on harness container, and it plans one backlog story of this repository with its activity logged

## Goal

S-0509 lets the planner, the orchestrator, the analyzer, and the curator run in a container (ADR-0247). It proved this with docker faked, since the host's flai serve could not run the story's code before its release (TH-0815). This story makes the live run: `flai serve` starts a planner on `harness: container`, and it plans one backlog story of this repository in the container. It shows that the mounts, the guard, and the activity entry hold on the real path.

## Acceptance criteria

- [x] With the host's container harness set (`flai serve agent harness container --login <path>`) and the planner's agent on `harness: container`, `flai serve` starts a planner in `flai-agent-claude-code:local` (or the project's layer). The planner plans one backlog story of this repository, and its plan's thread and drafted tasks are written through flai.
- [x] `wip/agents/planner.md` gains the run's activity entry, with duration and cost measured from the container's log as from a host run.
- [x] The guard refuses the planner inside the container what it refuses on the host, shown by at least one refusal in the run's log or by a call made by hand in the container.
- [x] Each limitation met is an issue under `design/issues`, and `design/system/strategic-agents.md` § In a container (S-0509) says what the live run showed.

## Tasks
- T-2498 The live planner run's preconditions are in place, and the operator has set the planner on harness container and asked for a plan of one backlog story
- T-2499 The planner's container run is watched to its end: it plans the chosen story through flai, and wip/agents/planner.md gains its entry with duration and cost from the container's log
- T-2500 A guard refusal of the planner inside its container is shown from the run's log, or from a call made by hand in the container with the planner's role
- T-2501 Each limitation the live run met is an issue, and strategic-agents.md § In a container (S-0509) says what the live run showed
- T-2502 The live planner run is checked: its container, its plan, its activity entry's duration and cost, and a guard refusal inside the container
- T-2503 Each limitation the live planner run met is an issue, and strategic-agents.md § In a container says what the run showed

## Notes

Waits for S-0509 (`after`), whose start paths and adapter it runs, and for a flai release that holds them. Before the run the operator sets the host's container harness login and image with `flai serve agent harness container`. An agent cannot do that, since the guard refuses it a host-configuration write (ADR-0220). Moved here from S-0509's criterion 4 on TH-0815.

### Live run

Checked by agent-S-0567 on T-2498, 2026-11-07:

| Check | Found |
|-------|-------|
| Installed flai | 1.83.0 at `~/.flai/bin/flai`, built 2026-11-03; S-0509 is in every release from `flai/v1.81.0`. `flai serve` 1.83.0 runs since 2026-11-03T18:12:07Z |
| Image | `flai-agent-sf-claude-code:local`, the layer `agent.image.layer` names, over `flai-agent-claude-code:local`. The images on the host were built 2026-11-02T18:51Z, before S-0509's pilot commit (18:56Z), so `make agent-images` rebuilt all three from main at `e41d6124e`: `sha256:521d859c7bca`, Claude Code 2.1.312, flai 1.84.0 inside |
| Host container harness | program `docker`, default bridge network; login not set |
| Planner's agent | `planning.agent` names a model and no harness |
| Story to plan | S-0559, the first story `flai plan --candidates` lists with no plan |

What the operator set, on TH-0843 (confirmed 2026-11-07T05:38:10Z):

- The container harness's login is `/home/alex/.claude-me`, mounted read-only. It is neither `~/.claude` nor `~/.claude-rmh`, the two folders the threads named.
- `planning.agent` is `harness: container`, `model: claude-opus-5-5`.
- alex asked for the plan of S-0559 at 2026-11-07T05:05:44Z. The journal says `started docker (container) to plan S-0559 as planner-S-0559`, with its log at `~/.flai/serve/agents/sf-planner-20261107T050545Z.log`.

The run, checked on T-2499:

| Check | Found |
|-------|-------|
| Harness | `~/.flai/serve/agents.json` records the run under `plans/S-0559`: `harness: container`, `command: docker`, `agent: planner-S-0559`, 05:05:45Z to 05:09:58Z, outcome `worked`, `harness_version: 2.1.312` |
| Image | The log opens `pilot: starting claude-code in a new session` and ends `pilot: claude-code exited with code 0`, so pilot, the image's entrypoint, ran it. Its Claude Code is 2.1.312, the image's, where the host's is 2.1.319. By `containerImageHost` in `flai/internal/serve/container.go`, which every container start takes, the image is the project's layer, `flai-agent-sf-claude-code:local`. No record names the image or the container, and Docker's events no longer hold them ([I-0274](../../../design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md)) |
| Mounts | The init event's `cwd` is the main checkout, `/home/alex/git/bytepunx/system-flow`. The role's cache folder `.flai-cache/agents/plan` was made at 05:05 and holds `claude`, `flai`, `pilot`, and `mcp.json`. The `.git/hooks` read-only mount is not in any record ([I-0274](../../../design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md)) |
| Plan | Written through flai's MCP tools alone: 2 `item_new` (T-2539, T-2540), 4 `item_edit` (S-0559's touches, forecast, cost of delay, notes), 1 `thread_open` (TH-0847). No `Edit` or `Write` call. The `planned` event names T-2539 and T-2540 |
| Activity entry | `wip/agents/planner.md`, 2026-11-07T05:09:49Z: `Items: S-0559, T-2539, T-2540`, `Seconds: 239`, `Cost: 2.1455 USD, estimated`. The log's result has `duration_ms: 239773` and `costUSD: 2.1454874`, so both are measured from the container's log, and the entry reads as a host run's does |

The guard, checked on T-2500:

- The run's log holds no refusal. The planner made no call the guard refuses, and no tool result is an error.
- By hand, `flai guard` was given four calls with `FLAI_ROLE=plan` and `FLAI_STARTED_BY=flai-serve`, inside `flai-agent-sf-claude-code:local` (its flai 1.84.0) and on the host (flai 1.84.0). Both gave the same verdicts, word for word:

| Call | Container | Host |
|------|-----------|------|
| `Write` to `design/system/x.md` | exit 2: "the planner cannot use Write" | exit 2, the same |
| `Bash` `git commit -m x` | exit 2: "it runs only git's reads" | exit 2, the same |
| `mcp__flai__item_move` S-0559 to ready | exit 2: "the planner cannot move S-0559 to ready" | exit 2, the same |
| `mcp__flai__item_get` S-0559 | exit 0 | exit 0 |

- In the run itself, the guard is the project's hooks, which run `scripts/flai.sh guard` from the checkout. A container started by hand from this agent's shell saw the checkout as an empty folder ([I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md)), so those hooks could not be replayed here. They also let a call through on any exit but 2, 126, and 127 ([I-0275](../../../design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md)).
