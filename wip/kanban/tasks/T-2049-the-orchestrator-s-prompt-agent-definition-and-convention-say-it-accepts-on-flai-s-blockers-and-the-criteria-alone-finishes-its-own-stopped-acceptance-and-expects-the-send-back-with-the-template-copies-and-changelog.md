---
id: T-2049
type: task
nature: improvement
title: The orchestrator's prompt, agent definition, and convention say it accepts on flai's blockers and the criteria alone, finishes its own stopped acceptance, and expects the send-back, with the template copies and changelog
status: backlog
parent: S-0467
owner: alex
created: 2026-10-22T21:00:46Z
updated: 2026-10-22T21:00:46Z
transitions: []
stream: S-0467
tags: [flai, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-2039, T-2042, T-2047]
---
# T-2049 The orchestrator's prompt, agent definition, and convention say it accepts on flai's blockers and the criteria alone, finishes its own stopped acceptance, and expects the send-back, with the template copies and changelog

## Work

Waits for the three code tasks, so that the words describe what flai does.

Change the accept step of `orchestratePrompt` in `flai/internal/harness/harness.go`, the matching step of `.claude/agents/orchestrator.md`, and "As the orchestrator" in `design/conventions/strategic-agents.md`, keeping the three in step, so that each says:

- Accept on flai's blockers and the acceptance criteria alone. A concern beyond them, such as a design change the story was not asked for or a gate not yet built, is not a reason to leave the story in review: accept it, and raise the concern on a thread on the story, or as an issue, for the operator. A thread it opened that nobody else has written on does not hold the story; one the operator or an agent wrote on does, until it is resolved.
- When a thread it opened at step 4 says what was missing and the next pass finds it supplied, it accepts; the acceptance resolves its thread.
- A dirty main checkout is no reason to wait: `flai accept` leaves the checkout's uncommitted files alone and refuses only a change that would overwrite one, which the dry run lists as a blocker; that blocker, and a resumed acceptance the operator made, are the operator's to clear, on a thread. An acceptance of its own that stopped before its commit it finishes with `flai accept <S-nnnn> --by orchestrator`, and never with `--yes`.
- An acceptance whose rebase stops on a conflict flai cannot settle sends the story back to in-progress itself and opens the thread: the orchestrator logs it and opens no thread of its own.

Keep the queue variant of the prompt true: with the queue running it still opens no thread for a blocker the queue settles. Update `TestTheOrchestratorsPromptSaysWhatEachPermissionDoes` and the other prompt tests in `flai/internal/harness/harness_test.go` for the new words. Copy the agent definition and the convention's baseline to `template/root/.claude/agents/orchestrator.md` and `template/root/design/conventions/strategic-agents.md`, and add the change to `template/CHANGELOG.md`. A write under `.claude/` is asked of the operator through `permission_prompt`; make it yourself, not from a sub-agent.

## Done when

- The prompt, the agent definition, the convention, and their template copies say the four points above in the same terms, and `TestTheOrchestratorsPromptSaysWhatEachPermissionDoes` covers the new words.
- `template/CHANGELOG.md` records the change, and `flai test` on the changed paths and the markdown lint pass.

## Notes

Layer 4, beside the design-and-docs task; no path in common with it. The convention's `updated` date moves with the edit.
