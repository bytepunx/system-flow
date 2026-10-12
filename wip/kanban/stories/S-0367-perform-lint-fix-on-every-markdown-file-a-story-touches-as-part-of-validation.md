---
id: S-0367
type: story
nature: improvement
title: Perform lint fix on every markdown file a story touches as part of validation
status: in-progress
owner: alex
created: 2026-10-12T16:48:56Z
updated: 2026-10-12T17:09:12Z
transitions:
  - to: ready
    at: 2026-10-12T16:48:58Z
    by: alex
  - to: in-progress
    at: 2026-10-12T16:49:57Z
    by: agent-S-0367
tags: [cli]
touches: [flai/cmd, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/verify/fix.go, flai/internal/verify/fix_test.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go, flai/internal/storygit/fix.go, flai/internal/storygit/fix_test.go, flai/internal/storygit/sync.go, design/adrs/0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md, design/adrs/README.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, system-flow.yaml, template/CHANGELOG.md, template/root/system-flow.yaml.tmpl, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1779
  estimated: true
  turns:
    - day: 2026-10-12
      ceremony: 2
      test_runs: 2
      hand_edits: 5
      work: 53
  models:
    - model: claude-opus-5-5
      input: 332
      output: 1941
      cache_read: 17839865
      cache_write: 616364
      cost: 8.3361
  strategic:
    - kind: orchestrator
      seconds: 443
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 110
          cache_read: 2015371
          cache_write: 4401
          cost: 0.5092
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
- [ ] All markdown files get a lint fix call made on them if they were modified as part of a story
- [ ] Any markdown files outside of a story that get added also get the lint fix call
- [ ] Lint fix should get called right after a rebase to ensure that old files with bad markdown don't continue polluting recently changed markdown files

## Tasks
- T-1455 A test tier takes fix_command, and verify.Fix runs it over the files the tier selects
- T-1456 flai stream sync runs the fix over what the branch changes right after its rebase, and commits what it fixed
- T-1457 The manifest's markdown tier fixes with lint-md.sh --fix, and the docs and an ADR say what fix_command does
- T-1458 flai task done runs the fix over the paths it commits, before it commits them
- T-1459 flai accept runs the fix over the markdown the main checkout changed, before it commits the acceptance
- T-1470 verify runs the check of the tiers that fix over the given paths and answers what the fix left as findings
- T-1471 The docs and ADR-0139 say that flai accept refuses on markdown the fix left, naming its files, lines, and errors

## Notes
