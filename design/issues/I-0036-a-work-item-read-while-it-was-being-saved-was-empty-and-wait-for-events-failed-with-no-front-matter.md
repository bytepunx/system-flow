---
id: I-0036
title: A work item read while it was being saved was empty, and wait_for_events failed with no front matter
class: defect
status: closed
count: 1
cost: 20m
first_reported: 2026-09-23T05:14:33Z
last_reported: 2026-09-23T05:14:33Z
updated: 2026-09-23T05:14:33Z
---

# I-0036 A work item read while it was being saved was empty, and wait_for_events failed with no front matter

## Description
A work item read while it was being saved was empty, and wait_for_events failed with no front matter

## Instances

### 2026-09-23T05:14:33Z
CI (flai workflow, publish of flai 1.12.1): TestWaitForEventsReportsAnEventWhileHeld failed, S-0001-story.md: no front matter. Repo.Save used os.WriteFile, which truncates before writing; the held wait read the file in between. The same could fail an agent's inbox or wait_for_events whenever the dashboard or the CLI saved an item. A race test counted 4,243 torn reads in 2,000 saves with os.WriteFile.

## Remediation
Closed 2026-09-23T05:14:33Z: fixed by S-0100: items, the board, narratives, and threads are replaced in one step (atomicfile.WriteFile)
