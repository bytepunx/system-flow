---
id: T-0432
type: task
nature: feature
title: Publish the flaiover image for linux/amd64 and linux/arm64
status: done
parent: S-0119
owner: arobson
created: 2026-09-26T03:17:08Z
updated: 2026-09-26T03:20:38Z
transitions:
  - to: ready
    at: 2026-09-26T03:17:26Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T03:17:26Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T03:20:38Z
    by: agent-S-0119
stream: S-0119
tags: []
---

# T-0432 Publish the flaiover image for linux/amd64 and linux/arm64

## Work

- `.github/workflows/release-flaiover.yml` sets up QEMU and builds and pushes `linux/amd64,linux/arm64` as one multi-platform manifest under the same tags.
- `flaiover/Dockerfile` runs the build stage on `$BUILDPLATFORM`, so only the runtime stage runs emulated. The production `node_modules` it copies are pure JavaScript.

## Done when

- A local `docker buildx build --platform linux/amd64,linux/arm64` of the Dockerfile succeeds, and the arm64 image starts and answers `/_health`.
- The workflow names both platforms.

## Notes
