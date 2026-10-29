---
id: T-2310
type: task
nature: improvement
title: The strategic-agents and work-management conventions, here and in the template, mark the sections each orchestrator reason needs with section topics, so a reason pack takes out the rest and every other pack keeps what it kept
status: done
parent: S-0537
owner: alex
created: 2026-10-28T22:07:13Z
updated: 2026-10-29T15:37:55Z
transitions:
  - to: ready
    at: 2026-10-29T15:25:22Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T15:25:23Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T15:37:55Z
    by: agent-S-0537
stream: S-0537
tags: [conventions, template, orchestrator]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/conventions/tooling.md]
usage:
  source: log
  seconds: 752
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 10408
      cache_read: 2079778
      cache_write: 67202
      cost: 1.0062
---
# T-2310 The strategic-agents and work-management conventions, here and in the template, mark the sections each orchestrator reason needs with section topics, so a reason pack takes out the rest and every other pack keeps what it kept

## Work

Criterion 2 on this repository, by ADR-0217 §2. The orchestrate role's conventions were 64,370 bytes of the general pack's 175,725 here on 2026-10-29. `work-management.md` is 24 KB and `strategic-agents.md` 22 KB. Today no convention heading carries a section topic, so every section inherits its file's `[all]`.

Every heading comment this task adds names `all` first, then the words: `<!-- topics: all, review -->`. `all` keeps the section in every pack that has it today. Only a reason pack reads the other words.

- `strategic-agents.md`:
  - Give `## As the planner` `all, planning, planned`, `## As the analyzer` `all, analysis`, and `## As the curator` `all, digest`.
  - Split `## As the orchestrator` into `###` sections, moving its bullets without changing their words. Keep the intro and the rules for every action (guard, figures, logging, never edit) under the `##` heading, untagged. Each `###` section takes `all` and the reasons that need it:
    - planning, and a story planner's threads: `planned`, `thread`
    - drafts: `draft`
    - promotion and order: `ready`
    - threads, answering and recommending: `thread`, `planned`
    - acceptance: `review`
    - releases: `release`
    - sub-projects: `thread`
  - Leave `## Rules`, `## Logging your activity`, `## Asking the operator`, and `## When in doubt` untagged.
- `work-management.md`: give `## Rules` `all, review, draft, planned`. A run for a thread, a release, or the ready column leaves it out; every other pack keeps it.
- Copy both baselines to `template/root/design/conventions/`. Add a line to `template/CHANGELOG.md`. Leave the `## Project additions` headings and what is below the marker alone.

No `after`. It runs beside T-2308, since they share no path. Until T-2316 lands, `flai check` warns `doc.topic` on the new words where no story declares them. T-2316 clears that, and the close-out runs after both.

## Done when

- `flai prime --story` on a story of this repository, and `flai prime --role plan`, `orchestrate`, `analyze`, and `curate`, keep every rule they kept before. Compare their text before and after: only headings differ.
- The template copies match the baselines above the marker, and the markdown lint is clean.
- `## Decisions` names the words and which sections take each.

## Notes
