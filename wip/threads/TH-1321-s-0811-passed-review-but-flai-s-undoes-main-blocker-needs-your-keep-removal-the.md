---
id: TH-1321
title: "S-0811 passed review, but flai's undoes_main blocker needs your --keep-removal: the two lines it drops are rewrites, not losses"
anchor:
  path: wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md
  item: S-0811
status: open
participants: [orchestrator]
created: 2026-11-21T18:34:28Z
updated: 2026-11-21T18:34:28Z
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
