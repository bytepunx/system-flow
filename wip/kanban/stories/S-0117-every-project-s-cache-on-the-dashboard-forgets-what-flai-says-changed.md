---
id: S-0117
type: story
nature: remediation
title: Every project's cache on the dashboard forgets what flai says changed
status: ready
parent: E-0003
owner: arobson
created: 2026-09-26T02:42:25Z
updated: 2026-09-26T02:42:25Z
transitions:
  - to: ready
    at: 2026-09-26T02:42:25Z
    by: claude
tags: [dashboard]
touches: [flaiover/src, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0117 Every project's cache on the dashboard forgets what flai says changed

## Goal

The dashboard shows what the repository says. Today it keeps one cache of flai's answers per project, but only the default project's cache listens for flai serve's change notices (`flaiover/src/hooks.server.ts` calls `watch()` on `repo()` once, outside any request). A request that names a project (`?project=sf`) gets its own Repo from `repo()` in `flaiover/src/lib/server/repo.ts`, and nothing ever calls `watch()` on it, so its answers are kept until the container restarts. Stories accepted and archived on another machine and pulled here stayed in the done lane: `/api/board?project=sf` listed S-0028, S-0115, and S-0116 while `/api/board`, `flai board`, and `wip/kanban/stories` showed none.

## Acceptance criteria
- [ ] every project's Repo made by `repo()` listens for flai's change, connected, and gone events, not only the default project's
- [ ] a change flai reports for a named project clears that project's cached answers, pinned by a test that asks the board for a project, reports a removed story file, and asks again
- [ ] a Repo given its own source (tests) still does not listen, and calling `watch()` twice still registers one set of listeners
- [ ] the defect is recorded in `design/issues/` and `summary.md`

## Tasks

## Notes

- Suggested fix: in `repo()`, call `void r.watch()` when a Repo is made for a key. `watch()` is idempotent and a no-op for a Repo with its own source.
- Workaround until released: `docker restart flaiover`.
