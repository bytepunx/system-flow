---
id: ADR-0155
title: "A rebased story branch that undoes main's change to a path is recorded at sync and blocks acceptance until the operator restores the path, keeps the removal with a reason, or sends the story back"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0069, ADR-0093, ADR-0098]
---

# ADR-0155 A rebased story branch that undoes main's change to a path is recorded at sync and blocks acceptance until the operator restores the path, keeps the removal with a reason, or sends the story back

## Context

[I-0124](../issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md): S-0232 changed `scripts/flai-snapshot.sh` and then, as agreed on MS-0014, reverted it to main of that time. S-0340, accepted later, made the same change to the same lines. When `flai stream sync` rebased `story/S-0232` onto main, git dropped the first commit, which had become empty, and replayed the revert cleanly over main's identical change. The branch now undid S-0340's change. No conflict showed; the only sign was the sync's note that the branch changed a path outside its touches, and that note leaves out every path inside the claim. Accepting the branch would have taken main's change away.

Before the rebase the branch's net change to the path is empty, so nothing about the branch alone shows the loss. After it, a later rebase starts from a base that already holds main's change, so a check made then cannot find the path either. The close-out syncs before review, so the acceptance's own rebase is usually a no-op.

The first plan told the agent to restore the file with `git checkout` and a commit. On TH-0420 (2026-10-14T18:14:21Z) the operator chose instead that the operator resolves it from the dashboard, path by path, with flai doing the git work: pushing the work down into git commands or hand-edited files is an anti-pattern.

## Decision

A path where a rebased story branch removes lines main added since the branch's old base is recorded when the sync finds it, and blocks the story's acceptance until the operator restores main's version with `flai stream restore`, keeps the removal with `--keep-removal` and a reason, or sends the story back.

- **The check.** `storygit.Rebase` reads the branch's merge base with main before `git rebase`. After a clean rebase, for each path that main changed since that old base and that the rebased branch still changes against main, it compares the lines main added since the old base with the lines the branch's net change against main removes, as a multiset of line text. A path where the branch removes any of them undoes main's change: the finding gives the path, main's lines the branch removes, and the old base. A path the branch deletes removes every line main added. A rebase whose old base is main's head has nothing to check.
- **The record.** Each finding is recorded per story when the sync makes it, in `.flai-cache/undone/S-nnnn.json` in the project's main checkout, beside `flai verify`'s last result: the story, and per path its lines, old base, and when it was found. A reader gives the recorded paths still undone, checking each against main now, so a path the branch no longer removes main's lines from drops out by itself. The record is host-local, like the verify result.
- **At sync.** `flai stream sync` and `flai task done` warn on each such path, inside the story's claim or outside it, name `flai stream restore S-nnnn --path <path>`, and give the findings as `undoes_main` in `--json`. Neither fails, as the checks of ADR-0046 do not.
- **The blocker.** The acceptance preview, `flai accept --dry-run` and the host read `accept.preview`, lists each recorded path not kept as a blocker, and gives it as data, `undoes_main`, with main's lines the branch removes, so the dashboard can show each path as a card. `flai accept` refuses on it before anything is merged; with `--json` it exits 4 with the paths as data. Its own rebase runs the check too, so a path only the acceptance's rebase finds is refused, merging nothing. This refines the conditions of ADR-0093 and the acceptance's rebase of ADR-0098.
- **Restore.** `flai stream restore S-nnnn --path <path>` commits main's version of each path onto `story/S-nnnn` in its worktree, or removes the file where main has none, refusing while the worktree has uncommitted changes or a rebase is unfinished, and for a path the record does not hold. The path then leaves the record. The dashboard runs it through the host write `stream.restore`. It is flai's git work, so ADR-0069's rule that agents never rebase or merge by hand stands.
- **Keep.** `flai accept S-nnnn --keep-removal <path> --reason "<why>"` passes that path, and records the reason under `### Removals kept` in the story's Notes, in the acceptance commit. `accept.run` takes `keep_removals`, each a path and a reason. flai and `accept.run` refuse a keep with `--by orchestrator`: an undone path waits for the operator.
- **Send back.** Moving the story back to in-progress with a reason, as today; its agent restores the path or redoes the change.
- **The dashboard.** The review page's cards and the board's badge are a story of their own, after this one, reading the data named here.

## Consequences

- A revert replayed over main's identical change is caught at the sync that makes it, and can no longer reach main unseen.
- The operator, not an agent and not the orchestrator, decides whether a story may remove what main gained, and the reason is kept with the story.
- A story that means to remove main's lines, such as one rewriting a section another story just added, needs the operator's keep at acceptance.
- The record lives on one host: an acceptance on another clone sees only what its own rebase finds.
- The check is by line text, so a line main added that the branch removes for another reason, such as a blank line or a brace, is reported too. The operator keeps such a removal with a reason.

## Alternatives considered

- **Tell the agent to restore with `git checkout <main> -- <path>` and a commit** (the first plan): it pushes the resolution into git commands, which the operator rejected on TH-0420.
- **Refuse the acceptance with no override** (the first plan): an intended removal could not be accepted short of restoring the lines and editing them again.
- **Let the orchestrator keep a removal**: it cannot tell an intended removal from a lost change; the operator keeps the decision.
- **Stop the sync on the finding**: a story may mean the removal, so the sync warns and acceptance is the gate.
- **Extend the outside-touches note**: it leaves out every path inside the claim, so it would miss the same revert of a claimed path.
- **Check only at acceptance**: by then the close-out's sync has moved the base past main's change, and the path cannot be found.
