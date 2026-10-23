---
id: D-0217
title: golangci-lint fails at once when another story's agent is running it on the same host
kind: cause
sources: [I-0101]
topics: [flai, cli, go, code, all]
paths: [flai/.golangci.yaml, flai/tests/integration/golangci_lock_test.go, design/system/devex.md, design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md, design/issues/summary.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0217 golangci-lint fails at once when another story's agent is running it on the same host

I-0101 was closed on 2026-10-07.

The cause:

```text
golangci-lint fails at once when another story's agent is running it on the same host
```

The fix:

```text
S-0308: flai/.golangci.yaml sets run.allow-parallel-runners, so a lint run beside another on the host no longer exits 3 on the shared lock; flai/tests/integration/golangci_lock_test.go holds the lock and lints with the config
```

Source: I-0101.
