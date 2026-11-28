---
id: T-3658
type: task
nature: remediation
title: flai kluster deploy --job, teardown --job, and flai checks run --job ask through the story's agent cache folder in its container, with a test that reproduces I-0405
status: backlog
parent: S-0849
owner: alex
created: 2026-11-28T00:05:09Z
updated: 2026-11-28T00:05:09Z
transitions: []
stream: S-0849
tags: [flai, kluster, cli]
touches: [flai/cmd/kluster.go, flai/cmd/checks.go, flai/cmd/kluster_container_test.go, docs/users/flai-reference.md]
after: [T-3656]
---
# T-3658 flai kluster deploy --job, teardown --job, and flai checks run --job ask through the story's agent cache folder in its container, with a test that reproduces I-0405

## Work

In `flai/cmd/kluster.go`, make `askKlusterJob` choose the requests folder it asks through.

- First write the test that reproduces I-0405 in `flai/cmd/kluster_container_test.go`, as that file's ADR-0319 tests set up a container: `FLAI_STORY` set, the configuration in the story's agent cache folder, and no `.flai-cache/kluster` in the main checkout. Run `flai kluster deploy S-nnnn --job` there with a goroutine that answers only through the story's agent folder. Without the fix, the request goes under `.flai-cache/kluster/requests` and the command ends with flai serve not running.
- Choose the folder as ADR-0319's readers choose the root. Use the main checkout's `requests/` when `.flai-cache/kluster` is there, as on the host. Otherwise use `kluster.AgentRequestsDir(root, story)`. Ask through that folder, and read the job state from it for `--wait`, `HasAskedJob`, and `foundState`.
- `flai checks run --job` in `flai/cmd/checks.go` goes through `askKlusterJob` already. Change its help, and `klusterJobHelp` in `kluster.go`, so they no longer name only `.flai-cache/kluster/requests`: in a story's container the request goes in the story's cache folder.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Add tests for `teardown --job` and `checks run --job` in the container, and for `deploy --job --wait` reading the job's end from the agent folder.

It waits for T-3656, whose `AgentRequestsDir` and folder-based `Ask` it calls. It shares no path with T-3657, so the two run together. This test answers through the folder itself, so it does not need T-3657's pass.

## Done when

- The reproducing test fails before the change and passes after it.
- The teardown, checks, and `--wait` tests in the container pass. The existing host tests of `--job` pass unchanged.
- The help of the three commands says where a container's request goes, and `docs/users/flai-reference.md` is current (`TestReferenceIsCurrent`).
- `flai test` on the four paths passes.

## Notes
