---
id: T-0306
type: task
nature: feature
title: One login token and one agent credential per user, kept in flai's home, and flai dashboard becomes multi-project aware
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:23Z
updated: 2026-09-20T21:17:33Z
transitions:
  - to: ready
    at: 2026-09-20T21:09:31Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T21:09:31Z
    by: system-flow
  - to: done
    at: 2026-09-20T21:17:33Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0306 One login token and one agent credential per user, kept in flai's home, and flai dashboard becomes multi-project aware

## Work
dashboard.token and dashboard.agent-key move from each project's .flai-cache to a shared location beside flai serve's state (the serve folder next to flai's config file). flai dashboard in any project ensures the one container (fixed name, one port) and the host flai are running, creates the shared secrets if missing, and registers the project; it prints that project's address. flai dashboard stop unregisters the project and stops the container only when it was the last one registered. flai dashboard status and flai dashboard token adapt. Existing per-project .flai-cache/dashboard.token and .agent-key are no longer read; flai dashboard says so once if it finds them.

## Done when
- Command tests: two projects share one container; stopping one leaves the other served; the shared secrets are created once and reused; flai dashboard status names every project the container serves
- The Go tests and lint pass

## Notes
