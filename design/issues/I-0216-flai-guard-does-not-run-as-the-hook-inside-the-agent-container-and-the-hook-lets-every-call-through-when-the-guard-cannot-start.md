---
id: I-0216
title: flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-10-28T12:30:51Z
last_reported: 2026-10-28T12:30:51Z
updated: 2026-10-28T12:30:51Z
---

# I-0216 flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start

## Description
flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start

## Instances

### 2026-10-28T12:30:51Z
Story: S-0505.
S-0505 T-2166: agent-S-0513 worked S-0513 in flai-agent-claude-code:local. This repository's .claude/settings.json runs the PreToolUse hook `out=$("$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard 2>&1); [ $? -eq 2 ] || exit 0; ...`. scripts/flai.sh builds bin/flai from the tree with go, which the image lacks, so the guard exits 127 (`go: not found`) on every call and the hook passes it on as exit 0: no call was checked, and the agent saw no hook error. The image carries flai on PATH, so a project made from the template, whose hook runs the installed `flai guard`, would run it; this repository's does not. The fail-open shape is the template's as well: any exit but 2, a missing binary included, lets the call through. Remediation: the container adapter (S-0508) runs the guard flai chooses, the image's flai, not the project's hook command, and a guard that cannot start refuses rather than allows; S-0510's layer would give this repository's hook its go.

## Remediation
