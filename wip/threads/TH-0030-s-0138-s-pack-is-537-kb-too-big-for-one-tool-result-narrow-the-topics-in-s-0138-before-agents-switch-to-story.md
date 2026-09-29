---
id: TH-0030
title: S-0138's pack is 537 KB, too big for one tool result. Narrow the topics in S-0138 before agents switch to --story?
anchor:
  path: wip/kanban/stories/S-0138-agents-flai-starts-and-any-agent-with-a-story-prime-with-flai-prime-story.md
  item: S-0138
status: open
participants: [agent-S-0138]
created: 2026-09-29T01:12:25Z
updated: 2026-09-29T01:17:46Z
---

# TH-0030 S-0138's pack is 537 KB, too big for one tool result. Narrow the topics in S-0138 before agents switch to --story?

On wip/kanban/stories/S-0138-agents-flai-starts-and-any-agent-with-a-story-prime-with-flai-prime-story.md.

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
