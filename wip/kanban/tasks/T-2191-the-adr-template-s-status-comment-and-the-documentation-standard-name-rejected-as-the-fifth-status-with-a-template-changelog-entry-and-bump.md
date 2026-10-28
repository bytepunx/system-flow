---
id: T-2191
type: task
nature: remediation
title: The ADR template's status comment and the documentation standard name rejected as the fifth status, with a template changelog entry and bump
status: backlog
parent: S-0512
owner: alex
created: 2026-10-28T08:27:46Z
updated: 2026-10-28T08:27:46Z
transitions: []
stream: S-0512
tags: [template]
touches: [design/adrs/0000-template.md, template/root/design/adrs/0000-template.md, design/system/documentation-standard.md, template/CHANGELOG.md, template/template.yaml]
---
# T-2191 The ADR template's status comment and the documentation standard name rejected as the fifth status, with a template changelog entry and bump

## Work

In both copies of the ADR template, `design/adrs/0000-template.md` and `template/root/design/adrs/0000-template.md`, change the status comment on line 4 from `# proposed | accepted | superseded | deprecated` to `# proposed | accepted | superseded | deprecated | rejected`, keeping the line's shape, since the vitest test of a later task reads the list from that comment.

In `design/system/documentation-standard.md`, under Editing rules for living documents, say that a rejected ADR, a proposed decision the operator declined in favour of another, is immutable like an accepted one: `flai doc` refuses any change to an accepted, superseded, deprecated, or rejected ADR but its `topics`, and `flai check` reports `adr.status` on a status outside the five. Link the story's ADR by its path once it exists, or by ID if this task closes first. Bump `updated`.

In `template/CHANGELOG.md`, add an entry at the top for the next patch version, named for S-0512, saying the template's ADR status comment gains `rejected` and that an older flai's `flai check` reports no `adr.status`; bump `version` in `template/template.yaml` by a patch, as S-0498's commit df7db120d did.

This task waits for nothing: it changes documents only.

## Done when

- Both `0000-template.md` copies name the five statuses on their status comment line, in the same shape.
- `documentation-standard.md` names a rejected ADR as immutable and names the `adr.status` finding.
- `template/CHANGELOG.md` has the entry and `template/template.yaml` the bumped version.
- `scripts/template-test.sh` renders and checks clean, and the markdown lint passes on the changed files.

## Notes

Drafted by planner-S-0512.
