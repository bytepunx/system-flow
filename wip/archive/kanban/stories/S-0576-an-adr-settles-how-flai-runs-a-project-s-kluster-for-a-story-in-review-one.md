---
id: S-0576
type: story
nature: research
title: "An ADR settles how flai runs a project's kluster for a story in review: one shared cluster, a namespace per story, the story's own settings, promotion on acceptance, and the manifest block a project fills in"
status: done
parent: E-0032
owner: alex
created: 2026-11-06T21:39:23Z
updated: 2026-11-07T00:13:07Z
transitions:
  - to: ready
    at: 2026-11-06T23:24:57Z
    by: alex
  - to: in-progress
    at: 2026-11-06T23:25:06Z
    by: agent-S-0576
  - to: review
    at: 2026-11-07T00:11:38Z
    by: agent-S-0576
  - to: done
    at: 2026-11-07T00:13:07Z
    by: orchestrator
blocked:
  - from: 2026-11-06T23:30:14Z
    until: 2026-11-06T23:57:21Z
    reason: awaiting the operator on TH-0832
  - from: 2026-11-07T00:02:11Z
    until: 2026-11-07T00:09:47Z
    reason: awaiting the operator on TH-0832
tags: [kluster, adr]
topics: [kluster, testing, deployment, manifest]
touches: [design/adrs, design/adrs/README.md, design/system/kluster.md, design/system/README.md, design/conventions/tooling.md, design/issues/I-0266-flai-check-finds-issues-summary-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [S-0575]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 710
  priced_by: harness
  turns:
    - day: 2026-11-06
      ceremony: 2
      hand_edits: 1
      work: 34
    - day: 2026-11-07
      ceremony: 1
      hand_edits: 2
      work: 17
  models:
    - model: claude-opus-5-5
      input: 136
      output: 49392
      cache_read: 14569521
      cache_write: 296546
      cost: 6.2279
  strategic:
    - kind: orchestrator
      seconds: 190
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 73
          output: 486
          cache_read: 2276491
          cache_write: 368781
          cost: 0.6358
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 12355
          cache_write: 32415
          cost: 0.0318
cost_of_delay:
  value: 163.4
  by: planner-E-0032
  at: 2026-11-06T21:50:00Z
forecast:
  duration: 1h
  delivery: 2026-11-07T06:23:00Z
  basis: "Its own forecast of 1h; 25th in the pull order with an in-progress limit of 5, behind S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556, S-0557, S-0558, S-0559, S-0560, S-0561, S-0562, S-0564, S-0565, S-0567, S-0568, S-0569, S-0570 and S-0571."
  by: flai
  at: 2026-11-06T23:20:09Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:21Z
---
# S-0576 An ADR settles how flai runs a project's kluster for a story in review: one shared cluster, a namespace per story, the story's own settings, promotion on acceptance, and the manifest block a project fills in

## Goal

The stories after this one build kluster support into flai and the dashboard. They need one decision to build on. The decision covers how a project tells flai about its cluster, its build, and its deploy; how a story in review gets a namespace of its own over shared infrastructure; and what happens to that namespace when the story leaves review or is accepted. This story records that decision in an ADR and a design document, from what S-0575's experiment found. It also brings this repository's tooling convention in line, which today says kluster is not used here.

## Acceptance criteria

- [x] An ADR under `design/adrs/`, accepted by the operator on a thread, decides:
  - one cluster per project, made from the kluster config file the project names, with the shared infrastructure its profile installs;
  - the namespace a story in review deploys into, named from its serial, and the default namespace main deploys into;
  - the `kluster:` block of `system-flow.yaml`: the config path, the build command, the deploy command, the teardown command, the default namespace, and which are optional;
  - the variables flai gives those commands, such as the story, the namespace, the image tag, and the kube context;
  - how a locally built image reaches the cluster while kluster has no import command;
  - what acceptance does: tear down the story's namespace and deploy main to the default namespace;
  - the host action that gates it, `kluster`, off by default, and what the review checks see while it is on.
- [x] The ADR's Consequences list each change it would ask of kluster's own repository, such as `status --json` or an image import. They are proposals for that repository, not work items here.
- [x] `design/system/kluster.md`, a new design document, describes how flai runs a project's kluster as the ADR decides. It is indexed in `design/system/README.md`.
- [x] The project addition in `design/conventions/tooling.md` that says k3d and kluster are not used here is reworded. system-flow itself still runs no cluster, and flai drives a project's kluster behind the `kluster` host action.
- [x] `flai check --strict` and the markdown lint pass on every file the story changes.

## Tasks

- T-2469 A proposed ADR decides the shared cluster, the story and default namespaces, the manifest's kluster block, the commands' variables, the image load, acceptance, and the kluster host action, and the operator accepts it on a thread
- T-2470 design/system/kluster.md describes how flai runs a project's kluster as the ADR decides, and the design index lists it
- T-2471 The tooling convention's project addition no longer says kluster is not used here: system-flow runs no cluster itself, and flai drives a project's kluster behind the kluster host action

## Notes

### Planning

- Touches:
  - The folder `design/adrs` is kept. The ADR's file name takes the number `flai adr new` gives it when the story runs, so no task can name it now.
  - `design/adrs/README.md` comes from the layout: flai generates the ADR index (ADR-0228).
  - `design/system/kluster.md` comes from the design: criterion 3 names it.
  - `design/system/README.md` comes from the layout: it is the design folder's index.
  - `design/conventions/tooling.md` comes from the design: criterion 4 names its project addition, which says k3d and kluster are not used here.
  - `flai touches suggest` also listed `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md`, and the dashboard documents, at 27% or less. They are left out: this story changes no command or page, and S-0577 to S-0581 take those documents.
- `after` S-0575: the ADR is decided from S-0575's results document, `design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md`. It reads the per-namespace settings, the times, how images were loaded, and kluster's gaps from there.
- Tasks: T-2469 (the ADR) is the first layer. T-2470 (the design document) and T-2471 (the convention) both wait for the accepted ADR. They share no path, so they form the second layer together.
- Forecast 1h. `flai forecast` gave 45m: 531 s per unit over three small research stories, at size 5, counted before the five touches were set. With the touches, size 10 would give 1h28m. It is held at 1h, since two of the touches are a generated index and one line of a convention. The operator's time to accept the ADR is not agent time and is not counted.
- Cost of delay 163.40 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration over the epic's seven open stories. It stands as flai gives it.

### Accepted by the orchestrator

- Verified: 2ccd15b269ef29a2964f1816465068a01a578b7d
- At: 2026-11-07T00:13:07Z

Verdict: pass, at commit 2ccd15b269ef29a2964f1816465068a01a578b7d; the operator accepted ADR-0254 on TH-0832.

- 1: design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md, design/adrs/README.md
- 2: design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md
- 3: design/system/kluster.md, design/system/README.md
- 4: design/conventions/tooling.md
- 5: design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md, design/adrs/README.md, design/conventions/tooling.md, design/system/README.md, design/system/kluster.md
