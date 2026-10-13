---
id: T-1515
type: task
nature: feature
title: flai doc lint and the host method doc.lint answer the markdown findings of content as it is edited
status: in-progress
parent: S-0378
owner: alex
created: 2026-10-13T16:56:32Z
updated: 2026-10-13T17:30:38Z
transitions:
  - to: ready
    at: 2026-10-13T17:30:37Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T17:30:38Z
    by: agent-S-0378
stream: S-0378
tags: [cli]
touches: [flai/cmd/doc_lint.go, flai/cmd/doc_lint_test.go, flai/internal/hostapi/docs.go, flai/internal/hostapi/docs_test.go]
---
# T-1515 flai doc lint and the host method doc.lint answer the markdown findings of content as it is edited

## Work

- Add `flai doc lint <path> --stdin --json` in a new `flai/cmd/doc_lint.go`, registered on the `doc` command from that file's own `init`, so it shares no file with T-1514. It lints the content given on standard input with `flai/internal/mdlint` and the project's markdownlint configuration (`mdlint.Load`), and answers each finding's line, rule, alias, and message.
- Add a read method `doc.lint` in `flai/internal/hostapi/docs.go`, taking `path` and `content`, with the same path guard as `doc.get`. It runs in process, since the editor calls it while the operator types.
- The findings are flai's subset of the markdownlint rules (ADR-0061), fast enough to run on each pause in typing. The markdown tier stays the authority at acceptance; the editor says so.
- Waits for nothing. It is the first layer, beside T-1514.

## Done when

- `flai doc lint docs/x.md --stdin --json` answers the findings of the content given, with line and rule, and an empty list for clean content.
- `doc.lint` answers the same over the host channel, and refuses a path outside the design, docs, and wip folders.
- `flai/cmd/doc_lint_test.go` and `flai/internal/hostapi/docs_test.go` cover both, and pass.

## Notes

- Criterion 2 of S-0378. No command lints arbitrary content today: `flai/internal/check` lints only the wip folder.
