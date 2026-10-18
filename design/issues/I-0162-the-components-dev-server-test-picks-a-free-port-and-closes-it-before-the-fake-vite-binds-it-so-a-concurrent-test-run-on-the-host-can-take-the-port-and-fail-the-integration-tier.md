---
id: I-0162
title: The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier
class: defect
status: closed
count: 1
cost: 7m
first_reported: 2026-10-17T18:10:30Z
last_reported: 2026-10-17T18:10:30Z
updated: 2026-10-18T02:57:12Z
---

# I-0162 The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier

## Description
The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier

## Instances

### 2026-10-17T18:10:30Z
Story: S-0395.
S-0395's close-out stopped at the integration tier: flai/internal/serve/components_test.go:154 TestComponentsStartsOneDevServerAndAnswersItAgain failed with 'listen tcp 127.0.0.1:42160: bind: address already in use'. freePort in flai/internal/components/devserver.go:304 listens on :0, closes, and hands the number to the dev server, which binds it later; another agent's tests ran on the host at the time. The rerun cost the integration tier's 5 minutes.

## Remediation

Story S-0447 remediates this issue, created from it at 2026-10-18T02:19:04Z.
Closed 2026-10-18T02:57:12Z: Fixed by S-0447: StartDevServer starts Vite again on another free port when Vite exits saying its port is already in use, three ports in all, so a port another process takes between freePort closing its listener and Vite binding it no longer fails the start. TestStartDevServerStartsViteAgainOnAnotherPortWhenItsPortWasTaken reproduces it by handing StartDevServer a held port, and fails on main.
