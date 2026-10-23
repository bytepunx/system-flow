---
id: D-0219
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
kind: cause
sources: [I-0058]
topics: [cli, flai, go, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go, design/system/flai-cli.md, design/system/agent-context.md, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md, design/issues/summary.md, flai/cmd/guard.go, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0111-flai-check-finds-narrative-state-outside-the-story-at-close-out.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0219 flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

I-0058 was closed on 2026-10-07.

The cause:

```text
flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
```

The fix:

```text
S-0246: flai guard reads a heredoc's body as input, not as commands, unless a shell on the line that opens it reads it (bash <<EOF, cat <<EOF | sh), so a sub-agent's heredoc whose text names a flai write is no longer refused
```

Source: I-0058.
