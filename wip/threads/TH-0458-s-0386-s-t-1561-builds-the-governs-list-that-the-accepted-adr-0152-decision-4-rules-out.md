---
id: TH-0458
title: "S-0386's T-1561 builds the governs list that the accepted ADR-0152 § Decision 4 rules out"
anchor:
  path: wip/kanban/stories/S-0386-a-story-s-agent-working-on-a-project-s-ui-is-primed-with-its-design-system-and-a-baseline-convention-tells-it-to-build-from-the-tokens-and-propose-changes-on-its-branch.md
  item: S-0386
status: resolved
participants: [orchestrator, agent-S-0386]
created: 2026-10-15T04:21:27Z
updated: 2026-10-15T04:39:50Z
---

# TH-0458 S-0386's T-1561 builds the governs list that the accepted ADR-0152 § Decision 4 rules out

On wip/kanban/stories/S-0386-a-story-s-agent-working-on-a-project-s-ui-is-primed-with-its-design-system-and-a-baseline-convention-tells-it-to-build-from-the-tokens-and-propose-changes-on-its-branch.md.

## Entries

### 2026-10-15T04:21:27Z orchestrator
**Recommendation:** drop T-1561 and the `init --project` governs list. Select stories by the guidance's `topics` alone, as ADR-0152 decides.

## What conflicts

- **ADR-0152 § Decision 4 is accepted:** "One list. The guidance's `topics` alone say which stories a design system reaches. There is no second record of governed sub-projects."
- **T-1561 builds that second record:** "flai design init --project records in the design system document which sub-projects it governs". It moved to in-progress at 04:21Z.
- **TH-0438 says the same.** S-0380's list of what each draft must change says, for S-0386: "Drop `init --project`, the governs list, and T-1561. The guidance's `topics` are the only list." It also says to drop the `flai/cmd/design.go` touches and add `context/story.go` and `conventions/conventions.go`. Its convention gets `roles: [story, verify]` (§ Decision 5).
- **S-0386 was never replanned.** The operator finalized it from the pre-decision draft, so its tasks still carry the old shape.

## Why this needs you

S-0386 is finalized, so cancelling T-1561 or changing its criteria to match ADR-0152 is yours. Its agent may decide to follow the accepted ADR and record that in its narrative's Decisions. Whatever you choose, the agent should stop building the governs list until you say.

### 2026-10-15T04:39:41Z agent-S-0386
S-0386's agent here. T-1561 builds no second list.

**What it does.** `flai design init --project <name>` checks that the name is a sub-project the manifest lists. It then writes that name into the guidance's `topics`, the one list ADR-0152 § Decision 4 keeps. With `--project` and no explicit `--topics`, the default `all` is dropped. Nothing else records governed sub-projects.

**Why keep it.** It meets criterion 1 as the operator finalized it, within the ADR. It also refuses a sub-project name the manifest does not know, which `--topics` cannot do.

**Recommendation:** keep T-1561 as built. This is recorded in the narrative's Decisions. If the operator would rather drop it, I will revert `--project` and its tests and tick criterion 1 on `--topics` alone.

Source: design/adrs/0152-agents-build-ui-from-a-tailwind-v4-stylesheet-flai-generates-from-the-tokens.md § Decision

### 2026-10-15T04:39:50Z orchestrator
Resolved: No conflict: T-1561's --project validates a sub-project name against the manifest and writes it into the guidance's topics, the single list ADR-0152 § Decision 4 keeps. No second record of governed sub-projects is made. My concern came from the task's title. The agent records this in its narrative's Decisions.
