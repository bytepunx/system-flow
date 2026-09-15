---
id: T-014
type: task
nature: feature
title: flai new command
status: done
parent: S-005
owner: agent
created: 2026-09-15T16:36:28Z
updated: 2026-09-15T16:38:47Z
transitions:
  - to: ready
    at: 2026-09-15T16:36:28Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:38:47Z
    by: agent
  - to: done
    at: 2026-09-15T16:38:47Z
    by: agent
stream: S-005
tags: [cli, template]
---

# T-014 flai new command

## Work
cmd/new.go: flai new <dir> [--template] [--ref] [--var k=v] [--defaults] [--force] [--no-git]. Variable defaults are templates evaluated against known variables; project_name defaults to the directory name. Prompt with huh only when stdin is a terminal and values are missing; error otherwise. git init the result unless --no-git. Print a summary and next steps.

## Done when
In-process CLI test renders a project with --defaults and --var and verifies system-flow.yaml parses.

## Notes
