---
id: T-046
type: task
nature: feature
title: gh repository creation and release bump rules in git.md
status: done
parent: S-010
owner: alex
created: 2026-09-15T23:01:12Z
updated: 2026-09-16T04:28:57Z
transitions:
  - to: ready
    at: 2026-09-15T23:01:12Z
    by: agent
  - to: in-progress
    at: 2026-09-15T23:01:12Z
    by: agent
  - to: done
    at: 2026-09-16T04:28:57Z
    by: agent
stream: S-010
tags: [conventions, release]
---

# T-046 gh repository creation and release bump rules in git.md

## Work
Add to git.md (template and here, above the marker): create remote repositories with gh when available, asking the operator for organization and visibility first; tag by delivery type: epic completion is a major release, a feature story a minor, remediation, improvement, and documentation-only stories a patch. Note the policy in design/tech/ci.md. Then create the GitHub repository for this monorepo with gh, push main, and push flai/v0.1.0 so the release workflow runs.

## Done when
Conventions updated in both places; repository exists on GitHub; the tag produced a release.

## Notes
