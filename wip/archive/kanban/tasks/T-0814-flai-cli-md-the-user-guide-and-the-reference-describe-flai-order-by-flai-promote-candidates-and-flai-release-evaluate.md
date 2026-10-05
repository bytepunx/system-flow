---
id: T-0814
type: task
nature: feature
title: flai-cli.md, the user guide, and the reference describe flai order --by, flai promote --candidates, and flai release --evaluate
status: done
parent: S-0217
owner: alex
created: 2026-10-04T19:09:37Z
updated: 2026-10-05T06:56:41Z
transitions:
  - to: ready
    at: 2026-10-05T06:25:32Z
    by: agent-S-0217
  - to: in-progress
    at: 2026-10-05T06:25:32Z
    by: agent-S-0217
  - to: done
    at: 2026-10-05T06:56:41Z
    by: agent-S-0217
stream: S-0217
tags: [flai]
touches: [design/system/flai-cli.md, design/system/server-performance.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-0813]
usage:
  source: log
  seconds: 1869
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 94
      output: 554
      cache_read: 4465653
      cache_write: 123585
      cost: 3.9536
---
# T-0814 flai-cli.md, the user guide, and the reference describe flai order --by, flai promote --candidates, and flai release --evaluate

## Work

Describe the three operations in `design/system/flai-cli.md` under `## Commands`. Cover each policy's arithmetic and how a missing figure and a tie are ordered, what makes a promotion candidate, how each release policy is evaluated, and the hostapi read and MCP tool of each. Say that they are the orchestrator's arithmetic, kept in flai so that the operator and the dashboard get the same answer without an agent (the designer's choice of 2026-10-02).

Describe them for the operator in `docs/users/flai.md`, next to `flai order` and `flai release`, and add the three tools wherever the guide lists the MCP tools.

Regenerate `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` with `scripts/flai-reference.sh` (`make flai-reference`). Do not edit them by hand.

This task waits for T-0813. The commands' flags, the reads, and the tools must be final before they are described and the reference regenerated.

## Done when

- `flai-cli.md` and `flai.md` describe the three commands, their policies, their reads, and their tools
- `flai-reference.md` and `settings.md` are what `make flai-reference` writes, with no diff after running it again
- the markdown lint and `flai check --strict` report nothing on these files

## Notes
