---
id: T-2874
type: task
nature: feature
title: kluster.md, flaiover-dashboard.md, flaiover.md, and flai.md say what the Deployments page, kluster.deployments, kluster.checks, and flai kluster deployments do for a story in progress
status: done
parent: S-0663
owner: alex
created: 2026-11-10T08:26:15Z
updated: 2026-11-13T15:19:19Z
transitions:
  - to: ready
    at: 2026-11-13T15:09:06Z
    by: agent-S-0663
  - to: in-progress
    at: 2026-11-13T15:09:06Z
    by: agent-S-0663
  - to: done
    at: 2026-11-13T15:19:19Z
    by: agent-S-0663
stream: S-0663
tags: [dashboard, cli]
touches: [design/system/kluster.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md]
after: [T-2873]
usage:
  source: log
  seconds: 613
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 91
      output: 32620
      cache_read: 7694385
      cache_write: 172099
      cost: 3.1857
---
# T-2874 kluster.md, flaiover-dashboard.md, flaiover.md, and flai.md say what the Deployments page, kluster.deployments, kluster.checks, and flai kluster deployments do for a story in progress

## Work

Write down what T-2868 to T-2873 built, as the code now does it.

- `design/system/kluster.md`: § Commands (S-0577, S-0578), the `flai kluster deployments` row and the record's fields (who asked, the column, the head, the last checks run, the limit and its refusals); § The dashboard (S-0580, S-0581) › The host methods, the `kluster.deployments` fields and a row for `kluster.checks` with its refusals in § Refusals; § The Deployments page, Overview's rows and actions for a story in progress, the count against the limit, and Access for a story in progress. Edit the sections S-0659 to S-0662 left there, and mark nothing built that is not.
- `design/system/flaiover-dashboard.md`: § Views, the Deployments row; § API (S-0011), a row for `/api/kluster/checks`, which `design-routes.test.ts` checks against the route.
- `docs/users/flaiover.md` § Ops › Deployments: what the operator sees and does for a story in progress, and the limit.
- `docs/users/flai.md` § Run a project's kluster: what `flai kluster deployments` now lists.

It waits for T-2873, the last change it describes.

## Done when

- The four documents say all of the story's criteria 1 to 5, and nothing the code does not do.
- `flai test` passes on the changed paths, the markdown lint and `design-routes` among them.

## Notes
