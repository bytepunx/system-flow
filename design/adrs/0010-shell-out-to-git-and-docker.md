---
id: ADR-0010
title: flai shells out to git and docker
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0010 flai shells out to git and docker

## Context

`flai` needs to clone template repos (with credentials, branches, and private hosts) and run the dashboard container. Both have mature CLIs that already handle the user's authentication and configuration.

## Decision

`flai` invokes `git` and `docker` as subprocesses and requires them on `PATH`. It does not link go-git or the Docker SDK. Missing binaries produce a clear error with an install hint.

## Consequences

- Private template repos work with whatever credential helper the user already has.
- `flai` binary stays small and dependency-light.
- Tests mock the process runner; a small integration suite runs real `git` when available.

## Alternatives considered

- go-git: no credential helper support, partial clone support, larger binary.
- Docker Engine SDK: pulls in a large dependency tree for two commands.
