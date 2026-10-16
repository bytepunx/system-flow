---
id: T-1835
type: task
nature: feature
title: The curator's prompt, definition, and guard take in absorption and compaction
status: backlog
parent: S-0432
owner: alex
created: 2026-10-16T06:07:08Z
updated: 2026-10-16T06:07:08Z
transitions: []
stream: S-0432
tags: [cli, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/guard/guard.go, flai/internal/guard/curate_test.go, ".claude/agents/curator.md", template/root/.claude/agents/curator.md]
after: [T-1833]
---
# T-1835 The curator's prompt, definition, and guard take in absorption and compaction

## Work

- `flai/internal/harness/harness.go`: `curatePrompt` adds, after taking in the harvest: for each active entry, look with `doc_search` on its title for a design section, convention, or ADR that now records it, read it with `doc_get`, and mark the entry absorbed with `flai digest absorb --by <path> --heading <h>` when it does; then run `flai digest compact`; name what it absorbed and compacted in its summary.
- `flai/internal/guard/guard.go`: `cliCurates` passes `digest absorb`, `cited`, and `compact`.
- `.claude/agents/curator.md` (template and here) says the same in short.

Waits for T-1833, whose commands the prompt names.

## Done when

`harness_test.go` and `curate_test.go` cover the new wording and the passes; `flai test flai/internal/harness flai/internal/guard` passes.

## Notes

The `.claude/agents/curator.md` write here goes through `permission_prompt`; the story's agent makes it itself and waits for the operator's answer.
