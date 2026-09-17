---
id: T-076
type: task
nature: feature
title: "publish package: clone remote branch, replace contents, commit, tag, push with verbatim git errors"
status: done
parent: S-021
owner: alex
created: 2026-09-17T03:34:08Z
updated: 2026-09-17T03:37:56Z
transitions:
  - to: ready
    at: 2026-09-17T03:37:55Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:37:55Z
    by: agent
  - to: done
    at: 2026-09-17T03:37:56Z
    by: agent
stream: S-021
tags: [cli, template]
---

# T-076 publish package: clone remote branch, replace contents, commit, tag, push with verbatim git errors

## Work
internal/publish: clone the remote into CacheDir/publish, checkout or create the branch, replace contents except .git, add, detect no changes, commit with the template version and a body naming the source, optional annotated tag v<version> refused when it exists on the remote, push branch and tag, --force for both; dry run reports files and message and removes the working clone; git errors returned verbatim.

## Done when
Tests on a bare remote cover dry run, push, nothing-to-push, tag, existing tag, branch creation, push after an external commit, missing remote, missing config.

## Notes
