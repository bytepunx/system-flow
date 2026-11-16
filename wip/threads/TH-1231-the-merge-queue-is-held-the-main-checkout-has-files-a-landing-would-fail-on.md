---
id: TH-1231
title: "The merge queue is held: the main checkout has files a landing would fail on"
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-11-15T21:46:13Z
updated: 2026-11-15T21:46:13Z
---

# TH-1231 The merge queue is held: the main checkout has files a landing would fail on

On system-flow.yaml.

## Entries

### 2026-11-15T21:46:13Z flai serve
The merge queue's preflight found what a landing would fail on in the main checkout, at main's head `0835f58c5205`:

- `design/workflows/roles/orchestrator.md`, line 15, `roles.file`: `front matter does not parse: cannot unmarshal map[string]interface {} into Go struct field Role.Trigger of type string; give each field the shape ADR-0207 §2 gives it` (flai check)
- `design/workflows/roles/planner.md`, line 28, `roles.file`: `front matter does not parse: cannot unmarshal map[string]interface {} into Go struct field Role.Tools of type string; give each field the shape ADR-0207 §2 gives it` (flai check)

The queue builds no batch and runs no tier until these are fixed. It changes none of these files: they are yours. Fix them in the main checkout. The queue looks again as the main checkout changes and as main's head moves, and the look that finds main clean clears the hold, resolves this thread, and builds batches again, with nothing else asked of you.
