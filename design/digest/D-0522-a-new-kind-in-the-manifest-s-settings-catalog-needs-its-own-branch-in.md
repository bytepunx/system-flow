---
id: D-0522
title: A new kind in the manifest's settings catalog needs its own branch in settingArg in flai/internal/hostapi/settings.go, or settings.manifest refuses the dashboard's value for it
kind: decision
sources: [S-0595, T-2598]
topics: [manifest, dashboard, flaiover, flai, go, code]
paths: [flai/internal/hostapi/settings.go, flai/internal/manifest/settings.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0522 A new kind in the manifest's settings catalog needs its own branch in settingArg in flai/internal/hostapi/settings.go, or settings.manifest refuses the dashboard's value for it

S-0595 added the settings catalog kind `components`, beside `text`, `command`, and `duration` ([ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md) § Consequences). Its narrative, `wip/archive/agents/S-0595.md`, records under Decisions:

> `hostapi/settings.go` gained a `components` branch in `settingArg` (T-2598), since `settings.manifest` would otherwise have refused the JSON list.

The dashboard's settings panels write through the host API's `settings.manifest`, which turns each value into the `key=value` argument `flai manifest set` reads. `settingArg` does that by the setting's kind, one branch per kind: true or false, digits, a line of text, an agent as compact JSON, a command as a JSON argument list. A kind with no branch there is refused before `flai manifest set` runs. So a story that adds a kind to `manifest.Settings` in `flai/internal/manifest/settings.go` also adds its branch to `settingArg`, with a test in `flai/internal/hostapi/settings_test.go`, in the same task.
