---
id: S-0520
type: story
nature: remediation
title: scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks
status: backlog
owner: alex
created: 2026-10-28T12:53:17Z
updated: 2026-10-28T13:58:16Z
transitions: []
tags: [cli, docker]
topics: [harness, docker]
touches: [design/adrs, design/adrs/README.md, flai/internal/storygit/mounts.go, flai/internal/storygit/mounts_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/cmd/verify_test.go, scripts/agent-container.sh, scripts/README.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/issues/I-0218-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md, design/issues/summary.md]
after: [S-0517, S-0518, S-0521]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 42
          cache_read: 591888
          cache_write: 5608
          cost: 0.1312
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-28T12:53:17Z
  value: 25
  by: planner-S-0520
  at: 2026-10-28T13:56:47Z
forecast:
  duration: 1h15m
  delivery: 2026-10-29T04:00:00Z
  basis: "flai forecast S-0520 gives 4m from size 2 (2 criteria, no touches; median 119 s per unit over 13 done remediation stories on claude-opus-5-5, small band); with 17 touches the median alone gives about 38m, raised to 1h15m for an ADR, a storygit and container-adapter change with tests, a verify fixture test, a check fallback, the hand-start script, docs in three places, a check by hand inside flai-agent-claude-code:local built from the tree, and the issue close; delivery is flai's play-out start, 2026-10-29T02:49Z less 4m, plus 1h15m, which the after on S-0517, S-0518, and S-0521 may move."
  by: planner-S-0520
  at: 2026-10-28T13:56:47Z
---
# S-0520 scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks

## Goal

This story remediates [I-0218](../../../design/issues/I-0218-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md), "scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0218 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0218 is closed with `flai issue close I-0218 --reason` saying what fixed it

## Tasks
- T-2232 An ADR refining ADR-0203 records that a story's container mounts the main checkout's wip and manifest read-only and its verify folder read-write, so flai verify and the close-out run inside it against the main checkout's state
- T-2233 flai verify reads from the main checkout only its .git, wip, and system-flow.yaml, the markdown lint configuration falling back to the worktree's, with a test whose fixture's main checkout holds nothing else
- T-2234 storygit's mounts helper gives the main checkout's wip and manifest read-only and .flai-cache/verify read-write beside the two git mounts, and the container adapter's Start passes them
- T-2235 scripts/agent-container.sh mounts the main checkout's wip and manifest read-only and .flai-cache/verify read-write, creating the verify folder first, and the runbook and scripts index say so
- T-2236 The docs say flai verify and the close-out run inside the container with the main checkout's wip, manifest, and verify folder mounted, the fix is checked once inside the container, and I-0218 is closed

## Notes

Cost of delay inputs set by flai from I-0218. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:31:00Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0520 on 2026-10-28 (the plan thread names the tasks, their layers, and the assumptions).

The cause, from the instance and the code. The close-out's first step is `scripts/flai.sh verify`, and `scripts/flai.sh` runs `go build` when the tree's `bin/flai` is missing or stale; a fresh worktree has no `bin/`, and `flai-agent-claude-code:local` has no `go`, so the step exits 127 before `flai verify` starts. Three stories already carry a part of the remedy the instance names: S-0518's T-2224 makes `scripts/flai.sh` run the installed flai, with a warning, when `go` is absent, which is the image's flai; S-0510's layer carries the tiers' tools (`go`, golangci-lint, Node, pnpm); and S-0521 makes `flai verify` refuse, naming the host's MCP tools, when the main checkout is absent, in place of taking the worktree for the project. What no story carries is what the refusal leaves: inside the container `flai verify` still has no main checkout, so the close-out still stops at its first step once S-0518 lands, with S-0521's refusal in place of exit 127. `flai verify` reads from the main checkout its `.git` (the base branch, the story's diff, the refs; mounted read-only by S-0517), its `wip/` (the narrative's `## Current state` and `## Next steps`, the board, the items, the threads), its `system-flow.yaml` (which is how `workitem.Open` knows the main checkout is there, and the queue check), and its markdownlint configuration (`check.go`'s `markdown`, from `MainRoot`), and it writes its report, which `flai verify --last` and the orchestrator read on the host, under `<main>/.flai-cache/verify/`.

The remediation proposed: a story's container mounts, beside the worktree and S-0517's two git mounts, the main checkout's `wip/` and `system-flow.yaml` read-only and its `.flai-cache/verify/` read-write, at their host paths, so that `flai verify` inside the container checks the same narrative, board, and branch as on the host and records its report where the host reads it; the markdown lint configuration falls back to the worktree's when the main checkout's cannot be read, so no single configuration file needs a mount; and a test with a fixture whose main checkout holds only `.git`, `wip/`, and `system-flow.yaml` guards that `flai verify` reads nothing else from it. An ADR refining ADR-0203 records the mounts, as S-0517's does for `.git`. Nothing in `scripts/close-out.sh` changes: its steps are right, and what they run is what the container lacked. Two parts stay other stories': the tiers' tools (S-0510) and the commit step, which fails inside the container when `--record-issues` recorded something, since S-0517 keeps `.git/objects` read-only and commits the host's; the plan thread raises that for S-0517's ADR.

Touches. `flai touches suggest S-0520` names no path on its own, since the story declared none; seeded with `scripts/close-out.sh`, `scripts/flai.sh`, `flai/cmd/verify.go`, `flai/internal/verify/story.go`, `scripts/agent-container.sh`, and `flai/internal/harness/adapters.go` (49 of 2478 commits) it lists `harness_test.go` (33%), `harness.go` (22%), `verify/story_test.go` (18%), `template/root/scripts/close-out.sh` (18%), `flai/cmd/verify_test.go` (16%), `scripts/README.md` (16%), `docs/users/flai-reference.md` (12%), `work-management.md`, `flai-cli.md`, `template/CHANGELOG.md`, `design/issues/summary.md`, `docs/users/flai.md`, `serve/agents.go`, `verify/verify.go`, `design/adrs/README.md`, `serve_actions.go`, and `cmd/test.go`. Kept from it (co-change): `flai/cmd/verify_test.go`, `scripts/README.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `design/issues/summary.md`, `design/adrs/README.md`. Left out: `harness.go` and `harness_test.go` (the adapter contract does not change; the mounts are the container adapter's own), `verify/story.go`, `story_test.go`, and `verify.go` (verify's steps do not change; only where `flai check` reads the lint configuration does), `close-out.sh` and the template's (unchanged, as above), `flai-reference.md` (no flag changes), `work-management.md` and `template/CHANGELOG.md` (no convention or template change), `serve/agents.go`, `serve_actions.go`, and `cmd/test.go` (the start path and `flai test` are not touched). The rest are the tasks' (layout):

- `design/adrs/README.md`, and the folder `design/adrs` (layout, folder kept): `flai adr new` in the first task writes the ADR under a number and slug no task can name yet; the folder is kept for that file alone, is in the manifest's shared paths, so it holds no ready story, and `flai task done` widens the touches with the committed path (ADR-0144).
- `flai/internal/check/check.go`, `check_test.go` (layout): the `markdown` check's configuration read, with its fallback to `Root`.
- `flai/cmd/verify_test.go` (co-change, layout): `verifyFixture` already builds a root with a linked worktree; the reproducing test strips the root to `.git`, `wip/`, and `system-flow.yaml`.
- `flai/internal/storygit/mounts.go`, `mounts_test.go`, `flai/internal/harness/container.go`, `container_test.go` (design: S-0517's helper and S-0508's adapter, not on main yet, which is why `after` names S-0517): the three mounts beside the two git mounts, and the adapter's `Start` passing them.
- `scripts/agent-container.sh`, `docs/operators/runbooks/agent-container.md`, `scripts/README.md` (design, layout): the hand start mounts the same, and the runbook's mount table and the scripts index say so.
- `docs/users/flai.md` § Verify a story before review, `design/system/flai-cli.md`, `design/system/agent-adapters.md` (co-change, design): what the container mounts for `flai verify` and what the close-out does inside it.
- `design/issues/I-0218-….md`, `design/issues/summary.md` (criterion 2): `flai issue close` writes both; `summary.md` is in the manifest's shared paths and holds nothing (ADR-0096).

`after`: S-0517 (the mounts go beside its `storygit` helper and into the container adapter's `Start`, on main only once S-0517 and, through it, S-0508 land), S-0518 (its T-2224 is what lets `scripts/flai.sh` reach the image's flai without `go`, the cause the issue names, which this story's check by hand runs through), and S-0521 (its `FLAI_CONFIG` and `FLAI_CACHE_DIR` in `scripts/agent-container.sh`, and the refusal this story turns into a pass). No `after` on S-0510: its layer is for the tiers, and its criterion 5, a close-out passing inside the container, needs this story's mounts rather than the other way round; the pull order settles which lands first. Overlaps on 2026-10-28: S-0521, in progress, touches `scripts/agent-container.sh`, `docs/operators/runbooks/agent-container.md`, and `docs/users/flai.md`; S-0414 and S-0502, in progress, touch `docs/users/flai.md`; S-0517 and S-0518 (backlog) share the script, the runbook, `scripts/README.md`, `container.go`, `mounts.go`, `design/adrs`, `agent-adapters.md`, and `docs/users/flai.md`, which `after` orders.

Forecast. `flai forecast S-0520` gives 4m (median 119 s per unit over 13 done remediation stories on claude-opus-5-5 in the small band; size 2: 2 criteria, 0 touches) and delivery 2026-10-29T02:49Z, 33rd in the pull order. Raised to 1h15m: the size flai measured had no touches, and with 17 the band's median alone gives about 38m; the story writes an ADR, changes the `storygit` helper and the container adapter with their tests, writes a verify fixture test and a check fallback, changes the hand-start script, checks the fix by hand inside `flai-agent-claude-code:local` (`scripts/agent-container.sh S-nnnn -- scripts/close-out.sh S-nnnn`, which needs the image built from the branch), changes three documents, and closes the issue. The delivery written is flai's play-out start, 2026-10-29T02:49Z less 4m, plus 1h15m: 2026-10-29T04:00Z; the `after` on three stories may move it, and flai's next replan plays it out again.

Cost of delay. `flai cod S-0520` gives 25.00 USD a week: 10m lost per 168h cycle at 150 USD an hour, from the input flai set from I-0218. Written as flai gives it: the issue has one instance, and the input is the operator's to raise if every container run of this repository stopping at its close-out weighs more than its one occurrence.
