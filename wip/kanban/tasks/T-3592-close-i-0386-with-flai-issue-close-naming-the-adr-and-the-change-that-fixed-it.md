---
id: T-3592
type: task
nature: improvement
title: Close I-0386 with flai issue close, naming the ADR and the change that fixed it
status: backlog
parent: S-0845
owner: alex
created: 2026-11-17T12:34:33Z
updated: 2026-11-17T12:34:33Z
transitions: []
stream: S-0845
tags: [issues]
touches: [design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md, design/issues/summary.md]
after: [T-3591]
---
# T-3592 Close I-0386 with flai issue close, naming the ADR and the change that fixed it

## Work

Waits for T-3591, the last change of the fix, so the reason names all of it. Run `flai issue close I-0386 --reason "<reason>"`, the reason naming T-3587's ADR, the orchestrator's `resolve` request under `accept_reviews`, the test T-3590 wrote, and that a story held by `undone_main` stays the operator's. Tick the story's second criterion with `flai criteria tick S-0845 2`.

## Done when

- I-0386 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
