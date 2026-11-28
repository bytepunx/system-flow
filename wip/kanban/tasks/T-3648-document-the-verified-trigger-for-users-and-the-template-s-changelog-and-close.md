---
id: T-3648
type: task
nature: improvement
title: Document the verified trigger for users and the template's changelog, and close I-0413 with what fixed it
status: backlog
parent: S-0854
owner: alex
created: 2026-11-27T23:59:54Z
updated: 2026-11-27T23:59:54Z
transitions: []
stream: S-0854
tags: [flai, template]
touches: [docs/users/flai.md, template/CHANGELOG.md, design/issues/I-0413-a-story-in-review-whose-close-out-or-resolve-job-passed-is-not-a-reason-that.md, design/issues/summary.md]
after: [T-3645, T-3646]
---
# T-3648 Document the verified trigger for users and the template's changelog, and close I-0413 with what fixed it

## Work

This task waits for T-3645 and T-3646. It documents what they built and closes I-0413 only once the fix and its test are in.

- `docs/users/flai.md`: add the `verified` event and the `passed-in-review` condition to the role trigger table. Say in the `review:S-nnnn` reason's row that a passed close-out or resolve job of a story in review gives that reason too, and that a verdict answered from the record does not. Link the ADR T-3641 wrote.
- `template/CHANGELOG.md`: under `## Unreleased`, one line saying that the orchestrator's role file gains the trigger, naming S-0854 and the ADR.
- Close I-0413 with `flai issue close I-0413 --reason "<reason>"`. The reason names the trigger, the ADR, and the test that reproduces the instance. `flai issue close` writes the issue file and `design/issues/summary.md`.

## Done when

- `docs/users/flai.md` and `template/CHANGELOG.md` describe the trigger and link the ADR.
- I-0413 is closed, with a reason saying what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
