---
id: S-0577
type: story
nature: feature
title: flai kluster up, down, and status run a project's kluster from the manifest's kluster block, flai check validates the block, and a kluster host action, off by default, appears on the Features page
status: backlog
parent: E-0032
owner: alex
created: 2026-11-06T21:39:42Z
updated: 2026-11-06T21:57:57Z
transitions: []
tags: [kluster, flai, cli]
topics: [kluster, manifest, host-actions]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/kluster_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/root.go, design/system/project-manifest.md, design/system/flai-cli.md, design/system/kluster.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, design/tech/kluster.md, design/tech/README.md]
after: [S-0576]
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
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 69
          cache_read: 416897
          cache_write: 53390
          cost: 0.1125
cost_of_delay:
  value: 128
  by: planner-E-0032
  at: 2026-11-06T21:50:05Z
forecast:
  duration: 47m
  delivery: 2026-11-07T12:22:00Z
  basis: "Its own forecast of 47m; 29th in the pull order with an in-progress limit of 5, behind S-0545, S-0520, S-0575, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556, S-0557, S-0558, S-0559, S-0560, S-0561, S-0562, S-0564, S-0565, S-0567, S-0568, S-0569, S-0570, S-0571 and S-0576."
  by: flai
  at: 2026-11-06T21:57:57Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:26Z
---
# S-0577 flai kluster up, down, and status run a project's kluster from the manifest's kluster block, flai check validates the block, and a kluster host action, off by default, appears on the Features page

## Goal

A project such as signet or authstar names its kluster in `system-flow.yaml`, and flai brings that cluster up, takes it down, and reports it. The story uses the `kluster:` block and the host action that S-0576's ADR decides. It is the base the deploy commands, flai serve's review deploys, and the dashboard's Ops pages build on.

## Acceptance criteria

- [ ] `system-flow.yaml` takes the `kluster:` block as the ADR from S-0576 names it. flai's manifest reads it, and `flai check` reports an invalid block, such as a config path that does not exist or an unknown key, as a finding of its own with the key and the reason.
- [ ] `flai kluster up` runs `kluster up` with the block's config. `flai kluster down` runs `kluster down` without kluster's own prompt. `flai kluster status` reports the cluster's name, whether it runs, and its age, parsed from `kluster status`, and `--json` prints them as an object.
- [ ] A missing `kluster` binary, or a project with no `kluster:` block, is refused with a sentence that says how to install kluster or which block to add, not a stack trace.
- [ ] A host action `kluster`, off by default, is in flai's list of host actions with a plain description. `flai serve enable kluster` turns it on for a project, and the Features page lists it with no change to the dashboard.
- [ ] The tests run with `kluster` faked on `PATH`. No test needs a real cluster or Docker.
- [ ] `design/system/project-manifest.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, the generated `docs/users/flai-reference.md`, and `docs/operators/settings.md` describe the block, the commands, and the host action. `design/tech/kluster.md`, a new tech file indexed in `design/tech/README.md`, records kluster, the version tested, where flai uses it, and why.

## Tasks

- T-2472 flai's manifest reads the kluster block the ADR names, and flai check reports an invalid one with the key and the reason
- T-2473 A kluster host action, off by default and seen by the dashboard, is in flai's list of host actions with a plain description
- T-2474 A kluster package runs kluster up, down, and status from the manifest's block, parses status into data, and says how to install kluster when it is missing
- T-2475 flai kluster up, down, and status are commands, status with --json, refusing a project with no kluster block with the block to add
- T-2476 The manifest, CLI, user, operator, and tech documents describe the kluster block, the three commands, and the kluster host action

## Notes

### Planning

- Touches, from the layout:
  - `flai/internal/manifest/manifest.go` and its test.
  - `flai/internal/check/check.go`, with a new `kluster_test.go`.
  - `flai/internal/hostapi/writes.go` and its test: the `Actions` map and `DashboardSees`, which the Features page lists.
  - `flai/cmd/root.go`, where commands are registered.
  - New files: `flai/internal/kluster/kluster.go`, `flai/cmd/kluster.go`, and their tests.
- Touches, from the design:
  - `design/system/project-manifest.md`, `design/system/flai-cli.md`, and `design/system/kluster.md`.
  - `docs/users/flai.md` and the generated `docs/users/flai-reference.md`, as `documentation.md`'s project additions ask for a changed command.
  - `docs/operators/settings.md`, which must index every manifest key.
  - `docs/operators/index.md`, where host actions are documented.
  - `design/tech/kluster.md` and `design/tech/README.md`, since a new tool flai calls needs a tech file.
- Touches left out:
  - `flai touches suggest` listed `flai/cmd/serve_actions.go` (6%). It is left out: `flai serve enable` takes any action in the `Actions` map.
  - It also listed `flaiover/src/lib/server/agent.ts` (6%). It is left out: the Features page lists what `DashboardSees` returns, so the dashboard needs no change.
- `after` S-0576: this story reads, from S-0576's accepted ADR and `design/system/kluster.md`, the `kluster:` block's keys and defaults, the kube context the commands use, and the host action's name, `kluster`. It also fills in a section of that design document.
- Tasks:
  - Layer 1: T-2472 (manifest and check) and T-2473 (host action). They share no path.
  - Layer 2: T-2474 (package), which waits for T-2472's manifest type.
  - Layer 3: T-2475 (commands), which waits for T-2474's package.
  - Layer 4: T-2476 (documents), which waits for T-2473 and T-2475.
- Forecast 47m: `flai forecast` as it stands, 107 s per unit over 110 large feature stories on claude-opus-5-5, at size 26.
- Cost of delay 128.00 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration. It stands as flai gives it.
