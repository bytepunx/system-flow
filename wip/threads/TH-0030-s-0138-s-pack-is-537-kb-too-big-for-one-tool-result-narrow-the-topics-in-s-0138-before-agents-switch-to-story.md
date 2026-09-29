---
id: TH-0030
title: S-0138's pack is 537 KB, too big for one tool result. Narrow the topics in S-0138 before agents switch to --story?
anchor:
  path: wip/kanban/stories/S-0138-the-mcp-prime-tool-returns-a-story-s-context-pack-built-by-the-same-code-as-flai-prime-story.md
  item: S-0138
status: resolved
participants: [agent-S-0138, alex]
created: 2026-09-29T01:12:25Z
updated: 2026-09-29T03:09:48Z
---

# TH-0030 S-0138's pack is 537 KB, too big for one tool result. Narrow the topics in S-0138 before agents switch to --story?

On wip/kanban/stories/S-0138-the-mcp-prime-tool-returns-a-story-s-context-pack-built-by-the-same-code-as-flai-prime-story.md.

## Entries

### 2026-09-29T01:12:25Z agent-S-0138
TH-0029 is still open, and S-0138 is the story that switches agents to the pack. Measured on this repository today:

| | Size |
|-|------|
| `flai prime --cat` | 57 KB |
| `flai prime --story S-0138` | 537 KB (about 134k tokens) |
| conventions in it | 48 KB |
| design by topics (24 sections) | 267 KB; `flai-cli.md` alone is 83 KB |
| ADRs by links (33) | 134 KB |
| ranked (10) | 60 KB; `flaiover-dashboard.md` alone is 53 KB |

This is more than an agent can load in one call. Claude Code keeps about 30 KB of a Bash result and puts the rest in a file to read in chunks. It caps an MCP tool result at 25k tokens by default, about 100 KB. So the new MCP `prime` tool cannot return this pack whole, and an agent told to "prime with flai prime --story" would take in about 134k tokens before starting work.

Options:
- **A (recommended).** Land S-0138 as written: the prompt, CLAUDE.md, template, and the `prime` tool. Before you accept it, I add a task to S-0138 that narrows the pack on this repository. It would put heading topics on `flai-cli.md` and the `[all]` design files, and load a ranked item only up to its section, not the whole document. The target is under 100 KB for a cli story. The switch then happens with a pack an agent can actually read. This goes beyond the story's current touches (design/system topics, and one ranking change in flai/internal/context).
- **B.** Land S-0138 as written and narrow in a separate story. Until that story lands, agents flai starts here would prime with a 537 KB pack.
- **C.** Hold the prompt and CLAUDE.md on `--cat` and ship only the MCP tool and docs in S-0138. The switch comes in a later story once the pack is small.

Meanwhile I am going on with everything that does not depend on the answer: the tool, the prompt text, the template, and the docs.

### 2026-09-29T01:17:46Z agent-S-0138
Status: everything the criteria name is built and committed on story/S-0138. That is the MCP `prime` tool, the prompt, CLAUDE.md, the template (1.0.22), `session-start.md`, and the docs. `make test`, `make integration`, markdown lint, the installer test, and golangci-lint pass.

Two things stop the move to review:
1. This question (A, B, or C).
2. `flai check --strict` warns, so the smoke tier's repository check fails. Neither warning comes from this story's code:
   - S-0144 went in progress at 01:13 while it was held. It overlaps S-0138 on `docs/users/flai.md` and `design/system/flai-cli.md`. That clears when one of the two is accepted.
   - TH-0029 is open on the archived S-0137. If you answer here, may I resolve TH-0029 pointing at this thread?

### 2026-09-29T02:52:12Z alex
A research story was completed introducing ADR 49. Read that and then let's look at whether we need new stories or if this story can be amended to cover the adapted use.

### 2026-09-29T02:52:44Z agent-S-0138
I read ADR-0049. It needs two new stories, and S-0138 should shrink rather than grow.

ADR-0049 names three pieces of work:
1. The budget, the briefs, and ranked cuts in `flai/internal/context`.
2. `doc_search`, plus `heading` on `doc_get` and `flai doc show`.
3. The wording in the prompt, `CLAUDE.md`, and `session-start.md`: the pack is a brief, so read a briefed body with `doc_get` before changing what it describes.

The ADR says S-0138 can take (3). But (3) only makes sense once (1) and (2) have landed. S-0138 can't wait for them in progress: it touches `flai/internal/context`, `flai/internal/mcpserver`, and `flai/cmd`, so it would hold both new stories off the board.

**Recommended:**
- **S-0138 is amended to land now, without the switch.** It keeps `ctxpack.ForStory`, the MCP `prime` tool, and the docs. The prompt, `CLAUDE.md`, the template, and `session-start.md` go back to `flai prime --cat`. The server instructions stop telling agents to call `prime`: at 537 KB it cannot return a pack whole, so the tool stays listed but is not yet recommended. I would rewrite criteria 1 and 3 to say this. The template bump goes, since the template no longer changes.
- **New story (1):** the pack fits `--budget` (80 KB default), with briefs and ranked cuts. Done when `flai prime --story S-0138` measures under the budget.
- **New story (2):** `doc_search`, and `heading` on `doc_get` and `flai doc show`.
- **New story (3):** switch the prompt, `CLAUDE.md`, the template, and `session-start.md` to `--story`, with ADR-0049's wording. It is `after` (1) and (2).

ADR-0049 is still `proposed`, and the stories build what it decides. Should I set it to accepted when I create them?

Say yes and I will amend S-0138, create the three stories under E-0010 as backlog with goals and criteria, and move S-0138 to review.

### 2026-09-29T03:07:30Z alex
yes

### 2026-09-29T03:09:48Z agent-S-0138
Resolved: Designer said yes: ADR-0049 accepted; S-0138 lands without the switch (prime tool, ForStory, docs); S-0146 (budget and briefs), S-0147 (doc_search, heading), and S-0148 (the switch, after both) created under E-0010.
