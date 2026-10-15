---
id: S-0386
type: story
nature: feature
title: A story's agent working on a project's UI is primed with its design system, and a baseline convention tells it to build from the tokens and propose changes on its branch
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:45:42Z
updated: 2026-10-15T01:10:41Z
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
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 1140
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 611
          cache_read: 6969486
          cache_write: 10695
          cost: 1.7539
cost_of_delay:
  value: 134.62
  by: planner-E-0020
  at: 2026-10-13T17:52:29Z
forecast:
  duration: 35m
  delivery: 2026-10-15T03:43:00Z
  basis: "Its own forecast of 35m; 6th in the pull order with an in-progress limit of 5, behind S-0381, S-0383, S-0382, S-0241, S-0384 and S-0385."
  by: flai
  at: 2026-10-15T01:10:41Z
finalized:
  by: alex
  at: 2026-10-15T00:58:30Z
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
- T-1561 flai design init --project records in the design system document which sub-projects it governs
- T-1562 flai prime --story briefs the design system and its convention for a story that reaches a sub-project it governs, and leaves both out where there is none
- T-1563 This repository's design records, and its conventions gain, a design-system convention: build UI from the tokens and propose changes on the story's branch
- T-1564 The template ships the design-system convention as baseline and lists it in its conventions README
- T-1565 The prime, flai design, and user docs and the reference describe governs, init --project, and how the design system reaches a story's agent

## Notes

- Drafted by the planner for E-0020. A baseline convention reaches every project made from the template with its next release, so this story changes the template as `CLAUDE.md` says: `design/system` first, then `template/`, in the same story.

### Planning

- **Touches.** None declared. All are the planner's, every one a file:
  - `flai/cmd/design.go` with its test: S-0381's `init`, extended with `--project`. Criterion 1.
  - `flai/internal/context/design.go` with its test: layout; it chooses which `design/system` documents a pack briefs (`folders`, `LoadDocs`). T-1562 may widen to the file that selects conventions if that is not `design.go`. Criteria 1 and 3.
  - `design/system/conventions.md`, `design/conventions/design-system.md`, `design/conventions/README.md`, `template/root/design/conventions/design-system.md`, and `template/root/design/conventions/README.md`: design and layout, the two copies of the conventions. Criterion 2.
  - `design/system/agent-context.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and the reference: design. Criterion 4.
- **Co-change.** `flai touches suggest` lists `template/CHANGELOG.md`, `docs/operators/settings.md`, the dashboard docs, and `work-management.md` in both copies. None is added: the changelog is written when a release is published, the story adds no setting, and it changes no other convention.
- **No folder touch.**
- **Layers.** T-1561, T-1562, and T-1563 together; then T-1564 and T-1565 together.
- **Forecast: 35m, delivery 2026-10-14T01:05:00Z.** `flai forecast` gave 27m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 17. Raised to 35m because a baseline convention is written in two copies, recorded in the design first, and checked on a project rendered from the template, which a rate per touch does not weigh. The delivery adds the 8m to flai's 2026-10-14T00:57Z.
- **Cost of delay: 134.62 USD a week.** `flai cod`: 35m of the 5h12m forecast over E-0020's seven open stories, of the epic's 1200 USD a week. It stands.
