---
id: D-0388
title: TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain failed now and then because endRun records the run's end before blockAsked writes the block
kind: cause
sources: [I-0220, S-0503]
topics: [flai, cli, go, serve, testing, code]
paths: [flai/internal/serve/agents_test.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0388 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain failed now and then because endRun records the run's end before blockAsked writes the block

I-0220, closed 2026-10-28: `TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain` (`flai/internal/serve/agents_test.go`) read the story's blocks once the run read as ended, before `endRun` wrote the block, so it failed now and then.

The fix, in S-0503 (T-2197): "the test waits for the block with waitFor instead of reading the story's blocks once the run reads as ended, since endRun records the end before blockAsked writes the block by design."

The lesson for later serve tests: a run that reads as ended has not yet written everything its end causes; wait for the effect itself.

Sources: I-0220, S-0503.
