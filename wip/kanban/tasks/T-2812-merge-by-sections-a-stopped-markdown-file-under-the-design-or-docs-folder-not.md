---
id: T-2812
type: task
nature: improvement
title: Merge by sections a stopped markdown file under the design or docs folder, not only one in claims.shared
status: done
parent: S-0637
owner: alex
created: 2026-11-09T22:09:19Z
updated: 2026-11-10T08:06:22Z
transitions:
  - to: ready
    at: 2026-11-10T07:59:51Z
    by: agent-S-0637
  - to: in-progress
    at: 2026-11-10T07:59:52Z
    by: agent-S-0637
  - to: done
    at: 2026-11-10T08:06:22Z
    by: agent-S-0637
stream: S-0637
tags: [flai, go, sync]
touches: [flai/internal/issues/generated.go, flai/internal/issues/generated_test.go]
after: [T-2811]
usage:
  source: log
  seconds: 390
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 128
      cache_read: 2327948
      cache_write: 8674
      cost: 0.8598
---
# T-2812 Merge by sections a stopped markdown file under the design or docs folder, not only one in claims.shared

## Work

In `flai/internal/issues/generated.go`, widen the matcher `FilesIn` gives `storygit.Files.Shared`: today `sharedMarkdown` takes a `.md` path only when `claims.Match` finds it in `claims.shared`. Take it also when it lies under the manifest's design or docs folder, relative to the repository root as git names it, read through the repository's layout rather than the literal names. Keep the issue files and the generated files out, as now, and every path under the wip folder. Rename the function and its comment to say what it now matches, citing the new ADR beside ADR-0175.

Every rebase flai runs takes its files from `FilesIn`: `flai stream sync`, `flai task done`'s sync, `flai accept`'s rebase, and the merge queue's batch, so this one change reaches all four.

Update `TestFilesInMergesTheSharedMarkdownFilesBySections` in `flai/internal/issues/generated_test.go`: `design/system/flaiover-dashboard.md` is now merged, and add a docs file outside the shared list, a wip work item and narrative, the repository's `README.md`, and a markdown file under `template/` that is not shared, each not merged. Add a case with renamed layout folders.

Waits for T-2811, whose ADR decides the scope this matcher implements.

## Done when

- The matcher takes the design and docs folders' markdown files and no wip, issue, generated, or other path.
- `flai test flai/internal/issues/generated.go flai/internal/issues/generated_test.go` passes.

## Notes
