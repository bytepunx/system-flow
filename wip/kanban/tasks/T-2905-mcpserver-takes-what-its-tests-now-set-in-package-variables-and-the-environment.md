---
id: T-2905
type: task
nature: improvement
title: mcpserver takes what its tests now set in package variables and the environment, the permission poll, the flai it runs, its version, the releases API and token, and a child's config, through Options
status: in-progress
parent: S-0684
owner: alex
created: 2026-11-10T08:47:16Z
updated: 2026-11-10T09:07:31Z
transitions:
  - to: ready
    at: 2026-11-10T09:07:30Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:07:31Z
    by: system-flow
stream: S-0684
tags: [flai, go, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/subproject.go, flai/internal/mcpserver/versions.go]
after: [T-2904]
---
# T-2905 mcpserver takes what its tests now set in package variables and the environment, the permission poll, the flai it runs, its version, the releases API and token, and a child's config, through Options

## Work

This task waits for the first task because both change `Options` in `server.go`. Besides the session, the serial tests set things the server reads from the process. Give each a field of `Options` with today's value as its default, so a test can pass it instead:

- `permissionPoll` (`permission.go:69`), which `permission_test.go` swaps for a fast poll.
- `flaiExecutable` (`subproject.go:48`), which `subproject_test.go` swaps for a built flai.
- The `FLAI_CONFIG` that `subproject_test.go` sets for the flai that `subproject_edit` runs. That child's environment is built at `subproject.go:152`, so the field adds to it rather than to the process.
- `buildinfo.Version`, which `outdated_test.go` sets. Take the version from `Options.Version`, or `s.version`, wherever the code under test reads `buildinfo.Version` through the server.
- `FLAI_RELEASES_API` and the token that `versions.go:55` reads. `versions_test.go` sets `FLAI_RELEASES_API`, `GH_TOKEN`, and `GITHUB_TOKEN`.

Change no behaviour: a zero field keeps what the server does today. `perf.SlowEnv`, which `server_test.go` sets, needs no new field, since `Options.Slow` exists already.

## Done when

- [ ] No test of the package needs to assign `permissionPoll`, `flaiExecutable`, or `buildinfo.Version`, or to set `FLAI_CONFIG`, `FLAI_RELEASES_API`, `GH_TOKEN`, or `GITHUB_TOKEN`, to get what it tests: an `Options` field gives it each.
- [ ] With every new field zero, the server reads what it read before.
- [ ] `flai test` passes on `flai/internal/mcpserver` and `flai/cmd`.

## Notes
