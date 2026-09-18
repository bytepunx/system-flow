---
id: I-0019
title: Acceptance from the dashboard is refused for files only the host's global gitignore hides, and the preview does not show it
class: defect
status: open
count: 1
cost: 8m
first_reported: 2026-09-18T20:54:17Z
last_reported: 2026-09-18T20:54:17Z
updated: 2026-09-18T20:54:17Z
---

# I-0019 Acceptance from the dashboard is refused for files only the host's global gitignore hides, and the preview does not show it

## Description
Acceptance from the dashboard is refused for files only the host's global gitignore hides, and the preview does not show it

## Instances

### 2026-09-18T20:54:17Z
S-0048: the operator dropped the card on done and the move was refused: 'working tree has uncommitted changes outside wip (.claude/)'. The file, .claude/settings.local.json, was created by the agent this session and is ignored on the host by ~/.config/git/ignore. The container has no global gitignore (HOME=/tmp, no core.excludesFile), so git there reports .claude/ as untracked. The acceptance preview did not warn, because the dashboard runs it with --yes while the move that follows runs without, and the dashboard has no way to pass --yes. Unblocked for this clone by adding the file to .git/info/exclude.

## Remediation
Queued as S-0051: the preview reports uncommitted paths outside wip and the confirmation lets the designer include them or cancel; the container sees the ignore rules the host does; the template's `.gitignore` covers agent-local files such as `.claude/settings.local.json`.
