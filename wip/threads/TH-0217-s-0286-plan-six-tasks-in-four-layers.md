---
id: TH-0217
title: "S-0286 plan: six tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0286-a-story-that-changes-a-path-claude-code-protects-is-flagged-at-review-and-accepted-only-by-the-operator-flai-checks-permission-prompt-against-each-new-claude-code-and-the-prompt-covers-every-protected-path-in-a-worktree-but-git.md
  item: S-0286
status: open
participants: [planner-S-0286]
created: 2026-10-06T22:55:25Z
updated: 2026-10-06T22:55:25Z
---

# TH-0217 S-0286 plan: six tasks in four layers

On wip/kanban/stories/S-0286-a-story-that-changes-a-path-claude-code-protects-is-flagged-at-review-and-accepted-only-by-the-operator-flai-checks-permission-prompt-against-each-new-claude-code-and-the-prompt-covers-every-protected-path-in-a-worktree-but-git.md.

## Entries

### 2026-10-06T22:55:25Z planner-S-0286
The plan for S-0286. The figures and where each touch came from are in the story's `### Planning` notes.

Tasks and layers:

1. **Layer 1:** T-1106, the ADR refining ADR-0086. Before the ADR relies on the claim that the running agent reads its configuration only from the main checkout, the task tests it with a hook changed in a worktree.
2. **Layer 2:** these two run together, since they share no file.
   - T-1110: the shared list of protected paths, and `permission_prompt` widened to it, still refusing `.git`.
   - T-1114: flai serve's check of each new Claude Code version.
3. **Layer 3:** T-1120, the acceptance gate. It waits for T-1110's list. The preview lists the protected files, and an acceptance by the orchestrator or by an agent on its own name is refused with the files named.
4. **Layer 4:**
   - T-1121, the review page. It waits for T-1120.
   - T-1122, the design, the guides, and the auto-approve description in `flai serve actions`. It waits for T-1114 and T-1120.

Figures:

- **Forecast:** 1h, delivery 2026-10-07T10:48Z. flai gave 43m; I raised it by 17m for the live hook test and the new check code.
- **Cost of delay:** 150 USD a week, from your 1h-per-cycle input on TH-0206.
- **Topics:** added `dashboard`.

Assumptions:

- The list of protected paths goes in a new leaf package, `flai/internal/protected`, because `mcpserver` does not import `preview`. The story's agent may place it elsewhere, as long as one list serves both.
- The dashboard's `accept.run` builds its refusal from the same preview as `flai accept`. So `hostapi/writes.go` changes only for the auto-approve description.
- S-0221 is done, so the orchestrator's acceptance path exists, and T-1120 tests it.
- T-1120 tells an agent's name from yours the way flai already does. If flai has no such way for a story's agent, T-1120 asks on this story before choosing one.
- The check's scratch project is a flai project with one in-progress story and its worktree. That is because `permission_prompt` allows only there.
- I kept the declared folder touches `flai/internal/serve`, `flai/internal/harness`, and `docs/operators`. The tasks name their files, and ADR-0096 narrows the claim to them. `design/adrs` stays a folder because the ADR's number is not known yet.

I have nothing to propose splitting or dropping.
