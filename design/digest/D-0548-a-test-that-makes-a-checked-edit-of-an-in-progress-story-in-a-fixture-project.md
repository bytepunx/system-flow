---
id: D-0548
title: A test that makes a checked edit of an in-progress story in a fixture project opens the story's narrative first, since flai's checked edit refuses an in-progress story with none (narrative.missing)
kind: decision
sources: [S-0622, T-2764]
topics: [testing, go]
paths: [flai/cmd/orchestrate_permissions_test.go]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0548 A test that makes a checked edit of an in-progress story in a fixture project opens the story's narrative first, since flai's checked edit refuses an in-progress story with none (narrative.missing)

S-0622's T-2764 wrote a test in `flai/cmd/orchestrate_permissions_test.go` that edits a story in a fixture project. From S-0622's Decisions (`wip/archive/agents/S-0622.md`):

```text
T-2764's test opens S-0007's narrative first, since flai's checked edit refuses an in-progress story with none (`narrative.missing`).
```

The checked edit runs `flai check` on what it wrote and restores the item when the check refuses it, and `narrative.missing` is an error for a story in progress. So a later test whose fixture puts a story in progress and then edits it, through `item_edit`, `flai edit`, `flai touches`, or the host API, opens the story's narrative in the fixture before the edit, or the edit is refused for a reason the test does not mean to test.

Sources: S-0622 (T-2764).
