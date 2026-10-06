---
id: T-1033
type: task
nature: improvement
title: The users' and operators' guides and the dashboard design describe the new hold rules, the shared paths, and how to edit them
status: done
parent: S-0295
owner: alex
created: 2026-10-06T12:16:47Z
updated: 2026-10-06T19:19:04Z
transitions:
  - to: ready
    at: 2026-10-06T19:14:48Z
    by: agent-S-0295
  - to: in-progress
    at: 2026-10-06T19:14:48Z
    by: agent-S-0295
  - to: done
    at: 2026-10-06T19:19:04Z
    by: agent-S-0295
stream: S-0295
tags: [flai, flaiover]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/system/agent-coordination.md, template/CHANGELOG.md, flai/internal/hostapi/settings.go]
after: [T-1026, T-1028, T-1029, T-1030, T-1031, T-1032]
usage:
  source: log
  seconds: 256
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 109
      output: 45162
      cache_read: 6085168
      cache_write: 178541
      cost: 3.1023
---
# T-1033 The users' and operators' guides and the dashboard design describe the new hold rules, the shared paths, and how to edit them

## Work

Bring the outward documentation in line with what the story built:

- `docs/users/flai.md`: in the section on holds and touches, say that a story in review no longer holds, that tasks narrow a folder touch, and that shared paths never hold. Add a section for `flai shared` and the MCP tools.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference`.
- `docs/operators/settings.md` § Project manifest: add the new key, its default, and its glob dialect, with the CLI, HTTP, MCP, and dashboard ways to change it.
- `docs/users/flaiover.md` § Settings and `design/system/flaiover-dashboard.md`: describe the shared paths section and the host API methods.
- `template/CHANGELOG.md`: add an entry for the manifest key and its default in new projects, below T-1026's entry.

This task waits for every task whose behaviour it documents: T-1028, T-1029, T-1030, T-1031, and T-1032. It also waits for T-1026, which writes `template/CHANGELOG.md` first.

## Done when

- Each document says what is above, and `flai-reference.md` matches `make flai-reference`.
- The settings doc test in `flai/cmd/settings_doc_test.go` passes with the new key listed.
- The markdown lint and `flai check --strict` pass.

## Notes
