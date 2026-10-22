---
id: T-2038
type: task
nature: improvement
title: flai parses, validates, and lists an example's answers, and the components.examples check reports a malformed one
status: backlog
parent: S-0476
owner: alex
created: 2026-10-22T20:59:16Z
updated: 2026-10-22T20:59:16Z
transitions: []
stream: S-0476
tags: [flai]
touches: [flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/components_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/internal/hostapi/components_test.go]
after: [T-2035]
---
# T-2038 flai parses, validates, and lists an example's answers, and the components.examples check reports a malformed one

## Work

Make flai accept the `answers` key the ADR of T-2035 fixes, in `flai/internal/components/examples.go`:

- Add `Answers []Answer` to `Example`, with `json:"answers,omitempty"`, and an `Answer` type with `Path`, `Method`, `Status`, `Headers`, `Body` (`json.RawMessage`), and `Events` (`[]AnswerEvent` with `Event` and `Data json.RawMessage`).
- Teach `oneExample` the key, so `ParseExamples` no longer reports it as one an example does not have, and report each malformed answer as a `Problem` with its line: a missing or empty `path`, a `method` that is not an HTTP method in upper case, a `status` outside 100–599 or not an integer, `headers` that is not an object of strings, an `events` entry without `data`, and any other key in an answer or an event.
- Leave `Validate` to the props: answers are not checked against the component.

`ListComponents` in `flai/internal/hostapi/components.go` and `flai components list --json` carry the answers through the struct's tags; update the command's help text in `flai/cmd/components.go` where it documents an example's keys, and `make flai-reference` if that text is generated into `docs/users/flai-reference.md`. The `components.examples` check reaches `ParseExamples` through `exampleFile` in `flai/internal/check/check.go`, which needs no change.

Tests: `examples_test.go` reads a file with answers and reports each malformed shape with its line; `flai/internal/check/components_test.go` has one case of a malformed answer reported as `components.examples`; `flai/internal/hostapi/components_test.go` and `flai/cmd/components_test.go` see `answers` in the listing.

This task waits for T-2035, which fixes the format it parses.

## Done when

- [ ] `ParseExamples` accepts `answers` with `path`, `method`, `status`, `headers`, `body`, and `events`, and reports each malformed answer with its line
- [ ] `flai components list --json` carries each example's `answers`
- [ ] `flai check` reports a malformed answer as `components.examples`, and passes a file whose answers are well formed
- [ ] The tests named above pass, with `go vet` and `golangci-lint` clean on the package

## Notes

Planned by planner-S-0476. The key's sub-keys follow T-2035's ADR; if that ADR settles on different ones, follow it.
