---
id: ADR-0148
title: "A design token's colour is in the srgb colour space with its hex or in oklch, the fixed colour roles are required with extra tokens allowed, an invalid token file is the check finding design.tokens, and contrast is reported, never a finding"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
topics: [design, cli, dashboard]
---

# ADR-0148 A design token's colour is in the srgb colour space with its hex or in oklch, the fixed colour roles are required with extra tokens allowed, an invalid token file is the check finding design.tokens, and contrast is reported, never a finding

## Context

[ADR-0146](0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md) makes a project's tokens a DTCG file. The Design page's sample components are compiled with the dashboard, so they can style only names they know in advance. flai computes contrast, which is defined on sRGB. Every finding a save introduces refuses the save ([ADR-0147](0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md)). The finding, [design-system-tool.md](../system/design-system-tool.md), recommended choices 2, 3, and 4 below; on TH-0445, on 2026-10-14, the operator answered only choice 1, so these stand as recommended.

## Decision

A design token's colour is in the `srgb` colour space, with its `hex`, or in `oklch`; the fixed colour roles are required and extra tokens are allowed; an invalid token file is the check finding `design.tokens`; and contrast is reported, never a finding.

1. **Colour spaces.** `srgb`, with its `hex`, and `oklch`. Alpha only in shadows and extra tokens. Any other colour space the DTCG Color module defines is a finding. An `oklch` colour outside sRGB is clipped for contrast, and the report says so.
2. **Roles.** Required: `background`, `surface`, `text`, `text-muted`, `border`, `primary`, `on-primary`, `secondary`, `on-secondary`, `focus`, and the states `success`, `warning`, `danger`, `info`, each with a `-soft` tint. Any other name is an extra token: validated, rendered, and shown as a swatch.
3. **`design.tokens`.** A `flai check` rule validates the token file against flai's subset of the format: it parses, its types and colour spaces are the ones flai reads, the required roles and steps are present, and the resolver is the `theme` modifier with `light` and `dark`. Its findings refuse a `file.save`, as every finding a save introduces does.
4. **Contrast.** flai computes the roles' pairs at the thresholds flaiover's own test uses: text, muted text, and each state on `background`, `surface`, or its tint at 4.5:1; `on-primary` and `on-secondary` on their fill at 4.5:1; `border` and `focus` on `background` at 3:1. `flai design` and the Design page report them. No check rule reads contrast.

## Consequences

- The operator can save a trial colour whose contrast fails, and sees the failure beside it.
- Converting `oklch` to sRGB is fixed arithmetic; it adds no dependency.
- A token file built in another tool in `display-p3` or `hsl` must be converted before flai accepts it.
- The gallery can rely on every role being present.

## Alternatives considered

- **Hex only.** Simpler, but `oklch` is what Tailwind v4's palette uses, and it steps evenly in lightness, which makes a dark value easier to derive.
- **Any colour.** Variables, mixes, and wide-gamut colours have no sRGB contrast flai can compute.
- **Free names.** The gallery could not style them.
- **Contrast as a check finding.** It would refuse every save of a trial colour that fails.
