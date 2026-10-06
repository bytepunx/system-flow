---
id: T-1023
type: task
nature: improvement
title: "An ADR refines ADR-0046: a story in review holds nothing, shared paths never hold, and tasks narrow a folder claim, and the design says so"
status: done
parent: S-0295
owner: alex
created: 2026-10-06T12:14:29Z
updated: 2026-10-06T18:26:24Z
transitions:
  - to: ready
    at: 2026-10-06T18:20:01Z
    by: agent-S-0295
  - to: in-progress
    at: 2026-10-06T18:20:01Z
    by: agent-S-0295
  - to: done
    at: 2026-10-06T18:26:24Z
    by: agent-S-0295
stream: S-0295
tags: [flai]
touches: [design/adrs, design/system/workflow.md, design/system/project-manifest.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 383
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 81
      output: 33819
      cache_read: 4556854
      cache_write: 133700
      cost: 2.3232
---
# T-1023 An ADR refines ADR-0046: a story in review holds nothing, shared paths never hold, and tasks narrow a folder claim, and the design says so

## Work

Write a new ADR with `flai adr new` that refines ADR-0046, from the three directions the operator chose on TH-0180:

- **Review holds nothing.** Only a story in progress holds a ready story. A story in review keeps its claim for `flai check`, the trial merge at sync, and the notice at acceptance, but no longer holds a ready story.
- **Shared paths never hold.** A manifest list of glob patterns names paths that many stories change in separate sections or new files. Propose `claims.shared` as its key. An overlap that lies wholly inside a shared path does not hold a story and is not a `wip.overlap`. The trial merge at sync and the notice at acceptance still report on it. Fix the glob dialect: `*` stays within one segment, `**` crosses segments, and a plain path matches itself and everything under it, so that `design/adrs` and `design/adrs/**` mean the same. The operator can edit the list through the CLI, HTTP, MCP, and the dashboard's project settings. Name the CLI (propose `flai shared list|add|remove|check`), the hostapi method (`settings.shared`), and the MCP tools, and say who may change the list through MCP: the operator's session, not a story's agent, its sub-agents, or the planner.
- **Tasks narrow a folder claim.** A story touch that is a folder holding at least one touch of the story's tasks is replaced, in the claim, by those task touches. Done tasks count, because the branch changed their files; cancelled ones do not. A folder touch that no task names inside stays whole.

Record the defaults. This project's list is `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `design/adrs`, `design/issues`, and `template/CHANGELOG.md`. Say whether a new project from the template starts with the same list or with an empty one; recommend the same list, since its paths come from the template's own layout.

Then update `design/system/workflow.md` § Branches and collisions, `design/system/project-manifest.md` (the new key), and `design/system/flai-cli.md` (the new command and tools). This task waits for nothing: every other task builds what it decides.

## Done when

- The ADR is accepted in the story branch with `refines: [ADR-0046]`. It names the key, the glob dialect, the defaults, the CLI, the HTTP method, the MCP tools, and who may use the MCP edit tool.
- `workflow.md`, `project-manifest.md`, and `flai-cli.md` describe the three rules and link the ADR.
- `flai check --strict` and the markdown lint pass.

## Notes
