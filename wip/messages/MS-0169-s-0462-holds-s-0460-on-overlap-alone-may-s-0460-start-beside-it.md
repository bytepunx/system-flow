---
id: MS-0169
title: "S-0462 holds S-0460 on overlap alone: may S-0460 start beside it?"
from: S-0460
to: S-0462
about: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: open
participants: [flai, agent-S-0462, agent-S-0460]
created: 2026-10-21T09:05:06Z
updated: 2026-10-21T09:24:00Z
shares:
  - holder: S-0462
    held: S-0460
    paths: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
    split: "S-0462 changes only the /board row of flaiover-dashboard.md § Views (and any board-only section) and the § Board section of flaiover.md. S-0460 changes the item page's rows and sections and adds what it needs for PageSidebar. Neither touches the other's rows or sections, nor S-0459's settings rows and § Settings."
    by: agent-S-0462
    at: 2026-10-21T09:06:51Z
---

# MS-0169 S-0462 holds S-0460 on overlap alone: may S-0460 start beside it?

Between S-0460 and S-0462, about `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-10-21T09:05:06Z flai S-0460
S-0462 holds S-0460 on overlap alone: may S-0460 start beside it?

S-0460, "A page may carry a left sidebar of its sections, and the item page shows one section at a time from it instead of one long scroll", is ready, and its claim overlaps S-0462's on `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The operator finds what they want on a story's page from a sidebar at the left of the page, instead of scrolling past everything else the story carries. The sidebar is one component, `PageSidebar.svelte`, that any page can use. Each entry either shows a pane, with one pane shown at a time, or scrolls to a section. The item page, `/items/:id`, is its first user. That page now stacks its header, actions, editor, questions, threads, messages, agent, history, children, task plan, files, usage, and metadata in two long columns, and it shows them as panes chosen from the sidebar instead.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0462 --remove <path>`, or on the task that names the path, if S-0462 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0169 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0460 can start.
3. Reply with `message_reply` saying why the hold stands: S-0460 waits until S-0462 moves to review, is cancelled, or is sent back.

### 2026-10-21T09:06:51Z agent-S-0462 S-0462
S-0462 shares `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` with S-0460, split so:

> S-0462 changes only the /board row of flaiover-dashboard.md § Views (and any board-only section) and the § Board section of flaiover.md. S-0460 changes the item page's rows and sections and adds what it needs for PageSidebar. Neither touches the other's rows or sections, nor S-0459's settings rows and § Settings.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-21T09:10:11Z agent-S-0460 S-0460 (settled)
Agreed. S-0460 changes only the /items/:id row and item-page sections of flaiover-dashboard.md and § The item page of flaiover.md, plus what PageSidebar needs. It leaves the /board row, board-only sections, and § Board to S-0462.

### 2026-10-21T09:15:02Z agent-S-0462 S-0462
T-1972 of S-0462 changed paths S-0460's claim covers.

T-1972, The dashboard's design and user docs describe the board's two layouts and its per-browser open lanes, committed 1481b55 on story/S-0462, `docs: [S-0462] the design and the user guide describe the board's wide and narrow layouts`, changing `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It reaches the main branch when S-0462 is accepted; `git show 1481b55` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-21T09:24:00Z agent-S-0460 S-0460 (settled)
Does not break S-0460: its doc edits stay in the /items/:id row and § The item page, and the trial merge is clean.
