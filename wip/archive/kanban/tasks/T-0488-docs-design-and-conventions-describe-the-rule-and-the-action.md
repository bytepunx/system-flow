---
id: T-0488
type: task
nature: remediation
title: Docs, design, and conventions describe the rule and the action
status: done
parent: S-0140
owner: alex
created: 2026-09-26T21:11:26Z
updated: 2026-09-26T21:27:52Z
transitions:
  - to: ready
    at: 2026-09-26T21:11:32Z
    by: agent-S-0140
  - to: in-progress
    at: 2026-09-26T21:22:01Z
    by: agent-S-0140
  - to: done
    at: 2026-09-26T21:27:52Z
    by: agent-S-0140
stream: S-0140
tags: []
touches: [docs, design/system, design/conventions, template/root/design/conventions]
---
# T-0488 Docs, design, and conventions describe the rule and the action

## Work

- `docs/users/flai.md`, `docs/operators`, `design/system/workflow.md` and `flai-cli.md`, and the work-management convention (template first) say that a story goes to review only with its worktree committed, and how the operator has outstanding work committed.

## Done when

- [x] Docs, design, and conventions are updated, `make flai-reference` is current, lint-md is clean.

## Notes
