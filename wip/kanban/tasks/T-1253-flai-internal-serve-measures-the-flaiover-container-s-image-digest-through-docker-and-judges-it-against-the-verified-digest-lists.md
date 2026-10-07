---
id: T-1253
type: task
nature: feature
title: flai/internal/serve measures the flaiover container's image digest through Docker and judges it against the verified digest lists
status: backlog
parent: S-0238
owner: alex
created: 2026-10-07T23:05:33Z
updated: 2026-10-07T23:05:33Z
transitions: []
stream: S-0238
tags: [cli]
touches: [flai/internal/serve/imagecheck.go, flai/internal/serve/imagecheck_test.go]
---
# T-1253 flai/internal/serve measures the flaiover container's image digest through Docker and judges it against the verified digest lists

## Work

Add the image check to `flai/internal/serve/imagecheck.go`, a new file. Nothing calls it yet. It waits for no task: it reads only S-0236's `flai/internal/dashboard` package, which is on main before this story starts.

- Through an `execx.Runner`, ask Docker about the container `flaiover`. Use `docker inspect flaiover` for its image ID and its `FLAIOVER_VERSION`. Use `docker image inspect <id>` for its `RepoDigests`.
- Compare each digest in `RepoDigests` with the verified lists cached beside `flai serve`'s state. Read them through S-0236's cache.
- When no list is cached for the running version, fetch and verify that version's list through S-0236's resolver, then compare.
- Return one verdict:
  - listed: the image and the digest that matched.
  - unlisted: the image and its digests. An image with no `RepoDigests`, such as a local build, is unlisted.
  - not measured: Docker is not on `PATH`, no container `flaiover` runs here, or the dashboard's URL is not on this host. The verdict carries the reason, and the caller falls back to the stamp check alone.
- Also return the container's ID and image ID, so a caller can tell when the container changed.
- Record in a comment which digest `RepoDigests` holds after a pull by digest, the index digest or the platform's. `release-signing.md` leaves this to check.

## Done when

- `imagecheck_test.go` covers each case with a fake runner: a listed digest, an unlisted one, a local image with no `RepoDigests`, a missing list fetched and verified, a list that does not verify, and the not-measured cases (no Docker, no container, a dashboard on another host).
- `flai test flai/internal/serve` passes.

## Notes
