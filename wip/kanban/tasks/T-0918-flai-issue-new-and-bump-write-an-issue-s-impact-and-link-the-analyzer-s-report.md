---
id: T-0918
type: task
nature: feature
title: flai issue new and bump write an issue's impact and link the analyzer's report
status: backlog
parent: S-0224
owner: alex
created: 2026-10-05T05:44:49Z
updated: 2026-10-05T05:44:49Z
transitions: []
stream: S-0224
tags: [flai]
touches: [flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/issues/record.go, flai/internal/issues/record_test.go, flai/cmd/issue.go, flai/cmd/issue_test.go]
---
# T-0918 flai issue new and bump write an issue's impact and link the analyzer's report

## Work

The first layer: everything else in S-0224 builds on what an issue holds after this task, so it waits for nothing.

- `flai issue new` takes the impact the analyzer measures: `--revenue-per-week`, `--penalty-per-week`, and `--time-lost-per-cycle`, and `--evidence` for the words behind them. Given any, the issue gets an `## Impact` section in the format `design/system/continuous-improvement.md` gives, one `- key: value` line per input, which `issues.impact` in `stories.go` already reads. An amount must be a number of zero or more, a duration a Go duration longer than zero; anything else is refused before the issue is written.
- `--report <path>` names the analysis report under `design/analysis/` that found it. The instance says `Report: <path>.`, and the `## Remediation` section, added when missing, gets a line linking the report by a relative path. A path that is not a markdown file under `design/analysis/` is refused.
- Deduplication: with `--report`, an open issue with the same title (`FindOpenByTitle`) is bumped instead of a new one being made, as `RecordOnce` does for close-out findings, and the result says which happened. `flai issue bump <id> --report <path>` links the report on an issue the analyzer judged to be the same finding under another title, and takes the impact flags too, replacing the `## Impact` lines it gives and keeping the others.
- `--json` returns the issue's ID, path, and whether it was new or bumped, so the analyzer can link it from its report.

## Done when

- `flai issue new` with the impact flags and `--report` writes an issue whose `## Impact` and `## Remediation` sections hold them, and `issues.impact` reads the inputs back
- With `--report`, an open issue of the same title is bumped and linked, not duplicated; `flai issue bump --report` links a report and updates the impact
- Bad amounts, durations, and report paths are refused with nothing written
- Tests in `flai/internal/issues` and `flai/cmd/issue_test.go` cover a new issue, a bumped duplicate, and each refusal; `scripts/flai-test.sh` passes

## Notes
