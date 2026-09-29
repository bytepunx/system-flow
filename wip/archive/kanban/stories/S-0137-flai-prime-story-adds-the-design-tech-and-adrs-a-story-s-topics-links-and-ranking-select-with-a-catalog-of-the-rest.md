---
id: S-0137
type: story
nature: feature
title: flai prime --story adds the design, tech, and ADRs a story's topics, links, and ranking select, with a catalog of the rest
status: done
parent: E-0010
owner: alex
created: 2026-09-26T17:46:25Z
updated: 2026-09-29T00:52:14Z
transitions:
  - to: ready
    at: 2026-09-26T22:40:37Z
    by: alex
  - to: in-progress
    at: 2026-09-29T00:33:25Z
    by: agent-S-0137
  - to: review
    at: 2026-09-29T00:49:23Z
    by: agent-S-0137
  - to: done
    at: 2026-09-29T00:52:14Z
    by: alex
tags: [cli]
touches: [flai/cmd/prime.go, flai/internal/context, flai/internal/search, flai/cmd/prime_test.go, design/system/flai-cli.md, design/system/agent-context.md, design/system/conventions.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0137 flai prime --story adds the design, tech, and ADRs a story's topics, links, and ranking select, with a catalog of the rest

## Goal

An agent starts its story with the design and decisions it builds on, and knows what else exists, as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides. After the conventions, `flai prime --story` prints, each headed with its path, heading path, and reason: the `design/system`, `design/tech`, and ADR sections whose topics match; every document or ADR the story, its epic, and its tasks link or name by ID, and one step further (ADRs a selected section links, ADRs a selected ADR refines), with a superseded ADR replaced by what supersedes it; the five ADRs and five design sections ranked highest by BM25 against the title, goal, and criteria among those not selected; then a catalog of every document not loaded and the heading outline of each loaded in part.

## Acceptance criteria
- [x] Each selection step is a function in the context package with behavior tests over a fixture repository: topics, links from story, epic, and tasks, one step of links, supersession, ranking, catalog, and no document printed twice (the first reason wins, the rest are listed).
- [x] Reasons read `topics: cli`, `linked from S-nnnn`, `linked from design/system/x.md § Heading`, `refined by ADR-nnnn`, `supersedes ADR-nnnn`, `rank n`.
- [x] `--json` returns every item with path, heading path, reason, and size; the header's size includes them.
- [x] Ranking reuses `flai/internal/search`, indexing sections rather than whole files, without a new dependency.
- [x] Run on this repository for S-0125 and for two archived stories, the output is attached to the story notes with its size, and every ADR those stories named that existed when they were created is either printed or in the catalog.
- [x] `design/system/flai-cli.md`, `design/system/agent-context.md`, and `docs/users/flai.md` describe the pack; all three test tiers and `flai check --strict` pass.

## Tasks
- T-0505 The context package loads design, tech, and ADRs and selects them by topics, links, one step of links, and supersession
- T-0506 Rank the sections not selected with flai/internal/search, and catalog the documents not loaded
- T-0507 flai prime --story prints the design items and the catalog, in text and --json
- T-0508 Replay the pack on S-0125 and two archived stories and attach the outputs to the story notes
- T-0509 Describe the pack in the design and user docs, and pass all three test tiers and flai check --strict

## Notes

After S-0136. No size budget (TH-0021 4b); the ranked step's fixed count of five and five is from the S-0125 replay. Sections are cut at headings only.

### Replay (T-0508), 2026-09-29

`flai prime --story` on this repository at `story/S-0137`. The full outputs are about 500 KB each, so they are not pasted here. This is what each one printed. `flai prime --story <id>` reproduces them, and `--json` gives every item.

| Story | Pack size | Design items | By topics | Linked | Ranked | Catalog, not loaded / in part |
|-------|-----------|--------------|-----------|--------|--------|-------------------------------|
| S-0125 (cli) | 489,627 bytes, 4,265 lines | 65, 414,794 bytes | 22 | 33 | 10 | 16 / 2 |
| S-0112 (cli) | 508,037 bytes, 4,266 lines | 64, 433,792 bytes | 22 | 32 | 10 | 18 / 1 |
| S-0116 (dashboard, cli) | 548,217 bytes, 4,619 lines | 71, 470,992 bytes | 29 | 32 | 10 | 10 / 2 |

- By topics, for all three: every `design/system` and `design/tech` file whose topics are `all`, `cli`, `go`, or `code`, whole. That is 16 system and 6 tech files; S-0116 adds `flaiover-dashboard.md` and the dashboard and sveltekit tech files.
- Linked: what each story names (S-0125 ADR-0047; S-0112 ADR-0041; S-0116 ADR-0042, replaced by ADR-0043, and ADR-0041 from its task T-0411). Then the ADRs the selected sections link, 27 to 30 of them, most from `flai-cli.md § Commands` and `dashboard-host-channel.md § Decision`. ADR-0022, 0034, 0038, and 0042 appear only as what supersedes them.
- Ranked, S-0125: `conventions.md` § Agent conventions, § Baseline topics, § Priming; `flaiover-dashboard.md` § Writes, § Views; ADR-0003, 0009, 0002, 0012, 0005.
- Ranked, S-0112: five `flaiover-dashboard.md` sections (Writes, API, Views, One dashboard for every project, The channel to flai on the host); ADR-0009, 0005, 0003, 0002, 0001.
- Ranked, S-0116: `conventions.md` § Topics, § Baseline topics, § Priming, § Agent conventions; `template.md` § Upgrading a project; ADR-0009, 0002, 0001, 0005, 0011.

Every ADR a story named that existed when it was created is printed or in the catalog:

| Story | Created | Named | Where |
|-------|---------|-------|-------|
| S-0125 | 2026-09-26 | ADR-0047 | printed, linked from S-0125 |
| S-0112 | 2026-09-24 | ADR-0038, ADR-0041 | ADR-0041 printed, linked from S-0112; ADR-0038 in the catalog, superseded by ADR-0041 |
| S-0116 | 2026-09-24 | ADR-0042, ADR-0043 | ADR-0043 printed, supersedes ADR-0042; ADR-0042 in the catalog, superseded by ADR-0043 |

Across the whole archive, 56 stories name 93 ADRs that existed when each was created. The counts below leave out the story's own links:

| Selection | ADRs reached |
|-----------|--------------|
| Topics alone (no ADR carries topics yet) | 0 |
| Topics, then one step from the selected sections | 91 |
| The above, then ranking | 93 |
| Topics and ranking, no step | 68 |
| Ranking alone, top 5, sections | 47 |

Ranking alone does worse here than the 80% of the S-0125 replay: it tends to pick the short early ADRs and the long `flaiover-dashboard.md` sections. Scoring ADRs whole, using the title and goal only, or matching exact terms only each came within two ADRs of 47. The pack's size, and whether to narrow it before S-0138, is asked on TH-0029.
