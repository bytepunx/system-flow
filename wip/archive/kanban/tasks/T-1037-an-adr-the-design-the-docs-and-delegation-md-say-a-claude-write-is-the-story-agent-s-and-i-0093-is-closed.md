---
id: T-1037
type: task
nature: remediation
title: An ADR, the design, the docs, and delegation.md say a .claude/ write is the story agent's, and I-0093 is closed
status: done
parent: S-0299
owner: alex
created: 2026-10-06T21:08:27Z
updated: 2026-10-06T22:53:04Z
transitions:
  - to: ready
    at: 2026-10-06T22:39:18Z
    by: agent-S-0299
  - to: in-progress
    at: 2026-10-06T22:39:18Z
    by: agent-S-0299
  - to: done
    at: 2026-10-06T22:53:04Z
    by: agent-S-0299
stream: S-0299
tags: [docs, adr]
touches: [design/adrs, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/settings.md, design/conventions/delegation.md, design/system/strategic-agents.md, design/issues/I-0093-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md, design/issues/summary.md]
after: [T-1035, T-1036]
usage:
  source: log
  seconds: 826
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 70
      output: 23148
      cache_read: 3205245
      cache_write: 104086
      cost: 1.7557
---
# T-1037 An ADR, the design, the docs, and delegation.md say a .claude/ write is the story agent's, and I-0093 is closed

## Work

This task waits for T-1035 and T-1036 because it records what they built, and it closes I-0093 only once both are done.

- Record the decision with `flai adr new`: while `auto-approve` is off, `flai guard` refuses a sub-agent's write to a file in a `.claude/` folder, and the story's agent makes that write itself. It refines ADR-0060, under which the guard did not guard files. It also refines ADR-0086, whose `permission_prompt` still decides the story agent's own write. Give it the I-0093 instance as its context.
- In `design/system/flai-cli.md`, add the rule to the `flai guard` row and to what the `auto-approve` host action does, citing the new ADR.
- In `docs/users/flai.md` § Writes under .claude/, say that a story's sub-agents cannot make these writes while `auto-approve` is off, and that the story's agent makes them after the layer or hands them over on a thread with `cp` commands. Say the same where `docs/operators/settings.md` describes `auto-approve`.
- Rewrite the project addition in `design/conventions/delegation.md` that says to edit files under `.claude/` with `Edit` or `Write` like any other file. Make it say that only the story's agent does so, once the layer's sub-agents are back, and that a sub-agent returns the file's content instead. Leave the baseline above the marker as it is.
- Close the issue with `flai issue close I-0093 --reason`, naming T-1035's guard rule and T-1036's start prompt as the fix. Run it from the story's worktree so that the change is on the story branch.

## Done when

- The new ADR is accepted, listed in `design/adrs/README.md`, and names ADR-0060 and ADR-0086.
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/settings.md`, and `design/conventions/delegation.md` describe the guard's refusal and the story agent's route, and none of them still tells a sub-agent to write under `.claude/`.
- I-0093 is closed with a reason that names the fix, and `design/issues/summary.md` shows it closed.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

Drafted by the planner for S-0299. The ADR number is not known until `flai adr new` runs, hence the folder touch on `design/adrs`.
