---
id: I-0359
title: flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T08:58:25Z
last_reported: 2026-11-13T08:58:25Z
updated: 2026-11-13T08:58:25Z
---

# I-0359 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1

## Description
flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1

## Instances

### 2026-11-13T08:58:25Z
Story: S-0764.
S-0764, the publish formula's publish step (T-3275): `flai release --pending --require-pending --json` in the main checkout exited 1 at 2026-11-13T08:57:33Z, because its `git add -A` could not create `.git/index.lock` ("File exists"). The orchestrator and the curator were running at that moment, and seconds later the lock was gone. Nothing was tagged, committed, or pushed. The step's stop lists only exit 2, so the run is blocked on a thread (TH-1142) until the operator answers, although running it again would most likely have passed. I-0100 (closed) recorded the same race in flai accept's commit step.

## Remediation
