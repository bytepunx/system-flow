---
id: T-1361
type: task
nature: remediation
title: flai host's watch restarts the dashboard from the image ID it recorded, not the tag
status: backlog
parent: S-0344
owner: alex
created: 2026-10-08T08:35:44Z
updated: 2026-10-08T08:35:44Z
transitions: []
stream: S-0344
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go]
after: [T-1360]
---
# T-1361 flai host's watch restarts the dashboard from the image ID it recorded, not the tag

## Work

Second layer: it waits for T-1360, because it takes the way T-1360 chose to run a container by image ID and still name its tag, and both change how `startDashboard` is called.

The watch has the same cause as I-0116. ADR-0118 §1 says a restart, `flai dashboard restart` or `flai host`'s watch, keeps the release that runs. Today `startDashboard` records the ref it was given, a tag after `flai dashboard` or `flai dashboard upgrade`, and `restartDashboard` starts that tag again, so a newer image a check pulled is started after a crash.

- Write the test first in `flai/cmd/dashboard_watch_test.go`: the dashboard starts from `latest` at one image ID, a check moves `latest`, the container goes; the watch's restart must run the recorded ID.
- In `flai/cmd/dashboard_watch.go`, have `startDashboard` record the image ID the container started from beside its ref (a new `dashboardRecord` field, read as absent from a record an older flai wrote), and `restartDashboard` start from that ID.
- When the recorded ID is no longer a local image, such as after `docker image prune`, start from the recorded ref and log it at warn, rather than leave the dashboard down.

## Done when

- The new test fails on the old watch and passes on the new one.
- A record an older flai wrote, with no image ID, still restarts from its ref.
- A recorded ID that is gone falls back to the ref with a warning, with a test.
- `flai test flai/cmd` passes.

## Notes

Drafted by the planner for S-0344. Included because criterion 1 asks that the cause no longer occur, and the watch restarts from a tag for the same reason.
