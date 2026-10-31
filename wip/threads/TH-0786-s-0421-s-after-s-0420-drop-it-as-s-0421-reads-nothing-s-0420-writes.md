---
id: TH-0786
title: "S-0421's after S-0420: drop it, as S-0421 reads nothing S-0420 writes"
anchor:
  path: wip/kanban/stories/S-0421-a-formula-is-a-markdown-file-under-design-workflows-formulas-that-flai-check-validates-and-flai-formula-run-makes-a-story-from-it-with-one-task-per-step-and-the-steps-after-between-them.md
  item: S-0421
status: open
participants: [agent-S-0551, orchestrator]
created: 2026-10-31T09:54:33Z
updated: 2026-10-31T09:56:18Z
---

# TH-0786 S-0421's after S-0420: drop it, as S-0421 reads nothing S-0420 writes

On wip/kanban/stories/S-0421-a-formula-is-a-markdown-file-under-design-workflows-formulas-that-flai-check-validates-and-flai-formula-run-makes-a-story-from-it-with-one-task-per-step-and-the-steps-after-between-them.md.

## Entries

### 2026-10-31T09:54:33Z agent-S-0551
**Recommendation:** drop S-0420 from S-0421's `after`:

```bash
flai edit S-0421 --after ""
```

## Why

- ADR-0230 (S-0551) says `after` names a story whose output this one reads, with the reason in Notes. The pull order sets an order of preference, and a trial merge now judges two stories that change the same files.
- S-0421 adds formula files under `design/workflows/formulas/`, `flai formula run`, and their check. S-0420 builds the dashboard's role page. S-0421 reads none of it: a role step runs through `role.run`, which S-0417 gave.
- S-0421's Notes name S-0420 only as a story replanned the same way, not as something it reads.
- While the `after` stays, S-0421 waits for S-0420's acceptance for no reason of its own.

The story's words are yours, so I have not edited them. Alternative: keep it, and add under `### Planning` what S-0421 reads from S-0420.

### 2026-10-31T09:56:18Z orchestrator (recommendation)
**Recommended:** drop S-0420 from S-0421's `after`, as agent-S-0551 proposes:

```bash
flai edit S-0421 --after ""
```

**Why:**

- ADR-0188 §6 says `after` names only a story whose accepted change the story needs on main, and never an order of preference.
- S-0421 needs nothing from S-0420. S-0421 adds formula files, `flai formula run`, and their check. S-0420 builds the dashboard's role page. A formula's role step runs through `role.run`, and S-0417 (done) added that.
- S-0421's Notes name S-0420 only as a story that was replanned the same way. That is a shared history, not a dependency.
- While the `after` stays, S-0421 waits for S-0420's acceptance with no reason of its own.

**Why you:** S-0421 is finalized, so its words are yours. My role cannot edit it.

Source: design/adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md § Decision
