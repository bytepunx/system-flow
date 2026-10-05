---
id: T-0870
type: task
nature: remediation
title: The design and the user docs describe the grown-claim notice and the longer wait, and I-0059 is closed
status: in-progress
parent: S-0244
owner: alex
created: 2026-10-05T03:16:19Z
updated: 2026-10-05T05:02:41Z
transitions:
  - to: ready
    at: 2026-10-05T05:02:41Z
    by: agent-S-0244
  - to: in-progress
    at: 2026-10-05T05:02:41Z
    by: agent-S-0244
stream: S-0244
tags: [docs, design]
touches: [design/system/agent-narrative.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md, design/issues/summary.md]
after: [T-0867, T-0868, T-0869]
---
# T-0870 The design and the user docs describe the grown-claim notice and the longer wait, and I-0059 is closed

## Work

- Extend the `overlapped` notice text with the notice T-0867 adds: its cause, its summary, and what the agent does. The text is in `design/system/workflow.md` (the S-0132 paragraph under Branches and collisions), `design/system/agent-narrative.md` (What an agent is told through MCP), and `design/system/flai-cli.md` (the rows for `flai task new`, `flai edit`, and `flai touches`).
- State the wait cap and the defaults T-0868 settled, and why. Put them in `agent-narrative.md` beside `wait_for_events`, and in `docs/users/flai.md` (the waits paragraph and the `wait_for_events` row) and `docs/users/flai-reference.md` where they list the tools.
- Close I-0059 with `flai issue close I-0059 --reason`, run in the story's worktree. Name what fixed each remediation: 1 by S-0249's scoped close-out, 2 by T-0867, 3 by T-0868, and 4 by T-0869. Say that 5 is not needed: once the close-out no longer stops on `wip.overlap`, two stories can no longer each wait for the other. Keep `design/issues/summary.md` current.

It waits for the other three, since it describes what they built and closes the issue on them.

## Done when

- [ ] The design and the user docs name the grown-claim `overlapped` notice and the wait cap and defaults, consistent with the code
- [ ] I-0059 is closed with a reason naming each remediation's fix, and `summary.md` agrees
- [ ] The markdown lint and `flai check --strict` scoped to S-0244 pass

## Notes
