---
id: S-0207
type: story
nature: improvement
title: The conventions and prime know the planner, orchestrator, and analyzer as roles
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-02T11:54:40Z
transitions: []
tags: [flai, template]
touches: [design/conventions/, template/root/design/conventions, flai/internal/conventions, flai/internal/context, flai/cmd/prime.go]
after: [S-0196]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0207 The conventions and prime know the planner, orchestrator, and analyzer as roles

## Goal

ADR-0068 makes a convention's `roles` the list of agents that read it. The planner, orchestrator, and analyzer are new agent kinds that prime without a story, and each needs the conventions that apply to it: safety and communication for all, decisions and work management for the orchestrator, documentation and continuous improvement for the analyzer, and a short convention of their own saying what each may and may not do.

## Acceptance criteria
- [ ] `plan`, `orchestrate`, and `analyze` are roles `flai check` accepts and `flai prime --role` primes; `flai prime --role plan --epic E-nnnn|--story S-nnnn`, `--role orchestrate`, and `--role analyze` return packs: the conventions whose roles are empty or list the role, what the item names (for the planner), and briefs of the design the role needs
- [ ] A baseline convention `strategic-agents.md` (here and in the template) says what each kind does, what it never does (the planner never moves items past backlog; the orchestrator never edits code; the analyzer never authors stories), how each logs its activity, and how each asks the operator
- [ ] The baseline's `roles` name the new roles where they apply, as the designer confirms on a thread
- [ ] `design/system/conventions.md` and the user guide describe the roles; tests cover each role's pack

## Tasks

## Notes
