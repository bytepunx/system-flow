---
id: TH-0676
title: Which story does Claude Code work inside the container for T-2166?
anchor:
  path: wip/kanban/stories/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md
  item: S-0505
status: answered
participants: [agent-S-0505, orchestrator, alex]
created: 2026-10-28T12:05:18Z
updated: 2026-10-28T12:11:11Z
---

# TH-0676 Which story does Claude Code work inside the container for T-2166?

On wip/kanban/stories/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md.

## Entries

### 2026-10-28T12:05:18Z agent-S-0505
S-0505's criterion 3 needs one story of this repository, chosen with you, worked inside the container from its first task to review. The story's notes ask for a documentation-only story, since the images carry no Go toolchain. No ready story is documentation-only: the backlog's are all Go or flaiover work.

**Recommendation:** let me create one small documentation-only story under E-0029 for it, not a draft, and move it to ready myself: "docs/operators/runbooks gains a runbook for building the agent images and running a story's agent in a container by hand" (touches a new `docs/operators/runbooks/agent-container.md` and `docs/operators/runbooks/README.md`, so it overlaps no open story). It is real work the epic needs and small enough for one session. The run uses your Claude login read-only from `~/.claude` and spends one session of your subscription, about 30 to 60 minutes.

Alternatives: name another story yourself, or let me make it a throwaway story that I cancel after the run.

Answer "as recommended", name a story, or say "throwaway". I am building the images (T-2149) and the container script (T-2159) meanwhile; T-2166 waits for this answer.

### 2026-10-28T12:05:33Z orchestrator (recommendation)
**Recommended for the operator:** answer "as recommended". Let agent-S-0505 create the documentation-only runbook story under E-0029, finalized, and move it to ready.

## Why I recommend it

- S-0505's notes ask for a documentation-only story unless you name another, because the bedrock and claude-code images carry no Go toolchain. No ready or backlog story is documentation-only today.
- A runbook for building the agent images and running an agent in a container by hand is work E-0029 needs anyway. A throwaway story would spend the same subscription session and leave nothing behind.
- Its two touches, `docs/operators/runbooks/agent-container.md` (new) and `docs/operators/runbooks/README.md`, overlap no open story, so it can start as soon as T-2166 needs it.

## Why it is yours

- Criterion 3 says the story is "chosen with the operator".
- The run spends one session of your Claude subscription, about 30 to 60 minutes, and mounts your login read-only.
- A new story is a change of scope. `work-management.md` § Rules keeps an agent from inventing scope.

**The alternatives:** name an existing story, or answer "throwaway" and the agent cancels the story after the run.

S-0505 goes on building the images (T-2149) and the container script (T-2159) meanwhile. Only T-2166 waits on your answer.

Source: design/conventions/work-management.md § Rules

### 2026-10-28T12:11:11Z alex
Confirmed the recommendation of 2026-10-28T12:05:33Z orchestrator.

Source: design/conventions/work-management.md § Rules
