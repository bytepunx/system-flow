---
id: T-0492
type: task
nature: feature
title: "Every convention in the template and here carries topics: [all], and the template's version records it"
status: done
parent: S-0134
owner: alex
created: 2026-09-27T03:58:47Z
updated: 2026-09-27T04:05:04Z
transitions:
  - to: ready
    at: 2026-09-27T03:58:55Z
    by: agent-S-0134
  - to: in-progress
    at: 2026-09-27T04:04:14Z
    by: agent-S-0134
  - to: done
    at: 2026-09-27T04:05:04Z
    by: agent-S-0134
stream: S-0134
tags: []
touches: [template/root/design/conventions, design/conventions, template/template.yaml, template/CHANGELOG.md]
---
# T-0492 Every convention in the template and here carries topics: [all], and the template's version records it

## Work

Add topics: [all] to every convention in template/root/design/conventions and design/conventions; bump template/template.yaml and add a CHANGELOG entry; a test shows flai upgrade brings the key to a project.

## Done when

grep shows topics: [all] in every file of both folders; the upgrade test passes; make smoke passes.

## Notes
