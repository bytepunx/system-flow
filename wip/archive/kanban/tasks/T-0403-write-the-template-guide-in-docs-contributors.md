---
id: T-0403
type: task
nature: feature
title: Write the template guide in docs/contributors
status: done
parent: S-0018
owner: alex
created: 2026-09-24T08:12:56Z
updated: 2026-09-24T08:14:59Z
transitions:
  - to: ready
    at: 2026-09-24T08:13:18Z
    by: agent-S-0018
  - to: in-progress
    at: 2026-09-24T08:13:18Z
    by: agent-S-0018
  - to: done
    at: 2026-09-24T08:14:59Z
    by: agent-S-0018
stream: S-0018
tags: []
---

# T-0403 Write the template guide in docs/contributors

## Work
Write `docs/contributors/template.md` for someone forking or editing the template: what is in the template repository, `template.yaml` field by field, variables and how their values are chosen, the built-in data and functions, the render rules, item templates, pointing flai at a fork, versioning and publishing, and testing. Every statement checked against `flai/internal/template` and `flai/cmd`.

## Done when
The guide covers template.yaml, variables, render rules, and testing, and each claim matches the code in this tree.

## Notes
