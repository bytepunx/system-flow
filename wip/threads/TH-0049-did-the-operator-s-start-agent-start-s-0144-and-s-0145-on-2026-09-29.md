---
id: TH-0049
title: Did the operator's Start agent start S-0144 and S-0145 on 2026-09-29?
anchor:
  path: wip/kanban/stories/S-0182-find-and-fix-how-flai-serve-1-23-0-started-an-agent-for-a-story-the-board-showed-held.md
  item: S-0182
status: resolved
participants: [agent-S-0182, alex]
created: 2026-10-01T09:15:36Z
updated: 2026-10-01T09:52:59Z
---

# TH-0049 Did the operator's Start agent start S-0144 and S-0145 on 2026-09-29?

On wip/kanban/stories/S-0182-find-and-fix-how-flai-serve-1-23-0-started-an-agent-for-a-story-the-board-showed-held.md.

## Entries

### 2026-10-01T09:15:36Z agent-S-0182
The 2026-09-29 instances ran on another host (pid 2077672, Linux). This host's `~/.flai/serve/journal.jsonl` starts 2026-09-26 and has no entry from 2026-09-29, so I cannot read that host's evidence.

What the code and the stories' transitions show:

- In flai 1.23.0 the launcher (`look`) checks `holds.Of` before every start, against the same `Repo.Holds` the board uses, and file-versus-directory claims overlap correctly (`PathsOverlap`). Every `agent not started ... held` line in this host's `serve.log` today names the right holder.
- `resume()` cannot explain either: S-0144 and S-0145 were new stories with no earlier run.
- The story page's Start agent (S-0115, S-0129) overrides a hold on purpose. It starts the agent with the same prompt as the launcher ("started by flai serve ... because it entered ready"), the same `FLAI_STARTED_BY=flai-serve`, and a `serve.agent` journal entry by `flai serve`. Its "started for a held story" warning goes to the `flai serve agent start` subprocess's log, not to the agent. So an agent started this way would report exactly what I-0050 records.
- S-0145 entered ready at 01:57:54Z and its agent started at 02:10:59Z, 3.5 minutes after S-0144's acceptance (02:07:25Z) and with S-0138 still in progress. The launcher looks on every kanban change and every minute, so it would have started S-0145 at 02:07:25Z if it was going to start it at all. A press of Start agent fits that gap. The launcher doesn't.

Questions:

1. Does that host's `~/.flai/serve/journal.jsonl` (or `serve.log`) still exist? An `agent.start` entry by you at about 01:10–01:14Z for S-0144 and at 02:10:59Z for S-0145 would settle it.
2. If not, do you recall pressing Start agent on either story?

My recommendation, whatever the answer: treat the operator's start as the cause unless the journal says otherwise. The fix would be:
- (a) An agent started past a hold or a full limit is told so in its prompt, with the hold's reason, so it neither mistakes the start for the launcher's nor works around the hold unasked.
- (b) `resume()` restarts only a story that is still open (in progress or in review). A story sent back to ready goes through the hold and the limit like any other.
- Regression tests for the launcher with file-versus-directory claims.

I'm going ahead with (a), (b), and the tests now. If you have the journal, I'll check it against them.

### 2026-10-01T09:52:59Z alex
Resolved.
