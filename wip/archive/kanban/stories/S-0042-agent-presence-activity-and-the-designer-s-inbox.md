---
id: S-0042
type: story
nature: feature
title: Agent presence, activity, and the designer's inbox
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-19T05:10:24Z
transitions:
  - to: ready
    at: 2026-09-19T01:56:17Z
    by: alex
  - to: in-progress
    at: 2026-09-19T04:44:55Z
    by: system-flow
  - to: review
    at: 2026-09-19T04:59:48Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:10:24Z
    by: alex
tags: [dashboard, cli]
---

# S-0042 Agent presence, activity, and the designer's inbox

## Goal
The designer sees which agents are active and what they are doing, and has one inbox of things that need a human: open questions, threads awaiting an answer, stories in review, blocked items, and overlapping touches.

## Acceptance criteria
- [x] An activity view built from `wip/agents/index.md` and the narratives: active streams, agent and session, last log entry and its age, current task, blocked flags
- [x] An inbox view listing open threads and questions, stories in review, blocked items, and `wip.overlap` warnings, each linking to its page; a badge in the navigation with the count
- [x] Optional desktop notifications and a webhook (`dashboard.notify_url`) for new inbox entries, off by default
- [x] Tests on fixtures; docs/users updated

## Tasks
- T-0191 Living design for activity, the designer's inbox, and notifications
- T-0192 Server: activity and inbox built from the files, with endpoints and fixture tests
- T-0193 Activity and inbox pages, and the inbox count in the navigation
- T-0194 Notifications for new inbox entries: desktop in the browser, and a webhook from dashboard.notify_url
- T-0195 User and operator documentation for activity, inbox, and notifications
- T-0196 Verify: all tiers, the pages in a real container, and the webhook against a local listener

## Notes
Presence is derived from files, so nothing new is written; an agent that stops logging simply ages out of the view.

Decided when pulled, 2026-09-19. Overlaps come from `flai check --json`, cached until the repository changes, so the rule has one implementation; without flai they are left out and the inbox says so. A thread awaits the designer while its last entry is not by the manifest's owner. Notifications, in the browser and by webhook, fire only for entries that appear after the page or the server started. The webhook is read once at server start. The activity view reads the narratives themselves; `wip/agents/index.md` is generated from the same files and adds nothing, so it is not read.

Verification, 2026-09-19. `make flai-test` passed (golangci-lint 0 issues; behavior, integration, smoke; markdown lint); flaiover lint, svelte-check 0 errors, 26 files and 162 tests, and a production build passed. In a container built from the branch, against a scratch project with one of each: `/api/activity` returned three streams with agent `claude`, session `abc123`, the task in progress, the blocked flag, and the last log entry; `/api/inbox` returned five entries, one of each kind, each linking to its page, the review entry to `/review/S-0003`. In a browser against that container: the inbox page grouped them with reviews first, the badge read 5, the desktop notification toggle was present and off; after the designer answered the thread from the host, the badge read 4 and the thread left the list without a reload; the activity page rendered; at 390 pixels nothing scrolled sideways, which also closes I-0018. Webhook, through the dev server on the host with `dashboard.notify_url` pointing at a local listener: nothing was posted for the four entries present at start; blocking a story produced exactly one POST with `{ project, entry: { key, kind, title, href, at } }` and no authorization or cookie header; an unrelated change produced none; the URL's query string did not appear in the server's log. Not exercised: a desktop notification actually raised by a browser, which needs a permission prompt a headless browser cannot answer; the detection it depends on is unit tested.
