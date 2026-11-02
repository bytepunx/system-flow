---
id: T-2438
type: task
nature: feature
title: The dashboard's acceptance drops the checkbox for uncommitted files, shows the shared set it commits, and shows a send back as such
status: done
parent: S-0566
owner: alex
created: 2026-11-02T18:36:14Z
updated: 2026-11-02T19:02:26Z
transitions:
  - to: ready
    at: 2026-11-02T18:36:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T18:36:55Z
    by: system-flow
  - to: done
    at: 2026-11-02T19:02:26Z
    by: agent-S-0566
stream: S-0566
tags: []
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/AcceptConfirm.svelte, flaiover/src/lib/components/AcceptConfirm.svelte.test.ts, "flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts", "flaiover/src/routes/api/items/[id]/move/+server.ts", flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/items/[id]/+page.svelte"]
usage:
  source: log
  seconds: 1492
  estimated: true
  priced_by: estimate
  models:
    - model: claude-fable-5-1
      input: 623
      output: 291
      cache_read: 2678293
      cache_write: 152010
      cost: 4.2922
---
# T-2438 The dashboard's acceptance drops the checkbox for uncommitted files, shows the shared set it commits, and shows a send back as such

## Work

ADR-0248 (`design/adrs/0248-*.md`) in the dashboard. Today `AcceptConfirm.svelte` and `Review.svelte` show the preview's `uncommitted` paths ("Uncommitted changes outside wip") with a checkbox "Include these files in the acceptance commit", and send `include_uncommitted` to `/api/items/[id]/accept` (and the move route for a drop on done), which the host write `accept.run` turns into `flai accept --yes`. After S-0566's Go tasks the preview answers `shared_uncommitted` (the main checkout's uncommitted files under wip, design, and the shared paths, which the acceptance commits) and `uncommitted` (the rest, left as they are, refused only when the story's change would overwrite one, as a `dirty_main` blocker in `blockers`).

- Remove the checkbox, the `include` state, `needsChoice`, and the `include_uncommitted` field from the request body in both components, the board page, the item page, and the two API routes (keep the routes tolerant of an absent field; drop the test cases that send it, or turn them into tests that it is no longer sent). The `onconfirm` callback of `AcceptConfirm` loses its `include` argument.
- Where the checkbox was, show `shared_uncommitted` as a muted note, "Committed with the acceptance:" with the paths, and `uncommitted` as an informational note, "Left as they are:" with the paths and one line saying the acceptance refuses only when the story would overwrite one. Neither is a warning.
- A refusal whose data carries `sent_back: true` (the JSON refusal `{"refused": {"id", "reason", "sent_back", "branch", "base", "conflicts", "thread"}}`, which reaches the review page as the `error` event's `data`): show it as a notice, not a danger alert, headed "Sent back to in-progress", with the reason, the conflicting paths, and a link to the thread (`/threads/TH-nnnn`, as the page links threads elsewhere), and refresh the item as a failure does today. The board page's drop on done gets the same refusal through the move route: show its reason.
- Update the tests in the four test files to the new props and texts, with vitest; add one for the send-back notice. Run `flai test` on the paths you changed (prettier, eslint, vitest, svelte-check).

A write to a path Claude Code protects is never yours: none is expected here. This task waits for no other: it is layer 1; the Go field names are fixed by the ADR and the task bodies.

## Done when

- No checkbox, `include` state, or `include_uncommitted` request field remains in `flaiover/src`; `AcceptConfirm`'s `onconfirm` takes no argument.
- The review page and the confirmation show `shared_uncommitted` as committed with the acceptance and `uncommitted` as left alone.
- A `sent_back` refusal shows as a notice naming the thread.
- `flai test` passes on the changed paths.

## Notes

Acceptance criteria 1, 2, and 3 of S-0566, in the dashboard.
