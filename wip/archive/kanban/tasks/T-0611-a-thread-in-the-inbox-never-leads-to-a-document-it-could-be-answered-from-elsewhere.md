---
id: T-0611
type: task
nature: feature
title: A thread in the inbox never leads to a document it could be answered from elsewhere
status: done
parent: S-0173
owner: alex
created: 2026-09-30T01:26:59Z
updated: 2026-10-01T07:40:58Z
transitions:
  - to: ready
    at: 2026-09-30T01:27:10Z
    by: agent-S-0173
  - to: in-progress
    at: 2026-10-01T07:40:58Z
    by: agent-S-0173
  - to: done
    at: 2026-10-01T07:40:58Z
    by: agent-S-0173
stream: S-0173
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 0
  models: []
---

# T-0611 A thread in the inbox never leads to a document it could be answered from elsewhere

## Work

A thread whose anchor has no item but whose path is a work item's file or a narrative leads to that item's page opened on the thread. A thread on another document leads where TH-0041 decides.

## Done when

No thread entry with an item, an item file, or a narrative leads to /docs; the TH-0041 case is handled as answered; tests cover each.

## Notes
