---
id: T-0489
type: task
nature: feature
title: A topics package parses file topics and splits a document into sections with effective topics
status: done
parent: S-0134
owner: alex
created: 2026-09-27T03:58:46Z
updated: 2026-09-27T03:59:53Z
transitions:
  - to: ready
    at: 2026-09-27T03:58:54Z
    by: agent-S-0134
  - to: in-progress
    at: 2026-09-27T03:58:55Z
    by: agent-S-0134
  - to: done
    at: 2026-09-27T03:59:53Z
    by: agent-S-0134
stream: S-0134
tags: []
touches: [flai/internal/topics]
---
# T-0489 A topics package parses file topics and splits a document into sections with effective topics

## Work

New package flai/internal/topics: file topics from front matter, heading comments `<!-- topics: a, b -->`, sections by heading with heading path and effective topics (own, else parent's, else the file's; the caller gives the default, [all] for conventions). Headings inside fenced code are not headings.

## Done when

Behavior tests cover nesting, inheritance, a heading comment at the end of the line, fenced code, and a file without topics; make test passes.

## Notes
