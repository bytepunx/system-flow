---
id: D-0242
title: "flai check finds `markdown.MD038` outside the story at close-out"
kind: cause
sources: [I-0096]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0242 flai check finds `markdown.MD038` outside the story at close-out

I-0096 was closed on 2026-10-08.

The cause:

```text
flai check finds `markdown.MD038` outside the story at close-out
```

The fix:

```text
Fixed by S-0318. The cause: S-0229's agent hand-edited its narrative's ## Decisions with a code span ending in a space, which broke MD038, and S-0227's close-out recorded it, though only S-0229 could fix it. ADR-0123 decides that a check scoped to a story leaves out a markdown finding on the narrative of another story that is neither done nor cancelled, so a close-out records none; that story's own close-out still finds it, and the unscoped check still warns. TestScopedCheckLeavesOutAMarkdownFindingOnAnotherOpenStorysNarrative and TestScopedCheckNotesEveryOtherMarkdownFindingOutsideTheStory in flai/internal/check/scope_test.go, and TestCheckRecordIssuesLeavesOutAMarkdownFindingOnAnotherOpenStorysNarrative in flai/cmd/check_test.go, reproduce it.
```

Source: I-0096.
