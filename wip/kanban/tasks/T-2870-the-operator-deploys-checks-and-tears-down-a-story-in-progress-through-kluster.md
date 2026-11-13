---
id: T-2870
type: task
nature: feature
title: The operator deploys, checks, and tears down a story in progress through kluster.deploy, a new kluster.checks, and kluster.teardown, and forwards into its namespace
status: in-progress
parent: S-0663
owner: alex
created: 2026-11-10T08:25:08Z
updated: 2026-11-13T14:25:24Z
transitions:
  - to: ready
    at: 2026-11-13T14:25:23Z
    by: agent-S-0663
  - to: in-progress
    at: 2026-11-13T14:25:24Z
    by: agent-S-0663
stream: S-0663
tags: [cli, dashboard]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/kluster_forward_test.go, flai/internal/hostapi/writes_test.go]
after: [T-2869]
---
# T-2870 The operator deploys, checks, and tears down a story in progress through kluster.deploy, a new kluster.checks, and kluster.teardown, and forwards into its namespace

## Work

Let the dashboard act on a story in progress's namespace, each through flai and behind the `kluster` host action, in `flai/internal/hostapi/kluster.go`.

- `kluster.deploy` takes a story in progress as it takes one in review (`klusterDeployDo`), through the same path S-0660's request runs, so the job counts against the limit, is refused over it with S-0660's sentence, and is recorded as asked by the operator.
- A new write, `kluster.checks`, with `{"story": "S-nnnn"}`, starts the project's checks run against the story's namespace through the path S-0661 built, recorded as asked by the operator. It refuses, in the order `kluster.deploy` does, with the host action off, no `kluster:` block, a story that is not an ID, a namespace not deployed at the story's head (saying to deploy first), and the `checks` the project names being none. It answers the run at once, as the other writes answer their job. Name it `kluster.checks` so the dashboard's route (T-2871) can be written beside this task.
- `kluster.teardown` of a story in progress needs no change but a test: it tears down any story with a record.
- `kluster.forward.start` and `kluster.forward.stop` do not look at the story's column (`flai/internal/hostapi/kluster_forward.go`). Add a test in `kluster_forward_test.go` that a forward into a story in progress's namespace starts and stops as one into a story in review does.
- Add `kluster.checks` to the table of writes and their host action in `writes_test.go`.

It waits for T-2869, which changes `kluster.go` and `kluster_test.go` first.

## Done when

- Tests show `kluster.deploy` and `kluster.checks` on a story in progress recorded as asked by the operator, each refusal of `kluster.checks`, a deploy refused over the limit, a teardown of a story in progress, and a forward into its namespace started and stopped.
- `flai test` passes on the changed paths, the `flaiover-flai` tier among them.

## Notes
