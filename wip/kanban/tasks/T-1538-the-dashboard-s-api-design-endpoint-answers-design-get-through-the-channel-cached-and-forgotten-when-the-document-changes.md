---
id: T-1538
type: task
nature: feature
title: The dashboard's /api/design endpoint answers design.get through the channel, cached and forgotten when the document changes
status: backlog
parent: S-0382
owner: alex
created: 2026-10-13T17:47:01Z
updated: 2026-10-13T17:47:01Z
transitions: []
stream: S-0382
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts]
---
# T-1538 The dashboard's /api/design endpoint answers design.get through the channel, cached and forgotten when the document changes

## Work

Add `design.get` to `REQUIRED_METHODS` in `agent.ts`, a `designSystem()` read in `repo.ts` with a cache key the `READS` table forgets on a `document` change, and `GET /api/design` in `flaiover/src/routes/api/design/+server.ts` answering it. Test the endpoint with a document, without one, and with flai's refusal. It waits for nothing in this story; it needs S-0381's `design.get` on the host.

## Done when

- `GET /api/design` answers the tokens, guidance, renderings, and problems, or that there is no document, and its test passes under `flai test`.

## Notes

- Layer 1, beside T-1539 and T-1540, with which it shares no file.
