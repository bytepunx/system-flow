---
id: T-3604
type: task
nature: feature
title: flai upgrade removes a retired template file a project has not changed, and reports without removing one with project additions
status: backlog
parent: S-0837
owner: alex
created: 2026-11-17T12:37:56Z
updated: 2026-11-17T12:38:40Z
transitions: []
stream: S-0837
tags: [flai]
touches: [flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, flai/cmd/upgrade.go, flai/cmd/upgrade_test.go]
after: [T-3600]
---
# T-3604 flai upgrade removes a retired template file a project has not changed, and reports without removing one with project additions

## Work

This task waits for T-3600, because it builds ADR-0316 §5's removal as accepted.

- `Compute` in `flai/internal/upgrade/upgrade.go` classifies only the paths the template renders, so a path the template no longer ships is never removed. Add a list of retired template paths, starting with the publish formula. Derive each path from the manifest's design layout, as `kept()` does.
- For each retired path the project holds: when its content matches the hash `system-flow.lock.yaml` recorded, mark it for removal. When the project changed it, as with text below the marker, keep it and report it. A path the project does not hold yields nothing.
- `Describe` and `flai/cmd/upgrade.go` print each removal, and each kept file with why it was kept. `Apply` deletes the removed files and drops them from the new lock.
- Test both cases and the absent case in `flai/internal/upgrade/upgrade_test.go`, on the model of the kept-file test. Test the command's output in `flai/cmd/upgrade_test.go`.

## Done when

- `flai upgrade` removes an unchanged retired file and drops it from the lock.
- It reports, without removing, a changed one, and the tests show both.
- `flai test` passes on the changed paths.

## Notes

- Drafted by planner-S-0837.
- T-3605 deletes the template's copy once this lands, so an upgrade over this template still knows the path.
