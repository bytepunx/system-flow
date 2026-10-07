---
id: T-1261
type: task
nature: feature
title: "flai dashboard --allow-unsigned, and the setting or --build, start the container with FLAIOVER_ALLOW_UNSIGNED=1"
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:17Z
updated: 2026-10-07T23:09:17Z
transitions: []
stream: S-0239
tags: [cli]
touches: [flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1259]
---
# T-1261 flai dashboard --allow-unsigned, and the setting or --build, start the container with FLAIOVER_ALLOW_UNSIGNED=1

## Work

The second half of criterion 1 and flai's container half of criterion 2. It waits for T-1259, which adds the setting it reads.

- `flai/cmd/dashboard.go`: `--allow-unsigned` on `flai dashboard`, which turns the allowance on for this one start without writing the configuration. When the allowance is in force for the project (T-1259's function), or the flag or `--build` is given, the container is started with `--env FLAIOVER_ALLOW_UNSIGNED=1`; otherwise the variable is not passed at all.
- `flai/cmd/dashboard_upgrade.go`: the `docker run` arguments are built here; `flai dashboard upgrade` and `restart` pass the variable on the same terms, so that a restart keeps what the start was given. Check whether `flai host`'s restart (ADR-0062) reuses the recorded arguments; if it rebuilds them in `dashboard_watch.go`, widen the touches to it.
- Regenerate the reference with `make flai-reference`, which rewrites `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`; `flai/cmd/reference_test.go` fails without it.
- Tests in `dashboard_test.go`, with the fake runner: the variable is passed with the setting on, with `--allow-unsigned`, and with `--build`, and is absent with none of them; `--allow-unsigned` leaves the configuration unchanged.

## Done when

- `flai dashboard --allow-unsigned`, the setting, and `--build` each start the container with `FLAIOVER_ALLOW_UNSIGNED=1`, and none of them leaves it out.
- The generated reference and flag index are current, and the tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner.
