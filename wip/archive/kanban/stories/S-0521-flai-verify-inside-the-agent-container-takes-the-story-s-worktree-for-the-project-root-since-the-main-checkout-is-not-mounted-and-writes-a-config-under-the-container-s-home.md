---
id: S-0521
type: story
nature: remediation
title: flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home
status: done
owner: alex
created: 2026-10-28T12:53:27Z
updated: 2026-10-28T14:05:53Z
transitions:
  - to: ready
    at: 2026-10-28T13:06:32Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T13:06:59Z
    by: agent-S-0521
  - to: review
    at: 2026-10-28T13:49:56Z
    by: agent-S-0521
  - to: done
    at: 2026-10-28T14:05:53Z
    by: orchestrator
tags: [cli, docker]
topics: [harness, docker]
touches: [flai/internal/workitem/repo.go, flai/internal/workitem/worktree_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/flai-cli.md, design/issues/I-0219-flai-verify-inside-the-agent-container-takes-the-story-s-worktree-for-the-project-root-since-the-main-checkout-is-not-mounted-and-writes-a-config-under-the-container-s-home.md, design/issues/summary.md, flai/cmd/stream.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2596
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 2
      hand_edits: 1
      work: 40
  models:
    - model: claude-opus-5-5
      input: 224
      output: 56734
      cache_read: 10398773
      cache_write: 597399
      cost: 6.8987
  strategic:
    - kind: orchestrator
      seconds: 249
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 37
          output: 473
          cache_read: 8895760
          cache_write: 24719
          cost: 1.9584
        - model: claude-sonnet-5-5
          input: 10
          output: 69
          cache_read: 110523
          cache_write: 38085
          cost: 0.1016
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-28T12:53:27Z
  value: 12.5
  by: planner-S-0521
  at: 2026-10-28T13:02:09Z
forecast:
  duration: 1h
  delivery: 2026-10-29T08:22:00Z
  basis: "Its own forecast of 1h; 35th in the pull order with an in-progress limit of 5, behind S-0502, S-0415, S-0414, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517, S-0518, S-0519 and S-0520."
  by: flai
  at: 2026-10-28T13:02:30Z
finalized:
  by: alex
  at: 2026-10-28T12:53:59Z
---
# S-0521 flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home

## Goal

This story remediates [I-0219](../../../design/issues/I-0219-flai-verify-inside-the-agent-container-takes-the-story-s-worktree-for-the-project-root-since-the-main-checkout-is-not-mounted-and-writes-a-config-under-the-container-s-home.md), "flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0219 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0219 is closed with `flai issue close I-0219 --reason` saying what fixed it

## Tasks
- T-2207 A Repo opened from a linked worktree whose main checkout is absent records the main checkout the .git file names and that it is not here, in place of taking the worktree for the project
- T-2209 scripts/agent-container.sh gives the container's flai its config and cache under the mounted agent folder with FLAI_CONFIG and FLAI_CACHE_DIR, so nothing is written under the container's home
- T-2211 flai verify, stream sync, task done, and check refuse in a worktree whose main checkout is absent, naming it and the host's MCP tools, with a test that reproduces I-0219
- T-2212 docs/users/flai.md and flai-cli.md say how flai finds a worktree's project and what it refuses when the main checkout is absent, the fix is checked once inside the container, and I-0219 is closed

## Notes

Cost of delay inputs set by flai from I-0219. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:31:04Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0521 on 2026-10-28 (the plan thread names the tasks, their layers, and the assumptions).

The cause, from the code. `workitem.Open` (`flai/internal/workitem/repo.go`) walks up to the nearest `system-flow.yaml`, which a story worktree carries, then `mainCheckout` reads the worktree's `.git` file, cuts the main checkout's path before `/.git/worktrees/`, and keeps it only when `<main>/system-flow.yaml` exists. Inside the container the main checkout's files are absent, so `MainRoot` falls back to the worktree, `WorktreePath` becomes `<worktree>/.flai-cache/worktrees/S-nnnn`, and `flai verify` reports the path I-0219 quotes; `flai check` reads the worktree's tracked `wip/` for the board's. The config comes from `config.Load`, which creates `$HOME/.flai/config.json` when neither `--config` nor `FLAI_CONFIG` names a file (ADR-0008); `scripts/agent-container.sh` passes neither, so the image's `HOME=/home/agent` takes it.

The remediation proposed: flai in a linked worktree whose main checkout is not here says so and refuses the commands that need it, in place of taking the worktree for the project, naming the main checkout and the host (the MCP tools `verify`, `task_done`, and `stream_sync`, which run there and which S-0505's agent used); and the container's flai is given its config and cache under the mounted agent folder, so nothing is written under the container's home. The design's rule stands: the container mounts the worktree alone, `wip/` is written on the host through flai mcp (agent-adapters.md § A container harness of flai's own, ADR-0019). Running `flai verify` inside the container with what it needs is S-0510's (its criterion 4), which builds on the refusal this story adds rather than on a silent wrong root.

Touches. `flai touches suggest S-0521` names no path on its own, since the story declares none; seeded with `repo.go`, `verify.go`, and `agent-container.sh` (21 of 2473 commits) it lists `design/system/flai-cli.md` (43%), `flai/cmd/verify_test.go` (33%), `docs/users/flai.md` (29%), `docs/users/flai-reference.md` (29%), `flai/internal/verify/story.go` (24%), `cmd/root.go`, `cmd/test.go`, `mcpserver/verify.go`, and the conventions. Kept from it (co-change): `flai-cli.md`, `verify_test.go`, `docs/users/flai.md`. Left out: `flai-reference.md` (generated from the command help, and no flag changes), `internal/verify/story.go` (the refusal sits in `cmd/verify.go` before the run starts, where `storyWorktree` is), `root.go`, `test.go`, `mcpserver/verify.go` (the MCP tool runs on the host, where the main checkout is), and the conventions. The rest are the tasks' (layout):

- `flai/internal/workitem/repo.go`, `worktree_test.go`: `mainCheckout` and `Repo` learn that the main checkout the `.git` file names is absent, and the one test of `mainCheckout` gains the case. T-A.
- `flai/cmd/verify.go`, `verify_test.go`, `flai/cmd/stream_sync.go`, `stream_sync_test.go`, `flai/cmd/check.go`, `check_test.go`, `flai/internal/taskdone/taskdone.go`, `taskdone_test.go`: the four commands the issue and S-0505's results name as needing the main checkout refuse when it is absent; `verify_test.go`'s `verifyFixture` already builds `root/.flai-cache/worktrees/S-004` with a `.git` file, so the reproducing test lives there. T-C.
- `scripts/agent-container.sh`, `docs/operators/runbooks/agent-container.md`: the hand start passes `FLAI_CONFIG` and `FLAI_CACHE_DIR` under the mounted agent folder, and the runbook's table of variables and mounts says so. T-B. The script is also S-0507's (T-2174, backlog); whichever starts second is held on the overlap, and the pull order settles which.
- `docs/users/flai.md`, `design/system/flai-cli.md` (co-change, design): how flai resolves a worktree's project and what it refuses when the main checkout is absent. T-D.
- `design/issues/I-0219-….md`, `design/issues/summary.md` (criterion 2): `flai issue close` writes both; `summary.md` is in the manifest's shared paths and holds nothing (ADR-0096). T-D.

No folder touch is kept. No `after`: the story needs no other story's change on main; S-0507 and S-0508 are told on the plan thread what pilot and the adapter should pass.

Forecast. `flai forecast S-0521` gives 5m (median 131 s per unit over 12 done remediation stories on claude-opus-5-5 in the small band; size 2: 2 criteria, 0 touches) and delivery 2026-10-29T08:42Z, 35th in the pull order. Raised to 1h: the size flai measured had no touches, and with 16 the band's median alone gives about 40m; the story changes `Repo` and four commands, writes a fixture with a worktree whose main checkout is absent, and checks the fix by hand inside `flai-agent-claude-code:local` (`scripts/agent-container.sh S-nnnn -- flai verify S-nnnn`), which needs the image built with the tree's flai. The delivery written is flai's play-out start, 2026-10-29T08:42Z less 5m, plus 1h: 2026-10-29T09:37Z.

Cost of delay. `flai cod S-0521` gives 12.50 USD a week: 5m lost per 168h cycle at 150 USD an hour, from the input flai set from I-0219. Written as flai gives it: the issue has one instance, and nothing argues for another input.

### Accepted by the orchestrator

- Verified: 5221dd99983adbc34dd89a46a53bcee5455ffc46
- Batch: 20261028-135048
- Enqueued verified: 5221dd99983adbc34dd89a46a53bcee5455ffc46
- At: 2026-10-28T14:05:53Z

Verdict: pass at 5221dd99983adbc34dd89a46a53bcee5455ffc46
- 1: flai/internal/workitem/repo.go, flai/internal/workitem/worktree_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/check.go, flai/cmd/check_test.go, flai/cmd/stream.go, flai/cmd/stream_sync_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, scripts/agent-container.sh
- 2: design/issues/I-0219-flai-verify-inside-the-agent-container-takes-the-story-s-worktree-for-the-project-root-since-the-main-checkout-is-not-mounted-and-writes-a-config-under-the-container-s-home.md, design/issues/summary.md
