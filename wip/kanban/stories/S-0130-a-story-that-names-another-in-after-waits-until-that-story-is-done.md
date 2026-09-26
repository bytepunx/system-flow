---
id: S-0130
type: story
nature: feature
title: "A story that names another in after: waits until that story is done"
status: backlog
parent: E-0009
owner: alex
created: 2026-09-26T07:59:22Z
updated: 2026-09-26T07:59:22Z
transitions: []
tags: [cli]
touches: [flai/internal/workitem, flai/internal/check, flai/internal/itemedit, flai/internal/serve, flai/internal/mcpserver, flai/cmd, flaiover, docs, design/system]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0130 A story that names another in after: waits until that story is done

## Goal

Add the explicit dependency of [ADR-0046](../../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md): a story may name, in `after:`, stories that must be done before it starts, for dependencies that are not about files. The hold, its reason, and the yellow card are the ones the overlap hold already has.

## Acceptance criteria

- [ ] Stories accept an optional `after: [S-nnnn, …]` front matter list, set and cleared with `flai edit --after` / `--clear-after`, MCP `item_edit`, and the dashboard's story editor.
- [ ] A ready story is held with reason `held (after): waits for S-nnnn (<status>); starts when S-nnnn is done` while any story it names is not done; a cancelled one keeps the hold and the reason says it was cancelled.
- [ ] `flai check` reports an `after:` entry that names no story, names the story itself, or forms a cycle.
- [ ] The launcher, `wait_for_work`, `flai board`, `inbox`, and the card treat an `after` hold as they treat an overlap hold.
- [ ] `design/system/work-hierarchy.md`, `workflow.md`, `flai-cli.md`, and `docs/users/flai.md` describe the field; `make test`, lint, and flaiover's tests pass.

## Tasks

## Notes

- Pulled after the two hold stories under E-0009, which it builds on.
