---
id: I-0275
title: This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would
class: impression
status: closed
count: 1
first_reported: 2026-11-07T05:43:56Z
last_reported: 2026-11-07T05:43:56Z
updated: 2026-11-10T11:24:32Z
---

# I-0275 This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would

## Description
This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would

## Instances

### 2026-11-07T05:43:56Z
Story: S-0567.
S-0567's live planner run: inside the container, this repository's PreToolUse hooks in .claude/settings.json run scripts/flai.sh guard, which builds bin/flai with go when the sources are newer than it (the layer image has go). The hooks refuse only on the guard's exit 2, and on 126 or 127 since ADR-0235; any other exit, such as 1 from scripts/flai-build.sh when go build fails under set -eu, lets the call through. The run's log cannot show which way its hooks went: the planner made no call the guard refuses, and no hook output is logged. By hand, the image's own flai guard refuses the planner a Write, git commit, and a move to ready exactly as the host does.

## Remediation

Story S-0666 remediates this issue, created from it at 2026-11-10T08:34:38Z.
Closed 2026-11-10T11:24:32Z: S-0666 (ADR-0280): when the build of bin/flai fails, scripts/flai.sh guard runs the older bin/flai, else the installed flai, and exits 127 with neither, so the guard hooks in .claude/settings.json refuse rather than pass the call; every other command still exits with the build's status. This repository's OpenCode plugin treats exit 126 or 127 as a guard that did not run. TestTheGuardHooksRefuseACallWhenTheBuildFails reproduces the failed build with a fake go, and TestRepoFlaiGuardPluginJudgesGuardExit covers the plugin.
