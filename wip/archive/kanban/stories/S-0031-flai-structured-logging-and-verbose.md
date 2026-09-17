---
id: S-0031
type: story
nature: feature
title: flai structured logging and --verbose
status: done
parent: E-0002
owner: alex
created: 2026-09-16T23:24:18Z
updated: 2026-09-16T23:51:18Z
transitions:
  - to: ready
    at: 2026-09-16T23:45:51Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:45:51Z
    by: agent
  - to: review
    at: 2026-09-16T23:48:50Z
    by: agent
  - to: done
    at: 2026-09-16T23:51:18Z
    by: alex
tags: [cli, logging]
---

# S-0031 flai structured logging and --verbose

## Goal
flai follows design/conventions/logging.md: structured events on stderr via log/slog, levels with the agreed meanings, `--verbose`, `LOG_LEVEL`, and `LOG_FORMAT`, with stdout reserved for command output.

## Acceptance criteria
- [x] A `--verbose` global flag selects debug; `LOG_LEVEL` and `LOG_FORMAT` are honoured; text on a terminal, JSON otherwise
- [x] Every existing stderr message (created config, fetching template, warnings) is a slog event with `component` and named fields, and stdout is unchanged so `--json` consumers are unaffected
- [x] Errors are logged once at the command boundary with `err`; exit codes unchanged
- [x] Tests pin field names for the events the CLI docs mention
- [x] docs/users/flai.md documents the flags and variables

## Tasks
- T-0057 logx package: slog setup, levels including fatal, env and flag
- T-0058 Replace stderr messages and the error boundary with slog events
- T-0059 Tests pinning field names, docs

## Notes
- Created by S-0030. Nature feature; bump flai minor on acceptance.
