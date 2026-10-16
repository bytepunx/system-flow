---
id: T-1684
type: task
nature: feature
title: flai check warns when a project has neither AGENTS.md nor CLAUDE.md at its root
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:07Z
updated: 2026-10-16T05:14:07Z
transitions: []
stream: S-0406
tags: [cli]
touches: [flai/internal/check/instructions.go, flai/internal/check/instructions_test.go, flai/internal/check/check.go]
---
# T-1684 flai check warns when a project has neither AGENTS.md nor CLAUDE.md at its root

## Work

- Add a rule in `flai/internal/check/instructions.go`, registered in `check.go`: a warning when the project's root has neither `AGENTS.md` nor `CLAUDE.md`, naming both and saying which harnesses read which.
- Test a root with each file, with both, and with neither.
- Waits for nothing.

## Done when

`flai check` warns on a project with neither file, says nothing on one with either, and the rule's tests pass.

## Notes
