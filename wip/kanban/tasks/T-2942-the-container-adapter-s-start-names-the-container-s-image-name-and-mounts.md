---
id: T-2942
type: task
nature: improvement
title: The container adapter's Start names the container's image, name, and mounts
status: backlog
parent: S-0665
owner: alex
created: 2026-11-10T08:50:21Z
updated: 2026-11-10T08:50:21Z
transitions: []
stream: S-0665
tags: [containers, harness]
touches: [flai/internal/harness/harness.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go]
---
# T-2942 The container adapter's Start names the container's image, name, and mounts

## Work

Give `harness.Start` (`flai/internal/harness/harness.go`) a field for a container start: the image the container runs, its name (`ContainerName`, `flai-agent-<subject>-<session>`), and its mounts, each as host path, container path, and whether it is read-only. Fill it in `flai/internal/harness/container.go` where the argv is built, so the field and the argv cannot disagree: `Start` for a story's container (the worktree, the git and main checkout mounts, and the cache), `roleStart` for a strategic agent's (the checkout, its hooks read-only, and the cache), and `containerStart` for the login mount and the image. Every other adapter leaves the field empty. The field carries no secret: the login mount names its path only, as the argv does.

This task waits for nothing. It is the first layer, because flai serve's record in the next task reads this field.

## Done when

- A story's and a strategic agent's container `Start` name the image, the container's name, and every mount the argv gives `-v`, in argv order, each with its read-only mark.
- Tests in `flai/internal/harness/container_test.go` check the field against the argv for a story's container and a strategic agent's, with and without a login and a layer image.
- `flai test` passes on the changed paths.

## Notes
