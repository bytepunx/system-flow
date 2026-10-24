---
id: S-0432
type: story
nature: feature
title: The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window
status: in-progress
parent: E-0025
owner: alex
created: 2026-10-16T05:59:39Z
updated: 2026-10-24T13:32:56Z
transitions:
  - to: ready
    at: 2026-10-24T12:40:36Z
    by: alex
  - to: in-progress
    at: 2026-10-24T13:32:56Z
    by: agent-S-0432
tags: [cli, template]
topics: [cli, planning, analysis, conventions]
touches: [".claude/agents/curator.md", design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/internal/check/check.go, flai/internal/check/digest_test.go, flai/internal/digest/decay.go, flai/internal/digest/decay_test.go, flai/internal/guard/curate_test.go, flai/internal/guard/guard.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/.claude/agents/curator.md, template/root/system-flow.yaml.tmpl, template/template.yaml]
after: [S-0431]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 519
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-24
      hand_edits: 1
      work: 34
  models:
    - model: claude-opus-5-5
      input: 116
      output: 771
      cache_read: 6933117
      cache_write: 304422
      cost: 2.4351
  strategic:
    - kind: orchestrator
      seconds: 3471
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 592
          cache_read: 4525871
          cache_write: 18493
          cost: 0.9984
        - model: claude-sonnet-5-5
          input: 4
          output: 11
          cache_read: 68727
          cache_write: 13659
          cost: 0.0563
cost_of_delay:
  value: 95.02
  by: planner-S-0432
  at: 2026-10-23T17:11:07Z
forecast:
  duration: 1h10m
  delivery: 2026-10-24T14:53:00Z
  basis: "Its own forecast of 1h10m; 2nd in the pull order with an in-progress limit of 5, behind S-0493 and S-0433."
  by: flai
  at: 2026-10-24T13:32:37Z
finalized:
  by: orchestrator
  at: 2026-10-23T17:11:31Z
---
# S-0432 The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window

## Goal

The digest stays small and true without the operator asking. With the `curate` action on, `flai serve` starts the curator after each acceptance, coalescing acceptances that come while a run is under way, and when `digest.schedule` comes round. Each run marks absorbed, with `flai digest absorb`, each entry a design document, a convention, or an ADR now records, naming that document, and so drops it from priming. `flai digest compact` turns each active entry that no narrative, thread, or work item has cited within `digest.compact_after` (TH-0476 question 4) into a one-line entry. Nothing is deleted: the state and the body change, and git keeps the rest.

## Acceptance criteria

- [ ] An ADR, refining the curator's, records the triggers (each acceptance, coalesced; `digest.schedule`, read as `analysis.schedule` is), how `cited` is found (an entry's ID in a narrative, a thread, or a work item, archived ones included, with the date of the latest), absorption, and compaction, and that neither deletes anything.
- [ ] `flai serve` starts the curator after an acceptance while `curate` is on and none runs, and once more when acceptances came during a run; `digest.schedule` starts it with the trigger `schedule <spec>`; each start is journalled with its trigger, and tests cover an acceptance, two during a run, and the schedule.
- [ ] `flai digest absorb <D-nnnn> --by <path> [--heading]` sets `status: absorbed` and `absorbed_by`, refusing a path that does not exist; an absorbed entry never loads in a pack (S-0430) and `flai digest list` shows it with the document.
- [ ] `flai digest cited [--json]` gives each entry its latest citation, and `flai digest compact [--dry-run] [--json]` turns each active entry uncited within `digest.compact_after` into `compacted`, its body reduced to its title, and lists what it compacted; a second run compacts nothing more; tests cover a cited entry kept and an uncited one compacted.
- [ ] The curator's prompt and `.claude/agents/curator.md` tell it to look for absorption in every run (with `doc_search` on each active entry's title) and to run `flai digest compact`; the guard passes `flai digest absorb`, `cited`, and `compact` to it.
- [ ] `digest.schedule` and `digest.compact_after` are manifest settings with checks (`manifest.digest`), the template's manifest documents them, and `docs/operators/settings.md`, `design/system/project-manifest.md`, `design/system/strategic-agents.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` describe the triggers, absorption, and compaction.

## Tasks
- T-1831 An ADR records when the curator runs on its own, how a citation is found, and how absorption and compaction change an entry without deleting it
- T-1832 digest.schedule and digest.compact_after are manifest settings with their check, and the template's manifest documents them
- T-1833 flai digest absorb, cited, and compact mark an entry absorbed, find each entry's latest citation, and compact the entries uncited within the window
- T-1834 flai serve starts the curator after each acceptance, coalesced, and when digest.schedule comes round, while the curate action is on
- T-1835 The curator's prompt, definition, and guard take in absorption and compaction
- T-1836 The design and the user and operator docs describe the curator's triggers, digest.schedule, digest.compact_after, absorption, and compaction

## Notes

- Drafted by the planner for E-0025, after S-0431, which gives the curator its run.
- The run after acceptance follows how `planning.replan` starts the planner when work ahead completes (ADR-0084); reuse its coalescing.

### Planning

- **Touches.** The story declared none when drafted; every touch comes from its six tasks, file by file, checked against the main checkout on 2026-10-23.
  - **From layout:** `flai/internal/serve/curate_schedule.go` with its test, on the pattern of `analysis_schedule.go`, wired in `serve.go` beside `newAnalysisScheduler`; `flai/internal/serve/curate.go` and its test (S-0431's), since the start it makes takes the trigger `accepted S-nnnn` or `schedule <spec>` and journals it, as `serve/analyze.go` does; `flai/internal/digest/decay.go` with its test; `flai/cmd/digest.go` and its test; `flai/internal/manifest/manifest.go` and `settings.go` with their tests; `flai/internal/check/check.go` and `flai/internal/check/digest_test.go`, because `manifest.analysis` and `manifest.planning` are reported from `check.go`, which loops the manifest's `Errors()`, so `manifest.digest` is reported there too, and its test follows `analysis_test.go`; `flai/internal/harness/harness.go`, `flai/internal/guard/guard.go`, their tests, and `.claude/agents/curator.md` in both copies; `template/root/system-flow.yaml.tmpl`.
  - **From design:** the ADR and its index; `docs/operators/settings.md`, `project-manifest.md`, `strategic-agents.md`, `flai-cli.md`, `docs/users/flai.md`, `template/CHANGELOG.md`, and `template/template.yaml` for criterion 6; `docs/users/flai-reference.md`, regenerated with the commands, which `reference_test.go` checks.
  - **Co-change** (`flai touches suggest`, 781 of 2289 commits): the strongest are `docs/operators/index.md` (15%), `design/system/flaiover-dashboard.md` (13%), and `docs/users/flaiover.md` (12%), which the dashboard and host-action stories drag in and this story does not change; the curate action's row in `docs/operators/index.md` is S-0431's. `flai/internal/check/check.go` (4%) is added from layout, above. `flai/internal/serve/agents.go` (4%) is not: `AgentRun.Trigger` already exists.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
  - Not on main yet: `flai/internal/digest/`, `flai/cmd/digest.go`, `serve/curate.go`, and `curator.md` come with S-0428 and S-0431, which this story waits for through `after`.
- **Forecast.** 1h10m, from flai's 1h1m (107 s per unit of size over 85 done feature stories on claude-opus-5-5 in the large band, size 34) raised for the four touches added, size 38 at the same rate; delivery is flai's play-out moved by the minutes added. Not raised further: the scheduler mirrors `analysis_schedule.go` and the `cadence` type in `replan.go`, and the decay commands are small.
- **Cost of delay.** 95.02 USD a week, `flai cod` on 2026-10-23: the story's 1h10m share of E-0025's 600.00 USD a week (`time_lost_per_cycle: 4h`, set on the epic after TH-0476 was confirmed) over 7h22m across its eight open stories without inputs. It moves when a sibling's forecast does.
