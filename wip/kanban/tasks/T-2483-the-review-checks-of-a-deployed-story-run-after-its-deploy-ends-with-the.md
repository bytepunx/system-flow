---
id: T-2483
type: task
nature: feature
title: The review checks of a deployed story run after its deploy ends, with the namespace and the kube context in their environment
status: backlog
parent: S-0579
owner: alex
created: 2026-11-06T21:44:28Z
updated: 2026-11-06T21:44:28Z
transitions: []
stream: S-0579
tags: [kluster, serve, checks]
touches: [flai/internal/serve/checks.go, flai/internal/serve/checks_test.go]
after: [T-2481]
---
# T-2483 The review checks of a deployed story run after its deploy ends, with the namespace and the kube context in their environment

## Work

It waits for T-2481, whose deploys and their records the checks wait for and read.

In `flai/internal/serve/checks.go`, while the `kluster` host action is on for the project and the block is set: hold a story's review checks until its deploy at the branch's head has ended. When it passed, run the checks with the variables S-0576's ADR names for them (such as the namespace and the kube context) in their environment. When it failed, run them as before and record in the run that no deployment was made. With the action off, or no block, run the checks as before, with nothing added.

Test in `checks_test.go`: checks held until a deploy ends, the variables present after a passing deploy, absent after a failing one, and absent with the action off.

## Done when

- A deployed story's checks see the namespace and context, and wait for the deploy.
- Every case above is tested, and `flai test` on the two paths passes.

## Notes
