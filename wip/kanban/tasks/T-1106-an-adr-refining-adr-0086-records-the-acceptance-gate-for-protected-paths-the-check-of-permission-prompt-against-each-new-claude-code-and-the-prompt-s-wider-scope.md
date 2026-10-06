---
id: T-1106
type: task
nature: improvement
title: An ADR refining ADR-0086 records the acceptance gate for protected paths, the check of permission_prompt against each new Claude Code, and the prompt's wider scope
status: backlog
parent: S-0286
owner: alex
created: 2026-10-06T22:53:24Z
updated: 2026-10-06T22:53:24Z
transitions: []
stream: S-0286
tags: [flai]
touches: [design/adrs]
---
# T-1106 An ADR refining ADR-0086 records the acceptance gate for protected paths, the check of permission_prompt against each new Claude Code, and the prompt's wider scope

## Work

Write the ADR with `flai adr new`. It refines ADR-0086 and records the story's three parts:

- **Acceptance is the gate.** A story whose branch changes a path Claude Code protects is accepted by the operator only. Say why: the agent runs in the main checkout and reads its hooks, settings, and agent definitions from there, so a write in a worktree changes nothing it obeys until the story is accepted.
- **The check.** When flai serve's Claude Code has a version flai has not checked, flai makes one write through `permission_prompt` in a scratch project. It records the version and the outcome, and opens a thread to the operator on a failure.
- **The scope.** `permission_prompt` handles every path Claude Code protects inside an in-progress story's worktree, except `.git`.

Name the protected paths as Claude Code's documentation lists them under "Protected paths" on the permission modes page. Name what is left out, and why.

The claim that the running agent reads its configuration only from the main checkout was never tested. Before the ADR relies on it, test it: change a hook in a scratch worktree under a headless `claude -p` started in the main checkout, and see whether the hook takes effect. Record the result and the Claude Code version in the ADR. If the claim fails, say so on a thread on S-0286 before going on.

Record the alternatives the story's Notes reject: `bypassPermissions`, moving the template's files out of `.claude/`, and a flai tool that writes the files itself.

No `after`: this task is the first layer, and every other task builds on its decision.

## Done when

- A proposed ADR under `design/adrs` refines ADR-0086 and records the three parts, the list of protected paths, and why acceptance is the gate.
- The ADR records the result of the worktree hook test and the Claude Code version it ran against.
- `flai check --strict` is clean.

## Notes

Drafted by planner-S-0286. The touch is the `design/adrs` folder because `flai adr new` assigns the ADR's number when it writes the file.
