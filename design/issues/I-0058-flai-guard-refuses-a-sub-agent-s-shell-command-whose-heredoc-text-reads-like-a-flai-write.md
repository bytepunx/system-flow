---
id: I-0058
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
class: efficiency
status: open
count: 1
cost: 2m
first_reported: 2026-10-02T17:23:40Z
last_reported: 2026-10-02T17:23:40Z
updated: 2026-10-02T17:23:40Z
---

# I-0058 flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

## Description
flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

## Instances

### 2026-10-02T17:23:40Z
S-0230 T-0711 (2026-10-02): a task sub-agent's Bash call, a python heredoc editing a Go comment that read 'written as having spent nothing', was refused by flai guard as if it were a flai command that changes work items; the sub-agent made the edit with Edit instead. The guard matches words inside a heredoc's text, not only the command it runs.

## Remediation
