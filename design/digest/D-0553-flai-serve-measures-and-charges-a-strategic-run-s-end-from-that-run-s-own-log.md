---
id: D-0553
title: flai serve measures and charges a strategic run's end from that run's own log, AgentRun.Log, not its kind's newest run log, so a run ending after a newer run of its kind started is charged right (I-0314)
kind: cause
sources: [I-0314, S-0653, T-2849, T-2851]
topics: [testing, flai, go, code]
paths: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/role_test.go, flai/internal/serve/analyze.go, design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md]
date: 2026-11-10
status: absorbed
absorbed_by: "design/system/strategic-agents.md#the-run-and-how-it-ended"
asked_again: []
cited: []
---

# D-0553 flai serve measures and charges a strategic run's end from that run's own log, AgentRun.Log, not its kind's newest run log, so a run ending after a newer run of its kind started is charged right (I-0314)

I-0314 was closed on 2026-11-10 by S-0653. `TestAFormulaRoleStepsEndIsToldToTheInstancesAgent` failed on its timing under load, in the go-test tier of any story that changes a package `flai/internal/serve` imports. It had two causes, one in flai's code and one in the test.

The defect, fixed by T-2849: `logRunEnd` and `logRunEndSaying` measured the newest run log of the run's kind. So a run whose end was logged after a newer run of its kind had started was measured and charged from the wrong log. They now take the ended run's own log, `AgentRun.Log`, in `roleEnded`, `planEnded`, `orchestrateEnded`, `curateEnded`, and `analyzeEnded`. From S-0653's Decisions (`wip/archive/agents/S-0653.md`): "A log that holds no run logs nothing rather than falling back to the newest, which would bring the defect back. `LogRunEnd`, the exported form only tests call, keeps measuring the kind's newest run." `TestARunEndedAfterANewerRunStartedIsMeasuredFromItsOwnLog` reproduces it without load.

The test's race, fixed by T-2851: `roleEnded` appends the `agent` event after it records the run as ended, so a single read of the inbox could miss it. The test now waits for each run's event with `waitFor`, and "keeps every look's events, since `inbox.CatchUp` advances the reader's cursor past what it reports."

So a later test that reads an inbox event flai serve appends after a state change polls for it and keeps each look's events, rather than reading once.

Sources: I-0314; S-0653 (`wip/archive/agents/S-0653.md`, Decisions); T-2849; T-2851.
