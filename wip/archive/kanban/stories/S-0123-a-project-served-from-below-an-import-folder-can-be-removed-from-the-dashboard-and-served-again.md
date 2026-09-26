---
id: S-0123
type: story
nature: improvement
title: A project served from below an import folder can be removed from the dashboard and served again
status: done
parent: E-0003
owner: alex
created: 2026-09-26T06:58:59Z
updated: 2026-09-26T07:27:34Z
transitions:
  - to: ready
    at: 2026-09-26T07:04:27Z
    by: alex
  - to: in-progress
    at: 2026-09-26T07:05:03Z
    by: agent-S-0123
  - to: review
    at: 2026-09-26T07:23:41Z
    by: agent-S-0123
  - to: done
    at: 2026-09-26T07:27:34Z
    by: alex
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
- [x] flai serve keeps a list of project folders the operator removed; a project in that list is not served from below an import folder or the start folder, and `flai serve project list` shows it as not served, with why
- [x] `flai serve project remove` on a project served from below a folder adds it to that list, and `flai serve project add` takes it off and serves it
- [x] The settings page offers Remove on a project served from below a folder, and Serve on one in the list; after Serve, the switcher gains it without a reload
- [x] Docs in `docs/users`, `docs/operators`, and `design/system`, and tried end to end in a browser with scratch repositories

## Tasks
- T-0451 flai serve keeps a list of the project folders the operator removed and does not serve them from below a folder
- T-0452 flai serve project remove adds a project served from below a folder to the removed list, and add takes it off
- T-0453 The settings page removes a project served from below a folder and serves one in the removed list
- T-0454 Docs for removing and serving again a project below a folder, and an end-to-end try in a browser with scratch repositories

## Notes

From TH-0016 on S-0122 (answer A), and TH-0014 before it, which left the list of removed roots for later work. Since S-0120, flai serve serves every project with a key below an import folder, so S-0122's Serve can only return flai's refusal for the projects listed as not served: no key, a key already served, or no dashboard known. This story gives Serve something to undo.

Tried end to end on 2026-09-26 (T-0454), in headless Chrome against a flaiover dev server built from `story/S-0123` on port 5199, and a scratch flai host and serve from the same branch, with their own `HOME`, configuration, and `FLAI_HOST_ADDR` (127.0.0.1:4291), and no `FLAI_HOST_*` in their environment. The scratch repositories were alpha, registered, and beta and gamma, imported with `flai import --commit` below a scratch import folder and not registered. `settings` was on for alpha and beta and off for gamma. The imported manifests first named port 4242, the operator's dashboard; the scratch serve dialled it for a few seconds and was refused at the credential proof, and the manifests were moved to 5199 before anything else. Every scratch process was stopped by PID afterwards; the operator's host and dashboard were untouched and answered healthy.

- Projects listed alpha, beta, and gamma, connected; beta and gamma said they are served because they are below the import folder, each with Remove, gamma's disabled with "run flai serve enable settings on the host, in <gamma's folder>".
- Remove on beta asked first ("It is below <folder>, so it is listed below as removed, and Serve there brings it back"). Once confirmed, the switcher went from alpha, beta, gamma to alpha, gamma, beta moved under "not served" as "Removed from the dashboard; Serve brings it back", `removed.json` held beta's folder, and beta's `git status` was clean.
- Serve on beta said "b is served, and in the switcher"; the switcher had alpha, beta, gamma again, beta was listed as connected, and `removed.json` was empty. The page did not reload at any point.
- On the shell, `flai serve project remove g` listed gamma as removed with how to serve it again, and `flai serve project add git/gamma` served it from below the folder again, connected within seconds, without registering it.
