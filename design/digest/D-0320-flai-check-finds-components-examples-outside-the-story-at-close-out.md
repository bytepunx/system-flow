---
id: D-0320
title: "flai check finds `components.examples` outside the story at close-out"
kind: cause
sources: [I-0161]
topics: [components, close-out, flai, cli, go, template, conventions, code, all]
paths: [flai/internal/check/check.go, flai/internal/check/components_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/system/component-library.md, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/issues/I-0161-flai-check-finds-components-examples-outside-the-story-at-close-out.md, design/issues/summary.md, design/system/workflow.md, flai/cmd/check_stats_test.go, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, docs/users/flai-reference.md]
date: 2026-10-18
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0320 flai check finds `components.examples` outside the story at close-out

I-0161 was closed on 2026-10-18.

The cause:

```text
flai check finds `components.examples` outside the story at close-out
```

The fix:

```text
S-0446 (ADR-0171): a sub-project whose svelte is not installed gets one advisory components.examples note on its examples folder instead of a warning on each examples file, and a check scoped to a story leaves out that note, and the one for a host without node, unless the story changes a file under the folder it names, so the close-out records nothing for them. Tests: TestComponentExamplesWithoutSvelteGetOneNote, TestScopedCheckLeavesOutAnUncheckedExamplesNoteOutsideTheStory, TestCheckRecordIssuesLeavesOutAnUncheckedExamplesNoteOutsideTheStory.
```

Source: I-0161.
