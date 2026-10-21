---
id: T-1844
type: task
nature: feature
title: flai stats prints the digest's weeks under a digest heading and carries them in --json
status: backlog
parent: S-0434
owner: alex
created: 2026-10-16T06:08:30Z
updated: 2026-10-20T19:01:12Z
transitions: []
stream: S-0434
tags: [cli]
touches: [flai/cmd/stats.go, flai/cmd/check_stats_test.go, flai/cmd/stats_test.go, docs/users/flai-reference.md]
after: [T-1843]
---
# T-1844 flai stats prints the digest's weeks under a digest heading and carries them in --json

## Work

In `flai/cmd/stats.go`, add `printDigest`: a `digest` heading, the entries by state on one line, then one line per week, `W42  threads 12 · stories 8 · 1.5 per story · 3 asked again`, the part for `per_story` left out when absent. `--json` carries `digest` as the report has it. Name the digest section in the command's help, as S-0404 named the merge queue's, and regenerate `docs/users/flai-reference.md` from it.

Waits for T-1843, whose report it prints. Runs beside the chart task.

## Done when

`check_stats_test.go` or `stats_test.go`, whichever holds the section tests by then, covers the text and JSON on a fixture with repeats and a week with no completed story; `docs/users/flai-reference.md` is regenerated; `flai test flai/cmd` passes.

## Notes

`docs/users/flai.md`'s description is the docs task's.
