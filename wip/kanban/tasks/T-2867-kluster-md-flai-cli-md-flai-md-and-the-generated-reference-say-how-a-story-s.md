---
id: T-2867
type: task
nature: feature
title: kluster.md, flai-cli.md, flai.md, and the generated reference say how a story's agent checks its namespace and runs its own commands there
status: in-progress
parent: S-0661
owner: alex
created: 2026-11-10T08:24:25Z
updated: 2026-11-13T13:14:24Z
transitions:
  - to: ready
    at: 2026-11-13T13:14:24Z
    by: agent-S-0661
  - to: in-progress
    at: 2026-11-13T13:14:24Z
    by: agent-S-0661
stream: S-0661
tags: [cli, docs]
touches: [design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-2863, T-2864, T-2865, T-2866]
---
# T-2867 kluster.md, flai-cli.md, flai.md, and the generated reference say how a story's agent checks its namespace and runs its own commands there

## Work

Write down what T-2863 to T-2866 built, as built:

- `design/system/kluster.md` § What the review checks see: the checks a story's agent asks for, the refusal when the namespace is not deployed at the head, the queue, no thread on a failure, and the record with who asked. § Commands: `flai kluster env` (or the form T-2865 built). Add the guard's namespace rule beside the line on `klusterWrites`. Mark built each part the ADR S-0659 accepted listed as not yet built.
- `design/system/flai-cli.md` § Commands: the new command, the `--job` flag, and the guard's rule.
- `docs/users/flai.md`: the MCP tool's row in the table under Serving agents over MCP, and a paragraph under Run a project's kluster on checking a story in progress and running a command with its variables.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference`, never by hand.

It waits for every other task, since it describes what they built.

## Done when

- The four files say all of the story's criteria 1 to 7 as built, and `make flai-reference` leaves `docs/users/flai-reference.md` unchanged.
- `flai check --strict` and the markdown lint pass on them.

## Notes
