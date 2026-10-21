---
id: S-0416
type: story
nature: feature
title: flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:39:46Z
updated: 2026-10-21T12:24:04Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, template]
touches: [".claude/settings.json", template/root/.claude/settings.json, design/workflows/roles/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/analyzer.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, design/adrs, design/adrs/README.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/guard/shared.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/harness/harness_test.go]
after: [S-0414]
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
      seconds: 1511
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 67
          output: 1162
          cache_read: 8560656
          cache_write: 20404
          cost: 1.8949
draft: true
cost_of_delay:
  value: 59.07
  by: planner-S-0416
  at: 2026-10-21T11:41:58Z
forecast:
  duration: 1h10m
  delivery: 2026-10-21T18:39:00Z
  basis: "Its own forecast of 1h10m; 11th in the pull order with an in-progress limit of 5, behind S-0461, S-0471, S-0473, S-0479, S-0482, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414 and S-0415."
  by: flai
  at: 2026-10-21T12:24:04Z
---
# S-0416 flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged

## Goal

The guard holds a role a project defines as it holds the three built-in roles today. With `FLAI_ROLE` naming a defined role, `flai guard` reads the role's file for the MCP tools, flai commands, and file edits it may make, each under the conditions the file names, and for the permissions that gate them, read from the manifest at each call; it refuses the rest, naming the role and the rule. The planner's, the orchestrator's, and the analyzer's rules move from Go into their role files, so one path decides for every role, and their guard tests pass unchanged.

## Acceptance criteria

- [ ] With `FLAI_ROLE` naming a defined role, `flai guard` decides the session's own calls from its role file: the reads every sub-agent has pass; each MCP tool, flai command, and file edit the file lists passes under its conditions and, when it names one, while the operator's permission for it is on in the manifest; every other write is refused with the role's name, the rule broken, and the advice to ask the operator with `thread_open`.
- [ ] The planner's, the orchestrator's, and the analyzer's rules are read from `design/workflows/roles/planner.md`, `orchestrator.md`, and `analyzer.md`, and `flai/internal/guard/guard_test.go`, `orchestrate_test.go`, `thread_test.go`, `authors_test.go`, `flai/cmd/guard_test.go`, `guard_permissions_change_test.go`, and `orchestrate_permissions_test.go` pass unchanged, but for the hook condition `flai/cmd/guard_test.go` holds, which changes with criterion 5.
- [ ] ADR-0143's own-name rule holds a defined role as it holds the planner and the analyzer: a `--by`, `--placed-by`, `FLAI_AGENT`, or `FLAI_SESSION` naming another is refused, its own name being the hook's `FLAI_AGENT` or the name `flai serve` gives the role's run.
- [ ] A role whose file is missing or does not parse is refused every write, naming the file and `flai check`; the guard fails closed.
- [ ] The `Edit|MultiEdit|Write|NotebookEdit` entry of `.claude/settings.json`, in the template and here, runs the guard for any `FLAI_ROLE`, not only `plan`, `orchestrate`, and `analyze`, so a defined role's file edits are held.
- [ ] `design/system/strategic-agents.md` § The guard, `design/system/flai-cli.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` say that the guard reads a role's rules from its file.

## Tasks
- T-1729 flai guard decides a defined role's own calls from its role file: tools, commands, and file edits under their conditions, gated by the manifest's permissions at each call
- T-1730 The planner's and the analyzer's guard rules are read from planner.md and analyzer.md through the role path, with their guard tests unchanged
- T-1731 The orchestrator's guard rules, its permissions and its thread rules, are read from orchestrator.md through the role path, with its permission tests unchanged
- T-1732 ADR-0143's own-name rule holds a defined role, its own name the hook's FLAI_AGENT or the name flai serve gives its run
- T-1733 The Edit, Write, and NotebookEdit hook in .claude/settings.json runs flai guard for any FLAI_ROLE, in the template and here
- T-1734 The design, the user guide, and the changelog say that flai guard reads every role's rules from its role file

## Notes

- Drafted by the planner for E-0024. Waits for the role file story; runs beside the prime story, with which it shares no file.
- The `.claude/settings.json` change is a protected write: the story's agent makes it through `permission_prompt`, which asks the owner (ADR-0086).
- Rules that are not plain allow-lists stay in Go as named conditions the role file references (the role file story's ADR): the planner's draft-only `item_new`, the moves to backlog, the analyzer's reports folder, and the orchestrator's thread and permission rules.

### Planning

- **TH-0471's answers (2026-10-21).** The plan follows them in two places.
  - **Answer 2** moved the role files from `design/roles/` to `design/workflows/roles/`, as S-0414 now says. Criterion 2 and the six role-file touches now give the new path. The tasks' titles and Work name the files without a folder, so they did not change.
  - **Answer 1:** E-0024 does not yet carry the input, `time_lost_per_cycle: 4h`.
  - Answers 3 to 6 change nothing the guard decides.
- **TH-0569's approval.** The orchestrator approved the plan and its three proposals.
  - **Criterion 2, narrowed.** Two tests assert the hook's condition word for word: `flai/cmd/guard_test.go` checks its `editGuardCondition`, and `flai/internal/harness/harness_test.go` checks the template's hook. Criterion 5 could not hold while criterion 2 held as written, so criterion 2 now excepts that one condition, which T-1733 changes.
  - **An ADR.** T-1732 records a short ADR that refines ADR-0143. ADR-0143's decision names the agents it holds, and a defined role is not among them.
  - **T-1734's Work** now covers `agent-adapters.md` and `flai guard`'s help.
- **Touches.** 27 of the 28 are files.
  - **Declared, kept:** the 15 earlier touches that are not role files.
  - **Declared, moved:** the six role files, now under `design/workflows/roles/` (answer 2).
  - **From layout, added:**
    - `flai/internal/guard/guard.go` on T-1729, which declares the `Guard` struct that the role path needs a field on.
    - `flai/internal/roles/conditions.go` and its test on T-1730 and T-1731, in case a built-in rule needs a named condition that S-0414's set lacks.
    - `flai/cmd/guard.go` and `docs/users/flai-reference.md` on T-1734, for the guard's help, which describes each built-in role's rules and the hook's condition.
  - **From the tests, added:** `flai/cmd/guard_test.go` and `flai/internal/harness/harness_test.go` on T-1733. Both pin the hook's condition.
  - **From design, added:**
    - `design/system/agent-adapters.md` on T-1734. It names the planner and the analyzer as the agents ADR-0143 holds, and it gives the hook's condition and the values of `FLAI_ROLE`.
    - `design/adrs/README.md` on T-1732, for the new ADR's entry.
  - **Folder touch kept:** `design/adrs`, on T-1732, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
  - **Co-change, not added:** the operator docs, the dashboard's design, `workflow.md`, `work-management.md`, and `flai/internal/harness/harness.go`. The guard changes none of them.
  - `.claude/settings.json` in both copies is a protected path, so T-1733 makes its write through `permission_prompt`.
- **Forecast.** flai's 1h1m is raised to 1h10m: size 34 at 107 s a unit.
  - The size counts files, and `guard.go` is one of them, at 1646 lines. The three built-in roles' rules move out of it into data, and every refusal text must stay word for word.
  - The two `conditions.go` touches may not be needed, so they overstate the size a little.
  - The delivery, 2026-10-21T17:30Z, is flai's play-out moved later by the 9 minutes of the raise.
- **Cost of delay.** 59.07 USD a week. `flai cod` refuses while E-0024 has no input, so the figure is worked out by hand from answer 1, as S-0414's and S-0415's planners did.
  - 4h lost per cycle, at `planning.hour_rate` 150 over the default 168h cycle, gives E-0024 600 USD a week.
  - S-0416 takes the share its 1h10m is of the 11h51m forecast for E-0024's twelve open stories: 70 / 711 × 600.
  - Once the input is on E-0024, `flai cod S-0416` gives the same figure while the forecasts stand.
