---
id: T-3536
type: task
nature: improvement
title: kluster.md says what the send-back tells the agent of its own deploy, and I-0357 is closed
status: done
parent: S-0831
owner: alex
created: 2026-11-17T08:29:56Z
updated: 2026-11-17T09:15:21Z
transitions:
  - to: ready
    at: 2026-11-17T09:14:43Z
    by: agent-S-0831
  - to: in-progress
    at: 2026-11-17T09:14:43Z
    by: agent-S-0831
  - to: done
    at: 2026-11-17T09:15:21Z
    by: agent-S-0831
stream: S-0831
tags: [design, kluster, issues]
touches: [design/system/kluster.md, design/issues/I-0357-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md]
after: [T-3533, T-3534]
usage:
  source: log
  seconds: 35
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 59
      cache_read: 1548655
      cache_write: 6856
      cost: 0.5933
---
# T-3536 kluster.md says what the send-back tells the agent of its own deploy, and I-0357 is closed

## Work

In `design/system/kluster.md` § The send-back (S-0694), rewrite what the thread's entry says and the bullet **The agent**, on what `harness.SentBack.Kluster` tells the agent, to match the texts T-3533 and T-3534 wrote. They say the agent may deploy its own story at its head with `kluster_deploy` or `flai kluster deploy S-nnnn --job` before review. They say a pass at the head resolves the thread, and that back in review only what has not passed at the head is deployed again. Make the bullet **Back in review** say that last point too, as the section on the namespace kept since S-0660 already does. Bump the document's `updated`.

Then close the issue with `flai issue close I-0357 --reason "<what fixed it>"`, naming S-0831, the two texts, and the tests that check them.

It waits for T-3533 and T-3534: it describes the words they write and names their tests.

## Done when

- `design/system/kluster.md` § The send-back says what both texts now tell the agent, and the markdown lint passes on it.
- `flai issue close I-0357` has run, its reason naming S-0831, `SentBack.then`, `heldText`, and their tests.

## Notes
