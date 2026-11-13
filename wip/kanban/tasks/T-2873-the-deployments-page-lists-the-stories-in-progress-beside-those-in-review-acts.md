---
id: T-2873
type: task
nature: feature
title: The Deployments page lists the stories in progress beside those in review, acts on them through DeploymentRow, and shows the count against the limit and the refused request
status: done
parent: S-0663
owner: alex
created: 2026-11-10T08:25:53Z
updated: 2026-11-13T15:07:48Z
transitions:
  - to: ready
    at: 2026-11-13T15:01:10Z
    by: agent-S-0663
  - to: in-progress
    at: 2026-11-13T15:01:10Z
    by: agent-S-0663
  - to: done
    at: 2026-11-13T15:07:48Z
    by: agent-S-0663
stream: S-0663
tags: [dashboard]
touches: [flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts]
after: [T-2870, T-2872]
usage:
  source: log
  seconds: 398
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 68
      output: 333
      cache_read: 3759106
      cache_write: 131912
      cost: 1.456
---
# T-2873 The Deployments page lists the stories in progress beside those in review, acts on them through DeploymentRow, and shows the count against the limit and the refused request

## Work

Change `flaiover/src/routes/ops/deployments/+page.svelte` so the operator sees and clears what agents have running.

- Overview: main's namespace first, then each story namespace through `DeploymentRow` (T-2872), in progress and in review alike, then the stories in progress or in review not yet deployed, each with its column and **Deploy**. Its handlers call the routes for `kluster.deploy`, `kluster.checks` (T-2871), and `kluster.teardown`, and show flai's refusals word for word.
- A badge in the row above the rail, and a line in Overview, say how many story namespaces live against the limit. When a request was refused for the limit, Overview names it: the story, who asked, and when.
- Access: the forwards of a story in progress's namespace start and stop as those of a story in review do. Widen to stories in progress what holds the page to review today: `redeployable`, `deployButton`, and the review-only rows.
- Logs opens at a deploy's or a checks run's step from the row's link.
- Keep the page's polling while a job runs, and the other panes, as they are.
- Update `deployments.svelte.test.ts` for each of these.

It waits for T-2870, the host writes it calls, and T-2872, the row it shows.

## Done when

- The page's tests show a story in progress's row with its column, head, deploy, checks, and who asked; **Deploy**, **Run checks**, and **Teardown** calling their routes; the count against the limit and a refused request; and a forward into a story in progress's namespace.
- `flai test` passes on the changed paths, `vitest`, `svelte-check`, and `design-routes` among them.

## Notes
