---
id: I-0342
title: "TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with \"the setup src/routes/layout.css does not compile: 404, 404 page not found\" while the main checkout's flaiover node_modules is intact"
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-10T10:40:56Z
last_reported: 2026-11-10T10:40:56Z
updated: 2026-11-10T10:54:41Z
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
