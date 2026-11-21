---
id: D-0664
title: A send back by hand opens its thread by flai with storygit.SentBackPrefix, so flai serve's startAcceptSentBack and resolveAcceptSentBack start and resolve it as an acceptance's send back, and tell the two apart by the title's by-hand suffix with sendback.ByHandMover
kind: decision
sources: [S-0823, I-0397]
topics: [flai, cli, go, code]
paths: [flai/internal/sendback/sendback.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_start.go]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: [S-0810]
---

# D-0664 A send back by hand opens its thread by flai with storygit.SentBackPrefix, so flai serve's startAcceptSentBack and resolveAcceptSentBack start and resolve it as an acceptance's send back, and tell the two apart by the title's by-hand suffix with sendback.ByHandMover

S-0823 fixed I-0397: a story the operator moved from review to in-progress by hand got no agent and no thread. Rather than add a loop of its own to `flai serve`, it made the by-hand thread look like an acceptance's send back. From S-0823's Decisions (`wip/archive/agents/S-0823.md`):

> The thread a send back by hand opens keeps storygit.SentBackPrefix and is by flai, so startAcceptSentBack and resolveAcceptSentBack, which already start and resolve such threads, take it up; T-3543 tells it apart by its title to word the prompt.

And, from the same Decisions:

> T-3543: serve tells a send back by hand from an acceptance's by its thread's title (sendback.ByHandMover) inside startAcceptSentBack, rather than a new loop, since that loop already starts and resolves every open Acceptance sent back thread by flai. A by-hand move that writes on an acceptance's thread still open keeps that thread's acceptance wording.

So a later story that adds another kind of send back can reuse that loop the same way: open the thread by `storygit.SentBackAuthor` with `storygit.SentBackPrefix`, and mark its kind in the title, as `sendback.ByHandSuffix` does in `flai/internal/sendback/sendback.go`. A story that changes the prefix or the author changes what the loop picks up. What the operator sees is in `design/system/flai-cli.md` (`flai move`) and `design/system/workflow.md` § Transitions and who makes them.
