---
id: T-3544
type: task
nature: remediation
title: kluster.md and the user guide say the review runs the checks alone when only they have not passed at the head, and I-0362 is closed
status: backlog
parent: S-0835
owner: alex
created: 2026-11-17T08:31:05Z
updated: 2026-11-17T08:31:05Z
transitions: []
stream: S-0835
tags: [flai, docs]
touches: [design/system/kluster.md, docs/users/flai.md, design/issues/I-0362-flai-serve-s-review-pass-deploys-again-a-story-in-review-whose-deploy-passed-at.md]
after: [T-3542]
---
# T-3544 kluster.md and the user guide say the review runs the checks alone when only they have not passed at the head, and I-0362 is closed

## Work

Waits for T-3542, whose behaviour this describes and whose fix the issue's close names.

1. `design/system/kluster.md` § What happens when: in the row "A story enters review otherwise", make the Built column say the checks run alone is built, naming S-0835, in place of the pass deploying again. Remove the "Not built" line that says the same further down. In § At the head, say that the review's checks run alone takes a record at the head as `kluster.AtHead` judges it, while a checks request a story's agent makes keeps the exact head. Bump `updated`.
2. `docs/users/flai.md`, where it says a story entering review is deployed again in place unless an asked-for deploy passed at its head with its checks: say that when only the checks have not passed at the head, the review runs the checks alone against the namespace as deployed. Bump `updated`.
3. Close the issue: `flai issue close I-0362 --reason` saying the review pass runs the checks alone when only they have not passed at the head, with the test that reproduces it, in S-0835.

## Done when

- `flai test` passes on the changed paths, the markdown lint among them.
- `flai issue show I-0362` reads closed with the reason, and no line of `design/system/kluster.md` says the checks run alone is not built for the review.

## Notes
