---
id: I-0330
title: "A flaiover design-page test, \"takes up the draft another tab writes\", times out at 5,000 ms under load and stops the flaiover tier"
class: defect
status: open
count: 1
cost: 16m
first_reported: 2026-11-10T05:37:30Z
last_reported: 2026-11-10T05:37:30Z
updated: 2026-11-10T05:37:30Z
---

# I-0330 A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier

## Description
A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier

## Instances

### 2026-11-10T05:37:30Z
Story: S-0640.
S-0640's close-out, job S-0640-20261110T052050Z at 9ed01712a, stopped at the flaiover tier after 15m45s. It failed on one test only: `flaiover/src/routes/design/design.svelte.test.ts` › editing the design system on the design page (S-0383) › "takes up the draft another tab writes, and drops its edits when another tab drops it". The test timed out at 5,000 ms, and 2,429 of 2,430 tests passed. The one-minute load average was about 108 to 111 on 24 cores. S-0640 changes no file under the Design page. The test settles only three times, and `open()` mounts the page once, so it does not have I-0319's cause of a settle per message. The time goes in mounting and rendering the design page at that load. I-0269 is the same kind of defect, in another file. The close-out was asked for again.

## Remediation
