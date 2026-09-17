---
id: T-066
type: task
nature: feature
title: "importer package: scan, detect folders and sub-projects, plan"
status: done
parent: S-006
owner: alex
created: 2026-09-17T00:15:09Z
updated: 2026-09-17T00:19:50Z
transitions:
  - to: ready
    at: 2026-09-17T00:19:49Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:19:49Z
    by: agent
  - to: done
    at: 2026-09-17T00:19:50Z
    by: agent
stream: S-006
tags: [cli, import]
---

# T-066 importer package: scan, detect folders and sub-projects, plan

## Work
internal/importer: Scan walks the repo (skipping .git, node_modules, vendor, and detected sub-projects) and reports existing layout folders, candidate documentation folders (adr, adrs, architecture, documentation, docs, design, wip), loose markdown outside the future structure, sub-projects by build file (go.mod, package.json with svelte.config for sveltekit, pyproject.toml, Cargo.toml), and whether the tree is a git repo. Plan turns the scan and chosen layout names into a proposal: folders to create, folder moves (adr* to design/adrs, documentation to docs/users), markdown to ask about, projects to record.

## Done when
Unit tests on a synthetic tree cover detection and the plan.

## Notes
