---
id: S-0432
type: story
nature: feature
title: The curator runs after each acceptance and on a schedule, marks an entry the design now records absorbed, and compacts an entry no one has cited within the window
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:59:39Z
updated: 2026-10-21T12:13:57Z
transitions: []
tags: [cli, template]
topics: [cli, planning, analysis, conventions]
touches: [".claude/agents/curator.md", design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/internal/digest/decay.go, flai/internal/digest/decay_test.go, flai/internal/guard/curate_test.go, flai/internal/guard/guard.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/serve.go, template/CHANGELOG.md, template/root/.claude/agents/curator.md, template/root/system-flow.yaml.tmpl, template/template.yaml]
after: [S-0431]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 55m
  delivery: 2026-10-22T05:53:00Z
  basis: "Its own forecast of 55m; 25th in the pull order with an in-progress limit of 5, behind S-0461, S-0471, S-0473, S-0474, S-0479, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430 and S-0431."
  by: flai
  at: 2026-10-21T12:13:57Z
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

- **Touches.** The story declared none. Every touch comes from its six tasks, file by file.
  - **From layout:** `flai/internal/serve/curate_schedule.go` with its test, on the pattern of `analysis_schedule.go`, and its wiring in `serve.go`; `flai/internal/digest/decay.go` with its test; `flai/cmd/digest.go` and its test; `flai/internal/manifest/manifest.go` and `settings.go` with their tests; `flai/internal/harness/harness.go`, `flai/internal/guard/guard.go`, their tests, and `.claude/agents/curator.md` in both copies; `template/root/system-flow.yaml.tmpl`.
  - **From design:** the ADR and its index; `docs/operators/settings.md`, `project-manifest.md`, `strategic-agents.md`, `flai-cli.md`, `docs/users/flai.md`, `template/CHANGELOG.md`, and `template/template.yaml` for criterion 6.
  - **Co-change, added:** `flai/internal/manifest/manifest.go` (75% of the commits that changed `analysis_schedule.go` or `settings.go`), `docs/operators/settings.md` (62%), `manifest_test.go` (62%), `docs/users/flai-reference.md` (38%), and `settings_test.go` (38%).
  - **Co-change, not added:** `flai/internal/check/check.go` (25%), since the manifest's own check reports `manifest.digest`, as it does `manifest.analysis`.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** Kept at flai's 55m. The scheduler mirrors the analysis scheduler and the decay commands are small. Delivery is flai's play-out moved by S-0429's raise.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
