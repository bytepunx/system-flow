---
id: T-2085
type: task
nature: feature
title: The Host page edits the accepted hosts through settings.dashboard_hosts, behind the settings host action, and says a restart applies them
status: done
parent: S-0489
owner: alex
created: 2026-10-23T18:08:42Z
updated: 2026-10-23T19:06:14Z
transitions:
  - to: ready
    at: 2026-10-23T18:55:17Z
    by: agent-S-0489
  - to: in-progress
    at: 2026-10-23T18:55:17Z
    by: agent-S-0489
  - to: done
    at: 2026-10-23T19:06:14Z
    by: agent-S-0489
stream: S-0489
tags: [flai, flaiover, dashboard]
touches: [flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/settings.ts, flaiover/src/lib/components/HostAddresses.svelte, flaiover/src/lib/components/HostAddresses.svelte.test.ts, flaiover/src/routes/host/+page.svelte, flaiover/src/routes/api/settings/+server.ts, flaiover/src/routes/api/settings/settings.test.ts, design/components/flaiover/src/lib/components/HostAddresses.examples.json, design/system/flaiover-dashboard.md, design/system/flai-cli.md, flaiover/src/lib/server/agent.ts]
after: [T-2083, T-2084]
usage:
  source: log
  seconds: 657
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 143
      output: 46355
      cache_read: 9658383
      cache_write: 269744
      cost: 4.4692
---
# T-2085 The Host page edits the accepted hosts through settings.dashboard_hosts, behind the settings host action, and says a restart applies them

## Work

Host API: add `settings.dashboard_hosts` to `settingsSpecs` in `flai/internal/hostapi/settings.go`, a host-wide write (`hostwide`, so it needs `flai serve enable settings --all-projects` like the dashboard token), taking `{hosts: string}` and running `flai config set dashboard.allowed_hosts <hosts>`; an entry with a space, a scheme, a path, or a character that is not a host name's is refused with `bad`, naming it, and an empty string clears the list. Cover it in `flai/internal/hostapi/writes_test.go` beside the other settings writes.

Dashboard: add `dashboard_hosts` to `SETTINGS_KINDS` in `flaiover/src/lib/settings.ts` and `allowed_hosts` to the `host` block of `SettingsView`; `POST /api/settings` in `flaiover/src/routes/api/settings/+server.ts` then passes it through as it does the others, with a test in `settings.test.ts`. Add `flaiover/src/lib/components/HostAddresses.svelte`, rendered on `/host` (`flaiover/src/routes/host/+page.svelte`) between the Dashboard and flai host areas, with the heading "Addresses": it reads `GET /api/settings` for the configured list and `GET /api/dashboard` (`dashboard.status`) for the list the container runs with, shows a text field with the configured list, one host per entry separated by commas, and a Save button; while the settings action is off host-wide it is read-only and names the command that allows it, as the settings pages do. A save posts `{kind: 'dashboard_hosts', hosts}`, shows flai's refusal under the field when it refuses, and on success says that the running container keeps its old list until it is restarted, pointing at Restart in the Dashboard area above, which the `dashboard` host action gives. When the configured and running lists differ, say so beside the field. Add its examples file under `design/components/` and tests in `HostAddresses.svelte.test.ts`: read-only, a save, a refusal, and the restart note.

Record the write and the page's new area in `design/system/flaiover-dashboard.md` § Writes and where the Host page is described.

Waits for T-2083, which gives the server the variable the saved list reaches it by, and T-2084, which gives flai the key this write sets, `settings.get`'s `allowed_hosts`, and `dashboard.status`'s running list.

## Done when

- With `settings --all-projects` on, saving `flow.example.com` on the Host page runs `flai config set dashboard.allowed_hosts flow.example.com` on the host, the page shows the saved list, and it says a restart applies it; a Restart from the Dashboard area then starts the container with `FLAIOVER_ALLOWED_HOSTS=flow.example.com`.
- With the action off, the area shows the list and the enabling command, and offers no Save.
- An entry that is not a host name is refused with flai's reason under the field, and nothing is written.
- The host API test, the component test, the settings route test, and `flai check` with the new examples file pass.

## Notes

Planned by the planner on 2026-10-23. The setting applies on a restart rather than at once, because the server reads the list from its environment (T-2083); applying it live would need the hub to carry host settings over the channel, which no page needs yet.
