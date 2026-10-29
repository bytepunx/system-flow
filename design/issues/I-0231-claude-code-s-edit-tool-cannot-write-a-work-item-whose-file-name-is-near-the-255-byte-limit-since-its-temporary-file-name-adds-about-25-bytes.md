---
id: I-0231
title: Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes
class: efficiency
status: closed
count: 1
cost: 3m
first_reported: 2026-10-29T14:15:00Z
last_reported: 2026-10-29T14:15:00Z
updated: 2026-10-29T15:32:41Z
---

# I-0231 Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes

## Description
Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes

## Instances

### 2026-10-29T14:15:00Z
Story: S-0536.
S-0536: Edit on wip/kanban/stories/S-0536-the-orchestrator-runs-for-a-reason-and-ends-...-adr-0209-is-accepted.md (a 230-byte name) failed with ENAMETOOLONG. The tool writes `<name>.tmp.<pid>.<hex>` beside the file, which is over 255 bytes. Item bodies are to be hand-edited (tooling.md), so the agent rewrote criterion 6 and Notes with a Python script instead. flai names items from their titles, which can be long, and nothing caps the slug so that a hand edit's temporary file still fits.

## Remediation

Story S-0542 remediates this issue, created from it at 2026-10-29T15:14:37Z.
Closed 2026-10-29T15:32:41Z: Fixed by S-0542 (T-2333, T-2334). Every file name flai makes from a title now cuts the title's slug at 80 bytes on a word, through one helper, template.CutSlug with template.SlugMax in flai/internal/template/render.go: work items (workitem.FileName, so flai story, task, and epic new, a retitle with flai edit, and an experiment's expected results document), issues, threads, and conversations, and, as before, ADRs and digest entries, whose hand-written cuts moved onto the helper. A name is then at most the ID, a dash, 80 bytes, and .md, so Claude Code's <name>.tmp.<pid>.<hex> beside it stays under the 255-byte limit. TestCreateCutsALongTitlesFileName in flai/internal/workitem/create_test.go reproduces the instance with S-0536's title; issues and threads have their own tests. Files named before keep their names until retitled. Documented in design/system/documentation-standard.md § Naming and docs/users/flai.md § Work items.
