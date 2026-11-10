---
id: T-3022
type: task
nature: remediation
title: An ADR, flai-cli.md, and flai.md say what the MCP tool test's answer leaves out
status: backlog
parent: S-0707
owner: alex
created: 2026-11-10T10:32:10Z
updated: 2026-11-10T10:32:10Z
transitions: []
stream: S-0707
tags: [flai, mcp, docs]
touches: [design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3019]
---
# T-3022 An ADR, flai-cli.md, and flai.md say what the MCP tool test's answer leaves out

## Work

Write a new ADR with `flai adr new` and accept it: the MCP tool `test` answers `flai test --json`'s result with a tier's `command` only for the tier that failed and a `files` count for the others, and fits 40,000 bytes, cutting `paths` with `paths_omitted`. Link ADR-0179 and ADR-0198 as the precedents, and I-0260 as the cause. flai writes the index row in `design/adrs/README.md`.

In `design/system/flai-cli.md`, change the `flai test` row's account of the MCP tool to match, linking the ADR. In `docs/users/flai.md`, section "Run the tests for what changed", say what the MCP tool's answer leaves out, and that `flai test --json` keeps every tier's whole command. It waits for T-3019 so that both describe the shape that task built.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`.
- `design/system/flai-cli.md` and `docs/users/flai.md` say what the MCP answer leaves out and where the whole commands are.
- The markdown lint passes on the changed files.

## Notes
