---
id: D-0384
title: "Plan for E-0029: six draft stories, their order, their tasks, and the assumptions: Recommended for the operator: approve the plan as drafted, and take the planner's three proposals"
kind: answer
sources: [TH-0659]
topics: [all]
paths: []
date: 2026-10-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0384 Plan for E-0029: six draft stories, their order, their tasks, and the assumptions: Recommended for the operator: approve the plan as drafted, and take the planner's three proposals

The operator, alex, answered TH-0659 at 2026-10-26T12:23:01Z.

What was asked:

```text
Six draft stories deliver the epic, in this order (`after`):
```

The answer:

```text
**Recommended for the operator:** approve the plan as drafted, and take the planner's three proposals.

- **The plan fits.** There are six stories, ordered by `after` only where one needs another's change on main. S-0506 and S-0507 can run beside each other after S-0505, and S-0509 and S-0510 after S-0508. Each story's tasks form layers that share no path. Each ADR question, `permission_prompt` over HTTP and the strategic agents' mount, is left to its story's ADR with a stated recommendation.
- **Proposal 1, a follow-up epic for flai's own loop with local models: yes, later.** Defer it until S-0508 shows the container path works.
- **Proposal 2, the signed `ghcr.io` image and its release workflow, and an OpenCode layer: yes, later.** Defer them until the local images are stable and E-0023 lands.
- **Proposal 3, drop the Docker Agent measurement from S-0505: yes.** ADR-0202 was rejected, so measuring Docker Agent as the inner loop answers a question already closed.

**Why it is yours.** Proposals 1 and 3 change scope, and finalizing E-0029's drafts waits on TH-0657's input.

**What I do after.** Once TH-0657 is answered and the planner writes the values, I finalize under finalize_drafts each draft that `flai promote --drafts` lists as complete and consistent.

Source: design/conventions/strategic-agents.md § As the planner
```

The operator, alex, answered TH-0659 again at 2026-10-26T12:26:58Z.

The answer:

```text
TH-0657 is answered
```

Source: TH-0659.
