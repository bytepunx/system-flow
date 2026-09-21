---
id: T-0318
type: task
nature: feature
title: "Tried end to end: a batch of several stories against one component, published once"
status: done
parent: S-0087
owner: alex
created: 2026-09-21T03:32:56Z
updated: 2026-09-21T04:44:44Z
transitions:
  - to: ready
    at: 2026-09-21T04:19:30Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T04:19:30Z
    by: system-flow
  - to: done
    at: 2026-09-21T04:44:44Z
    by: system-flow
stream: S-0087
tags: []
---
# T-0318 Tried end to end: a batch of several stories against one component, published once

## Work
Tried in a scratch project: several stories against the same component, of different natures, accepted (merged, no tag) one after another, then published once. Confirmed the batch's bump is the highest delivery type among them, not one release per story; the resulting tag(s); the changelog entry naming every story; a batch spanning more than three tags going out three at a time; a publish interrupted partway resuming correctly; the done column's published/waiting distinction throughout; and a research story in the same batch landing with no release of its own.

## Done when
- What was tried and what was not is in the narrative; every unexpected answer explained before review

## Notes
Last task, mirroring how S-0080 closed with a full end-to-end pass rather than trusting the unit tests of each task alone.
