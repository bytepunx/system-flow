---
id: T-0982
type: task
nature: remediation
title: The user guide says flai's lint checks bare email addresses
status: done
parent: S-0265
owner: alex
created: 2026-10-05T05:50:06Z
updated: 2026-10-07T08:54:44Z
transitions:
  - to: ready
    at: 2026-10-07T08:48:04Z
    by: agent-S-0265
  - to: in-progress
    at: 2026-10-07T08:48:05Z
    by: agent-S-0265
  - to: done
    at: 2026-10-07T08:54:44Z
    by: agent-S-0265
stream: S-0265
tags: [docs]
touches: [docs/users/flai.md]
usage:
  source: log
  seconds: 399
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 27
      output: 10773
      cache_read: 1319312
      cache_write: 36225
      cost: 0.7692
---
# T-0982 The user guide says flai's lint checks bare email addresses

## Work

In `docs/users/flai.md`, the `flai check` paragraph on the markdownlint configuration lists what flai checks as "bare URLs". Make it "bare URLs and email addresses".

Waits for nothing: it touches no path the lint task does, so the two run together.

## Done when

- The paragraph names bare email addresses among what flai checks.
- `scripts/lint-md.sh` passes on `docs/users/flai.md`.

## Notes
