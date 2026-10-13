---
id: T-1543
type: task
nature: feature
title: An edit model holds changed tokens apart from the saved ones and writes them back into the document's front matter, leaving the rest as it was
status: backlog
parent: S-0383
owner: alex
created: 2026-10-13T17:47:27Z
updated: 2026-10-13T17:47:27Z
transitions: []
stream: S-0383
tags: [dashboard]
touches: [flaiover/src/lib/design/tokens.ts, flaiover/src/lib/design/tokens.test.ts]
---
# T-1543 An edit model holds changed tokens apart from the saved ones and writes them back into the document's front matter, leaving the rest as it was

## Work

`flaiover/src/lib/design/tokens.ts` keeps the saved tokens and the edited ones, says which changed, checks a value as flai's `design.tokens` does so a bad one is marked before saving, discards back to the saved tokens, and composes the document's new raw text: the `tokens` block replaced in the front matter with the `yaml` package, every other key, its comments where it can, and the body left as they were. Tests cover each, with a round trip that changes nothing. It waits for nothing.

## Done when

- `tokens.test.ts` passes, including the unchanged round trip.

## Notes

- Layer 1, beside T-1544.
