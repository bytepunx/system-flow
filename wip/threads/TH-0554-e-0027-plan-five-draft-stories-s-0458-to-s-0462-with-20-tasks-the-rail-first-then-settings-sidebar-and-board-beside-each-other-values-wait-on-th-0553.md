---
id: TH-0554
title: "E-0027 plan: five draft stories, S-0458 to S-0462, with 20 tasks; the rail first, then settings, sidebar, and board beside each other; values wait on TH-0553"
anchor:
  path: wip/kanban/epics/E-0027-rework-main-layout.md
  item: E-0027
status: open
participants: [planner-E-0027, orchestrator]
created: 2026-10-21T07:06:42Z
updated: 2026-10-21T11:04:47Z
---

# TH-0554 E-0027 plan: five draft stories, S-0458 to S-0462, with 20 tasks; the rail first, then settings, sidebar, and board beside each other; values wait on TH-0553

On wip/kanban/epics/E-0027-rework-main-layout.md.

## Entries

### 2026-10-21T07:06:42Z planner-E-0027
E-0027 has five draft stories in the backlog, each with its tasks drafted. All of them pass `flai check --strict`. The forecasts are set. The cost of delay values wait on TH-0553 question 1.

## Stories and order

| Story | After | Forecast | What it delivers |
|-------|-------|----------|------------------|
| S-0458 | none | 1h30m | The nav rail of five groups and the slim top bar. `SiteMenu.svelte` goes, and below `md` the rail is a drawer. |
| S-0459 | S-0458 | 45m | `/settings` splits into the five Configuration pages. |
| S-0460 | S-0458 | 1h | `PageSidebar.svelte`, and the item page shown as panes. |
| S-0461 | S-0460 | 45m | The review, charts, and documents pages get the sidebar. |
| S-0462 | S-0458 | 50m | The board as one sideways-scrolling row from `sm`, and as lane headings that toggle below `sm`. |

S-0459, S-0460, and S-0462 wait only for S-0458. Their code paths do not overlap. All three update `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`, which are not shared paths, so while one is in progress the others are held on those two files.

## Tasks and layers

- **S-0458**
  - Layer 1: T-1953 (icon library) and T-1954 (`sitemenu.ts` data, settings anchors).
  - Layer 2: T-1955 (`NavRail.svelte`).
  - Layer 3: T-1956 (layout, top bar, drawer, remove `SiteMenu`).
  - Layer 4: T-1957 (Playwright) and T-1958 (docs).
- **S-0459**
  - Layer 1: T-1959 (`settingspages.ts`, `SettingsPanel` `page` prop).
  - Layer 2: T-1960 (routes, redirects, rail entries).
  - Layer 3: T-1961 (docs).
- **S-0460**
  - Layer 1: T-1962 (`PageSidebar`).
  - Layer 2: T-1963 (item page panes).
  - Layer 3: T-1964 (docs).
- **S-0461**
  - Layer 1: T-1965 (review), T-1966 (charts), and T-1967 (documents).
  - Layer 2: T-1968 (docs).
- **S-0462**
  - Layer 1: T-1969 (`lanecollapse.svelte.ts`).
  - Layer 2: T-1970 (board layouts, `max-w-6xl` opt-out).
  - Layer 3: T-1971 (Playwright) and T-1972 (docs).

## Assumptions

- TH-0553 questions 2 to 5 are taken as recommended unless you answer otherwise: `@lucide/svelte`, where Overview, Updates, and License go, the settings split, and the drawer below `md`.
- The rail's routes stay where they are: Orchestration, Planning, and Analysis remain under `/workflow/*`.
- Documentation's System, Conventions, and Issues entries open the existing document viewer at `design/system`, `design/conventions`, and `design/issues/summary.md`. No new Issues page is made.
- Workflow (coming soon) is a disabled entry, not a page.
- Pages keep `max-w-6xl` beside the rail. Only the board takes the whole width.
- On a phone the board opens Ready, In progress, and Review at first.
- The item page's eight panes are my reading of today's page. The story's agent may merge them.

## Proposals

- **Overlap with S-0420** (E-0024, backlog). It merges the Planner, Orchestrator, and Analyzer pages and touches `sitemenu.ts` and `SiteMenu.svelte`, which S-0458 removes. If S-0420 lands later, it edits `NavRail`'s data instead. I propose no change to either story, only that order.
- **Possible addition, not drafted:** a later story giving the sidebar to Inbox (by kind), Harnesses (by harness), and the Design page (by token group). I left it out to keep S-0461 to one or two sessions. Say if you want it drafted.
- **Archive the cancelled S-0452.** E-0027 replaces it, and `flai check` reports it as `item.archive`.

Nothing to split, merge, or drop among the five.

### 2026-10-21T11:04:47Z orchestrator (recommendation)
Recommended: resolve this thread. The plan was delivered as drafted.

- You finalized all five stories with the order and tasks shown here.
- S-0458, S-0459 and S-0462 are accepted and released.
- S-0460 is in progress, and S-0461 waits for it.
- The only open input is the cost of delay on TH-0553, question 1, and that thread carries it.

Nothing on this thread still needs a decision.

Source: design/conventions/strategic-agents.md § As the planner
