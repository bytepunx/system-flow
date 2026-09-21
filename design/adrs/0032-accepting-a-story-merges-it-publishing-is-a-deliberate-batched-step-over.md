---
id: ADR-0032
title: "Accepting a story merges it; publishing is a deliberate, batched step over everything accumulated"
status: accepted
date: 2026-09-21
supersedes: []
superseded_by: []
refines: [ADR-0019, ADR-0025, ADR-0031]
---

# ADR-0032 Accepting a story merges it; publishing is a deliberate, batched step over everything accumulated

## Context

Until now, accepting a story or an epic computed its release, tagged, and (when the push host action was enabled) pushed, all in the same step as the merge (ADR-0019: "`flai accept` rebases, merges ..., then tags and pushes as before"). Each accepted item cut its own release. A run of several small stories against the same component therefore cut a run of small releases — three remediations in a row meant three patch tags and three pushes — when one release covering all three would have served the operator just as well and cost three times less attention. The operator asked for this directly on 2026-09-20: moving a story to done should merge it and nothing else; a release should be a deliberate step, covering everything accumulated since the last one.

`flai accept`, `release.Compute`/`Apply`/`Tag`, and `flai release <id>` (which already computed and tagged one item without pushing) already kept merging, computing, and tagging as separate function calls glued together by `acceptItem`. Splitting them further, and extending the computation from one item to everything pending, was mostly a matter of using seams that already existed.

## Decision

**Acceptance merges, moves to done, archives, and commits. Nothing else.** `flai accept` (and `flai move <story> done`, which runs the same flow) no longer computes a bump, creates a tag, or pushes. It refuses an experiment before merging, exactly as before (ADR-0025); a research story is accepted the same as any other nature now, since nothing at accept time distinguishes a nature's release from another's any more — there is no release at accept time to distinguish. This is the part of ADR-0025's decision that changes: it is no longer meaningful to say a research story is "accepted **and pushed** without a release," because acceptance itself never pushes, whatever the nature. What ADR-0025 decided about *acceptance* — research goes through, an experiment is refused — is unchanged; what it said about *pushing* now describes every nature alike, not something particular to research.

**Publishing is a deliberate step over everything accumulated, not any one item.** `flai release --pending` computes, per component, every accepted item since that component's last tag (or ever, with none), and takes the **highest delivery level among them** — an epic outranks a feature story, which outranks a remediation or improvement — as the one bump for all of them together. Three patch-level stories and one feature story against the same component since its last tag now produce one minor release, not four. It bumps and commits the version files, creates the tag(s), and pushes the branch and every tag together, three tags to a push (`pending.Batches`, I-0026, which becomes load-bearing here rather than an edge case: a batched publish routinely produces more tags at once than a single acceptance did). A component nothing has touched since its last tag is left out of the batch entirely.

**Resumable, not transactional.** A publish that fails partway — one component's tag made, the next not; the push itself failing — leaves local state a second run can pick up from exactly where the first stopped: a tag that already exists is not recreated (`TagPending`'s existence check), and the push step re-asks what is locally ahead of the remote (the same `pending.Detect` the existing `flai push --pending` uses) rather than only what the failed run itself tagged, so a tag made before a crash is still found and pushed on the next attempt.

**The board's done column shows the distinction, and a Publish action does the same thing as the shell command.** A done card says whether its story is published or still waiting; a banner at the top of the column lists what publishing now would release — which components, which bump, which story IDs — and a Publish button when the operator has enabled the push host action, the same one that gated the old automatic push (ADR-0031, S-0078: `publish.run` reuses `ActionPush` rather than adding a second gate for what is, underneath, the same "push with the operator's credentials" capability). Off, the banner says what enables it and changes nothing; the dashboard cannot enable it itself, as with push before this story.

**A component's version, after this story, describes what was last cut from a batch, not from one story.** Its tag message and changelog entry name every item the batch bundled, not one item's title; a reader working out what shipped in `cli/v1.1.0` finds a list of stories, as they would have before, just possibly more than one.

## Consequences

- A story no longer determines its own release in isolation; what a batch publishes depends on what else was merged before the operator clicked Publish, which the operator now controls directly rather than it following automatically from acceptance order.
- The board's "accepted, not pushed" notice (`flai push --pending`, ADR-0026) keeps its old meaning — ordinary merged commits, released or not, that have not reached the remote — and now sits beside a second, distinct notice for releases specifically pending publication. The two can both be true at once: something can be merged and unpushed while also being unreleased.
- `flai release <id>` (one item, on request) is unchanged and still useful on its own, independent of the batch: `--pending` is a sibling flag on the same command, not a replacement.
- git.md's baseline "Tag on acceptance" rule is no longer true and is rewritten (template first, then here, in this story) to describe merge-then-publish; its stale note that `flai release` "computes the bump, tags, and pushes once it exists" is corrected too — it existed already, computing and tagging one item, before this story added `--pending` and the push.
- Old per-project or per-story habits built around "every acceptance releases" (release-note tooling watching every tag, for instance) will see fewer, larger releases instead of one per story; nothing here changes what a release tag or a changelog entry *is*, only when and how many land at once.

## Alternatives considered

- **A bump per story, summed rather than taking the highest.** Three patches would still make one release, but as `patch+patch+patch`, which is not how semver bumps compose and would produce version jumps a reader could not explain from the changelog alone. The highest level among the batch is what semver actually means by "what changed since the last release."
- **Let the operator choose which merged-and-unpublished stories to include in a given publish**, rather than always publishing everything accumulated. Discussed with the operator and declined: it adds a selection UI and a partial-batch state for a case ("hold one merged story back") that the workflow already has a place for — do not move it to done yet. One button, no picker, was chosen for simplicity.
- **Gate publishing behind a new host action**, separate from push. Rejected: publishing *is* pushing, with a bump and a tag added first; a second gate for the same underlying "act with the operator's credentials" capability would only be two settings to keep in sync instead of one.
