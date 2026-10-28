---
id: S-0517
type: story
nature: remediation
title: A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder
status: backlog
owner: alex
created: 2026-10-28T12:52:47Z
updated: 2026-10-28T18:46:21Z
transitions: []
tags: [cli, docker]
topics: [harness, docker]
touches: [design/adrs, design/adrs/README.md, flai/internal/storygit/mounts.go, flai/internal/storygit/mounts_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, scripts/agent-container.sh, scripts/README.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/agent-adapters.md, design/tech/docker.md, design/issues/I-0215-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md, design/issues/summary.md]
after: [S-0508]
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
      seconds: 390
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 594
          cache_read: 10804146
          cache_write: 11966
          cost: 2.3745
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:52:47Z
  value: 37.5
  by: planner-S-0517
  at: 2026-10-28T13:19:45Z
forecast:
  duration: 1h15m
  delivery: 2026-10-29T07:44:00Z
  basis: "Its own forecast of 1h15m; 29th in the pull order with an in-progress limit of 5, behind S-0415, S-0502, S-0531, S-0419, S-0423, S-0418, S-0421, S-0416, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515 and S-0516."
  by: flai
  at: 2026-10-28T18:46:21Z
finalized:
  by: alex
  at: 2026-10-28T12:53:25Z
---
# S-0517 A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder

## Goal

This story remediates [I-0215](../../../design/issues/I-0215-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md), "A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0215 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0215 is closed with `flai issue close I-0215 --reason` saying what fixed it

## Tasks
- T-2218 An ADR records the remediation proposed from I-0215's instance: the container mounts the main checkout's .git read-only and the story's worktree registration read-write, and commits stay the host's
- T-2219 storygit.GitMounts gives a story worktree's two git mounts, the common dir read-only and its registration read-write, with a real-git test that reproduces I-0215
- T-2220 scripts/agent-container.sh mounts the main checkout's .git read-only and the story's registration read-write by default, and FLAI_AGENT_GIT_MOUNT goes
- T-2221 The container adapter's Start mounts the main checkout's .git read-only and the story's registration read-write from storygit.GitMounts, and its argv test asserts nothing else under .git is writable
- T-2222 The docs say the container mounts .git read-only and the story's registration read-write and commits are the host's, the fix is checked inside the container, and I-0215 is closed

## Notes

Cost of delay inputs set by flai from I-0215. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:27:05Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0517 on 2026-10-28 (the plan thread names the tasks, their layers, and the assumptions).

The cause, from the instance and the code. `flai stream open` makes a story's worktree with `git worktree add` (`flai/internal/storygit/open.go`), which writes `<worktree>/.git` as a file, `gitdir: <root>/.git/worktrees/S-nnnn`, and in that folder a `commondir` of `../..`, the main checkout's `.git`. A container that mounts the worktree alone has neither, so every git command in it fails with `fatal: not a git repository`. S-0505 got round it by mounting the whole `.git` read-write (`FLAI_AGENT_GIT_MOUNT=1`), which lets the container write every branch, every other worktree's index, and the hooks and config of the whole repository, against ADR-0203's rule that the container holds no `.git` of the main checkout.

The remediation proposed, from what the container's git needs. Inside the container git only reads: `git status`, `diff`, `log`, and `show` for the agent's own work, and the index refresh `git status` writes. Every commit is the host's, through `task_done` and the close-out's MCP `verify` ([ADR-0203](../../../design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) § Decision, flai's MCP server; S-0521 makes `flai verify`, `stream sync`, and `task done` refuse inside the container). So the container mounts two more paths, at their host paths: the main checkout's `.git` read-only, and the story's own registration, `.git/worktrees/S-nnnn`, read-write over it, since git writes the worktree's index, `HEAD`, and locks there. The objects, refs, `packed-refs`, `config`, and hooks stay read-only, as ADR-0027 mounts the dashboard's hooks, config, and info; a `git commit` or `git add` in the container is refused by git with a plain write error, and the docs say the host commits. The alternative the issue names, a clone of the story branch fetched back, is the Kubernetes path ADR-0203 keeps open, and S-0505 measured that the bind mount is fast enough not to need it from the start. An ADR records the choice (first task), and the story's agent records another remediation instead if the plan's thread is answered with one.

Touches. The story declared none, so `flai touches suggest S-0517` ran on four seeds (21 of 2478 commits): `scripts/agent-container.sh`, `flai/internal/storygit/open.go`, `flai/internal/harness/adapters.go`, and `docs/operators/runbooks/agent-container.md`. Of what it listed, `docs/users/flai.md` (10%) and `design/issues/summary.md` (14%) are kept (co-change) for the reasons below. The rest are left out: `harness_test.go` (76%) and `harness.go` (52%) change with the adapters' contract, which this story keeps (the mounts are the container adapter's own, in `container.go`); `serve/agents.go`, `serve.go`, `cmd/serve_actions.go`, and `cmd/stream.go` change with how a start is driven, not with what is mounted; `flai-cli.md` and `agent-context.md` describe commands and priming, neither of which changes; `open_test.go` tests `flai stream open`, which this story leaves alone. The rest are the tasks' (layout):

- `design/adrs/README.md`, and the folder `design/adrs` (layout, folder kept): `flai adr new` in the first task writes the ADR under a number and slug no task can name yet; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144).
- `flai/internal/storygit/mounts.go`, `mounts_test.go` (layout): a helper beside `open.go` that reads a worktree's `.git` file and `commondir` and gives the two paths with their modes, where the real-git tests of the package live; `workitem.mainCheckout` (`repo.go`) parses the same file for another purpose and stays as it is.
- `flai/internal/harness/container.go`, `container_test.go` (design: S-0508's adapter): `Start` adds the two mounts to the `docker run` argv, and the argv test asserts the common dir is read-only and nothing under `.git` but the story's registration is writable.
- `scripts/agent-container.sh`, `docs/operators/runbooks/agent-container.md`, `scripts/README.md` (design, layout): the hand start mounts the two paths by default and `FLAI_AGENT_GIT_MOUNT` goes, with the runbook's mount and variable tables and the scripts index saying so.
- `docs/users/flai.md`, `design/system/agent-adapters.md`, `design/tech/docker.md` (co-change, design): the container's mounts as S-0508 documents them, with the two git mounts and that commits are the host's.
- `design/issues/I-0215-….md`, `design/issues/summary.md` (criterion 2): `flai issue close` writes both; `summary.md` is in the manifest's shared paths and holds nothing (ADR-0096).

`after: S-0508`: the mounts go into the container adapter's `Start`, in `flai/internal/harness/container.go`, which S-0508 writes, so this story needs S-0508's change on main; the issue itself names the remediation as the adapter's. Overlaps on 2026-10-28: S-0521, in progress, touches `scripts/agent-container.sh` and `docs/operators/runbooks/agent-container.md` (T-2209, done) and `docs/users/flai.md`; S-0414 and S-0502, in progress, touch `docs/users/flai.md`; this story is held in ready until they move to review (ADR-0096), and waits on S-0508 anyway. S-0507 (T-2174, backlog) touches the script too; whichever starts second syncs over the other.

Forecast. `flai forecast S-0517` gives 4m (median 119 s per unit of size over 13 done remediation stories on claude-opus-5-5 in the small band; size 2: 2 criteria, 0 touches) and delivery 2026-10-29T01:30Z, 30th in the pull order. Raised to 1h15m: the size flai measured had no touches, and with 14 the band's median alone gives about 30m; the story writes an ADR, a helper with a real-git test that reproduces the issue, changes the adapter's argv and its test and the script, checks the fix by hand inside `flai-agent-claude-code:local` (`scripts/agent-container.sh S-nnnn -- git status`, which needs the image built), and changes four documents. The delivery written is flai's play-out start, 2026-10-29T01:30Z less 4m, plus 1h15m: 2026-10-29T02:41Z.

Cost of delay. `flai cod S-0517` gives 37.50 USD a week: 15m of time lost per 168h cycle (the input flai set from I-0215) at 150 USD an hour, one cycle a week. Written as flai gives it: the issue has one instance, and nothing argues for another input.
