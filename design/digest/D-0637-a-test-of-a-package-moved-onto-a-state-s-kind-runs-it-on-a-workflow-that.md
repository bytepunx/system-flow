---
id: D-0637
title: A test of a package moved onto a state's kind runs it on a workflow that renames in-progress and review with no alias and adds an active state between them, as serve's flow_test.go renamedFlow does, so that no comparison by name passes by accident
kind: decision
sources: [S-0743, S-0739, S-0747, S-0838, I-0399]
topics: [flai, go, code, flaiover, dashboard]
paths: [flai/internal/serve/flow_test.go, flai/internal/workitem/rules_test.go, flai/internal/workitem/workflow_fixture_test.go, flai/internal/hostapi/flow.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/lanes.fixture.ts, "flaiover/src/routes/api/items/[id]/move/kinds.test.ts"]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0637 A test of a package moved onto a state's kind runs it on a workflow that renames in-progress and review with no alias and adds an active state between them, as serve's flow_test.go renamedFlow does, so that no comparison by name passes by accident

S-0743 moved flai serve's kluster deploys, forge pass, and merge queue off the names `review` and `done` and onto a state's kind and a transition's effect (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0743.md`) say how its tests prove it:

> The test workflow renames in-progress and review with no alias and adds a state of kind active between them, so that no comparison by name can pass by accident; workitem's own fixture keeps review as an alias.

The fixture is `renamedFlow` in `flai/internal/serve/flow_test.go`, parsed by `renamedWorkflow(t)` for `Options.Workflow`: `doing` and `awaiting` stand for in-progress and review, and `testing`, of kind `active`, sits between them. A pass that still reads a state's name deploys, pushes, queues, or sends back the wrong story under it.

S-0739 gave `flai/internal/workitem` a shared fixture of the same kind, in `flai/internal/workitem/workflow_fixture_test.go` (`wip/archive/agents/S-0739.md`, § Decisions):

> The shared test workflow renames `in-progress` to `doing` and adds `testing`, of kind `active` with a limit, between it and `review`, so the pull's target, two active states, and the epic's walk through an added state are all exercised.

Its helpers are `testFlow(t)`, which parses `testFlowFile` afresh so a test may change it, `withTestFlow(t, r)`, which sets a `Repo` to run it, and `inTestFlow`, which maps a baseline state's name to the fixture's. It renames `in-progress` with no alias, but keeps `review` by name.

S-0747 gave `flai/internal/hostapi` a hook of its own, because `workitem.Open` always gives the baseline (`wip/archive/agents/S-0747.md`, § Decisions):

> Tests need another workflow, and workitem.Open always gives the baseline: T-3228 adds flai/internal/hostapi/flow.go, an openRepo helper with a test hook that sets the opened Repo's Workflow (S-0737's hook), and the reads of T-3229 use it.

Every host API method opens its repository through `openRepo(root)` in `flai/internal/hostapi/flow.go`, which sets the package variable `testFlow` on the `Repo` when a test has set it. A test calls `withWorkflow(t, text)` in `flai/internal/hostapi/writes_test.go`, which parses the text and restores `testFlow` when the test ends. Its fixture, `testingWorkflow`, adds `testing` of kind `active`, renames review to `acceptance-review`, done to `shipped`, and cancelled to `dropped`, and puts `reason-given` on the cancel and send-back transitions as the baseline has them. It keeps `review` as an alias of `acceptance-review`.

So a later E-0034 story that moves a package onto kinds tests it on a board of that shape, not only on the baseline. In `workitem` it uses `withTestFlow`; in `hostapi` it opens the repository through `openRepo` and uses `withWorkflow`, and adds no second hook. Elsewhere it renames the waiting state too, as `renamedFlow` does, and does not lean on `repoWithWorkflow`'s fixture in `flai/internal/workitem/rules_test.go` alone (D-0627), whose `review` alias would let a name comparison pass. The `workitem` and `hostapi` fixtures keep that alias too.

The dashboard has fixtures of the same kind, in `flaiover/src/lib/lanes.fixture.ts`, the lanes and transitions `board.get` answers: `baselineFlow`, `testingFlow` (in-progress renamed `doing`, `testing` added after it, S-0739), `addedStateFlow` (`testing` added after in-progress, which keeps its name, S-0758), and, since S-0838, `renamedFlow`, the baseline with its pullable queue, accepted, and cancelled lanes renamed `queued`, `shipped`, and `dropped`. S-0838 fixed [I-0399](../issues/I-0399-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md), where the epic filters and the move route still compared names with ready, done, and cancelled, and its narrative (`wip/archive/agents/S-0838.md`, § Decisions) records that T-3568 added the renamed-lane fixture first, "since it adds the renamed-lane fixture flow the later tests read". The issue's fix says:

> renamedFlow in lanes.fixture.ts, with ready, done, and cancelled renamed, reproduces the cause in kinds.test.ts and in each view's test.

So a flaiover view or route moved onto a state's kind (`$lib/workflow`, [flaiover-dashboard.md § A state's kind](../system/flaiover-dashboard.md#a-states-kind-s-0758)) is tested on `renamedFlow` when it reads ready, done, or cancelled, and on `testingFlow` or `addedStateFlow` when it reads in-progress or an added state, rather than on a fixture of its own.

Sources: S-0743's, S-0739's, S-0747's, and S-0838's narratives, `wip/archive/agents/S-0743.md`, `wip/archive/agents/S-0739.md`, `wip/archive/agents/S-0747.md`, and `wip/archive/agents/S-0838.md`, § Decisions; I-0399's fix.
