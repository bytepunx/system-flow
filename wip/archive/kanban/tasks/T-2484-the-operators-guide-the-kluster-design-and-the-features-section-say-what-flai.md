---
id: T-2484
type: task
nature: feature
title: The operators guide, the kluster design, and the Features section say what flai serve runs while the kluster host action is on, and when
status: done
parent: S-0579
owner: alex
created: 2026-11-06T21:44:38Z
updated: 2026-11-07T03:17:34Z
transitions:
  - to: ready
    at: 2026-11-07T03:16:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T03:16:10Z
    by: system-flow
  - to: done
    at: 2026-11-07T03:17:34Z
    by: agent-S-0579
stream: S-0579
tags: [kluster, docs]
touches: [docs/operators/index.md, design/system/kluster.md, docs/users/flaiover.md]
after: [T-2482, T-2483]
usage:
  source: log
  seconds: 84
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 11141
      cache_read: 4067544
      cache_write: 49029
      cost: 1.3811
---
# T-2484 The operators guide, the kluster design, and the Features section say what flai serve runs while the kluster host action is on, and when

## Work

It waits for T-2482 and T-2483, whose behaviour it documents as built.

- `docs/operators/index.md`: a section on the kluster host action (S-0579). Say what flai serve deploys, tears down, and promotes, and when. Say how checks see the namespace, what a failure does, where the logs are, and how to turn the action on and off. Extend the short entry S-0577 wrote, rather than adding a second one.
- `design/system/kluster.md`: fill in its section on flai serve.
- `docs/users/flaiover.md` § Settings › Features: one sentence that the `kluster` host action is listed there and what turning it on does, linking the operators guide.

Bump `updated` on each.

## Done when

- The three documents say what the story built.
- `flai check --strict` and the markdown lint pass on the three paths.

## Notes
