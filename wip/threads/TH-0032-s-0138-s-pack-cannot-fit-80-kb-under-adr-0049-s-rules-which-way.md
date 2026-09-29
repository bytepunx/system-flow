---
id: TH-0032
title: "S-0138's pack cannot fit 80 KB under ADR-0049's rules: which way?"
anchor:
  path: wip/kanban/stories/S-0146-a-story-s-context-pack-fits-a-size-budget-with-briefs-of-what-its-topics-and-links-select.md
  item: S-0146
status: answered
participants: [agent-S-0146, alex]
created: 2026-09-29T03:27:13Z
updated: 2026-09-29T04:10:22Z
---

# TH-0032 S-0138's pack cannot fit 80 KB under ADR-0049's rules: which way?

On wip/kanban/stories/S-0146-a-story-s-context-pack-fits-a-size-budget-with-briefs-of-what-its-topics-and-links-select.md.

## Entries

### 2026-09-29T03:27:13Z agent-S-0146
The budget, briefs, and ranked fill are built. Measured with `flai prime --story S-0138` on the S-0146 branch (header included):

| Part | Size | Cut for the budget? |
|------|------|---------------------|
| Conventions | 48.0 KB | never |
| Open-issues table | 9.1 KB | never (part of the conventions step) |
| Named: ADR-0047, ADR-0049, `design/system/conventions.md` | 38.3 KB | never |
| Briefs: 22 design and tech files, 33 ADR decision lines | 23.0 KB | |
| Catalog, header | 4.9 KB | |
| Pack | 125 KB | nothing ranked; header says "over budget: named" |

Two things differ from ADR-0049's estimate. S-0138 was re-scoped after the estimate and now names ADR-0049 and `conventions.md` as well as ADR-0047 (38 KB, not 9 KB). And the estimate left out the 9 KB open-issues table. With only ADR-0047 named, the pack would still be about 95 KB. So the criterion "S-0138 measures under 80 KB" cannot be met without cutting something ADR-0049 says is never cut.

Options:

A (recommended). Keep ADR-0049's rules and change the criterion to "S-0138's pack is measured and the measurement is in the story notes, and a story that names less fits". I record S-0138 at 125 KB and a cli story that names one ADR under 80 KB with ranked sections. The room comes when you narrow the conventions' topics, as ADR-0049 plans, or from a later story that scopes the issues table to the story.

B. Raise the default budget to 128 KB. S-0138 would fit without ranking, but the pack would exceed an MCP tool result (about 100 KB), which is what the budget is for.

C. Named documents over a size (say 12 KB) load as a brief plus the sections the link points at, not whole. S-0138 fits. This reverses "never cut" in ADR-0049, so it needs a new ADR, and S-0146 would grow.

I carry on with the tests, the `flai check` warning, and the documents meanwhile.

### 2026-09-29T03:27:34Z agent-S-0146
A correction to option A. A cli story that names nothing does not fit 80 KB either. Conventions (48 KB), the issues table (9 KB), the briefs (23 KB), the catalog, and the header come to about 85 KB before anything named or ranked. ADR-0049 estimated the briefs at 16 KB and did not count the issues table or a header. Revised options:

A (recommended). Keep 80 KB and let the briefs yield to the budget, as ranked sections already do. Conventions and what is named are never cut. If the briefs do not fit, ADR briefs move to the catalog, the ones found last going first. If that is not enough, design briefs lose their first paragraph and keep their title and outline. The header says how many briefs were cut. A cli story that names up to about 20 KB then fits, with little or no room for ranking until the conventions are narrowed. S-0138 names 38 KB, so it stays over, and I record it as over in the notes.

B. Raise the default to 128 KB. Every brief stays, and S-0138 fits without ranking. `--json` is a little larger than the text, so the MCP `prime` result would exceed the roughly 100 KB a tool result carries.

C. Keep 80 KB and every brief, and accept that every cli story is over budget until you narrow the conventions' topics. The header says so on every pack.

Whichever you choose, the S-0138 criterion as written ("under 80 KB") can only be met by B.

### 2026-09-29T03:51:59Z alex
C - as long as the agent will read the full document if it determines from the brief that it is relevant

### 2026-09-29T03:52:31Z agent-S-0146
Done as C. The budget stays 80 KB and every brief is kept. When the briefs take a pack over the budget, `exceeded` is `briefs`, nothing is ranked, and the header says so and names narrowing the conventions' topics as the way to make room. On reading the full document: the `briefs` and `decisions` headings in the pack now say that a brief is not the document. When one bears on the story, the agent reads it whole with `doc_get` or `flai doc show` before relying on it or changing what it describes. The prompt, `CLAUDE.md`, and `session-start.md` wording is S-0148's, and I have noted there that it must say the same. I am rewording the S-0138 criterion to "measured, and the measurement is in the notes", per this answer, and recording the numbers.

### 2026-09-29T03:52:59Z agent-S-0146
One more finding from the measurements, for you to weigh. I am not waiting on it. What stories name is what takes packs over the budget, more than the briefs do:

| Story | Pack | Named | What |
|-------|------|-------|------|
| S-0138 | 125 KB | 36 KB | ADR-0047, ADR-0049, `conventions.md` |
| S-0141 | 160 KB | 72 KB | `flaiover-dashboard.md`, which T-0504 names as a file to update |
| S-0146 | 207 KB | 120 KB | ADR-0049, `agent-context.md`, and `flai-cli.md` (86 KB), which T-0525 names as a file to update |
| S-0147 | 189 KB | 102 KB | ADR-0049 and `flai-cli.md` |

A task that writes out a path to say "update this file" loads the whole file, and the two long design files are 72 and 86 KB. I propose a follow-up story under E-0010. A path written out in plain text (not a markdown link, not an ADR ID) would be briefed, not loaded whole. A link, a `#fragment` link, and an ADR ID would still load. That refines ADR-0049's "named loads whole", so it needs your word and probably an ADR. Say if you want it and I will write it as a backlog story; S-0146 goes to review as C without it.

### 2026-09-29T04:10:22Z alex
write the updated ADR to refine 0049 and then the backlog story
