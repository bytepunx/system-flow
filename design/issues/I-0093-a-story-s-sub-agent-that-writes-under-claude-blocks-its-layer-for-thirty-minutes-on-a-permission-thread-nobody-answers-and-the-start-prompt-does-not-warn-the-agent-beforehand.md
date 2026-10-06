---
id: I-0093
title: A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand
class: blocker
status: open
count: 1
cost: 6m
first_reported: 2026-10-06T20:53:14Z
last_reported: 2026-10-06T20:53:14Z
updated: 2026-10-06T20:53:14Z
---

# I-0093 A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

## Description
A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

## Instances

### 2026-10-06T20:53:14Z
Story: S-0223.
T-0957's sub-agent called Write on template/root/.claude/agents/analyzer.md and Edit on template/root/.claude/settings.json. The installed flai's permission_prompt opened TH-0192 and held the call. The owner was away, so the layer waited until the board watcher stopped and restarted the session about six minutes later. The warning on TH-0190 reached the agent after it had launched the layer. The fix was to write the files into .flai-cache/S-0223/ and ask alex to cp them in on TH-0196.

## Remediation
