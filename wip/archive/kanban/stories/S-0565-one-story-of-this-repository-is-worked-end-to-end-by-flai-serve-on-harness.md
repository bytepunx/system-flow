---
id: S-0565
type: story
nature: experiment
title: One story of this repository is worked end to end by flai serve on harness container, from story_start to review, once S-0508's container harness is in the installed flai
status: done
parent: E-0029
owner: alex
created: 2026-11-02T18:08:56Z
updated: 2026-11-07T07:55:05Z
transitions:
  - to: ready
    at: 2026-11-07T03:40:51Z
    by: alex
  - to: in-progress
    at: 2026-11-07T03:41:00Z
    by: agent-S-0565
  - to: review
    at: 2026-11-07T07:17:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T07:51:22Z
    by: flai
  - to: review
    at: 2026-11-07T07:53:27Z
    by: system-flow
  - to: done
    at: 2026-11-07T07:55:05Z
    by: orchestrator
blocked:
  - from: 2026-11-07T03:44:57Z
    until: 2026-11-07T04:19:08Z
    reason: awaiting the operator on TH-0836
  - from: 2026-11-07T04:48:18Z
    until: 2026-11-07T06:10:18Z
    reason: awaiting the operator on TH-0836
  - from: 2026-11-07T06:12:37Z
    until: 2026-11-07T06:18:00Z
    reason: awaiting the operator on TH-0836
tags: []
topics: [harness, docker]
touches: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/experiments/README.md, design/issues, design/issues/summary.md, design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md, design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md, design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md, design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md]
after: [S-0508, S-0510, S-0517, S-0520]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5810
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 9
      hand_edits: 4
      work: 321
  models:
    - model: claude-opus-5-5
      input: 694
      output: 110340
      cache_read: 121126871
      cache_write: 869326
      cost: 33.3896
  strategic:
    - kind: planner
      seconds: 761
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 168
          output: 47513
          cache_read: 12486805
          cache_write: 483216
          cost: 6.2582
    - kind: orchestrator
      seconds: 944
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 271
          output: 1827
          cache_read: 9287393
          cache_write: 1292419
          cost: 2.5616
        - model: claude-sonnet-5-5
          input: 12
          output: 93
          cache_read: 155277
          cache_write: 40894
          cost: 0.1412
cost_of_delay:
  value: 123.85
  by: planner-S-0565
  at: 2026-11-07T03:42:23Z
forecast:
  duration: 1h30m
  delivery: 2026-11-07T05:11:00Z
  basis: "Kept at 1h30m over flai's 27m: two proof runs watched end to end, a stop and a restart, the issues, and the results document, from the story's start at 03:41Z; holds only once the host runs flai 1.84.0 and the proof stories reach ready."
  by: planner-S-0565
  at: 2026-11-07T03:42:23Z
finalized:
  by: alex
  at: 2026-11-07T03:40:41Z
---
# S-0565 One story of this repository is worked end to end by flai serve on harness container, from story_start to review, once S-0508's container harness is in the installed flai

## Goal

Prove the container harness S-0508 built in the place it is meant for: `flai serve` on this host, the installed release, starts one story of this repository on `harness: container` and its agent works it from `story_start` to review inside the container, with its usage measured and its card's status shown as for a host run. Each limitation met becomes an issue.

## Acceptance criteria

- [x] One story of this repository is worked end to end by `flai serve` with `harness: container` on its agent, from `story_start` to review, its usage measured and its card's status shown as for a host run, and each limitation met is an issue under `design/issues`
- [x] A stop from the dashboard and a restart after a thread is answered are each tried once on that story's container, and what happened is recorded in this story's notes
- [x] One story of this repository that changes Go is worked end to end by `flai serve` on `harness: container` with the project's image layer (`flai/images/system-flow/Dockerfile`, S-0510), and `scripts/close-out.sh` runs every tier the manifest's `tests` select and passes inside the container, and each limitation met is an issue under `design/issues`

## Tasks
- T-2449 The results document is started with the hypothesis, the success measure, and the host's preconditions, and the two proof stories carry harness container
- T-2450 A stop from the dashboard and a restart after an answered thread are each tried once on the first proof story's container, and what happened is in S-0565's Notes
- T-2451 The first proof story is worked from story_start to review in its container, its usage measured and its card's status shown as for a host run, and each limitation met is an issue
- T-2452 The Go proof story is worked to review in its container on this repository's image layer, its close-out passing every tier the manifest's tests select inside the container, and each limitation met is an issue
- T-2453 The results document gives the results against each criterion and a recommendation, and the experiments index lists it

## Notes

Moved out of S-0508 (its criterion 6 and the live-run half of T-2169) on TH-0812, as D-0461 records for S-0423. `after: S-0508`: the `container` adapter, its host settings (`flai serve agent harness container --image/--network/--login`), and the serve start path exist only on `story/S-0508`; they reach this host's `flai serve`, the installed release, only once S-0508 is accepted and published and the host upgraded.

The work, from T-2169: with the host settings in place (`flai serve agent harness container --login ~/.claude`) and the images built (`make agent-images`), set `harness: container` on one small story of this repository (`flai edit S-nnnn --harness container`), and let `flai serve` start it. The proof story should change documentation or shell only: the Go tiers and the close-out need the project's image layer (S-0510, S-0520), and git does not work in the container yet (S-0517). Record each limitation met with `flai issue new`, naming the story of E-0029 that carries it on; I-0253 and I-0254 are S-0508's already.

Criterion 3 moved here from S-0510 on TH-0814 (its criterion 5, and the close-out half of its criterion 4), as D-0466 records for S-0508. `after` gains `S-0510`, whose layer reaches this host's `flai serve` only once it is accepted and published, and `S-0517`, which makes git work in the container so the close-out can sync, commit, and check the branch (I-0215). S-0510 ran `flai test` by hand inside the layered image on its own worktree, and every per-path tier passed (S-0510's notes); what is left is a Go-changing story worked by `flai serve` with the layer, its close-out passing inside the container. Mind I-0259: installing `flaiover/node_modules` inside the Alpine container leaves only musl bindings, so the host's flaiover tier then fails.

### The stop and the restart (T-2450)

Tried on S-0556's container, not S-0558's. S-0558 reached review six minutes after its start, five minutes after the operator was asked for the stop, and its agent asked nothing on a thread (TH-0836). S-0556 was also worked on `harness: container` from `story_start` on. The operator added a line to its Notes asking its agent to confirm on a thread before its close-out, so that it would end on a question.

The stop from the dashboard:

- 06:17:41Z: `flai serve` started S-0556's container, `flai-agent-S-0556-fc3d3d78-a68a-4209-9f07-edf43d26d84b`, on `flai-agent-sf-claude-code:local`, in session `fc3d3d78-…`. Its agent moved S-0556 to in-progress at 06:18:19Z.
- 06:24:57Z: the operator pressed Stop. The `agent` event was `stopped`, why "stopped by the operator at 2026-11-07T06:24:57Z". The run ended at 06:25:00Z, outcome `stopped`, exit 143. The journal says "ended its process group … its container … had ended with it". `docker ps -a` no longer listed the container at 06:25:08Z. A forked task sub-agent (T-2459) was mid-call and ended with it.
- The card and the story's agent page showed the agent "stopped by the operator". The page offered no Retry; the card's right-click menu did (the operator's report on TH-0836). Filed as I-0278. Both use `retryable()` in `flaiover/src/lib/activity.ts`, which does not look at the harness.
- 06:25:30Z: Retry from the card's menu. The event was `started`, why "the operator restarted it: its last agent stopped by the operator at 2026-11-07T06:24:57Z". It started a new container in a new session, `c8242e4c-…` ("pilot: starting claude-code in a new session"). That is what Retry does on a host too: a new agent, told how the last one ended, that goes on from the narrative (`docs/users/flai.md`). The new agent went on with T-2459 and closed it at 06:32:09Z.

The restart after an answered thread:

- 06:50:26Z: S-0556's agent asked on TH-0849, as its Notes asked. Its run ended at 06:51:36Z, why "waiting for an answer to TH-0849". `flai serve` blocked S-0556, reason "awaiting the operator on TH-0849" (ADR-0205), and the container was gone. The session's store stayed in the story's writable Claude Code folder, `.flai-cache/agents/S-0556/claude/projects/…/c8242e4c-….jsonl` (ADR-0203).
- 06:51:54Z: the operator confirmed the orchestrator's recommendation on TH-0849. At 06:51:55Z `flai serve` unblocked S-0556 and started its container again, why "started again in its session, TH-0849 answered". The container kept the session's name, `flai-agent-S-0556-c8242e4c-…`, and the log says "pilot: starting claude-code resuming its session", with session `c8242e4c-…` throughout.
- The agent went on from where it stopped: its first calls read TH-0849 and its inbox, and it said it would sync and ask for the close-out.

Limitations met: I-0278 (no Retry on the agent page after a stop). The restart in the same session and the stop of the container needed nothing a host run does not.

### Planning

Planned by planner-S-0565 on 2026-11-03, and planned again on 2026-11-07 after the operator finalized the story and moved it to ready. The plan's thread TH-0821 named the tasks, their layers, and the assumptions, and was resolved with S-0558 and S-0556 as the proof stories, option (a) for criterion 3, and S-0520 in `after`. The replan's thread on S-0565 says what changed since then.

What this story's own branch changes. The proof stories are other stories, worked on their own branches by their own agents in the container; this story's agent sets them up, watches their runs, and records what it found. Its branch holds the results document, the issues it files, and the experiments index.

Touches. The story declared none at its first plan, so `flai touches suggest S-0565` ran on five seeds (`design/experiments/README.md`, `docs/operators/runbooks/agent-container.md`, `design/tech/docker.md`, `design/system/agent-adapters.md`, `design/issues/summary.md`; 749 of 2940 commits). On 2026-11-07 it ran again on the four touches the story declares (775 of 3012 commits): its top paths are `design/system/flai-cli.md` and `docs/users/flai.md` (5% each), then `design/adrs/README.md` and `design/system/flaiover-dashboard.md` (3% each), and each changes with the commands, the dashboard, or an ADR, which this story does not change; the limitations it finds become issues, and the fixes are stories of E-0029. So the touches stay as they were:

- `design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md` (layout): the results document an experiment story is accepted with, named for the story as `design/experiments/README.md` asks, started from `template.md`. The story's agent may choose another slug; `flai task done` widens the touches with the committed path (ADR-0144).
- `design/experiments/README.md` (layout): the index row for the results document.
- `design/issues`, folder kept (layout): criteria 1 and 3 make each limitation an issue with `flai issue new`, which names each file from its title; no task can name those files yet, so the folder stays on the story alone, and `flai task done` widens the touches with each committed issue.
- `design/issues/summary.md` (layout, co-change): `flai issue new` regenerates it; it is in the manifest's shared paths and holds nothing (ADR-0096).

The story's own Notes, where criterion 2's record goes, are written with `flai edit` in the main checkout and are no touch of its branch.

`after`. S-0508: the `container` harness adapter, its host settings, and the start path in `flai serve` are S-0508's output, which this story runs. S-0510: criterion 3 runs the project's image layer, `flai/images/system-flow/Dockerfile`, and the manifest key that names it, which S-0510 writes. S-0517: the close-out inside the container syncs and checks the branch with git, which works there only through the two git mounts S-0517 added (ADR-0246). S-0520, added by the planner: `flai verify`, the close-out's first step, reads the main checkout's `wip/` and `system-flow.yaml` and writes its report under `.flai-cache/verify/`, none of which the container mounts without S-0520 (I-0218); criterion 3's close-out inside the container reads those mounts. All four are accepted. The first release holding each is flai 1.81.0 for S-0508 and S-0517, 1.83.0 for S-0510, and 1.84.0 for S-0520, published 2026-11-07. On 2026-11-07 the host's installed flai was still 1.83.0, so T-2449's first precondition holds only once the host is upgraded to 1.84.0 or later; that is the operator's.

Forecast. `flai forecast S-0565` gave 9m at the first plan and 27m on 2026-11-07 (median 227 s per unit of size over 5 done experiment stories, size 7: 3 criteria, 4 touches). Kept at 1h30m of agent time: the size counts no waiting, and the work is two proof runs watched end to end, a stop and a restart arranged with the operator, each limitation filed, and the results document written; S-0505, one hand-built run with a results document, took 2014 s. The proof stories' own agent time is theirs, not this story's. The delivery written is the story's start, 2026-11-07T03:41Z, plus 1h30m: 2026-11-07T05:11Z. It holds only if the host is upgraded and the proof stories are moved to ready (TH-0826, TH-0828) without a wait.

Cost of delay. `flai cod S-0565` gives 123.85 USD a week on 2026-11-07, up from 52.12: S-0565's share of E-0029's 150 USD a week penalty input, 1h30m of the 1h49m forecast over the epic's two open stories without inputs, S-0565 and S-0567, now that the epic's other stories are done. Written as flai gives it; the story has no inputs of its own, and nothing argues for another share.
- 2026-11-07T07:51:22Z: moved to in-progress: the rebase of story/S-0565 onto main stopped on conflicts in design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, which the acceptance of S-0565 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0565

### Accepted by the orchestrator

- Verified: 437677d100cefcca0095dda76fb600cbb1b83996
- At: 2026-11-07T07:55:05Z

Verdict: pass at 437677d100cefcca0095dda76fb600cbb1b83996
- 1: design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md, design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md, design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md, design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md, design/issues/summary.md
- 2: design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md
- 3: design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md, design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md, design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md
