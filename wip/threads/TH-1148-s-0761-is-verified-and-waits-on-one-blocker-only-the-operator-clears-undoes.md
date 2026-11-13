---
id: TH-1148
title: "S-0761 is verified and waits on one blocker only the operator clears: undoes_main on design/system/flai-cli.md"
anchor:
  path: wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md
  item: S-0761
status: resolved
participants: [orchestrator, alex]
created: 2026-11-13T09:34:11Z
updated: 2026-11-13T11:42:32Z
---

# TH-1148 S-0761 is verified and waits on one blocker only the operator clears: undoes_main on design/system/flai-cli.md

On wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md.

## Entries

### 2026-11-13T09:34:11Z orchestrator
**Recommended:** accept S-0761 by hand, keeping the removal:

```bash
flai accept S-0761 --keep-removal design/system/flai-cli.md --reason "the flai check row keeps every word of main's; S-0761 only adds to it"
```

**Where it stands.**

- The close-out passed at `09d428b`, the head of `story/S-0761`, with every tier, integration and smoke included.
- The verifier matched both criteria to the changed files. Criterion 1: `flai/internal/check/scope.go`, `flai/internal/check/scope_test.go`, `flai/cmd/check.go`, `flai/cmd/check_test.go`. Criterion 2: `design/issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md`. S-0678's blocked interval is closed on main.
- The acceptance preview names one blocker, `undoes_main`. It says `story/S-0761` removes 1 line main added to `design/system/flai-cli.md`: the `flai check` row of § Commands.

**Why it is a false positive (I-0360).** I split main's row and the branch's row at each sentence and compared them. The branch keeps every sentence of main's and adds two things: the `item.blocked-closed` sentence, and `item.blocked-closed` with ADR-0296 in the `--record-issues` list. The check compares whole lines, and that row is one line. `flai stream restore` would drop S-0761's own edit.

**Why it is yours.** ADR-0155 makes keeping a removal the operator's decision. The orchestrator's `accept` takes no `--keep-removal`. My recommendation on TH-1147 said this before the story came back to review. flai serve has since resolved that thread.

I leave S-0761 in review.

### 2026-11-13T11:42:32Z alex
Resolved: S-0761 was accepted
