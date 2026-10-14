---
id: T-1516
type: task
nature: feature
title: flai accept's refusal for markdown the fix left carries each finding's file, line, rule, and checkout as data, through accept.run to the dashboard
status: done
parent: S-0378
owner: alex
created: 2026-10-13T16:56:46Z
updated: 2026-10-13T18:06:15Z
transitions:
  - to: ready
    at: 2026-10-13T17:54:58Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T17:54:59Z
    by: agent-S-0378
  - to: done
    at: 2026-10-13T18:06:15Z
    by: agent-S-0378
stream: S-0378
tags: [cli]
touches: [flai/cmd/accept_markdown.go, flai/cmd/accept.go, flai/cmd/accept_fix_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md]
after: [T-1514]
usage:
  source: log
  seconds: 676
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 28195
      cache_read: 4296879
      cache_write: 115178
      cost: 2.0769
---
# T-1516 flai accept's refusal for markdown the fix left carries each finding's file, line, rule, and checkout as data, through accept.run to the dashboard

## Work

- The ADR-0140 refusal in `flai/cmd/accept_markdown.go` is text today: a `rule:` error whose lines are `path:line: message (tier)`. The markdown tier's format is `plain`, so its `verify.Finding`s carry no path or line, only markdownlint-cli2's lines in the message.
- Parse each markdownlint-cli2 line (`path:line[:column] MDnnn/alias message`) of a failing fix tier into a finding with `path`, `line`, `rule`, and `message`, and keep the text refusal as it is.
- Name where each finding is to be corrected: `checkout: main` for the main checkout's uncommitted markdown, `checkout: story` and the story ID for the files the branch changes, so the dashboard opens the right file (T-1514's `--story`).
- With `--json`, print the refusal as `{refused: {id, reason, findings: [...]}}` and exit 4, as `flai doc save` does. Give `accept.run` in `flai/internal/hostapi/writes.go` `exits: {4: Refused}`, so `outcome()` passes the findings to the dashboard as the error's data.
- Waits for T-1514: both change `flai/internal/hostapi/writes.go` and its test.

## Done when

- `flai accept S-nnnn --json` on a branch with markdown the fix cannot fix exits 4 and prints each finding with its path, line, rule, message, and checkout.
- `accept.run` answers that refusal as a `Refused` error whose data holds the findings, and its NDJSON progress still ends with the error event.
- `flai/cmd/accept_fix_test.go` and `flai/internal/hostapi/writes_test.go` cover a main checkout finding and a branch finding, and pass.

## Notes

- Criterion 4 of S-0378. The dry run does not run the fix's check (ADR-0140), so the refusal is the only place the findings come from.
