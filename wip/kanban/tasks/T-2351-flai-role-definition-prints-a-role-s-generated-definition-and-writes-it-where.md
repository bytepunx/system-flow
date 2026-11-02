---
id: T-2351
type: task
nature: improvement
title: flai role definition prints a role's generated definition and writes it where the ADR puts it
status: done
parent: S-0553
owner: alex
created: 2026-10-30T20:47:27Z
updated: 2026-11-02T00:44:05Z
transitions:
  - to: ready
    at: 2026-11-02T00:33:40Z
    by: agent-S-0553
  - to: in-progress
    at: 2026-11-02T00:33:42Z
    by: agent-S-0553
  - to: done
    at: 2026-11-02T00:44:05Z
    by: agent-S-0553
stream: S-0553
tags: [cli]
touches: [flai/cmd/role.go, flai/cmd/role_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, docs/operators/settings.md]
after: [T-2347, T-2348, T-2350]
usage:
  source: log
  seconds: 623
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 69
      output: 21321
      cache_read: 5251268
      cache_write: 111406
      cost: 2.2013
---
# T-2351 flai role definition prints a role's generated definition and writes it where the ADR puts it

## Work

Add `flai role definition <name>` to `flai/cmd/role.go`. It prints the harness definition generated from the role's file as a Claude Code agent file, front matter and prompt. With `--write`, it writes it where the ADR puts it, outside `.claude/`, and says what the operator does once. Document the subcommand in `docs/users/flai.md`, regenerate `docs/users/flai-reference.md`, and add it to `design/system/flai-cli.md`.

Second layer: it waits for the generation task, whose definition it prints, and for the ADR, which says where it writes. It shares no path with the check task.

## Done when

`role_test.go` covers the print for a built-in role and a defined one, `--write` to the ADR's place, and an unknown name refused. The reference test passes and `flai check --strict` reports nothing.

## Notes

Drafted by the planner of E-0024 on 2026-10-30.
