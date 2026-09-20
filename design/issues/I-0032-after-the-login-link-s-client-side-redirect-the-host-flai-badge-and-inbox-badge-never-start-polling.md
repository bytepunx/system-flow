---
id: I-0032
title: After the login link's client-side redirect, the host flai badge and inbox badge never start polling
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-20T21:41:46Z
last_reported: 2026-09-20T21:41:46Z
updated: 2026-09-20T21:41:46Z
---

# I-0032 After the login link's client-side redirect, the host flai badge and inbox badge never start polling

## Description
After the login link's client-side redirect, the host flai badge and inbox badge never start polling

## Instances

### 2026-09-20T21:41:46Z
2026-09-20, found while checking S-0080's project switcher in a browser. +layout.svelte's onMount guards hostFlai.start(), inboxState.start(), and (since S-0080) projectState.refresh() with signedIn = page.url.pathname !== '/login', evaluated once at mount. The layout mounts once for the whole SPA; login/+page.svelte's submit() calls goto() (a client-side transition, not a hard reload) to leave /login, so onMount never re-runs and signedIn stays false from the moment of the first visit. Following the login link directly (`http://host/login#token=...`) hits this every time; a bookmark or a hard reload of any other page does not, because the layout then mounts with signedIn already true. Reproduced on story/S-0080's branch and confirmed present on main (same goto() call, unrelated to this story). Pre-existing since S-0035/S-0042 (whichever first added a background poller started from +layout.svelte's onMount), not something S-0080 introduced; recorded because the new project switcher hit the same gap. Remedy to consider: re-evaluate signedIn reactively (e.g. from page.url via $effect) instead of once in onMount, or have the login page's redirect be a hard navigation.

## Remediation
