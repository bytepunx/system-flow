---
id: T-0493
type: task
nature: feature
title: Rules and docs say an accepted ADR may gain topics, describe the field and the heading comment, and document flai adr topics
status: done
parent: S-0134
owner: alex
created: 2026-09-27T03:58:48Z
updated: 2026-09-27T04:09:21Z
transitions:
  - to: ready
    at: 2026-09-27T03:58:55Z
    by: agent-S-0134
  - to: in-progress
    at: 2026-09-27T04:05:04Z
    by: agent-S-0134
  - to: done
    at: 2026-09-27T04:09:21Z
    by: agent-S-0134
stream: S-0134
tags: []
touches: [design/conventions, template/root/design/conventions, CLAUDE.md, template/root/CLAUDE.md.tmpl, design/system, docs/users/flai.md]
---
# T-0493 Rules and docs say an accepted ADR may gain topics, describe the field and the heading comment, and document flai adr topics

## Work

decisions.md (template, then here), CLAUDE.md and its template, design/system/documentation-standard.md allow topics besides superseded_by; design/system/conventions.md describes topics and the heading comment; docs/users/flai.md documents flai adr topics; regenerate the reference if it is generated.

## Done when

make test, make integration, make smoke, lint-md, and flai check --strict pass.

## Notes
