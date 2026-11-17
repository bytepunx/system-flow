---
id: T-3558
type: task
nature: remediation
title: flai digest new and edit refuse an entry whose markdown the project's lint rejects
status: done
parent: S-0836
owner: alex
created: 2026-11-17T11:33:11Z
updated: 2026-11-17T11:42:15Z
transitions:
  - to: ready
    at: 2026-11-17T11:33:17Z
    by: agent-S-0836
  - to: in-progress
    at: 2026-11-17T11:33:17Z
    by: agent-S-0836
  - to: done
    at: 2026-11-17T11:42:15Z
    by: agent-S-0836
stream: S-0836
tags: [flai, digest]
touches: [flai/internal/digest/write.go, flai/internal/digest/write_test.go, flai/cmd/digest.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 538
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 40
      output: 251
      cache_read: 5448596
      cache_write: 17189
      cost: 2.0876
---
# T-3558 flai digest new and edit refuse an entry whose markdown the project's lint rejects

## Work

- In `flai/internal/digest/write.go`, `New` and `Edit` lint the entry they wrote with the project's markdownlint configuration, through `workitem.Repo.LintGuard`, as flai's wip writers do (ADR-0061, S-0179). Each finding the write brings is a `Problem` with rule `markdown.MDnnn`, its line, and markdownlint's message, so the write is refused as any other problem is: `New` removes the file, and `Edit` puts it back unless the finding was there before the edit.
- I-0363's cause: the curator wrote D-0592 through `flai digest new`, which checked its front matter and sources and never its markdown, and `flai check` lints only wip/, so an entry markdownlint rejects reached main and failed `TestRepositoryLintsClean`.
- A test in `write_test.go` writes D-0592's line 17 as the curator first wrote it under a configuration with MD038, and is refused with MD038; `Edit` that brings it is refused and the file put back.
- `flai/cmd/digest.go`'s help for `new` and `edit`, `docs/users/flai.md`, the generated `docs/users/flai-reference.md`, and the digest section of `design/system/agent-context.md` say so. Record the rule as an ADR.
- Waits for no task: it shares no path with T-3529, and is in layer 1.

## Done when

- `New` and `Edit` refuse an entry that brings a markdown finding, and a project with no markdownlint configuration is not linted.
- `flai test` on the changed paths passes.

## Notes
