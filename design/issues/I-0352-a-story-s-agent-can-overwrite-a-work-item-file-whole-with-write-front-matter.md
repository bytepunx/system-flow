---
id: I-0352
title: A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it
class: impression
status: open
count: 1
cost: 1m
first_reported: 2026-11-11T22:53:06Z
last_reported: 2026-11-11T22:53:06Z
updated: 2026-11-11T22:53:06Z
---

# I-0352 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it

## Description
A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it

## Instances

### 2026-11-11T22:53:06Z
Story: S-0717.
Found by T-3120. On 2026-11-10 at 11:03:59.967Z S-0665's agent overwrote the main checkout's T-3079 file, front matter included, with a Claude Code Write of the content PLACEHOLDER, six seconds after flai task new had created it, and wrote it back from .flai-cache/items-wip-kanban-tasks.json at 11:04:19.632Z (~/.flai/serve/agents/sf-S-0665-20261110T102103Z.log, lines 426 to 441). For those 20 seconds every Repo.List failed (I-0345, I-0344, before S-0716). flai guard and the tooling convention let a story's agent hand-edit an item's body; nothing refuses or warns of a Write that replaces a whole work item file. A guard rule that refuses a story's agent a Write to a file under wip/kanban that leaves no front matter, pointing it to Edit or flai edit, would stop the window at its source.

## Remediation
