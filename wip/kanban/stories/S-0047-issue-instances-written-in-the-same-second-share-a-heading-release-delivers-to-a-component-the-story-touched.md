---
id: S-0047
type: story
nature: remediation
title: "Issue instances written in the same second share a heading; release delivers to a component the story touched"
status: backlog
parent: E-0001
owner: alex
created: 2026-09-18T17:11:45Z
updated: 2026-09-18T17:11:45Z
transitions: []
tags: [cli]
touches: [flai/internal/issues, flai/internal/release]
---

# S-0047 Issue instances written in the same second share a heading; release delivers to a component the story touched

## Goal
Two small flai defects that each broke something on main are fixed at the root: `flai issue bump` and `flai issue new` never write duplicate same-second headings, and the release never awards the delivery bump to a component the story's commits did not touch.

## Acceptance criteria
- [ ] Issue instances recorded within the same second share one heading, the way narrative log entries do since I-0011; a test bumps twice at one instant and the file passes the duplicate-heading rule
- [ ] The release picks the delivered component from the story's tags only among components its commits touched; when the tags name several, the one with the most touched files wins; a tag naming an untouched component never delivers; `--deliver` still overrides; tests cover a [dashboard, cli] story that only touched the CLI
- [ ] The remaining E-0006 stories carry tags that say where they deliver
- [ ] I-0011 and I-0016 are closed against this story

## Tasks

## Notes
- I-0016: S-0039 was pure flai work tagged [dashboard, cli]; flaiover took a minor with zero files and flai a patch. S-0035 and S-0037 had the same shape. Versions stand as cut.
- I-0011 second occurrence: duplicate instance headings in I-0016 failed the markdown lint on main on 2026-09-18.
