---
id: S-0518
type: story
nature: remediation
title: flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start
status: backlog
owner: alex
created: 2026-10-28T12:52:56Z
updated: 2026-10-29T04:06:06Z
transitions: []
tags: [cli, template]
topics: [harness]
touches: [".claude/settings.json", template/root/.claude/settings.json, template/CHANGELOG.md, scripts/flai.sh, scripts/README.md, flai/cmd/guard_test.go, flai/internal/harness/harness_test.go, flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, design/system/agent-adapters.md, design/adrs/README.md, design/adrs, design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md, design/issues/summary.md]
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
      seconds: 801
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 70
          output: 1042
          cache_read: 19646500
          cache_write: 23639
          cost: 4.3183
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:52:56Z
  value: 37.5
  by: planner-S-0518
  at: 2026-10-28T13:33:45Z
forecast:
  duration: 1h15m
  delivery: 2026-10-29T15:35:00Z
  basis: "Its own forecast of 1h15m; 23rd in the pull order with an in-progress limit of 5, behind S-0417, S-0539, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0508, S-0509, S-0510, S-0511, S-0512, S-0515, S-0516 and S-0517."
  by: flai
  at: 2026-10-29T04:06:06Z
---
# S-0518 flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start

## Goal

This story remediates [I-0216](../../../design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md), "flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0216 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0216 is closed with `flai issue close I-0216 --reason` saying what fixed it

## Tasks
- T-2223 An ADR refining ADR-0060 and ADR-0203 records that the guard hook refuses a call when flai guard cannot start, and that inside the container the guard is the image's flai through the project's hook
- T-2224 scripts/flai.sh runs the installed flai, with a warning, when go is absent and bin/flai is not built, so this repository's guard hook reaches the image's flai inside the container
- T-2225 The guard hooks in .claude/settings.json and the template's refuse a call when flai guard cannot start, and a test reproduces I-0216 against both files' commands
- T-2226 The guard's docs say the hook refuses a call when the guard cannot start, flai guard's help names it, the fix is checked once inside the image, and I-0216 is closed

## Notes

Cost of delay inputs set by flai from I-0216. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:30:51Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0518 on 2026-10-28 (TH-0688 names the tasks, their layers, and the assumptions).

The cause, from the code. The guard itself exits only 0 (allow) or 2 (refuse): `flai/cmd/guard.go` passes nothing else, and no code in flai handles a guard that cannot start, since that case is outside the binary. The hook's shell line decides: `out=$(<guard> 2>&1); [ $? -eq 2 ] || exit 0; ...` in both `.claude/settings.json` files allows anything that is not exactly 2. This repository's `<guard>` is `"$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard`, which runs `scripts/flai-build.sh` (`go build`) when `bin/flai` is missing or stale; a fresh worktree has no `bin/`, the image `flai-agent-claude-code:local` has no `go`, so the build exits 127 and the hook passes every call. The template's hook runs the installed `flai`, which the image carries on `PATH`, so a project made from the template is guarded inside the container today; the fail-open shape is the template's too, and a host without `flai` on `PATH` is unguarded in the same way.

The remediation proposed, in two parts, recorded by an ADR refining ADR-0060 and ADR-0203. First, the hook refuses the call, exit 2 with `flai guard did not run (exit N): <output>`, when the guard command exits 126 or 127 (not found or not executable), and keeps allowing on exit 0 and on any other exit, as ADR-0060 decided: a fail-closed on every non-zero exit would lock a story's agent of this repository out of `Bash` and `Edit` whenever the tree's build fails (both hooks run the guard), and would refuse everything to an older flai without `guard`. Second, inside the container the guard is the image's flai through the project's hook: the template's hook names `flai` on `PATH` already, and `scripts/flai.sh` runs the installed flai, with a warning, when `go` is absent and `bin/flai` is not built. pilot (S-0507) passes no hook of its own: hooks from several settings sources add up, so a `--settings` file of flai's would run beside the project's hook rather than replace it, and excluding the project's settings would drop its permissions too. S-0510's layer gives this repository's hook its `go` later, which makes `scripts/flai.sh` build again as on the host.

Touches. `flai touches suggest S-0518` names no path on its own, since the story declares none; seeded with the two settings files, `guard_test.go`, `scripts/flai.sh`, and `agent-adapters.md` (48 of 2478 commits) it lists `flai/cmd/guard.go` (31%), `flai/internal/guard/guard.go` (31%), `guard/guard_test.go` (27%), `docs/users/flai-reference.md` (25%), `flai-cli.md` (15%), `docs/operators/index.md` (12%), `design/adrs/README.md` (10%), `docs/users/flai.md` (10%), `docs/operators/settings.md` (8%), `template/CHANGELOG.md` (8%), `template/template.yaml` (8%), `design/issues/summary.md` (6%), `strategic-agents.md` (6%), and `harness_test.go` (6%). Kept from it (co-change): `guard.go` (its help text describes what the hook does with its exit codes), `flai-reference.md` (generated from that help by `make flai-reference`), `design/adrs/README.md`, `docs/users/flai.md`, `docs/operators/settings.md`, `template/CHANGELOG.md`, `design/issues/summary.md`, `harness_test.go` (its template test requires the hook's `flai guard 2>&1` substring). Left out: `guard/guard.go` and `guard/guard_test.go` (the policy and the readers do not change, only the shell line around the guard), `flai-cli.md` and `docs/operators/index.md` (neither describes the hook's exit codes), `template/template.yaml` (bumped at release), `strategic-agents.md` (its guard section is per role, not the hook's shape). The rest are the tasks' (layout):

- `.claude/settings.json`, `template/root/.claude/settings.json` (layout): the two `PreToolUse` commands. T-2225. Both are protected paths: the agent's writes go through `permission_prompt`, which asks the operator on a thread (ADR-0086), and the story is accepted by the operator alone (ADR-0106).
- `scripts/flai.sh`, `scripts/README.md` (layout): the fallback to the installed flai and the index's line on it. T-2224.
- `flai/cmd/guard_test.go` (layout): `checkGuardHooks` already reads both settings files and runs the edit hook's condition under `sh`; the test that reproduces I-0216 runs each file's Bash hook command with a fake guard and with none. T-2225.
- `design/system/agent-adapters.md` (design): § The guard says the guard fails open on input it cannot read; it gains what the hook does when the guard cannot start. T-2226.
- `design/issues/I-0216-….md` (criterion 2): `flai issue close` writes it and `summary.md`; `summary.md` is in the manifest's shared paths and holds nothing (ADR-0096). T-2226.
- The folder `design/adrs` (layout, folder kept): `flai adr new` in T-2223 names the ADR's file by its number and title, which no task can name; the folder is kept for that file alone, is in the manifest's shared paths, so it holds no ready story, and `flai task done` widens the touches with the committed path (ADR-0144).

Overlaps on 2026-10-28: `docs/users/flai.md` is touched by S-0414, S-0502, and S-0521 (in progress), `docs/operators/settings.md` by S-0414 and S-0502, and `harness_test.go` by S-0502, so the story is held until they move to review or are accepted; S-0416 (ready, E-0024) touches the two settings files, `guard.go`, `guard_test.go`, `harness_test.go`, `docs/users/flai.md`, `agent-adapters.md`, and `template/CHANGELOG.md`, so whichever of the two starts first holds the other, and the pull order settles which. No `after`: the story needs no other story's change on main, and S-0507 and S-0508 are told on TH-0688 what the container relies on.

Forecast. `flai forecast S-0518` gives 4m (median 119 s per unit over 13 done remediation stories on claude-opus-5-5 in the small band; size 2: 2 criteria, 0 touches) and delivery 2026-10-29T01:41Z, 31st in the pull order. Raised to 1h15m: the size flai measured had no touches, and with 16 the band's median alone gives about 36m; the story writes an ADR, a shell-level test matrix over two settings files and a copy of the scripts tree, two protected writes that each wait on the operator's answer on a thread and are made again after it, the generated reference, and a check by hand inside `flai-agent-claude-code:local`, which needs `make agent-images` built from the tree. The delivery written is flai's play-out start, 2026-10-29T01:41Z less 4m, plus 1h15m: 2026-10-29T02:52Z.

Cost of delay. `flai cod S-0518` gives 37.50 USD a week: 15m lost per 168h cycle at 150 USD an hour, from the input flai set from I-0216. Written as flai gives it: the issue has one instance, and the input is the operator's to raise if a guard silently off in every container run of this repository weighs more than its one occurrence.
