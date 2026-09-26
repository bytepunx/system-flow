---
id: S-0122
type: story
nature: improvement
title: The dashboard's settings page serves, removes, and shows the health of every project the host knows
status: in-progress
parent: E-0003
owner: alex
created: 2026-09-26T05:18:07Z
updated: 2026-09-26T05:59:36Z
transitions:
  - to: ready
    at: 2026-09-26T05:24:31Z
    by: alex
  - to: in-progress
    at: 2026-09-26T05:59:36Z
    by: agent-S-0119
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
- [ ] The settings page lists every served project with its key, name, root, and connected state or last error from `flai serve`
- [ ] A system-flow project under an import root that is not served is shown with a Serve action, which registers it through the host channel. The switcher gains it without a reload
- [ ] Each served project has a Remove action, confirmed, that unregisters it through the channel, touches none of its files, and says how to add it back
- [ ] Serve and Remove are hostapi writes gated by the `settings` host action and journalled like import, and are refused with the reason when the action is not enabled for the project
- [ ] The switcher tells a project that is served but disconnected apart from one that is connected, and links to the settings page for why
- [ ] Tried end to end in a browser with scratch repositories. Docs in `docs/users` and `design/system/flaiover-dashboard.md`

## Tasks

## Notes

Builds on the host channel methods that S-0121 adds, so pull S-0121 first. Related to [I-0047](../../../design/issues/I-0047-a-project-imported-with-flai-import-on-the-command-line-is-neither-served-nor-offered-so-the-dashboard-never-shows-it.md).
