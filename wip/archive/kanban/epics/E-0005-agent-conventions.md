---
id: E-0005
type: epic
nature: feature
title: Agent conventions
status: done
owner: alex
created: 2026-09-15T18:26:31Z
updated: 2026-09-17T00:04:08Z
transitions:
  - to: ready
    at: 2026-09-17T00:04:08Z
    by: alex
  - to: in-progress
    at: 2026-09-17T00:04:08Z
    by: alex
  - to: review
    at: 2026-09-17T00:04:08Z
    by: alex
  - to: done
    at: 2026-09-17T00:04:08Z
    by: alex
tags: []
---

# E-0005 Agent conventions

## Outcome
Every system-flow project ships a `design/conventions` folder: one markdown file per topic area that tells any agent how to work here, so the operator never starts a project by re-establishing standards, norms, and ways of working. The baseline set comes from the template, projects extend it below a marker, and the template's CLAUDE.md instructs the agent to prime each session with the conventions before touching anything. This repository adopts the same conventions. This epic takes priority over all other open work.

## Stories
- S-0022 Define the conventions folder standard
- S-0023 Write the baseline convention set
- S-0024 Prime agent sessions from CLAUDE.md
- S-0025 Tooling support for conventions
- S-0026 Adopt the conventions in this repository
- S-0027 flai issue commands and check rules for design/issues
- S-0030 Logging and telemetry conventions
- S-0026 Adopt the conventions in this repository
- S-0027 flai issue commands and check rules for design/issues
- S-0030 Logging and telemetry conventions
- S-0027 flai issue commands and check rules for design/issues
- S-0030 Logging and telemetry conventions
- S-0030 Logging and telemetry conventions

## Notes
Requested 2026-09-15. Supersedes the previous pull order: S-0022 to S-0025 come first. S-0010 stays in review awaiting the human's tag push.
