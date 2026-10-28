---
id: T-2197
type: task
nature: feature
title: TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain waits for the block endRun writes after the end, so the merge queue's integration tier passes
status: done
parent: S-0503
owner: alex
created: 2026-10-28T12:34:09Z
updated: 2026-10-28T12:35:31Z
transitions:
  - to: ready
    at: 2026-10-28T12:34:14Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-28T12:34:15Z
    by: agent-S-0503
  - to: done
    at: 2026-10-28T12:35:31Z
    by: agent-S-0503
stream: S-0503
tags: []
touches: [flai/internal/serve/agents_test.go, design/issues/I-0220-testastorywhoseagentendedaskingisblockeduntiltheanswerstartsitagain-reads-the-story-s-blocks-once-the-run-reads-as-ended-before-endrun-writes-the-block-so-it-fails-now-and-then.md, design/issues/summary.md]
usage:
  source: log
  seconds: 76
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 735
      cache_read: 203087
      cache_write: 6502
      cost: 0.1074
---

# T-2197 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain waits for the block endRun writes after the end, so the merge queue's integration tier passes

## Work

## Done when

## Notes
