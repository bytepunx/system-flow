---
id: ADR-0066
title: An experiment story is accepted like any other and records its results in a document under design/experiments
status: accepted
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0025, ADR-0032]
---

# ADR-0066 An experiment story is accepted like any other and records its results in a document under design/experiments

## Context

ADR-0025 decided that an `experiment` story stays on its branch: acceptance refuses it before anything is merged, and `flai accept --no-release` is the deliberate way to land one. ADR-0032 kept that refusal when acceptance stopped releasing. In practice an experiment's outcome is lost when its branch is: S-0176, an experiment in progress, would be refused at review, and nothing records what it found except its narrative. The designer decided on 2026-10-02 that experiments are accepted, and that their results are recorded in a document of their own.

## Decision

An `experiment` story is accepted like any other nature: merged, moved to done, archived, and committed. It records its results in a document under `design/experiments/`, one per experiment, named for the story, which says the hypothesis, the success measure, what was done, the results, and a recommendation (adopt, adapt, or drop). Acceptance refuses an experiment that has no such document. Publishing gives no component a bump on an experiment's account, whatever it touched, as for research (ADR-0025, ADR-0032). This supersedes ADR-0025's experiment clause; its research clause, as ADR-0032 refined it, stands.

## Consequences

- `flai accept` and the dashboard's acceptance stop refusing an experiment by nature, and refuse one without its results document instead; `--no-release` is no longer the way to land one.
- The release plan treats `experiment` as it treats `research`: it contributes nothing to a component's next release.
- `design/experiments/` is a new folder with a `README.md`, listed in the repository layout, in `CLAUDE.md`, and in the template.
- An experiment's code reaches `main` unreleased; the next story that delivers to the component releases it, as for research.
- `design/system/workflow.md`, `work-hierarchy.md`, `repository-layout.md`, and the user guide change with the code.

## Alternatives considered

- Keep experiments on their branch: their results are lost with the branch, and the operator cannot accept the story.
- Accept experiments and record results only in the narrative: narratives are archived with the story and are not where a reader looks for what the system has learned.
- Record results in `design/system/`: that folder says how the system is, and an experiment that is dropped does not describe it.
