---
id: TH-0035
title: "The story page already pages threads (S-0133): what did you see stacked?"
anchor:
  path: wip/kanban/stories/S-0153-threads-in-the-story-pages-are-supposed-to-be-paged.md
  item: S-0153
status: resolved
participants: [agent-S-0153, alex]
created: 2026-09-29T06:46:39Z
updated: 2026-09-29T06:59:49Z
---

# TH-0035 The story page already pages threads (S-0133): what did you see stacked?

On wip/kanban/stories/S-0153-threads-in-the-story-pages-are-supposed-to-be-paged.md.

## Entries

### 2026-09-29T06:46:39Z agent-S-0153
The story page, the review page, and document pages have shown one thread at a time since S-0133 (26 Sep). The running dashboard image (built 29 Sep 00:33Z) has it. It shows `2 of 3 ← →` in small text above and below the thread, and it only appears when there are two or more threads.

So S-0153's three criteria read as already met, and I cannot tell what you saw. I could not look myself: the dashboard on :4242 reports no host flai connected right now.

My guess is the pager is too easy to miss. Unless you say otherwise, I will:

1. Put the count between the arrows: `← 2 of 3 →`.
2. Make the arrows button-sized, next to the Threads heading, so they're hard to miss.
3. Let the left and right arrow keys page while focus is in the pager.

Is that it? Or did you see threads stacked on a particular page (its URL, or the story ID)? If so, tell me which one and I'll fix that instead.

### 2026-09-29T06:48:48Z alex
The story I last saw it on is already archived. I am not sure I remember for sure but maybe 147 or 149

### 2026-09-29T06:51:29Z agent-S-0153
Thanks. Neither story had more than one thread showing:

- S-0149's page has one thread, TH-0032, with 11 entries.
- S-0147 has one thread, TH-0033, which is resolved and so hidden unless you tick "show resolved".

So the long scroll was probably TH-0032's 11 entries, not several threads stacked.

The pager change is done on the branch: the page now shows `← 1 of 3 →` beside the Threads heading and again under the thread, and the arrow keys page.

For long threads, which of these do you want?

1. **Collapse older entries (recommended).** A thread shows its last two entries, with "show 9 earlier entries" above them to expand the rest. You keep the conversation's context and land on what is new.
2. **Page the entries.** One entry at a time, with its own `← n of m →`.
3. **Nothing more.** The thread pager is what you meant; I move S-0153 to review.

If I hear nothing, I'll do 1 as a second task.

### 2026-09-29T06:57:09Z alex
1

### 2026-09-29T06:59:49Z agent-S-0153
Done as option 1 (bc47b14): a thread of more than two entries shows its last two, with show N earlier entries above them. S-0153 is going to review.

Resolved.
