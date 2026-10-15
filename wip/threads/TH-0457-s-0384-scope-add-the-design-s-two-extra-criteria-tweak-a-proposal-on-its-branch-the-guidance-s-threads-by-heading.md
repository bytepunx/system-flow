---
id: TH-0457
title: "S-0384 scope: add the design's two extra criteria (tweak a proposal on its branch, the guidance's threads by heading)?"
anchor:
  path: wip/kanban/stories/S-0384-an-agent-s-proposal-for-the-design-system-is-a-story-whose-branch-changes-it-and-the-design-page-previews-that-branch-beside-the-main-branch-s-before-the-operator-accepts-it.md
  item: S-0384
status: open
participants: [agent-S-0384, orchestrator]
created: 2026-10-15T03:17:24Z
updated: 2026-10-15T03:17:34Z
---

# TH-0457 S-0384 scope: add the design's two extra criteria (tweak a proposal on its branch, the guidance's threads by heading)?

On wip/kanban/stories/S-0384-an-agent-s-proposal-for-the-design-system-is-a-story-whose-branch-changes-it-and-the-design-page-previews-that-branch-beside-the-main-branch-s-before-the-operator-accepts-it.md.

## Entries

### 2026-10-15T03:17:24Z agent-S-0384
The design is larger than the story you finalized. `design/system/design-system-tool.md` § "S-0384 Proposals and threads" asks the replan to change S-0384 in three ways. The planner did not apply them.

**Refinements. I am applying these now, because they fit the criteria as written:**

- Criterion 2: a proposal is a branch that changes the token file, the guidance, or both.
- Criterion 4: the form's story gets fixed criteria (the tokens as asked, the stylesheet in step, the UI's tests passing). Its touches are the token file, plus the stylesheet once the manifest names one (S-0385). It is moved to ready and is not a draft.

**Two new criteria, a task each. Not in the story:**

1. On the Design page, the operator edits a proposal's tokens on its branch through `file.save` with `story`. This touches `flaiover/src/routes/api/design/tokens/+server.ts`.
2. The Design page lists the guidance's threads by heading and opens one with `thread.new`. This touches `flaiover/src/lib/components/Threads.svelte`.

**Recommendation: leave both out of S-0384 and make each a backlog story under E-0020.** S-0384 is already eight tasks in five layers. Both are additive, and neither blocks the comparison or the form.

Answer one of:

- **leave out**: I finish S-0384 as written. You or the planner make the two stories.
- **add both** (or name one): I add the criteria and a task for each, and widen the touches.

### 2026-10-15T03:17:34Z orchestrator (recommendation)
**Recommendation:** answer "leave out". S-0384 finishes as written, with the two refinements its agent is applying. E-0020's planner then drafts the two criteria as backlog stories when it replans the epic, which TH-0438 already asks for.

## Why

- **The first one cannot be checked end to end yet.** Tweaking a proposal's tokens on its branch goes through `file.save` with `story` (§ S-0384 Proposals and threads). No flai has `file.save` until the `file.get`/`file.save` story lands, and that story is not drafted yet. In S-0384 it could only be built against a stub, as S-0383's save was.
- **Both are additive.** Neither blocks the comparison or the request form, which are what S-0384's criteria ask for.
- **S-0384 is already eight tasks in five layers.** Adding two widens its claim to `api/design/tokens/+server.ts` and `Threads.svelte`, which would hold other stories.

## Why this needs you

- S-0384 is finalized, so adding criteria is yours.
- The two new stories come from E-0020's planner. I cannot start it: `flai plan --candidates` does not list E-0020. Run `flai plan E-0020` once, and it drafts these two along with the file story and the component-design story.

Source: design/system/design-system-tool.md § S-0384 Proposals and threads
