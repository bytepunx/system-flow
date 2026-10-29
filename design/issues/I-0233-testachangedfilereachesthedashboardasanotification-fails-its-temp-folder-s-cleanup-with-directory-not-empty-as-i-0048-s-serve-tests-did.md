---
id: I-0233
title: "TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with \"directory not empty\", as I-0048's serve tests did"
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-29T16:37:28Z
last_reported: 2026-10-29T16:37:28Z
updated: 2026-10-29T16:37:28Z
---

# I-0233 TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did

## Description
TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did

## Instances

### 2026-10-29T16:37:28Z
Story: S-0537.
2026-10-29: S-0537's close-out stopped at the go-test tier after 8m5s on flai/internal/serve TestAChangedFileReachesTheDashboardAsANotification (serve_test.go:337). The failure: "TempDir RemoveAll cleanup: unlinkat /tmp/TestAChangedFileReachesTheDashboardAsANotification.../001/serve: directory not empty". Something still writes to the serve folder after the test returns, as I-0048 found for the stub agents (closed by S-0122's newAgentLab wait). S-0537 changed flai/internal/serve/orchestrate.go and orchestrate_reasons.go, the prime forms of the orchestrator's reasons, and not the dashboard notification path the test covers. The close-out was run again.

## Remediation
