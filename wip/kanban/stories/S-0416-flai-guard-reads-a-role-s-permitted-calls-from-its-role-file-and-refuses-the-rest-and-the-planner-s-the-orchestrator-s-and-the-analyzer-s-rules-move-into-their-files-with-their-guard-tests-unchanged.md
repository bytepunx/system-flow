---
id: S-0416
type: story
nature: feature
title: flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:39:46Z
updated: 2026-10-21T09:37:18Z
transitions: []
tags: [cli, template]
topics: [cli, planning, orchestration, template]
touches: [".claude/settings.json", design/roles/analyzer.md, design/roles/orchestrator.md, design/roles/planner.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, flai/cmd/guard.go, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/guard/shared.go, template/CHANGELOG.md, template/root/.claude/settings.json, template/root/design/roles/analyzer.md, template/root/design/roles/orchestrator.md, template/root/design/roles/planner.md, docs/users/flai-reference.md]
after: [S-0414]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h
  delivery: 2026-10-21T14:56:00Z
  basis: "Its own forecast of 1h; 12th in the pull order with an in-progress limit of 5, behind S-0460, S-0461, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414 and S-0415."
  by: flai
  at: 2026-10-21T09:37:18Z
---
# S-0416 flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged

## Goal

The guard holds a role a project defines as it holds the three built-in roles today. With `FLAI_ROLE` naming a defined role, `flai guard` reads the role's file for the MCP tools, flai commands, and file edits it may make, each under the conditions the file names, and for the permissions that gate them, read from the manifest at each call; it refuses the rest, naming the role and the rule. The planner's, the orchestrator's, and the analyzer's rules move from Go into their role files, so one path decides for every role, and their guard tests pass unchanged.

## Acceptance criteria

- [ ] With `FLAI_ROLE` naming a defined role, `flai guard` decides the session's own calls from its role file: the reads every sub-agent has pass; each MCP tool, flai command, and file edit the file lists passes under its conditions and, when it names one, while the operator's permission for it is on in the manifest; every other write is refused with the role's name, the rule broken, and the advice to ask the operator with `thread_open`.
- [ ] The planner's, the orchestrator's, and the analyzer's rules are read from `design/roles/planner.md`, `orchestrator.md`, and `analyzer.md`, and `flai/internal/guard/guard_test.go`, `orchestrate_test.go`, `thread_test.go`, `authors_test.go`, `flai/cmd/guard_test.go`, `guard_permissions_change_test.go`, and `orchestrate_permissions_test.go` pass unchanged.
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

- **Touches.** The story declared none. Every touch comes from its six tasks.
  - **From layout:** a new `flai/internal/guard/roles.go` with its test; `call.go`, `guard.go`, `authors.go` with its test, and `shared.go`; `flai/cmd/guard.go`; the three role files in both copies; and `.claude/settings.json` in both copies.
  - **From design:** `strategic-agents.md`, `flai-cli.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` for criterion 6.
  - **Co-change, added:** `docs/users/flai-reference.md` (24%), regenerated from `flai guard`'s help, on the first task.
  - **Co-change, not added:** the operator docs, the dashboard's design, and `work-management.md`. The guard's tests are named in criterion 2 to pass unchanged, so they are not touched.
  - No folder touch. `.claude/settings.json` is a protected path: its task makes the write through `permission_prompt`.
- **Forecast.** Raised from flai's 43m to 1h. The size counts files, and `guard.go` is one file of about 1200 lines whose rules move into data with every refusal text kept. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
