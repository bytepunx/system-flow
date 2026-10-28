---
id: D-0319
title: The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier
kind: cause
sources: [I-0162]
topics: [flai, cli, go, code, all]
paths: [flai/internal/components/devserver.go, flai/internal/components/devserver_test.go, design/system/component-library.md, design/issues/I-0162-the-components-dev-server-test-picks-a-free-port-and-closes-it-before-the-fake-vite-binds-it-so-a-concurrent-test-run-on-the-host-can-take-the-port-and-fail-the-integration-tier.md, design/issues/summary.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md]
date: 2026-10-18
status: absorbed
absorbed_by: "design/system/component-library.md#the-dev-server"
asked_again: []
cited: []
---

# D-0319 The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier

I-0162 was closed on 2026-10-18.

The cause:

```text
The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier
```

The fix:

```text
Fixed by S-0447: StartDevServer starts Vite again on another free port when Vite exits saying its port is already in use, three ports in all, so a port another process takes between freePort closing its listener and Vite binding it no longer fails the start. TestStartDevServerStartsViteAgainOnAnotherPortWhenItsPortWasTaken reproduces it by handing StartDevServer a held port, and fails on main.
```

Source: I-0162.
