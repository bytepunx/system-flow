---
id: S-0459
type: story
nature: improvement
title: The settings page becomes five Configuration pages, Features, Agents, File System, Projects, and Harnesses, each reached from the rail
status: backlog
parent: E-0027
owner: alex
created: 2026-10-21T06:59:57Z
updated: 2026-10-21T07:19:47Z
transitions: []
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/settingspages.ts, flaiover/src/lib/settingspages.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/routes/settings/+page.svelte, flaiover/src/routes/settings/+page.ts, "flaiover/src/routes/settings/[page]/+page.svelte", "flaiover/src/routes/settings/[page]/+page.ts", flaiover/src/routes/settings/settings.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/operators/settings.md]
after: [S-0458]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 45m
  delivery: 2026-10-22T01:43:00Z
  basis: "Its own forecast of 45m; 31st in the pull order with an in-progress limit of 5, behind S-0453, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433, S-0434, S-0435, S-0456 and S-0458."
  by: flai
  at: 2026-10-21T07:19:47Z
---
# S-0459 The settings page becomes five Configuration pages, Features, Agents, File System, Projects, and Harnesses, each reached from the rail

## Goal

The operator finds a setting under the Configuration page it belongs to, instead of scrolling one page of ten sections. The rail's Configuration group links five pages:

| Page | Sections |
|------|----------|
| Features | Host actions, Checks |
| Agents | Default agent for new stories, Providers, Planning again |
| File System | Shared paths, Import folders |
| Projects | Projects, Tokens |
| Harnesses | Starting agents, one block per harness |

The split is the operator's answer to TH-0553 question 4. Each page shows only its own sections. What each section reads, saves, and gates is unchanged.

## Acceptance criteria

- [ ] `$lib/settingspages.ts` maps each of the five pages to its slug, its title, and its sections, as TH-0553 settles them. Each of `SettingsPanel.svelte`'s ten sections belongs to exactly one page, and a test fails when a section belongs to none or to two.
- [ ] `/settings/<page>` renders the page's title and only its sections, for `features`, `agents`, `file-system`, `projects`, and `harnesses`. An unknown slug is a 404. `/settings` redirects to `/settings/features`, and a link to one of the old `/settings` anchors lands on the page that holds that section.
- [ ] Each section reads, saves, and refuses as it did on the one page. The settings host action's gate, the read-only dashboard, and a host flai that keeps no settings still show their notice on each page. The existing `SettingsPanel` tests pass, run against each section's page.
- [ ] The rail's five Configuration entries link `/settings/<page>`. `locate` names each as the current page, so the right entry is `aria-current` there.
- [ ] `docs/users/flaiover.md` § Settings, `design/system/flaiover-dashboard.md` § Views, and the Settings page references in `docs/operators/settings.md` name the five pages and what each holds.

## Tasks
- T-1959 settingspages.ts maps the five Configuration pages to their sections, and SettingsPanel renders only the sections of the page it is told
- T-1960 /settings/<page> serves each Configuration page, /settings and its old anchors redirect to the page that holds them, and the rail links the five pages
- T-1961 The user, design, and operator docs name the five Configuration pages and what each holds

## Notes

- Drafted by the planner for E-0027. Second, after S-0458, whose rail it points at its pages.
- The sections stay in `SettingsPanel.svelte`, told which page to render, rather than being split into ten components. Their shared load, gate, and save state is one script, and splitting it is not this story's outcome.

### Planning

- **Touches.** None were declared, so every touch is the planner's. All are files, and no folder touch was kept.
  - **Layout:** `flaiover/src/lib/components/SettingsPanel.svelte` (1049 lines, ten `<h2>` sections) and its test. `flaiover/src/lib/settingspages.ts` and its test are new. `flaiover/src/routes/settings/+page.svelte` and a new `+page.ts` handle the redirect. `settings/[page]/+page.svelte` and `+page.ts` are new, with `settings.svelte.test.ts`. `sitemenu.ts` and its test retarget the five entries.
  - **Design:** `docs/users/flaiover.md` § Settings and `design/system/flaiover-dashboard.md` § Views. `docs/operators/settings.md` names "the dashboard's Settings page" for host actions and shared paths. It is also the strongest co-change outside the hubs, at 7%.
  - **Co-change:** the rest of `flai touches suggest`'s list is hub documents and flai's Go code, none of which this story reaches. None was added.
- **Forecast.** 45m against flai's 29m. Every `SettingsPanel` case is re-run on its section's page, and an anchor redirect must run in the browser, because the server never sees the hash. The delivery is flai's play-out, moved by its own 16m raise and S-0458's 43m.
- **Cost of delay.** Not yet set. E-0027 has no inputs, so `flai cod` refuses. The input is asked on TH-0553 question 1.
