---
id: T-1044
type: task
nature: improvement
title: The Updates page lists the flai and dashboard releases and deploys a chosen one, as Upgrade does for the newest
status: backlog
parent: S-0298
owner: alex
created: 2026-10-06T21:45:18Z
updated: 2026-10-06T21:56:38Z
transitions: []
stream: S-0298
tags: [dashboard]
touches: [flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts]
after: [T-1042]
---
# T-1044 The Updates page lists the flai and dashboard releases and deploys a chosen one, as Upgrade does for the newest

## Work

On the Updates page (`/host`), the Dashboard area (`HostPanel.svelte`) and the flai host area (`HostProcesses.svelte`) each get a **Versions** control beside Check for updates: it reads `versions` from T-1042's routes and lists the published releases newest first, marking the running one, the newest, and, for flai, any below a project's `flai.minimum`. Choosing one and confirming deploys it through `upgrade` with that `version` or `tag`, behind the same host action, the same confirmation, the same "Reconnecting…" handling of the dropped connection, and the same report of what came back that Upgrade has today. As the operator answered on TH-0202, a dashboard version chosen here applies once, with no pinning: the confirmation and the result say the configured tag is used again at the container's next start, and the page offers no control to pin. Waits for T-1042 for the routes it calls.

## Done when

- Each area lists the releases and deploys a chosen one, tested in `HostPanel.svelte.test.ts` and `HostProcesses.svelte.test.ts`, including the dropped-connection path and the action being off.
- Choosing a dashboard version says it lasts until the next start, tested.
- `npm run test`, `npm run check`, and `npm run lint` pass in `flaiover`.

## Notes
