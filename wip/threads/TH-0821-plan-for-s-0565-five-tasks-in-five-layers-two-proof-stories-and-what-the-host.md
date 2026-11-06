---
id: TH-0821
title: "Plan for S-0565: five tasks in five layers, two proof stories, and what the host needs first"
anchor:
  path: wip/kanban/stories/S-0565-one-story-of-this-repository-is-worked-end-to-end-by-flai-serve-on-harness.md
  item: S-0565
status: resolved
participants: [planner-S-0565, orchestrator, flai serve]
created: 2026-11-03T01:26:01Z
updated: 2026-11-05T01:32:35Z
---

# TH-0821 Plan for S-0565: five tasks in five layers, two proof stories, and what the host needs first

On wip/kanban/stories/S-0565-one-story-of-this-repository-is-worked-end-to-end-by-flai-serve-on-harness.md.

## Entries

### 2026-11-03T01:26:01Z planner-S-0565
Plan for S-0565, written by planner-S-0565.

## Tasks and layers

The work is one chain: each run depends on what the one before it set up or found.

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-2449 Results document started; host preconditions checked; both proof stories set to `harness: container` | none | results document |
| 2 | T-2450 A stop from the dashboard and a restart after an answered thread, each tried once on the first proof story's container, recorded in S-0565's Notes | T-2449 | S-0565's story file |
| 3 | T-2451 First proof story worked to review in its container; usage, card status, and issues recorded | T-2450 | results document, `design/issues/summary.md` |
| 4 | T-2452 Go proof story worked to review on this repository's image layer; close-out passes every selected tier inside the container | T-2451 | results document, `design/issues/summary.md` |
| 5 | T-2453 Results and recommendation written; experiments index row added | T-2451, T-2452 | results document, `design/experiments/README.md` |

## Story figures

- Touches: the results document `design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md`, `design/experiments/README.md`, `design/issues/summary.md`, and the folder `design/issues`. The folder is kept because `flai issue new` names the issue files.
- Forecast: 1h30m, raised from flai's 9m. flai's size counted no touches and no waiting.
- Delivery: 2026-11-03T08:48Z.
- Cost of delay: 52.12 USD a week, flai's share of E-0029's penalty input.
- The reasons are under `### Planning` in the story's Notes.

## Assumptions

1. The proof stories are other stories, worked by their own agents in the container. S-0565's agent sets them up, watches them, and records what happened.
2. Two proof stories, run one after the other. The first changes only documentation or shell and covers criteria 1 and 2. The second changes Go and covers criterion 3. One Go story could cover all three criteria, but a failure would then hide which layer broke.
3. Before S-0565 can start, the installed `flai` must be a published release that holds S-0508, S-0510 (in progress), and S-0517, and the host must be upgraded. The images must be built (`make agent-images` and this repository's layer). The host's container harness must be set (`flai serve agent harness container --login ~/.claude`, the image, the network). These are the operator's to do: `flai guard` refuses an agent every write to the host configuration (ADR-0220). T-2449 checks each one and asks on this thread for any that is missing.
4. Moving a proof story to ready is the operator's or the orchestrator's job. S-0565's agent only sets its harness.
5. The stop is the operator's, from the dashboard. The restart after a thread needs the proof story's agent to end its session on a question. If it asks nothing, S-0565's agent asks here how to bring one about.

## Questions for the operator

1. Which proof stories? My recommendation: S-0558 (safety.md against git.md, documentation only) for criteria 1 and 2, and S-0543 (a small Go test fix in `flai/internal/serve`) for criterion 3. Both are drafts, so you would finalize them, or name others.
2. S-0520 (the close-out stops inside the container because `scripts/flai.sh` builds flai with go) is not in `after`. I assume S-0510's layer carries go and supersedes it. If it does not, criterion 3 needs S-0520 first, and I would add it to `after`.

Nothing else is proposed: no task to split, merge, or drop.

### 2026-11-03T01:27:54Z orchestrator (recommendation)
**Recommended:** approve the plan with three changes.

1. Add `S-0520` to S-0565's `after`.
2. Use S-0556 as the Go proof story for criterion 3, not S-0543.
3. Keep S-0558 for criteria 1 and 2.

Plan and finalize both proof stories before S-0565 starts.

## Question 2: S-0520 is still needed

S-0510's layer carries go (`flai/images/system-flow/Dockerfile` copies `/usr/local/go`). S-0518 is accepted, so `scripts/flai.sh` also runs without go. That clears I-0218's exit 127.

It does not let the close-out pass inside the container:

- S-0510's Notes say `scripts/close-out.sh` was "Not run inside", because "only the host commits (S-0517, S-0520)".
- S-0520's planning says what is still missing: `flai verify` inside the container has no main checkout. S-0520 mounts its `wip/`, its `system-flow.yaml`, and its `.flai-cache/verify/`.

Without S-0520, criterion 3 stops at the close-out's first step.

## Question 1: S-0543 cannot be the Go proof

On TH-0813 I recommended cancelling S-0543 as delivered by S-0418 and S-0548. I-0233 is closed. Its one task would change no Go, so it cannot show a Go close-out passing.

Don't use the timing remediations (S-0560, S-0561, S-0569) either. Their tests are the ones that fail under load, so a failure would not show which layer broke.

S-0556 is a better fit: a small, deterministic change to what `flai guard` lists among a sub-agent's reads.

S-0558 and S-0556 are both drafts with no plan yet. Each needs its planner, then finalizing.

## The rest of the plan fits

- Five tasks in one chain. Each run depends on what the previous one found.
- Touches: the results document, the index, and `design/issues` with its summary.
- Forecast 1h30m and value 52.12 USD a week, both as flai gives them.
- **One small fix:** T-2450 touches S-0565's own story file under `wip/kanban`. The Notes say criterion 2's record is written with `flai edit` in the main checkout and is "no touch of its branch". T-2450 should drop that touch, or the Notes should change.

## What stays yours

Assumption 3: publishing a release with S-0508, S-0510, and S-0517, upgrading the host, building the images, and setting the container harness.

**Why I recommend rather than answer:** choosing the proof stories and adding an `after` change S-0565's scope. That is your call.

Source: wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md

### 2026-11-03T01:28:14Z planner-S-0565
Update from planner-S-0565: question 2 is settled from the code, and criterion 3 raises a new question.

## S-0520 added to `after`

S-0510's layer carries Go, golangci-lint, ShellCheck, and pnpm. It does not make the close-out run inside the container. S-0510's own notes say `scripts/close-out.sh` was not run inside, and leave that to S-0517 and S-0520. The close-out's first step, `flai verify`, reads the main checkout's `wip/` and `system-flow.yaml`, and writes under `.flai-cache/verify/`. The container mounts none of these until S-0520 adds them (I-0218). Criterion 3 reads that output, so S-0565's `after` is now S-0508, S-0510, S-0517, and S-0520. The reason is under `### Planning`.

## Question 3: where criterion 3's close-out runs

Under `flai serve`, the story's agent asks for the close-out with the MCP tool `verify`. `flai serve` then runs `scripts/close-out.sh` as a job on the host, in the story's worktree (ADR-0229, `flai/internal/serve/verifyjob.go`). S-0510 does not change that. No open story makes the job run inside the container. So "passes inside the container" cannot happen under `flai serve` as built.

My recommendation is (a):

- **(a) Two checks, no new story.** The Go proof story's close-out runs as `flai serve`'s usual host job. T-2452 then also runs `scripts/close-out.sh` once by hand inside the layered image on that story's worktree, with `scripts/agent-container.sh S-nnnn -- scripts/close-out.sh S-nnnn` and S-0520's mounts in place. It records the tiers that ran and their outcome. The criterion's words stay as they are.
- **(b) A new story first.** A new E-0029 story makes `flai serve` run a container story's close-out job inside its container, on the layer image. S-0565 waits for it in `after`. I would draft it on your word.

T-2452's Work already allows for (a). On (b) I would add the new story to `after` and rewrite T-2452.

### 2026-11-03T01:28:46Z planner-S-0565
The orchestrator's points are taken into the plan:

- S-0520 is in S-0565's `after`. T-2449 now also checks that the installed release holds S-0520.
- T-2450 no longer touches S-0565's story file. Its record goes into S-0565's Notes through `flai edit` in the main checkout, as the Notes say.
- T-2452 now runs `scripts/close-out.sh` once by hand inside the layered image, as option (a) describes, unless the answer here names a story that moves the close-out job into the container.
- The tasks name no proof story. T-2449 takes whichever ones this thread settles: S-0558 and S-0556 as the orchestrator recommends, or others you name.

Open for you: question 1 (the proof stories), question 3 (option (a) or (b)), and the host preconditions in assumption 3.

### 2026-11-03T01:30:00Z orchestrator (recommendation)
**Recommended on question 3:** (a), two checks and no new story. The rest of my recommendation of 01:27:54Z stands: S-0558 for criteria 1 and 2, S-0556 for criterion 3, both planned and finalized before S-0565 starts.

## Why (a)

- ADR-0229 runs the close-out as a flai serve job in the story's worktree, on the host. Nothing open changes that, so criterion 3's "passes inside the container" cannot come from the job.
- (a) meets the criterion's words without new scope: the Go story is worked by `flai serve` in its container on the layer, and `scripts/close-out.sh` also runs once by hand inside the layered image on that story's worktree. T-2452 already says this.
- (b) adds a story and another `after` to an experiment whose job is to find limitations. If the host-side job turns out to be one, T-2452 files it as an issue, and that issue can become the story (b) describes.

## What changed since my last entry

- S-0520 is in `after` now. That settles my point 1.
- T-2450 no longer touches S-0565's story file. That settles the small fix I raised.

## Why I recommend rather than answer

Choosing between (a) and (b), and the proof stories, is S-0565's scope, so it is your call. I leave S-0565 a draft until this thread is answered.

Source: design/adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md § Decision

### 2026-11-05T01:30:53Z flai serve
Confirmed the recommendation of 2026-11-03T01:30:00Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md § Decision

### 2026-11-05T01:32:25Z orchestrator
The plan is approved as confirmed. flai serve confirmed my recommendation of 2026-11-03T01:30:00Z.

## What is settled

- Question 1: S-0558 is the proof story for criteria 1 and 2. S-0556 is the Go proof story for criterion 3.
- Question 2: S-0520 is in S-0565's `after`. That was done already.
- Question 3: option (a). T-2452 runs `scripts/close-out.sh` once by hand inside the layered image, and no new story is added.
- The tasks, the touches, the 1h30m forecast, and the 52.12 USD a week value fit the story.

## What I did

- Started the planner for S-0558 (planner-S-0558) and for S-0556 (planner-S-0556). flai plan --candidates lists both as having no plan.

## What waits

- S-0565 stays a draft until both proof stories are planned and finalized, as the confirmed answer asks.
- Assumption 3 is still yours: publish a release that holds S-0508, S-0510, S-0517, and S-0520, upgrade the host, build the images, and set the container harness. T-2449 checks each one.

Source: design/adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md § Decision

### 2026-11-05T01:32:35Z orchestrator
Resolved: Plan approved. flai serve confirmed the recommendation: proof stories S-0558 and S-0556, option (a), and S-0520 in after. The planners for S-0558 and S-0556 have started.
