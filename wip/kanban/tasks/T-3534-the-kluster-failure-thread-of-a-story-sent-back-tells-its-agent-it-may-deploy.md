---
id: T-3534
type: task
nature: improvement
title: The kluster failure thread of a story sent back tells its agent it may deploy the story before review
status: backlog
parent: S-0831
owner: alex
created: 2026-11-17T08:29:48Z
updated: 2026-11-17T08:29:48Z
transitions: []
stream: S-0831
tags: [flai, serve, kluster]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go]
---
# T-3534 The kluster failure thread of a story sent back tells its agent it may deploy the story before review

## Work

Rewrite the sent-back branch of `heldText` in `flai/internal/serve/kluster.go`, which the deploy failure and the checks failure threads end with. Keep what it says of the kept namespace, of a cause outside the story, and of cancelling or accepting. In place of "Once S-nnnn is back in review … flai serve deploys it again in place", say that the agent may ask for its own deploy at the branch's head with `kluster_deploy` or `flai kluster deploy <story> --job`, which runs the checks once the deploy passes while `checks` is on. Say that a pass at the head resolves this thread, and that back in review flai serve deploys again in place only when that has not passed at the head ([kluster.md](../../../design/system/kluster.md) § The send-back, ADR-0292 §3). Leave the branch for a story not sent back as it is: that story stays in review.

In `flai/internal/serve/kluster_test.go`, make `TestAFailedDeployOpensOneThreadSendsTheStoryBackAndAPassResolvesIt` and `TestFailingChecksAgainstTheNamespaceSendTheStoryBack` check that the thread's entry names `kluster_deploy` and `flai kluster deploy <story> --job`. That reproduces I-0357 on the thread.

It waits for no task: it shares no path with the start prompt.

## Done when

- The deploy and checks failure threads of a story sent back name `kluster_deploy` and `flai kluster deploy <story> --job`, say a pass at the head resolves the thread, and say the review deploys again only what has not passed at the head.
- The two tests fail on the old text and pass on the new one, and `flai test` passes on both files.

## Notes
