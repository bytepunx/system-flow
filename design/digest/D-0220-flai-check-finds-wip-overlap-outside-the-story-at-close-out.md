---
id: D-0220
title: "flai check finds `wip.overlap` outside the story at close-out"
kind: cause
sources: [I-0076]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/summary.md, design/system/work-hierarchy.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0220 flai check finds `wip.overlap` outside the story at close-out

I-0076 was closed on 2026-10-07.

The cause:

```text
flai check finds `wip.overlap` outside the story at close-out
```

The fix:

```text
S-0279 fixed it as ADR-0115 decides: a close-out's scoped flai check keeps only a wip.overlap that names its story, as a note, and --record-issues records no wip.overlap (T-1145); wip.overlap compares the claims of the stories in progress, as the pull hold does, one finding per pair of stories (T-1143).
```

Source: I-0076.
