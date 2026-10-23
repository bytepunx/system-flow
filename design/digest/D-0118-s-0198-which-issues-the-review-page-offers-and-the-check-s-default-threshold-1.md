---
id: D-0118
title: "S-0198: which issues the review page offers, and the check's default threshold: 1 - every open issue that this story touched should show first in the list with a checked checkbox; all issues that…"
kind: answer
sources: [TH-0074]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/issues, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/check, flai/internal/manifest, flai/internal/mcpserver, flai/internal/harness, flai/internal/hostapi, flaiover/src/lib, flaiover/src/routes/api/issues, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/system/project-manifest.md, design/system/flaiover-dashboard.md, design/conventions/continuous-improvement.md, template/root/design/conventions/continuous-improvement.md, design/system/conventions.md, docs/users/conventions.md]
date: 2026-10-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0118 S-0198: which issues the review page offers, and the check's default threshold: 1 - every open issue that this story touched should show first in the list with a checked checkbox; all issues that…

The operator, alex, answered TH-0074 at 2026-10-03T01:46:03Z.

What was asked:

```text
Two choices change the work. I am going on with everything else while they are open.
```

The answer:

```text
1 - every open issue that this story touched should show first in the list with a checked checkbox; all issues that don't have a story for their remediation or improvement should appear in the list, unchecked. If any issues are checked, the `Accept` button text changes to `Accept and Create Stories` in which the action should perform both tasks (create stories only takes place if the accept succeeds, otherwise, do not proceed to the create stories step). 2 - proceed with the recommendation AND create a story for issue 7 so that flai check does not fail.
```

Source: TH-0074.
