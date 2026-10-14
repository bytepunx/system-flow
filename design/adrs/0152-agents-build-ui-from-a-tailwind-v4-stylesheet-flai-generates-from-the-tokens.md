---
id: ADR-0152
title: "Agents build UI from a Tailwind v4 stylesheet flai generates from the tokens and design.stale keeps in step, and prime gives the stories the guidance's topics reach the guidance and the token file whole, with a baseline design-system convention"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0049, ADR-0050]
topics: [design, cli, template]
---

# ADR-0152 Agents build UI from a Tailwind v4 stylesheet flai generates from the tokens and design.stale keeps in step, and prime gives the stories the guidance's topics reach the guidance and the token file whole, with a baseline design-system convention

## Context

An agent building UI needs the token names and the guidance on when to use each. It should get them without being told, and its code should not drift from the tokens. [ADR-0049](0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) loads whole what a story names and briefs what its topics select; [ADR-0050](0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md) briefs a large file named only by its path. Prime reads only the markdown files directly in `design/system`, `design/tech`, and `design/adrs`. [ADR-0146](0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md) puts the tokens in a DTCG file and the guidance in `design/system/design-system.md`. The finding, [design-system-tool.md](../system/design-system-tool.md), recommended choice 10 below; on TH-0445, on 2026-10-14, the operator answered only choice 1, so it stands as recommended.

## Decision

Agents build UI from a Tailwind v4 stylesheet flai generates from the tokens and `design.stale` keeps in step, and prime gives the stories the guidance's topics reach the guidance and the token file whole, with a baseline design-system convention.

1. **The stylesheet.** flai generates a Tailwind v4 `@theme` stylesheet from the tokens, at a path the manifest names, with `flai design css --write`. The UI's own stylesheet imports it whole, and nothing hand-written sits in it. How it selects the dark values is S-0385's to settle.
2. **`design.stale`.** A `flai check` finding when the stylesheet differs from the tokens. It names the token file as well as the stylesheet, so a branch that changed only the token file stops at its close-out.
3. **Prime.** For every story the guidance document's `topics` reach, prime loads the guidance and the token file whole, as if the story named them. Prime reads the token file at the path the manifest names, beside the markdown it reads today. Above ADR-0050's threshold it briefs them, and the convention tells the agent to read them with `doc_get` and `flai file show`.
4. **One list.** The guidance's `topics` alone say which stories a design system reaches. There is no second record of governed sub-projects.
5. **The convention.** A baseline convention, `design-system.md`, in the template and here, with `roles: [story, verify]`. Prime gives it only to the stories the guidance reaches, and to none where a project has no design system. It says: build UI from the token utilities; no arbitrary value or hex for colour, type, spacing, radius, or shadow; when the UI needs a value the tokens lack, change the token file on the story's branch, widen the touches to it and the stylesheet, run `flai design css --write`, and say why in the guidance; ask on a thread first when the change alters what a role means.
6. **No lint yet.** A lint of raw values in components is deferred until `design/issues` shows raw values getting through, and then belongs in the project's own linter rather than in flai.

## Consequences

- Utilities exist for the tokens and no others of the project's, so the tokens are the easy path. Arbitrary values such as `bg-[#3b82f6]` still compile; the convention and the verifier catch them.
- Every UI story's pack grows by the guidance and the token file, an estimated eighth of the budget or less.
- Prime gains its first read of a file that is not markdown.
- A baseline convention ships with a rule that keeps it out of projects without a design system.

## Alternatives considered

- **Briefs and the convention only.** A brief gives the guidance's outline, not the tokens, and leaves reading the rest to the agent's memory, the failure the operator reports.
- **A record of the sub-projects a design system governs, written by `flai design init --project`.** A second list beside `topics` that can disagree with it.
- **The lint now.** False positives on layout values, SVG, and chart palettes, and flai reads no code today.
