---
id: T-0701
type: task
nature: feature
title: Main's markdown lint passes again and I-0056 records why flai let TH-0067's bare email through
status: done
parent: S-0191
owner: arobson
created: 2026-10-02T16:11:54Z
updated: 2026-10-02T16:11:58Z
transitions:
  - to: ready
    at: 2026-10-02T16:11:58Z
    by: agent-S-0191
  - to: in-progress
    at: 2026-10-02T16:11:58Z
    by: agent-S-0191
  - to: done
    at: 2026-10-02T16:11:58Z
    by: agent-S-0191
stream: S-0191
tags: [dashboard]
touches: [design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md, design/issues/summary.md]
after: [T-0700]
usage:
  source: log
  seconds: 0
  models: []
---
# T-0701 Main's markdown lint passes again and I-0056 records why flai let TH-0067's bare email through

## Work

S-0191's close-out stopped at the markdown lint on `wip/threads/TH-0067-…`, which S-0231's acceptance committed to main with a bare email address four times (MD034). Wrap the addresses as autolinks in one `chore` commit on main, under TH-0017's answer A, and record the occurrence as an issue, as TH-0056 asks for findings outside a story. It waits for T-0700 because the close-out of T-0700's work is what found it.

## Done when

- [x] TH-0067 on main has no bare email address, and the markdown lint passes on it.
- [x] An issue records the occurrence, the cause in flai's MD034, and a remedy to consider.

## Notes
