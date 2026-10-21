---
id: T-1986
type: task
nature: remediation
title: Pin both runs of TestASettledReplyStartsTheAgentWaitingOnItOnce to one second so the first look never finds the message new to the other story
status: cancelled
parent: S-0473
owner: alex
created: 2026-10-21T12:06:45Z
updated: 2026-10-21T12:11:05Z
transitions:
  - to: cancelled
    at: 2026-10-21T12:11:05Z
    by: agent-S-0473
stream: S-0473
tags: [flai, serve, test]
touches: [flai/internal/serve/agents_test.go]
---
# T-1986 Pin both runs of TestASettledReplyStartsTheAgentWaitingOnItOnce to one second so the first look never finds the message new to the other story

## Work

The planner's reading of the cause, to confirm before the fix: the test calls `lab.putEnded` twice, for `other` and then for `id`, and each stamps its run's `Started` and `Ended` with the wall-clock second (`agents_test.go`, `putEnded`). The message from `id` to `other` is sent at `id`'s `Started`. `messaged` in `flai/internal/serve/agents.go` looks for messages strictly after `max(Ended-1s, Started)` of `other`'s run, through `messages.ToSince`, which compares with `After`. When the two `putEnded` calls fall in the same second, the message is not new to `other`, and the first look starts nothing. When they straddle a second boundary, the message is one second after `other`'s run, so it is new and owes a reply, and the first look starts `other`'s agent (S-0001 on MS-0001), which fails at `started before the reply` (`agents_test.go:1347`). A loaded integration tier makes the straddle likelier. That is why it fails on branches that do not touch serve.

1. Reproduce it on purpose: shift `other`'s run one second earlier than `id`'s, with `lab.endedAgo(other, time.Second)` after the second `putEnded`. Run the test and see it fail at `started before the reply`. Record the result in the narrative, then take the shift out.
2. Fix the test, not `messaged`: the production rule is right, since a message that comes after a run ended is new to it. Make the scenario the test means explicit: `other`'s run was already over the message when it ended. Pin `other`'s run's `Started` and `Ended` to the message's second (`started`), or send the message at a time no later than `other`'s `Ended`. Keep every assertion the test makes. Say in a comment on the test why the times are pinned, naming I-0158.
3. If step 1 does not reproduce the failure, the cause is elsewhere. Stop, record what you found in the narrative, and plan again before you change `agents.go`. Widen this task's touches if you do.

No `after`: this is the story's first layer.

## Done when

- The narrative records the deliberate one-second shift failing the test at `started before the reply`, or why the cause turned out to be elsewhere.
- The test pins `other`'s run to the message's second, and a comment names I-0158 and says why.
- `go test ./internal/serve -run '^TestASettledReplyStartsTheAgentWaitingOnItOnce$' -count=200`, run from `flai/`, passes, and so does `flai test` on `flai/internal/serve/agents_test.go`.

## Notes

Drafted by the planner from I-0158's five instances and a reading of `agents.go` (`messaged`) and `messages.go` (`ToSince`).
- 2026-10-21T12:11:05Z: moved to cancelled: duplicate of T-1992, which did the fix and closed I-0158
