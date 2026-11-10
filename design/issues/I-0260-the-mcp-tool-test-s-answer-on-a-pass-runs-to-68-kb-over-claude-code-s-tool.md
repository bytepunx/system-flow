---
id: I-0260
title: The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-11-02T19:44:19Z
last_reported: 2026-11-02T19:44:19Z
updated: 2026-11-10T09:11:42Z
---

# I-0260 The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given

## Description
The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given

## Instances

### 2026-11-02T19:44:19Z
Story: S-0425.
S-0425: `test` on nine paths (seven Go packages, flaiover/src/lib, flaiover/src/routes/charts) passed every tier, but its answer was 68,784 characters, so Claude Code saved it to a file and the agent read it back with jq. The tiers' `command` fields carry the file list: gofmt 14 KB, flaiover-lint 10 KB, vitest 10 KB. A passing run needs only each tier's name, state, and duration. Afterwards the agent ran `scripts/flai.sh test --json` into a file and summarized it with jq.

## Remediation

Story S-0707 remediates this issue, created from it at 2026-11-10T09:11:42Z.
