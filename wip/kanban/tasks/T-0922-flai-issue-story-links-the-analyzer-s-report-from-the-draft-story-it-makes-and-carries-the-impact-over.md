---
id: T-0922
type: task
nature: feature
title: flai issue story links the analyzer's report from the draft story it makes and carries the impact over
status: backlog
parent: S-0224
owner: alex
created: 2026-10-05T05:44:57Z
updated: 2026-10-05T05:44:57Z
transitions: []
stream: S-0224
tags: [flai]
touches: [flai/internal/issues/stories.go, flai/internal/issues/stories_test.go]
after: [T-0918]
---
# T-0922 flai issue story links the analyzer's report from the draft story it makes and carries the impact over

## Work

Waits for T-0918: it reads the `Report:` lines and the `## Impact` section that task writes, and its tests make their issues with T-0918's `issues.New`.

- `ForStory` finds the reports an issue names, in its instances' `Report:` lines and its `## Remediation` links, and the draft story it makes links each by a relative path from the story's file, in the goal beside the issue's link, and says in `## Notes` that the analyzer found it. An issue that names no report makes the story it makes today.
- The cost of delay inputs already carry over from `## Impact` (S-0203); check that an issue written by `flai issue new` with the impact flags and `--report` gives a story whose `cost_of_delay.inputs` hold every figure, and that a bumped duplicate's latest Impact wins.

## Done when

- A story made from an issue the analyzer filed links the issue and its report, and its `cost_of_delay.inputs` hold the issue's impact with `by: flai`
- Tests in `stories_test.go` cover a new analyzer issue, a bumped duplicate naming two reports, and the inputs carried over; `scripts/flai-test.sh` passes

## Notes
