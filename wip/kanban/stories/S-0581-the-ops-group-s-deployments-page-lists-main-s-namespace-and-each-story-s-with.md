---
id: S-0581
type: story
nature: feature
title: The Ops group's Deployments page lists main's namespace and each story's, with its commit, outcome, log, and how to reach it for a demo, and deploys, tears down, or promotes through flai
status: backlog
parent: E-0032
owner: alex
created: 2026-11-06T21:40:41Z
updated: 2026-11-06T22:37:33Z
transitions: []
tags: [kluster, flaiover, hostapi]
topics: [kluster, dashboard, deployment]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/deployments/+server.ts, flaiover/src/routes/api/kluster/deploy/+server.ts, flaiover/src/routes/api/kluster/teardown/+server.ts, flaiover/src/routes/api/kluster/promote/+server.ts, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md]
after: [S-0578, S-0580]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 23
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 44
          cache_read: 263972
          cache_write: 39239
          cost: 0.0725
cost_of_delay:
  value: 114.38
  by: planner-E-0032
  at: 2026-11-06T21:50:24Z
forecast:
  duration: 42m
  delivery: 2026-11-07T14:06:00Z
  basis: "Its own forecast of 42m; 30th in the pull order with an in-progress limit of 5, behind S-0520, S-0575, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556, S-0557, S-0558, S-0559, S-0560, S-0561, S-0562, S-0564, S-0565, S-0567, S-0568, S-0569, S-0570, S-0571, S-0576, S-0577, S-0578, S-0579 and S-0580."
  by: flai
  at: 2026-11-06T22:37:33Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:46Z
---
# S-0581 The Ops group's Deployments page lists main's namespace and each story's, with its commit, outcome, log, and how to reach it for a demo, and deploys, tears down, or promotes through flai

## Goal

The operator demos or tries a story's build from the dashboard, without waiting for a full build and deploy loop. The Deployments page in the Ops group lists what runs on the project's kluster: main in the default namespace, and each story in its own. From there the operator deploys a story, tears one down, or promotes main through flai on the host, and sees how to reach a deployment for a demo.

## Acceptance criteria

- [ ] flai's host API gains a read, `kluster.deployments`, and three writes, `kluster.deploy` (with a story), `kluster.teardown` (with a story), and `kluster.promote`, which run what S-0578's commands run. The writes are refused, with flai's own sentence, while the project's `kluster` host action is off.
- [ ] The rail's Ops group gains Deployments, at `/ops/deployments`, after Klusters. `sitemenu.test.ts` pins it.
- [ ] The page lists the default namespace first, then each story's: the story's ID and title linking its item page, the namespace, the commit, when it was deployed, the outcome, and a link to each step's log. While a deploy runs, the page shows it running and refreshes until it ends.
- [ ] Each row has the actions that fit it. A story in review that is not deployed, or deployed at an older commit, can be deployed. A story's namespace can be torn down, with confirmation first. Main can be promoted. flai's refusal or error text is shown verbatim.
- [ ] Each deployment shows how to reach it for a demo, as S-0576's ADR decides, such as the hosts its services answer on and the port-forward command to run.
- [ ] Page, route, and host API tests cover the list, each write, a refusal while the action is off, a running deploy, and an empty cluster.
- [ ] `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md` describe the Deployments page.

## Tasks

- T-2489 flai's host API answers kluster.deployments and runs kluster.deploy, kluster.teardown, and kluster.promote in the background, refusing the writes while the kluster host action is off
- T-2490 The rail's Ops group gains Deployments at /ops/deployments, after Klusters
- T-2491 The Deployments page lists main's namespace and each story's with its commit, outcome, logs, and demo access, and deploys, tears down, or promotes through the dashboard's API routes
- T-2492 The dashboard design, the dashboard guide, and the kluster design describe the Deployments page

## Notes

### Planning

- Touches, from the layout:
  - `flai/internal/hostapi/kluster.go` and its test, which S-0580 adds.
  - `flaiover/src/lib/sitemenu.ts` and its test.
  - `flaiover/src/lib/server/agent.ts` and `repo.ts`.
  - New files: the routes under `flaiover/src/routes/api/kluster/` (`deployments`, `deploy`, `teardown`, `promote`) and the page under `flaiover/src/routes/ops/deployments/`, each with its test.
- Touches, from the design: `design/system/flaiover-dashboard.md` § Views, `docs/users/flaiover.md`, and `design/system/kluster.md`.
- Touches left out: `flai touches suggest` listed `flai/internal/hostapi/writes.go` (14%). It is left out, since S-0580 merges the kluster specs into the method table and this story adds to `kluster.go` alone.
- `after`:
  - S-0578: this story reads S-0578's deploy, teardown, promote, and list in `flai/internal/kluster/deploy.go`, the deployment record's JSON fields, and its lock, which these writes share with flai serve.
  - S-0580: this story reads S-0580's Ops group in `flaiover/src/lib/sitemenu.ts` and its `flai/internal/hostapi/kluster.go`, which it extends. It also follows the routes under `flaiover/src/routes/api/kluster/` and the `repo.ts` functions S-0580 writes.
- Tasks:
  - Layer 1: T-2489 (host API) and T-2490 (rail). They share no path.
  - Layer 2: T-2491 (page and routes), which waits for both.
  - Layer 3: T-2492 (documents).
- Forecast 42m: `flai forecast` as it stands, 107 s per unit over 110 large feature stories on claude-opus-5-5, at size 23.
- Cost of delay 114.38 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration. It stands as flai gives it.
