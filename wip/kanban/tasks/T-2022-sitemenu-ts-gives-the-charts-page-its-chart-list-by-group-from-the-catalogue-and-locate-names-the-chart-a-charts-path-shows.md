---
id: T-2022
type: task
nature: improvement
title: sitemenu.ts gives the Charts page its chart list by group from the catalogue, and locate names the chart a /charts path shows
status: done
parent: S-0484
owner: alex
created: 2026-10-21T20:57:53Z
updated: 2026-10-22T20:26:23Z
transitions:
  - to: ready
    at: 2026-10-22T20:19:31Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:23:30Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:26:23Z
    by: agent-S-0484
stream: S-0484
tags: [dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts]
after: [T-2018]
usage:
  source: log
  seconds: 173
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 7841
      cache_read: 1227244
      cache_write: 43999
      cost: 0.6702
---
# T-2022 sitemenu.ts gives the Charts page its chart list by group from the catalogue, and locate names the chart a /charts path shows

## Work

In `flaiover/src/lib/sitemenu.ts`, let a `LinkedPage` carry optional `children`: a list of labelled groups, each `{ key, label, items: { key, label, href }[] }`, and give the Charts page its children from `CHART_GROUPS` and `TITLES` in `viz/catalog.ts`, each item's `href` `resolve('/charts/[kind]', { kind })`. Keep the page's own `href` at cycle time for `locate`'s callers and the tests that pin hrefs. Make `locate` return a third field, `child`, the chart kind when the path is `/charts/<kind>` and the kind is in `KINDS`, else `null`, with `group` and `page` as before.

Second layer: waits for T-2018, which gives it the catalogue to import without the builders.

## Done when

`sitemenu.test.ts` pins that Charts' children are the four groups in order with every kind once, that each item's href is `/charts/<kind>`, and that `locate('/charts/cycle-time')` gives `{ group: 'flow', page: 'charts', child: 'cycle-time' }`, `locate('/charts/nope')` gives `child: null`, and other pages give `child: null`. `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. `sitemenu.ts` imports from `viz/catalog.ts`, never from `viz/charts.ts`, so the layout's bundle stays free of the builders.
