---
id: T-2606
type: task
nature: feature
title: The kluster design, the dashboard design, and the user and operator guides describe components, pending steps, kluster.log, and the Deployments page's rail
status: done
parent: S-0595
owner: alex
created: 2026-11-07T22:08:38Z
updated: 2026-11-07T22:49:19Z
transitions:
  - to: ready
    at: 2026-11-07T22:10:22Z
    by: agent-S-0595
  - to: in-progress
    at: 2026-11-07T22:45:27Z
    by: agent-S-0595
  - to: done
    at: 2026-11-07T22:49:19Z
    by: agent-S-0595
stream: S-0595
tags: [cli, dashboard]
touches: [design/system/kluster.md, design/system/flaiover-dashboard.md, design/system/project-manifest.md, docs/users/flaiover.md, docs/users/flai.md, docs/operators/index.md]
after: [T-2600, T-2604, T-2605]
usage:
  source: log
  seconds: 232
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 22002
      cache_read: 10754786
      cache_write: 84857
      cost: 3.1444
---
# T-2606 The kluster design, the dashboard design, and the user and operator guides describe components, pending steps, kluster.log, and the Deployments page's rail

## Work

ADR-0263, as built by T-2598 to T-2605. In `design/system/kluster.md`: the block table gains `components`, with an example; § The steps and their variables names `build-<name>`, `FLAI_KLUSTER_COMPONENT`, and the per-component images file; the record's step outcomes gain `pending` and the plan written before the first step; § The host methods gains `kluster.log` with its fields; § The Klusters page says step 4 is the shared form; § The Deployments page describes the rail's five panes, the Build pane's panels, the live feed, and the Setup link in place of a refused Deploy. Link ADR-0263 where ADR-0259 is linked. In `design/system/flaiover-dashboard.md` § Views, the Deployments row or section follows. In `design/system/project-manifest.md`, where the `kluster:` block's keys are listed, add `components`. In `docs/users/flaiover.md` § Deployments and § Klusters, say what the operator sees and does; in `docs/users/flai.md`, where the `kluster:` block and `flai manifest set` for its keys are shown, add `components`; in `docs/operators/index.md` near the Deployments paragraphs, the live logs. Bump each file's `updated`. Waits for the code tasks it describes.

## Done when

- Each document says what the code does, with no copy of another's text beyond a link; the markdown lint passes on them with `flai test`.

## Notes
