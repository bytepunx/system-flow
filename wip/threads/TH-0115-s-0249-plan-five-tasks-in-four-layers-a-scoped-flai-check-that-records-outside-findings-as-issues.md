---
id: TH-0115
title: "S-0249 plan: five tasks in four layers, a scoped flai check that records outside findings as issues"
anchor:
  path: wip/kanban/stories/S-0249-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md
  item: S-0249
status: resolved
participants: [planner-S-0249, alex]
created: 2026-10-05T00:19:57Z
updated: 2026-10-05T00:23:18Z
---

# TH-0115 S-0249 plan: five tasks in four layers, a scoped flai check that records outside findings as issues

On wip/kanban/stories/S-0249-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md.

## Entries

### 2026-10-05T00:19:57Z planner-S-0249
**Solution proposed from I-0057's 19 instances and your answer on TH-0056.** `flai check --story S-nnnn` reports findings outside the story as notes that do not fail the run. `--record-issues` opens or bumps one issue per rule for them. The close-out scripts, here and in the template, pass both flags. The unscoped check, as CI runs it, still gates everything.

Tasks and layers:

1. Layer 1: **T-0840**, the `--story` scope. A finding is the story's when its path is the story's file or its tasks' files, its narrative, its threads, or a path in its diff against main. The task includes the reproduction test.
2. Layer 2: **T-0841**, `--record-issues`, which opens or bumps one issue per rule and writes each story's instance once. It waits for T-0840: it reads T-0840's outside marker, and both change `flai/cmd/check.go`.
3. Layer 3, two tasks that run together because they share no path. Both wait for T-0841.
   - **T-0844**: `close-out.sh` exports the story, `check.sh` passes the flags, and `TestMonorepoIsClean` is scoped the same way. Applies here and in the template.
   - **T-0845**: an ADR refining ADR-0073, `flai-cli.md`, `continuous-improvement.md`, `docs/users/flai.md`, both copies of `work-management.md`, and `template/CHANGELOG.md`.
4. Layer 4: **T-0847**, a close-out rehearsal past main's findings, then closing I-0057. It waits for T-0844 and T-0845.

Figures:

- Forecast 1h, delivery 2026-10-05T01:45Z. flai gave 30m; I raised it because of the five tasks across the template, docs, and an ADR.
- Cost of delay 17.50 USD/week, from your 7m on TH-0112.
- 19 touches, with tags `flai` and `template` and topics `cli`, `conventions`, and `template`.

Assumptions (correct any of them here):

- Errors outside the story are notes too, not only warnings. S-0200 and S-0203 stopped on main's `issues.duplicate-id`, which no story could clear.
- A `wip.overlap` with another open story always counts as outside the story, even when it is reported on the story's own file. This follows remediation 1 of I-0059/S-0244.
- Out-of-story findings get one issue per rule, such as "flai check finds `story.unaccepted` outside the story at close-out". The alternative was one issue per rule and path, which would make many issues out of `wip.overlap` pairs.
- The issues are written in the story's worktree, and the close-out commits them with the story.

Related stories (no change made):

- **S-0252**, in progress, shares `flai/internal/issues` and `design/issues`. The touches hold keeps S-0249 behind it, which also keeps issue numbering safe across branches. I did not add an `after`; say if you want one.
- **S-0244**: S-0249 delivers its remediation 1, so the close-out no longer fails on `wip.overlap`. I propose narrowing S-0244 to its remediations 2–5 once S-0249 is accepted. I have not edited it.
- **S-0266 and S-0267** also change the close-out and verifier flow. Their touches will tell whether they collide with T-0844.

### 2026-10-05T00:23:18Z alex
Resolved.
