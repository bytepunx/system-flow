---
id: T-3529
type: task
nature: remediation
title: mdlint parses the code spans of D-0592's line 17 as markdownlint-cli2 does and reports no MD038 there, with a fixture that reproduces it
status: backlog
parent: S-0836
owner: alex
created: 2026-11-17T08:29:26Z
updated: 2026-11-17T08:29:26Z
transitions: []
stream: S-0836
tags: [flai, mdlint]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/rules.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/code-spans.md, flai/internal/mdlint/testdata/cases/expected.txt]
---
# T-3529 mdlint parses the code spans of D-0592's line 17 as markdownlint-cli2 does and reports no MD038 there, with a fixture that reproduces it

## Work

- Establish the divergence first. Lint the paragraph at line 17 of `design/digest/D-0592-flai-stream-sync-json-prints-a-map-that-flai-cmd-stream-go-builds-by-hand-not.md` with markdownlint-cli2 0.20.0 (through `npx`, as `scripts/mdlint-fixtures.sh` runs it) and with `flai/internal/mdlint`. Find which code span mdlint reports and where its span boundaries differ from micromark's. CommonMark gives a backslash no escaping power inside a code span, while one before a backtick outside a span escapes it, so check the opening of the next span and the padding rule in `md038`, not only `codeClose`. Record the cause in the narrative's Decisions.
- Add the line, or the smallest text that reproduces the finding, to `flai/internal/mdlint/testdata/cases/code-spans.md`, and regenerate `expected.txt` with `scripts/mdlint-fixtures.sh`.
- Add a case to the MD038 table in `flai/internal/mdlint/mdlint_test.go` that fails before the fix.
- Fix the parse in `flai/internal/mdlint/inline.go` (`parseRange`, `codeClose`), or the rule in `rules.go` (`md038`), whichever holds the cause. Widen the touches if the cause lies elsewhere in the package.
- Waits for no task: it is the first layer.

## Done when

- `TestFixturesMatchMarkdownlint` passes with the new fixture, and the new MD038 case passes.
- `TestRepositoryLintsClean` reports nothing on D-0592.
- `flai test` on the changed paths passes.

## Notes
