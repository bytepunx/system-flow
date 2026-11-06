---
id: T-2461
type: task
nature: remediation
title: Say in the design, the delegation convention, and the template's changelog that the guard refuses a sub-agent the close-out job
status: backlog
parent: S-0556
owner: alex
created: 2026-11-05T01:37:28Z
updated: 2026-11-05T01:37:28Z
transitions: []
stream: S-0556
tags: [flai, template]
touches: [design/system/flai-cli.md, design/system/strategic-agents.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md]
after: [T-2460]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 65968
          cache_write: 1221
          cost: 0.016
---
# T-2461 Say in the design, the delegation convention, and the template's changelog that the guard refuses a sub-agent the close-out job

## Work

Bring the documents in step with what T-2459 and T-2460 built:

- `design/system/flai-cli.md`, the `flai guard` row: take `verify` (S-0270) out of the list of `guard.MCPReads`, say that the story's agent and the orchestrator call it, and that a sub-agent may run `flai verify` without `--record-issues` and without `--job`, naming this story
- `design/system/strategic-agents.md`: where it lists the orchestrator's tools and says how the guard holds it, name `verify` among the tools its role file lists; where it lists the planner's, the analyzer's, and the curator's reads, take `verify` out if it is there
- `design/conventions/delegation.md` and its template copy `template/root/design/conventions/delegation.md`, kept identical above the marker: the sub-agent rule "Never ask flai serve for the story's close-out job ... unless your prompt asks" becomes a rule the guard enforces, so say that `flai guard` refuses it, and drop "unless your prompt asks"
- `template/CHANGELOG.md`: one line under `## Unreleased` for S-0556, covering the orchestrator's role file naming `verify` and the delegation rule

Bump `updated` on each design file you change. Read each section with `doc_get` and its heading before you change it. This task waits for T-2460, because it describes the refusals that task builds.

## Done when

- No design document or convention says a sub-agent may call the MCP `verify` tool or ask for the close-out job
- `flai check --strict` and the markdown lint pass on the changed files

## Notes
