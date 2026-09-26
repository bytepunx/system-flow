---
id: T-0481
type: task
nature: feature
title: Acceptance records an overlap notice for each open story whose claim covers a changed path
status: done
parent: S-0132
owner: alex
created: 2026-09-26T18:17:26Z
updated: 2026-09-26T18:19:58Z
transitions:
  - to: ready
    at: 2026-09-26T18:17:40Z
    by: agent-S-0132
  - to: in-progress
    at: 2026-09-26T18:17:41Z
    by: agent-S-0132
  - to: done
    at: 2026-09-26T18:19:58Z
    by: agent-S-0132
stream: S-0132
tags: []
touches: [flai/cmd, flai/internal/itemedit]
---
# T-0481 Acceptance records an overlap notice for each open story whose claim covers a changed path

## Work

After the story branch is merged, take the paths main gained (`git diff --name-only` from main's head before the merge to after it). For every other story in progress or in review, compare them with its claim (`workitem.Holds.Claim`, `PathsOverlap`); a story with no touches overlaps everything. Append a notice naming the open story, the accepted story, who accepted, and the overlapping paths to a notice log of its own under `.flai-cache` (not `edits.jsonl`, which an older flai reads as edits). `flai accept` prints whom it told, and returns it with `--json`.

## Done when

- A test accepts a story beside an open story that overlaps it and one that does not: only the first gets a notice, naming the overlapping paths.
- `go test -short ./cmd/... ./internal/itemedit/...` passes.

## Notes
