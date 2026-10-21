---
id: S-0428
type: story
nature: feature
title: A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:58:10Z
updated: 2026-10-21T13:25:33Z
transitions: []
tags: [cli, template]
topics: [cli, conventions, planning, analysis]
touches: [design/adrs, design/adrs/README.md, flai/internal/digest/digest.go, flai/internal/digest/digest_test.go, flai/internal/digest/write.go, flai/internal/digest/write_test.go, flai/internal/digest/secret.go, flai/internal/digest/secret_test.go, flai/internal/check/check.go, flai/internal/check/digest_test.go, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/cmd/root.go, docs/users/flai-reference.md, flai/internal/context/sections_test.go, flai/internal/hostapi/docs_test.go, design/digest/README.md, template/root/design/digest/README.md, design/README.md, template/root/design/README.md.tmpl, design/system/repository-layout.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/conventions.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/CHANGELOG.md, template/template.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 51m
  delivery: 2026-10-22T03:05:00Z
  basis: "Its own forecast of 51m; 21st in the pull order with an in-progress limit of 5, behind S-0479, S-0483, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424 and S-0425."
  by: flai
  at: 2026-10-21T13:25:33Z
---
# S-0428 A digest entry is one markdown file under design/digest that flai writes, lists, and checks, and an ADR records where it lives and how it is written

## Goal

The digest has a home and a format before anything fills it. An entry is one fact a story learned: a decision from a narrative's `## Decisions`, an answer the operator gave on a thread, a recommendation they confirmed or overruled and why, or an issue's cause and fix. Each entry is a markdown file under `design/digest/` (or the folder TH-0476 settles) whose title states the fact in one sentence and whose front matter gives its kind, the stories, threads, and issues it comes from, the topics and paths it bears on, its date, and its state. flai writes an entry with `flai digest new`, prints the digest with `flai digest list` and `flai digest show`, and `flai check` refuses an entry that does not parse, cites a source that does not exist, or carries what looks like a secret's value. `doc_search` finds each entry and the dashboard's documents page renders it, as they do an ADR.

## Acceptance criteria

- [ ] An ADR records where the digest lives within the layout (its folder following `layout.design`, as `design/analysis` does), an entry's file name and ID (`D-nnnn`), its front matter (`id`, `title`, `kind`, `sources`, `topics`, `paths`, `date`, `status` of `active`, `absorbed`, or `compacted`, `absorbed_by`, `asked_again`, `cited`), its body, and that nothing is deleted: absorption and compaction change an entry's state and body, and git keeps the rest.
- [ ] `flai/internal/digest` reads every entry of a project, and `flai digest list [--status] [--topic] [--json]` and `flai digest show <D-nnnn> [--json]` print them with their sources, topics, paths, date, and state.
- [ ] `flai digest new "<fact>" --kind <kind> --source <id>... [--topic] [--path] [--json]` allocates the next ID, writes the file in the checkout it runs in, and runs the check on it, refusing and writing nothing when the check fails; `flai digest edit <D-nnnn>` changes an entry's body, topics, paths, sources, and `asked_again` the same way.
- [ ] `flai check` reports `digest.entry` for an entry whose front matter does not parse or misses a field, whose kind or state is unknown, whose source names no work item, thread, or issue in the project or its archive, or whose `absorbed_by` names no file; and `digest.secret` for an entry whose text matches the shape of a key or token, so that an entry says where a secret lives and never its value.
- [ ] `doc_search` and `flai doc search` return an entry's section, and the host API's documents tree, which the dashboard renders, lists `design/digest/` and returns an entry, shown by tests.
- [ ] `design/digest/README.md` in the template and here says what the folder is for and that only flai and the digest's keeper write in it; `design/README.md` and its template, `CLAUDE.md` and its template, `design/system/repository-layout.md`, `docs/users/conventions.md`, `design/system/flai-cli.md`, `docs/users/flai.md` with `flai-reference.md` regenerated, and `template/CHANGELOG.md` describe the folder and the commands.

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

- **Touches.** The story declared none. Every touch comes from its six tasks, file by file.
  - **From design:** the ADR and its index; for criterion 6, `repository-layout.md`, `design/README.md` and its template, `docs/users/conventions.md`, `flai-cli.md`, `docs/users/flai.md`, and `template/CHANGELOG.md`.
  - **From layout:** the new package `flai/internal/digest` (`digest.go`, `write.go`, `secret.go`, and their tests) on the pattern of `flai/internal/analysis`; `flai/cmd/digest.go` and `root.go`; `checker.digest` in `flai/internal/check/check.go` with a new `digest_test.go`, beside `checker.analysis`; tests in `flai/internal/context/sections_test.go` and `flai/internal/hostapi/docs_test.go`, since `doc_search` and the documents tree already walk `design/`; the two READMEs; `template/template.yaml` for the version.
  - **Co-change, added:** `design/system/flai-cli.md` (32% of the seed commits), `docs/users/flai.md` (29%), `flai/cmd/root.go` (17%), `design/adrs/README.md` (10%), `docs/users/flai-reference.md`, `design/README.md`, and `CLAUDE.md` with `template/root/CLAUDE.md.tmpl` (12%), whose layout table lists `design/analysis/` and so takes `design/digest/`.
  - **Co-change, not added:** `flai/internal/check/check_test.go` (a new `digest_test.go` instead); `design/system/project-manifest.md`, `docs/operators/settings.md`, and `flai/internal/manifest/manifest.go`, since the folder follows `layout.design` and adds no setting, which is why the sixth criterion no longer names them; `design/system/flaiover-dashboard.md`, since the documents page needs no change.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** Kept at flai's 51m. The template copies inflate the size about as much as the ADR and the new package add. Delivery is flai's play-out.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
