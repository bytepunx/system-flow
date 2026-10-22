---
id: T-2024
type: task
nature: improvement
title: NavRail.svelte draws Charts as an expandable item listing the chart groups and their charts, the chart shown current, and the e2e layout test opens a chart from the rail
status: done
parent: S-0484
owner: alex
created: 2026-10-21T20:58:29Z
updated: 2026-10-22T20:31:01Z
transitions:
  - to: ready
    at: 2026-10-22T20:19:31Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:26:34Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:31:01Z
    by: agent-S-0484
stream: S-0484
tags: [dashboard]
touches: [flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/e2e/layout.e2e.ts]
after: [T-2019, T-2022]
usage:
  source: log
  seconds: 267
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 10608
      cache_read: 1660286
      cache_write: 59525
      cost: 0.9067
---
# T-2024 NavRail.svelte draws Charts as an expandable item listing the chart groups and their charts, the chart shown current, and the e2e layout test opens a chart from the rail

## Work

In `flaiover/src/lib/components/NavRail.svelte`, draw a page with `children` as a button, not a link: its icon and label, `aria-expanded`, `aria-controls="nav-rail-page-<key>"`, `data-rail-page=<key>`, `onclick` `rail.toggle(p.key)`, the icon in the accent colour while expanded and muted while collapsed, and a chevron as a group has. Below it, hidden while collapsed, a nested `<ul>` with one `<li>` per child group: a muted label (`data-rail-label`) that is not a link, then its items as links, each `withProject(item.href, projectState.current)`, `data-rail-chart=<kind>`, bold in the accent colour with `aria-current="page"` while `here.child` is its kind, calling `onnavigate` on click so the drawer closes. The Charts button is marked current, bold in the accent colour, while `here.page === 'charts'`. The arrival effect expands the page as well as the group when `here.child` or `here.page === 'charts'`. Keep `design/components/flaiover/src/lib/components/NavRail.examples.json` in step; its props do not change. In `flaiover/e2e/layout.e2e.ts`, add `/charts/cycle-time` to the pages the rail is checked on, and a case that expands Charts and opens Throughput from the rail, in the docked rail and in the drawer.

Third layer: waits for T-2019, the state that keeps Charts' expansion, and T-2022, the children and `locate`'s `child`.

## Done when

`NavRail.svelte.test.ts` covers Charts as a button with its `aria` attributes, the groups' labels and the charts' links in catalogue order, the current chart and Charts both marked, expansion on arrival at a chart, the stored toggle, and that the existing page-count assertions still hold with Charts' links counted apart. `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. Charts links nowhere itself: the operator picks a chart, as the epic asks. `/charts` still redirects to cycle time for a typed address.
