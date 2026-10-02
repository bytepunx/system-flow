---
id: S-0196
type: story
nature: improvement
title: A convention's roles name every agent that reads it, the story's agent included
status: done
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-02T23:49:45Z
transitions:
  - to: ready
    at: 2026-10-02T23:24:57Z
    by: alex
  - to: in-progress
    at: 2026-10-02T23:25:32Z
    by: agent-S-0196
  - to: review
    at: 2026-10-02T23:49:29Z
    by: agent-S-0196
  - to: done
    at: 2026-10-02T23:49:45Z
    by: alex
tags: [flai, template]
topics: [conventions]
touches: [flai/internal/conventions, flai/internal/context, flai/cmd/prime.go, design/conventions, template/root/design/conventions, design/system/conventions.md, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/conventions.md, docs/users/flai-reference.md, flai/cmd/prime_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1477
  models:
    - model: claude-haiku-4-5-20251001
      input: 186
      output: 7458
      cache_read: 1005625
      cache_write: 60240
      cost: 0.2133
    - model: claude-opus-5-5
      input: 196
      output: 48195
      cache_read: 8934342
      cache_write: 203820
      cost: 4.2097
    - model: claude-sonnet-5-5
      input: 32
      output: 7536
      cache_read: 498387
      cache_write: 90945
      cost: 0.4025
---
# S-0196 A convention's roles name every agent that reads it, the story's agent included

## Goal

[ADR-0068](../../../design/adrs/0068-a-convention-s-roles-list-every-agent-that-reads-it-the-story-s-agent-included.md) (2026-10-02) refines ADR-0059: a convention's `roles` lists every agent that reads it, `story` included, and a convention without `roles` is read by all agents. Today `roles` names only sub-agents (`flai/internal/conventions/conventions.go` ~36: `explore`, `verify`), the story's agent reads every convention, and `flai check` warns about `story`.

## Acceptance criteria
- [x] `flai check` accepts `story`, `explore`, `verify`, `orchestrator`, `planner`, and `analyzer` in a convention's `roles`, and warns about any other value
- [x] `flai prime --story S-nnnn` (and the MCP `prime`) prints only the conventions whose `roles` are empty or list `story`; `--role explore|verify` those whose roles are empty or list the role; `--cat` prints every convention; topics apply as before
- [x] The baseline's `roles`, here and in the template, are set as the designer chose on 2026-10-02, with `story` kept on every file the story's agent reads today unless the designer chose otherwise: `code-quality` [story, verify]; `communication` [story, explore, verify] or none; `continuous-improvement` [story]; `decisions` [story]; `delegation` [story, explore, verify]; `documentation` none; `git` [story]; `logging` [story, verify]; `safety` none; `session-start` [story]; `tooling` [story, explore, verify] or none; `work-management` [story]; `telemetry` and `README` none. Confirm the open choices with the designer
- [x] `design/system/conventions.md` (its Roles section and table) and the user guide describe the new meaning
- [x] Tests cover each role's selection, no roles, and the check

## Tasks
- T-0714 flai primes and checks roles as ADR-0068 says
- T-0715 The design and the user guide describe roles as ADR-0068 means them
- T-0716 The baseline's roles name every agent that reads each convention

## Notes

The designer set some of these values in the template's copies on 2026-10-02; they were held back from the conventions until this story, because `flai check --strict` warns on `story` until it lands.

On TH-0071 (2026-10-02) the designer chose no roles for `communication` and `tooling`, so every agent reads them.
