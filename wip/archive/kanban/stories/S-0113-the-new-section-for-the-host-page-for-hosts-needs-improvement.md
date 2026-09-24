---
id: S-0113
type: story
nature: improvement
title: The new section for the host page for hosts needs improvement
status: done
parent: E-0003
owner: alex
created: 2026-09-24T07:40:17Z
updated: 2026-09-24T07:52:26Z
transitions:
  - to: ready
    at: 2026-09-24T07:40:35Z
    by: alex
  - to: in-progress
    at: 2026-09-24T07:41:50Z
    by: agent-S-0113
  - to: review
    at: 2026-09-24T07:50:39Z
    by: agent-S-0113
  - to: done
    at: 2026-09-24T07:52:26Z
    by: alex
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0113 The new section for the host page for hosts needs improvement

## Goal

Reduce the width of the Process column, increase the width of the other columns and provide enough space between them.

## Acceptance criteria
- [x] State, Version, and Restarts columns should have plenty of horizontal space between them

## Tasks
- T-0395 The host page's process table gives State, Version, and Restarts the width and a narrow Process column
- T-0396 The git convention's trailer names the model that authored the change, as the operator decided in TH-0009

## Notes

- Verified in a browser against the dev server with `/api/host` stubbed, at 1100px: Process 112px; State, Version, and Restarts 138px each with the actions shown and 207px each without, with 2rem of padding between them. At 600px they narrow to 71px each with the actions shown and still do not overlap.
- Both host-page cards widen from `max-w-2xl` to `max-w-3xl`: at the old width the three columns would get about 107px each with the actions shown, and comma lists like `running, stopped` would wrap.
