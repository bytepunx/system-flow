---
id: S-031
type: story
nature: feature
title: flai structured logging and --verbose
status: backlog
parent: E-002
owner: alex
created: 2026-09-16T23:24:18Z
updated: 2026-09-16T23:24:18Z
transitions: []
tags: [cli, logging]
---

# S-031 flai structured logging and --verbose

## Goal
flai follows design/conventions/logging.md: structured events on stderr via log/slog, levels with the agreed meanings, `--verbose`, `LOG_LEVEL`, and `LOG_FORMAT`, with stdout reserved for command output.

## Acceptance criteria
- [ ] A `--verbose` global flag selects debug; `LOG_LEVEL` and `LOG_FORMAT` are honoured; text on a terminal, JSON otherwise
- [ ] Every existing stderr message (created config, fetching template, warnings) is a slog event with `component` and named fields, and stdout is unchanged so `--json` consumers are unaffected
- [ ] Errors are logged once at the command boundary with `err`; exit codes unchanged
- [ ] Tests pin field names for the events the CLI docs mention
- [ ] docs/users/flai.md documents the flags and variables

## Tasks

## Notes
- Created by S-030. Nature feature; bump flai minor on acceptance.
