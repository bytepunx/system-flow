---
id: T-2015
type: task
nature: improvement
title: An ADR refining ADR-0094 records that under publish the orchestrator starts the publish formula in place of release_publish, which stays for a project with none
status: backlog
parent: S-0423
owner: alex
created: 2026-10-21T13:26:22Z
updated: 2026-10-21T13:26:22Z
transitions: []
stream: S-0423
tags: [cli, template]
touches: [design/adrs, design/adrs/README.md]
---
# T-2015 An ADR refining ADR-0094 records that under publish the orchestrator starts the publish formula in place of release_publish, which stays for a project with none

## Work

ADR-0094 has the orchestrator publish "only through `release_publish`". This story changes that: under `publish`, when `release_evaluate` is met, the orchestrator starts the baseline `publish` formula with `flai formula run publish --ready`, and `release_publish` refuses while the project defines a `publish` formula, staying for a project with none. Write the ADR with `flai adr new`, refining ADR-0094 and ADR-0067: the decision, that each publish is then a story worked by a story's agent, measured, and archived, whose main-checkout steps run `flai release` with the host's credentials; the `push` host action as the gate it keeps; and the alternatives, among them keeping `release_publish` beside the formula. Record what TH-0583 settles on where a main-checkout step runs and where the `push` gate lives, once the operator has answered. Add it to `design/adrs/README.md`.

Waits for nothing: first layer, beside T-2013. The code and the role file follow it.

## Done when

`flai check --strict` reports nothing on the ADR, its `## Decision` opens with the decision sentence, and `design/adrs/README.md` lists it.

## Notes

Asked for by the orchestrator on TH-0583. The folder touch `design/adrs` is kept because the ADR's file name is known only once `flai adr new` makes it; the folder is in `claims.shared`, so it holds no story.
