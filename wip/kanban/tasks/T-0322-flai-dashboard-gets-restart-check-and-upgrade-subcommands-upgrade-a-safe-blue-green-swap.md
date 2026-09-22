---
id: T-0322
type: task
nature: feature
title: flai dashboard gets restart, check, and upgrade subcommands, upgrade a safe blue-green swap
status: done
parent: S-0081
owner: alex
created: 2026-09-22T21:12:33Z
updated: 2026-09-22T21:22:27Z
transitions:
  - to: ready
    at: 2026-09-22T21:13:12Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T21:13:12Z
    by: system-flow
  - to: review
    at: 2026-09-22T21:22:27Z
    by: system-flow
  - to: done
    at: 2026-09-22T21:22:27Z
    by: system-flow
stream: S-0081
tags: []
---
# T-0322 flai dashboard gets restart, check, and upgrade subcommands, upgrade a safe blue-green swap

## Work
Add `flai dashboard restart`, `flai dashboard check`, and `flai dashboard upgrade [--tag TAG]` in `flai/cmd/dashboard.go`. Refactor the container-args building out of `runDashboard` into a shared helper so all three paths (initial start, restart, upgrade) mount the same two secrets the same way.

- `restart`: if the container is not running, starts it (same as bare `flai dashboard`, no pull). If running, captures the image it is actually running (`docker inspect <name> --format '{{.Image}}'` by way of `.Config.Image`, i.e. the ref, not a fresh pull), stops it, starts it again with that same ref. Registrations in `flai serve`'s state are untouched either way — restart never unregisters a project.
- `check`: pulls the configured (or `--tag`-overridden) ref, compares the freshly pulled image's ID (`docker image inspect <ref> --format '{{.Id}}'`) to the running container's image ID (`docker inspect <name> --format '{{.Image}}'`); reports whether an upgrade is available, with both refs/tags. Changes nothing. Not running is reported plainly, not as an error.
- `upgrade`: pulls the target ref; if not running, starts it directly (nothing to swap). If running and the pulled image's ID matches the running one, reports "already up to date" and does nothing else. If they differ: starts a temporary container from the new image (`--rm`, same two secrets, published on an ephemeral loopback port via `-p 127.0.0.1::<containerPort>`, discovered with `docker port`), polls its `/_health` over plain HTTP for up to ~10s; on success stops the temporary container, then stops the running one and starts the final container with the new image at the real name and port; on failure or timeout, stops the temporary container and leaves the running one untouched, returning a failure that says why. The running container is never stopped until the replacement has proven healthy.

Verified by hand against real Docker on a scratch image (not `flaiover:s80final` or any name touching the operator's dashboard) before writing the unit tests, so the tests encode confirmed behavior of `docker pull --quiet`, `docker inspect`, `docker port`, and `docker image inspect --format '{{.Id}}'` rather than assumptions about their output.

## Done when
`--json` output on all three names the container, the before/after image, and, for `check`/`upgrade`, whether an upgrade is or was available. Unit tests (extending `cmd/dashboard_test.go`'s fake runner) cover: restart while running and while stopped; check with and without an available upgrade; upgrade when already current; upgrade that swaps successfully; upgrade whose temporary container never becomes healthy, asserting the real container was never stopped. `go test -race -short ./cmd/...` and `golangci-lint run ./...` clean.

## Notes
No flaiover or hostapi wiring in this task — CLI and its tests only. The blue-green swap is the load-bearing design decision of this story: it is the only way "an upgrade that fails leaves the previous container running" (S-0081's third criterion) can be true rather than best-effort.
