---
id: S-0207
type: story
nature: improvement
title: The conventions and prime know the planner, orchestrator, and analyzer as roles
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-03T17:49:39Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:36Z
    by: alex
  - to: in-progress
    at: 2026-10-03T07:06:31Z
    by: agent-S-0207
  - to: review
    at: 2026-10-03T07:38:32Z
    by: agent-S-0207
  - to: done
    at: 2026-10-03T17:49:39Z
    by: alex
tags: [flai, template]
touches: [design/conventions, template/root/design/conventions, flai/internal/conventions, flai/internal/context, flai/cmd/prime.go, flai/internal/mcpserver, docs/operators/settings.md, docs/users/flai-reference.md, template/CHANGELOG.md, template/template.yaml, design/system/conventions.md, design/system/agent-context.md, design/system/flai-cli.md, design/system/metrics.md, design/adrs, docs/users/flai.md, docs/users/conventions.md, design/issues]
after: [S-0196]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1957
  models:
    - model: claude-haiku-4-5-20251001
      input: 266
      output: 8184
      cache_read: 1439085
      cache_write: 65297
      cost: 0.2667
    - model: claude-opus-5-5
      input: 400
      output: 120772
      cache_read: 23562581
      cache_write: 479333
      cost: 10.1479
    - model: claude-sonnet-5
      input: 92
      output: 20833
      cache_read: 3707154
      cache_write: 152784
      cost: 1.3319
---
# S-0207 The conventions and prime know the planner, orchestrator, and analyzer as roles

## Goal

ADR-0068 makes a convention's `roles` the list of agents that read it. The planner, orchestrator, and analyzer are new agent kinds that prime without a story, and each needs the conventions that apply to it: safety and communication for all, decisions and work management for the orchestrator, documentation and continuous improvement for the analyzer, and a short convention of their own saying what each may and may not do.

## Acceptance criteria
- [x] `plan`, `orchestrate`, and `analyze` are roles `flai check` accepts and `flai prime --role` primes; `flai prime --role plan --epic E-nnnn|--story S-nnnn`, `--role orchestrate`, and `--role analyze` return packs: the conventions whose roles are empty or list the role, what the item names (for the planner), and briefs of the design the role needs
- [x] A baseline convention `strategic-agents.md` (here and in the template) says what each kind does, what it never does (the planner never moves items past backlog; the orchestrator never edits code; the analyzer never authors stories), how each logs its activity, and how each asks the operator
- [x] The baseline's `roles` name the new roles where they apply, as the designer confirms on a thread
- [x] `design/system/conventions.md` and the user guide describe the roles; tests cover each role's pack

## Tasks
- T-0746 flai check and flai prime know the plan, orchestrate, and analyze roles
- T-0747 A baseline strategic-agents.md convention, and the baseline's roles name the new roles
- T-0751 The design, an ADR, and the user guide describe the strategic roles and their packs

## Notes

- Close-out (2026-10-03): the lint, the whole suite, and the template smoke pass. `flai check --strict` stopped on 10 warnings, none from S-0207's changes. Nine are `wip.overlap` with S-0200 and its task T-0752, in progress side by side, which TH-0084 settles: whichever story is accepted second renumbers its ADR and keeps both changes. One is `story.unaccepted` on S-0173's unmerged branch. As TH-0056 decided, this is recorded as a bump of I-0057, and S-0244 in the backlog covers it.
