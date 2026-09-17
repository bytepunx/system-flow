# Changelog

## 1.0.2 - 2026-09-17

- S-0033 Four-digit work item IDs (patch).

## 1.0.1 - 2026-09-17

- S-0021 flai template push publishes template changes to a remote (patch).

## 1.0.0 - 2026-09-16

Epic E-0005 (agent conventions) complete: twelve baseline conventions under
`design/conventions`, `design/issues` with `flai issue`, `scripts/` behind
the Makefile with three test tiers, and a `CLAUDE.md` that primes every
session from the conventions. First major release of the template.

## 0.4.0 - 2026-09-16

- `design/conventions/logging.md` and `telemetry.md` are active: five log levels including fatal, structured events, `/_health`, `/_ready`, `/metrics`, OpenTelemetry, golden signals (S-0030).

## 0.3.0 - 2026-09-16

Note: under the refined release rule (incidental additive touches are a patch) this would have been 0.2.1; the version stands as cut.

- `design/conventions/logging.md` and `telemetry.md` added as drafts, indexed at 110 and 120 (S-0025).
- Markdownlint excludes test fixtures and `bin/`.

## 0.2.0 - 2026-09-16

- `design/conventions/`: ten baseline agent norms with a project-additions marker (S-0023, S-0026).
- `design/issues/` skeleton and `scripts/` with test tier stubs and Makefile delegation (S-0026).
- `CLAUDE.md` opens with a priming section and is a map; norms live in the conventions (S-0024).
- Rendered projects are markdownlint-clean; lint config aligned with the item format.

## 0.1.0 - 2026-09-15

Initial prototype. Layout, work item schema, baseline CLAUDE.md, devex files.
