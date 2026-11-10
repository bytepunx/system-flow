---
id: T-2975
type: task
nature: improvement
title: pilot loads a defined role's file by its name from the mounted checkout and starts the inner harness as that role
status: backlog
parent: S-0705
owner: alex
created: 2026-11-10T10:29:04Z
updated: 2026-11-10T10:29:04Z
transitions: []
stream: S-0705
tags: [flai, harness]
touches: [flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go]
after: [T-2973]
---
# T-2975 pilot loads a defined role's file by its name from the mounted checkout and starts the inner harness as that role

## Work

Second layer. It waits for T-2973, whose start request field (`role_name`, and the inputs) pilot reads here.

- In `flai/internal/pilot/run.go`, beside `builtin`, give a request that names a defined role its file: read the manifest at the request's root for the design folder, load the role files under it (`roles.Load`), and set `RoleFile` to the one whose name is `RoleName`. Refuse, with one line saying what to fix, a name no file has, a file with a problem `roles.Load` reports, and a file whose verb is not the request's `Role`. Run it in `prepare` before the adapter's `Start`, so that the inner harness's definition, `--agent <name>`, `FLAI_ROLE`, `FLAI_ITEM`, and the start prompt with the inputs are what they are on the host.
- In `flai/internal/pilot/request.go`, `Check` refuses a `role_name` with no `Role`, with a story, or on a probe.
- Tests in `flai/internal/pilot/run_test.go` and `request_test.go`: a checkout with a role file of its own under `design/workflows/roles/` and a start request naming it starts the fake inner harness with that role's agent, verb, item, and a prompt that holds the inputs; each refusal above names its fix.

## Done when

- pilot starts a defined role from its start request and the mounted checkout's role file, and refuses each case above with its fix.
- `flai test` on the changed paths passes.

## Notes
