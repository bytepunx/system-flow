---
id: ADR-0078
title: "An accepted epic contributes no release bump of its own: its stories carry theirs"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0076]
---

# ADR-0078 An accepted epic contributes no release bump of its own: its stories carry theirs

## Context

`git.md` said an epic contributes a major bump, and `release.LevelFor` returned major for any epic. Accepting an epic was a deliberate act until S-0200, and no epic had been accepted here (flai was at 1.29.1). With ADR-0076 an epic is accepted with its last open story, so the end of every epic would have cut a major release of each component it delivered to, without anyone asking for one. The designer chose, on TH-0082, that it should not.

## Decision

An accepted epic contributes no bump of its own, whether it follows its last story into done or is accepted by hand with `flai accept E-nnnn`. Its stories carry their bumps when they are accepted: a feature story a minor, a remediation or improvement story a patch. A major release is a breaking change, and a breaking change has its own story against the component (git.md). `release.LevelFor` gives an epic `none`, as it gives research and experiments, and `flai release <epic>` says why it releases nothing.

## Consequences

- No version is bumped to a new major by an epic's end; a major release comes only from a story that makes a breaking change.
- Accepting an epic is bookkeeping: done and archived, no publish follows from it alone.
- git.md's baseline release rule changes in the template, a template patch release.

## Alternatives considered

- Keep the major bump and have the acceptance say "the next publish is a major release": the end of an epic is no longer a deliberate act, so it would cut major releases by accident.
- A minor bump for an epic: its feature stories already contribute minors, so it would add nothing.
