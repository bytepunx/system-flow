---
id: T-1143
type: task
nature: improvement
title: wip.overlap compares the claims of the stories in progress, one finding per pair of stories, as the pull hold does
status: backlog
parent: S-0279
owner: alex
created: 2026-10-07T01:20:15Z
updated: 2026-10-07T01:20:40Z
transitions: []
stream: S-0279
tags: [flai]
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go]
after: [T-1141]
---
# T-1143 wip.overlap compares the claims of the stories in progress, one finding per pair of stories, as the pull hold does

## Work

Rewrite `overlap()` in `flai/internal/check/check.go` as T-1141's ADR decides. It waits for T-1141 because that ADR settles the rule. T-1145 waits for this task, because it scopes the findings by the story IDs that this task puts on them.

Today `overlap()` pairs every in-progress item by its own raw touches. It skips only an item and its parent. The recommended rule changes four things:

- Compare stories, not items. Each story in progress is compared by its claim, `Holds.Claim`, built from a `Holds` that knows the items. Today it is `NewHolds(nil, ...)`, which has no tasks. In the claim, a folder touch is narrowed to the touches its tasks name inside it, and an open task's touch outside the story's touches is added (ADR-0096 §3).
- Report each pair of stories once, on the first story's item file at its `touches` line. A story is no longer reported against another story's task.
- Carry the two story IDs on each `wip.overlap` finding, in a field of `Finding`, so that T-1145 can scope it without parsing the message. Leave the field out of `--json` unless the ADR says otherwise.
- Keep the shared-path exception: `Holds.Overlaps` with `WithShared`.

`check.Overlaps` gives the same findings to the designer's inbox (`flai/internal/hostapi/people.go`, S-0158). `flai/internal/itemedit/itemedit.go` compares a write's new `wip.overlap` findings with those before it. Run their tests. When one of them needs a change, widen this task's touches with `flai touches --add` first.

## Done when

- A test in `flai/internal/check/check_test.go` reproduces the S-0223 instance and passes. A story in progress touches a folder, and its task names one file inside the folder. Another story in progress touches a different file in that folder. No `wip.overlap` is reported.
- Another test, from the S-0262 and T-0877 instance, passes. Two stories in progress overlap, and a task of one touches the same path. Exactly one finding names the two stories, and it carries both IDs.
- `TestTouchesOverlap` and `TestTouchesOverlapInsideASharedPath` pass, updated where the rule changes them.
- `flai test flai/internal/check flai/internal/itemedit flai/internal/hostapi` passes.

## Notes
