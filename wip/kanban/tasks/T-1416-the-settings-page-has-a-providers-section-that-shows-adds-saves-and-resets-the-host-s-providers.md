---
id: T-1416
type: task
nature: feature
title: The Settings page has a Providers section that shows, adds, saves, and resets the host's providers
status: backlog
parent: S-0359
owner: alex
created: 2026-10-08T08:53:16Z
updated: 2026-10-08T08:54:29Z
transitions: []
stream: S-0359
tags: [dashboard]
touches: [flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/server/agent.ts]
after: [T-1415]
---
# T-1416 The Settings page has a Providers section that shows, adds, saves, and resets the host's providers

## Work

- `settings.ts`: the `Provider` type and `providers` under `agent`, `provider` among the settable names and in `HOSTWIDE`.
- `lib/server/agent.ts`: `settings.provider` among the methods the dashboard may call.
- `SettingsPanel.svelte`, after the harnesses: one box per provider with `api` (a select of the three), `base_url`, `key_env`, and `models` one `alias=name` a line, a line saying `key set` or `<VAR> is not set in flai serve's environment`, Save and Reset; and a form to add one by name. Read-only with the enabling command while `settings` is off, as the other sections are.
- The default agent's `AgentFields` on this page binds T-1417's `provider` and sends it through `settings.default_agent`.
- `SettingsPanel.svelte.test.ts`: the section read-only and writable, a save's payload, a reset, the unset key line, and a default agent with a provider.

It waits for T-1415, whose read and write it calls. It runs together with T-1417, which touches other files; the binding of the default agent's provider works once both are in.

## Done when

- `flai test flaiover/src/lib/settings.ts flaiover/src/lib/server/agent.ts flaiover/src/lib/components/SettingsPanel.svelte` passes.

## Notes

Layer 2 of S-0359.
