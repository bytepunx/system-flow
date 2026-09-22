---
id: T-0324
type: task
nature: feature
title: A host page shows the dashboard's own version and update state, with restart/upgrade/stop actions
status: backlog
parent: S-0081
owner: alex
created: 2026-09-22T21:12:56Z
updated: 2026-09-22T21:12:56Z
transitions: []
stream: S-0081
tags: []
---
# T-0324 A host page shows the dashboard's own version and update state, with restart/upgrade/stop actions

## Work
New route `/host` (`flaiover/src/routes/host/`) and `/api/dashboard` (`GET` for status, `POST` with `{action: "check"|"restart"|"upgrade"|"stop", tag?}` for the rest, mirroring `/api/publish`'s shape). The header's host flai badge (`HostFlai.svelte`) links to it; add "Host" to the nav in `+layout.svelte`.

The page shows: the running container's image and tag (from `dashboard.status`, polled like `hostFlai` already is), a "Check for updates" button (`dashboard.check`) that reports up to date or what is available without changing anything, and Restart / Upgrade / Stop buttons. Each write goes through the same request-ID/host-action-disabled handling `PublishBanner.svelte` already established (a `request_id`, a `Disabled`-coded refusal shown as "ask the operator to run `flai serve enable dashboard`", not a raw error). While a restart or upgrade runs, the page shows that it is in progress (the write is long-running: real Docker time, not instant); when the `EventSource` this page or the header holds reconnects after a gap, the page refreshes its status and says a restart or upgrade happened, distinguishing "reconnected, version unchanged" from "reconnected, now running v<new>" using the version it had before the write against the version `dashboard.status` reports after reconnecting.

## Done when
New and existing component/unit tests pass (`pnpm run test:unit`, `pnpm run check`, `pnpm run lint`); the disabled-action and refusal paths are covered the same way `PublishBanner.svelte.test.ts` covers them for publish.

## Notes
Depends on T-0323's hostapi methods. No new capability is added to the existing `/api/agent` polling; this page uses its own `dashboard.status` poll, on its own route, so other pages are not slowed by it.
