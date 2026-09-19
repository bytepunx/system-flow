---
id: T-0233
type: task
nature: feature
title: "flai adr new and flai adr accept: number from the files, slug, front matter, index row, supersedes and refines, body on standard input, checked before kept, committed on its own"
status: done
parent: S-0060
owner: alex
created: 2026-09-19T09:39:14Z
updated: 2026-09-19T09:45:51Z
transitions:
  - to: ready
    at: 2026-09-19T09:39:14Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:39:15Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:45:51Z
    by: system-flow
stream: S-0060
tags: []
---

# T-0233 flai adr new and flai adr accept: number from the files, slug, front matter, index row, supersedes and refines, body on standard input, checked before kept, committed on its own

## Work
A package `flai/internal/adr` and the commands. `flai adr new "<title>" [--status proposed|accepted] [--supersedes N]... [--refines N]... [--body-stdin] [--autocommit] [--trailer]`: the next number is the highest `NNNN-*.md` present plus one (gaps are not filled, a counter in a document is not consulted); the file is `NNNN-slug.md` with `id`, `title`, `status`, `date`, `supersedes`, `superseded_by`, and `refines` when given; the body is the project's `0000-template.md` sections unless given on standard input; the index row is added to `design/adrs/README.md` with the status note the index already uses ("accepted, refines 0018", "supersedes 0007"); each superseded ADR gets `superseded_by` set, the one edit the conventions allow to an accepted ADR, and its index row is noted. With a body or `--autocommit` it is one step as for items (S-0059): `flai check` before and after, a refusal with exit 4 and the findings, every file restored, then one `docs:` commit of the new file, the index, and the superseded files. `--print-body` prints the template's sections. `flai adr accept N` sets a proposed ADR to accepted with today's date, updates its index row, and commits the same way. Tests with real git: numbering with gaps, slug, index row, supersedes and refines, refusal leaving nothing behind, commit contents, accept.

## Done when
- The commands are tested end to end with real git
- `make flai-test` passes

## Notes
