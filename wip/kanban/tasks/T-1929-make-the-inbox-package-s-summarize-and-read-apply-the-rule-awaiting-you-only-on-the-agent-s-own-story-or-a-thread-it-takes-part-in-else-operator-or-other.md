---
id: T-1929
type: task
nature: remediation
title: "Make the inbox package's Summarize and Read apply the rule: awaiting you only on the agent's own story or a thread it takes part in, else operator or other"
status: done
parent: S-0453
owner: alex
created: 2026-10-20T22:20:23Z
updated: 2026-10-21T00:54:16Z
transitions:
  - to: ready
    at: 2026-10-21T00:41:48Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T00:41:49Z
    by: system-flow
  - to: done
    at: 2026-10-21T00:54:16Z
    by: system-flow
stream: S-0453
tags: [flai]
touches: [flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go]
after: [T-1927]
usage:
  source: log
  seconds: 747
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 10220
      cache_read: 1515019
      cache_write: 72566
      cost: 0.9158
---
# T-1929 Make the inbox package's Summarize and Read apply the rule: awaiting you only on the agent's own story or a thread it takes part in, else operator or other

## Work

In `flai/internal/inbox/inbox.go`, give `Summarize` what the rule needs: the agent's own story and the operator's names, the manifest's `owner` and the thread's story's `owner`. Set `awaiting` as the ADR of T-1927 records:

- `"you"` when the agent is in `participants` or the thread is on its own story or a task of it (`threads.StoryOf`), and the last entry is not the agent's, ADR-0090's pending-recommendation exception kept.
- `"operator"` for any other open thread that awaits the operator: a recommendation pending, or the last entry by neither operator name.
- `"other"` otherwise.

Make `Read` count only `"you"` in `awaiting_you`, list only those without `all`, and take an option that lists the `"operator"` threads too, for the orchestrator's role. Keep the story-scoped embedded read and `threads_omitted` that S-0450 added counting only threads awaiting the agent.

Add the test that reproduces the measurement of 2026-10-20: many open threads awaiting the operator, none on the reading agent's story. The agent gets no thread listed, `awaiting_you` 0, and `threads_omitted` 0, while a thread on its story whose last entry is the operator's still awaits it. Update `TestThreadsAwaitingTheAgent` to the rule.

It waits for T-1927 so that it builds what the ADR records.

## Done when

- `Summarize` gives `"you"`, `"operator"`, and `"other"` by the rule, and `Read` counts and lists by it.
- The reproduction test and the updated tests pass with `flai test flai/internal/inbox`.

## Notes

The other callers of `Summarize`, in the MCP server, are T-1930's; `story_start` and `task_done` are T-1932's.
