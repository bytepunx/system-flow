---
id: T-0550
type: task
nature: feature
title: An inbox thread entry links to its story's page naming the thread
status: done
parent: S-0155
owner: alex
created: 2026-09-29T07:02:04Z
updated: 2026-09-29T07:03:01Z
transitions:
  - to: ready
    at: 2026-09-29T07:02:22Z
    by: agent-S-0155
  - to: in-progress
    at: 2026-09-29T07:02:22Z
    by: agent-S-0155
  - to: done
    at: 2026-09-29T07:03:01Z
    by: agent-S-0155
stream: S-0155
tags: []
usage:
  source: log
  seconds: 39
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 19
      output: 3867
      cache_read: 933582
      cache_write: 19203
      cost: 0.4178
---

# T-0550 An inbox thread entry links to its story's page naming the thread

## Work

`hrefFor` in `flaiover/src/lib/server/inbox.ts` sends a thread entry with an item to `/items/<item>?thread=<TH-id>`, the thread ID taken from the entry's key. A thread on a document with no item keeps its document page.

## Done when

`lib/server/inbox.test.ts` pins the new link for a thread on a story, and the document link for a thread on a document; `npm run check`, lint, and the unit tests pass.

## Notes
