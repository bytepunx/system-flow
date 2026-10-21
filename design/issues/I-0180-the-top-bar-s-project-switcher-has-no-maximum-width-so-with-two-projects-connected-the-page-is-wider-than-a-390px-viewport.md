---
id: I-0180
title: The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-21T09:26:50Z
last_reported: 2026-10-21T09:26:50Z
updated: 2026-10-21T09:26:50Z
---

# I-0180 The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport

## Description
The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport

## Instances

### 2026-10-21T09:26:50Z
Story: S-0462.
S-0462's Playwright test first answered the board from a fake flai under a second project key. With two projects known to the server, layout.e2e.ts › at 390px › / failed: the page was 42px wider than the viewport, because ProjectSwitcher.svelte's select grows to fit its longest option ("E2E board (unsigned allowed: flai and dashboard)"). The suite misses it only because no run connects two projects. S-0462's test now answers /api/board by route interception instead; the defect stands.

## Remediation
