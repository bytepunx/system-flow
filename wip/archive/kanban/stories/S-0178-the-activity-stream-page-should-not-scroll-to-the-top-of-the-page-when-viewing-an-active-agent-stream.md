---
id: S-0178
type: story
nature: remediation
title: The activity stream page should not scroll to the top of the page when viewing an active agent stream
status: done
owner: alex
created: 2026-10-01T07:38:11Z
updated: 2026-10-02T16:52:58Z
transitions:
  - to: ready
    at: 2026-10-01T07:39:43Z
    by: alex
  - to: in-progress
    at: 2026-10-02T16:36:00Z
    by: agent-S-0178
  - to: review
    at: 2026-10-02T16:52:40Z
    by: agent-S-0178
  - to: done
    at: 2026-10-02T16:52:58Z
    by: alex
tags: [dashboard]
topics: [client-side-activity]
touches: [flaiover/src/lib/components/AgentStream.svelte, flaiover/src/lib/components/AgentStream.svelte.test.ts, flaiover/src/routes/activity, flaiover/src/lib/project.svelte.ts, flaiover/src/lib/project.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1018
  models:
    - model: claude-opus-5-5
      input: 208
      output: 85428
      cache_read: 7150130
      cache_write: 308453
      cost: 5.0751
    - model: claude-sonnet-5-5
      input: 12
      output: 2887
      cache_read: 148298
      cache_write: 41870
      cost: 0.1632
---
# S-0178 The activity stream page should not scroll to the top of the page when viewing an active agent stream

## Goal

Prevent the activity page from scrolling automatically (this causes the operator to have to constantly scroll back to the activity window they'd been watching).

## Acceptance criteria
- [x] When viewing an agent's activity stream, the page does not scroll automatically
- [x] If the user scrolls away from the end of the stream, the page should not scroll back down automatically until the operator returns to the bottom again

## Tasks
- T-0704 A stream starts over only for another story or run, and opens without closing itself
- T-0705 The activity page keeps what it shows when a reload fails or flai is away
- T-0706 An empty project list keeps the project chosen
- T-0707 The dashboard's documents say the activity page holds still

## Notes

How the criteria were verified: the cause was reproduced in headless Chromium against the activity page with every `/api/*` mocked. Reloads of the agents (every 15 s while one runs) and of the activity emptied every stream window and pulled the window up (1555 → 793 px, to 0 when the emptied page fit the viewport), and put a box scrolled away from its end back at its end. With the stream following the values of its story and run (T-0704), the same reloads moved the window 0 px and a box scrolled up stayed where it was through 14 s of appends. The rest is covered by the component tests, since jsdom lays nothing out: a reload passing the same story and run changes nothing, a box away from its end is left there, a stream that opened stays open when its agent ends, a failed reload keeps the cards, an answer from no flai keeps the windows, and an empty project list keeps the project. Not tried in Firefox or Safari, nor against a live flai.
