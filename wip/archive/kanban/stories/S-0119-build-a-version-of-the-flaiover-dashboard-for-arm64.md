---
id: S-0119
type: story
nature: feature
title: Build a version of the flaiover dashboard for arm64
status: done
owner: alex
created: 2026-09-26T03:14:08Z
updated: 2026-09-26T03:23:38Z
transitions:
  - to: ready
    at: 2026-09-26T03:15:49Z
    by: alex
  - to: in-progress
    at: 2026-09-26T03:16:18Z
    by: agent-S-0119
  - to: review
    at: 2026-09-26T03:23:13Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T03:23:38Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/Dockerfile, ".github/workflows/release-flaiover.yml", flai/cmd, design/system/flaiover-dashboard.md, docs]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0119 Build a version of the flaiover dashboard for arm64

## Goal

CI does not presently build an arm64 compatible docker image. This forces people to build their own image locally or do without and breaks the dashboard --pull command.

## Acceptance criteria
- [x] CI builds include an arm64 image build
- [x] flai dashboard --pull fetches the correct image from ghcr.io based on the platform

## Tasks
- T-0432 Publish the flaiover image for linux/amd64 and linux/arm64
- T-0433 flai dashboard pulls the image for the Docker daemon's platform
- T-0434 Document the multi-platform image and the platform-aware pull

## Notes

- How the criteria were verified. `release-flaiover.yml` now builds `linux/amd64,linux/arm64` with QEMU and pushes one manifest. That workflow runs only on main and on tags, so it has not run yet; its first run is the push after acceptance. The same Dockerfile was built locally with `docker buildx` for both platforms on an arm64 host, and both images answered `/_health` 200 (amd64 under emulation).
- `flai dashboard` pulls with `--platform` set to the daemon's `os/arch`. It pulls again when the present image is for another platform, and names the platform when a tag has none. `TestDashboardPullsTheDaemonsPlatform` covers this with the fake runner. Real Docker on this host answered `docker version` with `linux/arm64`, accepted `docker pull --platform`, and gave the `no matching manifest` text the error check matches.
- Tags published before this story stay amd64 only. `flai dashboard --pull` on arm64 needs a tag published after the first run of the new workflow (`latest` is republished on every push to main).
