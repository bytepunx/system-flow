---
id: T-3417
type: task
nature: remediation
title: flai verify --last and the host's verify.status say a report is from before the head moved and name the close-out job running at the head
status: backlog
parent: S-0799
owner: alex
created: 2026-11-16T01:46:07Z
updated: 2026-11-16T01:46:07Z
transitions: []
stream: S-0799
tags: [flai, verify]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go]
after: [T-3414]
---
# T-3417 flai verify --last and the host's verify.status say a report is from before the head moved and name the close-out job running at the head

## Work

The review page's Verification row reads `verify.status` (`flai/internal/hostapi/writes_host.go`, `verifyStatus`, line 379). It answers `verify.LastReport` as is, so the page shows the old commit and "Stopped at sync" while a close-out runs at the new head. `flai verify S-nnnn --last` (`flai/cmd/verify.go`) prints the same report.

- In `verifyStatus`, read the story's branch head and call `verify.LastReportAt` from T-3414, so the answer carries `moved` and `running`. When the head cannot be read, answer the report as today and log a warning, as the kluster part does.
- In `flai verify --last`, text and `--json`, do the same. After the verdict, print a line when the report is from before the head moved, naming both commits, and a line naming the close-out job running at the head and its step.
- Test it: in `writes_test.go`, beside the `verify.status` case (line 2691), a report at an old commit and a running job at the head answer `moved` and `running`. In `flai/cmd/verify_test.go`, `--last` prints both lines. Each fails before the change: it reproduces I-0384 on the review page's source.

This task waits for T-3414, whose reader it calls. It shares no file with the other tasks of its layer.

## Done when

- `verify.status` and `flai verify --last --json` carry `moved` and `running` as T-3414 defines them, and `--last` prints them.
- `flai test flai/cmd/verify.go flai/internal/hostapi/writes_host.go` passes.

## Notes
