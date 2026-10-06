---
id: TH-0216
title: "S-0270 plan: seven tasks in four layers on S-0273's verify package, and three assumptions to confirm"
anchor:
  path: wip/kanban/stories/S-0270-verification-is-a-module-flai-verify-runs-the-tiers-the-diff-selects-and-answers-structured-findings-replacing-the-verifier-sub-agent-s-run.md
  item: S-0270
status: resolved
participants: [planner-S-0270, alex]
created: 2026-10-06T22:54:34Z
updated: 2026-10-06T23:32:31Z
---

# TH-0216 S-0270 plan: seven tasks in four layers on S-0273's verify package, and three assumptions to confirm

On wip/kanban/stories/S-0270-verification-is-a-module-flai-verify-runs-the-tiers-the-diff-selects-and-answers-structured-findings-replacing-the-verifier-sub-agent-s-run.md.

## Entries

### 2026-10-06T22:54:34Z planner-S-0270
I've planned S-0270: seven tasks in four layers, on top of S-0273's `flai/internal/verify` package.

The tasks, by layer:

1. **Layer 1:** T-1071. `Verify` runs the close-out's checks in order and stops at the first failure: the unfinished rebase, the tiers the diff selects, the scoped `flai check --strict` (outside findings are notes), the narrative, and contains-main. It answers each step's state, duration, and findings, and stores the last result in `.flai-cache/verify/<story>.json`.
2. **Layer 2, after T-1071:**
   - T-1075 adds the `flai verify S-nnnn` command, with text and `--json`, an exit status that follows the outcome, and `--record-issues`.
   - T-1076 adds the MCP tool `verify` and the host methods `verify.run` and `verify.status`.
3. **Layer 3:**
   - T-1082 (after T-1075): both close-out scripts call `flai verify`. They keep the commit, the sync check after it, and the clean-worktree check.
   - T-1089 (after T-1076): `Review.svelte` shows the last result, through a new `/api/items/[id]/verify` route.
   - T-1096 (after T-1075): a new ADR refining ADR-0059, plus the story agent's and orchestrator's prompts, `delegation.md`, `work-management.md`, both `verifier.md`, and the template copies with a CHANGELOG line.
4. **Layer 4:** T-1108 (after all the others) updates `flai-cli.md`, `devex.md`, `dashboard-host-channel.md`, the user guides, and the regenerated reference.

The figures:

- **Forecast:** 90m, against flai's 50m. Delivery is 2026-10-07T10:05Z, which is S-0273's forecast delivery plus 90m.
- **Cost of delay:** 43 USD a week stands, against flai's 164.92.

The reasons are in the story's `### Planning`.

The assumptions, each with my recommended answer. Reply only to change one:

1. S-0273's tier declaration carries the paths each tier applies to, so `flai verify` can pick tiers by diff, as `close-out.sh` picks the flai, template, and flaiover tiers today. If S-0273 lands without that, T-1071 adds it.
2. `verify.run` sits behind the existing `checks` host action, not a new one.
3. "The dashboard's story page" in criterion 2 means the review page `/review/[id]`, where `Review.svelte` already shows checks. It shows the result only, with no run button.

One proposal: drop the declared touches `flai/internal/serve/checks.go` and `checks_config.go`. No task changes them, because `flai verify` runs S-0273's tiers, not the review checks' commands. While S-0270 is in progress they hold any ready story that touches those files. I kept them, because they were declared. Reply `drop` and I'll remove them.

Two other points:

- S-0273 is still a draft, and S-0270 waits for it.
- T-1096's `.claude/agents/verifier.md` edit will ask you to allow it on a thread.

### 2026-10-06T23:32:31Z alex
Resolved.
