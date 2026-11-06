---
id: S-0578
type: story
nature: feature
title: flai kluster deploy builds a story's image, loads it into the project's kluster, and runs the project's deploy command in the story's namespace, and teardown and promote remove it and deploy main to the default namespace
status: backlog
parent: E-0032
owner: alex
created: 2026-11-06T21:39:56Z
updated: 2026-11-06T23:26:38Z
transitions: []
tags: [kluster, flai, cli]
topics: [kluster, deployment, testing]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
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
      seconds: 35
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 62
          cache_read: 381640
          cache_write: 53273
          cost: 0.104
cost_of_delay:
  value: 73.53
  by: planner-E-0032
  at: 2026-11-06T21:50:10Z
forecast:
  duration: 27m
  delivery: 2026-11-07T06:38:00Z
  basis: "Its own forecast of 27m; 26th in the pull order with an in-progress limit of 5, behind S-0576, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556, S-0557, S-0558, S-0559, S-0560, S-0561, S-0562, S-0564, S-0565, S-0567, S-0568, S-0569, S-0570, S-0571 and S-0577."
  by: flai
  at: 2026-11-06T23:26:38Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:31Z
---
# S-0578 flai kluster deploy builds a story's image, loads it into the project's kluster, and runs the project's deploy command in the story's namespace, and teardown and promote remove it and deploy main to the default namespace

## Goal

A story's change can be tried on a kluster without a full build and deploy loop. One command builds the story's image in its worktree, loads it into the project's cluster, and deploys it into the story's own namespace with the project's own deploy command. Two more commands remove the namespace, and deploy main to the default namespace as an acceptance does. flai serve's review deploys and the dashboard's Deployments page run these commands.

## Acceptance criteria

- [ ] `flai kluster deploy S-nnnn` runs the block's build command in the story's worktree, loads each image it names into the cluster as S-0576's ADR decides, and runs the deploy command with the namespace, the story, the image tag, and the kube context the ADR names. It records the deployment (story, namespace, commit, images, time, outcome) in a file under `.flai-cache`, and `--json` prints it.
- [ ] Each step's output goes to a log file the record names. A step that fails stops the deploy, and flai names the step and its exit code.
- [ ] `flai kluster teardown S-nnnn` runs the teardown command, or deletes the namespace when the block names none, and removes the record. `flai kluster promote` deploys main's head from the main checkout into the default namespace and records it.
- [ ] `flai kluster deployments` lists the records, the default namespace's first, and `--json` prints them as data the dashboard reads.
- [ ] A deploy refuses a story that is not in progress or in review, and a second deploy while one runs for the project.
- [ ] The tests run with `kluster`, `docker`, and the project's commands faked. No test needs a real cluster or Docker.
- [ ] `design/system/kluster.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated `docs/users/flai-reference.md` describe the commands, the record, and the logs.

## Tasks

- T-2477 The kluster package deploys a story: it builds in the story's worktree, loads the images, runs the deploy command in the story's namespace, logs each step, and records the deployment under .flai-cache
- T-2478 The kluster package tears a story's namespace down, promotes main's head to the default namespace, and lists the deployments with the default namespace first
- T-2479 flai kluster deploy, teardown, promote, and deployments are commands, each with --json
- T-2480 The kluster design, the CLI design, and the user guide describe deploy, teardown, promote, deployments, the record, and the logs

## Notes

### Planning

- Touches:
  - From the layout: `flai/internal/kluster/deploy.go` and its test, new beside S-0577's package, and `flai/cmd/kluster.go` and its test, the command S-0577 adds.
  - From the design: `design/system/kluster.md` and `design/system/flai-cli.md`, and `docs/users/flai.md` with the generated reference, for the changed commands.
  - `flai touches suggest` listed `docs/operators/index.md` and `docs/operators/settings.md` (23% each). They are left out: this story adds no setting and no host action behaviour, and S-0579 documents flai serve's.
- `after` S-0577: this story reads S-0577's `flai/internal/kluster` package, for the cluster's config and kube context, and the manifest's `kluster:` block type. It extends the `flai kluster` command in `flai/cmd/kluster.go` that S-0577 writes.
- Tasks: one task a layer, since each builds on the last.
  - T-2477 deploy and T-2478 teardown, promote, and list share `deploy.go`.
  - T-2479 (commands) calls both.
  - T-2480 (documents) describes them.
- Forecast 27m: `flai forecast` as it stands, 107 s per unit over 110 large feature stories on claude-opus-5-5, at size 15.
- Cost of delay 73.53 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration. It stands as flai gives it.
