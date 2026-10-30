---
id: ADR-0223
title: "An ADR's status is one of proposed, accepted, superseded, deprecated, or rejected, a rejected ADR is immutable as an accepted one is, and flai check reports adr.status on any other"
status: accepted
date: 2026-10-30
supersedes: []
superseded_by: []
refines: [ADR-0001]
---

# ADR-0223 An ADR's status is one of proposed, accepted, superseded, deprecated, or rejected, a rejected ADR is immutable as an accepted one is, and flai check reports adr.status on any other

## Context

An ADR's status was named in one place only, the comment on the `status` line of `design/adrs/0000-template.md`: `proposed | accepted | superseded | deprecated`. Nothing checked it. ADR-0202 was set `rejected` by hand on main, in favour of [ADR-0203](0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md), a status the template's comment, the documentation standard, `flai doc`'s refusal of a change to an accepted ADR, and the dashboard's `search.test.ts` did not name. `flai check` accepted it without a finding, and every story's vitest tier failed until S-0501 added `rejected` to the test's own list (I-0208).

The cause is a status landing on main that nothing validates, with the list kept in four copies that drift apart.

## Decision

An ADR's status is one of `proposed`, `accepted`, `superseded`, `deprecated`, or `rejected`; a rejected ADR is a proposed decision the operator declined, recorded in favour of another, and is immutable once set as an accepted one is; and `flai check` reports `adr.status` on an ADR whose status is outside the five.

- **One list.** flai keeps the five statuses in one exported list that `flai check`, `flai doc`, and `flai adr` read. `flai adr new --status` still takes only `proposed` or `accepted`.
- **Immutable.** `flai doc` refuses any change to an accepted, superseded, deprecated, or rejected ADR but its `topics`, as it did for the first three ([ADR-0047](0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md)). A rejection is written on a proposed ADR, which is still a draft and editable; from then on the ADR is read-only.
- **The finding.** `adr.status` is a warning on the `status` line of an ADR whose status is outside the five, naming the status found and the five allowed. The template `0000-template.md`, whose status is `template`, is passed over.
- **The template.** The status comment in both copies of `0000-template.md` names the five, and the dashboard's test reads its allowed statuses from that comment rather than from a list of its own.

## Alternatives

- **Keep four statuses and record a rejection as a superseding ADR.** That loses the record that nothing took the rejected decision's shape, and ADR-0202 already carries `rejected` on main.
- **Validate nothing and let each test keep its own list.** The cause stays: the next status written by hand fails a later story's tiers rather than the close-out of the change that wrote it.

## Consequences

- An unknown status stops a close-out at `flai check`, before the vitest tier, on the story that wrote it.
- An older flai reports no `adr.status`; a project on it is checked by the template's comment alone.
- The documentation standard, `flai-cli.md`, and the user docs name the five statuses and the finding.
