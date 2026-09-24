---
id: S-0109
type: story
nature: remediation
title: A detached host write is journalled as started before it acts, so a retry after serve went away is answered, not repeated
status: done
owner: alex
created: 2026-09-24T03:26:14Z
updated: 2026-09-24T07:52:20Z
transitions:
  - to: ready
    at: 2026-09-24T03:52:38Z
    by: alex
  - to: in-progress
    at: 2026-09-24T07:35:00Z
    by: agent-S-0109
  - to: review
    at: 2026-09-24T07:46:38Z
    by: agent-S-0109
  - to: done
    at: 2026-09-24T07:52:20Z
    by: alex
tags: [cli]
touches: [flai/internal/hostapi]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0109 A detached host write is journalled as started before it acts, so a retry after serve went away is answered, not repeated

## Goal

A dashboard write that ends the flai serve answering it (a serve restart, an upgrade) is not done twice when the dashboard retries it (TH-0006).

## Acceptance criteria

- [x] hostapi records a detached host write under its request ID before it acts, so the same request asked again after serve went away is answered with what happened, or that it is under way, instead of acting again
- [x] a serve restart from the dashboard, retried after the connection drops, restarts serve once, tried live

## Tasks
- T-0392 hostapi records a detached write before it acts, beside flai serve's state
- T-0393 Docs say a repeat of a detached write is answered from a record that outlives flai serve
- T-0394 A serve restart from the dashboard, retried after the connection drops, restarts serve once, tried live

## Notes

Found by agent-S-0107 while trying S-0107 live, on a scratch host built from story/S-0106. A serve restart from the dashboard ran twice, 0.3 s apart. flaiover's `repo.write` retries a write once when the connection is lost and flai returns within 5 s. The serve that took the request died before its journal recorded it, so the new serve had no record and acted again. S-0107 sends `host.process` for serve or all, and `host.upgrade`, with `retry: false` (be41439); this story fixes the cause in flai. The journal is in memory today: surviving the restart means writing it beside flai serve's state.

Verified live on 2026-09-24, on a scratch `flai host` built from `story/S-0109` (its own configuration, address 127.0.0.1:4291, and project) and a scratch flaiover dev server for that project, with the retry turned on for `host.process` serve in the dev server's copy only (not committed). `POST /api/host {"action":"restart","process":"serve"}`:

- The host log has one `process restart` for serve and one new serve (pid 2690756, then 2691311), and none in the ten seconds after.
- The dashboard lost the connection at 07:43:42.755Z, the new serve connected at 07:43:42.914Z, and the retry was answered 200 in 0.33 s with `{"repeat": {"under_way": false}}` and the line "started ... by a flai serve (pid 2690756) that ended before it could say how it went; it was not done again".
- `serve/requests.json` held the record the first serve wrote before acting, `done: false`, pid 2690756.
- Seen on the way, not changed here: a serve restart leaves no line in `serve/journal.jsonl`, before this story as after, because the serve that would write it is the one that ends. The record in `requests.json` is the only trace.
