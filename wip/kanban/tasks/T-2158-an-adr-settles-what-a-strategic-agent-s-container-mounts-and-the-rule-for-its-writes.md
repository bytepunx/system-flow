---
id: T-2158
type: task
nature: feature
title: An ADR settles what a strategic agent's container mounts and the rule for its writes
status: backlog
parent: S-0509
owner: alex
created: 2026-10-26T04:48:25Z
updated: 2026-10-26T04:48:25Z
transitions: []
stream: S-0509
tags: [cli, docker]
touches: [design/adrs/README.md, design/system/strategic-agents.md]
---
# T-2158 An ADR settles what a strategic agent's container mounts and the rule for its writes

## Work

Write an ADR with `flai adr new`, refining ADR-0203 and the ADR S-0507 wrote, that settles the mounts of a strategic agent's container. The planner, the orchestrator, the analyzer, and the curator run in the main checkout today (`flai/internal/serve/plan.go`, `orchestrate.go`, `analyze.go`, `curate.go` start them in `entry.Root`): the planner and the orchestrator write only through flai's MCP tools and CLI, the analyzer writes its report under `design/analysis/` and commits nothing, and the curator writes the digest through `flai digest`. Weigh two choices: the main checkout mounted read-write at its host path, with the main checkout's `.git` inside the mount and the role's guard refusing the writes it refuses today; or the main checkout read-only with every write over the HTTP MCP server, which the analyzer's report and the curator's `flai digest` do not reach today. Recommended: read-write for the four roles, since their writes are the ones the guard already bounds and the container still mounts no home, no `~/.ssh`, and no `~/.flai`, so no credential leaves the host; name the risk (a role writing the main checkout's `.git`) and what the guard refuses inside the container for each role. Record the role's cache folder, `.flai-cache/agents/<role>`, for its session. Add a paragraph to `design/system/strategic-agents.md` under each role's "Its agent" section. This task waits for nothing.

## Done when

The ADR is accepted and indexed in `design/adrs/README.md`, `design/system/strategic-agents.md` says for each role what its container mounts, and `flai check --strict` passes.

## Notes

The ADR's file is named by `flai adr new`; the story keeps the `design/adrs` folder touch for it.
