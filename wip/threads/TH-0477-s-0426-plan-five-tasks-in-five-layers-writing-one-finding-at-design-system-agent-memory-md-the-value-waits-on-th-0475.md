---
id: TH-0477
title: "S-0426 plan: five tasks in five layers, writing one finding at design/system/agent-memory.md; the value waits on TH-0475"
anchor:
  path: wip/kanban/stories/S-0426-agent-memory-system.md
  item: S-0426
status: open
participants: [planner-S-0426]
created: 2026-10-16T05:59:09Z
updated: 2026-10-16T05:59:09Z
---

# TH-0477 S-0426 plan: five tasks in five layers, writing one finding at design/system/agent-memory.md; the value waits on TH-0475

On wip/kanban/stories/S-0426-agent-memory-system.md.

## Entries

### 2026-10-16T05:59:09Z planner-S-0426
S-0426 is planned with five research tasks in the backlog. They run one at a time because all five write `design/system/agent-memory.md`.

| Layer | Task | Writes | After |
|-------|------|--------|-------|
| 1 | T-1795 How flai recalls and searches today, its cost, and the headings candidates are judged under | `agent-memory.md`, `design/system/README.md` | none |
| 2 | T-1796 Open source candidates: memory systems, embedded stores, search engines | `agent-memory.md` | T-1795 |
| 3 | T-1797 State of the art: hybrid retrieval, RAPTOR and GraphRAG, MemGPT and A-MEM, caching and incremental indexes | `agent-memory.md` | T-1796 |
| 4 | T-1798 Trial of up to three approaches against flai's BM25 on this repository's corpus | `agent-memory.md` | T-1797 |
| 5 | T-1799 Recommendation, follow-on stories, operator decisions, link from `agent-context.md` | `agent-memory.md`, `design/system/agent-context.md` | T-1798 |

## Story figures

- Touches: `design/system/agent-memory.md`, `design/system/README.md`, `design/system/agent-context.md`. No folder touches.
- Topics: `cli` added.
- Forecast: 1h30m, raised from flai's 8m. Delivery 2026-10-16T10:00Z.
- Cost of delay: no value yet. It waits on TH-0475.

## Hold

S-0426 is held on `design/system/README.md` until S-0413 reaches review. flai asked S-0413's agent in MS-0127 to share the file, with one row for each story.

## Assumptions

- The finding is one design document, like S-0125's `agent-context.md` and S-0413's `agent-prose.md`. It writes no ADR and drafts no story. T-1799 names the follow-on stories and puts the decisions to you on a thread.
- The trial (T-1798) runs only in a scratch directory with local models: no paid service, no account, and no install on the host. If it needs one of those, the agent asks first. If you would rather skip the trial, drop T-1798. The recommendation would then rest on documentation alone.
- E-0025 (the digest) and S-0427 (the SQLite store) are related. The finding names what it means for each. It does not plan them.

## Proposals

None to split, merge, or drop.
