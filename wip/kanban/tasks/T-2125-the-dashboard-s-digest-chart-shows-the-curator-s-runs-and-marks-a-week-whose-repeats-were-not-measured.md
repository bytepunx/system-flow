---
id: T-2125
type: task
nature: feature
title: The dashboard's Digest chart shows the curator's runs and marks a week whose repeats were not measured
status: done
parent: S-0504
owner: alex
created: 2026-10-26T04:10:23Z
updated: 2026-10-26T04:39:42Z
transitions:
  - to: ready
    at: 2026-10-26T04:32:58Z
    by: agent-S-0504
  - to: in-progress
    at: 2026-10-26T04:32:59Z
    by: agent-S-0504
  - to: done
    at: 2026-10-26T04:39:42Z
    by: agent-S-0504
stream: S-0504
tags: [flaiover, digest, charts]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-2123]
usage:
  source: log
  seconds: 403
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 65
      output: 23237
      cache_read: 3364504
      cache_write: 144599
      cost: 2.0349
---
# T-2125 The dashboard's Digest chart shows the curator's runs and marks a week whose repeats were not measured

## Work

Make the Digest chart read what T-2123 sends, so that the dashboard no longer draws an unmeasured repeat count as 0:

- `flaiover/src/lib/viz/charts.ts`: `DigestWeek` gains `curator_runs` and takes `asked_again` as optional; the report's `digest` gains `curator` with `runs` and `last_run`, optional. In `digestChart`, a week without `asked_again` draws its whole `per_story` as `other threads` and no `asked again` part, and its tooltip says the repeats were not measured, because no curator run had reached the week, with the week's `curator_runs`. The line under the chart that gives `digest.entries` gains the curator: `Curator: 3 runs in the window · last run <time>`, or `Curator: never ran on this project`, when `digest.curator` is absent. Cover each in `charts.test.ts`.
- `flaiover/src/lib/viz/explain.ts`: the Digest explainer says what `asked again` means and that a week shows it only once a curator run has reached it. Update `explain.test.ts`.
- `flaiover/src/routes/charts/[kind]/+page.svelte`: the Digest table gains a `curator runs` column, and shows `–` for a week with no `asked_again`; `charts.svelte.test.ts` covers the table and the line under the chart with a report that has `curator` and one that has none.
- `design/system/flaiover-dashboard.md`: the `/charts/digest` row says what the chart now draws, the new tooltip text, the table's column, and the curator line, with the ADR linked. `docs/users/flaiover.md`: the Digest row of the charts table says the same for the user. Bump each `updated`.

A flai on the host older than this release sends no `curator_runs` and no `curator`: treat a missing `curator_runs` as 0 and a missing `curator` as never ran, so the page stays whole.

This waits on T-2123 because the chart reads the JSON that task defines, and it is checked against the flai of this tree.

## Done when

- `pnpm vitest` in `flaiover` passes with the new cases in `charts.test.ts`, `explain.test.ts`, and `charts.svelte.test.ts`.
- `flai test` on the changed paths passes: prettier, eslint, and vitest on them.
- Against the flai of this tree with no `wip/agents/curator.md`, `/charts/digest` draws each week's bar without an `asked again` part, its tooltip says the repeats were not measured, the table shows `–` under asked again and 0 under curator runs, and the line under the chart says the curator never ran.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the chart as it is.

## Notes

Drafted by planner-S-0504. `flaiover/src/lib/viz/palette.ts` is left alone unless the chart needs a new series colour; none is added as planned, since the `asked again` series is only left out of a week.
