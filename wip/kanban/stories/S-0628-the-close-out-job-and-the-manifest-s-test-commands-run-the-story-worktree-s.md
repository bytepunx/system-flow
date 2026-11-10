---
id: S-0628
type: story
nature: improvement
title: The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host
status: in-progress
owner: alex
created: 2026-11-07T23:40:34Z
updated: 2026-11-10T05:52:09Z
transitions:
  - to: ready
    at: 2026-11-09T22:04:53Z
    by: alex
  - to: in-progress
    at: 2026-11-10T05:52:09Z
    by: agent-S-0628
tags: [flai, security]
topics: [harness, docker, verify]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/container.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, design/adrs/0272-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in-that-container.md, design/adrs/README.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/issues/I-0295-the-close-out-job-and-the-manifest-s-test-commands-run-the-story-worktree-s.md, design/issues/summary.md]
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
      seconds: 41
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 38
          cache_read: 193878
          cache_write: 6360
          cost: 0.0533
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
forecast:
  duration: 2h
  delivery: 2026-11-10T02:07:00Z
  basis: "Its own forecast of 2h; 8th in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0633, S-0634, S-0645, S-0644, S-0638, S-0636, S-0640, S-0625 and S-0630."
  by: flai
  at: 2026-11-09T22:33:17Z
finalized:
  by: alex
  at: 2026-11-09T22:04:48Z
---
# S-0628 The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host

## Goal

This story remediates [I-0295](../../../design/issues/I-0295-the-close-out-job-and-the-manifest-s-test-commands-run-the-story-worktree-s.md), "The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0295 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0295 is closed with `flai issue close I-0295 --reason` saying what fixed it

## Tasks
- T-2772 Record in an ADR that a container story's close-out runs its verify in that container, and what still runs on the host
- T-2775 Give the container adapter a docker run argv for one command in a story's container with its mounts
- T-2777 Run a container story's close-out job with its verify in the story's container and its commit and sync on the host through flai alone
- T-2782 Run the operator's checks for a container story in the story's container
- T-2784 Say in the design and the operators' and users' guides what a container story runs inside and what still runs on the host
- T-2785 Close I-0295 with what fixed it

## Notes

The analyzer found I-0295: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0295: analyzer 41s, 0.8156 USD.

### Planning

**Remedy proposed.** For a story whose agent runs on the `container` harness, the close-out job runs `flai verify --record-issues` in the story's container, with the story's mounts. The operator's checks for that story run there too. The commit and the sync stay on the host, through flai's own code and git, with the main checkout's hooks. What still runs on the host by design is named in an ADR and in the guides. The plan's thread asks the operator to confirm this remedy over documenting the boundary alone.

**Touches.** No touches were declared. `flai touches suggest` with three seed paths listed eleven co-changing files. I kept the three test files beside their code. I left out `flai/cmd/checks.go`, `harness.go`, `agents.go`, and `serve.go`, at 18% each: no task changes them. A task that does will widen the touches when it closes.

| Touch | From |
|-------|------|
| `flai/internal/harness/container.go`, `flai/internal/serve/verifyjob.go`, `flai/internal/serve/checks.go` | design: the issue's Impact names them |
| `flai/internal/harness/container_test.go`, `flai/internal/serve/verifyjob_test.go`, `flai/internal/serve/checks_test.go` | co-change, 55%, 18%, and 18% |
| `flai/internal/serve/container.go` | layout: `containerImageHost` decides the story's image |
| `design/adrs/0272-...md`, `design/adrs/README.md` | layout: the next free ADR number and its generated index; the slug is a prediction |
| `design/system/flai-cli.md`, `design/system/strategic-agents.md`, `docs/operators/index.md`, `docs/operators/runbooks/agent-container.md`, `docs/users/flai.md` | design: the sections that describe the container harness and the close-out job, which the issue says no document covers |
| the I-0295 file, `design/issues/summary.md` | criterion 2 |

There are no folder touches.

**Forecast.** `flai forecast` gave 30m: size 18, at 98 s per unit, over 142 done large-band improvement stories. I raised it to 2h. For a container story the job reimplements the close-out script's commit step and its three sync rounds in Go, and its tests fake docker. S-0618, a container story of the same kind, was forecast at 2h. Delivery is a day out, since the story is a draft until the operator finalizes it.

**Cost of delay.** None yet. The story has no inputs and no epic, so `flai cod` refuses it. The plan's thread asks the operator for an input. I will work out the value once it is set.

**After.** None. The story reads no other open story's output.

**Out of scope.** The analysis also notes the symlink hygiene in `flai/internal/docedit/file.go`, `FileReach` and `readText`: they check paths as strings and then follow symlinks. That is another cause. The plan's thread offers it as a story of its own.
