---
id: S-0122
type: story
nature: improvement
title: The dashboard's settings page serves, removes, and shows the health of every project the host knows
status: done
parent: E-0003
owner: alex
created: 2026-09-26T05:18:07Z
updated: 2026-09-26T07:04:22Z
transitions:
  - to: ready
    at: 2026-09-26T05:24:31Z
    by: alex
  - to: in-progress
    at: 2026-09-26T05:59:36Z
    by: agent-S-0119
  - to: review
    at: 2026-09-26T06:59:48Z
    by: agent-S-0122
  - to: done
    at: 2026-09-26T07:04:22Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover, flai/internal/hostapi]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0122 The dashboard's settings page serves, removes, and shows the health of every project the host knows

## Goal

The operator manages the projects the dashboard shows from the dashboard: they see each one's health, add one that is already a system-flow project, and remove one they no longer want, without a shell.

## Acceptance criteria
- [x] The settings page lists every served project with its key, name, root, and connected state or last error from `flai serve`
- [x] A system-flow project under an import root that is not served is shown with why, and with a Serve action that asks flai serve to register it through the host channel and shows flai's answer, a refusal with its reason for every case flai has today. A project flai serve comes to serve joins the switcher without a reload. Reworded on TH-0016 (answer A); a Serve that undoes a Remove is S-0123
- [x] Each served project has a Remove action, confirmed, that unregisters it through the channel, touches none of its files, and says how to add it back
- [x] Serve and Remove are hostapi writes gated by the `settings` host action and journalled like import, and are refused with the reason when the action is not enabled for the project
- [x] The switcher tells a project that is served but disconnected apart from one that is connected, and links to the settings page for why
- [x] Tried end to end in a browser with scratch repositories. Docs in `docs/users` and `design/system/flaiover-dashboard.md`

## Tasks

- T-0445 settings.get reports every project flai serve serves and every one below an import folder it does not, and settings.serve and settings.unserve register and unregister one
- T-0446 The settings page lists the served projects with their health, and Serve and a confirmed Remove change them
- T-0447 The switcher tells a served project that is not connected apart from a connected one, and links to the settings page
- T-0448 Docs: the users guide, the dashboard design, and the channel design say how projects are served and removed from the settings page
- T-0449 Tried end to end in a browser with scratch repositories
- T-0450 flai's serve tests wait for the stub agents they start to end before their temp folders are removed

## Notes

Builds on the host channel methods that S-0121 adds, so pull S-0121 first. Related to [I-0047](../../../design/issues/I-0047-a-project-imported-with-flai-import-on-the-command-line-is-neither-served-nor-offered-so-the-dashboard-never-shows-it.md).

Tried end to end on 2026-09-26 (T-0449), in headless Chrome against a flaiover dev server built from `story/S-0122` on port 5199, and a scratch flai host and serve from the same branch, with their own `HOME`, configuration, and `FLAI_HOST_ADDR`, and no `FLAI_HOST_*` in their environment. The scratch repositories were alpha, delta, and zeta, registered, and beta, gamma (no key), epsilon (key alpha), and later eta, below a scratch import folder. `settings` was on for alpha, beta, zeta, and gamma's folder, and off for delta and epsilon. Every scratch process was stopped by PID afterwards; the operator's host and dashboard were untouched.

- Projects listed alpha, beta, delta, and zeta, each with key, name, folder, and "connected since". beta said it is served because it is below the import folder and had no Remove. Under "Below the import folders, not served" were epsilon (its key alpha is served already, for alpha's folder) and gamma (no key).
- Remove on zeta asked first ("None of its files is touched… Serve it again with flai serve project add <folder>"). Once confirmed, zeta left the list and the switcher without a reload, and the page said how to serve it again. zeta's committed files were unchanged, and the journal had `settings.unserve` for zeta: `flai serve project remove -- <folder>`.
- With `settings` off for delta, its Remove was disabled with "run flai serve enable settings on the host, in <delta's folder>". A POST of `unserve` for delta, through alpha, got 403: `the host action "settings" is not enabled for <delta's folder>`.
- Serve on gamma returned flai's refusal: "not served: its system-flow.yaml has no key". Serve cannot succeed today; see TH-0016.
- With delta's folder moved away, the switcher showed "delta (served, not connected)", titled "the folder is gone", and beside it "1 not connected: why?" linking to `/settings`. Following the link from delta's board switched to alpha and opened Settings.
- eta, committed below the import folder, joined the switcher without a reload after flai serve's next scan.

