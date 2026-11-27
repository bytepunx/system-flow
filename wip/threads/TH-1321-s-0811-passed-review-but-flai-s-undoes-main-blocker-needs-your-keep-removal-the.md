---
id: TH-1321
title: "S-0811 passed review, but flai's undoes_main blocker needs your --keep-removal: the two lines it drops are rewrites, not losses"
anchor:
  path: wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md
  item: S-0811
status: resolved
participants: [orchestrator, agent-S-0811, alex]
created: 2026-11-21T18:34:28Z
updated: 2026-11-27T22:47:30Z
---

# TH-1321 S-0811 passed review, but flai's undoes_main blocker needs your --keep-removal: the two lines it drops are rewrites, not losses

On wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md.

## Entries

### 2026-11-21T18:34:28Z orchestrator
**Recommended:** accept S-0811 yourself, keeping the two removals:

```sh
flai accept S-0811 --by alex --keep-removal design/system/workflow.md --keep-removal docs/operators/runbooks/change-the-board.md --reason "the branch rewrites main's 'nothing else reads the keys' sentence to add S-0811's publish; every mention in main's line is kept"
```

**Why I cannot accept it:** `accept`'s dry run at 1c58c4be57b5e8bb6559fbdca4f66ff21b24907d, the head of `story/S-0811`, names one blocker, `undoes_main`, in two files. My `accept` tool takes no `--keep-removal`, and the other way out, `flai stream restore`, would bring back main's stale sentence beside the story's new one.

**What the two removals are:** each is the "Nothing else reads the keys" sentence, which the story rewrites. I-0408, which the story records, is about exactly this sentence.

- `design/system/workflow.md`: the new line keeps S-0807, S-0809, S-0810, and S-0808 and S-0812 as still to come. It moves S-0811 to done.
- `docs/operators/runbooks/change-the-board.md`: the new lines keep "the acceptance's gates and CI run as they did". They add the publish (S-0811) and correct main's stale claim that the batch runs as it did, which has not been true since S-0809.

**What passed at 1c58c4b:**

- `flai verify S-0811 --last`: every step passed, integration and smoke included.
- The verifier found all four criteria met, and every changed file is within the story's touches.
  - 1: `flai/cmd/release.go`, `flai/cmd/release_publish_suites.go`, `flai/cmd/release_publish_suites_test.go`, `flai/internal/hostapi/writes_outcome.go`, `flai/internal/hostapi/writes_actions.go`, `flai/internal/hostapi/writes_test.go`
  - 2: `template/root/design/workflows/formulas/publish.md`, `flai/cmd/publish_formula_suites_test.go`, `template/CHANGELOG.md`
  - 3: `flai/cmd/release_publish_suites_test.go` (`TestReleasePendingRunsNoTierWithoutPublish`)
  - 4: `design/system/pushing-from-the-board.md`, `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`

**Alternative:** send S-0811 back to in-progress so its agent keeps main's line word for word. That costs another close-out and changes nothing a reader needs.

S-0811 stays in review. I changed nothing.

### 2026-11-21T18:38:30Z agent-S-0811
Agreed: both removals are deliberate rewrites, so keep them. After TH-1322 the branch is rebased again, onto 574abdc7 (S-0812's acceptance), and it gives the same two removals. This time the line removed from each file is S-0812's version of the sentence. It lists the publish (S-0811) as still to come, so it cannot be kept word for word.

- `design/system/workflow.md`: the new sentence names the batch (S-0809), a state's runs (S-0810), the publish (S-0811), and `flai test --at ci` (S-0812) as done, and the acceptance's gates (S-0808) as still to come. S-0812's `flai test --at` paragraph and S-0811's publish paragraph follow it, both whole.
- `docs/operators/runbooks/change-the-board.md`: S-0812's sentence keeps every word but its last clause, which now adds the publish since S-0811 and leaves only the acceptance's gates running as they did.

S-0811 goes back to review once the close-out passes at the new head. The `flai accept S-0811 --keep-removal design/system/workflow.md --keep-removal docs/operators/runbooks/change-the-board.md --reason "..."` you recommended still applies. I recorded the tool's missing `--keep-removal` as I-0410.

### 2026-11-21T18:39:22Z orchestrator
Noted. My recommendation to the operator stands, for the new head once it is in review.

- S-0811 is in progress now. Its branch is being rebased onto 574abdc7 after TH-1322.
- Once it is back in review and its close-out passes at the new head, I verify it again and run `accept`'s dry run there.
- If `undoes_main` names only these two rewrites, the operator's command above still applies. My `accept` tool still takes no `--keep-removal` (I-0410).

**What the operator needs to do:** nothing yet. Accept S-0811 with the command above once I report that it passed at the new head.

### 2026-11-27T22:47:30Z alex
Resolved: S-0811 was accepted
