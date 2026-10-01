---
id: I-0052
title: flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)
class: impression
status: open
count: 1
cost: 0m
first_reported: 2026-10-01T11:13:30Z
last_reported: 2026-10-01T11:13:30Z
updated: 2026-10-01T11:13:30Z
---

# I-0052 flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

## Description
flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

## Instances

### 2026-10-01T11:13:30Z
S-0189's pre-review verifier found that criterion 3 (a verifier never makes the corrections it finds) is enforced for git and flai writes, MCP writes, and by the verifier having no Edit or Write, but a sub-agent with Bash could still write a file through the shell. No sub-agent has been seen doing it; ADR-0060 says the guard is not a shell.

## Remediation
