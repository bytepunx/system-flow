---
id: S-0460
type: story
nature: feature
title: A page may carry a left sidebar of its sections, and the item page shows one section at a time from it instead of one long scroll
status: backlog
parent: E-0027
owner: alex
created: 2026-10-21T07:00:19Z
updated: 2026-10-21T07:05:58Z
transitions: []
tags: [dashboard]
topics: [dashboard, design]
touches: [flaiover/src/lib/components/PageSidebar.svelte, flaiover/src/lib/components/PageSidebar.svelte.test.ts, flaiover/src/lib/pagesidebar.svelte.ts, flaiover/src/lib/pagesidebar.svelte.test.ts, design/components/flaiover/src/lib/components/PageSidebar.examples.json, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0458]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h
  delivery: 2026-10-22T01:32:00Z
  basis: "flai's 26m from 107s per unit of size counts nine touches, but one is the 669-line item page restructured into panes with its tests rewritten; delivery is flai's play-out moved by its own 34m raise and S-0458's 43m."
  by: planner-E-0027
  at: 2026-10-21T07:05:58Z
---
# S-0460 A page may carry a left sidebar of its sections, and the item page shows one section at a time from it instead of one long scroll

## Goal

The operator finds what they want on a story's page from a sidebar at the left of the page, instead of scrolling past everything else the story carries. The sidebar is one component, `PageSidebar.svelte`, that any page can use. Each entry either shows a pane, with one pane shown at a time, or scrolls to a section. The item page, `/items/:id`, is its first user. That page now stacks its header, actions, editor, questions, threads, messages, agent, history, children, task plan, files, usage, and metadata in two long columns, and it shows them as panes chosen from the sidebar instead.

## Acceptance criteria

- [ ] `PageSidebar.svelte` takes a list of entries, each with a label, an id, an optional count, and a mode: pane or scroll. It renders them as a left-aligned list, with `aria-current` on the entry shown.
  - A pane entry shows its pane and hides the others.
  - A scroll entry scrolls to its section, and the entry is current while its section is in view.
  - The chosen entry is kept in the URL's hash, so back, forward, and a shared link open the same pane.
  - From Tailwind's `md` width the sidebar is sticky beside the content. Below `md` it is a single row above the content that never makes the page scroll sideways.
  - It is reachable and usable by keyboard alone.
- [ ] On the item page, the title, kind chips, status line, notices, and action buttons stay above the panes on every pane. The panes are these, in this order:
  - **Story**: the body in `ItemEditor`, with the review block.
  - **Tasks**: children and the task plan, for an item with children.
  - **Questions and threads**: `OpenQuestions` and `Threads`, with the count of open ones.
  - **Messages**, for a story with conversations.
  - **Agent**: `StoryAgent` and `ForgeStatus`.
  - **History**.
  - **Files**.
  - **Planning and usage**: cost of delay, forecast, expected cost, usage, and the metadata block.

  A pane with nothing to show for the item is left out of the sidebar.
- [ ] Story is the pane shown on arrival. A link with `?thread=<TH-id>` or `?question=<key>`, as the Inbox makes, opens Questions and threads with that thread or question in view. A link with a pane's hash opens that pane.
- [ ] The item page's existing tests pass, each finding its content on the pane that now holds it. New tests cover the panes left out, the arrival pane, and the inbox links.
- [ ] `PageSidebar.svelte` has its examples file under `design/components/`. `design/system/flaiover-dashboard.md` § Views (`/items/:id`) and `docs/users/flaiover.md` § The item page describe the sidebar and the panes.

## Tasks
- T-1962 PageSidebar.svelte lists a page's sections at its left, each showing a pane or scrolling to a section, with the choice kept in the URL's hash
- T-1963 The item page keeps its header and actions on top and shows its sections as panes chosen from the sidebar, with the Inbox's thread and question links opening their pane
- T-1964 The dashboard's design and user docs describe the page sidebar and the item page's panes

## Notes

- Drafted by the planner for E-0027. Third, after S-0458, whose layout gives the main area the sidebar sits in. It can run beside S-0459 and S-0462, which touch none of its code, though all three update the same two documents.
- The pane list is the planner's reading of today's page (`flaiover/src/routes/items/[id]/+page.svelte`). The story's agent may merge or rename panes when the page reads better, and should say why in the narrative.
