---
id: I-0180
title: The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-10-21T09:26:50Z
last_reported: 2026-10-21T09:26:50Z
updated: 2026-11-07T04:01:14Z
---

# I-0180 The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport

## Description
The top bar's project switcher has no maximum width, so with two projects connected the page is wider than a 390px viewport

## Instances

### 2026-10-21T09:26:50Z
Story: S-0462.
S-0462's Playwright test first answered the board from a fake flai under a second project key. With two projects known to the server, layout.e2e.ts › at 390px › / failed: the page was 42px wider than the viewport, because ProjectSwitcher.svelte's select grows to fit its longest option ("E2E board (unsigned allowed: flai and dashboard)"). The suite misses it only because no run connects two projects. S-0462's test now answers /api/board by route interception instead; the defect stands.

## Remediation

Story S-0557 remediates this issue, created from it at 2026-11-02T00:04:25Z.
Closed 2026-11-07T04:01:14Z: S-0557 fixed it. ProjectSwitcher.svelte's select is capped at 12rem, 20rem from sm, and its wrapper shrinks and wraps, so a long label can no longer widen the top bar; the current project's whole label leads the select's title (T-2507). flaiover/e2e/project-switcher.e2e.ts reproduces I-0180 with two long-named projects at 390px and 1440px: it fails on the old switcher, the select ending at 548px, and passes with the cap (T-2509).
