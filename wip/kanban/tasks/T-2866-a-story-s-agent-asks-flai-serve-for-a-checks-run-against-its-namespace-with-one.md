---
id: T-2866
type: task
nature: feature
title: A story's agent asks flai serve for a checks run against its namespace with one flai command and one MCP tool, and gets each check's result and log
status: backlog
parent: S-0661
owner: alex
created: 2026-11-10T08:24:11Z
updated: 2026-11-13T12:42:04Z
transitions: []
stream: S-0661
tags: [cli]
touches: [flai/cmd/checks.go, flai/cmd/checks_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go]
after: [T-2863, T-2865]
---
# T-2866 A story's agent asks flai serve for a checks run against its namespace with one flai command and one MCP tool, and gets each check's result and log

## Work

Give the agent the request T-2863's pass takes, in the shape S-0660 gave the deploy request: the same hand-off to flai serve, the same wait, and the same answer when the wait runs out. Use the command and tool names the ADR S-0659 accepted. This plan assumes `flai checks run S-nnnn --job` in `flai/cmd/checks.go`, and the MCP tool `kluster_checks` beside S-0660's deploy tool in `flai/internal/mcpserver/kluster.go`, registered in `folder.go` and listed in `server_test.go`'s `TestToolsAreAdvertised`.

- The request names the story and is refused for another story than the agent's own. The refusals T-2863 gives reach the agent whole, saying to deploy first where it is not deployed at the head.
- The answer holds the outcome, the commit, the namespace, each check's name, outcome, exit, and log, and the run's log. A request that outlasts its wait says the run is queued or running, and that `flai checks status S-nnnn` reads it later.
- `flai checks run S-nnnn` without `--job` keeps running in the caller as today, for the operator.
- A sub-agent is refused the tool and `--job`, as S-0660 has it refused the deploy request. S-0660 found a tool no read table lists is refused a sub-agent already; change `flai/internal/guard/guard_tables.go` only if `--job` needs it.

Tests in `checks_test.go` and `kluster_test.go` fake flai serve's answer. They cover a passing run, a failing one, the not-deployed refusal, another story refused, and a wait that runs out. It waits for T-2863, whose request and answer it carries.

## Done when

- `flai checks run S-nnnn --job` and the MCP tool ask flai serve for the checks against the story's namespace, and answer each check's result and log, or the refusal, or that the run is still queued or running.
- A sub-agent is refused both; `go test ./cmd/ ./internal/mcpserver/ ./internal/guard/` passes.

## Notes
