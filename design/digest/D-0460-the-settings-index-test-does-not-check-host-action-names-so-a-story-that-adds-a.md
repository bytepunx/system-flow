---
id: D-0460
title: The settings-index test does not check host action names, so a story that adds a host action adds it to the host_actions row of docs/operators/settings.md itself
kind: decision
sources: [S-0422, T-1779]
topics: [docs, cli, flai]
paths: [docs/operators/settings.md, docs/operators/index.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0460 The settings-index test does not check host action names, so a story that adds a host action adds it to the host_actions row of docs/operators/settings.md itself

S-0422 added the `formula` host action. Its row in `docs/operators/settings.md` was left to the docs task, T-1779, because no test would fail without it.

From S-0422's narrative, `## Decisions`:

```text
The `formula` host action's row in `docs/operators/settings.md` waits for T-1779: the settings-index test does not cover host actions.
```

The settings-index test fails on a configuration key, a manifest key, or an environment variable with no row (`design/conventions/documentation.md` § Project additions). Host actions all share one key, `host_actions.<name>`, whose row lists each action by name. A story that adds a host action names it in that row, and links its section in `docs/operators/index.md`, by hand; nothing checks it.

Sources: S-0422 (wip/archive/agents/S-0422.md), T-1779.
