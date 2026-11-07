---
id: I-0308
title: flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T22:56:50Z
last_reported: 2026-11-07T22:56:50Z
updated: 2026-11-07T22:56:50Z
---

# I-0308 flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier

## Description
flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier

## Instances

### 2026-11-07T22:56:50Z
Story: S-0595.
S-0595 added kluster.components to the settings catalog in flai/internal/manifest/settings.go (T-2598). flai test on the changed paths passed, but docs/users/flai-reference.md lists the catalog's keys, and TestReferenceIsCurrent (flai/cmd/reference_test.go:131) runs only in the integration tier. The close-out stopped at integration after 6m26s, asking for make flai-reference. A tier like S-0592's settings-index, selected by flai/internal/manifest/settings.go, that runs ^TestReferenceIsCurrent under -short would have caught it at the task.

## Remediation
