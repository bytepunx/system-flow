---
id: D-0273
title: "flai check finds `issues.no-story` outside the story at close-out"
kind: cause
sources: [I-0134]
topics: [flai, cli, go, template, conventions, code, all]
paths: [design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md, design/issues/summary.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0273 flai check finds `issues.no-story` outside the story at close-out

I-0134 was closed on 2026-10-12.

The cause:

```text
flai check finds `issues.no-story` outside the story at close-out
```

The fix:

```text
Fixed by S-0365 (ADR-0136): a check scoped to a story, as the close-out runs it, leaves out every issues.no-story outside the story and takes back its counts, so --record-issues no longer records an open issue no story links, such as I-0052; one on an issue file the story changes stays inside, and the main checkout's flai check still warns on every one. TestScopedCheckLeavesOutAnIssuesNoStoryOutsideTheStory and TestCheckRecordIssuesLeavesOutAnIssuesNoStoryOutsideTheStory reproduce it.
```

Source: I-0134.
