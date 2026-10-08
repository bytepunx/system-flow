---
id: T-1360
type: task
nature: remediation
title: flai dashboard restart starts the image the container runs, by its image ID, never the tag again
status: backlog
parent: S-0344
owner: alex
created: 2026-10-08T08:35:35Z
updated: 2026-10-08T08:35:35Z
transitions: []
stream: S-0344
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go, docs/users/flai-reference.md]
---
# T-1360 flai dashboard restart starts the image the container runs, by its image ID, never the tag again

## Work

First layer: it waits for no task. Record in the narrative's `## Decisions` the fix chosen, since the story asks for one proposed from I-0116's instance before it is built.

- Write the test first in `flai/cmd/dashboard_test.go`: the fake docker runs `flaiover` from `ghcr.io/bytepunx/flaiover:latest` at one image ID, then `flai dashboard check` moves `latest` to a newer ID; `flai dashboard restart` must run the old ID and pull nothing. It fails on the current code, which restarts from `{{.Config.Image}}`.
- In `runDashboardRestart` (`flai/cmd/dashboard_upgrade.go`), start from `containerImageID` (`{{.Image}}`), as S-0316's upgrade fallback does, when the container runs; keep the configured ref when it does not run.
- Keep what `flai dashboard status`, `flai dashboard versions`, and restart's own output name as the image: the tag, not a bare `sha256:` ID. Recommended: `containerArgs` gives the container a label with the ref it was started for, and `containerInfo` (`flai/cmd/dashboard.go`) prefers that label over `{{.Config.Image}}`. Another way that keeps the tag shown is fine; say why in `## Decisions`.
- Keep restart's `Long` text true, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- The new test fails on the old restart and passes on the new one.
- `flai dashboard restart` runs the image ID the container ran, and its output and `flai dashboard status` still name the tag.
- `flai test flai/cmd` passes and `docs/users/flai-reference.md` matches the command help.

## Notes

Drafted by the planner for S-0344. The cause is read from the code (I-0116's instance): `containerInfo` returns `{{.Config.Image}}`, the tag, and `flai dashboard check` pulls that tag.
