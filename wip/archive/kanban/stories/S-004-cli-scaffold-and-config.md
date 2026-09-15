---
id: S-004
type: story
nature: feature
title: CLI scaffold and config
status: done
parent: E-002
owner: agent
created: 2026-09-15T16:10:00Z
updated: 2026-09-15T18:05:00Z
transitions:
  - to: ready
    at: 2026-09-15T17:05:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:06:00Z
    by: agent
  - to: review
    at: 2026-09-15T17:55:00Z
    by: agent
  - to: done
    at: 2026-09-15T18:05:00Z
    by: alex
tags: []
---

# S-004 CLI scaffold and config

## Goal
A flai binary with cobra command tree, version command, ~/.flai/config.json read and write, and CI for lint, test, and build.

## Acceptance criteria
- [x] flai version prints version, commit, date
- [x] flai config get and set round-trip the config file
- [x] First run creates the config with defaults
- [x] golangci-lint and go test -race pass in CI (workflow written; verified locally with golangci-lint v2.5.0 and go test -race, CI itself runs once the repo is on GitHub)

## Tasks
- T-008 Go module and cobra root command
- T-009 flai version
- T-010 Config read and write
- T-011 Lint config and CI workflow

## Notes
- cobra v1.10.2 pinned. errcheck excludes fmt.Fprint* in .golangci.yaml.
- Host golangci-lint is v1.64.8; v2.5.0 was installed to a scratch GOBIN for the local check. Install v2 to ~/go/bin to lint normally.
