---
id: T-2175
type: task
nature: feature
title: A planner run in a container plans one story of this repository with its activity logged, each limitation is an issue, and the docs say how a strategic agent runs in a container
status: done
parent: S-0509
owner: alex
created: 2026-10-26T04:51:44Z
updated: 2026-11-02T19:07:33Z
transitions:
  - to: ready
    at: 2026-11-02T19:02:42Z
    by: agent-S-0509
  - to: in-progress
    at: 2026-11-02T19:02:44Z
    by: agent-S-0509
  - to: done
    at: 2026-11-02T19:07:33Z
    by: agent-S-0509
stream: S-0509
tags: [cli, docker]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/summary.md, design/system/agent-adapters.md]
after: [T-2170]
usage:
  source: log
  seconds: 289
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 18107
      cache_read: 5886670
      cache_write: 99354
      cost: 2.194
---
# T-2175 A planner run in a container plans one story of this repository with its activity logged, each limitation is an issue, and the docs say how a strategic agent runs in a container

## Work

Since TH-0815, the live planner run is S-0567's, after S-0509's release. This task records as issues the limitations ADR-0247 names and the work met: the container's flai has no host configuration, `.git/config` stays writable, and a role a project defines is refused on the container harness. It then writes the documentation: `design/system/flai-cli.md` (the strategic start paths through the adapter, beside § The container harness (S-0508)), `docs/users/flai.md`, and `docs/operators/index.md`, following `design/system/strategic-agents.md` § In a container (S-0509). Waits for T-2170, which finishes the start paths the documents describe.

## Done when

Each limitation has an issue, the three documents describe the behaviour the tests check, and `flai check --strict` and the markdown lint pass.

## Notes

The proof run costs a planner session; name the story it planned in the story's notes.
