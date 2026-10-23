---
id: D-0151
title: flai check --strict stops a story's close-out on wip/ findings outside the story
kind: cause
sources: [I-0057]
topics: [cli, conventions, template, flai, go, code, all]
paths: [flai/internal/check, flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/issues, flai/internal/storygit, scripts/close-out.sh, scripts/check.sh, scripts/README.md, template/root/scripts/close-out.sh, template/root/scripts/check.sh, template/root/scripts/README.md, template/CHANGELOG.md, template/template.yaml, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/adrs, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md, design/issues]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0151 flai check --strict stops a story's close-out on wip/ findings outside the story

I-0057 was closed on 2026-10-05.

The cause:

```text
flai check --strict stops a story's close-out on wip/ findings outside the story
```

The fix:

```text
S-0249: flai check --story S-nnnn reports the findings outside the story as notes that neither an error nor --strict fails on, and --record-issues records each rule's in one open issue, opened or bumped once per story and findings (ADR-0085). scripts/close-out.sh exports CLOSE_OUT_STORY, so scripts/check.sh and TestMonorepoIsClean scope the close-out's check to the story, here and in the template; CI's unscoped check still gates every finding.
```

Source: I-0057.
