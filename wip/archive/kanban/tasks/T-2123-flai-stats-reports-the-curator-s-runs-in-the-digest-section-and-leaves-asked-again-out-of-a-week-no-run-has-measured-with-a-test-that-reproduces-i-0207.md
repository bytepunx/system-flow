---
id: T-2123
type: task
nature: feature
title: flai stats reports the curator's runs in the digest section and leaves asked_again out of a week no run has measured, with a test that reproduces I-0207
status: done
parent: S-0504
owner: alex
created: 2026-10-26T04:09:47Z
updated: 2026-10-26T04:25:50Z
transitions:
  - to: ready
    at: 2026-10-26T04:14:07Z
    by: agent-S-0504
  - to: in-progress
    at: 2026-10-26T04:14:08Z
    by: agent-S-0504
  - to: done
    at: 2026-10-26T04:25:50Z
    by: agent-S-0504
stream: S-0504
tags: [flai, digest, metrics]
touches: [flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flai/cmd/check_stats_test.go]
after: [T-2121]
usage:
  source: log
  seconds: 702
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 6011
      cache_read: 1560072
      cache_write: 48955
      cost: 0.7701
---
# T-2123 flai stats reports the curator's runs in the digest section and leaves asked_again out of a week no run has measured, with a test that reproduces I-0207

## Work

Build what T-2121's ADR and `design/system/metrics.md` § Digest define, in `flai/internal/metrics/digest.go`:

- `digestSection` takes the curator's activity document as well: `metrics.go` passes `opt.Activities`, which `statsread` already loads from `wip/agents/*.md`, and the section picks the one whose kind is the curator's (`workitem.ActivityCurator`). Nothing else in `statsread` changes.
- `DigestWeek` gains `CuratorRuns int` (`curator_runs`), the activities whose time falls in the week, and `AskedAgain` becomes `*int` (`asked_again,omitempty`), nil in a week no activity's time reaches: the newest activity is before the week's `start`, or there is no document. A week with such an activity keeps the count ADR-0200 defines.
- `Digest` gains `Curator *DigestCurator` (`curator,omitempty`) with `Runs int` (`runs`) over the window and `LastRun` (`last_run`, RFC 3339 UTC, omitted while empty), nil while the document does not exist.

Write the test that reproduces I-0207 first, in `digest_test.go`: threads opened to the operator and entries that list none of them in `asked_again`, and no curator activity. Before the change it reads `asked_again: 0` in each week; after it, no `asked_again`, `curator_runs: 0`, and no `curator`. Add a test with one activity in the second week: the first week has no `asked_again`, the second has its count and `curator_runs: 1`, and `curator` has `runs: 1` and that `last_run`. Keep the existing tests green, changing their expectations only where the ADR changes them.

In `flai/cmd/stats.go`, `printDigest` prints the text the ADR fixes: `asked again –` for a nil count, `curator runs n` on each week's line, and the `curator` line, or `curator  never ran`. Update `stats_test.go` and the digest assertions in `check_stats_test.go`, which asserts the printed section's text.

This waits on T-2121 because the JSON names and the printed text are the ADR's.

## Done when

- `flai stats --json` on a project with no `wip/agents/curator.md` reports `digest.weeks[]` with `curator_runs: 0` and no `asked_again`, and no `digest.curator`; with the document, each week's `curator_runs`, `asked_again` only in weeks a run has reached, and `digest.curator` with `runs` and `last_run`.
- `digest_test.go` has the test that reproduces I-0207 and the test with one run, and they pass.
- `flai stats` prints the text the ADR fixes, and `stats_test.go` and `check_stats_test.go` pass.
- `flai test` on the changed paths passes: gofmt, vet, golangci-lint, and the packages' tests.

## Notes

Drafted by planner-S-0504. `metrics.Options.Activities` is loaded at `flai/internal/statsread/statsread.go` through `repo.Activities()`, and `strategic.go` already reads the curator's document from it for `strategic[]`, so the section needs no new reader.
