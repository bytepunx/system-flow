---
id: T-0960
type: task
nature: feature
title: design/analysis holds the analyzer's reports and their index, flai check validates them, and the layout tables list the folder
status: in-progress
parent: S-0223
owner: alex
created: 2026-10-05T05:46:40Z
updated: 2026-10-06T20:15:10Z
transitions:
  - to: ready
    at: 2026-10-06T20:15:09Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:15:10Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [flai/internal/analysis/analysis.go, flai/internal/analysis/analysis_test.go, flai/internal/check/check.go, flai/internal/check/analysis_test.go, design/analysis/README.md, template/root/design/analysis/README.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, design/README.md, template/root/design/README.md.tmpl, design/system/repository-layout.md, docs/users/conventions.md, docs/users/index.md]
after: [T-0949]
---
# T-0960 design/analysis holds the analyzer's reports and their index, flai check validates them, and the layout tables list the folder

## Work

Add `flai/internal/analysis`, as `flai/internal/experiment` is for `design/experiments`: `Folder` (`analysis`), `Dir()` from the manifest's design folder, the report name `<date>-<focus>.md` (`ReportPath(date, focus)`), the focuses (`bottlenecks`, `intent`, `risk`, and `all` for a run with none), and `Validate`, which reads a report's front matter (`title`, `updated`, `status`, `focus`, and the window, `from` and `to`) and names what is missing or bad.

`flai check` validates every report in the folder and reports a report that `design/analysis/README.md` does not list (`flai/internal/check/check.go`, beside `c.experiments()`), and treats the folder as it treats `experiments` among the design folders (line 771).

Ship the folder's `README.md`, which says what a report holds and indexes the reports, in the repository and in the template (`template/root/design/analysis/README.md`). List the folder in the layout tables: `CLAUDE.md` and `template/root/CLAUDE.md.tmpl` (you may edit only your report), `design/README.md` and its template, `design/system/repository-layout.md` (the tree and a paragraph as `design/experiments/` has), and `docs/users/conventions.md` and `docs/users/index.md`.

It waits for T-0949, which changes `check.go` before it.

## Done when

- `analysis.Validate` names each missing or bad front matter field in a test, and `ReportPath` gives `<date>-<focus>.md`
- `flai check` reports a bad report and an unlisted one on a fixture, and passes on this repository
- the README is in both places and every layout table lists `design/analysis/`
- `go test ./internal/analysis/ ./internal/check/` passes

## Notes
