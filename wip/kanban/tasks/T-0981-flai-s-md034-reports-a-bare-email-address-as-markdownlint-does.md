---
id: T-0981
type: task
nature: remediation
title: flai's MD034 reports a bare email address as markdownlint does
status: backlog
parent: S-0265
owner: alex
created: 2026-10-05T05:50:04Z
updated: 2026-10-05T05:50:04Z
transitions: []
stream: S-0265
tags: [flai, mdlint]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases]
---
# T-0981 flai's MD034 reports a bare email address as markdownlint does

## Work

- Add a fixture, `flai/internal/mdlint/testdata/cases/email.md`, of bare email addresses: alone, mid-sentence, before a full stop or a comma, with `.`, `-`, `_`, and `+` in the local part, in a list item and a blockquote, and the cases that must not be reported: inside a code span, inside `<…>`, as a link destination, and in link text. Include the shape of TH-0067's entries of 2026-10-02.
- Regenerate `expected.txt` with `scripts/mdlint-fixtures.sh`, so that markdownlint-cli2 at the version `lint-md.sh` runs decides which lines MD034 reports.
- In `inline.go`, recognise GFM's extended email autolink beside `bareURL` (alphanumerics, `.`, `-`, `_`, `+`, then `@`, then a domain with at least one `.` whose last character is neither `-` nor `_`), skip it where `bareURL` would, and add it to the URLs `md034` reports.
- Add a regression test, `TestBareEmailOfI0056`, linting a thread entry with a bare address and asserting MD034 on its line.

Waits for nothing: it is the fix.

## Done when

- `TestFixturesMatchMarkdownlint` passes with the new fixture, so flai and markdownlint-cli2 report MD034 on the same lines.
- `TestBareEmailOfI0056` fails without the change to `inline.go` and passes with it.
- `TestRepositoryLintsClean` and `scripts/flai-test.sh` pass, and `flai check --strict` reports no new `markdown.MD034` in the main checkout.

## Notes
