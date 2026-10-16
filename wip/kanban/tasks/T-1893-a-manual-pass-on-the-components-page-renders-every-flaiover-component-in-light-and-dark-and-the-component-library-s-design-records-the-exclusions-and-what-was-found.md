---
id: T-1893
type: task
nature: feature
title: A manual pass on the Components page renders every flaiover component in light and dark, and the component library's design records the exclusions and what was found
status: backlog
parent: S-0444
owner: alex
created: 2026-10-16T21:40:21Z
updated: 2026-10-16T21:40:21Z
transitions: []
stream: S-0444
tags: [dashboard]
touches: [design/system/component-library.md]
after: [T-1891, T-1892]
---
# T-1893 A manual pass on the Components page renders every flaiover component in light and dark, and the component library's design records the exclusions and what was found

## Work

Open every listed component on the Components page, in light and dark, from its defaults and each example. Record in the story's notes any that render wrongly, and what is wrong. Write a section of `design/system/component-library.md` on this repository's use: flaiover's setup, each exclusion with its reason from T-1890's decisions, and what the pass found.

Waits for T-1891 and T-1892, so that the pass sees every example and the test is in place.

## Done when

The notes record the pass, the design names each exclusion, and `flai check --strict` passes.

## Notes

None.
