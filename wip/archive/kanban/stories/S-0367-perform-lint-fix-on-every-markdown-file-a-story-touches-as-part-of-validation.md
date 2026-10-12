---
id: S-0367
type: story
nature: improvement
title: Perform lint fix on every markdown file a story touches as part of validation
status: done
owner: alex
created: 2026-10-12T16:48:56Z
updated: 2026-10-12T17:44:30Z
transitions:
  - to: ready
    at: 2026-10-12T16:48:58Z
    by: alex
  - to: in-progress
    at: 2026-10-12T16:49:57Z
    by: agent-S-0367
  - to: review
    at: 2026-10-12T17:32:10Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T17:44:30Z
    by: alex
tags: [cli]
touches: [flai/cmd, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/verify/fix.go, flai/internal/verify/fix_test.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go, flai/internal/storygit/fix.go, flai/internal/storygit/fix_test.go, flai/internal/storygit/sync.go, design/adrs/0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md, design/adrs/README.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, system-flow.yaml, template/CHANGELOG.md, template/root/system-flow.yaml.tmpl, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, design/adrs/0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2545
  turns:
    - day: 2026-10-12
      ceremony: 3
      test_runs: 3
      hand_edits: 7
      work: 68
  models:
    - model: claude-opus-5-5
      input: 380
      output: 148086
      cache_read: 23402093
      cache_write: 647705
      cost: 11.6639
  strategic:
    - kind: orchestrator
      seconds: 669
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 135
          cache_read: 4255078
          cache_write: 12714
          cost: 1.0759
        - model: claude-sonnet-5-5
          input: 23
          output: 120
          cache_read: 416330
          cache_write: 92035
          cost: 0.4108
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-10-12T16:48:56Z
---
# S-0367 Perform lint fix on every markdown file a story touches as part of validation

## Goal

Run `scripts/lint-md.sh --fix [markdown file]` for every markdown file that gets touched as part of a story to prevent bad markdown from slipping into main and breaking build checks.

## Acceptance criteria
- [x] All markdown files get a lint fix call made on them if they were modified as part of a story
- [x] Any markdown files outside of a story that get added also get the lint fix call
- [x] Lint fix should get called right after a rebase to ensure that old files with bad markdown don't continue polluting recently changed markdown files
- [x] Markdown the fix cannot fix blocks the review from being accepted, naming the files, lines, and errors, so that the operator can correct them (TH-0413)

## Tasks
- T-1455 A test tier takes fix_command, and verify.Fix runs it over the files the tier selects
- T-1456 flai stream sync runs the fix over what the branch changes right after its rebase, and commits what it fixed
- T-1457 The manifest's markdown tier fixes with lint-md.sh --fix, and the docs and an ADR say what fix_command does
- T-1458 flai task done runs the fix over the paths it commits, before it commits them
- T-1459 flai accept runs the fix over the markdown the main checkout changed, before it commits the acceptance
- T-1470 verify runs the check of the tiers that fix over the given paths and answers what the fix left as findings
- T-1471 The docs and ADR-0139 say that flai accept refuses on markdown the fix left, naming its files, lines, and errors

## Notes

- The fix is a test tier's `fix_command`. This repository's `markdown` tier runs `[scripts/lint-md.sh, --fix, "{files}"]` (ADR-0139).
- When the fix runs:
  - `flai task done` runs it before it commits.
  - Every clean `flai stream sync`, and `flai accept`'s rebase, runs it right after the rebase over every file the branch changes, and commits what it rewrote.
  - `flai accept` runs it over the main checkout's uncommitted files.
- Criterion 2 is the markdown `flai accept` commits from the main checkout, as the operator confirmed on TH-0413.
- Criterion 4 is the operator's decision on TH-0413, added here (ADR-0140). The markdown editor they asked for alongside it is the draft story S-0378.
