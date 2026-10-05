---
id: S-0266
type: story
nature: improvement
title: The verifier reads a close-out run's exit status without piping it away, and the story agent tells it where the run is expected to stop
status: ready
owner: alex
created: 2026-10-05T00:06:34Z
updated: 2026-10-05T00:34:46Z
transitions:
  - to: ready
    at: 2026-10-05T00:13:54Z
    by: alex
tags: [flai, template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md, ".claude/agents/verifier.md", template/root/.claude/agents/verifier.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/flai-cli.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/system/workflow.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: planner-S-0266
    at: 2026-10-05T00:22:37Z
  value: 15
  by: planner-S-0266
  at: 2026-10-05T00:22:55Z
forecast:
  duration: 40m
  delivery: 2026-10-05T01:20:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 3, behind S-0249."
  by: flai
  at: 2026-10-05T00:34:46Z
---
# S-0266 The verifier reads a close-out run's exit status without piping it away, and the story agent tells it where the run is expected to stop

## Goal

S-0248's first verifier ran `scripts/close-out.sh` three times, about two minutes each: the first run piped the output through `tail`, which lost the script's exit status; the second printed `$?` after the pipe, which was tail's; the third redirected the output into a file, which the verifier's own rules forbid. The verifier's definition (`.claude/agents/verifier.md`, and its copy in `template/root/.claude/agents/verifier.md`) bans writing files but does not say how to read a long script's status, and `scripts/close-out.sh` prints its success line only when every step passes, so a run that stops gives nothing to tell a stop from an interrupted tail. The story agent's prompt (`flai/internal/harness/harness.go`) also hands the verifier the run without naming the step it is known to stop at, I-0057's `flai check --strict`, so the verifier treated the expected stop as something to reproduce. One run, read once, is the aim: the agent log `~/.flai/serve/agents/sf-S-0248-20261004T231346Z.log` from 23:18:07Z to 23:23:54Z shows the waste.

## Acceptance criteria
- [ ] The verifier definition says how to run a long script and read its exit status in one command without a pipe or a file, with the timeout to give it, in `.claude/agents/verifier.md` and `template/root/.claude/agents/verifier.md`
- [ ] `scripts/close-out.sh` ends every run with one line that says its outcome and the step it stopped at, so that the verifier never needs to run it again to learn why it stopped
- [ ] The story agent's prompt, when it hands a close-out run to the verifier, names any step it already knows will stop and why, so that the verifier reports it and goes on to the steps after it rather than re-running
- [ ] `design/conventions/delegation.md` and the template's copy say the same, and the design (`design/system/devex.md` or where close-out is described) records it

## Tasks
- T-0849 close-out.sh ends every run with one line naming its outcome and the step it stopped at
- T-0850 The verifier's definition says how to run a long script once and read its exit status without a pipe or a file
- T-0851 The story agent's prompt names the close-out steps it knows will stop when it hands the run to the verifier
- T-0852 delegation.md, the design, and the user docs say how the verifier runs the close-out once and what the story agent tells it

## Notes

Found by the operator's review of the S-0248 agent log on 2026-10-04. Companion to S-0249, which removes the expected stop itself.

### Planning

The story declared no touches, so `flai touches suggest` was started from the paths that the goal and the criteria name. Where each touch came from:

- **Design (the goal and the criteria name them):**
  - `scripts/close-out.sh`
  - `.claude/agents/verifier.md` and `template/root/.claude/agents/verifier.md`
  - `flai/internal/harness/harness.go`
  - `design/conventions/delegation.md` and `template/root/design/conventions/delegation.md`
- **Co-change, from `flai touches suggest`:**
  - `flai/internal/harness/harness_test.go`, changed with the seeds in 71% of their commits. The prompt's sentences are asserted there.
  - `docs/users/flai.md` (13%), which describes the close-out run.
  - The suggestion also listed `template/CHANGELOG.md` and `template/template.yaml`. They are left out because the release tooling writes them at acceptance.
- **Layout:**
  - `template/root/scripts/close-out.sh`, the template's own close-out. The template's verifier definition will tell every project's verifier to read the last line, so that script needs the line too.
  - `scripts/README.md` and `template/root/scripts/README.md`, which describe each `close-out.sh`.
  - `design/system/workflow.md`, where close-out is described. The criterion allows it in place of `devex.md`. `devex.md` is left to S-0267, which rewrites the test tiers there.
  - `design/system/flai-cli.md`, which describes the prompt that flai serve builds.

**Tags:** `[flai, template]`. `flai check --strict` asked for a component tag, and the story delivers to both components.

**Forecast:** 40m, adjusted from `flai forecast`'s 26m. That figure is the median 89 s per unit of size over 17 done large improvement stories, at size 17 (4 criteria, 13 touches). It was raised for two reasons:

- The close-out runs the flai tier, about two minutes a run, and the verifier runs it twice.
- `.claude/agents/verifier.md` is likely refused to an agent that flai serve starts (I-0069, S-0257). That costs a thread round trip with the operator.

Delivery moves by the same 14m, from 01:06Z to 01:20Z.

**Cost of delay:** 15.00 USD/week, as `flai cod` computes it and unchanged. It is 6m of time lost per 168h cycle at 150 USD an hour. The operator gave the input on TH-0113: the verifier's whole window in the S-0248 log, 23:18:07Z to 23:23:54Z.

**Overlap:** the board's holds will sequence this story against three others:

- S-0267, in progress: `scripts/close-out.sh` and `scripts/README.md`
- S-0249: both `close-out.sh` scripts, `flai-cli.md`, and `docs/users/flai.md`
- S-0268: `delegation()` in `harness.go` and `delegation.md`, though it declares no touches yet.

TH-0116 proposed `after: [S-0249]`, so that the close-out's last line names S-0249's scoped check. It is not set, because the order is the operator's.
