---
id: I-0231
title: Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-10-29T14:15:00Z
last_reported: 2026-10-29T14:15:00Z
updated: 2026-10-29T14:15:00Z
---

# I-0231 Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes

## Description
Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes

## Instances

### 2026-10-29T14:15:00Z
Story: S-0536.
S-0536: Edit on wip/kanban/stories/S-0536-the-orchestrator-runs-for-a-reason-and-ends-...-adr-0209-is-accepted.md (a 230-byte name) failed with ENAMETOOLONG. The tool writes `<name>.tmp.<pid>.<hex>` beside the file, which is over 255 bytes. Item bodies are to be hand-edited (tooling.md), so the agent rewrote criterion 6 and Notes with a Python script instead. flai names items from their titles, which can be long, and nothing caps the slug so that a hand edit's temporary file still fits.

## Remediation
