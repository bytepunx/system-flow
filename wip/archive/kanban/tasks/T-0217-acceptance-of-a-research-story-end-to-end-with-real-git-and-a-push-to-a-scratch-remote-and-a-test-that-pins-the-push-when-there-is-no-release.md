---
id: T-0217
type: task
nature: feature
title: Acceptance of a research story end to end with real git and a push to a scratch remote, and a test that pins the push when there is no release
status: done
parent: S-0053
owner: alex
created: 2026-09-19T07:15:38Z
updated: 2026-09-19T07:20:22Z
transitions:
  - to: ready
    at: 2026-09-19T07:20:22Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:20:22Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:20:22Z
    by: system-flow
stream: S-0053
tags: []
---

# T-0217 Acceptance of a research story end to end with real git and a push to a scratch remote, and a test that pins the push when there is no release

## Work
A behaviour test with real git and a bare scratch remote: a research story in review with a story branch is accepted with `flai accept`; the branch is merged, the story archived, the acceptance committed, no tag created, no version file or changelog changed, and the remote's main has the acceptance commit. The same through `flai move <story> done`. A second test pins that the push happens whenever a plan cuts no release (a story that touched no component), since nothing pins it today. An experiment story is refused before anything is merged (the I-0021 rule: refuse before acting).

## Done when
- The end to end tests pass and fail without the change
- `make flai-test` passes

## Notes
