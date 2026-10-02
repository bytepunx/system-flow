---
id: T-0707
type: task
nature: remediation
title: The dashboard's documents say the activity page holds still
status: done
parent: S-0178
owner: arobson
created: 2026-10-02T16:46:51Z
updated: 2026-10-02T16:51:12Z
transitions:
  - to: ready
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: in-progress
    at: 2026-10-02T16:50:54Z
    by: agent-S-0178
  - to: done
    at: 2026-10-02T16:51:12Z
    by: agent-S-0178
stream: S-0178
tags: []
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-0704, T-0705, T-0706]
usage:
  source: log
  seconds: 18
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 11
      output: 4568
      cache_read: 382345
      cache_write: 16494
      cost: 0.2714
---
# T-0707 The dashboard's documents say the activity page holds still

## Work

Say in the user guide's Activity section, and in the dashboard design where it describes the stream window, what the reader can now rely on: the page does not move under them as streams, agents, or projects are read again; a stream window follows its end only while the reader is at it; it opens when its agent starts and closes only when the reader closes it.

Waits for T-0704, T-0705, and T-0706: it describes what they settle.

## Done when

- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` say what the three tasks made true, and nothing they no longer do.

## Notes
