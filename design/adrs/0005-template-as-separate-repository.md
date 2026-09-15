---
id: ADR-0005
title: The template is a separate git repository rendered by flai
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0005 The template is a separate git repository rendered by flai

## Context

The standard must be applied to new and existing repos consistently, be versionable, and be forkable by teams with their own conventions. The CLI must not need a release to pick up a template change.

## Decision

The template is a git repository with a `template.yaml` manifest and a `root/` tree. `flai` clones it into a cache keyed by repo and ref, renders `.tmpl` files with Go `text/template`, copies everything else, and records the template repo, ref, and version in the project manifest. The default source is `https://github.com/bytepunx/system-flow-template`; users override it in `~/.flai/config.json`, and a local path is accepted for development. The prototype lives in `./template` in this monorepo until the first publish.

## Consequences

- Template changes ship independently of `flai` releases. `min_flai` in the manifest guards against incompatibility.
- Two repos to maintain. A CI job in this monorepo renders and checks the prototype so it never drifts from the standard.
- Go `text/template` syntax in template files is visible to anyone reading the template repo directly; the `.tmpl` suffix keeps the rendered files clean.

## Alternatives considered

- Embedding the template in the `flai` binary: simplest, but every convention tweak needs a release and forks are impossible.
- Cookiecutter or copier: a Python dependency for a Go tool.
