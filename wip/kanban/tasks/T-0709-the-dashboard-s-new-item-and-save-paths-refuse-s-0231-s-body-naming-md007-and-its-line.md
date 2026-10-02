---
id: T-0709
type: task
nature: remediation
title: The dashboard's new-item and save paths refuse S-0231's body naming MD007 and its line
status: in-progress
parent: S-0240
owner: arobson
created: 2026-10-02T16:50:06Z
updated: 2026-10-02T16:56:55Z
transitions:
  - to: ready
    at: 2026-10-02T16:56:55Z
    by: agent-S-0240
  - to: in-progress
    at: 2026-10-02T16:56:55Z
    by: agent-S-0240
stream: S-0240
tags: []
touches: [flai/internal/itemnew, flai/internal/hostapi, flaiover/src/lib/server]
after: [T-0708]
usage:
  source: log
  seconds: 43
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 214
      cache_read: 1090266
      cache_write: 39765
      cost: 0.4557
---
# T-0709 The dashboard's new-item and save paths refuse S-0231's body naming MD007 and its line

## Work

Prove the dashboard's two paths refuse what the CLI refuses, now that the lint knows MD007: `item.new` (`flaiover/src/routes/api/items` → host action → `flai story new --body-stdin` → `itemnew.Create`) and `doc.save` (an edit that follows, `docedit.Save`). Add a test in `flai/internal/itemnew` that creating a story with S-0231's original body is refused with a `markdown.MD007` finding at its line and leaves nothing behind, and one in `flaiover/src/lib/server/writes.test.ts` that the `item.new` and `doc.save` actions are refused with the rule and line in their findings. Confirm the new-item form and the editor show a refusal's findings with rule and line, and fix it if they do not. Waits for T-0708, which adds the rule these tests need.

## Done when

- The itemnew test and the flaiover writes tests pass, and each fails without T-0708's rule.
- The dashboard form shows a refusal's rule and line, as the CLI prints them.

## Notes
