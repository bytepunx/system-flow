---
id: T-1517
type: task
nature: feature
title: The dashboard's document endpoints pass a story to doc.show and doc.save, and /api/docs/lint answers doc.lint
status: done
parent: S-0378
owner: alex
created: 2026-10-13T16:56:55Z
updated: 2026-10-13T18:06:19Z
transitions:
  - to: ready
    at: 2026-10-13T17:54:59Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T17:55:00Z
    by: agent-S-0378
  - to: done
    at: 2026-10-13T18:06:19Z
    by: agent-S-0378
stream: S-0378
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/routes/api/docs/edit/+server.ts, flaiover/src/routes/api/docs/file/+server.ts, flaiover/src/routes/api/docs/file/save.test.ts, flaiover/src/routes/api/docs/lint/+server.ts, flaiover/src/routes/api/docs/lint/lint.test.ts]
after: [T-1514, T-1515]
usage:
  source: log
  seconds: 679
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 90
      output: 442
      cache_read: 3665270
      cache_write: 121263
      cost: 1.7333
---
# T-1517 The dashboard's document endpoints pass a story to doc.show and doc.save, and /api/docs/lint answers doc.lint

## Work

- `GET /api/docs/edit?path=<p>&story=<S-nnnn>` and `PUT /api/docs/file` with `story` in its body pass `story` to `doc.show` and `doc.save` (T-1514). Without it they behave as today.
- Add `POST /api/docs/lint` with `{path, content}`, answering `doc.lint`'s findings (T-1515), and the function behind it in `flaiover/src/lib/server/repo.ts`.
- Treat `doc.lint` as optional in `REQUIRED_METHODS` (`flaiover/src/lib/server/agent.ts`), so a dashboard newer than the host's flai still connects. Without it the endpoint answers 501 and the editor shows no live findings.
- Waits for T-1514 and T-1515, whose host methods it calls.

## Done when

- `save.test.ts` covers a save with `story`, and its 409 and 422 answers. `lint.test.ts` covers findings, clean content, and a host without `doc.lint`.
- `scripts/flai.sh test` on these paths passes.

## Notes

- T-1518, the editor, uses these endpoints with the shapes above, so the two can run in the same layer.
