---
id: T-0856
type: task
nature: remediation
title: mdlint reports MD038, spaces inside a code span, as markdownlint does
status: backlog
parent: S-0262
owner: alex
created: 2026-10-05T03:12:53Z
updated: 2026-10-05T03:12:53Z
transitions: []
stream: S-0262
tags: [flai]
touches: [flai/internal/mdlint]
---
# T-0856 mdlint reports MD038, spaces inside a code span, as markdownlint does

## Work

- Add MD038 (`no-space-in-code`, tags `whitespace` and `code`) to the rule table in `flai/internal/mdlint/mdlint.go`, implemented in `rules.go` over the code spans `inline.go` already finds (`codeClose`). Report a code span whose content starts or ends with a space where markdownlint-cli2 0.20.0 does, and nothing where it does not: a span of only spaces, and a span padded by one space on both sides, as CommonMark strips, are not findings. ADR-0061 says a case mdlint cannot place is not reported.
- Add a fixture under `testdata/cases` with I-0072's instance, T-0822's original body (`git show 38d51a2^` on its archived path), whose code span holding a list marker and the word Trigger ends with a space before its closing backtick. Add code-span cases beside it: a leading space, a trailing space, both, one space padding both sides, a span of spaces, a double-backtick span holding a backtick, and spans in a list item, a heading, link text, and a table cell.
- Regenerate `testdata/cases/expected.txt` with `scripts/mdlint-fixtures.sh`, so the test compares MD038's rule and line with markdownlint's.
- Add a test that `mdlint.Guard`, the step `item_new`, `flai task new`, and thread entries take, refuses I-0072's body naming MD038 and its line.
- Waits for nothing: it is the first layer.

## Done when

- `flai/internal/mdlint` reports MD038 on I-0072's body at the line markdownlint-cli2 reports, and agrees with markdownlint-cli2 on every fixture case
- The guard test refuses the body, naming MD038
- `scripts/flai-test.sh` passes, its integration tier linting every markdown file in the repository with no finding

## Notes

S-0240 added MD007 the same way (3446ac5): rule, fixture of the body that slipped through, regenerated reference, test.
