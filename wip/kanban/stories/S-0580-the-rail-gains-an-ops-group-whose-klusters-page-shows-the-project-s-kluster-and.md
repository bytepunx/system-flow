---
id: S-0580
type: story
nature: feature
title: The rail gains an Ops group whose Klusters page shows the project's kluster and brings it up or takes it down through flai, behind the kluster host action
status: backlog
parent: E-0032
owner: alex
created: 2026-11-06T21:40:27Z
updated: 2026-11-06T21:58:12Z
transitions: []
tags: [kluster, flaiover, hostapi]
topics: [kluster, dashboard, host-actions]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/+server.ts, flaiover/src/routes/api/kluster/up/+server.ts, flaiover/src/routes/api/kluster/down/+server.ts, flaiover/src/routes/api/kluster/kluster.test.ts, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md]
after: [S-0577]
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
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 60
          cache_read: 381639
          cache_write: 53273
          cost: 0.104
cost_of_delay:
  value: 117.1
  by: planner-E-0032
  at: 2026-11-06T21:50:19Z
forecast:
  duration: 43m
  delivery: 2026-11-07T13:14:00Z
  basis: "Its own forecast of 43m; 32nd in the pull order with an in-progress limit of 5, behind S-0545, S-0520, S-0575, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556, S-0557, S-0558, S-0559, S-0560, S-0561, S-0562, S-0564, S-0565, S-0567, S-0568, S-0569, S-0570, S-0571, S-0576, S-0577, S-0578 and S-0579."
  by: flai
  at: 2026-11-06T21:58:12Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:41Z
---
# S-0580 The rail gains an Ops group whose Klusters page shows the project's kluster and brings it up or takes it down through flai, behind the kluster host action

## Goal

The operator manages a project's kluster from the dashboard. A new Ops group in the rail holds the pages for it, as E-0032's outcome asks. Its first page, Klusters, shows whether the project's cluster runs and brings it up or takes it down through flai on the host. S-0577's commands do the work. The Deployments page that a later story adds sits in the same group.

## Acceptance criteria

- [ ] flai's host API gains a read, `kluster.status`, and two writes, `kluster.up` and `kluster.down`, which run what `flai kluster status`, `up`, and `down` run. The writes are refused, with flai's own sentence, while the project's `kluster` host action is off.
- [ ] The rail (`$lib/sitemenu.ts`) gains an Ops group, with a lucide icon, between Design and Documentation. Its first page is Klusters, at `/ops/klusters`. `sitemenu.test.ts` pins the new order.
- [ ] The Klusters page shows the cluster's name, profile, provider, whether it runs, and its age. It has an Up button and a Down button, and Down asks for confirmation first. While an action runs, the page says so, and it shows flai's refusal or error text verbatim.
- [ ] With the host action off, or no `kluster:` block, the page says which, and names the Features page or the block to add, in place of the buttons.
- [ ] Page, route, and host API tests cover the status read, each write, a refusal while the action is off, and the empty states.
- [ ] `design/system/flaiover-dashboard.md` § Views (the rail's groups and the routes) and `docs/users/flaiover.md` describe the Ops group and the Klusters page.

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
