---
id: S-0155
type: story
nature: improvement
title: Inbox items take users to the story page, not document pages
status: done
parent: E-0013
owner: alex
created: 2026-09-29T05:57:30Z
updated: 2026-09-29T07:06:37Z
transitions:
  - to: ready
    at: 2026-09-29T06:04:02Z
    by: alex
  - to: in-progress
    at: 2026-09-29T07:00:34Z
    by: agent-S-0155
  - to: review
    at: 2026-09-29T07:06:28Z
    by: agent-S-0155
  - to: done
    at: 2026-09-29T07:06:37Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 376
  models:
    - model: claude-opus-5-5
      input: 132
      output: 27367
      cache_read: 6607347
      cache_write: 135907
      cost: 2.9566
---
# S-0155 Inbox items take users to the story page, not document pages

## Goal

Right now, clicking some parts of an inbox notification for open questions in story threads can take the operator to a documentation page instead of the related story.

Inbox items that require the operator to handle a thread on a story page should always take the operator to the story's page.

## Acceptance criteria
- [x] Clicking the inbox item where a user needs to respond to a question in a thread always takes the operator to the story page and displays the correct thread.

## Tasks
- T-0550 An inbox thread entry links to its story's page naming the thread
- T-0551 The story page opens on the thread its link names
- T-0552 The dashboard design and user guide say where an inbox thread leads

## Notes

- Verified by tests, not in a live dashboard: `hrefFor` sends a thread entry on an item to `/items/<item>?thread=<TH-id>` (`flaiover/src/lib/server/inbox.test.ts`), and `Threads.svelte` opens on the thread `select` names and scrolls it into view (`Threads.svelte.test.ts`); the item page passes `?thread=` as `select`. Running the branch's dashboard needs a second `flai serve`, which is not started here.
- Every story thread in `wip/threads` today is anchored on its story's item, so each one's entry now opens that story's page at that thread.
- Decision: a hand-written open question (kind `question`) still links to its narrative. It is not a thread, the story page does not show it, and the inbox answers it in place (S-0090).
- Decision: a thread on a task opens the task's page at that thread, because the story page lists only the threads anchored on the story (`threads.For`). Showing a task's threads on its story's page is a change to flai's `threads.list`, outside this story; no task-anchored thread exists today.
- Decision: a thread anchored on a document path with no item keeps its document page, where it is shown.
