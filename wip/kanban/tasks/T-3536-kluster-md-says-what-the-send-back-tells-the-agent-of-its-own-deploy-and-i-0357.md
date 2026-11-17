---
id: T-3536
type: task
nature: improvement
title: kluster.md says what the send-back tells the agent of its own deploy, and I-0357 is closed
status: backlog
parent: S-0831
owner: alex
created: 2026-11-17T08:29:56Z
updated: 2026-11-17T08:29:56Z
transitions: []
stream: S-0831
tags: [design, kluster, issues]
touches: [design/system/kluster.md, design/issues/I-0357-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md]
after: [T-3533, T-3534]
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
