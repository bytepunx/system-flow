---
id: ADR-0312
title: "A tier belongs to a suite, and the workflow file names the suites a transition needs, a state runs on entry with the cluster's deploy and checks, and the batch, the publish, and CI run"
status: proposed
date: 2026-11-16
supersedes: []
superseded_by: []
refines: [ADR-0174, ADR-0254, ADR-0291, ADR-0301]
---

# ADR-0312 A tier belongs to a suite, and the workflow file names the suites a transition needs, a state runs on entry with the cluster's deploy and checks, and the batch, the publish, and CI run

## Context

Nothing ties a test tier to a workflow step. Which tiers run is decided by each caller's Go code and four flags on the tier in the manifest (, , , ), and the workflow file ([ADR-0301](0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md)) cannot name a tier, a verification, or a cluster run: a transition holds , , , , , and , its rules are a fixed vocabulary none of which reads the verify report, and the parser refuses any other key.

So in this repository a story that changes  runs the same suites three or four times:

- at its close-out before review: the static tiers, , and , about twelve minutes, since  runs the  tiers there even while the queue batches ([ADR-0233](0233-while-the-merge-queue-can-batch-the-close-out-runs-a-queue-tier-whose-when.md));
- again in review whenever main brought in an overlapping path ([ADR-0291](0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md), [ADR-0303](0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md)), which on 2026-11-15 held every story in review but one;
- again in the merge queue's batch, which runs every tier its paths select ([ADR-0174](0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md));
- again in CI once main is published, from commands the workflow files hold apart from flai.

Publishing runs no tier. The cluster's deploy and checks run on entering a state of kind  ([ADR-0254](0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md)), which is code, not the file. A project cannot say that review needs only its behaviour tests, that an integration state between review and done deploys and runs the slow suites, or that publish runs them.

The operator asked on 2026-11-16 that this be possible, in the workflow file, with the template giving each step flai ships its default, and decided: steps name suite labels, not tiers, so that a project can add suites of its own; a cluster step is skipped where no cluster or deploy is configured; and this makes the customisation possible without changing what this repository runs.

## Decision

A tier belongs to a suite, and the workflow file says which suites each step needs and runs.

1. **Suites.** A tier in the manifest's  may give , a kebab-case label; a tier that gives none is the suite of its own name. Several tiers may share a suite. A suite is whatever a project calls it; flai gives no suite a meaning. The template's manifest labels its tiers , , , and .  runs a suite's tiers. A  tier belongs to its suite and still runs only when named.
2. **What a move needs.** A transition may give , a list of suites, or . A story's move by that transition needs its last  to have passed, at the head of its branch, every tier of those suites that its paths select. It is a rule of the move, , with a severity as every rule has,  unless the file says . The close-out an agent asks for before a move runs the suites the story's next transitions need, in place of every tier.
3. **What the acceptance needs.** The  of the transition with the  effect is what the acceptance's gates judge: the merge queue's eligibility,  ([ADR-0174](0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md)), and  ([ADR-0291](0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md), [ADR-0303](0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md)) read only tiers of the suites it lists. A suite the acceptance does not need never holds a story in review or sends it back.
4. **What a state runs.** A state may give , with , a list or , and , any of  and . When a story enters the state, flai serve rebases its branch onto the main branch and runs them as a job, with its output kept and its progress given to the dashboard, as the Rebase and resolve job does (S-0804). A job that passes leaves the story in the state, for the transition out of it; a rebase that stops or a step that fails sends the story back by the state's  transition, with the thread its agent is started on. The cluster's deploy and checks run where a state's  names them, no longer on every state of kind .
5. **No cluster.** A  step is skipped, and the job's record says so, when the project configures no deploy, the kluster host action is off, or no agent namespace is available. It never holds a story for want of a cluster.
6. **Points outside an item's moves.** The file's top-level  gives suites, or , for three points no transition stands for: , what the merge queue's batch runs over the merge of its stories; , what published , pushed to origin runs at the main branch's head before it tags, refusing to publish when one fails; and , what  runs, for a project's CI to call in place of commands of its own.
7. **Reuse.** The verify report records each suite's tiers and the head and base they passed at, so a later step runs again only what the paths since send back ([ADR-0156](0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md)).
8. **The baseline is today's behaviour.** flai's baseline and the template's file give  at  on the move to review,  on the acceptance, , , no , and  on . A project that edits nothing runs what it ran before. A tier's  and  keep their meaning until a project's file names suites for , and are then retired by a story of their own.
9. **The check.** design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md:8: warning: issues.no-story: I-0327 has been open since 2026-11-09 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0327
wip/kanban/stories/S-0753-the-dashboard-s-board-draws-the-columns-their-kinds-and-the-moves-between-them.md:1: warning: wip.after-hub: S-0753 is named in after by 3 open stories, S-0754, S-0756 and S-0758, and each of them waits until it is accepted; after names only a story whose accepted change the story needs on main, and the pull order sets an order of preference, so drop S-0753 from the after of each that does not need its change (flai edit <id> --after) (ADR-0188)
wip/kanban/stories/S-0804-rebase-and-resolve-rebases-the-story-s-branch-and-runs-its-close-out-as-a-job.md:19: warning: touches.folder: S-0804's claim keeps the folder template, which holds 87 files, and holds no ready story now; name the files its tasks change with flai touches, or say under ### Planning in S-0804's Notes why files may be added there (ADR-0188)
wip/threads/TH-1232-the-merge-queue-hold-th-1231-comes-from-the-host-s-flai-1-98-0-which-cannot.md:7: warning: threads.archived: TH-1232 is open but S-0748 is archived; resolve it with flai thread resolve TH-1232
system-flow.yaml:45: warning: agent.image: the harness image flai-agent-claude-code:local of the inner harness claude-code is stale: built for flai 1.86.12, not dev; flai serve renews it itself at the next start that needs it, unless `flai serve agent image auto off` turned its image job off, or renew it now with `flai image ensure`; see `flai image status`
system-flow.yaml:45: warning: agent.image: the image flai-agent-sf-claude-code:local of agent.image.layer flai/images/system-flow/Dockerfile is stale: built for flai 1.86.12, not dev; carries no dev.system-flow.recipe label, so what it was built from is not known; carries no dev.system-flow.base label, so what it was built from is not known; flai serve renews it itself at the next start that needs it, unless `flai serve agent image auto off` turned its image job off, or renew it now with `flai image ensure`; see `flai image status`
4254 items checked, 0 errors, 6 warnings (4 that --strict passes over: review over its limit, an epic behind its stories, and a built-in role file the project lacks, which only the operator clears, or a narrative not yet written, which flai stream state writes, or component examples left unchecked for want of node, a Svelte sub-project, or a sub-project's svelte, or the project's image layer not built or left unchecked for want of docker, which only building it clears, or a story in progress keeping a folder touch or a story three or more open stories name in after, which only their stories clear, or two stories in progress whose claims overlap, which the trial merge judged (ADR-0230)) reports a suite the file names that no tier has, as an error, and a suite of a tier that is not  which no transition, state, or point names, as a warning, so that no suite stops running unseen.

## Consequences

- A project can keep its close-out to its fast suites and run the slow ones once: in an integration state between review and done, in the batch, or at publish. This repository changes nothing until its own file says so.
- The acceptance's gates become data. A project whose acceptance needs only  is not rebased and verified again in review for a change to main that only  would catch; that change is caught where the slow suites next run.
- A suite left to a later step finds a break later, on the merge of several stories, and the batch or the publish that finds it sends more back at once. A project chooses that with its file.
- The workflow file gains keys, so a flai older than this refuses a file that uses them ([ADR-0301](0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md)'s strict parse); the template's file must not use them before the release that reads them is the oldest a project runs.
- The cluster runs move from a kind to the file. The  rule, which the move does not evaluate today, becomes the  rule's neighbour and is checked where a transition names it.
- CI stays the forge's; flai gives it one command, and a project's workflow files call it or do not.

## Alternatives considered

- **Steps name tiers.** Rejected by the operator: the template's file cannot know a project's tiers, and a project could not add a suite without editing every step.
- **A stage key on each tier in the manifest, in place of the workflow file.** Rejected: the manifest does not know the project's states, and one tier may be needed at several steps.
- **Rules alone, with no .** Rejected: a rule checks and cannot run, so an integration state would need an agent or the operator to ask for each run.
- **Hold a story when a state names a cluster step and no cluster is there.** Rejected by the operator: the template's default would stop every project without a cluster.
- **Change this repository's steps now.** Rejected by the operator for now:  runs the slow suites at the close-out because a fifth of the batches sent a story back ([ADR-0233](0233-while-the-merge-queue-can-batch-the-close-out-runs-a-queue-tier-whose-when.md)).
