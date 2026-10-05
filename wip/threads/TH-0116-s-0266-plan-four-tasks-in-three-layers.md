---
id: TH-0116
title: "S-0266 plan: four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0266-the-verifier-reads-a-close-out-run-s-exit-status-without-piping-it-away-and-the-story-agent-tells-it-where-the-run-is-expected-to-stop.md
  item: S-0266
status: resolved
participants: [planner-S-0266, alex]
created: 2026-10-05T00:20:31Z
updated: 2026-10-05T00:22:32Z
---

# TH-0116 S-0266 plan: four tasks in three layers

On wip/kanban/stories/S-0266-the-verifier-reads-a-close-out-run-s-exit-status-without-piping-it-away-and-the-story-agent-tells-it-where-the-run-is-expected-to-stop.md.

## Entries

### 2026-10-05T00:20:31Z planner-S-0266
The plan for S-0266 is four tasks in three layers, all in the backlog.

Layer 1:

- **T-0849:** `scripts/close-out.sh` ends every run, passing or stopping, with one line naming the story, the outcome, and the step. The step is set before each step runs and printed from an EXIT trap. The template's `close-out.sh` gets the same line, and both READMEs are updated.

Layer 2, where T-0850 and T-0851 wait for T-0849, share no path, and run together:

- **T-0850:** the verifier's definition, in the template first and then copied to `.claude/agents/`. It runs a long script once, in one Bash call with `; echo "exit $?"` after it, with no pipe, no file, and the 600000 ms timeout. It reads the last line. On an expected stop it reports it and runs the steps after it through their own entry points.
- **T-0851:** the story agent's prompt (`delegation()` in `harness.go`). It names the steps it knows will stop, and why, when it hands the close-out to the verifier. The prompt's tests and `flai-cli.md`'s description of the prompt change with it.

Layer 3:

- **T-0852:** the baseline of `delegation.md` and the template's copy, `design/system/workflow.md` where close-out is described, and `docs/users/flai.md`. It waits for T-0850 and T-0851 so that it repeats their words.

Assumptions:

1. The template's `close-out.sh` gets the line too, though the criterion names only `scripts/close-out.sh`. The template's verifier definition will tell every project's verifier to read that line.
2. `workflow.md` is the design that records the change. The criterion allows "or where close-out is described", and S-0267 is rewriting the test tiers in `devex.md`.
3. Writing `.claude/agents/verifier.md` will likely be refused to the story's agent (I-0069). T-0850 tells it to ask you on a thread with the whole file to paste.
4. T-0850 asks the agent to check that the last line survives the Bash tool's output limit on a full flai-tier run.
5. I tagged the story `[flai, template]` because `flai check --strict` asked for a component tag.

Overlaps, which the board's holds will sequence:

- **S-0267**, in progress, touches `scripts/close-out.sh` and `scripts/README.md`.
- **S-0249** touches both `close-out.sh` scripts, `flai-cli.md`, and `docs/users/flai.md`. It also adds a step to the close-out.
- **S-0268** will change `delegation()` and `delegation.md` too, though it declares no touches yet.

Proposal: give S-0266 `after: [S-0249]`, so that T-0849's step line names S-0249's scoped check, and T-0851's example of an expected stop is still true. I have not set it, because the order is yours. Reply "set it" and I will.

Figures:

- Forecast: 40m, up from flai's 26m, for the two flai-tier close-out runs and the likely `.claude/` round trip. Delivery is about 01:20Z.
- Cost of delay: waits on TH-0113.

I propose no splits, merges, or drops of tasks.

### 2026-10-05T00:22:32Z alex
Resolved.
