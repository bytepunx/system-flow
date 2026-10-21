---
id: I-0145
title: flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words
class: impression
status: closed
count: 1
cost: 10m
first_reported: 2026-10-13T18:47:09Z
last_reported: 2026-10-13T18:47:09Z
updated: 2026-10-21T08:08:39Z
---

# I-0145 flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words

## Description
flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words

## Instances

### 2026-10-13T18:47:09Z
Story: S-0378.
S-0378 made `flai accept --json` refuse markdown the lint fix left with exit 4 and each finding as data (path, line, rule, checkout), which `accept.run` passes to the review page. The board's move to done and the item page go through `item.move` (`flai move`), which answers the same refusal as text with exit 1. To open each named file on the right checkout, `flaiover/src/lib/filelinks.ts` `linkSegments` drops the story when the text says "the main checkout's uncommitted files", matching flai's wording in `flai/cmd/accept_markdown.go`. A change of that wording would silently send main-checkout files to the story's branch. Remedy: let `flai move <story> done --json` print the same refusal (`acceptRefused`) and give `item.move` `exits: {4: Refused}`, then have the board and item page read the findings as data, as Review does.

## Remediation

Story S-0465 remediates this issue, created from it at 2026-10-21T07:24:59Z.
Closed 2026-10-21T08:08:39Z: S-0465: flai move <story> done --json now refuses as flai accept --json does (exit 4, the refused object with each finding's checkout), item.move answers it as Refused with that data, and the board and the item page list the findings through AcceptFindings, so linkSegments no longer reads flai's wording to pick a file's checkout.
