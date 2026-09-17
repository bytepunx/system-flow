---
id: T-0057
type: task
nature: feature
title: "logx package: slog setup, levels including fatal, env and flag"
status: done
parent: S-0031
owner: alex
created: 2026-09-16T23:45:51Z
updated: 2026-09-16T23:48:49Z
transitions:
  - to: ready
    at: 2026-09-16T23:48:49Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:48:49Z
    by: agent
  - to: done
    at: 2026-09-16T23:48:49Z
    by: agent
stream: S-0031
tags: [cli, logging]
---

# T-0057 logx package: slog setup, levels including fatal, env and flag

## Work
internal/logx: New(w, opts) builds a slog.Logger; level from --verbose, then LOG_LEVEL (debug, info, warn, error, fatal), default info; format from LOG_FORMAT (text, json), default text when w is a terminal else json; custom Fatal level above error rendered as FATAL in both handlers; ts, level, component, msg keys as the convention names them.

## Done when
Unit tests cover level and format resolution and the fatal rendering.

## Notes
