---
id: D-0388
title: TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain failed now and then because endRun records the run's end before blockAsked writes the block
kind: cause
sources: [I-0220, S-0503, I-0248, S-0561]
topics: [flai, cli, go, serve, testing, code]
paths: [flai/internal/serve/agents_test.go, flai/internal/serve/stop_test.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0388 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain failed now and then because endRun records the run's end before blockAsked writes the block

I-0220, closed 2026-10-28: `TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain` (`flai/internal/serve/agents_test.go`) read the story's blocks once the run read as ended, before `endRun` wrote the block, so it failed now and then.

The fix, in S-0503 (T-2197): "the test waits for the block with waitFor instead of reading the story's blocks once the run reads as ended, since endRun records the end before blockAsked writes the block by design."

I-0248, closed 2026-11-07, is the same order seen from the event log. `TestTheOperatorsStopIsOneStoppedAgentEvent`'s waiting-for-an-answer case (`flai/internal/serve/stop_test.go`) stopped the agent as soon as the run's outcome read asked. Under a loaded close-out the stopped event was appended before the ended event. The fix, in S-0561 (T-2550, 375012db1): "the waiting-for-an-answer case ... now waits for the ended event before it stops the agent. The launcher's endRun records the run's outcome, then appends its ended event." S-0561 reproduced it with a temporary 300ms sleep between `l.ended(run)` and `appendAgentEnd`, and committed no seam for it. No product code changed: "the event log is news, not state (ADR-0211 §4)", so `endRun` keeps its order.

The lesson for later serve tests: a run that reads as ended has not yet written everything its end causes, its block or its ended event; wait for the effect itself.

Sources: I-0220, S-0503, I-0248, S-0561.
