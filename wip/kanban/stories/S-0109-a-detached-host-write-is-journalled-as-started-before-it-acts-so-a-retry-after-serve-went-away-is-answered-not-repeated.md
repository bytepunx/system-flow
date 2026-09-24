---
id: S-0109
type: story
nature: remediation
title: A detached host write is journalled as started before it acts, so a retry after serve went away is answered, not repeated
status: ready
owner: alex
created: 2026-09-24T03:26:14Z
updated: 2026-09-24T06:35:16Z
transitions:
  - to: ready
    at: 2026-09-24T03:52:38Z
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

- [ ] hostapi records a detached host write under its request ID before it acts, so the same request asked again after serve went away is answered with what happened, or that it is under way, instead of acting again
- [ ] a serve restart from the dashboard, retried after the connection drops, restarts serve once, tried live

## Tasks

## Notes

Found by agent-S-0107 while trying S-0107 live, on a scratch host built from story/S-0106. A serve restart from the dashboard ran twice, 0.3 s apart. flaiover's `repo.write` retries a write once when the connection is lost and flai returns within 5 s. The serve that took the request died before its journal recorded it, so the new serve had no record and acted again. S-0107 sends `host.process` for serve or all, and `host.upgrade`, with `retry: false` (be41439); this story fixes the cause in flai. The journal is in memory today: surviving the restart means writing it beside flai serve's state.
