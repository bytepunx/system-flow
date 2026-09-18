---
id: T-0150
type: task
nature: improvement
title: Verify all tiers and strict check on this repository and a rendered template
status: done
parent: S-0049
owner: alex
created: 2026-09-18T18:08:30Z
updated: 2026-09-18T18:19:43Z
transitions:
  - to: ready
    at: 2026-09-18T18:15:03Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:15:04Z
    by: alex
  - to: done
    at: 2026-09-18T18:19:43Z
    by: alex
stream: S-0049
tags: []
---

# T-0150 Verify all tiers and strict check on this repository and a rendered template

## Work
Run `make test`, `make integration`, and `make smoke`, the linter through `make flai-test`, and `scripts/flai.sh check --strict`. Exercise the rule end to end in a scratch project rendered from the template: a story with criteria and no tasks moves to ready and to in-progress, passes `flai check --strict`, is refused at review, and is accepted at review once a task exists. Tick the story's acceptance criteria only for what was observed.

## Done when
- All three tiers and lint pass, with output recorded in the narrative
- The scratch walk-through behaved as the criteria say
- Every acceptance criterion on S-0049 is checked, or left unchecked with the reason in the story notes

## Notes
