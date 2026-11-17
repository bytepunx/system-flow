---
id: T-3529
type: task
nature: remediation
title: mdlint and markdownlint-cli2 agree on D-0592's line 17 in both its wordings, pinned by the code span fixture
status: done
parent: S-0836
owner: alex
created: 2026-11-17T08:29:26Z
updated: 2026-11-17T11:32:08Z
transitions:
  - to: ready
    at: 2026-11-17T11:17:03Z
    by: agent-S-0836
  - to: in-progress
    at: 2026-11-17T11:17:03Z
    by: agent-S-0836
  - to: done
    at: 2026-11-17T11:32:08Z
    by: agent-S-0836
stream: S-0836
tags: [flai, mdlint]
touches: [flai/internal/mdlint/testdata/cases/code-spans.md, flai/internal/mdlint/testdata/cases/expected.txt]
usage:
  source: log
  seconds: 905
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 7241
      cache_read: 4718125
      cache_write: 59038
      cost: 1.5609
---
# T-3529 mdlint and markdownlint-cli2 agree on D-0592's line 17 in both its wordings, pinned by the code span fixture

## Work

- Establish the divergence first: lint line 17 of `design/digest/D-0592-flai-stream-sync-json-prints-a-map-that-flai-cmd-stream-go-builds-by-hand-not.md` with markdownlint-cli2 0.20.0 and with `flai/internal/mdlint`, as the line is on main and as the curator first wrote it (commit 8ed4cd414). Record the cause in the narrative's Decisions.
- Add both wordings to `flai/internal/mdlint/testdata/cases/code-spans.md` and regenerate `expected.txt` with `scripts/mdlint-fixtures.sh`.
- Planned as a parser fix in `inline.go` or `rules.go`. The divergence was not there: markdownlint-cli2 reports the same MD038 on the first wording (17:146, context `` `). But ` ``) and nothing on main's, as mdlint does. So no parser change, and no MD038 table case: the fixture pins the agreement against markdownlint's own findings.
- Waits for no task: it is the first layer.

## Done when

- `TestFixturesMatchMarkdownlint` passes with both wordings in the fixture: one MD038 on the first, none on main's.
- `flai test` on the changed paths passes.

## Notes

- `flai test` on the two fixture paths passed (vet, golangci-lint, go-test), 2026-11-17.
