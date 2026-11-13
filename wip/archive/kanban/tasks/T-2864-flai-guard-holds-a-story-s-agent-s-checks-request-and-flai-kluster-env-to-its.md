---
id: T-2864
type: task
nature: feature
title: flai guard holds a story's agent's checks request and flai kluster env to its own story's namespace, refusing every sub-agent the request
status: done
parent: S-0661
owner: alex
created: 2026-11-10T08:23:35Z
updated: 2026-11-13T12:58:30Z
transitions:
  - to: ready
    at: 2026-11-13T12:42:10Z
    by: agent-S-0661
  - to: in-progress
    at: 2026-11-13T12:42:10Z
    by: agent-S-0661
  - to: done
    at: 2026-11-13T12:58:30Z
    by: agent-S-0661
stream: S-0661
tags: [cli]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go]
usage:
  source: log
  seconds: 980
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 18341
      cache_read: 3499757
      cache_write: 104734
      cost: 1.6931
---
# T-2864 flai guard holds a story's agent's checks request and flai kluster env to its own story's namespace, refusing every sub-agent the request

## Work

ADR-0292 §2 decides what `flai guard` refuses: `--job` requests, `env`, `exec`, and `log` for another story's namespace or the default one. S-0660 built the rule for `flai kluster deploy|teardown S-nnnn --job` and the tools `kluster_deploy` and `kluster_teardown` (`klustering`, `klusterTool`), and S-0763 the rule for `exec` and `log` (`klusterReaches`, `reaching`) and the kubeconfig and `kluster` program refusals. The API server's refusal through the namespace's own kubeconfig is the boundary, so no rule reads `kubectl -n` or `--namespace`: that was the planner's draft before the ADR, and the ADR rejects it.

In `flai/internal/guard/guard_rules.go`, extend the existing rules, no wider:

- `flai checks run S-nnnn --job`: a story's agent, served or not, with no role, may run it for its own story alone, written plainly, as it may `flai kluster deploy S-nnnn --job`, with flags in any order and form, flai's own flags before `checks`, through `scripts/flai.sh` or a path, after `cd <dir> &&` or in `bash -c`. Refused: another story, none, two, a line not written plainly, `--job` set false, and every sub-agent, with the reason that the request is the story's agent's. `flai checks run` without `--job` keeps whatever rule it has today for each session.
- The MCP tool `kluster_checks`: as `kluster_deploy` in `klusterTool`, its own story or none for the story's agent, refused to every sub-agent and for another story.
- `flai kluster env S-nnnn`: add `env` to `klusterReaches`, so a story's agent and its sub-agents run it for their own story alone, as they run `exec` and `log`, and a session with no story of its own other than the operator's is refused it.

Update `guard_tables.go` only where its tables must name the new tool or subcommand. Each refusal says what to write instead.

Tests in `kluster_test.go` beside S-0660's and S-0763's: the own story passing by command and by tool, another story and the default namespace refused, a sub-agent refused the request and the tool but let through `env` for its own story, `--job=false` refused, and the operator not refused. It waits for no task: the guard imports neither `serve` nor `kluster`, and names the tool and the command by their words.

## Done when

- `flai guard` lets a story's agent run `flai checks run S-nnnn --job`, call `kluster_checks`, and run `flai kluster env S-nnnn` for its own story, and refuses them for another story or the default namespace, and the request and tool to every sub-agent, each with what to do instead.
- `flai test flai/internal/guard` passes.

## Notes
