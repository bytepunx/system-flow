---
id: I-0052
title: flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)
class: impression
status: closed
count: 1
cost: 0m
first_reported: 2026-10-01T11:13:30Z
last_reported: 2026-10-01T11:13:30Z
updated: 2026-10-12T17:37:21Z
---

# I-0052 flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

## Description
flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

## Instances

### 2026-10-01T11:13:30Z
S-0189's pre-review verifier found that criterion 3 (a verifier never makes the corrections it finds) is enforced for git and flai writes, MCP writes, and by the verifier having no Edit or Write, but a sub-agent with Bash could still write a file through the shell. No sub-agent has been seen doing it; ADR-0060 says the guard is not a shell.

## Remediation

Story S-0368 remediates this issue, created from it at 2026-10-12T16:52:41Z.
Closed 2026-10-12T17:37:21Z: Fixed by S-0368 (ADR-0141). flai guard now refuses a read-only sub-agent a shell command that writes a file outside the temporary folder. A read-only sub-agent is one whose .claude/agents definition grants neither Edit nor Write, such as the explorer and the verifier, or Claude Code's built-in Explore and Plan. The refused forms are a redirection, tee, sed -i, perl -i, cp, mv, rm, touch, mkdir, rmdir, truncate, ln, install, patch, dd of=, and a formatter's write flag. TestDecideReadOnlyShellWrites in flai/internal/guard/guard_test.go reproduces each command the issue names, and so do writes_test.go and the hook-JSON case in flai/cmd/guard_test.go. Task sub-agents still write through the shell. A write by an interpreter or a script is not seen, as ADR-0141 records.
