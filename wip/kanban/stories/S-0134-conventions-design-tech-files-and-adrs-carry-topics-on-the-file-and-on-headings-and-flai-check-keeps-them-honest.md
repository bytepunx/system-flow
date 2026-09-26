---
id: S-0134
type: story
nature: feature
title: Conventions, design, tech files, and ADRs carry topics on the file and on headings, and flai check keeps them honest
status: backlog
parent: E-0010
owner: alex
created: 2026-09-26T17:46:02Z
updated: 2026-09-26T17:46:02Z
transitions: []
tags: [cli, template]
touches: [flai/internal/conventions, flai/internal/adr, flai/internal/docedit, flai/internal/check, flai/cmd, template/root/design, design/conventions, design/system]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0134 Conventions, design, tech files, and ADRs carry topics on the file and on headings, and flai check keeps them honest

## Goal

Documents can say which stories they are for, as [ADR-0047](../../../design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) decides: `topics: [...]` in the front matter of conventions, `design/system` and `design/tech` files, and ADRs, and `<!-- topics: a, b -->` on any heading line. Every convention, in the template and here, starts as `topics: [all]` for the designer to narrow. Nothing reads topics yet to select context; this story makes them parseable, settable, and checked.

## Acceptance criteria
- [ ] A package parses a document's file topics and splits its body into sections by heading, each with its heading path and effective topics (its own, else its parent's, else the file's; a convention without `topics` is `[all]`), with behavior tests for nesting, inheritance, and a heading comment at the end of the line.
- [ ] `flai adr topics ADR-nnnn <topics…>` sets `topics` on an ADR of any status and changes nothing else; `flai doc save` and the dashboard's editor still refuse any other change to an accepted ADR, and `topics` in the front matter of an accepted ADR is no longer a refusal on its own.
- [ ] `flai check` warns on a `design/system` or `design/tech` file without `topics`, and on a topic used by a document that no sub-project (name, tags, kind), `code`, `all`, story, or epic uses.
- [ ] Every file in `template/root/design/conventions/` and `design/conventions/` has `topics: [all]`; `flai upgrade` brings it to a project; the template's version and changelog record it.
- [ ] The baseline rule that an accepted ADR is edited only to set `superseded_by` also allows `topics`, in `decisions.md`, `CLAUDE.md` and its template, and `design/system/documentation-standard.md`; `design/system/conventions.md` describes the field and the heading comment.
- [ ] `docs/users/flai.md` documents `flai adr topics`; `make test`, `make integration`, `make smoke`, and `flai check --strict` pass.

## Tasks

## Notes

First of the E-0010 stories from S-0125; S-0135 to S-0138 build on it, in that order. Topic vocabulary and inheritance are in ADR-0047 "Topics on documents". The heading comment was chosen because it does not render and MD033 is off.
