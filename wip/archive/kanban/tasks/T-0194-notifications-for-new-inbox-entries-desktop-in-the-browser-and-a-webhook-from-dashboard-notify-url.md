---
id: T-0194
type: task
nature: feature
title: "Notifications for new inbox entries: desktop in the browser, and a webhook from dashboard.notify_url"
status: done
parent: S-0042
owner: alex
created: 2026-09-19T04:46:17Z
updated: 2026-09-19T04:53:26Z
transitions:
  - to: ready
    at: 2026-09-19T04:50:43Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:50:44Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:53:26Z
    by: system-flow
stream: S-0042
tags: []
touches: [flaiover/src, flai/internal/manifest]
---

# T-0194 Notifications for new inbox entries: desktop in the browser, and a webhook from dashboard.notify_url

## Work
Desktop: a toggle on the inbox page, stored in the browser, off by default; when on and permission is granted, a new inbox entry (a key not seen before in this tab) raises a `Notification` with the entry's title; never for entries present when the page loaded. Webhook: add `notify_url` to the manifest's dashboard section in Go (and nowhere else in flai); in the dashboard's server start a notifier at init when the manifest has the key: on repository change events, debounce, rebuild the inbox, and POST `{ project, entry: { key, kind, title, href, at } }` as JSON for each key not seen before, with a short timeout, no retries, a warning log on failure, and nothing sent for entries present at start. The token, file contents, and paths outside the entry's link are never sent. Tests: the notifier against a local HTTP server on an ephemeral port (new entry posted once, existing ones not, failure logged and not thrown), and the client's new-entry detection as a pure function.

## Done when
- The notifier and detection tests pass
- With no `notify_url`, no notifier starts and no request is made (tested)

## Notes
