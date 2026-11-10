---
id: T-2901
type: task
nature: remediation
title: The guard counts the read forms of flai serve among the reads
status: done
parent: S-0680
owner: alex
created: 2026-11-10T08:46:44Z
updated: 2026-11-10T09:01:30Z
transitions:
  - to: ready
    at: 2026-11-10T08:52:28Z
    by: agent-S-0680
  - to: in-progress
    at: 2026-11-10T08:52:29Z
    by: agent-S-0680
  - to: done
    at: 2026-11-10T09:01:30Z
    by: agent-S-0680
stream: S-0680
tags: [flai, guard]
touches: [flai/internal/guard/guard_tables.go, flai/internal/guard/guard.go, flai/internal/guard/hostconfig.go, flai/internal/guard/guard_test.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/orchestrate_test.go]
usage:
  source: log
  seconds: 541
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 22
      output: 142
      cache_read: 3088711
      cache_write: 19188
      cost: 1.1436
---
# T-2901 The guard counts the read forms of flai serve among the reads

## Work

Propose the fix on the story before building it, as its goal asks: the plan below is the proposal, so record it in the narrative's `## Decisions` and go on unless the operator's answer on the plan's thread says otherwise.

`reads` in `flai/internal/guard/guard.go` finds no `serve` entry in `cliReads` or `flagReads` (`flai/internal/guard/guard_tables.go`), so the guard refuses the orchestrator, every strategic agent, and every sub-agent `flai serve status`, though `askForHost` in `flai/internal/guard/hostconfig.go` tells them it is open, and ADR-0220 lists it among the reads.

- Add `serve` to `flagReads` with a form function, `readsServe`, beside `readsRoles`, that passes exactly the read forms ADR-0220 lists: `status`, `actions`, `journal`, `agent` bare and `agent show`, `agent harness` and `agent provider` bare or with a name alone, `checks show`, `components show`, `forge` bare or with a forge alone, `notify show`, `mcp show`, and `import list`, each with no setting flag of `hostSettings` or `forgeSettings`. Reuse `subcommands` with `hostValues` and `writesHost`, so the two tables never disagree.
- Pass nothing else: `flai serve` alone, `start`, `stop`, `project`, `usage`, `agent start`, `restart`, `commit`, `stream`, `orchestrate`, and every form `writesHost` finds stay refused as they are now.
- Make `askForHost` name only what `readsServe` passes.
- Tests: in `flai/internal/guard/orchestrate_test.go`, `TestARedirectionNeitherAllowsNorRefusesTheOrchestratorsCommand` takes `flai serve status` as passing, with and without `2>&1`, and its comment says so. In `flai/internal/guard/guard_test.go`, a table test that a read-only sub-agent (explorer) and the orchestrator pass each read form, and are refused `flai serve`, `flai serve start`, `flai serve agent start S-0001`, `flai serve stop`, and a setting form such as `flai serve agent provider litellm --api x`. In `flai/internal/guard/hostconfig_test.go`, the `reads` list there also passes `reads`.

Waits for no task: it is the first layer.

## Done when

- The guard passes `flai serve status`, with and without a redirection, for the orchestrator and for a read-only sub-agent, and the tests above pin it.
- Every `flai serve` form that starts, stops, or writes is refused as before, and `flai test flai/internal/guard` passes.

## Notes
