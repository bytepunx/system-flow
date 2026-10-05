---
id: T-0877
type: task
nature: remediation
title: The ADR, design, and users' docs describe permission_prompt and auto-approve, and I-0069 is closed
status: done
parent: S-0257
owner: alex
created: 2026-10-05T04:15:28Z
updated: 2026-10-05T04:30:50Z
transitions:
  - to: ready
    at: 2026-10-05T04:15:53Z
    by: agent-S-0257
  - to: in-progress
    at: 2026-10-05T04:27:27Z
    by: agent-S-0257
  - to: done
    at: 2026-10-05T04:30:50Z
    by: agent-S-0257
stream: S-0257
tags: []
touches: [design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/delegation.md, design/issues]
after: [T-0874, T-0875, T-0876]
usage:
  source: log
  seconds: 203
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 84
      output: 449
      cache_read: 3168555
      cache_write: 84691
      cost: 1.3432
---
# T-0877 The ADR, design, and users' docs describe permission_prompt and auto-approve, and I-0069 is closed

## Work

1. Write a new ADR, refining ADR-0038 and naming ADR-0067. It records three things:
   - why the headless claude-code agent gets a permission handler;
   - what `permission_prompt` approves, and what it never approves;
   - the shell-only `auto-approve` host action.
2. Update the documents that describe this:
   - `design/system/flai-cli.md`: the MCP tools, flai serve's claude-code start, and the host actions;
   - `docs/users/flai.md` and `docs/users/flai-reference.md`;
   - the project addition in `design/conventions/delegation.md` about `.claude/`: an agent edits these files and `permission_prompt` asks the operator, or allows the write at once under `auto-approve`.
3. Close I-0069 with `flai issue close I-0069 --reason`.

## Done when

The documents agree with the code, the markdown lint is clean, and I-0069 is closed.

## Notes
