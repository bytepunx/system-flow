---
id: T-3161
type: task
nature: feature
title: scripts/agent-image-digests.sh writes the agent images' digest list in the ADR's format, with a test in the smoke tier
status: backlog
parent: S-0734
owner: alex
created: 2026-11-12T15:57:36Z
updated: 2026-11-12T16:07:54Z
transitions: []
stream: S-0734
tags: [cli]
touches: [scripts/agent-image-digests.sh, scripts/agent-image-digests-test.sh, scripts/smoke.sh, scripts/README.md]
---
# T-3161 scripts/agent-image-digests.sh writes the agent images' digest list in the ADR's format, with a test in the smoke tier

## Work

Write `scripts/agent-image-digests.sh`, POSIX `sh` under `set -eu`, modelled on `scripts/flaiover-digests.sh`. Given the version and each image's repository and index digest, it writes the digest list in the format S-0733's ADR names. If the ADR extends flaiover's format with a kind for agent images, follow that, so `flai/internal/dashboard/digests.go`'s parser, or its successor, reads it.

Write `scripts/agent-image-digests-test.sh`, which runs it over fixed digests and compares its output with the expected list. Add the test to `scripts/smoke.sh` beside the other script tests, and both scripts to `scripts/README.md`.

It waits for no task.

## Done when

- `scripts/agent-image-digests-test.sh` passes, and fails when the output changes.
- ShellCheck passes on both scripts.
- `scripts/smoke.sh` runs the test, and `scripts/README.md` lists both.

## Notes
