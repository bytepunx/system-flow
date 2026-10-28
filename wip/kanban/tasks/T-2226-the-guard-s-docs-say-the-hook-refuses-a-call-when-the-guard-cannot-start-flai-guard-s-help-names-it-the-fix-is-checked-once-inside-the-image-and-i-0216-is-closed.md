---
id: T-2226
type: task
nature: remediation
title: The guard's docs say the hook refuses a call when the guard cannot start, flai guard's help names it, the fix is checked once inside the image, and I-0216 is closed
status: backlog
parent: S-0518
owner: alex
created: 2026-10-28T13:34:33Z
updated: 2026-10-28T13:34:33Z
transitions: []
stream: S-0518
tags: [cli, template]
touches: [docs/operators/settings.md, docs/users/flai.md, design/system/agent-adapters.md, flai/cmd/guard.go, docs/users/flai-reference.md, design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md, design/issues/summary.md]
after: [T-2223, T-2225]
---
# T-2226 The guard's docs say the hook refuses a call when the guard cannot start, flai guard's help names it, the fix is checked once inside the image, and I-0216 is closed

## Work

Describe the hook's new shape where the old one is described, linking the ADR T-2223 wrote. `docs/operators/settings.md`: the `hooks.PreToolUse` row says the hook refuses the call, with the reason, when `flai guard` cannot start (exit 126 or 127), and still passes a guard that ran and allowed or exited otherwise. `docs/users/flai.md` § Sub-agents, where it says the guard needs a harness that runs it before each call: what the agent sees when flai is not on `PATH`, and that inside the agent container the image's flai is the guard, this repository's through `scripts/flai.sh` (T-2224). `design/system/agent-adapters.md` § The guard: beside "It fails open: input it cannot read passes", the hook's refusal when the guard cannot start; and in § A container harness of flai's own, the sentence of "What the spike settled" that says this repository's hook let every call through gains that S-0518 closed it. `flai/cmd/guard.go`: one sentence in the command's long help on what the hook does with its exit codes, then `make flai-reference` to regenerate `docs/users/flai-reference.md`. Check the fix once inside `flai-agent-claude-code:local`, built from the branch with `make agent-images`: in the story's worktree run the hook command by hand through `scripts/agent-container.sh S-0518 -- sh -c '<the Bash hook command>'` with a sub-agent's call on standard input, as S-0505 did, and see it refuse with the guard's reason through the image's flai, and with `PATH` emptied see it refuse with `did not run (exit 127)`. Then close the issue: `flai issue close I-0216 --reason "..."`, the reason naming the ADR, the hook's refusal on exit 126 and 127, and `scripts/flai.sh`'s fallback to the installed flai; the command writes the issue file and `design/issues/summary.md`. Tick the story's criteria you verified with `flai criteria tick`. Waits for T-2223, whose ADR the docs link, and T-2225, whose hooks they describe.

## Done when

The three documents and the help describe the hook's refusal when the guard cannot start, `docs/users/flai-reference.md` is regenerated and `TestSettingsIndexIsComplete` and the reference test pass, the check inside the image is recorded in the narrative with the commands and their exit codes, I-0216 is closed with its reason, `design/issues/summary.md` no longer lists it, and `flai check --strict` is clean.

## Notes

`docs/users/flai.md` and `docs/operators/settings.md` are touched by stories in progress on the day of planning (S-0414, S-0502, S-0521); the sync before this task brings their changes if they land first.
