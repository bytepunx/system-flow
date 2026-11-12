---
id: T-3169
type: task
nature: feature
title: The agent images are pulled by the digest their release's signed list names, verified as flaiover's list is
status: backlog
parent: S-0740
owner: alex
created: 2026-11-12T15:59:26Z
updated: 2026-11-12T15:59:26Z
transitions: []
stream: S-0740
tags: [cli]
touches: [flai/internal/dashboard/digests.go, flai/internal/dashboard/digests_test.go, flai/internal/dashboard/resolve.go, flai/internal/agentimage/pull.go, flai/internal/agentimage/pull_test.go]
---
# T-3169 The agent images are pulled by the digest their release's signed list names, verified as flaiover's list is

## Work

Let the digest-list code read the agent images' list in the format S-0733's ADR names and S-0734's `scripts/agent-image-digests.sh` writes:

- `Parse` and `Verify` in `flai/internal/dashboard/digests.go` (lines 92 to 125);
- `Resolve` and `FetchList` in `resolve.go` (lines 94 to 182), with the `flai/vX.Y.Z` release in place of `flaiover/vX.Y.Z`.

Keep flaiover's behaviour and refusal kinds as they are.

`flai/internal/agentimage/pull.go` resolves the running flai's version (`buildinfo`), fetches and verifies the list, and pulls each image by `<repository>@<digest>` with `docker pull`. It retries a network failure up to the ADR's count, then tags the image with the name the start runs. Tests run over a fake release client and docker, covering:

- a good list;
- a bad signature;
- an unknown key;
- a digest not on the list;
- a network failure retried, then given up.

It waits for no task.

## Done when

- `go test ./internal/dashboard/ ./internal/agentimage/` passes, and the dashboard's tests are unchanged.
- A pull refused for each refusal kind names the kind and what to do.

## Notes
