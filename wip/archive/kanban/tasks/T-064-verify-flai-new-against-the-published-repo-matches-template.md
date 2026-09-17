---
id: T-064
type: task
nature: feature
title: Verify flai new against the published repo matches ./template
status: done
parent: S-003
owner: alex
created: 2026-09-17T00:07:13Z
updated: 2026-09-17T00:11:37Z
transitions:
  - to: ready
    at: 2026-09-17T00:11:37Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:11:37Z
    by: agent
  - to: done
    at: 2026-09-17T00:11:37Z
    by: agent
stream: S-003
tags: [template]
---

# T-064 Verify flai new against the published repo matches ./template

## Work
Render a project from the published repo at v1.0.0 and from ./template with the same variables; diff the two trees; run flai check --strict on both.

## Done when
Trees identical apart from template.repo in the manifest; both pass check.

## Notes
