---
id: S-0386
type: story
nature: feature
title: A story's agent working on a project's UI is primed with its design system, and a baseline convention tells it to build from the tokens and propose changes on its branch
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:45:42Z
updated: 2026-10-13T17:45:42Z
transitions: []
tags: [cli, template]
topics: [design, conventions]
touches: [flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/context/design.go, flai/internal/context/design_test.go, template/root/design/conventions/design-system.md, template/root/design/conventions/README.md, design/conventions/design-system.md, design/conventions/README.md, design/system/conventions.md, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0384, S-0385]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0386 A story's agent working on a project's UI is primed with its design system, and a baseline convention tells it to build from the tokens and propose changes on its branch

## Goal

Agents leverage the design system without being told each time: a story's agent that works on a sub-project whose UI the design system governs finds the design system in its prime pack, and a baseline convention tells it to build that UI from the tokens, to add no colour, size, or font outside them, and to propose a change to the design system on its own branch, where the operator previews it on the Design page before accepting.

## Acceptance criteria

- [ ] `flai design init --project <name>` records in the design system document which sub-projects it governs, and `flai prime --story` briefs it for every story whose topics or touches reach one of them.
- [ ] A baseline convention, `design-system.md`, in the template and in this repository, reaches the same stories' agents and says: build UI from the design system's tokens through the stylesheet `flai design css --write` keeps; add no colour, size, font, radius, or shadow outside them; make a change the UI needs to the design system document on the story's branch, so the operator previews it on the Design page before accepting; and record the reason for a design decision in the document's guidance.
- [ ] A project with no design system document primes as before, and the convention says nothing to its agents.
- [ ] `design/system/conventions.md`, `design/system/agent-context.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated reference say how.

## Tasks

## Notes

- Drafted by the planner for E-0020. A baseline convention reaches every project made from the template with its next release, so this story changes the template as `CLAUDE.md` says: `design/system` first, then `template/`, in the same story.
