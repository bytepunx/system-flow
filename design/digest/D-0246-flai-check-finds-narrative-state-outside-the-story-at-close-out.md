---
id: D-0246
title: "flai check finds `narrative.state` outside the story at close-out"
kind: cause
sources: [I-0109]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0109-flai-check-finds-narrative-state-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0111-flai-check-finds-narrative-state-outside-the-story-at-close-out.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0246 flai check finds `narrative.state` outside the story at close-out

I-0109 was closed on 2026-10-08.

The cause:

```text
flai check finds `narrative.state` outside the story at close-out
```

The fix:

```text
S-0323 (ADR-0125): a check scoped to a story now leaves out a narrative.state finding on the narrative of another story neither done nor cancelled, as it does a markdown finding there, so a close-out no longer records the placeholder Current state and Next steps of a story just started; that story's own close-out stops on them.
```

Source: I-0109.
