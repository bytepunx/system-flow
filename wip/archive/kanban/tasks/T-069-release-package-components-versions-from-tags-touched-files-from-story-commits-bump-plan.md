---
id: T-069
type: task
nature: feature
title: "release package: components, versions from tags, touched files from story commits, bump plan"
status: done
parent: S-029
owner: alex
created: 2026-09-17T02:04:19Z
updated: 2026-09-17T02:09:09Z
transitions:
  - to: ready
    at: 2026-09-17T02:09:08Z
    by: agent
  - to: in-progress
    at: 2026-09-17T02:09:08Z
    by: agent
  - to: done
    at: 2026-09-17T02:09:09Z
    by: agent
stream: S-029
tags: [cli, release]
---

# T-069 release package: components, versions from tags, touched files from story commits, bump plan

## Work
internal/release: semver parse and bump; LevelFor (epic major, feature minor, remediation and improvement patch, research and experiment refuse); Commits by [ID] in messages; TouchedFiles via git diff-tree; CurrentVersion from <name>/v* tags (numeric sort) or template.yaml for template components; Compute maps touched files to manifest projects, picks the delivered component from --deliver, item tags, parent tags, or the only touched one, and plans delivered-level plus incidental patches; Apply bumps template version and changelog; Tag creates annotated tags on HEAD.

## Done when
Unit tests on a synthetic git repo cover every branch.

## Notes
