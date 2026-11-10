---
id: T-3054
type: task
nature: improvement
title: flai serve agent harness takes --outside folders, checks and stores them with the harness's settings, and shows them
status: backlog
parent: S-0701
owner: alex
created: 2026-11-10T10:35:09Z
updated: 2026-11-10T10:35:09Z
transitions: []
stream: S-0701
tags: [serve, config]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_agent.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
after: [T-3048]
---
# T-3054 flai serve agent harness takes --outside folders, checks and stores them with the harness's settings, and shows them

## Work

Add `outside`, a list of folders, to a harness's host settings in `flai/internal/config/config.go`, beside `network` and `login`. It is kept for a project under `agent.projects` and for every project in `agent.harnesses` with `--all-projects`. `flai serve agent harness <name> --outside <dir>` takes the flag again for each folder, `--outside ""` clears the list, and `--reset` clears it too. Each folder is an existing folder, given as an absolute path or under `~/`, and stored absolute. The validation refuses home itself, `~/.flai`, and any folder that holds the project root or that the project holds, naming the folder and the flag that fixes it. `flai serve agent harness` and `flai serve settings` show the list.

Add the setting's row to `docs/operators/settings.md`, which `TestSettingsIndex` needs. Regenerate `docs/users/flai-reference.md` with `make flai-reference`, which `TestReferenceIsCurrent` needs.

This task waits for T-3048, whose ADR decides the setting's name and its rules.

## Done when

- Tests in `flai/internal/config/config_test.go` cover a valid list, a relative path, a missing folder, home, and a folder that holds the project.
- A test in `flai/cmd/serve_actions_test.go` sets, shows, and clears the list.
- `flai test` passes on the changed paths.

## Notes
