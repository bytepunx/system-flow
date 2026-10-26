---
id: S-0434
type: story
nature: feature
title: flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it
status: done
parent: E-0025
owner: alex
created: 2026-10-16T06:00:19Z
updated: 2026-10-24T22:56:21Z
transitions:
  - to: ready
    at: 2026-10-24T12:41:07Z
    by: alex
  - to: in-progress
    at: 2026-10-24T21:38:09Z
    by: agent-S-0434
  - to: review
    at: 2026-10-24T22:41:24Z
    by: agent-S-0434
  - to: done
    at: 2026-10-24T22:56:21Z
    by: orchestrator
tags: [cli, dashboard]
topics: [cli, dashboard, analysis]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/metrics.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, flai/cmd/check_stats_test.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/explain.test.ts, flaiover/src/lib/viz/explain.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/template.yaml, design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md, design/issues/summary.md]
after: [S-0433]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3820
  priced_by: harness
  turns:
    - day: 2026-10-24
      ceremony: 2
      hand_edits: 2
      work: 47
  models:
    - model: claude-opus-5-5
      input: 346
      output: 111529
      cache_read: 18677528
      cache_write: 928056
      cost: 11.225
  strategic:
    - kind: orchestrator
      seconds: 13271
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 187
          output: 2979
          cache_read: 31990183
          cache_write: 62421
          cost: 7.0292
        - model: claude-sonnet-5-5
          input: 14
          output: 82
          cache_read: 192345
          cache_write: 73756
          cost: 0.1822
cost_of_delay:
  value: 92.98
  by: planner-S-0434
  at: 2026-10-23T17:26:19Z
forecast:
  duration: 1h15m
  delivery: 2026-10-24T23:05:00Z
  basis: "Its own forecast of 1h15m; 2nd in the pull order with an in-progress limit of 5, behind S-0435."
  by: flai
  at: 2026-10-24T21:37:58Z
finalized:
  by: orchestrator
  at: 2026-10-23T17:35:18Z
---
# S-0434 flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it

## Goal

The analyzer and the operator can see whether the digest cut the questions that wait on the operator (I-0151). `metrics.md` defines, by an ADR, the digest's measure, and `flai stats` reports it per ISO week of the window: the threads opened to the operator, the stories completed, the threads per completed story, and the threads that asked what an entry already answered, those an entry dated before the thread lists in its `asked_again`. The dashboard draws it as a chart beside the agent-waiting chart.

## Acceptance criteria

- [x] An ADR, refining ADR-0081 and ADR-0091, defines in `design/system/metrics.md` a thread opened to the operator (opened by an author other than the operator and other than flai's own notices, such as the upgrade thread), the week it counts in (its first entry's), the completed stories of a week (as `waiting.weeks[]` counts them), the rate, and a repeat (a thread an entry lists in `asked_again` whose `date` is before the thread's first entry).
- [x] `flai stats --json` carries `digest.weeks[]` with `week`, `start`, `threads_to_operator`, `stories_completed`, `per_story` (absent with no completed story), and `asked_again`, and `digest.entries` with the counts of active, compacted, and absorbed entries; `flai stats` prints one line per week under a `digest` heading; tests cover a week with repeats, a week with no completed story, and a project with no digest.
- [x] The dashboard's charts page has a digest chart of threads per completed story per week with the repeats stacked, and a test renders it from a fixture, including a flai that does not report `digest`.
- [x] `design/system/metrics.md` § Charts, `docs/users/flai.md`, and `docs/users/flaiover.md` describe the measure and the chart; the convention `strategic-agents.md` § As the analyzer names the measure for a bottlenecks run.

## Tasks
- T-1842 An ADR defines the digest's measure in metrics.md: threads opened to the operator per completed story per week, and the threads that asked what an entry already answered
- T-1843 flai/internal/metrics computes digest.weeks and digest.entries from the threads, the completed stories, and the entries
- T-1844 flai stats prints the digest's weeks under a digest heading and carries them in --json
- T-1845 The dashboard's charts page draws threads per completed story per week with the repeats stacked
- T-1846 The user docs and the dashboard's design describe the digest's measure and chart, and the analyzer is told to read it on a bottlenecks run

## Notes

- Drafted by the planner for E-0025, after S-0433 in the epic's order. It reads only S-0428's front matter, so it could follow S-0428; the plan thread proposes that, since the weeks before the digest are what the experiment (S-0435) compares against and are already in the threads.
- `metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR (CLAUDE.md).

### Planning

Replanned on 2026-10-23 by planner-S-0434, with S-0428 in review and E-0025's input set. The five tasks were revisited against the criteria: each criterion has its task and their `after` stand; T-1845 and T-1846 gained touches, below. The plan is on TH-0616.

- **Touches.** The story declared none at first; all 27 come from its five tasks, file by file, and `flai touches suggest` lists them as the claim.
  - **From layout:** `flai/internal/metrics/digest.go` with its test, on the pattern of `queue.go` (S-0404): an input on `Options`, a field on the report, a call from `Compute`; `metrics.go`; `flai/internal/statsread/statsread.go` with its test, which already loads the threads and the manifest and so can load the entries with `digest.List(repo)` and pass `repo.Manifest.Owner`; `flai/cmd/stats.go` with `check_stats_test.go` and `stats_test.go`; `flaiover/src/lib/viz/charts.ts` and the charts page with their tests.
  - **From layout, added by this plan:** `flaiover/src/lib/viz/catalog.ts`, since S-0484 moved the chart kinds there and its `TITLES` is a `Record<Kind, string>`, so a new kind does not compile without a title; `flaiover/src/lib/viz/explain.ts`, whose `EXPLAINERS` is a `Record<Kind, ...>` too; and `explain.test.ts`, which runs every kind against a heading of `metrics.md`. S-0404 predates the catalogue, which is why the first plan lacked them.
  - **From design:** the ADR and its index, `metrics.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, `flaiover-dashboard.md`, the convention `strategic-agents.md` in both copies, `template/CHANGELOG.md`, and `template/template.yaml` for criterion 4.
  - **Co-change, added:** `design/system/flai-cli.md` (38%, the highest `flai touches suggest` lists): its `flai stats` row names each section, `queue` since S-0404, which changed it; the digest section belongs there too, on T-1846. Earlier plans added `charts.test.ts`, the charts page and its test, `check_stats_test.go`, `metrics.go`, `metrics.md`, `docs/users/flai.md`, `flaiover-dashboard.md`, `docs/users/flaiover.md`, and, after S-0404 was accepted, `stats_test.go`, `statsread_test.go`, and `docs/users/flai-reference.md`, which a new section's help regenerates.
  - **Co-change, not added:** `docs/operators/index.md` and `settings.md` (14%) and `project-manifest.md`, since no setting is added; `workflow.md`, `work-management.md`, `work-hierarchy.md`, `agent-narrative.md`, and `issues/summary.md`, co-changes of the shared docs this story has no reason to edit; `flai/cmd/root.go`, `harness.go`, `hostapi/writes.go`, and `mcpserver/folder.go`, since no command, prompt, or write is added.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** flai's 49m (size 27: 4 criteria, 23 touches; 26th in the pull order) raised to 1h15m, delivery 2026-10-24T13:23Z, 26m after flai's play-out. S-0404, of the same shape (an ADR on `metrics.md`, a stats section and its print, a chart, the docs), was forecast 37m and took 1h14m of agent time (4464 s in its archived record), and this story adds a package input and the three catalogue files above.
- **Cost of delay.** 92.98 USD a week: S-0434's share of E-0025's 600 USD a week, 1h15m of the 8h4m forecast over its eight stories, from `flai cod` on 2026-10-23 with the forecast above. The epic's input, `time_lost_per_cycle: 4h`, is the answer the operator confirmed on TH-0476 (2026-10-21). Kept at flai's figure.

### Accepted by the orchestrator

- Verified: aa9635cf0d90fe0aae3e534300db4e6697e4518e
- Batch: 20261024-224219
- Enqueued verified: aa9635cf0d90fe0aae3e534300db4e6697e4518e
- At: 2026-10-24T22:56:21Z

Verdict: S-0434 meets all four criteria at aa9635cf, every changed file under its touches and none under .claude/; flai verify passed at that head with integration and smoke left to the queue.
- 1: design/adrs/0200-flai-stats-reports-per-week-the-threads-opened-to-the-operator-the-stories.md, design/adrs/README.md, design/system/metrics.md
- 2: flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flai/cmd/check_stats_test.go, docs/users/flai-reference.md
- 3: flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, flaiover/src/routes/charts/[kind]/+page.svelte, flaiover/src/routes/charts/[kind]/charts.svelte.test.ts
- 4: design/system/metrics.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml
