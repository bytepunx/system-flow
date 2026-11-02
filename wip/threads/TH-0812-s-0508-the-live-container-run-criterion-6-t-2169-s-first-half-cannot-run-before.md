---
id: TH-0812
title: "S-0508: the live container run (criterion 6, T-2169's first half) cannot run before S-0508 is accepted and published; move it to a follow-up story?"
anchor:
  path: wip/kanban/stories/S-0508-flai-serve-runs-a-story-s-agent-in-a-container-through-a-container-harness-adapter-with-the-worktree-the-story-s-cache-folder-and-the-login-mounted-the-keys-and-the-mcp-token-by-environment-and-the-image-and-network-the-host-sets.md
  item: S-0508
status: resolved
participants: [agent-S-0508, orchestrator, alex]
created: 2026-11-02T06:52:58Z
updated: 2026-11-02T18:09:37Z
---

# TH-0812 S-0508: the live container run (criterion 6, T-2169's first half) cannot run before S-0508 is accepted and published; move it to a follow-up story?

On wip/kanban/stories/S-0508-flai-serve-runs-a-story-s-agent-in-a-container-through-a-container-harness-adapter-with-the-worktree-the-story-s-cache-folder-and-the-login-mounted-the-keys-and-the-mcp-token-by-environment-and-the-image-and-network-the-host-sets.md.

## Entries

### 2026-11-02T06:52:58Z agent-S-0508
**Recommended:** move criterion 6 and the live run in T-2169 to a draft follow-up story in E-0029, with `after: S-0508`.

**Why it cannot run here.** Criterion 6 asks for one story worked end to end by `flai serve` on `harness: container`. The `flai serve` on this host is the installed release (`~/.flai/bin/flai`). The `container` adapter, its settings, and the serve start path exist only on `story/S-0508`. They reach that `flai serve` only once S-0508 is accepted and published, and the host upgraded. Starting a second `flai serve` from this tree is not an option on this host.

**What I would do on "yes":**

- Remove criterion 6 from S-0508, and narrow T-2169 to the documentation.
- Write a draft story in E-0029 with criterion 6's words and the live-run half of T-2169's Work. Its Notes say it waits on S-0508 for the adapter on the installed flai.
- Record the move in S-0508's Notes.

**Precedent:** you answered the same question this way for S-0423 (TH-0803, D-0461), S-0536 (D-0422), and S-0537 (D-0423).

**Alternative:** keep criterion 6. S-0508 then cannot meet the definition of done and stays in progress until after its own release.

Everything else in S-0508 goes on meanwhile: the adapter, the settings, the protected list, the serve start path, and the documents. Only the live run waits on this answer.

### 2026-11-02T06:54:14Z orchestrator (recommendation)
**Recommended:** say yes, as you answered for S-0423 on TH-0803 (D-0461).

- S-0508's agent removes criterion 6 from S-0508 and narrows T-2169 to the documentation.
- It writes a draft story in E-0029 with `after: S-0508`. The draft holds criterion 6's words and the live-run half of T-2169's Work. Its Notes say why it waits: the adapter must reach the installed flai first.
- S-0508's Notes record the move.

**Why:** this is the same case as D-0461 and D-0422. The `container` adapter, its host settings, and the serve start path exist only on `story/S-0508`. The `flai serve` on this host is the installed release, so no story can run on `harness: container` until S-0508 is accepted, published, and the host upgraded.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and adds a story. That changes scope, so it is your call. Finalizing the draft is yours too.

**Alternative:** keep criterion 6. I don't recommend it: S-0508 could never meet the definition of done.

Source: design/digest/D-0461-s-0423-the-live-publish-run-t-1783-criterion-4-cannot-finish-before-s-0423-is.md

### 2026-11-02T18:08:27Z alex
Confirmed the recommendation of 2026-11-02T06:54:14Z orchestrator.

Source: design/digest/D-0461-s-0423-the-live-publish-run-t-1783-criterion-4-cannot-finish-before-s-0423-is.md

### 2026-11-02T18:09:37Z agent-S-0508
Resolved: Criterion 6 moved to S-0565, a draft in E-0029 with after: S-0508; S-0508's Notes record the move.
