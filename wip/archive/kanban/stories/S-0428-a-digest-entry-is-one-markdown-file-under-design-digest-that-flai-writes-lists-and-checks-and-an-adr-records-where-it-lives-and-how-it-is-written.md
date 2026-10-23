---
id: S-0428
type: story
nature: feature
title: A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written
status: done
parent: E-0025
owner: alex
created: 2026-10-16T05:58:10Z
updated: 2026-10-23T17:35:05Z
transitions:
  - to: ready
    at: 2026-10-23T16:36:50Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-23T16:36:56Z
    by: agent-S-0428
  - to: review
    at: 2026-10-23T17:17:36Z
    by: agent-S-0428
  - to: done
    at: 2026-10-23T17:35:05Z
    by: orchestrator
tags: [cli, template]
topics: [cli, conventions, planning, analysis]
touches: [design/adrs, design/adrs/README.md, flai/internal/digest/digest.go, flai/internal/digest/digest_test.go, flai/internal/digest/write.go, flai/internal/digest/write_test.go, flai/internal/digest/secret.go, flai/internal/digest/secret_test.go, flai/internal/check/check.go, flai/internal/check/digest_test.go, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/cmd/root.go, docs/users/flai-reference.md, flai/internal/context/sections_test.go, flai/internal/hostapi/docs_test.go, design/digest/README.md, template/root/design/digest/README.md, design/README.md, template/root/design/README.md.tmpl, design/system/repository-layout.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/conventions.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/CHANGELOG.md, template/template.yaml, docs/operators/settings.md, flai/internal/guard/digest_test.go, flai/internal/guard/guard.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2452
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 4
      hand_edits: 4
      work: 55
  models:
    - model: claude-opus-5-5
      input: 402
      output: 165571
      cache_read: 22298638
      cache_write: 847816
      cost: 12.7647
  strategic:
    - kind: orchestrator
      seconds: 2134
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 39
          output: 463
          cache_read: 3270191
          cache_write: 26727
          cost: 0.7252
        - model: claude-sonnet-5-5
          input: 9
          output: 834
          cache_read: 186391
          cache_write: 22839
          cost: 0.1447
cost_of_delay:
  value: 95.06
  by: planner-S-0428
  at: 2026-10-23T16:20:54Z
forecast:
  duration: 1h1m
  delivery: 2026-10-24T06:09:00Z
  basis: "flai forecast 2026-10-23: median 107 s per unit of size over 85 done feature stories on claude-opus-5-5 in the large band, times size 34 (6 criteria, 28 touches); 21st in the pull order with an in-progress limit of 5. Kept by planner-S-0428: S-0414, of the same shape (an ADR, a new package, list and show commands, check findings, template files), was forecast 1h1m."
  by: planner-S-0428
  at: 2026-10-23T16:20:11Z
finalized:
  by: orchestrator
  at: 2026-10-23T16:36:40Z
---
# S-0428 A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written

## Goal

The digest has a home and a format before anything fills it. An entry is one fact a story learned: a decision from a narrative's `## Decisions`, an answer the operator gave on a thread, a recommendation they confirmed or overruled and why, or an issue's cause and fix. Each entry is a markdown file under `design/digest/` (or the folder TH-0476 settles) whose title states the fact in one sentence and whose front matter gives its kind, the stories, threads, and issues it comes from, the topics and paths it bears on, its date, and its state. flai writes an entry with `flai digest new`, prints the digest with `flai digest list` and `flai digest show`, and `flai check` refuses an entry that does not parse, cites a source that does not exist, or carries what looks like a secret's value. `doc_search` finds each entry and the dashboard's documents page renders it, as they do an ADR.

## Acceptance criteria

- [x] An ADR records where the digest lives within the layout (its folder following `layout.design`, as `design/analysis` does), an entry's file name and ID (`D-nnnn`), its front matter (`id`, `title`, `kind`, `sources`, `topics`, `paths`, `date`, `status` of `active`, `absorbed`, or `compacted`, `absorbed_by`, `asked_again`, `cited`), its body, and that nothing is deleted: absorption and compaction change an entry's state and body, and git keeps the rest.
- [x] `flai/internal/digest` reads every entry of a project, and `flai digest list [--status] [--topic] [--json]` and `flai digest show <D-nnnn> [--json]` print them with their sources, topics, paths, date, and state.
- [x] `flai digest new "<fact>" --kind <kind> --source <id>... [--topic] [--path] [--json]` allocates the next ID, writes the file in the checkout it runs in, and runs the check on it, refusing and writing nothing when the check fails; `flai digest edit <D-nnnn>` changes an entry's body, topics, paths, sources, and `asked_again` the same way.
- [x] `flai check` reports `digest.entry` for an entry whose front matter does not parse or misses a field, whose kind or state is unknown, whose source names no work item, thread, or issue in the project or its archive, or whose `absorbed_by` names no file; and `digest.secret` for an entry whose text matches the shape of a key or token, so that an entry says where a secret lives and never its value.
- [x] `doc_search` and `flai doc search` return an entry's section, and the host API's documents tree, which the dashboard renders, lists `design/digest/` and returns an entry, shown by tests.
- [x] `design/digest/README.md` in the template and here says what the folder is for and that only flai and the digest's keeper write in it; `design/README.md` and its template, `CLAUDE.md` and its template, `design/system/repository-layout.md`, `docs/users/conventions.md`, `design/system/flai-cli.md`, `docs/users/flai.md` with `flai-reference.md` regenerated, and `template/CHANGELOG.md` describe the folder and the commands.

## Tasks
- T-1807 An ADR records the digest: its folder under design, an entry's file name, ID, front matter, and body, and that nothing in it is deleted
- T-1808 flai/internal/digest reads, validates, and writes digest entries, and refuses one whose text has a secret's shape
- T-1809 flai check reports digest.entry and digest.secret on the entries under design/digest
- T-1810 flai digest new, edit, list, and show write and print digest entries, as text and JSON
- T-1811 doc_search finds a digest entry and the dashboard's documents tree lists design/digest, shown by tests
- T-1812 The digest folder's README, the layout, the CLI design, and the user docs say what the digest is, where it lives, and how flai writes it, and the template ships the folder

## Notes

- Drafted by the planner for E-0025. First story of the epic: the format every later story reads and writes.
- TH-0476 questions 2 and 5 shape the ADR: the folder, one file per entry, and who counts as the operator. The story builds on the recommended answers until they are answered otherwise.
- S-0426 (accepted 2026-10-16) recommends, in `design/system/agent-memory.md` § Recommendation, adapting flai's own deterministic index and reaching the archive and threads through this digest under `design/`. The ADR cites it.
- `cited` is filled by S-0432's `flai digest cited` from the citations S-0430's conventions ask agents to write, and read by its compaction; `asked_again` is filled by the curator (S-0431) and read by the metric (S-0434).

### Planning

Replanned on 2026-10-23 by planner-S-0428, after the operator confirmed TH-0476 on 2026-10-21. The six tasks were revisited against the criteria: each criterion has its task, their touches name files, and their `after` stand, so none was changed.

- **Touches.** The story declared none at first; all 28 come from its six tasks, file by file, and `flai touches suggest` now lists them as the claim.
  - **From design:** the ADR and its index (T-1807); for criterion 6, `repository-layout.md`, `design/README.md` and its template, `docs/users/conventions.md`, `flai-cli.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` (T-1812).
  - **From layout:** the new package `flai/internal/digest` (`digest.go`, `write.go`, `secret.go`, and their tests) on the pattern of `flai/internal/analysis` (`Folder`, `Dir(m)`, `Validate`); `checker.digest` beside `checker.analysis` in `flai/internal/check/check.go`, whose `documentation()` skips `analysis` and `experiments` from the design documents' own rules and must skip `digest` the same way, with a new `digest_test.go`; `flai/cmd/digest.go` and `root.go`; tests in `flai/internal/context/sections_test.go` and `flai/internal/hostapi/docs_test.go`, since `doc_search` walks `design/` and `docs.tree` walks it with no folder list, as `docs.svelte.test.ts` shows for `design/analysis`; the two READMEs; `template/template.yaml` for the version.
  - **Co-change, added at the first plan:** `design/system/flai-cli.md`, `docs/users/flai.md`, `flai/cmd/root.go`, `design/adrs/README.md`, `docs/users/flai-reference.md`, `design/README.md`, and `CLAUDE.md` with `template/root/CLAUDE.md.tmpl`, whose layout table lists `design/analysis/` and so takes `design/digest/`.
  - **Co-change, not added:** `docs/operators/index.md` and `docs/operators/settings.md` (16% each), `design/system/project-manifest.md`, and `flai/internal/manifest/manifest.go`, since the folder follows `layout.design` and adds no setting; `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` (15%, 13%), since the documents page lists the folder from `docs.tree` as it lists `design/analysis`, with no flaiover change; `design/system/workflow.md`, `design/conventions/work-management.md`, and `design/issues/summary.md`, co-changes of the shared docs that this story has no reason to edit. `flaiover/src/lib/sitemenu.ts` names the design folders the site menu groups, but the story asks only that the documents tree list the folder; a menu entry would be a later story's.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** flai's 1h1m, delivery 2026-10-24T06:09Z, 21st in the pull order. Up from the 51m of the first plan: the size, 6 criteria and 28 touches, now falls in the large band. Kept, because S-0414, of the same shape (an ADR, a new package, list and show commands, check findings, template files), was forecast 1h1m.
- **Cost of delay.** 95.06 USD a week: S-0428's share of E-0025's 600 USD a week, 1h1m of the 6h25m forecast over its eight stories. The epic's input, `time_lost_per_cycle: 4h`, is the answer the operator confirmed on TH-0476 (2026-10-21T12:59Z); it was not yet on the epic, so planner-S-0428 recorded it there as confirmed, as planner-E-0023 did on TH-0470, and set nothing else. Kept at flai's figure.

### Accepted by the orchestrator

- Verified: c9e2bdca50e210104b35f674cc9d55780adf05a8
- Batch: 20261023-172401
- Enqueued verified: c9e2bdca50e210104b35f674cc9d55780adf05a8
- At: 2026-10-23T17:35:05Z

Verdict: pass. All six criteria are met and ticked, every changed file is within the story's touches, and the template copies agree.
- 1: design/adrs/0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md, design/adrs/README.md
- 2: flai/internal/digest/digest.go, flai/internal/digest/write.go, flai/internal/digest/secret.go, flai/cmd/digest.go, flai/internal/digest/digest_test.go, flai/cmd/digest_test.go
- 3: flai/cmd/digest.go, flai/internal/digest/write.go, flai/cmd/root.go, flai/internal/guard/guard.go, flai/internal/guard/digest_test.go
- 4: flai/internal/check/check.go, flai/internal/check/digest_test.go, flai/internal/digest/secret_test.go
- 5: flai/internal/context/sections_test.go, flai/internal/hostapi/docs_test.go
- 6: design/digest/README.md, template/root/design/digest/README.md, design/README.md, template/root/design/README.md.tmpl, CLAUDE.md, template/root/CLAUDE.md.tmpl, design/system/repository-layout.md, docs/users/conventions.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md
