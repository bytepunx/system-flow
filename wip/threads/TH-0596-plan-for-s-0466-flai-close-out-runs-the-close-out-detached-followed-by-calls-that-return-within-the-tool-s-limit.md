---
id: TH-0596
title: "Plan for S-0466: flai close-out runs the close-out detached, followed by calls that return within the tool's limit"
anchor:
  path: wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md
  item: S-0466
status: resolved
participants: [planner-S-0466, orchestrator]
created: 2026-10-22T21:02:16Z
updated: 2026-10-23T03:21:59Z
---

# TH-0596 Plan for S-0466: flai close-out runs the close-out detached, followed by calls that return within the tool's limit

On wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md.

## Entries

### 2026-10-22T21:02:16Z planner-S-0466
The plan, six tasks in five layers, and the solution it assumes. The story's goal asks for a proposal before building: T-2030 records it as an ADR, and the operator may redirect it there or here.

**The proposed fix.** The close-out moves into flai as `flai close-out <story> [-m | -F] [--wait 9m] [--status] [--cancel]`, running the same rounds as `scripts/close-out.sh` (verify, commit, sync check, stream sync, at most three rounds, ADR-0156). The run is detached from the call that starts it, with its state and log under `.flai-cache/close-out/`, as `flai checks run` and the host's `verify.run` already detach. The caller follows it and returns at `--wait` (default 9m, under the Bash tool's 10 minutes) with a still-running line; the agent calls it again until the last line names the outcome. A caller the tool kills leaves the run alive, so a second round of integration and smoke (12 minutes) finishes on its own. `scripts/close-out.sh`, here and in the template, becomes a wrapper on it. Alternatives for the ADR to weigh: the agent's `run_in_background` Bash call with a 2h timeout (harness-specific; a headless session ends ten minutes after its turn), stopping after round 1 to be run again (the second round still needs 12 minutes), passing over main's commits (rejected by ADR-0156). ADR-0174 expected I-0148 to close once the merge queue ran the `queue` tiers; the operator kept it open, so this fix holds with the queue off.

**Tasks and layers.**

1. T-2030: the ADR (design/adrs/, design/adrs/README.md).
2. T-2031, after T-2030: `flai close-out` runs the rounds in Go, foreground, with the script's last line (flai/cmd/closeout.go, flai/internal/closeout/closeout.go, tests).
3. T-2033, after T-2031: the run detached, state and log under .flai-cache/close-out, `--wait`, `--status`, `--cancel`, with the test that reproduces I-0148 (caller killed with SIGKILL mid-tier, run finishes, a second call prints the outcome). T-2046, after T-2031: both `scripts/close-out.sh` call `flai close-out`, READMEs updated. The two share no path and run together.
4. T-2050, after T-2033 and T-2046: delegation.md and work-management.md (both trees), the harness prompt in harness.go and its test, `.claude/agents/verifier.md` (both trees), docs/users/flai.md, flai-reference.md (generated), flai-cli.md, workflow.md.
5. T-2054, after T-2050: `flai issue close I-0148 --reason`, summary.md, criteria ticked.

**Assumptions.**

- The default wait is a flag's default (9m), not a manifest or config setting, so docs/operators/settings.md gains no row. The ADR may decide otherwise; then T-2050 adds the row.
- The MCP tool `verify` stays synchronous and is not changed: the close-out is the agent's shell call, and the MCP tool is not what I-0148's instance ran.
- `--cancel` is in scope since `flai checks cancel` already has the pattern; drop it from T-2033 if the ADR finds it unneeded.
- The template's `scripts/close-out.sh` needs the flai release that ships `flai close-out`, named in its header and the changelog, as ADR-0156 did for tier reuse; no fallback loop is kept.
- The `design/adrs/` folder touch stands for the new ADR, whose number `flai adr new` allocates; `flai task done` widens the touches to the file.
- Tags `flai, template` were added to the story: `flai check` refused a task while the story touched both components with no component tag.

**Figures.** Cost of delay 62.50 USD a week, flai's figure from the operator's input (25m a cycle), unchanged. Forecast raised from flai's 42m (size 28, 89 s a unit over 87 done improvement stories) to 3h: a Go reimplementation of the close-out against git fixtures, process-group detachment with signal tests, an ADR that waits for the operator's agreement, and two `.claude/` writes that each need a permission prompt are not counted by touches.

### 2026-10-22T21:02:24Z orchestrator
Approved under plan_backlog_stories, as a plan: six tasks in five layers, with the ADR (T-2030) first, and T-2033 and T-2046 running together. The forecast is raised from 42m to 3h, for the Go rewrite, the detachment and signal tests, and the two `.claude/` writes. The topics and tags added so `flai check` passes stand.

S-0466 is in progress, and its goal asks for a proposal before building. The fix you propose, `flai close-out` run detached and followed by calls that return within `--wait`, is a new command and a decision not yet recorded. That is the operator's, so a recommendation follows and the thread stays open. Until it is answered, the story's agent can draft the ADR (T-2030) and should build nothing past it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-22T21:02:29Z orchestrator (recommendation)
Recommended for the operator: take the planner's proposal. Move the close-out into flai as `flai close-out <story>`, run detached under `.flai-cache/close-out/`, with each call returning within `--wait`, 9m by default. `scripts/close-out.sh` becomes a thin wrapper in both trees.

- **It fixes the cause in I-0148.** Two rounds of the integration and smoke tiers take about 24 minutes. A Bash call ends at 10 minutes. A detached run survives the caller, and a second call reads the outcome.
- **It follows a pattern flai already has.** `flai checks run` and the host's `verify.run` detach the same way, so there is precedent and little new machinery.
- **The alternatives are weaker.**
  - A `run_in_background` Bash call ties the fix to one harness, and a headless session ends ten minutes after its turn.
  - Stopping after one round still leaves the 12-minute second round.
  - Passing over main's commits reverses ADR-0156.

Two smaller choices are bundled in: the wait is a flag, not a manifest setting, and `--cancel` is in scope. Say if you want either changed, and the planner adjusts T-2033 or T-2050.

The forecast rises from 42m to 3h, at 62.50 USD a week.

Source: design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md

### 2026-10-22T21:07:17Z orchestrator (recommendation)
Withdrawn: my earlier recommendation here is replaced by the one on TH-0594, which I made after reading S-0466's agent's finding there.

The detached `flai close-out` this plan builds is not what fixes I-0148 in this repo. The tree's dev config lacks the `merge-queue` host action, so close-outs run integration and smoke that ADR-0174 already leaves to the queue. Recommended: answer TH-0594 with option 3. S-0466 then does the small fix, and the agent cancels T-2030, T-2031, and T-2033. This plan becomes a backlog story for projects with the queue off.

Source: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md

### 2026-10-23T03:21:59Z orchestrator
Resolved: Settled by the operator on TH-0594 (option 3, confirmed 2026-10-23T03:21Z): S-0466 does the small fix (the dev config's merge-queue action, flai verify naming the off condition, the dev docs), and the detached flai close-out this plan drafted becomes a later backlog story for projects with the queue off.
