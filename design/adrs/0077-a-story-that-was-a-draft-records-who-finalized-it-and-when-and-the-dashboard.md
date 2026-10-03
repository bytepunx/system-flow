---
id: ADR-0077
title: "A story that was a draft records who finalized it and when, and the dashboard finalizes through item.finalize"
status: proposed
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0074]
---

# ADR-0077 A story that was a draft records who finalized it and when, and the dashboard finalizes through item.finalize

## Context

ADR-0074 gave stories a `draft` flag, set on a story an agent wrote and cleared when the operator finalizes it, with `flai edit --no-draft` or `flai move <story> ready --yes`. Nothing recorded who cleared it: `updated` is a time with no author, `flai edit` notes its author only in a git-ignored file beside the MCP cursors, and transitions record moves between states, which a finalize in the backlog is not. Once the flag was cleared, nothing in the item said it had been a draft. The metrics (S-0205) need to tell a planner's draft the operator took as written from one they rewrote, and the dashboard needed a write that finalizes without moving the story (S-0201). The question is TH-0083.

## Decision

A story that was a draft keeps a `finalized` block, `by` and `at`, written when its draft flag is cleared, and the dashboard finalizes through the hostapi write `item.finalize`.

- `flai edit --no-draft` stamps `finalized` with the editor (`--by`, else `FLAI_AGENT`, else the config author), as ADR-0074 stamps cost of delay and forecast. A finalizing move, `flai move <story> ready --yes` or `item.move` `finalize`, stamps it with the move's `by` and time, the same as its transition.
- Clearing the flag of a story that is not a draft stamps nothing and leaves an earlier block as it is. `flai edit --draft` removes the block: the story is a draft again.
- `Validate`, and through it `flai check` (`item.front-matter`), refuses a `finalized` block on an epic or task, one without `by` or a UTC `at`, and one on a story that is still a draft.
- `item.finalize` takes `{id}` and runs `flai edit <id> --no-draft --by=<owner>`, committed as `item.edit` is. It refuses an ID that is not a story as invalid, and a story that is not a draft as a rule (`Rule`), because there is nothing to finalize.
- The key is listed in `front-matter-fields.txt`, so the release that carries it raises `flai.minimum`, as ADR-0074's fields did.

## Consequences

- Whether a story was a draft, who took it as their own, and when, stays in the item after the flag is gone. Whether the body changed between the draft and `finalized.at` is in the git history, which the metrics can read.
- The board's card menu (S-0202) reuses `item.finalize`.
- A finalize that comes second, after someone else's, is refused rather than overwriting the first record.
- One more front matter key: a host whose flai is older stops before reading the project until it is upgraded.

## Alternatives considered

- **Finalize moves the story to ready, so its transition's `by` records who.** No new key, but a story could not be finalized and kept in the backlog, and nothing would say it had been a draft.
- **A transition entry for the finalize.** Transitions are moves between states; one that is not would skew time in state and every metric that reads them.
- **`updated` plus the edit note beside the cursors.** The note is git-ignored and local to a host, so the metrics could not read it.
- **A `drafted` block as well, kept from creation.** Who wrote the draft is already the first transition's `by` and the git history.
