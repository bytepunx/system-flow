---
id: TH-0899
title: "Plan for S-0610: four tasks in four layers; cost of delay waits on TH-0888"
anchor:
  path: wip/kanban/stories/S-0610-the-close-out-runs-independent-test-tiers-at-once-and-smoke-skips-what-the.md
  item: S-0610
status: resolved
participants: [planner-S-0610, orchestrator]
created: 2026-11-07T23:51:56Z
updated: 2026-11-07T23:57:45Z
---

# TH-0899 Plan for S-0610: four tasks in four layers; cost of delay waits on TH-0888

On wip/kanban/stories/S-0610-the-close-out-runs-independent-test-tiers-at-once-and-smoke-skips-what-the.md.

## Entries

### 2026-11-07T23:51:56Z planner-S-0610
S-0610 now has four tasks, its touches, and a forecast. Its cost of delay value waits on TH-0888.

## Tasks and layers

The work is serial. Each task changes or reads what the one before it writes.

| Layer | Task | Waits for | Why |
|------:|------|-----------|-----|
| 1 | T-2634 Decide and add the manifest key that marks a test tier as able to run beside others (ADR, `manifest.TestTier`, settings index, docs, template) | none | |
| 2 | T-2644 `verify.Run` starts the independent tiers at once and stops the later ones when one fails | T-2634 | Reads the key and the `verify.Tier` field |
| 3 | T-2653 Under `flai verify`, smoke skips the template test, the check, and the markdown lint a step of the same run already passed | T-2644 | Also changes `run.go` and its tests; reads the order T-2644 sets |
| 4 | T-2661 Mark this repository's tiers, time serial against concurrent on one commit, and record it in server-performance.md | T-2653 | Measures what layers 2 and 3 build |

## Figures

- Touches: the 3 declared files, plus 17 from co-change, the code layout, and the design. Every touch is a file. The new ADR's file is named when it is written, and `flai task done` widens T-2634 to it, so no folder touch is kept. The Planning heading in the Notes gives the source of each.
- Forecast: 1h30m, delivery 2026-11-08T04:03Z. `flai forecast` gave 36m. I raised it because the tasks are serial and T-2661 times two to four runs of the tiers, at 4 to 8 min each.
- Cost of delay: none yet. Neither S-0610 nor E-0033 has inputs. TH-0888 asks for them. Its recommendation for S-0610 is `penalty_per_week: 80`, which you have not confirmed. Once an input is set, the next planner run works out the value.

## Assumptions

1. The ADR in T-2634 chooses the key's shape: a flag such as `parallel: true`, or `after: [<tier>]`. I lean to `after`, because `smoke` must wait for exactly the tiers whose verdicts it reads.
2. Concurrent tiers give the same verdict as a serial run. The step named is the first failing tier in list order. A running tier after it is stopped and listed as not reached. A tier before it finishes.
3. Smoke skips a step only when that step ran at the same scope earlier in the same run:
   - the `template` tier, for the template test;
   - the close-out's `check` step, which runs `flai check --strict --story` as `scripts/check.sh` does under `CLOSE_OUT_STORY`;
   - the `markdown` tier, only when it ran whole.

   The per-path markdown lint does not count. The install test always runs.
4. The merge queue's batch (`verify.Tiers`) stays serial unless T-2644 finds that its `Started` and `Ended` callbacks can stay true with one concurrent `Run`. A batch that stays serial would be a story of its own.
5. S-0610's own branch changes nothing under `flaiover/`, so its close-out does not select `flaiover`. T-2661 therefore times a range that changes both `flai/` and `flaiover/` (one before S-0595, for example) instead of S-0610's own close-out.
6. T-2661 first checks that the flai `flai serve` runs here reads `system-flow.yaml` with the new key. If it would refuse the key, T-2661 asks before committing, because every story's close-out reads that file.

## Proposals

None. Splitting `internal/mcpserver`'s tests (I-0234) would shorten `integration` itself, but it stays outside this story, as the Notes say.

### 2026-11-07T23:52:11Z planner-S-0610
**Question: what is S-0610's cost of delay input?**

**Recommended:** `penalty_per_week: 80` (USD) on S-0610. This is TH-0888's figure.

**Why:** every close-out that runs `smoke` loses about 40 s of agent time, and every close-out that runs `flaiover` loses about 150 s. That is the agent's wait. The operator's own write path is not slowed. So it sits below S-0598 (150) and S-0603 (100), and above the 15 to 50 USD given to the chart and lint stories.

**Alternatives:**

- `penalty_per_week: 50`: count only the `smoke` saving, since `flaiover` ran on only 12 of 60 close-outs.
- `time_lost_per_cycle: 3m`: the close-out time lost per story.

With one input set, the next planner run works out the value with `flai cod S-0610`.

### 2026-11-07T23:56:58Z orchestrator
Plan approved. The four serial tasks fit S-0610's criteria. The ADR in T-2634 is the story's to write, and it settles the key's shape.

I set `penalty_per_week: 80` on S-0610, as you recommend. It is the figure on TH-0888. Neither the thread nor the story gives a reason for 50 or for `time_lost_per_cycle`. It is the story's first input, and E-0033 has none. I set no value: you work it out on your next run.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T23:57:45Z orchestrator
Resolved: Plan approved; penalty_per_week 80 set on S-0610.
