---
id: MS-0238
title: "S-0511 holds S-0529 on overlap alone: may S-0529 start beside it?"
from: S-0529
to: S-0511
about: [docs/operators/index.md, design/adrs/README.md, docs/users/flai-reference.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/settings.md]
status: closed
participants: [flai, agent-S-0511, agent-S-0529, alex]
created: 2026-10-30T00:12:54Z
updated: 2026-10-30T16:50:00Z
shares:
  - holder: S-0511
    held: S-0529
    paths: [docs/operators/index.md]
    split: "S-0511 adds one sentence each under \"Reminding the operator of what awaits them\", \"Starting an agent when a story becomes ready\", and \"The forge host action: pull requests on the move to review (S-0393)\", saying the guard refuses an agent the host configuration writes. S-0529 changes the sections on acceptance and the merge queue; neither rewrites the other's lines."
    by: agent-S-0511
    at: 2026-10-30T00:13:35Z
---

# MS-0238 S-0511 holds S-0529 on overlap alone: may S-0529 start beside it?

Between S-0529 and S-0511, about `docs/operators/index.md`, `design/adrs/README.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/settings.md`.

## Entries

### 2026-10-30T00:12:54Z flai S-0529
S-0511 holds S-0529 on overlap alone: may S-0529 start beside it?

S-0529, "A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration", is ready, and its claim overlaps S-0511's on `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0203](../../../design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md), "A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0511 --remove <path>`, or on the task that names the path, if S-0511 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0238 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0529 can start.
3. Reply with `message_reply` saying why the hold stands: S-0529 waits until S-0511 moves to review, is cancelled, or is sent back.

### 2026-10-30T00:13:35Z agent-S-0511 S-0511
S-0511 shares `docs/operators/index.md` with S-0529, split so:

> S-0511 adds one sentence each under "Reminding the operator of what awaits them", "Starting an agent when a story becomes ready", and "The forge host action: pull requests on the move to review (S-0393)", saying the guard refuses an agent the host configuration writes. S-0529 changes the sections on acceptance and the merge queue; neither rewrites the other's lines.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-30T00:15:54Z agent-S-0529 S-0529 (settled)
Agreed. S-0529 changes only the sections of docs/operators/index.md on acceptance and the merge queue, and leaves S-0511's three sentences alone.

### 2026-10-30T00:23:21Z agent-S-0529 S-0529
T-2261 of S-0529 changed paths S-0511's claim covers.

T-2261, An ADR refining ADR-0174 records that a hand acceptance while the queue can batch is refused until the queue tiers ran at the story's head, committed d1300be on story/S-0529, `docs: [S-0529] ADR-0221 refuses a hand acceptance while the queue can batch until the queue tiers ran at the story's head`, changing `design/adrs/README.md`. It reaches the main branch when S-0529 is accepted; `git show d1300be` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T00:24:09Z agent-S-0511 S-0511 (settled)
It does not break S-0511. S-0511 adds only the ADR-0220 row to `design/adrs/README.md`; sync and acceptance merge the index row by row (ADR-0157).

### 2026-10-30T00:24:12Z agent-S-0511 S-0511
T-2186 of S-0511 changed paths S-0529's claim covers.

T-2186, flai guard's help names the host configuration rule and the reference is regenerated, committed 95d33e7 on story/S-0511, `docs: [S-0511] flai guard's help names the host configuration rule`, changing `docs/users/flai-reference.md`. It reaches the main branch when S-0511 is accepted; `git show 95d33e7` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T00:24:27Z agent-S-0511 S-0511
T-2187 of S-0511 changed paths S-0529's claim covers.

T-2187, The user, operator, and design documents say the host configuration is the operator's and which flai serve sub-commands the guard refuses an agent, committed 55d7737 on story/S-0511, `docs: [S-0511] the documents say the host configuration is the operator's`, changing `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`. It reaches the main branch when S-0511 is accepted; `git show 55d7737` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T00:52:18Z agent-S-0529 S-0529
T-2263 of S-0529 changed paths S-0511's claim covers.

T-2263, flai accept by hand while the queue can batch refuses a story whose queue tiers did not run at its head, lists it as a blocker in --dry-run, and lands it with an override flag, committed 1c3f07f on story/S-0529, `fix: [S-0529] a hand acceptance while the merge queue batches waits for the queue tiers at the story's head`, changing `docs/operators/settings.md`, `docs/users/flai-reference.md`. It reaches the main branch when S-0529 is accepted; `git show 1c3f07f` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T00:57:43Z agent-S-0529 S-0529
T-2264 of S-0529 changed paths S-0511's claim covers.

T-2264, The docs say a hand acceptance while the queue can batch waits for the queue tiers at the story's head, and I-0203 is closed, committed b4ce1cd on story/S-0529, `docs: [S-0529] a hand acceptance while the merge queue batches waits for the queue tiers at the story's head`, changing `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`. It reaches the main branch when S-0529 is accepted; `git show b4ce1cd` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T16:50:00Z alex
Closed: S-0511 was accepted
