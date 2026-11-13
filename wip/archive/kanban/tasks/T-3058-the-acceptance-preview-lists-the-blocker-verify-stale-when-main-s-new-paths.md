---
id: T-3058
type: task
nature: remediation
title: The acceptance preview lists the blocker verify_stale when main's new paths send back a tier the story's last verify ran
status: done
parent: S-0683
owner: alex
created: 2026-11-10T10:35:22Z
updated: 2026-11-12T16:18:28Z
transitions:
  - to: ready
    at: 2026-11-12T16:01:23Z
    by: agent-S-0683
  - to: in-progress
    at: 2026-11-12T16:01:23Z
    by: agent-S-0683
  - to: done
    at: 2026-11-12T16:18:28Z
    by: agent-S-0683
stream: S-0683
tags: [flai, acceptance]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, design/adrs/0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md]
after: [T-3051]
usage:
  source: log
  seconds: 1025
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 66
      output: 18997
      cache_read: 6325610
      cache_write: 119326
      cost: 2.4984
---
# T-3058 The acceptance preview lists the blocker verify_stale when main's new paths send back a tier the story's last verify ran

## Work

It waits for T-3051, whose ADR decides the gate this task builds.

In `flai/internal/verify/story.go`, export a function for a story's last passing report and main's current head. It answers the tiers that main's paths since the report's base would send back, and the paths that reach them. It reuses `sentBack`'s selection (`SelectStory` with `{importers}`, and the tier that covers each) and leaves out `wip/` paths, as ADR-0135 and the close-out's sync check do. It answers nothing when the report's base is main's head, or when main's new paths select no tier the report ran.

In `flai/internal/preview/accept.go`, add a gate beside `queueTiersGate`. It sets the blocker `verify_stale`, with its message, on every acceptance of a story with a branch: the operator's, the board's, and the orchestrator's, with the queue on or off. Leave out the queue's landing (`opts.Batch`), a resumed acceptance, and a story done already. Give the result the fields `verify_stale_tiers` and `verify_stale_paths` for `--json` and the dashboard. The message names the tiers and the paths, and says to ask for the close-out again with `flai verify <id> --job`, or the MCP tool `verify`. Table tests cover each case in both files.

## Done when

- Unit tests show the blocker for a report whose base is behind a main that changed a package the story's tested package imports.
- Unit tests show no blocker when main is at the report's base, when main changed only `wip/` paths, or when the paths select no tier the report ran, and none under `--batch`.
- `flai test` passes on the four files.

## Notes
