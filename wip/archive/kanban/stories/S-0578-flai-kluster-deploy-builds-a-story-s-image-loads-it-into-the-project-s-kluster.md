---
id: S-0578
type: story
nature: feature
title: flai kluster deploy builds a story's image, loads it into the project's kluster, and runs the project's deploy command in the story's namespace, and teardown and promote remove it and deploy main to the default namespace
status: done
parent: E-0032
owner: alex
created: 2026-11-06T21:39:56Z
updated: 2026-11-07T02:35:38Z
transitions:
  - to: ready
    at: 2026-11-07T00:51:57Z
    by: alex
  - to: in-progress
    at: 2026-11-07T01:32:23Z
    by: agent-S-0578
  - to: review
    at: 2026-11-07T02:32:52Z
    by: agent-S-0578
  - to: done
    at: 2026-11-07T02:35:38Z
    by: orchestrator
tags: [kluster, flai, cli]
topics: [kluster, deployment, testing]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/kluster/kluster.go, docs/operators/settings.md, design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md, design/issues/summary.md]
after: [S-0577]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3650
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 3
      hand_edits: 3
      work: 78
  models:
    - model: claude-opus-5-5
      input: 288
      output: 111511
      cache_read: 42124180
      cache_write: 480351
      cost: 14.1241
  strategic:
    - kind: orchestrator
      seconds: 128
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 36
          output: 257
          cache_read: 1321934
          cache_write: 116521
          cost: 0.3465
        - model: claude-sonnet-5-5
          input: 6
          output: 46
          cache_read: 44228
          cache_write: 33275
          cost: 0.0551
cost_of_delay:
  value: 73.53
  by: planner-E-0032
  at: 2026-11-06T21:50:10Z
forecast:
  duration: 27m
  delivery: 2026-11-07T02:05:00Z
  basis: "Its own forecast of 27m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-11-07T01:32:10Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:31Z
---
# S-0578 flai kluster deploy builds a story's image, loads it into the project's kluster, and runs the project's deploy command in the story's namespace, and teardown and promote remove it and deploy main to the default namespace

## Goal

A story's change can be tried on a kluster without a full build and deploy loop. One command builds the story's image in its worktree, loads it into the project's cluster, and deploys it into the story's own namespace with the project's own deploy command. Two more commands remove the namespace, and deploy main to the default namespace as an acceptance does. flai serve's review deploys and the dashboard's Deployments page run these commands.

## Acceptance criteria

- [x] `flai kluster deploy S-nnnn` runs the block's build command in the story's worktree, loads each image it names into the cluster as S-0576's ADR decides, and runs the deploy command with the namespace, the story, the image tag, and the kube context the ADR names. It records the deployment (story, namespace, commit, images, time, outcome) in a file under `.flai-cache`, and `--json` prints it.
- [x] Each step's output goes to a log file the record names. A step that fails stops the deploy, and flai names the step and its exit code.
- [x] `flai kluster teardown S-nnnn` runs the teardown command, or deletes the namespace when the block names none, and removes the record. `flai kluster promote` deploys main's head from the main checkout into the default namespace and records it.
- [x] `flai kluster deployments` lists the records, the default namespace's first, and `--json` prints them as data the dashboard reads.
- [x] A deploy refuses a story that is not in progress or in review, and a second deploy while one runs for the project.
- [x] The tests run with `kluster`, `docker`, and the project's commands faked. No test needs a real cluster or Docker.
- [x] `design/system/kluster.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated `docs/users/flai-reference.md` describe the commands, the record, and the logs.

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

### As built

- Criterion 3: teardown runs the block's `teardown` when it names one and then always deletes the namespace, as `design/system/kluster.md` § The steps says (step 4), so a block with no `teardown` gets only the delete. Promote builds and deploys main's head from a detached worktree of it under the main checkout, `.flai-cache/kluster/main`, as ADR-0254 §5 says, not from the main checkout's own files, which may hold uncommitted work.
- Records are one file per namespace, `.flai-cache/kluster/deployments/<namespace>.json`, as ADR-0254 §1 says, rather than one `deployments.json`.

### Accepted by the orchestrator

- Verified: 9728de82247905c29e267c4f006dd86631db1f7d
- At: 2026-11-07T02:35:38Z

Verdict: pass — all seven criteria are met; the close-out passed at 74529b1 and again at 9728de8, the same change rebased over a wip-only commit of main, and no convention breach was found.
- 1: flai/internal/kluster/deploy.go, flai/internal/kluster/kluster.go, flai/cmd/kluster.go, flai/internal/kluster/deploy_test.go
- 2: flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go
- 3: flai/internal/kluster/deploy.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go
- 4: flai/internal/kluster/deploy.go, flai/cmd/kluster.go
- 5: flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go
- 6: flai/internal/kluster/deploy_test.go, flai/cmd/kluster_test.go
- 7: design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md
