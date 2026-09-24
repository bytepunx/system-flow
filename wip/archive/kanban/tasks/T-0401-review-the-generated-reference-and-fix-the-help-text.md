---
id: T-0401
type: task
nature: feature
title: Review the generated reference and fix the help text
status: done
parent: S-0017
owner: alex
created: 2026-09-24T08:12:55Z
updated: 2026-09-24T08:17:41Z
transitions:
  - to: ready
    at: 2026-09-24T08:13:04Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T08:16:02Z
    by: system-flow
  - to: done
    at: 2026-09-24T08:17:41Z
    by: system-flow
stream: S-0017
tags: [cli]
touches: [flai/cmd, docs/users]
---
# T-0401 Review the generated reference and fix the help text

## Work
Generate `docs/users/flai-reference.md`, read it end to end, and fix defects in the cobra help text it exposes (typos such as "storys", stale wording, missing descriptions, inconsistent flag help). Regenerate until it reads correctly.

## Done when
Every command in the tree appears in the reference with a description, the reviewed help text is committed with the regenerated file, and tests pass.

## Notes
