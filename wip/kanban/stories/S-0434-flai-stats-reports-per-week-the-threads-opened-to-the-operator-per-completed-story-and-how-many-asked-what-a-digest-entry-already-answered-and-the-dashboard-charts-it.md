---
id: S-0434
type: story
nature: feature
title: flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:19Z
updated: 2026-10-17T03:34:53Z
transitions: []
tags: [cli, dashboard]
topics: [cli, dashboard, analysis]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/system/flaiover-dashboard.md, design/system/metrics.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/check_stats_test.go, flai/cmd/stats.go, flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/charts.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/template.yaml]
after: [S-0433]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 39m
  delivery: 2026-10-18T02:18:00Z
  basis: "Its own forecast of 39m; 43rd in the pull order with an in-progress limit of 5, behind S-0442, S-0444, S-0443, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432 and S-0433."
  by: flai
  at: 2026-10-17T03:34:53Z
---
# S-0434 flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it

## Goal

The analyzer and the operator can see whether the digest cut the questions that wait on the operator (I-0151). `metrics.md` defines, by an ADR, the digest's measure, and `flai stats` reports it per ISO week of the window: the threads opened to the operator, the stories completed, the threads per completed story, and the threads that asked what an entry already answered, those an entry dated before the thread lists in its `asked_again`. The dashboard draws it as a chart beside the agent-waiting chart.

## Acceptance criteria

- [ ] An ADR, refining ADR-0081 and ADR-0091, defines in `design/system/metrics.md` a thread opened to the operator (opened by an author other than the operator and other than flai's own notices, such as the upgrade thread), the week it counts in (its first entry's), the completed stories of a week (as `waiting.weeks[]` counts them), the rate, and a repeat (a thread an entry lists in `asked_again` whose `date` is before the thread's first entry).
- [ ] `flai stats --json` carries `digest.weeks[]` with `week`, `start`, `threads_to_operator`, `stories_completed`, `per_story` (absent with no completed story), and `asked_again`, and `digest.entries` with the counts of active, compacted, and absorbed entries; `flai stats` prints one line per week under a `digest` heading; tests cover a week with repeats, a week with no completed story, and a project with no digest.
- [ ] The dashboard's charts page has a digest chart of threads per completed story per week with the repeats stacked, and a test renders it from a fixture, including a flai that does not report `digest`.
- [ ] `design/system/metrics.md` § Charts, `docs/users/flai.md`, and `docs/users/flaiover.md` describe the measure and the chart; the convention `strategic-agents.md` § As the analyzer names the measure for a bottlenecks run.

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

- **Touches.** The story declared none. Every touch comes from its five tasks, file by file.
  - **From layout:** a new `flai/internal/metrics/digest.go` with its test, on the pattern of `waiting.go`; `metrics.go` for the report's field; `flai/internal/statsread/statsread.go` to load the entries; `flai/cmd/stats.go` and `check_stats_test.go`; `flaiover/src/lib/viz/charts.ts`, `charts.test.ts`, and the charts page with its test.
  - **From design:** the ADR, its index, and `metrics.md`; `docs/users/flai.md`, `docs/users/flaiover.md`, `flaiover-dashboard.md`, the convention `strategic-agents.md` in both copies, `template/CHANGELOG.md`, and `template/template.yaml` for criterion 4.
  - **Co-change, added:** `flaiover/src/lib/viz/charts.test.ts` (51% of the commits that changed `waiting.go`, `stats.go`, or `charts.ts`), the charts page (31%) and its test (27%), `flai/cmd/check_stats_test.go` (29%), `flai/internal/metrics/metrics.go` (27%), `design/system/metrics.md` (24%), `docs/users/flai.md` (22%), `flaiover-dashboard.md` (13%), and `docs/users/flaiover.md` (13%).
  - **Co-change, not added:** `docs/users/flai-reference.md` (18%), since no flag is added; `usage.go`, `spend.go`, `strategic.go`, `SpendTable.svelte.test.ts`, and `palette.ts`, since neither usage nor spend changes and the chart reads the palette as it is.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** Kept at flai's 39m. It is a metric on the waiting section's pattern and one chart. Delivery is flai's play-out moved by S-0429's raise.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
