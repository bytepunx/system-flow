---
id: ADR-0017
title: Work item IDs are zero-padded to four digits
status: accepted
date: 2026-09-17
supersedes: []
superseded_by: []
refines: [ADR-0003]
---

# ADR-0017 Work item IDs are zero-padded to four digits

## Context

ADR-0003 fixed IDs as `E-nnn`, `S-nnn`, `T-nnn`. This repository passed T-0100 within two days of use, so a three-digit width stops sorting lexically and reads unevenly within a project's first months. ADR-0003 is refined, not reversed: the type letter, the dash, per-type sequences, and never reusing a number all stand.

## Decision

`flai` allocates IDs zero-padded to four digits (`E-0001`, `S-0034`, `T-0101`). The next number is one more than the highest existing number regardless of width. Three- and four-digit IDs are both valid and may coexist; sorting is numeric; every command accepts an ID in any padding and resolves it to the item's own ID. `flai migrate ids` widens an existing repository: it renames items and narratives with `git mv` and rewrites references in the layout folders, the root markdown and yaml files, and each project's root markdown files, with a dry run. Accepted ADRs that mention an item keep their meaning through the mechanical rename; that rewrite is the one edit permitted on them.

## Consequences

- IDs are stable for ten thousand items per type; a fifth digit is allowed by the parser and would trigger another refinement.
- Fixtures under `testdata/` keep three-digit IDs so the side-by-side path stays tested.
- Issue IDs (`I-nnn`) are unchanged; they are not kanban items and stay small.

## Alternatives considered

- Unpadded IDs: file listings and board order would not sort without a numeric-aware sort everywhere.
- Five or six digits: reads as noise for the projects this standard targets.
- Project-prefixed IDs: still deferred, as in ADR-0003.
