---
id: T-0077
type: task
nature: feature
title: flai template push command and publish section in template.yaml
status: done
parent: S-0021
owner: alex
created: 2026-09-17T03:34:08Z
updated: 2026-09-17T03:37:56Z
transitions:
  - to: ready
    at: 2026-09-17T03:37:56Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:37:56Z
    by: agent
  - to: done
    at: 2026-09-17T03:37:56Z
    by: agent
stream: S-0021
tags: [cli, template]
---

# T-0077 flai template push command and publish section in template.yaml

## Work
template.Manifest gains publish.repo and publish.ref; the prototype's template.yaml names bytepunx/system-flow-template main; flai template push [dir] [--remote] [--ref] [--tag] [--force] [--dry-run] [--json] with the directory defaulting to config template.repo when local.

## Done when
In-process tests cover dry run, push with tag and branch creation, nothing to push, refused tag, and the config-remote error.

## Notes
