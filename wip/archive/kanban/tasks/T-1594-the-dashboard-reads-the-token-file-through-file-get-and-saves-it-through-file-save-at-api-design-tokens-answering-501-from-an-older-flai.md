---
id: T-1594
type: task
nature: feature
title: The dashboard reads the token file through file.get and saves it through file.save at /api/design/tokens, answering 501 from an older flai
status: done
parent: S-0383
owner: alex
created: 2026-10-15T02:38:58Z
updated: 2026-10-15T02:49:03Z
transitions:
  - to: ready
    at: 2026-10-15T02:39:21Z
    by: agent-S-0383
  - to: in-progress
    at: 2026-10-15T02:39:22Z
    by: agent-S-0383
  - to: done
    at: 2026-10-15T02:49:03Z
    by: agent-S-0383
stream: S-0383
tags: []
touches: [flaiover/src/routes/api/design/tokens/+server.ts, flaiover/src/routes/api/design/tokens/tokens.test.ts, flaiover/src/lib/server/repo.ts]
usage:
  source: log
  seconds: 581
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 20803
      cache_read: 1894798
      cache_write: 74213
      cost: 1.204
---
# T-1594 The dashboard reads the token file through file.get and saves it through file.save at /api/design/tokens, answering 501 from an older flai

## Work

`flaiover/src/routes/api/design/tokens/+server.ts` is the token file's own endpoint (ADR-0147):

- `GET` reads it through the host method `file.get` `{path}`, at the path `design.get` names, and answers its raw content and SHA-256.
- `PUT` saves it through `file.save` `{path, content, hash, message?}` on the main branch.

`flaiover/src/lib/server/repo.ts` gains the two calls, mapped as the other host methods are:

- 501 when the host flai is older than `file.get` and `file.save`
- 409 for a stale hash, with the current content
- 422 for a refusal, with flai's findings, such as `design.tokens`

It waits for nothing in this story. The flai side of `file.get` and `file.save` is not in this story (TH-0454).

## Done when

- `tokens.test.ts` beside the endpoint passes for a read, a save, a 501, a conflict, and a refusal, against a stubbed host.

## Notes

- Added at the start because the story's save cannot go through `doc.save`: the tokens are a DTCG file (ADR-0146). T-1546 waits for it.
