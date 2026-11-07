---
id: T-2479
type: task
nature: feature
title: flai kluster deploy, teardown, promote, and deployments are commands, each with --json
status: done
parent: S-0578
owner: alex
created: 2026-11-06T21:43:39Z
updated: 2026-11-07T01:56:41Z
transitions:
  - to: ready
    at: 2026-11-07T01:49:05Z
    by: agent-S-0578
  - to: in-progress
    at: 2026-11-07T01:49:05Z
    by: agent-S-0578
  - to: done
    at: 2026-11-07T01:56:41Z
    by: agent-S-0578
stream: S-0578
tags: [kluster, flai, cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go]
after: [T-2478]
usage:
  source: log
  seconds: 456
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 38
      output: 235
      cache_read: 5686223
      cache_write: 50351
      cost: 2.0467
---
# T-2479 flai kluster deploy, teardown, promote, and deployments are commands, each with --json

## Work

It waits for T-2478, whose teardown, promote, and list the commands call, after T-2477's deploy.

In `flai/cmd/kluster.go`, add four subcommands beside S-0577's:

- `deploy S-nnnn` deploys the story from its worktree, streaming each step's name and printing the record;
- `teardown S-nnnn` tears it down, asking first unless `--yes` is given;
- `promote` deploys main's head to the default namespace;
- `deployments` lists the records in a table.

Each prints its record or records as JSON with `--json`, and refuses a project with no `kluster:` block as S-0577's do. A failing step exits non-zero and names the step, its exit code, and its log. Test each, a failure, and `--json` in `flai/cmd/kluster_test.go` with the fakes.

## Done when

- The four subcommands work as above and are listed in `flai kluster --help`.
- A failure's step, exit code, and log are printed, and tested.
- `flai test` on the two paths passes.

## Notes
