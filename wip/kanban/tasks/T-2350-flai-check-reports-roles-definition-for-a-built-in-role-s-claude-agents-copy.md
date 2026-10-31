---
id: T-2350
type: task
nature: improvement
title: flai check reports roles.definition for a built-in role's .claude/agents copy that is neither the ADR's link nor the definition its role file generates
status: backlog
parent: S-0553
owner: alex
created: 2026-10-30T20:47:22Z
updated: 2026-10-30T20:47:22Z
transitions: []
stream: S-0553
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/roles_test.go]
after: [T-2347, T-2348]
---
# T-2350 flai check reports roles.definition for a built-in role's .claude/agents copy that is neither the ADR's link nor the definition its role file generates

## Work

In `flai/internal/check/check.go`, beside `roles.baseline`, compare each built-in role's `.claude/agents/<name>.md`, when it exists, with what the ADR accepts: the link it names, or a file equal to the definition generated from the role file. Report `roles.definition` naming the file and the first field that differs. A missing file is not a finding.

Second layer: it waits for the generation task, whose definition it compares, and for the ADR, which says what a copy may be.

## Done when

`roles_test.go` covers a matching copy, a drifted description, a drifted prompt, the accepted link, and a missing file. `flai check --strict` on this repository reports what the `.claude/agents` pair task then clears.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. No code compares role files with `.claude/agents` today; D-0411 has kept them in step by hand.
