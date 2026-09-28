---
id: S-0134
type: story
nature: feature
title: Conventions, design, tech files, and ADRs carry topics on the file and on headings, and flai check keeps them honest
status: done
parent: E-0010
owner: alex
created: 2026-09-26T17:46:02Z
updated: 2026-09-28T22:39:28Z
transitions:
  - to: ready
    at: 2026-09-26T22:40:26Z
    by: alex
  - to: in-progress
    at: 2026-09-27T03:56:51Z
    by: agent-S-0134
  - to: review
    at: 2026-09-28T22:05:54Z
    by: agent-S-0134
  - to: done
    at: 2026-09-28T22:39:28Z
    by: alex
tags: [cli, template]
touches: [flai/internal/conventions, flai/internal/adr, flai/internal/docedit, flai/internal/check, flai/internal/topics, flai/internal/upgrade, flai/internal/lock, flai/internal/metrics/testdata, flai/cmd, template/root/design, template/root/CLAUDE.md.tmpl, template/template.yaml, template/CHANGELOG.md, design/conventions, design/system, design/tech, design/adrs/README.md, design/issues, CLAUDE.md, docs/users, docs/contributors, docs/operators/settings.md]
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
- [x] A package parses a document's file topics and splits its body into sections by heading, each with its heading path and effective topics (its own, else its parent's, else the file's; a convention without `topics` is `[all]`), with behavior tests for nesting, inheritance, and a heading comment at the end of the line.
- [x] `flai adr topics ADR-nnnn <topics…>` sets `topics` on an ADR of any status and changes nothing else; `flai doc save` and the dashboard's editor still refuse any other change to an accepted ADR, and `topics` in the front matter of an accepted ADR is no longer a refusal on its own.
- [x] `flai check` warns on a `design/system` or `design/tech` file without `topics`, and on a topic used by a document that no sub-project (name, tags, kind), `code`, `all`, story, or epic uses.
- [x] Every file in `template/root/design/conventions/` and `design/conventions/` has `topics: [all]`; `flai upgrade` brings it to a project; the template's version and changelog record it.
- [x] The baseline rule that an accepted ADR is edited only to set `superseded_by` also allows `topics`, in `decisions.md`, `CLAUDE.md` and its template, and `design/system/documentation-standard.md`; `design/system/conventions.md` describes the field and the heading comment.
- [x] `flai upgrade` keeps the `topics` a project set on a convention and takes the template's on one whose topics the project left as the template gave them; `system-flow.lock.yaml` records the template's topics per marker file (TH-0028, answer C).
- [x] `docs/users/flai.md` documents `flai adr topics`; `make test`, `make integration`, `make smoke`, and `flai check --strict` pass.

## Tasks
- T-0489 A topics package parses file topics and splits a document into sections with effective topics
- T-0490 flai adr topics sets topics on an ADR of any status, and doc save lets topics alone change on an accepted ADR
- T-0491 flai check warns on design and tech files without topics and on topics nothing uses
- T-0492 Every convention in the template and here carries topics: [all], and the template's version records it
- T-0493 Rules and docs say an accepted ADR may gain topics, describe the field and the heading comment, and document flai adr topics
- T-0494 flai upgrade keeps a project's own topics on a convention and takes the template's otherwise

## Notes

Verified on story/S-0134: `make flai-test` (lint, behaviour, integration, smoke) and `flai check --strict` pass. The template's version and changelog are written by the release tooling at publish, not by hand (`git.md`): `flai release S-0134 --dry-run` shows `template/template.yaml` and `template/CHANGELOG.md` bumped. Stories and epics have no `topics` key until S-0135, so the check's vocabulary is `all`, `code`, and the manifest's words for now; S-0135 adds the items' topics. `flai doc save` takes a topics-only change to an accepted ADR; the dashboard's editor still shows accepted ADRs read-only.

First of the E-0010 stories from S-0125; S-0135 to S-0138 build on it, in that order. Topic vocabulary and inheritance are in ADR-0047 "Topics on documents". The heading comment was chosen because it does not render and MD033 is off.
