---
id: S-0196
type: story
nature: improvement
title: A convention's roles name every agent that reads it, the story's agent included
status: backlog
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-02T09:46:13Z
transitions: []
tags: [flai, template]
topics: [conventions]
touches: [flai/internal/conventions, flai/internal/context, flai/cmd/prime.go, design/conventions/, template/root/design/conventions]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0196 A convention's roles name every agent that reads it, the story's agent included

## Goal

[ADR-0068](../../../design/adrs/0068-a-convention-s-roles-list-every-agent-that-reads-it-the-story-s-agent-included.md) (2026-10-02) refines ADR-0059: a convention's `roles` lists every agent that reads it, `story` included, and a convention without `roles` is read by all agents. Today `roles` names only sub-agents (`flai/internal/conventions/conventions.go` ~36: `explore`, `verify`), the story's agent reads every convention, and `flai check` warns about `story`.

## Acceptance criteria
- [ ] `flai check` accepts `story`, `explore`, and `verify` in a convention's `roles`, and warns about any other value
- [ ] `flai prime --story S-nnnn` (and the MCP `prime`) prints only the conventions whose `roles` are empty or list `story`; `--role explore|verify` those whose roles are empty or list the role; `--cat` prints every convention; topics apply as before
- [ ] The baseline's `roles`, here and in the template, are set as the designer chose on 2026-10-02, with `story` kept on every file the story's agent reads today unless the designer chose otherwise: `code-quality` [story, verify]; `communication` [story, explore, verify] or none; `continuous-improvement` [story]; `decisions` [story]; `delegation` [story, explore, verify]; `documentation` none; `git` [story]; `logging` [story, verify]; `safety` none; `session-start` [story]; `tooling` [story, explore, verify] or none; `work-management` [story]; `telemetry` and `README` none. Confirm the open choices with the designer
- [ ] `design/system/conventions.md` (its Roles section and table) and the user guide describe the new meaning
- [ ] Tests cover each role's selection, no roles, and the check

## Tasks

## Notes

The designer set some of these values in the template's copies on 2026-10-02; they were held back from the conventions until this story, because `flai check --strict` warns on `story` until it lands.
