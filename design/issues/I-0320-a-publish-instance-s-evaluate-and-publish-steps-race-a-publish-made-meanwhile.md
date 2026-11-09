---
id: I-0320
title: A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-09T16:58:32Z
last_reported: 2026-11-09T16:58:32Z
updated: 2026-11-09T21:22:33Z
---

# I-0320 A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make

## Description
A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make

## Instances

### 2026-11-09T16:58:32Z
Story: S-0641.
S-0641, the publish formula's instance made by the orchestrator for S-0610 and S-0615: T-2733's `flai release --evaluate --require-met` exited 0 at about 16:57:50Z with 2 stories pending, then commit 7029b028a, "chore: publish flai 1.86.5, template 1.1.28", landed on main under the operator's name at 16:57:57Z and pushed flai/v1.86.5, and T-2734's `flai release --pending --json` ran after it, exiting 0 with `"reason": "nothing pending"`. The release the instance was made for happened, but not by its step, and nothing in the instance holds a second publish back while it runs, nor tells its step's exit 0 with nothing pending apart from a publish it made. The step's own outcome cannot say which release covered the pending stories.

## Remediation

Story S-0643 remediates this issue, created from it at 2026-11-09T21:22:33Z.
