---
id: T-3160
type: task
nature: feature
title: The bedrock Dockerfile takes flai and pilot already built, both images carry the label that ties them to their flai, and every base is pinned by digest
status: done
parent: S-0734
owner: alex
created: 2026-11-12T15:57:32Z
updated: 2026-11-13T08:02:52Z
transitions:
  - to: ready
    at: 2026-11-13T07:56:24Z
    by: agent-S-0734
  - to: in-progress
    at: 2026-11-13T07:56:24Z
    by: agent-S-0734
  - to: done
    at: 2026-11-13T08:02:52Z
    by: agent-S-0734
stream: S-0734
tags: [cli]
touches: [flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/ci-pins.sh]
usage:
  source: log
  seconds: 388
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 6019
      cache_read: 1393433
      cache_write: 38404
      cost: 0.6706
---
# T-3160 The bedrock Dockerfile takes flai and pilot already built, both images carry the label that ties them to their flai, and every base is pinned by digest

## Work

In `flai/images/bedrock/Dockerfile`, add a way to take `flai` and `pilot` already built, such as a build argument or a named build context. The release workflow then copies the release's own binaries rather than building them again. Keep the stage that builds them from the tree, for `scripts/agent-image.sh` and a development tree.

Give both images the label S-0733's ADR names, from build arguments holding the version and commit. The bedrock image has `dev.system-flow.flai` at line 36 today; the claude-code image's label is at lines 10 to 11.

Pin every `FROM` by tag and index digest, as `flaiover/Dockerfile` is. Make `scripts/ci-pins.sh` (lines 78 to 99) fail an image Dockerfile whose `FROM` has no digest.

`flai/images` is a protected folder for a container's agent (`flai/internal/protected/protected.go`, lines 75 to 78), so the story's agent's writes there ask the owner. It waits for no task.

## Done when

- `scripts/agent-image.sh` still builds both images from the tree, and `scripts/agent-image.sh --check` passes.
- A build given binaries already built copies them, and `docker image inspect` shows the label with the version and commit given.
- `scripts/ci-pins.sh` passes, and fails on a test `FROM` with no digest.

## Notes
