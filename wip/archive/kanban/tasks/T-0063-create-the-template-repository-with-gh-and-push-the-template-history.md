---
id: T-0063
type: task
nature: feature
title: Create the template repository with gh and push the template history
status: done
parent: S-0003
owner: alex
created: 2026-09-17T00:07:13Z
updated: 2026-09-17T00:11:37Z
transitions:
  - to: ready
    at: 2026-09-17T00:11:36Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:11:37Z
    by: agent
  - to: done
    at: 2026-09-17T00:11:37Z
    by: agent
stream: S-0003
tags: [template, devops]
---

# T-0063 Create the template repository with gh and push the template history

## Work
Ask the operator for visibility (org is bytepunx per the manifest); create github.com/bytepunx/system-flow-template with gh; push the history of ./template as its main branch using git subtree split (no rewrite of this repository); tag v1.0.0 from template.yaml.

## Done when
Repository exists with template.yaml and CHANGELOG.md at its root and a v1.0.0 tag.

## Notes
