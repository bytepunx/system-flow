---
id: D-0280
title: "S-0367 criterion 2: which markdown files are \"outside of a story\": that's correct"
kind: answer
sources: [TH-0413]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/cmd, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/verify/fix.go, flai/internal/verify/fix_test.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go, flai/internal/storygit/fix.go, flai/internal/storygit/fix_test.go, flai/internal/storygit/sync.go, design/adrs/0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md, design/adrs/README.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, system-flow.yaml, template/CHANGELOG.md, template/root/system-flow.yaml.tmpl, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, design/adrs/0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0280 S-0367 criterion 2: which markdown files are "outside of a story": that's correct

The operator, alex, answered TH-0413 at 2026-10-12T16:58:39Z.

What was asked:

```text
**Recommendation:** read "markdown files outside of a story that get added" as the markdown that `flai accept` commits from the main checkout, the `wip/` items, narratives, threads, and messages flai writes there. `flai accept` runs the markdown fix over the main checkout's changed markdown just before its `git add -A` and commit.
```

The answer:

```text
that's correct. I do have a one adjustment, however. If the markdown can't be fixed, it should block the review being accepted with the files, lines, and errors so that the user can correct them. Lets also make a story that adds a markdown editor that can be launched from inlined file links in reports and other areas of the site so that the operator can correct the issue and unblock the review. I understand this is a significant increase in scope, so the editor can become a separate story for the planner.
```

Source: TH-0413.
