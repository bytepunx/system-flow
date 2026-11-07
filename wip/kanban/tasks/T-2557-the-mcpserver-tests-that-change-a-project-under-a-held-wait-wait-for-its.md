---
id: T-2557
type: task
nature: remediation
title: The mcpserver tests that change a project under a held wait wait for its heartbeat instead of sleeping 150 ms, and hold up to 30 s
status: backlog
parent: S-0560
owner: alex
created: 2026-11-07T08:20:14Z
updated: 2026-11-07T08:20:14Z
transitions: []
stream: S-0560
tags: [flai, tests]
touches: [flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/orchestrator_threads_test.go]
after: [T-2555]
---
# T-2557 The mcpserver tests that change a project under a held wait wait for its heartbeat instead of sleeping 150 ms, and hold up to 30 s

## Work

Waits for T-2555, because a first heartbeat only means the wait is holding once that task has put `keepAlive` after the snapshot in `waitForWork`.

S-0569 fixed this race in the folder fixture's tests only (T-2522): `callHeld` waits for the wait's first heartbeat and holds up to `heldWait`, 30 s. The single-project fixture keeps the old shape. Its tests sleep 150 ms, change the project, and give the wait `timeout_seconds: 3` under `MaxWait` 3 s, and `answered` gives up after 5 s. `TestWaitForWorkPassesOverAHeldStory`, in I-0246's instances, is one of them. On a host at load 87 those bounds are reached, and the wait answers `timed_out`.

Give the fixture in `flai/internal/mcpserver/server_test.go` what the folder fixture has: a short `Heartbeat`, a progress handler, and a held call that returns once the wait holds. Keep one `heldWait`. Move these 19 tests to it, in place of their sleep:

| File | Tests |
|------|-------|
| `work_test.go` | `TestWaitForWorkWaitsForAStoryToBeReady`, `TestWaitForWorkWaitsForRoomUnderTheLimit`, `TestWaitForWorkWaitsWhileReviewIsFull`, `TestWaitForWorkPassesOverAHeldStory`, `TestWaitForWorkWakesForAThreadWrittenWhileWaiting`, `TestARecommendationIsNoAnswerToTheAgentThatAskedUntilConfirmed`, `TestTheOperatorsThreadsAwaitTheOperatorNotTheAgent` |
| `server_test.go` | `TestDesignerReplyReachesAWaitingAgent`, `TestWaitForEventsReportsAnEventWhileHeld`, `TestAnEditBySomeoneElseReachesTheAgent`, `TestAnOverlapAtAcceptanceReachesTheAgent`, `TestAnEmptyEventLogDoesNotWakeWaitForEventsBeforeItsEvent`, `TestALoggedEventWakesWaitForEvents`, `TestAnAgentEventIsReportedOnceWhileWaitForEventsIsHeld` |
| `messages_test.go` | `TestAMessageWakesAWaitingAgent`, `TestASettledReplyWakesAWaitingAgentOwingNone` |
| `orchestrator_threads_test.go` | `TestTheOrchestratorRecommendsAndTheOperatorConfirms`, `TestTheOrchestratorAnswersCitingItsSource`, `TestTheOrchestratorEscalatesWhatItCannotSource` |

Raise each test's hold and its own select to the bound of `heldWait`: a longer bound costs only a run that fails. Leave a wait that a test expects to time out on its short `timeout_seconds`, and raise only the bound of the select that reads it. Keep each test's assertions, and check that one measuring how soon the wait answers still means what it did. Run the package once afterwards and record its time in the narrative, against the 136 s S-0548 measured under `-race`.

## Done when

- [ ] None of the 19 tests sleeps before changing the project, and each holds up to `heldWait`
- [ ] Each test asserts no less than before
- [ ] `GOMAXPROCS=2 go test -race -count=5 ./internal/mcpserver`, run from `flai/`, passes beside a CPU load, recorded in the narrative with the package's time

## Notes
