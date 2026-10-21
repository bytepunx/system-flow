---
id: S-0458
type: story
nature: feature
title: The dashboard's pages are reached from a left nav rail of five collapsible groups, and the top bar holds only the name, the project picker, the theme, and the status pills
status: done
parent: E-0027
owner: alex
created: 2026-10-21T06:59:38Z
updated: 2026-10-21T08:32:26Z
transitions:
  - to: ready
    at: 2026-10-21T07:54:23Z
    by: alex
  - to: in-progress
    at: 2026-10-21T07:55:01Z
    by: agent-S-0458
  - to: review
    at: 2026-10-21T08:28:02Z
    by: agent-S-0458
  - to: done
    at: 2026-10-21T08:32:26Z
    by: orchestrator
tags: [dashboard]
topics: [dashboard, design]
touches: [flaiover/package.json, flaiover/pnpm-lock.yaml, design/tech/icons.md, design/tech/README.md, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/navrail.svelte.test.ts, design/components/flaiover/src/lib/components/NavRail.examples.json, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, flaiover/src/routes/layout.css, flaiover/src/lib/components/SiteMenu.svelte, flaiover/e2e/layout.e2e.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/components/SiteMenu.svelte.test.ts, design/issues/I-0177-pnpm-add-run-bare-in-a-story-worktree-installs-from-the-default-store-so-the-next-flaiover-install-sh-aborts-with-err-pnpm-aborted-remove-modules-dir-no-tty.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1997
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 5
      hand_edits: 2
      work: 51
  models:
    - model: claude-opus-5-5
      input: 344
      output: 141420
      cache_read: 17688967
      cache_write: 612489
      cost: 10.1539
  strategic:
    - kind: orchestrator
      seconds: 501
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 891
          cache_read: 6308091
          cache_write: 28755
          cost: 1.5888
        - model: claude-sonnet-5-5
          input: 12
          output: 76
          cache_read: 205218
          cache_write: 57700
          cost: 0.2152
forecast:
  duration: 1h30m
  delivery: 2026-10-21T10:19:00Z
  basis: "Its own forecast of 1h30m; 3rd in the pull order with an in-progress limit of 5, behind S-0456, S-0463, S-0464, S-0465, S-0459 and S-0460."
  by: flai
  at: 2026-10-21T07:54:30Z
finalized:
  by: alex
  at: 2026-10-21T07:29:11Z
---
# S-0458 The dashboard's pages are reached from a left nav rail of five collapsible groups, and the top bar holds only the name, the project picker, the theme, and the status pills

## Goal

The operator moves between flaiover's sections from a nav rail docked at the left, not from the two-tier site menu in the header. The rail groups the pages by concern into five groups:

- **Flow**: Board, Inbox, Charts, and Workflow, marked coming soon.
- **Agents**: Orchestration, Planning, Analysis, Activity, Threads, and Messages.
- **Design**: System and Components.
- **Documentation**: System, Conventions, ADRs, Issues, and Search.
- **Configuration**: Features, Agents, File System, Projects, and Harnesses.

Each group expands or collapses to show the icons of its pages. An expanded group's icon is highlighted, and a collapsed group's icon is dimmed. The top bar keeps only the flaiover name, the project picker, the theme setting, and the status pills. This story replaces `SiteMenu.svelte`. It builds the shell that the epic's other stories hang their pages and sidebars on.

## Acceptance criteria

- [x] The top bar holds only these, in this order: the flaiover name, which links to Overview (`/`); the project picker (`ProjectSwitcher.svelte`); the theme setting; and the status pills (`HostFlai.svelte`, which links to `/host`). The site menu's two rows are gone, and `SiteMenu.svelte` is removed.
- [x] `$lib/sitemenu.ts` defines the rail as data: the five groups in the epic's order, each with a name, an icon, and its pages in the epic's order. Each page has a label, an icon, and a route:
  - Board `/board`, Inbox `/inbox` with its count badge, and Charts `/charts`.
  - Orchestration `/workflow/orchestrator`, Planning `/workflow/planner`, and Analysis `/workflow/analyzer`.
  - Activity, Threads, and Messages at their own routes.
  - Design System `/design` and Components `/design/components`.
  - Documentation System, Conventions, and Issues at their documents under `/docs` (`design/system`, `design/conventions`, and `design/issues/summary.md`), and ADRs `/adrs` and Search `/search`.
  - The five Configuration pages at their sections of `/settings`, each by an anchor that `SettingsPanel.svelte` now carries.
  - Workflow is shown as coming soon. It is not a link, and it says so to assistive technology.

  `locate` still names the group and page for every route. `/items`, `/review`, and `/new` belong to Flow, and `/edit` to Documentation.
- [x] Each group is a button with `aria-expanded`. Pressing it shows or hides its pages' icons, each with its label and an accessible name. An expanded group's icon is drawn in the accent colour, and a collapsed group's icon is dimmed with the muted text colour. The current page has `aria-current="page"`. The group holding the current page is expanded on arrival. The groups the operator expanded or collapsed are kept per browser in `localStorage`.
- [x] A collapsed group still shows its pages' indicators beside its icon, so Flow carries the Inbox count while it is collapsed.
- [x] From Tailwind's `md` width the rail is docked at the left of every page, beside the main area. Below `md` it is a drawer: a menu button at the start of the top bar opens it, choosing a page or pressing Escape closes it, and focus returns to the button. A Playwright test at 390px and at 1440px shows that the page never scrolls sideways.
- [x] The icons come from the library the operator chose on TH-0553 question 2, recorded in `design/tech`. No icon is drawn from a colour or size outside the token utilities.
- [x] `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md` § The header describe the top bar and the rail in place of the site menu. `NavRail.svelte` has its examples file under `design/components/`.

## Tasks
- T-1953 flaiover has the icon library TH-0553 settles, recorded in design/tech
- T-1954 sitemenu.ts defines the rail's five groups, their icons, and their pages, and the settings sections carry anchors for the Configuration entries
- T-1955 NavRail.svelte draws the five groups, expanding and collapsing each with its icon highlighted or dimmed, and keeps the choice per browser
- T-1956 The layout's top bar holds the name, the project picker, the theme, and the status pills, with the rail docked from md and a drawer below it, and SiteMenu.svelte is gone
- T-1957 A Playwright test shows the rail docked at 1440px and as a drawer at 390px, with no page scrolling sideways
- T-1958 The dashboard's design and user docs describe the top bar and the nav rail in place of the site menu

## Notes

- Drafted by the planner for E-0027. First of five: the epic's other stories need the rail's data and the layout's main area.
- The Configuration entries point at anchors on `/settings` until the next story splits that page into five pages and points the entries at them. The Documentation System, Conventions, and Issues entries open the existing document viewer, so no new page is made for them.
- S-0420 (E-0024, backlog) merges the Planner, Orchestrator, and Analyzer pages into one role page, and touches `sitemenu.ts` and `SiteMenu.svelte`. Whichever of the two lands second moves the Agents group's three entries to the routes it then finds.
- Replaces the cancelled S-0452 (keep the top menu visible). With the rail docked, no menu scrolls away.

### Planning

- **Touches.** None were declared, so every touch is the planner's. All are files, and no folder touch was kept.
  - **Layout:** `flaiover/src/routes/+layout.svelte`, which holds the header today (lines 68 to 90), and `layout.css`. `flaiover/src/lib/sitemenu.ts` and its test hold `SITE_MENU` and `locate`. `SiteMenu.svelte` is removed. `NavRail.svelte`, `navrail.svelte.ts`, and their tests are new, following `boardtypes.svelte.ts` for the per-browser state. `SettingsPanel.svelte` gets the Configuration anchors. `flaiover/e2e/layout.e2e.ts` is new beside `live-events.e2e.ts`. `flaiover/package.json` and `pnpm-lock.yaml` take the icon library.
  - **Design:** `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md` § The header describe the site menu now. `design/tech/icons.md` is new, with its `README.md` row. `NavRail.examples.json` is under `design/components/`, as ADR-0164 requires.
  - **Co-change:** `flai touches suggest` lists hub documents (`design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md`), none of which this story changes. `items/[id]/+page.svelte` and `board/+page.svelte` are left to S-0460 and S-0462. None was added.
- **Forecast.** 1h30m against flai's 47m. flai's figure counts 7 criteria and 19 touches at the median for feature stories. This story also adds a dependency, builds a drawer whose behaviour is checked in a browser at two widths, and writes the project's second Playwright test, so the raise is about double. The delivery is flai's play-out, 30th in the pull order, moved by the 43m raise.
- **Cost of delay.** Not yet set. E-0027 has no inputs, so `flai cod` refuses. The input is asked on TH-0553 question 1.

### Accepted by the orchestrator

- Verified: 97ed4e1ccc30b49c81ac07b9837061a625279994
- Batch: 20261021-082935
- Enqueued verified: 97ed4e1ccc30b49c81ac07b9837061a625279994
- At: 2026-10-21T08:32:26Z

Verdict: S-0458 meets all seven acceptance criteria at 97ed4e1c; flai verify passed every step there (the flaiover tier covering flaiover-lint and vitest). Non-blocking: design/system/design-system-tool.md § The site menu and a comment in flaiover/src/routes/license/+page.svelte still describe the removed SiteMenu.
- 1: flaiover/src/routes/+layout.svelte, flaiover/src/lib/components/SiteMenu.svelte, flaiover/src/lib/components/SiteMenu.svelte.test.ts
- 2: flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts
- 3: flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/navrail.svelte.test.ts, flaiover/src/lib/components/NavRail.svelte.test.ts
- 4: flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/sitemenu.ts
- 5: flaiover/src/routes/+layout.svelte, flaiover/e2e/layout.e2e.ts
- 6: design/tech/icons.md, design/tech/README.md, flaiover/package.json, flaiover/pnpm-lock.yaml
- 7: design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/components/flaiover/src/lib/components/NavRail.examples.json
