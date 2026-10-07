---
id: I-0058
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
class: efficiency
status: closed
count: 6
cost: 3m
first_reported: 2026-10-02T17:23:40Z
last_reported: 2026-10-06T22:01:45Z
updated: 2026-10-07T09:03:38Z
---

# I-0058 flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

## Description
flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

## Instances

### 2026-10-02T17:23:40Z
S-0230 T-0711 (2026-10-02): a task sub-agent's Bash call, a python heredoc editing a Go comment that read 'written as having spent nothing', was refused by flai guard as if it were a flai command that changes work items; the sub-agent made the edit with Edit instead. The guard matches words inside a heredoc's text, not only the command it runs.

### 2026-10-04T04:07:29Z
Story: S-0209.
T-0791's task sub-agent was refused a python heredoc editing guard_test.go: the text held apostrophes and the words flai story new, read as a flai write; it went on with the Edit tool.

### 2026-10-06T04:56:08Z
Story: S-0282.
T-0996's task sub-agent had a Bash heredoc refused because a '#' comment line inside it read as a flai command; it worked around it with direct edits.

### 2026-10-06T10:54:30Z
Story: S-0285.
T-1004's task sub-agent had a Bash heredoc refused: a Go comment inside it read 'flai guard records'. It wrote the file with the Edit tool instead.

### 2026-10-06T11:24:31Z
Story: S-0221.
T-0903's task sub-agent wrote a Go test through a bash heredoc whose comment named flai accept --by; the guard refused it, and it wrote the file with Edit instead. The orchestrator's acceptance (ADR-0093) passes its evidence through a heredoc too, so a Verdict line naming a flai write would be refused the same way.

### 2026-10-06T22:01:45Z
Story: S-0227.
T-0942's task sub-agent edited design/system/strategic-agents.md through a shell script; the guard refused it because the script's text held the words of a flai issue command. It redid the edits with the Edit tool.

## Remediation
Closed 2026-10-07T09:03:38Z: S-0246: flai guard reads a heredoc's body as input, not as commands, unless a shell on the line that opens it reads it (bash <<EOF, cat <<EOF | sh), so a sub-agent's heredoc whose text names a flai write is no longer refused
