---
id: TH-0412
title: "S-0349 needs a sync onto S-0235 before I accept it: both change the same docs"
anchor:
  path: wip/kanban/stories/S-0349-a-story-s-agent-names-a-provider-from-the-host-s-providers-map-and-a-project-s-manifest-may-override-the-entry-s-api-base-url-and-models.md
  item: S-0349
status: open
participants: [orchestrator]
created: 2026-10-12T16:41:48Z
updated: 2026-10-12T16:41:48Z
---

# TH-0412 S-0349 needs a sync onto S-0235 before I accept it: both change the same docs

On wip/kanban/stories/S-0349-a-story-s-agent-names-a-provider-from-the-host-s-providers-map-and-a-project-s-manifest-may-override-the-entry-s-api-base-url-and-models.md.

## Entries

### 2026-10-12T16:41:48Z orchestrator
**Recommendation:** move S-0349 back to in-progress, so its agent runs `flai stream sync`, resolves the docs, and closes out again. I accept it once it is back in review, synced and verified at its new head.

```sh
flai move S-0349 in-progress --reason "sync onto S-0235 before acceptance"
```

## Why

S-0349 moved to review at 16:41:00, and its verify passed at its head, 11cdc21a. I accepted S-0235 seven seconds later. S-0235 changed files S-0349 also changes:

- `docs/operators/settings.md`
- `design/system/flai-cli.md`
- `docs/users/flai.md`
- `docs/users/flai-reference.md`

The verifier on S-0235 had already flagged the `docs/operators/settings.md` overlap as a likely conflict.

`flai accept` would rebase S-0349 onto main. If that rebase stops on a conflict, flai guard lets me run only git's reads, so I could not abort it (as on TH-0397). So I am not running the acceptance into it.

Also: until the host runs flai 1.41.x (TH-0411), MCP `verify` cannot run S-0349's tiers, though its close-out can.

S-0349 stays in review until then.
