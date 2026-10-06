---
id: T-0998
type: task
nature: remediation
title: permission_prompt takes an answer from the story's owner or the project's owner, with a test reproducing I-0081
status: backlog
parent: S-0284
owner: alex
created: 2026-10-06T06:24:14Z
updated: 2026-10-06T06:24:14Z
transitions: []
stream: S-0284
tags: [flai]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go]
after: [T-0997]
---
# T-0998 permission_prompt takes an answer from the story's owner or the project's owner, with a test reproducing I-0081

## Work

Build the fix T-0997's ADR records in `flai/internal/mcpserver/permission.go`.

- `askOperator` gathers who may answer: the story's owner and `s.repo.Manifest.Owner`, each when set and not the asking agent, without duplicates.
- `awaitAnswer` takes that set in place of the one owner: an entry after the question by anyone in it is the answer; with the set empty, any entry not by the agent is, as now.
- `permissionRequest` names who may answer: one name when the two are the same, both when they differ, and says a reply by anyone else is not an answer.
- Update the doc comments on `askOperator` and `awaitAnswer`, and `permissionPromptDescription` if it names the owner.

In `flai/internal/mcpserver/permission_test.go`, reproduce I-0081: a story owned by `arobson` in a project whose manifest owner is `alex`, and a reply `allow` by `alex` lets the write through and settles the thread as allowed by `alex`. Keep or add the cases that a reply by the story's owner still answers, a reply by another agent is not an answer, and the request names both names when they differ.

## Done when

- the I-0081 test fails on the code before the change and passes after it
- the existing permission tests pass, updated only where the request's wording changed
- `scripts/flai-test.sh` passes

## Notes
