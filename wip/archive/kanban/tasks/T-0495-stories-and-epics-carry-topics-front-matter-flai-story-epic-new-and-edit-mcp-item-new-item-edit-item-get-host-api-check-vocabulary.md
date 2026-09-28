---
id: T-0495
type: task
nature: feature
title: "Stories and epics carry topics: front matter, flai story/epic new and edit, MCP item_new/item_edit/item_get, host API, check vocabulary"
status: done
parent: S-0135
owner: alex
created: 2026-09-28T22:41:23Z
updated: 2026-09-28T22:45:40Z
transitions:
  - to: ready
    at: 2026-09-28T22:41:39Z
    by: agent-S-0135
  - to: in-progress
    at: 2026-09-28T22:41:40Z
    by: agent-S-0135
  - to: done
    at: 2026-09-28T22:45:40Z
    by: agent-S-0135
stream: S-0135
tags: []
touches: [flai/internal/workitem, flai/cmd, flai/internal/itemedit, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/check, flai/internal/topics]
---

# T-0495 Stories and epics carry topics: front matter, flai story/epic new and edit, MCP item_new/item_edit/item_get, host API, check vocabulary

## Work

Add an optional `topics:` list to `workitem.Item` for stories and epics: parsed strictly, validated as one word each (the rule `topics.Valid` applies), written after `tags` by `Marshal`, refused on tasks. Set it with `--topics` on `flai story new` and `flai epic new`, `--topics`/`--clear-topics` on `flai edit`, `topics` on MCP `item_new` and `item_edit`, shown by `item_get` and `flai edit --show`, and passed through the host API's `item.new` and `item.edit`. `flai check`'s topic vocabulary gains the topics stories and epics declare.

## Done when

Behavior tests cover create, edit, clear, a task refused, an invalid word refused, MCP and host API plumbing, and the vocabulary; the round-trip test is green.

## Notes
