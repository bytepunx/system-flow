---
id: E-005
type: epic
nature: feature
title: Agent conventions
status: backlog
owner: alex
created: 2026-09-15T18:26:31Z
updated: 2026-09-15T18:26:32Z
transitions: []
tags: []
---

# E-005 Agent conventions

## Outcome
Every system-flow project ships a `design/conventions` folder: one markdown file per topic area that tells any agent how to work here, so the operator never starts a project by re-establishing standards, norms, and ways of working. The baseline set comes from the template, projects extend it below a marker, and the template's CLAUDE.md instructs the agent to prime each session with the conventions before touching anything. This repository adopts the same conventions. This epic takes priority over all other open work.

## Stories
- S-022 Define the conventions folder standard
- S-023 Write the baseline convention set
- S-024 Prime agent sessions from CLAUDE.md
- S-025 Tooling support for conventions

## Notes
Requested 2026-09-15. Supersedes the previous pull order: S-022 to S-025 come first. S-010 stays in review awaiting the human's tag push.
