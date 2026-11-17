---
id: T-3606
type: task
nature: improvement
title: The work-management and strategic-agents conventions, here and in the template, describe a publish as a release_publish call, and the template's changelog records the change
status: backlog
parent: S-0837
owner: alex
created: 2026-11-17T12:38:19Z
updated: 2026-11-17T12:38:19Z
transitions: []
stream: S-0837
tags: [template, docs]
touches: [design/conventions/work-management.md, design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-3601, T-3602, T-3603, T-3605]
---
# T-3606 The work-management and strategic-agents conventions, here and in the template, describe a publish as a release_publish call, and the template's changelog records the change

## Work

This task waits for T-3601, T-3602, T-3603, and T-3605, because it describes what they built: the tool's fetch and `covered_by`, the refusal of every formula, the role's words, and the formula's removal.

- `work-management.md`, the rule on publishing (around line 19): an agent the operator asks to publish still runs `git fetch` and `flai release --pending`. Drop the link to the publish formula and the words that make a publish its steps.
- `strategic-agents.md`:
  - The rule that a strategic agent never runs a formula "but the orchestrator's `publish`" (around line 23): the orchestrator runs no formula.
  - The rule on accepting and publishing (around line 24): keep it as it is unless it names the formula.
  - § As the orchestrator, `publish` (around line 52): when `release_evaluate` is met, call `release_publish` with the met figure, or the reasoning under `judgement`, as its reason. Drop `flai formula run publish`, the `judgement` input, and the quoting rule ADR-0273 needed.
- Make each change in the template's copy above the baseline marker. Edit this repository's copy to match.
- Add one `## Unreleased` entry to `template/CHANGELOG.md` covering the template's whole change: the formula and `scripts/fetch.sh` retired, the role file, the conventions, and that `flai upgrade` removes an unchanged copy.

## Done when

- Neither convention, in either copy, describes a publish as a story or names the publish formula.
- The changelog has the entry.
- `flai check --strict` and the markdown lint pass.

## Notes

- Drafted by planner-S-0837.
