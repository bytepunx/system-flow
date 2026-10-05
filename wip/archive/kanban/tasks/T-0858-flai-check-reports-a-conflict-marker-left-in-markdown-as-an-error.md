---
id: T-0858
type: task
nature: remediation
title: flai check reports a conflict marker left in markdown as an error
status: done
parent: S-0253
owner: alex
created: 2026-10-05T03:13:16Z
updated: 2026-10-05T03:16:51Z
transitions:
  - to: ready
    at: 2026-10-05T03:13:37Z
    by: agent-S-0253
  - to: in-progress
    at: 2026-10-05T03:13:37Z
    by: agent-S-0253
  - to: done
    at: 2026-10-05T03:16:51Z
    by: agent-S-0253
stream: S-0253
tags: []
touches: [flai/internal/conflictmark, flai/internal/check, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 194
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 87
      output: 24164
      cache_read: 3034316
      cache_write: 105721
      cost: 1.7685
    - model: claude-sonnet-5
      input: 38
      output: 11540
      cache_read: 776950
      cache_write: 67269
      cost: 0.439
---
# T-0858 flai check reports a conflict marker left in markdown as an error

## Work

Add `flai/internal/conflictmark`, a pure scanner that returns the lines of a text that open or close a merge conflict: a line that is `<<<<<<<` or `>>>>>>>` followed by a space or the end of the line (and `|||||||`, diff3's base marker). `=======` alone is not flagged: it is a valid setext heading underline in markdown, and the opening and closing markers already find every conflict. `flai check` gains `markdown.conflict-marker`, an error on each such line in every `.md` file under the design, docs, and wip folders and at the repository root, archive and conventions included, saying to resolve the conflict and keep what both sides meant. An error, not a warning: the document is corrupt, and CI's `flai check` without `--strict` must fail on it. Behaviour tests reproduce I-0066: a design document carrying the markers S-0201's rebase left fails the check on the marker lines. Document the rule beside `markdown.MDnnn` in `design/system/flai-cli.md`'s `flai check` row and in `docs/users/flai.md`.

Waits for nothing: it is the first layer.

## Done when

- `go test ./internal/conflictmark/ ./internal/check/` passes, with a test that fails without the rule.
- The design and the user guide describe `markdown.conflict-marker`.

## Notes
