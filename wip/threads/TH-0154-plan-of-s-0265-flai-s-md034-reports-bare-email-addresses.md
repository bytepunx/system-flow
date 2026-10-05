---
id: TH-0154
title: "Plan of S-0265: flai's MD034 reports bare email addresses"
anchor:
  path: wip/kanban/stories/S-0265-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md
  item: S-0265
status: resolved
participants: [planner-S-0265, alex]
created: 2026-10-05T05:50:21Z
updated: 2026-10-05T05:52:18Z
---

# TH-0154 Plan of S-0265: flai's MD034 reports bare email addresses

On wip/kanban/stories/S-0265-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md.

## Entries

### 2026-10-05T05:50:21Z planner-S-0265
Plan for S-0265, written by planner-S-0265.

Remedy proposed, from I-0056's instance: flai's MD034 (`flai/internal/mdlint/inline.go`, `bareURL`) knows only `http(s)://` literals. Recognise GFM's extended email autolink in the same inline parse and report it as MD034. markdownlint-cli2's output on a new fixture decides the edge cases.

Tasks, in two layers:

- Layer 1, running together because they share no path and wait for nothing:
  - T-0981 flai's MD034 reports a bare email address as markdownlint does: an email fixture, regenerated `expected.txt`, the change to `inline.go`, and the regression test `TestBareEmailOfI0056`.
  - T-0982 The user guide says flai's lint checks bare email addresses: one phrase in `docs/users/flai.md`.
- Layer 2:
  - T-0983 Close I-0056 saying what fixed it: after T-0981 and T-0982.

Enrichment: touches are `inline.go`, `mdlint_test.go`, `testdata/cases`, `docs/users/flai.md`, and the I-0056 file with `design/issues/summary.md`. I added the topic `cli`. The forecast is 25m, up from flai's 5m, going by the 16m to 23m that S-0262, S-0240, and S-0258 took, with delivery 2026-10-06T00:17Z. Cost of delay is 25 USD a week, the value flai works out from its own inputs, which I left as they were.

Assumptions:

- markdownlint-cli2 at the version `scripts/lint-md.sh` pins is the reference, through `scripts/mdlint-fixtures.sh`. That needs `npx` on the host that works the story.
- Only the bare `local@domain` form needs fixing. `mailto:` and `xmpp:` literals are fixed too only if the fixture shows markdownlint reporting them.
- Running the repository's wip through the stricter rule turns up no bare addresses: I-0056 says the ones in TH-0067 were fixed on main. If T-0981's `flai check --strict` finds any, they are findings outside the story, and the close-out records them as such.

T-0982 could be folded into T-0981. I kept it separate so the docs change runs alongside the fix. Tell me if you would rather merge them.

### 2026-10-05T05:52:18Z alex
Resolved.
