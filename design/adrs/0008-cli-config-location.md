---
id: ADR-0008
title: CLI configuration lives in ~/.flai/config.json
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0008 CLI configuration lives in ~/.flai/config.json

## Context

The user must be able to change the template source to a fork or branch, and the dashboard image and port, without flags on every call.

## Decision

`flai` reads and writes `~/.flai/config.json`. The template cache lives in `~/.flai/cache`. `FLAI_CONFIG` and `--config` override the file path. The format is plain JSON parsed with `encoding/json`; no config framework.

## Consequences

- One fixed location on every OS, easy to document and to back up. Not XDG-compliant on Linux, which was accepted for simplicity and because the brief specified the path.
- Config is small and human-editable; `flai config get|set` exists for convenience only.

## Alternatives considered

- XDG base directories: correct on Linux, different on every OS, more to explain.
- YAML config: JSON was specified and is unambiguous.
