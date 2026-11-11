---
id: S-0721
type: story
nature: improvement
title: "Every MCP server flai builds resolves all its tools' JSON schemas again, which is 56% of the mcpserver tests' CPU: share one SchemaCache"
status: done
owner: alex
created: 2026-11-11T05:24:30Z
updated: 2026-11-11T05:47:58Z
transitions:
  - to: ready
    at: 2026-11-11T05:24:47Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T05:27:10Z
    by: agent-S-0721
  - to: review
    at: 2026-11-11T05:45:25Z
    by: agent-S-0721
  - to: done
    at: 2026-11-11T05:47:58Z
    by: alex
tags: [mcp, tests]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, design/tech/go-libraries.md, flai/internal/mcpserver/schema_cache_test.go, wip/agents/S-0721.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1108
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 3
      test_runs: 2
      work: 33
  models:
    - model: claude-opus-5-5
      input: 106
      output: 21852
      cache_read: 12571761
      cache_write: 319699
      cost: 5.419
  strategic:
    - kind: planner
      seconds: 249
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 58
          output: 17810
          cache_read: 1954424
          cache_write: 92618
          cost: 1.4883
    - kind: orchestrator
      seconds: 137
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 44
          output: 359
          cache_read: 1526591
          cache_write: 236759
          cost: 0.4893
        - model: claude-sonnet-5-5
          input: 6
          output: 35
          cache_read: 44354
          cache_write: 30821
          cost: 0.0572
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: alex
    at: 2026-11-11T05:30:40Z
  value: 25
  by: planner-S-0721
  at: 2026-11-11T05:31:11Z
forecast:
  duration: 30m
  delivery: 2026-11-11T05:58:00Z
  basis: flai forecast's 13m (median 107 s per unit of size over 34 done improvement stories, size 7), raised to 30m for the four timed test runs, two of them under -race at about a minute each plus their builds, and the close-out's -race tier on mcpserver
  by: planner-S-0721
  at: 2026-11-11T05:29:22Z
---
# S-0721 Every MCP server flai builds resolves all its tools' JSON schemas again, which is 56% of the mcpserver tests' CPU: share one SchemaCache

## Goal

Every MCP server flai builds in one process takes its tools' JSON schemas from one shared `mcp.SchemaCache`, so a tool's schema is reflected and resolved once per process, not once per server. Nothing else changes: this story is the cache alone.

## Acceptance criteria
- [x] Both `mcp.NewServer` calls, in `flai/internal/mcpserver/server.go` (`New`) and `flai/internal/mcpserver/folder.go` (`newFolderServer`), set `ServerOptions.SchemaCache` to one package-level `mcp.NewSchemaCache()`.
- [x] A test builds two servers from the same options and shows `tools/list` answers the same tools with the same input and output schemas from both, and the same as a server built without the cache; and that servers built from differing options (single project and folder, each role) still list the schemas they listed before. If any tool's schema depends on its options and not on its Go type alone, that tool is kept out of the cache or the story stops and says so on a thread.
- [x] The story's narrative records the package's CPU time before and after, measured as in Notes, with and without `-race`; the `-race` run's user CPU is under half the before figure.
- [x] `design/tech/go-libraries.md` says of the go-sdk row that flai shares one `SchemaCache`; no dependency is added or upgraded.

## Tasks
- T-3083 Set one package-level mcp.SchemaCache on both MCP servers, with a test that their tool schemas are unchanged
- T-3084 Say in go-libraries.md's go-sdk row that flai shares one SchemaCache
- T-3085 Measure the mcpserver package's CPU before and after the cache, with and without -race, and record it in the narrative

## Notes

Measured 2026-11-10 on main at 596e5975d, this host (24 cores, loaded), `./internal/mcpserver` alone, 337 tests and subtests, `-short -count=1 -skip ^TestGateway`:

| | Wall | User CPU |
|---|---|---|
| `-race` (what a close-out runs) | 56 s | 608 s |
| no `-race` | 10 s | 65 s |

To repeat it: `go test -race -short -c -o <bin> ./internal/mcpserver`, then from `flai/internal/mcpserver` run `time <bin> -test.short -test.count=1 -test.skip '^TestGateway' -test.cpuprofile <prof>` and `go tool pprof -top -cum <bin> <prof>`.

What the profile shows: `mcpserver.New` is 56% of the plain run's CPU (25.7 s of 46 s sampled) and 35% of the race run's, nearly all of it `addProjectTools` → `mcp.AddTool` → `jsonschema.(*Schema).Resolve`. Each test's fixture builds a server, and each server registers 53 tools, reflecting over and resolving the same input and output types every time. Under `-race` that reflection-heavy walk is what the detector instruments most (`__tsan_read`, `checkptr*`), which is how 65 s becomes 608 s. `os/exec` is under 1%. Three stories' close-outs ran this at once on 2026-11-10 and took the host to a load of 84.

The SDK in use, go-sdk v1.8.0, has the cache already (`mcp/schema_cache.go`, `ServerOptions.SchemaCache`); flai does not set it. Its documented trade-offs: it is unbounded, one entry per Go type or schema pointer, which is fixed here; and a pre-defined schema is cached by pointer identity, so a schema changed in place behind the same pointer would be served stale. It is safe for concurrent use, which the parallel tests need.

The saving is an estimate from the profile, not a measurement: the cache was not tried. It also helps `flai mcp http`, which builds a server per session.

Out of scope, each its own story if this one's figures leave it worth doing: building fewer servers in the tests, running fewer packages under `-race` in a close-out, and limiting how many close-outs run at once.

### Planning

Touches:

| Path | Source | Why |
|---|---|---|
| `flai/internal/mcpserver/server.go` | declared | `New`'s `mcp.NewServer` call and the package-level cache (T-3083) |
| `flai/internal/mcpserver/folder.go` | declared | `newFolderServer`'s `mcp.NewServer` call; `permission_prompt`'s pre-defined schema is checked here (T-3083) |
| `design/tech/go-libraries.md` | declared | the go-sdk row (T-3084) |
| `flai/internal/mcpserver/schema_cache_test.go` | design | criterion 2's test, in a new file rather than `server_test.go` or `folder_test.go`, which `flai touches suggest` lists as co-changing in 42% and 26% of commits and which other stories touch often |
| `wip/agents/S-0721.md` | design | criterion 3's figures go in the narrative (T-3085) |

Left out of what `flai touches suggest` lists: `docs/users/flai.md`, `design/system/flai-cli.md`, `docs/operators/settings.md`, and `flai/cmd/mcp*.go`. The cache changes no command, setting, or behaviour a user sees.

Forecast: `flai forecast` gave 13m, from the median of 107 s per unit of size over 34 done improvement stories, times size 7. I raised it to 30m. Four timed runs of the package are needed, two of them under `-race` at about a minute each plus their builds, and the close-out runs the `-race` tier on this package too. Delivery counts from the start at 2026-11-11T05:28:00Z.

Cost of delay: 25 USD a week, as `flai cod` works it out from the operator's input on TH-1067. That input is 10m of time lost per 168h cycle, at 150 USD an hour and one cycle a week. It stands unadjusted: the figure is the operator's, and it matches D-0441 for the same kind of cost.

No `after`: the story reads nothing another open story writes.
