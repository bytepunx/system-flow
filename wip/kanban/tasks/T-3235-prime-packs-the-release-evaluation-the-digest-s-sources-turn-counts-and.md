---
id: T-3235
type: task
nature: feature
title: Prime packs, the release evaluation, the digest's sources, turn counts, and statsread read a state's kind
status: backlog
parent: S-0749
owner: alex
created: 2026-11-12T16:08:29Z
updated: 2026-11-12T16:08:29Z
transitions: []
stream: S-0749
tags: [cli]
touches: [flai/internal/context/reasons.go, flai/internal/context/reasons_figures.go, flai/internal/context/story.go, flai/internal/release/evaluate.go, flai/internal/digest/sources.go, flai/internal/usage/turns.go, flai/internal/statsread/statsread.go]
---
# T-3235 Prime packs, the release evaluation, the digest's sources, turn counts, and statsread read a state's kind

## Work

Replace the one to three references to a state's name in each of these files with a read of the kind: the orchestrator's reasons and a story's pack pick their items by kind, the release evaluation counts the stories in accepted states, the digest harvests from accepted and cancelled items, and the turn counts and statsread tell open from closed by kind.

## Done when

These files name no state, and the existing tests of these packages pass unchanged.

## Notes
