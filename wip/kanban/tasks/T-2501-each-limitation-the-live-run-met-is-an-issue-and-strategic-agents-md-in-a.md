---
id: T-2501
type: task
nature: feature
title: "Each limitation the live run met is an issue, and strategic-agents.md § In a container (S-0509) says what the live run showed"
status: backlog
parent: S-0567
owner: alex
created: 2026-11-07T03:43:51Z
updated: 2026-11-07T03:44:02Z
transitions: []
stream: S-0567
tags: [docker, docs]
touches: [design/system/strategic-agents.md, design/issues/summary.md]
after: [T-2499, T-2500]
---
# T-2501 Each limitation the live run met is an issue, and strategic-agents.md § In a container (S-0509) says what the live run showed

## Work

Waits for T-2499 and T-2500: it files and writes up what they recorded in S-0567's Notes under `### Live run`.

For each limitation recorded there, file an issue with `flai issue new`, naming S-0567 and the story of E-0029 that would carry its fix. Where an open issue holds it already, such as I-0256 (no host configuration in the container), I-0257 (`.git/config` writable), or I-0258 (a project's own role refused), add the instance with `flai issue bump` instead. `flai issue` regenerates `design/issues/summary.md`.

Add to `design/system/strategic-agents.md` § In a container (S-0509) a paragraph on the live run (S-0567): the release and image it ran, the mounts seen, the guard refusal shown, the activity entry and how its cost was priced, and each limitation with a link to its issue. Bump the document's `updated`. Say only what the run showed; ADR-0247 stays as it is.

## Done when

- Every limitation in S-0567's `### Live run` notes names an issue under `design/issues`, new or bumped.
- `design/system/strategic-agents.md` § In a container (S-0509) says what the live run showed, with each issue linked, and passes the markdown lint.

## Notes
