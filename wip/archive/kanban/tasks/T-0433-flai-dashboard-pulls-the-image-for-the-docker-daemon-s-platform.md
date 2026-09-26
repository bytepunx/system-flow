---
id: T-0433
type: task
nature: feature
title: flai dashboard pulls the image for the Docker daemon's platform
status: done
parent: S-0119
owner: arobson
created: 2026-09-26T03:17:08Z
updated: 2026-09-26T03:21:40Z
transitions:
  - to: ready
    at: 2026-09-26T03:17:26Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T03:20:38Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T03:21:40Z
    by: agent-S-0119
stream: S-0119
tags: []
---

# T-0433 flai dashboard pulls the image for the Docker daemon's platform

## Work

- `flai/cmd/dashboard.go` asks the daemon for its platform (`docker version --format '{{.Server.Os}}/{{.Server.Arch}}'`) and passes it to `docker pull --platform`. It does not use flai's own GOARCH, because an amd64 flai under Rosetta talks to an arm64 daemon.
- A present image built for another platform is pulled again instead of reused.
- A registry with no image for the platform gets an error that names the platform and suggests `--build`.

## Done when

- Behavior tests cover: the pull names the daemon's platform; a present image of the wrong platform is pulled; a present image of the right platform is not; the missing-manifest error.
- `make test` and lint pass.

## Notes
