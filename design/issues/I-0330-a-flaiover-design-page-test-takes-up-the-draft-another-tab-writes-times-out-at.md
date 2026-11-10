---
id: I-0330
title: "A flaiover design-page test, \"takes up the draft another tab writes\", times out at 5,000 ms under load and stops the flaiover tier"
class: defect
status: closed
count: 1
cost: 16m
first_reported: 2026-11-10T05:37:30Z
last_reported: 2026-11-10T05:37:30Z
updated: 2026-11-10T07:31:27Z
---

# I-0330 A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier

## Description
A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier

## Instances

### 2026-11-10T05:37:30Z
Story: S-0640.
S-0640's close-out, job S-0640-20261110T052050Z at 9ed01712a, stopped at the flaiover tier after 15m45s. It failed on one test only: `flaiover/src/routes/design/design.svelte.test.ts` › editing the design system on the design page (S-0383) › "takes up the draft another tab writes, and drops its edits when another tab drops it". The test timed out at 5,000 ms, and 2,429 of 2,430 tests passed. The one-minute load average was about 108 to 111 on 24 cores. S-0640 changes no file under the Design page. The test settles only three times, and `open()` mounts the page once, so it does not have I-0319's cause of a settle per message. The time goes in mounting and rendering the design page at that load. I-0269 is the same kind of defect, in another file. The close-out was asked for again.

## Remediation

Story S-0655 remediates this issue, created from it at 2026-11-10T05:54:35Z.
Closed 2026-11-10T07:31:27Z: S-0655 (T-2838, d171add4b) gives the S-0383 block of flaiover/src/routes/design/design.svelte.test.ts a 15 s timeout. Timed phase by phase, the test's time is the page's own mount and its re-renders into and out of Edit mode, 71 to 622 ms per settle, in each settle's first timer turn. The other 19 turns take about 1 ms each, so the test cannot cut that time. The test took 1.0 to 3.2 s alone at a load of 28 to 63. At a load of about 110 it passed 5 s, and 15 s leaves three times the default. Every test in the block mounts the page the same way, so the timeout covers the block.
