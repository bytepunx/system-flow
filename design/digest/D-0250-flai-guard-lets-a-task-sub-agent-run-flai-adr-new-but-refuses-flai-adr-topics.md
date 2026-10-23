---
id: D-0250
title: flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
kind: cause
sources: [I-0062]
topics: [cli, conventions, template, flai, go, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go, design/system/agent-context.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, docs/users/flai.md, design/issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md, design/issues/summary.md, design/adrs/0127-flai-guard-lets-a-story-s-sub-agent-write-an-adr-with-flai-adr-new-topics-and.md, design/adrs/README.md, template/CHANGELOG.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0250 flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

I-0062 was closed on 2026-10-08.

The cause:

```text
flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
```

The fix:

```text
S-0287 fixed it (ADR-0127): flai guard now lets a sub-agent run flai adr new, flai adr topics, and flai adr accept, and call adr_new, so a task sub-agent records its ADR through flai rather than copying internal/adr by hand or guessing the number. It refuses only --commit, --autocommit, and adr_new's commit, since the commit stays the story's agent's; TestGuardLetsASubAgentWriteAnADRButNotCommitIt in flai/cmd/guard_test.go reproduces the instances. delegation.md and its template copy say so.
```

Source: I-0062.
