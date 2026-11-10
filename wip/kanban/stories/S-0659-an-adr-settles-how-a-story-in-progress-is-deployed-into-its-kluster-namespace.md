---
id: S-0659
type: story
nature: feature
title: "An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once"
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:07Z
updated: 2026-11-10T08:09:07Z
transitions: []
tags: [cli]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once

## Goal

An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request, so that the stories after it build one design. It amends [ADR-0254](../../../design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md), and [kluster.md](../../../design/system/kluster.md) says the same.

## Acceptance criteria

- [ ] An accepted ADR says how a story's agent asks for a deploy and a teardown of its own story's namespace, that flai serve runs the job, and that the `kluster` host action gates it
- [ ] It says what the agent may run against the namespace: the project's checks, its own commands with the namespace's variables, and what `flai guard` refuses, another story's namespace and the default one among them
- [ ] It says when a namespace is made, deployed again, and torn down across in-progress, review, a send back, acceptance, and cancellation, in a table that replaces ADR-0254 §4's
- [ ] It sets how many story namespaces may live at once, where the project sets it, and what a request over the limit gets
- [ ] It says how the agent waits for a job behind the queue, what it is told on a failure, and when a failure reaches the operator
- [ ] It says what is recorded of a deploy and a checks run before review, and where the move to review, the close-out, and an acceptance read it
- [ ] ADR-0254 names it in `superseded_by` or the new ADR names the sections it amends, as `design/conventions/decisions.md` has it, and `design/system/kluster.md` describes the design as decided and marks each part not yet built

## Tasks

## Notes

- A design story: it changes no code.
- The hand-off `flai verify --job` makes to flai serve is the pattern to weigh first.
