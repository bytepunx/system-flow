---
id: T-0431
type: task
nature: feature
title: The repo_url prompt defaults to the origin remote as an https URL and refuses a value that is not a URL
status: done
parent: S-0117
owner: alex
created: 2026-09-26T05:25:44Z
updated: 2026-09-26T05:34:41Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:14Z
    by: agent-S-0117
  - to: in-progress
    at: 2026-09-26T05:32:26Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T05:34:41Z
    by: agent-S-0117
stream: S-0117
tags: []
---

# T-0431 The repo_url prompt defaults to the origin remote as an https URL and refuses a value that is not a URL

## Work

The template's `repo_url` default is the origin remote made an https URL (`git@github.com:owner/repo.git` and `ssh://` forms become `https://github.com/owner/repo`), for `flai new` and `flai import`. A value typed at the prompt or given with `--var repo_url=` that is not empty and not an http or https URL with a host and a path is refused with the reason. `https://github.com:owner/repo` is refused.

## Done when

Tests cover the conversion of the common remote forms and the refusal; `make test` passes.

## Notes
