---
id: T-2310
type: task
nature: improvement
title: The strategic-agents and work-management conventions, here and in the template, mark the sections each orchestrator reason needs with section topics, so a reason pack takes out the rest and every other pack keeps what it kept
status: backlog
parent: S-0537
owner: alex
created: 2026-10-28T22:07:13Z
updated: 2026-10-28T22:07:13Z
transitions: []
stream: S-0537
tags: [conventions, template, orchestrator]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
---
# T-2310 The strategic-agents and work-management conventions, here and in the template, mark the sections each orchestrator reason needs with section topics, so a reason pack takes out the rest and every other pack keeps what it kept

## Work

Criterion 2 on this repository. The orchestrate role's conventions are about 64 KB of the general pack's 170 KB here (measured on 2026-10-28 with `flai prime --role orchestrate --json`): loaded whole, they alone are over a quarter of it. `work-management.md` is 24 KB and `strategic-agents.md` 22 KB, and today no convention heading carries a section topic (`<!-- topics: ... -->`, `flai/internal/topics/topics.go`), so every section inherits its file's `[all]`.

- In `strategic-agents.md`, give `## As the planner`, `## As the analyzer`, and `## As the curator` the section topics of their roles' packs (`planning`, `analysis`, `digest`), and split `## As the orchestrator` into `###` sections per reason, the conventions on threads (answering, recommending, a story planner's threads), acceptance, finalizing drafts, promotion and order, release and sub-projects, each with its reason's topic and `orchestration`.
- In `work-management.md`, group the rules the orchestrator acts on (definitions of ready and done, publishing, cancelling) under `###` headings whose topics keep them in every pack that has them today and name the reasons that need them, so that a reason pack can leave out the rest. The words of the rules do not change: only headings and the order of the bullets under them.
- Copy both to the template and add a line to `template/CHANGELOG.md`.

The section topics are the vocabulary T-2308 filters by; agree the names in the narrative's `## Decisions` before either task closes. When a section topic would cut a rule from a story's pack, as a section topic without `all` would, choose the rule T-2308 applies instead (record it in `## Decisions`, and in an ADR under `design/adrs` if it changes how a pack selects sections, refining ADR-0049 and ADR-0075).

No `after`: it runs beside T-2308.

## Done when

- `flai prime --story` on a story of this repository and `flai prime --role plan`, `analyze`, and `curate` keep every rule they kept before, compared before and after, with only headings differing.
- The template copies match, and `scripts/flai.sh check --strict` and the markdown lint are clean.
- `## Decisions` names the section topics and which reasons take each.

## Notes
