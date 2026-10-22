---
id: T-2050
type: task
nature: improvement
title: The conventions, the agents' prompts, the user docs, and the design describe the detached close-out and the calls that follow it, in place of one 600000 ms call
status: backlog
parent: S-0466
owner: alex
created: 2026-10-22T21:00:52Z
updated: 2026-10-22T21:00:52Z
transitions: []
stream: S-0466
tags: [cli, template]
touches: [design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/verifier.md", template/root/.claude/agents/verifier.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/workflow.md]
after: [T-2033, T-2046]
---
# T-2050 The conventions, the agents' prompts, the user docs, and the design describe the detached close-out and the calls that follow it, in place of one 600000 ms call

## Work

Write what T-2031, T-2033, and T-2046 built into every place that tells an agent or a reader how the close-out runs, in both trees where a file has a template copy:

- `design/conventions/delegation.md`, the close-out bullets under `## Rules` and the one under `## As an explorer or a verifier`, which say to run it "as one command with the longest timeout the harness allows (Claude Code's Bash tool: 600000 ms)": say instead that `scripts/close-out.sh S-nnnn -m "<message>"` starts the close-out detached and follows it for up to its wait, and that the agent calls `flai close-out S-nnnn` again, as often as its last line says the run is still going, until that line names the outcome and the step; and `template/root/design/conventions/delegation.md` the same.
- `design/conventions/work-management.md`, the definition-of-done bullet and the "Before review" bullet that describe what the close-out does and when it stops: add how it is followed; and its template copy.
- `flai/internal/harness/harness.go`, the delegation prompt that repeats the 600000 ms call, and `flai/internal/harness/harness_test.go`, which asserts its words.
- `.claude/agents/verifier.md` and `template/root/.claude/agents/verifier.md`, step 4, which names the 600000 ms call. A write under `.claude/` asks the operator through `permission_prompt`; make it yourself, not through a sub-agent (ADR-0102).
- `docs/users/flai.md`: a section for `flai close-out` beside § Verify a story before review, with its flags, the state and log under `.flai-cache/close-out/`, the still-running line, and a second call attaching; and the close-out sentences near the `flai task done` section and under § Resuming across a clean rebase. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- `design/system/flai-cli.md`, a row for `flai close-out` in the `## Commands` table, naming the ADR and this story; `design/system/workflow.md` § Branches and collisions, where the close-out's rounds are described.

Waits for T-2033 and T-2046, which settle the command's behaviour and the scripts' words.

## Done when

- No file under `design/conventions/`, `template/root/design/conventions/`, `.claude/agents/`, `template/root/.claude/agents/`, or `flai/internal/harness/` tells an agent to run the close-out as one 600000 ms call; each tells it to call `flai close-out` again until the last line names the outcome.
- `docs/users/flai.md` and `docs/users/flai-reference.md` document `flai close-out` and its flags, `design/system/flai-cli.md` lists it, and `design/system/workflow.md` says the close-out runs detached.
- `flai test` on the changed paths, the markdown lint among its tiers, passes; `flai check --strict` is clean; the template tier of `make flai-test` passes.

## Notes

The fourth layer. The `docs/users/flai-reference.md` is generated; do not edit it by hand.
