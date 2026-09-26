---
id: S-0123
type: story
nature: improvement
title: A project served from below an import folder can be removed from the dashboard and served again
status: backlog
parent: E-0003
owner: alex
created: 2026-09-26T06:58:59Z
updated: 2026-09-26T06:58:59Z
transitions: []
tags: [cli, dashboard]
touches: [flai/internal/serve, flai/cmd, flaiover]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0123 A project served from below an import folder can be removed from the dashboard and served again

## Goal

The operator can stop the dashboard showing a project that is served because it is below an import folder, without removing the whole folder, and can bring it back with Serve on the settings page.

## Acceptance criteria
- [ ] flai serve keeps a list of project folders the operator removed; a project in that list is not served from below an import folder or the start folder, and `flai serve project list` shows it as not served, with why
- [ ] `flai serve project remove` on a project served from below a folder adds it to that list, and `flai serve project add` takes it off and serves it
- [ ] The settings page offers Remove on a project served from below a folder, and Serve on one in the list; after Serve, the switcher gains it without a reload
- [ ] Docs in `docs/users`, `docs/operators`, and `design/system`, and tried end to end in a browser with scratch repositories

## Tasks

## Notes

From TH-0016 on S-0122 (answer A), and TH-0014 before it, which left the list of removed roots for later work. Since S-0120, flai serve serves every project with a key below an import folder, so S-0122's Serve can only return flai's refusal for the projects listed as not served: no key, a key already served, or no dashboard known. This story gives Serve something to undo.
