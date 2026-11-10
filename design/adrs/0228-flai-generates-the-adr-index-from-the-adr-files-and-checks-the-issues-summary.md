---
id: ADR-0228
title: "flai generates the ADR index from the ADR files and checks the issues summary against its folder, CHANGELOG.md files merge by union under one Unreleased section, and the acceptance notice leaves out the paths in claims.shared"
status: accepted
date: 2026-11-10
supersedes: [ADR-0157]
superseded_by: []
refines: [ADR-0096, ADR-0098, ADR-0175]
---

# ADR-0228 flai generates the ADR index from the ADR files and checks the issues summary against its folder, CHANGELOG.md files merge by union under one Unreleased section, and the acceptance notice leaves out the paths in claims.shared

## Context

Every story edits the same few files. S-0550's Goal counts, from 2026-10-15 to 2026-10-29, 326 touches of `design/issues/summary.md`, 242 of `docs/users/flai.md`, 223 of `design/system/flai-cli.md`, 147 of `docs/users/flai-reference.md`, and 99 of `design/adrs/README.md`, and most stories also add lines to `template/CHANGELOG.md`. It counts 248 overlap notices over 97 acceptances and 62 conflicts, almost all on those files. Each notice made an open story's agent run `flai stream sync` again.

The recount in [workflow.md](../system/workflow.md) § "Hub files: the fortnight before S-0550", from `.flai-cache/overlaps.jsonl`, `.flai-cache/conflicts.jsonl`, and the acceptance commits, narrows those figures to the fortnight itself:

- 131 acceptances sent 190 notices, 1.45 per acceptance. The other 58 of the 248 were grown-claim notices.
- 180 of the 190 named a path inside `claims.shared`, and 123 named only such paths.
- 42 conflicts were recorded in the fortnight. 25 were on a shared path, 20 only on shared paths. `template/CHANGELOG.md` had 9.

Earlier decisions settle part of this:

- [ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md) writes the issues summary again when a rebase stops on it. Nothing checks that the committed summary is what its folder renders.
- [ADR-0157](0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md) merges the ADR index row by row and keeps it hand-editable. It rejected regenerating the index because some status cells carry notes the front matter does not hold, such as "superseded in part by 0016".
- [ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md) §2 frees the paths in `claims.shared` from the hold, but keeps them in the acceptance notice.
- `template/CHANGELOG.md` has no merge rule of its own beyond the section merge of [ADR-0175](0175-a-rebase-that-stops-on-a-markdown-file-inside-the-manifest-s-shared-paths.md), which leaves a stop when both sides change the same lines of one section.

Read against its folder, the index differs in 11 of its 223 status cells. Titles and links match in every row. Of the 11:

- Three are wrong today: 0067 and 0077 read `proposed` but are accepted, and 0202 reads `proposed` but is rejected.
- Eight carry notes the front matter does not hold: "in part" on 0007, 0016, 0038, and 0041, "refined by 0022" on 0019, and refines that the files of 0022, 0023, and 0024 do not list.

Front matter cannot tell a supersession in part from a whole one. `flai adr new --supersedes` sets only `superseded_by`, and the conventions forbid any other edit to an accepted ADR. So every superseded ADR here keeps `status: accepted`: the two superseded in part (0007, 0038) and the six superseded whole (0022, 0024, 0026, 0027, 0034, 0042).

## Decision

**flai generates the ADR index from the ADR files and checks the issues summary against its folder, `CHANGELOG.md` files merge by union under one `## Unreleased` section, and the acceptance notice leaves out the paths in `claims.shared`.** This supersedes ADR-0157 in part, and refines ADR-0096, ADR-0098, and ADR-0175.

### 1. The ADR index is generated

- `flai adr index` writes `design/adrs/README.md` from the ADR files. `flai adr new`, `flai adr accept`, and the supersede that `adr new --supersedes` does write the whole index again, in place of editing one row.
- The table has one row per `NNNN-*.md` file in the folder, `0000-template.md` left out, in number order. A row is `| [NNNN](<file>) | <title> | <status cell> |`, the title from front matter with `|` escaped.
- The text above the table and below it is kept as written. Only the table is rewritten.
- `flai check` reports `adr.index`, a warning, when the index differs from what its folder renders, on the first line that differs. It takes the place of today's two `adr.index` findings, a row with no file and a file with no row, which are cases of the same difference.
- A rebase that stops on the index writes it again from the rebased tree's ADR files, as ADR-0098 does for the summary. The index joins the generated files `issues.Generated` names. The row-by-row merge of ADR-0157 is removed. As before, a stop on any other path is left for the agent, and the index is written once nothing else is in conflict.
- The rest of ADR-0157 stands: the trial merge leaves the index out of a pair's conflicts. Because the index is in `issues.Generated`, `flai verify`'s patch id leaves it out with no further change ([ADR-0192](0192-flai-verify-s-patch-id-leaves-out-the-files-a-sync-s-rebase-writes-again-so-a.md)), and the merge queue's batch writes it again through the same steps (ADR-0175).

### 2. The status cell

The cell is these parts, in this order, each after the first preceded by a comma and a space:

1. The `status` word, always.
2. `supersedes <list>`, when `supersedes` is not empty.
3. `refines <list>`, when `refines` is not empty.
4. `superseded by <list>`, when `superseded_by` is not empty.

A list is the ADR numbers, four digits without `ADR-`, in ascending order, the last two joined by "and", the others by a comma, with no comma before "and": `0085, 0122 and 0133`.

| Front matter | Cell |
|---|---|
| `status: accepted` | accepted |
| `status: proposed`, `refines: [ADR-0130, ADR-0131]` | proposed, refines 0130 and 0131 |
| `status: accepted`, `supersedes: [ADR-0024]`, `refines: [ADR-0020, ADR-0029]` | accepted, supersedes 0024, refines 0020 and 0029 |
| `status: accepted`, `refines: [ADR-0029, ADR-0037]`, `superseded_by: [ADR-0041]` | accepted, refines 0029 and 0037, superseded by 0041 |
| `status: superseded`, `superseded_by: [ADR-0031]` | superseded, superseded by 0031 |

The cell never says "in part". The status word says whether the older ADR still stands: `accepted` with a `superseded_by` means the newer ADR replaces some or all of it, and the newer ADR says which. An ADR whose status is `superseded` is wholly replaced. No ADR here has that status today.

Notes the front matter does not hold are dropped at the first regeneration: the four "in part" notes, "refined by 0022", and the refines of 0022, 0023, and 0024. The three wrong statuses are corrected.

### 3. The issues summary is checked against its folder

`flai check` reports `issues.summary`, a warning, when `design/issues/summary.md` differs from what the issue files render. Its `updated:` line is left out of the comparison, since it records when the file was written.

`flai check --fix` rewrites the issues summary and the ADR index from their folders, then checks. It rewrites generated files only, and fixes no other finding.

### 4. Changelogs merge by union

- `.gitattributes` gives `CHANGELOG.md` `merge=union`, in this repository and in the template's `template/root/.gitattributes`. A merge or rebase keeps both sides' lines and never stops on the file.
- A story writes its changelog lines under a `## Unreleased` section at the top, never under a version heading.
- `flai release` names that section with the version and date, and leaves a fresh, empty `## Unreleased` above it.
- The digest takes no attribute. It has no index file ([design/digest/README.md](../digest/README.md)): each entry is a file of its own, so two stories never edit the same digest file.

### 5. The acceptance notice leaves out the shared paths

`flai accept`'s overlap notice to an open story lists only the changed paths outside every `claims.shared` pattern, read as ADR-0096 §2 reads them. An open story whose covered paths all lie inside shared patterns gets no notice. This reverses ADR-0096 §2 for the acceptance notice alone.

The trial merge at `flai stream sync` still reports shared paths, and so does `flai task done`'s message, which is about a task's commit while both stories are open. The ready-story hold and `flai serve` already skip shared paths (ADR-0096 §2), and still do.

The reason: the trial merge already decides whether a shared edit conflicts. A notice for a shared path that merges cleanly only makes an agent sync again for nothing.

## Consequences

- No story edits the ADR index or the issues summary by hand, and no rebase stops on either. Two stories that supersede the same ADR still stop on that ADR's `superseded_by`, which is data both sides write.
- Four "in part" notes, one "refined by", and three refines lists leave the index. A reader finds the extent of a supersession in the newer ADR.
- A hand edit to the table is undone at the next regeneration, and `flai check` reports it until then.
- Two stories that add changelog lines merge without a stop. Union keeps both sides of a line both changed, so a released entry is never rewritten, and a duplicated line is removed by hand.
- Fewer overlap notices per acceptance: in the fortnight before, 123 of the 190 notices named only shared paths, and would not have been sent. A real textual conflict on a shared path is still caught by the trial merge and the rebase at the next sync.
- S-0550's tasks implement each part, and the design in `design/system`, the conventions, and the guides describe them in the same story. The fortnight after landing is compared with the figures above.

## Alternatives considered

- **Keep the index hand-edited, merged row by row (ADR-0157).** It still stops a rebase on a status both sides change, and leaves three wrong statuses uncorrected. The notes it keeps are eight cells out of 223.
- **Read "accepted with a `superseded_by`" as "superseded in part".** It would mislabel the six ADRs superseded whole, which also keep `status: accepted`, and those files cannot be changed.
- **A new front matter key for a supersession in part.** The two ADRs it would mark are accepted and cannot gain the key. It adds a flag and a key for a fact the newer ADR already states in prose.
- **Keep hand notes the front matter lacks, by merging them into the generated cells.** That is ADR-0157's merge again, with its stops.
- **A custom merge driver for the index or the changelog.** Each clone must name it in its `git config`, which flai does not install (ADR-0098). `union` is built into git and needs no configuration.
- **Let `flai check --fix` fix every finding it can.** The other findings need judgement. Only a generated file has one right content.
- **Drop shared paths from the trial merge too.** It is what catches a real conflict on a shared file. It stays.
