---
id: T-2575
type: task
nature: remediation
title: Say in the operators guide and flai-cli.md that a story already in ready gets its agent once the agent action is turned on
status: backlog
parent: S-0589
owner: alex
created: 2026-11-07T10:06:34Z
updated: 2026-11-07T10:06:34Z
transitions: []
stream: S-0589
tags: [docs, serve]
touches: [docs/operators/index.md, design/system/flai-cli.md]
after: [T-2574]
---
# T-2575 Say in the operators guide and flai-cli.md that a story already in ready gets its agent once the agent action is turned on

## Work

Bring the documentation in step with T-2574's fix, in the same story as `documentation.md` asks. It waits for T-2574 because what the documents say, at once or within a minute, depends on the mechanism that task chose.

- `docs/operators/index.md` § "Starting an agent when a story becomes ready": say that a story already in ready when the operator turns the `agent` action on, as after the board's Import, gets its agent on the look that follows, with no Start agent needed.
- `design/system/flai-cli.md`, the `flai serve agent` row of the commands table: name what makes the launcher look, adding the turn-on of the `agent` action, and say it with S-0589.

Change no ADR: ADR-0041 and ADR-0043 already decide this behaviour.

## Done when

- Both documents say what happens to a story already in ready when the `agent` action is turned on, as T-2574 built it.
- The markdown lint passes on both files.

## Notes
