---
id: T-1534
type: task
nature: feature
title: flai check reports a design system token that does not parse as design.tokens with its key and line
status: backlog
parent: S-0381
owner: alex
created: 2026-10-13T17:46:35Z
updated: 2026-10-14T23:53:27Z
transitions: []
stream: S-0381
tags: [cli]
touches: [flai/internal/check/design.go, flai/internal/check/design_test.go, flai/internal/check/check.go]
after: [T-1590, T-1592]
---
# T-1534 flai check reports a design system token that does not parse as design.tokens with its key and line

## Work

Add the rule `design.tokens` in `flai/internal/check/design.go`, called from `Run` in `check.go`: one error per problem `designsystem.Validate` returns, on the document's path and the token's line. A project with no design system document has no finding. Test both, and a valid document. It waits for T-1532, whose `Validate` it reports.

## Done when

- `go test ./internal/check/` passes with the three cases.
- `flai check --strict` on this repository is unchanged.

## Notes

- Layer 2, beside T-1533.
