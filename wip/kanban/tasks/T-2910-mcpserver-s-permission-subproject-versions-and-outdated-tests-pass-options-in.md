---
id: T-2910
type: task
nature: improvement
title: mcpserver's permission, subproject, versions, and outdated tests pass Options in place of package variables and run in parallel, and together goes
status: done
parent: S-0684
owner: alex
created: 2026-11-10T08:47:46Z
updated: 2026-11-10T09:42:22Z
transitions:
  - to: ready
    at: 2026-11-10T09:32:23Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:32:24Z
    by: system-flow
  - to: done
    at: 2026-11-10T09:42:22Z
    by: agent-S-0684
stream: S-0684
tags: [flai, go, mcp, testing]
touches: [flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/subproject_test.go, flai/internal/mcpserver/versions_test.go, flai/internal/mcpserver/outdated_test.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/subproject.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go]
after: [T-2905, T-2906]
usage:
  source: log
  seconds: 598
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 68
      output: 509
      cache_read: 11830639
      cache_write: 438462
      cost: 4.5164
---
# T-2910 mcpserver's permission, subproject, versions, and outdated tests pass Options in place of package variables and run in parallel, and together goes

## Work

This task waits for the second task, whose `Options` fields it passes. It also waits for the third task, which owns `server_test.go` until it is done and sets git's identity once for the package.

- `permission_test.go`: pass the fast poll through `Options` in place of assigning `permissionPoll`. Make its `together` subtests top-level parallel tests again.
- `subproject_test.go`: pass the built flai and the child's config through `Options` in place of `flaiExecutable` and `t.Setenv("FLAI_CONFIG")`. Pass `StartedBy` in place of `t.Setenv(guard.StartedByEnv)`. Rely on the package's git identity in place of `gittest.Identity`.
- `versions_test.go`: pass the releases API and the token through `Options` in place of `t.Setenv` of `FLAI_RELEASES_API`, `GH_TOKEN`, and `GITHUB_TOKEN`.
- `outdated_test.go`: pass the version through `Options.Version` in place of assigning `buildinfo.Version`.
- `server_test.go`: remove `together` and `inEnv` once nothing calls them. Remove `sessionEnv` too, unless `TestMain` still needs it to keep the session of the agent running the tests away from `SessionFromEnv`'s test.

## Done when

- [ ] No test file of `flai/internal/mcpserver` calls `t.Setenv` or `together`, or assigns a package variable, but `session_test.go`'s one test of `SessionFromEnv` and `TestMain`.
- [ ] Every top-level test in the files this task touches calls `t.Parallel()`.
- [ ] `go test -race -count=1 ./internal/mcpserver` passes from `flai/`, and so does `flai test` on the package.

## Notes
