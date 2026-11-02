---
id: I-0216
title: flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start
class: defect
status: closed
count: 2
cost: 10m
first_reported: 2026-10-28T12:30:51Z
last_reported: 2026-11-02T19:06:46Z
updated: 2026-11-02T20:35:21Z
---

# I-0216 flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start

## Description
flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start

## Instances

### 2026-10-28T12:30:51Z
Story: S-0505.
S-0505 T-2166: agent-S-0513 worked S-0513 in flai-agent-claude-code:local. This repository's .claude/settings.json runs the PreToolUse hook `out=$("$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard 2>&1); [ $? -eq 2 ] || exit 0; ...`. scripts/flai.sh builds bin/flai from the tree with go, which the image lacks, so the guard exits 127 (`go: not found`) on every call and the hook passes it on as exit 0: no call was checked, and the agent saw no hook error. The image carries flai on PATH, so a project made from the template, whose hook runs the installed `flai guard`, would run it; this repository's does not. The fail-open shape is the template's as well: any exit but 2, a missing binary included, lets the call through. Remediation: the container adapter (S-0508) runs the guard flai chooses, the image's flai, not the project's hook command, and a guard that cannot start refuses rather than allows; S-0510's layer would give this repository's hook its go.

### 2026-11-02T19:06:46Z
Story: S-0509.
S-0509 runs the planner, the orchestrator, the analyzer, and the curator in the same container, with the guard as the inner Claude Code's hook from the mounted checkout's `.claude/settings.json`. In this repository that hook runs `scripts/flai.sh guard`, which needs `go`, absent from the image. So a strategic agent's guard does not start there either. ADR-0247 rests each role's bounds on that guard, so this issue now covers strategic agents as well as stories.

## Remediation

Story S-0518 remediates this issue, created from it at 2026-10-28T12:52:56Z.
Closed 2026-11-02T20:35:21Z: S-0518 fixed it (ADR-0235): the PreToolUse guard hooks in this repository's and the template's .claude/settings.json refuse a call with 'flai guard did not run (exit N)' when the guard exits 126 or 127, not found or not executable, where they let every exit but 2 through; and scripts/flai.sh, with no go, runs bin/flai as it is, else the installed flai with a warning, else exits 127, so inside the agent image this repository's hook reaches the image's flai. That covers S-0509's instance too: a strategic agent's container runs the same hook from the mounted checkout. TestTheGuardHooksRefuseACallWhenTheGuardCannotStart reproduces the case; the check inside flai-agent-claude-code:local is in S-0518's narrative.
