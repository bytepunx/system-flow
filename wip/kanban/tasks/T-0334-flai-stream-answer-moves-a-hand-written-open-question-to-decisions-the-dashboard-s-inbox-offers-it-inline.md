---
id: T-0334
type: task
nature: feature
title: "flai stream answer moves a hand-written open question to Decisions; the dashboard's inbox offers it inline"
status: done
parent: S-0090
owner: alex
created: 2026-09-23T00:27:54Z
updated: 2026-09-23T00:28:28Z
transitions:
  - to: ready
    at: 2026-09-23T00:28:11Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T00:28:13Z
    by: system-flow
  - to: done
    at: 2026-09-23T00:28:28Z
    by: system-flow
stream: S-0090
tags: []
---

# T-0334 flai stream answer moves a hand-written open question to Decisions; the dashboard's inbox offers it inline

## Work
A hand-written bullet under a narrative's `## Open questions` has no thread of its own to reply to, so it was only ever answerable by editing the document by hand. `flai stream answer <story-id> "<question>" "<answer>"` (`internal/workitem/narrative.go`'s `AnswerOpenQuestion`, exported `OpenQuestions` alongside it) removes the matching bullet — exactly as `OpenQuestions` reads it, continuation lines joined — and records it, with the answer, under `## Decisions` (`design/system/agent-narrative.md`: "Answered questions move to Decisions"), refusing when nothing matches. `stream.answer` exposes it over the channel the same shape as `thread.reply` (`--by=owner(p)`, `--` before the free text).

`inbox.designer`'s question entries now carry `Item` (the story), so the dashboard can target the right narrative without parsing it out of `path`. `InboxView.svelte` gets an inline answer form per question entry, gated on `writable`, posting to a new `POST /api/streams/:id/answer`; on success it calls `inboxState.refresh()` so the entry drops off without waiting for the file-change SSE round trip.

## Done when
Criteria 1 and 3 (criterion 2 is S-0089's, not retested here — see the story's Notes): a Go test proves `AnswerOpenQuestion` moves the bullet to Decisions and refuses an unmatched or already-answered question; `writes_test.go`'s `good`/`refused` tables cover `stream.answer`'s command line; a real end-to-end test (`writes.test.ts`, the real flai binary) proves the narrative file changes; `Inbox.svelte.test.ts` covers the form hidden when not writable, a successful answer posting and refreshing, and a refusal shown verbatim. `go test ./...`, `pnpm run check`/`test:unit`/`lint`, `flai check --strict` clean.

## Notes
`flai stream answer` and `stream.answer` were named to match the existing `flai stream log`/`stream.log` pair (a "narrative.answer" name was considered and dropped for consistency).
