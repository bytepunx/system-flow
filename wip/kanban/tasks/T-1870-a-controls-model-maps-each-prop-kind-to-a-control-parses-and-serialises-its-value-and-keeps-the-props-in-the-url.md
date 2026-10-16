---
id: T-1870
type: task
nature: feature
title: A controls model maps each prop kind to a control, parses and serialises its value, and keeps the props in the URL
status: backlog
parent: S-0440
owner: alex
created: 2026-10-16T21:35:26Z
updated: 2026-10-16T21:35:26Z
transitions: []
stream: S-0440
tags: [dashboard]
touches: [flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/explorer/controls.test.ts]
---
# T-1870 A controls model maps each prop kind to a control, parses and serialises its value, and keeps the props in the URL

## Work

`controls.ts`: `controlFor(prop)` by kind, `parse(prop, text)` with its error, `initial(prop)` from the default as written where it is a literal, and `toQuery` and `fromQuery` for a component's props in the URL (JSON, base64url-encoded, under one parameter).

Waits for nothing.

## Done when

Tests cover every kind, a default that is not a literal, invalid JSON, and the URL's round trip. `flai test flaiover/src/lib/explorer` passes.

## Notes

None.
