---
id: T-1058
type: task
nature: improvement
title: The host channel method stream.state writes a narrative's Current state and Next steps
status: backlog
parent: S-0271
owner: alex
created: 2026-10-06T22:49:56Z
updated: 2026-10-06T22:49:56Z
transitions: []
stream: S-0271
tags: [hostapi]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1055]
---
# T-1058 The host channel method stream.state writes a narrative's Current state and Next steps

## Work

Criterion 2, on the host channel.

- In `flai/internal/hostapi/writes.go`, add `stream.state` beside `stream.log` and `stream.answer`. It takes `{id, current?, next?}`.
- Build `flai stream state <id> --current=<text> --next=<text> --json`. Every value is a `--flag=value`, as the other writes do.
- Refuse a call that gives neither text as invalid params before anything runs.
- Map flai's refusal exit to Refused (-32010), as `stream.log` does.
- Who acted is the manifest's owner.

## Done when

- [ ] `writes_test.go` lists a valid call and the refusals for `stream.state`. The test that fails a method with neither passes.
- [ ] `flai hostapi stream.state '{...}'` writes the two sections in a fixture.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner for S-0271. It waits for the command, which it runs. The dashboard does not call `stream.state`: criterion 2 asks only that the dashboard tick, which it already does through `item.criteria`.
