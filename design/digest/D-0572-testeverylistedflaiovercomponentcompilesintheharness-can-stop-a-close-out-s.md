---
id: D-0572
title: "TestEveryListedFlaioverComponentCompilesInTheHarness can stop a close-out's integration tier with Go's 404 for flaiover's layout.css while other close-outs run, with node_modules intact; I-0342 records it and the close-out runs it again"
kind: cause
sources: [I-0342, S-0673]
topics: [testing, go, flaiover, code]
paths: [flai/internal/components/flaiover_test.go, flai/internal/components]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: [S-0714]
---

# D-0572 TestEveryListedFlaioverComponentCompilesInTheHarness can stop a close-out's integration tier with Go's 404 for flaiover's layout.css while other close-outs run, with node_modules intact; I-0342 records it and the close-out runs it again

S-0673's first close-out, at e4b487a, stopped after 8m39s at the integration tier in `flai/internal/components/flaiover_test.go`. I-0342 (`design/issues/I-0342-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md`) records it:

```text
"the setup src/routes/layout.css does not compile: 404" with the body "404 page not found". That body is Go's default 404, not vite's. S-0673 changes nothing under flai/internal/components or flaiover. The main checkout's flaiover/node_modules was last written on 2026-10-21, so this is not I-0179's reinstall. Several other stories' close-outs ran at the same time. A guess, not checked: the harness URL reached another process's server on the port.
```

It is not the missing-dependency cause D-0350 records, since `node_modules` was intact. S-0673 asked for the close-out again, and it passed every step at 4d3800e (`wip/archive/agents/S-0673.md`, Decisions). So a story that changes nothing under `flai/internal/components` or `flaiover` and meets this failure bumps I-0342 and runs its close-out again, as D-0543 and D-0524 say for their tests. S-0714 remediates I-0342.

Sources: I-0342, S-0673.
