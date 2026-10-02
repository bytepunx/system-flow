---
id: T-0708
type: task
nature: remediation
title: mdlint reports MD007 ul-indent as markdownlint does
status: done
parent: S-0240
owner: arobson
created: 2026-10-02T16:50:06Z
updated: 2026-10-02T16:56:51Z
transitions:
  - to: ready
    at: 2026-10-02T16:50:32Z
    by: agent-S-0240
  - to: in-progress
    at: 2026-10-02T16:50:32Z
    by: agent-S-0240
  - to: done
    at: 2026-10-02T16:56:51Z
    by: agent-S-0240
stream: S-0240
tags: []
touches: [flai/internal/mdlint]
usage:
  source: log
  seconds: 379
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 67
      output: 20579
      cache_read: 2529133
      cache_write: 69658
      cost: 1.3602
---
# T-0708 mdlint reports MD007 ul-indent as markdownlint does

## Work

flai's markdown lint has no MD007 (ul-indent), so S-0231's goal list, indented one space, passed `flai story new` and the dashboard's `item.new` alike, and markdownlint-cli2 stopped S-0193's close-out on it (I-0055). Implement MD007 in `flai/internal/mdlint` as markdownlint 0.20.0's cli2 runs it: an unordered list item's indent is `start_indent` (when `start_indented`) plus `indent` per enclosing unordered list, only where no enclosing list is ordered, read from the project's markdownlint file with markdownlint's options and alias `ul-indent`. Add a fixture that reproduces S-0231's original body, regenerate the reference with `scripts/mdlint-fixtures.sh`, and keep mdlint silent where markdownlint is. Waits for nothing: the other tasks test what this adds.

## Done when

- `mdlint` reports MD007 on S-0231's original body at the four lines markdownlint-cli2 reports, and the agreement test passes with the regenerated `expected.txt`.
- A test fails with the rule taken out.
- `go test ./internal/mdlint` passes, and linting every markdown file in this repository still reports nothing.

## Notes
