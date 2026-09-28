---
id: T-0491
type: task
nature: feature
title: flai check warns on design and tech files without topics and on topics nothing uses
status: done
parent: S-0134
owner: alex
created: 2026-09-27T03:58:47Z
updated: 2026-09-27T04:04:14Z
transitions:
  - to: ready
    at: 2026-09-27T03:58:55Z
    by: agent-S-0134
  - to: in-progress
    at: 2026-09-27T04:01:36Z
    by: agent-S-0134
  - to: done
    at: 2026-09-27T04:04:14Z
    by: agent-S-0134
stream: S-0134
tags: []
touches: [flai/internal/check, design/system, design/tech, template/root/design]
---
# T-0491 flai check warns on design and tech files without topics and on topics nothing uses

## Work

check: doc.topics warns on a design/system or design/tech file without topics; doc.topic warns on a topic in a convention, design, tech, or ADR file (front matter or heading) outside the vocabulary: every sub-project's name, tags, and kind, code, all. Give every design/system and design/tech file here, and the template's overview, topics so the strict check stays clean.

## Done when

Check tests cover both rules including a heading topic; flai check --strict passes here and on a rendered template.

## Notes
