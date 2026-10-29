---
id: TH-0753
title: "Criterion 2: a real thread's pack is 31% of the general pack, under a quarter only for the test's small fixture thread"
anchor:
  path: wip/kanban/stories/S-0537-flai-prime-role-orchestrate-takes-the-run-s-reasons-and-primes-for-them-the-role-s-conventions-once-and-per-reason-what-it-needs-instead-of-the-general-pack.md
  item: S-0537
status: resolved
participants: [agent-S-0537, orchestrator, alex]
created: 2026-10-29T16:28:31Z
updated: 2026-10-29T17:13:14Z
---

# TH-0753 Criterion 2: a real thread's pack is 31% of the general pack, under a quarter only for the test's small fixture thread

On wip/kanban/stories/S-0537-flai-prime-role-orchestrate-takes-the-run-s-reasons-and-primes-for-them-the-role-s-conventions-once-and-per-reason-what-it-needs-instead-of-the-general-pack.md.

## Entries

### 2026-10-29T16:28:31Z agent-S-0537
Criterion 2's test passes. A pack primed for one thread is under a quarter of the general pack, measured over this repository's conventions and design. But its fixture thread has one entry, on a three-line story. Packs for real threads here are about a third. I have left criterion 2 unticked until you choose.

| Reason, measured 2026-10-29 | Pack | Of the general pack (176 KB) |
|---|---|---|
| `thread:TH-0750` | 54.6 KB | 31% |
| `thread:TH-0719` | 58.1 KB | 33% |
| `ready` | 47.4 KB | 27% |
| `release` | 29.8 KB | 17% |
| `draft:S-0542` | 63.0 KB | 36% |
| `review:S-0536` | 80.2 KB | 45% |

The TH-0719 pack is made of:

- the conventions cut by reason, 27.4 KB, down from 64.7 KB;
- the thread, 6.3 KB;
- its story whole, 11.0 KB;
- the digest's answers, 8.4 KB, within ADR-0194's eighth of the budget;
- the header and left-out list, about 5 KB.

**Recommended: trim the two largest per-reason parts, and make the test's thread a real one.**

- A thread loads its item as its title, status, goal, and acceptance criteria, not whole. That leaves out the notes, tasks, and planning, which `item_get` gives when the orchestrator needs them. This saves about 7 KB on a planned story.
- A reason pack's digest is capped at a sixteenth of the budget, about 5 KB, still within ADR-0194's eighth. This saves about 4 KB.
- The test measures a pack for a copy of a real thread of this repository on its story, not a one-entry fixture.
- Real thread packs would then be about 42 to 44 KB, just under a quarter. I would record the change in a new ADR refining ADR-0217 §3 and §4, and tick criterion 2 on the new test.

**Or: keep the pack as it is.** I would tick criterion 2 on the test as written. The table above is already recorded in `design/system/agent-context.md`, and I would note the real-thread figure in S-0537's notes.

Either way, criteria 1, 3, and 4 are done. I am running the close-out meanwhile.

### 2026-10-29T16:29:43Z orchestrator (recommendation)
**Recommended: trim the pack as S-0537's agent proposes, and tick criterion 2 on a test over a real thread.**

- **The item a thread is on loads as its title, status, goal, and acceptance criteria.** Its notes, tasks, and planning come from `item_get` when a reply needs them. ADR-0209 §3 asks for "the thread, its item", not the item whole.
- **The digest in a reason pack is capped at a sixteenth of the budget.** That stays within ADR-0194's eighth.
- **The test measures a copy of a real thread of this repository on its story,** not a one-entry fixture.
- **Record it in a new ADR refining ADR-0217 §3 and §4.** `design/adrs` is already in S-0537's touches.

**Why:** criterion 2 is there so that a run for one thread reads much less than the general pack (ADR-0209 §3). A test that passes on a three-line story while real threads measure 31% to 33% meets the criterion's words but not its purpose. The trim reaches about 42 to 44 KB, under a quarter, and changes no other pack.

**Or:** keep the pack as it is. The agent then ticks criterion 2 on the test as written and notes the real-thread figures in S-0537's Notes.

The review pack, at 45%, is not criterion 2's concern. It can go to S-0536's follow-up story under E-0030 if it matters there.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

### 2026-10-29T16:31:31Z alex
Confirmed the recommendation of 2026-10-29T16:29:43Z orchestrator.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

### 2026-10-29T17:13:08Z agent-S-0537
Done as confirmed, in T-2339 and ADR-0218.

- A thread's item now loads as its title, status, goal, and criteria.
- A reason pack's digest is capped at a sixteenth of the budget.
- The test measures copies of a real thread, TH-0750, and its story, S-0542.

Measured against the general pack of 177 KB:

| Thread | Thread size | Pack before | Pack after |
|---|---|---|---|
| TH-0750 | 4.9 KB | 31.0% | 23.5% |
| TH-0719 | 7.0 KB | 33.6% | 25.2% |
| TH-0717 | 8.9 KB | 35.4% | 27.0% |

The thread itself is what varies. Of the repository's 754 threads, the median is 2.8 KB and four in five are under 5 KB, so a pack for a thread of up to about 6.5 KB comes in under a quarter, roughly nine in ten threads. I took TH-0750 for the test because it is larger than four in five threads, not because it passes. `design/system/agent-context.md` records the table and the distribution. Criterion 2 is ticked.

### 2026-10-29T17:13:14Z agent-S-0537
Resolved: The operator confirmed the trim at 16:31Z. T-2339 and ADR-0218 did it, and criterion 2 is ticked on the test over TH-0750.
