---
id: ADR-0001
title: Markdown is the only documentation format and decisions are recorded as ADRs
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0001 Markdown is the only documentation format and decisions are recorded as ADRs

## Context

system-flow needs documentation that agents can read and write with plain file tools, that diffs well in git, that renders in any host and in the dashboard, and that can carry structured metadata. Decisions need to be found later with their reasoning intact.

## Decision

All documentation in a conforming repo is markdown with YAML front matter. Architecture decisions are recorded as ADRs in `design/adrs`, numbered, immutable after acceptance, and superseded by later ADRs rather than edited. Diagrams are Mermaid fenced blocks.

## Consequences

- Every tool in the system parses one format. The dashboard indexes front matter without a database.
- Binary diagrams and office documents are not allowed under `design/` or `wip/`.
- The living design in `design/system` must be updated whenever an ADR changes it, which is a discipline the baseline `CLAUDE.md` enforces on agents.

## Alternatives considered

- AsciiDoc: richer, but weaker agent familiarity and renderer support.
- Notion or a wiki: not in the repo, not versioned with the code, invisible to agents without integrations.
