---
id: TH-0119
title: "S-0268 plan: three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0268-the-story-agent-s-prompt-says-to-batch-independent-edits-and-commands-in-one-turn-to-take-a-new-task-through-ready-in-one-command-and-to-prove-a-test-fails-without-the-shared-stash.md
  item: S-0268
status: resolved
participants: [planner-S-0268, alex]
created: 2026-10-05T00:30:19Z
updated: 2026-10-05T01:34:34Z
---

# TH-0119 S-0268 plan: three tasks in two layers

On wip/kanban/stories/S-0268-the-story-agent-s-prompt-says-to-batch-independent-edits-and-commands-in-one-turn-to-take-a-new-task-through-ready-in-one-command-and-to-prove-a-test-fails-without-the-shared-stash.md.

## Entries

### 2026-10-05T00:30:19Z planner-S-0268
**The plan for S-0268.** The cost of delay question on TH-0117 is still open.

- **Touches (8):**
  - `flai/internal/harness/harness.go` and `flai/internal/harness/harness_test.go`
  - `design/conventions/git.md` and `delegation.md`, with their template copies
  - `design/system/flai-cli.md` and `docs/operators/index.md`
  - `### Planning` in the Notes says where each one came from.
- **Tags and topics:** tags `[flai, template]`, topics `[cli, conventions]`.
- **Forecast:** 25m, delivery 01:47Z.
  - `flai forecast` gave 16m. I raised it for two close-out runs of the flai tier.
  - The delivery counts from S-0266's forecast acceptance at 01:22Z. S-0268 shares four paths with S-0266, so it is held until then.

**Tasks**, in two layers:

1. **Layer 1**, which can run together because the two tasks share no path:
   - **T-0853**: the prompt's two sentences in `harness.go`, asserted in `harness_test.go`.
   - **T-0854**: `git.md` gets the rule on showing a test fails, and the ban on `git stash`. `delegation.md` gets the batching and the one-command move. Both have template copies.
2. **Layer 2:**
   - **T-0855**: `flai-cli.md` and the operator docs record both. It waits for T-0853 and T-0854.

**Assumptions:**

- "In one command" means one shell call that chains two moves (`flai move T-nnnn ready && flai move T-nnnn in-progress`), not a change to `flai move`. If you want `flai move` to take a backlog task straight to in-progress, that is a CLI change, and I would make it a story of its own.
- `delegation.md` gets the batching sentences as well as the prompt, because the goal names it and an agent started without flai serve's prompt reads only the conventions. The criteria require only the prompt.
- `code-quality.md` already says a change gets a test that fails without it. I left it alone, and `git.md` says how to show that, as the criterion asks.
- `template/CHANGELOG.md` and `template/template.yaml` are left to the release tooling, as for S-0266.

### 2026-10-05T01:34:34Z alex
Resolved.
