---
id: ADR-0156
title: "A close-out that main outran syncs and verifies again, at most three rounds, and flai verify reuses across that clean rebase the tiers main's new paths do not select"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0110, ADR-0135]
---

# ADR-0156 A close-out that main outran syncs and verifies again, at most three rounds, and flai verify reuses across that clean rebase the tiers main's new paths do not select

## Context

The close-out runs `flai verify <story> --record-issues`, commits what is outstanding, and then checks with `flai verify <story> --sync-only` that the branch still contains the main branch ([ADR-0110](0110-a-story-s-agent-runs-the-close-out-which-runs-flai-verify-itself-before-review.md)). That check passes over commits that change only `wip/` paths the branch does not change ([ADR-0135](0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md)). An acceptance or a release during the run still fails it.

I-0136 has eight instances, from the close-outs of S-0232 (twice), S-0337 (twice), S-0373, S-0360, S-0352, S-0379, and S-0378, on 2026-10-12 and 2026-10-13. In each, every tier passed, an acceptance landed on main while the tiers ran, and the last check failed. Each branch then rebased onto main cleanly. Each failure meant that the agent ran `flai stream sync` and then the whole close-out again, about six minutes of tiers each time. `flai verify` resumes only at an unchanged head and base commit (S-0341). A rebase changes both, so every tier ran again.

## Decision

When main gains commits while a close-out runs, the close-out syncs and verifies again on its own, at most three rounds, and `flai verify` reuses across that clean rebase the tiers the branch's unchanged patch passed, but for those main's new paths select.

1. **The close-out syncs and verifies again.** When its last check, `flai verify <story> --sync-only`, fails, `scripts/close-out.sh` runs `flai stream sync <story>` itself. When the sync ends clean, the close-out runs `flai verify <story> --record-issues`, the commit step, and the last check again. It runs at most three rounds, then stops at the last check as it does today, with the same hint. A sync that stops on a conflict, or a verify that fails, stops the close-out at that step as today. Each round is printed, and the last line names the round the close-out ended on. Both the repository's script and the template's do this.
2. **`flai verify` reuses tiers across a clean rebase.** The report records the patch id of the branch's own change, `git patch-id --stable` of its diff from the merge base with the main branch. It records the merge base too. The patch id is recorded only when the worktree has nothing uncommitted. A later run may then reuse the tiers that report passed, even though its head and base commit differ, when three things hold:
   - its patch id is the same;
   - its fingerprint is the same;
   - the report lists the same tiers in the same order.

   The paths changed between the report's merge base and this run's are main's new paths. A tier that those paths select, through the same selection the branch's own paths go through (`{importers}` included), runs again. Every tier after it runs again too, as resume does today. The `rebase`, `sync`, `narrative`, and `check` steps always run. `--fresh` reuses nothing, as today. The step's text names the tiers reused across the rebase and the tiers that main's paths sent back.
3. **Main's commits are not passed over.** A change on main outside `wip/` can break the branch without touching any of its files, for example a Go package the branch imports. The run's verdict has to hold for the head that goes to review. So the tiers that main's change can reach run again on that head, and only the tiers it cannot reach are reused.

This refines ADR-0135, whose last check stops the close-out on any other commit the branch lacks, and ADR-0110, whose close-out the agent ran again by hand.

## Consequences

- An acceptance during a close-out no longer needs the agent: the close-out syncs, verifies, and checks again on its own. It stops only after three rounds, on a conflict, or on a failing step.
- How much a round saves depends on what main changed and which tiers select it. In this repository the integration and smoke tiers select `flai/**`. An acceptance that changes flai therefore runs them again, and saves only the agent's round trip. An acceptance that changes no tier's paths, or only a cheap tier's, reuses everything else.
- When a round's commit step commits issues that `--record-issues` recorded, the branch's patch changes. The next round then reuses nothing. A report taken with uncommitted changes records no patch id, so the next round reuses nothing either.
- The template's `scripts/close-out.sh` loops with any flai that has `--sync-only`. It reuses tiers across the rebase only with a flai that records the patch id, and its changelog names that version. With an older flai every round runs every tier.
- The definition of done and the close-out rule in `work-management.md` say what the close-out now does on its own, and when it still stops.

## Alternatives considered

- **Pass over commits on main whose paths the branch does not change.** I-0136's first instance suggests it. A change elsewhere can break the branch, as in decision 3. ADR-0135 kept the same safeguard when it widened the step to `wip/`-only commits.
- **Hold acceptances while a close-out runs.** That delays the operator's acceptances for minutes, and couples every acceptance to every story's close-out.
- **Compare the branch's own change path by path, rather than by one patch id.** That would let a round reuse tiers after its commit step committed recorded issues. The operator confirmed the patch-id gate on TH-0424, and the gate is simpler to reason about.
- **Sync and run everything again in the close-out, without reuse.** It removes the agent's round trip but runs every tier each round, so it does not cut the cost I-0136 measures.
