---
id: D-0313
title: "S-0384 scope: add the design's two extra criteria (tweak a proposal on its branch, the guidance's threads by heading): Recommendation: answer \"leave out\""
kind: answer
sources: [TH-0457]
topics: [design, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/designsystem/proposals.go, flai/internal/designsystem/proposals_test.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts, flaiover/src/routes/api/design/proposals/+server.ts, flaiover/src/routes/api/design/proposals/proposals.test.ts, flaiover/src/lib/design/diff.ts, flaiover/src/lib/design/diff.test.ts, flaiover/src/lib/components/design/ProposalCompare.svelte, flaiover/src/lib/components/design/ProposalCompare.svelte.test.ts, flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, flai/internal/designsystem/css.go, docs/operators/settings.md]
date: 2026-10-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0313 S-0384 scope: add the design's two extra criteria (tweak a proposal on its branch, the guidance's threads by heading): Recommendation: answer "leave out"

The operator, alex, answered TH-0457 at 2026-10-15T04:12:07Z.

What was asked:

```text
The design is larger than the story you finalized. `design/system/design-system-tool.md` § "S-0384 Proposals and threads" asks the replan to change S-0384 in three ways. The planner did not apply them.
```

The answer:

```text
**Recommendation:** answer "leave out". S-0384 finishes as written, with the two refinements its agent is applying. E-0020's planner then drafts the two criteria as backlog stories when it replans the epic, which TH-0438 already asks for.

## Why

- **The first one cannot be checked end to end yet.** Tweaking a proposal's tokens on its branch goes through `file.save` with `story` (§ S-0384 Proposals and threads). No flai has `file.save` until the `file.get`/`file.save` story lands, and that story is not drafted yet. In S-0384 it could only be built against a stub, as S-0383's save was.
- **Both are additive.** Neither blocks the comparison or the request form, which are what S-0384's criteria ask for.
- **S-0384 is already eight tasks in five layers.** Adding two widens its claim to `api/design/tokens/+server.ts` and `Threads.svelte`, which would hold other stories.

## Why this needs you

- S-0384 is finalized, so adding criteria is yours.
- The two new stories come from E-0020's planner. I cannot start it: `flai plan --candidates` does not list E-0020. Run `flai plan E-0020` once, and it drafts these two along with the file story and the component-design story.

Source: design/system/design-system-tool.md § S-0384 Proposals and threads
```

Source: TH-0457.
