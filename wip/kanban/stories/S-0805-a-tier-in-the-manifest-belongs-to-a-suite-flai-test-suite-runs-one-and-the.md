---
id: S-0805
type: story
nature: feature
title: A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:41:56Z
updated: 2026-11-16T23:41:56Z
transitions: []
tags: [flai, template]
topics: [verify, workflow]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §1. A tier in the manifest's `tests:` may give `suite: <name>`; a tier that gives none is the suite of its own name. This is the label every later story of the ADR reads. Accept the ADR, or record what the operator changes in it, as this story's first task.

## Acceptance criteria

- [ ] ADR-0312 is accepted, or superseded by the ADR this story writes from the operator's changes.
- [ ] A tier's `suite` is read from the manifest, kebab-case, defaulting to the tier's name, and `flai check` refuses one that is not kebab-case.
- [ ] `flai test --suite <name>` runs the tiers of a suite that the paths select, and with `--all` every tier of it but a `manual` one.
- [ ] The template's manifest labels its tiers `static`, `behavior`, `integration`, and `smoke`, and this repository's manifest labels its own without changing which tiers any command runs.
- [ ] The manifest design, the reference, and the user guide describe `suite`.

## Tasks

## Notes
