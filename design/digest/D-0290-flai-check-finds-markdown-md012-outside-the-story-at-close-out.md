---
id: D-0290
title: "flai check finds `markdown.MD012` outside the story at close-out"
kind: cause
sources: [I-0139]
topics: [markdown, narrative, planning, flai, cli, go, code, all]
paths: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go, design/system/flai-cli.md, design/issues/I-0139-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md, design/issues/summary.md, docs/users/flai.md]
date: 2026-10-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0290 flai check finds `markdown.MD012` outside the story at close-out

I-0139 was closed on 2026-10-13.

The cause:

```text
flai check finds `markdown.MD012` outside the story at close-out
```

The fix:

```text
S-0376: flai stream answer (AnswerOpenQuestion) now leaves one blank line where it removes an open question and saves through the markdown lint guard, and flai archive writes each narrative it moves with every run of blank lines outside fenced code collapsed to one, so no archived narrative keeps an MD012 finding that no story's close-out can clear. TestAnswerOpenQuestionLeavesOneBlankLine and TestArchiveCollapsesANarrativesBlankRuns reproduce the S-0232 case.
```

Source: I-0139.
