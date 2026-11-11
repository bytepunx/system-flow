---
id: I-0342
title: "TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with \"the setup src/routes/layout.css does not compile: 404, 404 page not found\" while the main checkout's flaiover node_modules is intact"
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-10T10:40:56Z
last_reported: 2026-11-10T10:40:56Z
updated: 2026-11-11T19:01:24Z
---

# I-0342 TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with "the setup src/routes/layout.css does not compile: 404, 404 page not found" while the main checkout's flaiover node_modules is intact

## Description
TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with "the setup src/routes/layout.css does not compile: 404, 404 page not found" while the main checkout's flaiover node_modules is intact

## Instances

### 2026-11-10T10:40:56Z
Story: S-0673.
S-0673's close-out job S-0673-20261110T103135Z at e4b487a stopped after 8m39s at the integration tier, at flai/internal/components/flaiover_test.go:89: "the setup src/routes/layout.css does not compile: 404" with the body "404 page not found". That body is Go's default 404, not vite's. S-0673 changes nothing under flai/internal/components or flaiover. The main checkout's flaiover/node_modules was last written on 2026-10-21, so this is not I-0179's reinstall. Several other stories' close-outs ran at the same time. A guess, not checked: the harness URL reached another process's server on the port.

## Remediation

Story S-0714 remediates this issue, created from it at 2026-11-10T10:54:41Z.
Closed 2026-11-11T19:01:24Z: S-0714: the components dev server now takes its port at random from 10000 to 65535 outside the host's ephemeral range, so no port 0 listener or outgoing connection of another close-out can be given it once vite stops listening (T-3114); the harness's plugin sets X-Flai-Harness to a per-start token on every answer, and StartDevServer takes a 200 only with it, with a test that puts a Go server on the picked port (T-3115); and TestEveryListedFlaioverComponentCompilesInTheHarness fails an answer without the token as another server's on vite's port, not as a compile failure (T-3116).
