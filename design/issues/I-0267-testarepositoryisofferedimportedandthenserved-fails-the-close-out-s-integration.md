---
id: I-0267
title: "TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with \"failed to read frame header: EOF\" under host load, on a branch that does not touch flai/internal/serve"
class: defect
status: open
count: 1
cost: 9m
first_reported: 2026-11-07T02:16:58Z
last_reported: 2026-11-07T02:16:58Z
updated: 2026-11-10T08:55:40Z
---

# I-0267 TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve

## Description
TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve

## Instances

### 2026-11-07T02:16:58Z
Story: S-0580.
S-0580's second close-out, at 13346d1, stopped at integration on TestARepositoryIsOfferedImportedAndThenServed (flai/internal/serve/candidates_test.go:125: failed to get reader: failed to read frame header: EOF). The first close-out, at 27220c8, ran the same code (the two commits differ only in design/issues) and passed this test, stopping on I-0246's test instead. S-0580 does not change flai/internal/serve. The host's load average was 56 over five minutes, with S-0578's agent and close-out running too.

## Remediation

Story S-0703 remediates this issue, created from it at 2026-11-10T08:55:40Z.
