---
id: ADR-0120
title: "Agents of two open stories message each other in conversations kept under wip/messages, apart from the operator's threads, and closed when either story leaves"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0020]
---

# ADR-0120 Agents of two open stories message each other in conversations kept under wip/messages, apart from the operator's threads, and closed when either story leaves

## Context

Stories run in parallel, and their agents change paths that meet (`design/system/agent-coordination.md`). The only channel between two agents today is a thread ([ADR-0020](0020-files-plus-mcp.md)), and a thread is the designer's: it is anchored on one document or item, every thread an agent writes awaits the operator in `inbox` and in the dashboard's designer inbox, it is mirrored into the story's `## Open questions`, and acceptance or `flai archive` resolves it ([ADR-0109](0109-archiving-an-item-resolves-the-threads-still-open-on-it.md)). Nothing in a thread names the other story's agent as the one to answer. Two agents that want to agree on who changes what before their branches conflict either ask the operator or say nothing. E-0018 gives them a channel of their own; this decision fixes how it is addressed, kept, and ended.

## Decision

A message is sent from one open story to another, and the messages between two stories make a conversation that is kept apart from threads.

1. **Addressed from story to story.** `flai message send <S-nnnn> "<text>" --from <S-nnnn>` starts a conversation between the sender's story (`--from`, else `FLAI_STORY`, else the story in an `agent-S-nnnn` `FLAI_AGENT`) and the addressee. Both must be stories in progress or in review, and not the same story; a send to or from a story in any other state is refused with its state. `--about` names the repository paths the conversation is about. A reply names the conversation and comes from one of its two stories; the other story is the one it awaits. Each entry records its author, the agent, and its story.
2. **Kept as files under `wip/messages`.** One markdown file per conversation, `wip/messages/MS-nnnn-<slug>.md`, written in the main checkout as threads are, numbered one more than the highest there. Its front matter holds `id`, `title` (the first line of the first message), `from`, `to`, `about`, `status` (`open` or `closed`), `participants`, `created`, and `updated`; its body holds the dated entries under `## Entries`, each headed `### <time> <author> <story>`. flai lints what it writes with the project's markdownlint configuration ([ADR-0061](0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md)), and `flai check --strict` validates the files: front matter, unique IDs, file names, that both stories exist, and the dated entries.
3. **Closed when either story leaves.** A conversation reads as closed when its status is `closed` or either of its stories is done, cancelled, or archived; a closed one takes no reply. `flai accept` and `flai archive` close every open conversation of each story they archive with an entry `Closed: <S-nnnn> was accepted` (or `was archived`), and acceptance commits the files with the rest. A cancelled story's conversations read as closed without an entry. A closed conversation is never reopened: a new message starts a new one.
4. **Apart from threads.** Messages have their own folder, ID prefix, package, and command. No message appears in `flai thread list`, among the threads awaiting the operator in `inbox`, in the dashboard's designer inbox, or in a narrative's `## Open questions`. The operator reads them with `flai message list --all` and `flai message show`. A question for the designer is still a thread.

`flai message` (`send`, `reply`, `list`, `show`, each with `--json`) is the first surface. The MCP tools, the agents' inbox, and the dashboard follow in later stories of E-0018.

## Consequences

- Two agents can coordinate without the operator seeing every exchange as work awaiting them, and without asking the operator to relay.
- The operator's inbox keeps its meaning: what awaits the operator is a thread.
- `wip` gains a folder that `flai check`, acceptance, and archive must know; the dashboard and the MCP server read it in their own stories.
- A conversation in flight when its story is accepted closes with the reason and keeps every entry.

## Alternatives considered

- **Threads with a story-to-story anchor.** It would reuse the thread package and files, but every reader of `wip/threads`, the operator's inbox, the designer inbox, and the narrative mirror, would need to tell the two kinds apart, and ADR-0109's resolution on one anchored item does not fit a conversation between two.
- **Messages to agents rather than to stories.** An agent's name changes when flai serve restarts it or the operator works the story by hand; the story is what lasts and what both claims name.
- **Messages kept in the narratives.** A narrative belongs to one story and is rewritten by `flai stream state`; a conversation belongs to two.
