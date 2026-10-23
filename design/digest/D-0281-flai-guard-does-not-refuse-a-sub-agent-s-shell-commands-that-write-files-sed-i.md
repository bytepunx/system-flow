---
id: D-0281
title: flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)
kind: cause
sources: [I-0052]
topics: [cli, conventions, flai, go, code, all]
paths: [design/adrs/README.md, design/system/agent-adapters.md, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/writes.go, flai/internal/guard/writes_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md, design/issues/summary.md, design/adrs/0141-flai-guard-refuses-a-read-only-sub-agent-a-shell-command-that-writes-a-file.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0281 flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

I-0052 was closed on 2026-10-12.

The cause:

```text
flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)
```

The fix:

```text
Fixed by S-0368 (ADR-0141). flai guard now refuses a read-only sub-agent a shell command that writes a file outside the temporary folder. A read-only sub-agent is one whose .claude/agents definition grants neither Edit nor Write, such as the explorer and the verifier, or Claude Code's built-in Explore and Plan. The refused forms are a redirection, tee, sed -i, perl -i, cp, mv, rm, touch, mkdir, rmdir, truncate, ln, install, patch, dd of=, and a formatter's write flag. TestAReadOnlySubAgentWritesNoFileThroughTheShell in flai/internal/guard/guard_test.go reproduces each command the issue names, and so do writes_test.go and TestGuardRefusesAReadOnlySubAgentsShellWrite in flai/cmd/guard_test.go. Task sub-agents still write through the shell. A write by an interpreter or a script is not seen, as ADR-0141 records.
```

Source: I-0052.
