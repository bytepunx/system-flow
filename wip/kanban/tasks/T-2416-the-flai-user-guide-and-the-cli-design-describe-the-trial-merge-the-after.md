---
id: T-2416
type: task
nature: improvement
title: The flai user guide and the CLI design describe the trial merge, the after advisory, and the new stats
status: done
parent: S-0551
owner: alex
created: 2026-10-31T07:44:56Z
updated: 2026-10-31T10:18:12Z
transitions:
  - to: ready
    at: 2026-10-31T10:07:19Z
    by: agent-S-0551
  - to: in-progress
    at: 2026-10-31T10:07:21Z
    by: agent-S-0551
  - to: done
    at: 2026-10-31T10:18:12Z
    by: agent-S-0551
stream: S-0551
tags: [flai, docs]
touches: [docs/users/flai.md, design/system/flai-cli.md, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/cmd/promote.go, flai/cmd/story_start.go, design/system/workflow.md, docs/users/flai-reference.md]
after: [T-2393, T-2395, T-2404, T-2410]
usage:
  source: log
  seconds: 651
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 131
      output: 41673
      cache_read: 8194638
      cache_write: 300546
      cost: 4.1476
---
# T-2416 The flai user guide and the CLI design describe the trial merge, the after advisory, and the new stats

## Work

Update `docs/users/flai.md`: `### Touches` and the "Held" row of `#### Starting a story` for the trial merge and the conflicting paths in a hold's reason; `### Waiting for another story` for the `after` rule and the advisory; `## Flow metrics` for the in-progress share and the ready wait. Update `design/system/flai-cli.md` where it describes the holds, `flai check`'s findings, and `flai stats`.

It waits for the tasks that change the behaviour it documents: the serve and caller wiring, the check advisory, and the stats.

## Done when

- Both documents describe the trial merge, the advisory, and the new stats as built.
- The markdown lint passes on both.

## Notes

Both files are in the manifest's shared paths, so other stories may change them meanwhile; sync before editing.
