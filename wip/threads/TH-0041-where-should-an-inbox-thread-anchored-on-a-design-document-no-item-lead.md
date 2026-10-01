---
id: TH-0041
title: Where should an inbox thread anchored on a design document (no item) lead?
anchor:
  path: wip/kanban/stories/S-0173-inbox-open-narratives-navigate-to-the-wrong-location.md
  item: S-0173
status: resolved
participants: [agent-S-0173, alex]
created: 2026-09-30T01:26:55Z
updated: 2026-10-01T07:47:06Z
---

# TH-0041 Where should an inbox thread anchored on a design document (no item) lead?

On wip/kanban/stories/S-0173-inbox-open-narratives-navigate-to-the-wrong-location.md.

## Entries

### 2026-09-30T01:26:55Z agent-S-0173
S-0173's second criterion says no inbox entry may ever lead to documentation. Open questions in narratives will lead to their story's page, which will show them and let you answer them there. Threads on an item, or on an item's file or narrative, will lead to that item's page opened on the thread.

One case has no item: a thread anchored on a document, such as TH-0017 on design/issues/I-0027-....md. Such a thread lives on that document; the document's page is the only place that shows it and lets you answer it.

Recommended: such a thread leads to the document's page, opened on that thread and scrolled to it (today the document page ignores which thread was named). That is the thread's own place, not the narrative problem this story is about.

Alternative: a new thread page (/threads/TH-nnnn) that shows one thread and its reply box, linked to from the inbox for every thread without an item. Bigger, and a new view.

Which one?

### 2026-10-01T07:23:46Z alex
lets add a "threads" page that makes it possible to see all open threads with the same navigation behavior that the story page has and navigate to that when there is no parent item. the threads page belongs under the Workflow parent menu item.

### 2026-10-01T07:46:53Z agent-S-0173
Done on `story/S-0173` (T-0613). The site menu's Workflow group now has a `/threads` page. It shows every open thread with the story page's thread view: one at a time with the pager, **show resolved**, reply, and resolve. It opens and scrolls to the thread `?thread=` names.

An inbox thread with no item leads there. A thread on an item's file or narrative still goes to that item's page. Each thread on the page links to what it is anchored on. New threads are started on the item or document they are about, since a thread needs an anchor.

### 2026-10-01T07:47:06Z agent-S-0173
Resolved: Decided by alex: a Threads page under Workflow. Built in S-0173 (T-0613); an inbox thread with no item leads there.
