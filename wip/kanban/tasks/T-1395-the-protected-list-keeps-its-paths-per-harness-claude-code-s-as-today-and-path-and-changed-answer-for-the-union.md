---
id: T-1395
type: task
nature: improvement
title: The protected list keeps its paths per harness, Claude Code's as today, and Path and Changed answer for the union
status: backlog
parent: S-0354
owner: alex
created: 2026-10-08T08:48:08Z
updated: 2026-10-08T08:48:08Z
transitions: []
stream: S-0354
tags: [cli]
touches: [flai/internal/protected/protected.go, flai/internal/protected/protected_test.go]
---
# T-1395 The protected list keeps its paths per harness, Claude Code's as today, and Path and Changed answer for the union

## Work

- `protected.go`: the folders and files become a map from harness name to its list, `claude-code` holding today's entries, plus `For(harness string) []string` and `Harnesses()`. `Path`, `Git`, and `Changed` keep their signatures and answer for the union, so `guard`, `permission.go`, and `preview/accept.go` need no change.
- The doc comment says the list is flai's policy (ADR-0130), that a harness's own files join it with its adapter, and why `.git` stays apart.
- `protected_test.go` pins Claude Code's list and checks the union.

First layer: it runs together with T-1394, which touches other files.

## Done when

- `flai test flai/internal/protected/ flai/internal/preview/` passes with no expectation changed.

## Notes

Layer 1 of S-0354.
