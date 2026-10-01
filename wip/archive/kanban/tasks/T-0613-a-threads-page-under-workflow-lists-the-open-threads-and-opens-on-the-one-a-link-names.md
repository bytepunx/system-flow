---
id: T-0613
type: task
nature: feature
title: A Threads page under Workflow lists the open threads and opens on the one a link names
status: done
parent: S-0173
owner: arobson
created: 2026-10-01T07:35:14Z
updated: 2026-10-01T07:40:57Z
transitions:
  - to: ready
    at: 2026-10-01T07:35:29Z
    by: agent-S-0173
  - to: in-progress
    at: 2026-10-01T07:40:57Z
    by: agent-S-0173
  - to: done
    at: 2026-10-01T07:40:57Z
    by: agent-S-0173
stream: S-0173
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 0
  models: []
---

# T-0613 A Threads page under Workflow lists the open threads and opens on the one a link names

## Work

TH-0041: a /threads page, in the site menu's Workflow group, shows every open thread of the project with the story page's thread view: one at a time with the pager, show resolved, reply and resolve, and opened and scrolled to the thread ?thread= names. Each thread links to what it is anchored on.

## Done when

The Workflow group lists Threads; /threads?thread=TH-nnnn opens on that thread and lets the operator answer it; tests cover the menu and the page.

## Notes
