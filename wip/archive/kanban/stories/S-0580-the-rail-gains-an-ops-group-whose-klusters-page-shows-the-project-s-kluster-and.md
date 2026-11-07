---
id: S-0580
type: story
nature: feature
title: The rail gains an Ops group whose Klusters page shows the project's kluster and brings it up or takes it down through flai, behind the kluster host action
status: done
parent: E-0032
owner: alex
created: 2026-11-06T21:40:27Z
updated: 2026-11-07T02:25:41Z
transitions:
  - to: ready
    at: 2026-11-07T00:52:04Z
    by: alex
  - to: in-progress
    at: 2026-11-07T01:32:36Z
    by: agent-S-0580
  - to: review
    at: 2026-11-07T02:23:36Z
    by: agent-S-0580
  - to: done
    at: 2026-11-07T02:25:41Z
    by: orchestrator
tags: [kluster, flaiover, hostapi]
topics: [kluster, dashboard, host-actions]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/+server.ts, flaiover/src/routes/api/kluster/up/+server.ts, flaiover/src/routes/api/kluster/down/+server.ts, flaiover/src/routes/api/kluster/kluster.test.ts, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md, flaiover/src/lib/components/NavRail.svelte.test.ts, design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md, design/issues/summary.md, design/issues/I-0267-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md]
after: [S-0577]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3079
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 4
      hand_edits: 2
      work: 55
  models:
    - model: claude-opus-5-5
      input: 308
      output: 92252
      cache_read: 42761275
      cache_write: 514393
      cost: 13.9359
  strategic:
    - kind: orchestrator
      seconds: 154
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 303
          cache_read: 1344267
          cache_write: 189013
          cost: 0.3693
        - model: claude-sonnet-5-5
          input: 4
          output: 62
          cache_read: 12601
          cache_write: 29773
          cost: 0.0301
cost_of_delay:
  value: 117.1
  by: planner-E-0032
  at: 2026-11-06T21:50:19Z
forecast:
  duration: 43m
  delivery: 2026-11-07T02:24:00Z
  basis: "Its own forecast of 43m; 3rd in the pull order with an in-progress limit of 5, behind S-0578 and S-0579."
  by: flai
  at: 2026-11-07T01:32:23Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:41Z
---
# S-0580 The rail gains an Ops group whose Klusters page shows the project's kluster and brings it up or takes it down through flai, behind the kluster host action

## Goal

The operator manages a project's kluster from the dashboard. A new Ops group in the rail holds the pages for it, as E-0032's outcome asks. Its first page, Klusters, shows whether the project's cluster runs and brings it up or takes it down through flai on the host. S-0577's commands do the work. The Deployments page that a later story adds sits in the same group.

## Acceptance criteria

- [x] flai's host API gains a read, `kluster.status`, and two writes, `kluster.up` and `kluster.down`, which run what `flai kluster status`, `up`, and `down` run. The writes are refused, with flai's own sentence, while the project's `kluster` host action is off.
- [x] The rail (`$lib/sitemenu.ts`) gains an Ops group, with a lucide icon, between Design and Documentation. Its first page is Klusters, at `/ops/klusters`. `sitemenu.test.ts` pins the new order.
- [x] The Klusters page shows the cluster's name, profile, provider, whether it runs, and its age. It has an Up button and a Down button, and Down asks for confirmation first. While an action runs, the page says so, and it shows flai's refusal or error text verbatim.
- [x] With the host action off, or no `kluster:` block, the page says which, and names the Features page or the block to add, in place of the buttons.
- [x] Page, route, and host API tests cover the status read, each write, a refusal while the action is off, and the empty states.
- [x] `design/system/flaiover-dashboard.md` § Views (the rail's groups and the routes) and `docs/users/flaiover.md` describe the Ops group and the Klusters page.

## Tasks

- T-2485 flai's host API answers kluster.status and runs kluster.up and kluster.down in the background, refusing the writes while the kluster host action is off
- T-2486 The rail's site menu gains an Ops group, between Design and Documentation, whose first page is Klusters at /ops/klusters
- T-2487 The Klusters page shows the project's cluster, brings it up or takes it down through the dashboard's API routes, and says when the host action is off or the block is missing
- T-2488 The dashboard design, the dashboard guide, and the kluster design describe the Ops group and the Klusters page

## Notes

### Planning

- Touches, from the layout:
  - `flai/internal/hostapi/kluster.go` and its test, new beside `components.go`. `flai/internal/hostapi/writes.go` merges the specs into the method table.
  - `flai/cmd/serve_actions.go`: the components writes run in flai serve's own process, wired there, and so will the kluster jobs.
  - `flaiover/src/lib/sitemenu.ts` and its test: the rail's groups.
  - `flaiover/src/lib/navrail.svelte.ts`: its `keys` list remembers a group's expanded state.
  - `flaiover/src/lib/server/agent.ts` and `repo.ts`: the method allow-list and the client.
  - New files: the routes under `flaiover/src/routes/api/kluster/` and the page under `flaiover/src/routes/ops/klusters/`, each with its test, named as `routes/api/components` and `routes/design/components` are.
- Touches, from the design: `design/system/flaiover-dashboard.md` § Views, `docs/users/flaiover.md`, and `design/system/kluster.md`.
- Touches left out: `flai touches suggest` listed `design/system/flai-cli.md` (32%), `docs/users/flai.md` (30%), and `flai/internal/hostapi/writes_test.go` (20%). The story adds no command, and its host API tests are in `kluster_test.go`.
- `after` S-0577: this story reads S-0577's `flai/internal/kluster` package (up, down, and status) and its `kluster` host action in `flai/internal/hostapi/writes.go`, which gates the writes.
- Tasks:
  - Layer 1: T-2485 (host API) and T-2486 (rail). They share no path.
  - Layer 2: T-2487 (page and routes), which waits for both.
  - Layer 3: T-2488 (documents).
- Forecast 43m: `flai forecast` as it stands, 107 s per unit over 110 large feature stories on claude-opus-5-5, at size 24.
- Cost of delay 117.10 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration. It stands as flai gives it.

### Accepted by the orchestrator

- Verified: 53751382c278486e7dc7663ac794bd1032c99810
- At: 2026-11-07T02:25:41Z

Verdict: pass at 53751382c278486e7dc7663ac794bd1032c99810
- 1: flai/internal/hostapi/kluster.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go, flai/internal/hostapi/kluster_test.go
- 2: flaiover/src/lib/sitemenu.ts, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/NavRail.svelte.test.ts
- 3: flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/api/kluster/+server.ts, flaiover/src/routes/api/kluster/up/+server.ts, flaiover/src/routes/api/kluster/down/+server.ts, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts
- 4: flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/lib/server/repo.ts
- 5: flai/internal/hostapi/kluster_test.go, flaiover/src/routes/api/kluster/kluster.test.ts, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts
- 6: design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md
