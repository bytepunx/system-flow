---
id: T-1634
type: task
nature: feature
title: flai serve merges a story's pull request on the forge when asked, by the project's method, deletes the remote story branch, and answers whether it merged or why not
status: done
parent: S-0397
owner: alex
created: 2026-10-16T04:39:59Z
updated: 2026-10-18T02:41:51Z
transitions:
  - to: ready
    at: 2026-10-18T02:28:30Z
    by: agent-S-0397
  - to: in-progress
    at: 2026-10-18T02:28:31Z
    by: agent-S-0397
  - to: done
    at: 2026-10-18T02:41:51Z
    by: agent-S-0397
stream: S-0397
tags: [cli]
touches: [flai/internal/serve/forgemerge.go, flai/internal/serve/forgemerge_test.go, flai/internal/forge/forge.go, flai/internal/forge/github.go, flai/internal/forge/github_test.go, flai/internal/forge/state.go, flai/internal/serve/forge.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forge_test.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/serve.go]
usage:
  source: log
  seconds: 800
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 129
      output: 755
      cache_read: 8598519
      cache_write: 210425
      cost: 3.8289
---
# T-1634 flai serve merges a story's pull request on the forge when asked, by the project's method, deletes the remote story branch, and answers whether it merged or why not

## Work

- **Contract.** Add `PullRequest(ctx, number int) (PullRequest, error)` to `forge.Forge`: the pull request by number, in any state, `ErrNotFound` when absent. In the GitHub adapter it is `GET /repos/{o}/{r}/pulls/{n}`. Make `Merge`'s `ErrMergeRefused` error carry the forge's own message. Update every `Forge` implementer in the tests.
- **Merged elsewhere.** `forge.Status` gains `State`, the pull request's state as last read. When no open pull request is found and the kept number is known, the read asks `PullRequest(number)` and keeps `merged` or `closed`. `PushStory` neither pushes again nor opens a new pull request for a story whose pull request is merged.
- **Merge.** In `flai/internal/serve/forgemerge.go`, `MergeStory(ctx, ForgeMerge) MergeAnswer`, for one story under pull-request acceptance:
  - Find its pull request by the kept number, else by its head branch.
  - Merged already: answer `already-merged`.
  - Closed without a merge: refuse.
  - Open with a head other than the one the gate passed at: refuse, naming both.
  - Otherwise `Merge` by `acceptance.merge_method`, or the method given. A refusal is answered with the forge's reason.
  - After a merge, or one found made, delete `story/S-nnnn` on `origin` with `git push origin --delete`, when origin still has it. A failed delete is reported in the answer and does not fail it.
  - The answer: story, pull request number and URL, method, outcome (`merged`, `already-merged`, `refused`), merge commit, reason, branch deleted or why not, and when.
- **Hand-off.** The request is `.flai-cache/forge/merge/S-nnnn.request.json`: story, head, method, by.
  - `RequestMerge` writes it atomically and waits 30s for flai serve to take it, by renaming it to `S-nnnn.taken.json`. When the request is still there, it removes it and refuses, saying to start flai serve with the forge host action on. Otherwise it waits 5m for `S-nnnn.answer.json`, reads it, and cleans up.
  - flai serve's merge pass, started per project beside the forge pass, polls the folder every second. It takes each request and answers it: refused, naming `flai serve enable forge`, while the host action is off; otherwise with `MergeStory` and serve's token.
  - Waits and intervals are options, so the tests run fast. Never log or answer the token.
- **Tests.** Against the fake forge and a scratch bare origin: a merge by each method; already merged; refused for a conflict and for a rule, with the forge's reason; a moved head; the branch deleted; the hand-off round trip; the host action off; no flai serve taking the request; a merged pull request neither pushed nor reopened.
- **Layer.** First, beside T-1633: it waits for no task.

## Done when

`MergeStory` and the hand-off answer each case above. A pull request someone else merged is kept as merged and never reopened. The forge and serve tests pass.

## Notes
