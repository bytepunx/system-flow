---
id: D-0572
title: "TestEveryListedFlaioverComponentCompilesInTheHarness can stop a close-out's integration tier with Go's 404 for flaiover's layout.css while other close-outs run, with node_modules intact; I-0342 records it and the close-out runs it again"
kind: cause
sources: [I-0342, S-0673, S-0714, T-3114, T-3115, T-3116]
topics: [testing, go, flaiover, code]
paths: [flai/internal/components/flaiover_test.go, flai/internal/components, flai/internal/components/devserver.go, flai/internal/components/harness/plugin.mjs, flai/cmd/components_test.go, flai/internal/serve/components_test.go, design/system/component-library.md]
date: 2026-11-10
status: absorbed
absorbed_by: "design/system/component-library.md#the-dev-server"
asked_again: []
cited: [S-0714]
---

# D-0572 TestEveryListedFlaioverComponentCompilesInTheHarness can stop a close-out's integration tier with Go's 404 for flaiover's layout.css while other close-outs run, with node_modules intact; I-0342 records it and the close-out runs it again

S-0673's first close-out, at e4b487a, stopped after 8m39s at the integration tier in `flai/internal/components/flaiover_test.go`. I-0342 (`design/issues/I-0342-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md`) records it:

```text
"the setup src/routes/layout.css does not compile: 404" with the body "404 page not found". That body is Go's default 404, not vite's. S-0673 changes nothing under flai/internal/components or flaiover. The main checkout's flaiover/node_modules was last written on 2026-10-21, so this is not I-0179's reinstall. Several other stories' close-outs ran at the same time. A guess, not checked: the harness URL reached another process's server on the port.
```

It is not the missing-dependency cause D-0350 records, since `node_modules` was intact. S-0673 asked for the close-out again, and it passed every step at 4d3800e (`wip/archive/agents/S-0673.md`, Decisions).

S-0714 fixed I-0342 on 2026-11-11: "the components dev server now takes its port at random from 10000 to 65535 outside the host's ephemeral range, so no port 0 listener or outgoing connection of another close-out can be given it once vite stops listening (T-3114); the harness's plugin sets X-Flai-Harness to a per-start token on every answer, and StartDevServer takes a 200 only with it, with a test that puts a Go server on the picked port (T-3115); and TestEveryListedFlaioverComponentCompilesInTheHarness fails an answer without the token as another server's on vite's port, not as a compile failure (T-3116)." This host's narrow ephemeral range is D-0596.

Two of its decisions bear on later changes (wip/archive/agents/S-0714.md, Decisions). "StartDevServer refuses options with no Token, rather than treating an empty token as no check: an optional check would leave the I-0342 cause open for a future caller." And the fake vites in `flai/cmd/components_test.go` and `flai/internal/serve/components_test.go` "now read the token from the config --config names and send the header, as the plugin does. Without that, 13 tests there failed on the new check." A new fake vite in a test does the same.

`design/system/component-library.md` § "The dev server" now states the port pick, the token, and the test's verdict.

Sources: I-0342, S-0673, S-0714 (T-3114, T-3115, T-3116).
