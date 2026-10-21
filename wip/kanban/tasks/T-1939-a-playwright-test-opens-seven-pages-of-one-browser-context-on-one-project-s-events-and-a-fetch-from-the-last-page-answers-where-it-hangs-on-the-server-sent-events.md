---
id: T-1939
type: task
nature: remediation
title: A Playwright test opens seven pages of one browser context on one project's events and a fetch from the last page answers, where it hangs on the server-sent events
status: backlog
parent: S-0454
owner: alex
created: 2026-10-20T22:22:29Z
updated: 2026-10-20T22:22:29Z
transitions: []
stream: S-0454
tags: [dashboard, testing]
touches: [flaiover/playwright.config.ts, flaiover/e2e/live-events.e2e.ts, scripts/flaiover-e2e.sh, Makefile]
after: [T-1934, T-1936]
---
# T-1939 A Playwright test opens seven pages of one browser context on one project's events and a fetch from the last page answers, where it hangs on the server-sent events

## Work

Write `flaiover/e2e/live-events.e2e.ts`. In one Chromium browser context it opens at least seven pages, each following the same project's events through `listen` (a page such as the board, or one that only listens), waits until each has its stream open, and then fetches from the last page a route that answers at once, such as `/_health` or `/api/project`. The test passes when that fetch answers within a few seconds. Chromium's limit of six connections to one host holds under Playwright as in Brave, so seven event streams over HTTP/1.1 leave the fetch queued.

- **The server.** `flaiover/playwright.config.ts` runs `npm run build && npm run preview`, which is `vite preview`: it loads neither `flaiover/server.js` nor the dev server's upgrade plugin, so no WebSocket would upgrade. Run the built app with `node server.js` instead, on the port the config waits for, with the settings it needs to start without a login token, or with one the test sets as its cookie.
- **The project.** `/api/events` serves the project the URL names, and the dashboard shows nothing without a flai on the host (ADR-0029). Use the fake flai the server tests use (`flaiover/src/lib/server/project-events.test.ts`) dialing `/agent`, or a real `flai serve` against the template's sample repository; say which in the narrative's Decisions.
- **Running it.** Add `scripts/flaiover-e2e.sh`, which installs Playwright's Chromium when missing and runs `playwright test`, and a `flaiover-e2e` target in the `Makefile` that calls it. No tier in `system-flow.yaml` is added: that file's tiers are flai's to own.
- **Against the old transport.** Run the test once with `flaiover/src/lib/events.ts` from the main branch, in a scratch copy, and record in the narrative that the fetch hangs there. Then run it on the story's branch, where it passes.

It waits for T-1934 and T-1936, the server and the client it drives. The server-sent events must still exist for the run against them, so the removal waits for this task.

## Done when

- `scripts/flaiover-e2e.sh` passes on the story's branch.
- The narrative records the run against the server-sent events and that it hung.

## Notes

If Chromium under Playwright does not hold to six connections per host, the test proves nothing; check it in the run against the old transport before relying on it.
