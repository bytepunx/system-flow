---
id: S-0581
type: story
nature: feature
title: The Ops group's Deployments page lists main's namespace and each story's, with its commit, outcome, log, and how to reach it for a demo, and deploys, tears down, or promotes through flai
status: done
parent: E-0032
owner: alex
created: 2026-11-06T21:40:41Z
updated: 2026-11-07T03:14:42Z
transitions:
  - to: ready
    at: 2026-11-07T00:52:07Z
    by: alex
  - to: in-progress
    at: 2026-11-07T02:35:50Z
    by: agent-S-0581
  - to: review
    at: 2026-11-07T03:13:08Z
    by: agent-S-0581
  - to: done
    at: 2026-11-07T03:14:42Z
    by: orchestrator
tags: [kluster, flaiover, hostapi]
topics: [kluster, dashboard, deployment]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/deployments/+server.ts, flaiover/src/routes/api/kluster/deploy/+server.ts, flaiover/src/routes/api/kluster/teardown/+server.ts, flaiover/src/routes/api/kluster/promote/+server.ts, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md, flai/internal/kluster/deploy.go]
after: [S-0578, S-0580]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2263
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 1
      hand_edits: 2
      work: 39
  models:
    - model: claude-opus-5-5
      input: 208
      output: 108275
      cache_read: 28961228
      cache_write: 473552
      cost: 11.2831
  strategic:
    - kind: orchestrator
      seconds: 40
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 13
          output: 88
          cache_read: 530372
          cache_write: 56190
          cost: 0.141
cost_of_delay:
  value: 114.38
  by: planner-E-0032
  at: 2026-11-06T21:50:24Z
forecast:
  duration: 42m
  delivery: 2026-11-07T03:25:00Z
  basis: "Its own forecast of 42m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-11-07T02:35:38Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:46Z
---
# S-0581 The Ops group's Deployments page lists main's namespace and each story's, with its commit, outcome, log, and how to reach it for a demo, and deploys, tears down, or promotes through flai

## Goal

The operator demos or tries a story's build from the dashboard, without waiting for a full build and deploy loop. The Deployments page in the Ops group lists what runs on the project's kluster: main in the default namespace, and each story in its own. From there the operator deploys a story, tears one down, or promotes main through flai on the host, and sees how to reach a deployment for a demo.

## Acceptance criteria

- [x] flai's host API gains a read, `kluster.deployments`, and three writes, `kluster.deploy` (with a story), `kluster.teardown` (with a story), and `kluster.promote`, which run what S-0578's commands run. The writes are refused, with flai's own sentence, while the project's `kluster` host action is off.
- [x] The rail's Ops group gains Deployments, at `/ops/deployments`, after Klusters. `sitemenu.test.ts` pins it.
- [x] The page lists the default namespace first, then each story's: the story's ID and title linking its item page, the namespace, the commit, when it was deployed, the outcome, and a link to each step's log. While a deploy runs, the page shows it running and refreshes until it ends.
- [x] Each row has the actions that fit it. A story in review that is not deployed, or deployed at an older commit, can be deployed. A story's namespace can be torn down, with confirmation first. Main can be promoted. flai's refusal or error text is shown verbatim.
- [x] Each deployment shows how to reach it for a demo, as S-0576's ADR decides, such as the hosts its services answer on and the port-forward command to run.
- [x] Page, route, and host API tests cover the list, each write, a refusal while the action is off, a running deploy, and an empty cluster.
- [x] `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md` describe the Deployments page.

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

### Accepted by the orchestrator

- Verified: 409d27b117c88e4fc443e1c9ca1bf5d050795525
- At: 2026-11-07T03:14:42Z

Verdict: pass — all seven criteria are met, the close-out passed at 409d27b, and every changed file is inside the story's touches.
- 1: flai/internal/hostapi/kluster.go, flai/internal/kluster/deploy.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/deployments/+server.ts, flaiover/src/routes/api/kluster/deploy/+server.ts, flaiover/src/routes/api/kluster/teardown/+server.ts, flaiover/src/routes/api/kluster/promote/+server.ts
- 2: flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/routes/ops/deployments/+page.svelte
- 3: flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/lib/server/repo.ts, flai/internal/hostapi/kluster.go
- 4: flaiover/src/routes/ops/deployments/+page.svelte
- 5: flai/internal/hostapi/kluster.go, flaiover/src/routes/ops/deployments/+page.svelte, design/system/kluster.md
- 6: flai/internal/hostapi/kluster_test.go, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/lib/sitemenu.test.ts
- 7: design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md
