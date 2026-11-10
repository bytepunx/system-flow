---
id: D-0532
title: "S-0603 criterion 3: agent.status is not under 200 ms with 13 projects at once on this host: Recommended: option 1, reword criterion 3, with one addition: the 200 ms target goes to an issue, because…"
kind: answer
sources: [TH-0911, S-0603, TH-0910, S-0601]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [flai/internal/serve/agents.go, flai/cmd/serve_actions.go, flai/internal/serve/agents_test.go, flai/internal/messages/messages.go, flai/internal/workitem/hold.go, flai/internal/workitem/share.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions_test.go, design/system/server-performance.md, design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md, design/issues/summary.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: [S-0636, T-2804, TH-0954, TH-0961, TH-1032]
---

# D-0532 S-0603 criterion 3: agent.status is not under 200 ms with 13 projects at once on this host: Recommended: option 1, reword criterion 3, with one addition: the 200 ms target goes to an issue, because…

The operator, alex, answered TH-0911 at 2026-11-09T07:21:22Z.

What was asked:

```text
**Recommended:** option 1. Reword criterion 3 to what the story showed, and leave the 200 ms target to S-0602 (cause 13).
```

The answer:

```text
**Recommended:** option 1, reword criterion 3, with one addition: the 200 ms target goes to an issue, because S-0602 does not measure it.

**Why option 1:**

- The design never expected cause 12's fix to reach 200 ms alone. server-performance.md § Where the time goes, cause 12, estimates "15 to 150 ms and the lock queue; with cause 13, agent.state's CPU from about 450 to under 60 ms". S-0603's own Notes say the same.
- The story did what cause 12 proposes: `agents.json` kept parsed, the agentsMu queue gone, the conversations read once a call. The measured p50 of a call, 13 at once, fell from 369–463 ms to 221–283 ms. This repository's call fell from 2,092–2,952 ms to 1,097–1,107 ms.
- What is left is cause 13, the thread and conversation reads, plus a host at load 95 to 101 on 24 cores.
- This is the same case as TH-0910 for S-0601, and the same kind of answer as D-0466: a criterion that cannot be met until another story lands moves out rather than holding a finished story.

**The addition:** S-0602's criterion 3 times `threads.list`, `inbox.designer`, and `stats.get`. It does not time `agent.status`, so nothing would measure the 200 ms once both stories are in. The S-0603 agent records that target as an issue, "agent.status with 13 projects at once is not under 200 ms", citing causes 12 and 13 and this thread, and names it in S-0603's Notes. Its title's "takes under 200 ms" becomes "takes about 250 ms, not 0.8 to 2.4 s", to match the reworded criterion.

**Then:** the agent ticks criterion 3 and moves S-0603 to review. I verify it and accept or enqueue it under accept_reviews.

**Why I recommend rather than answer:** rewording a finalized story's criterion changes its scope, which is your call.

**Alternatives:** option 2 holds a finished story behind S-0602 and a quieter host. Option 3 leaves a criterion unticked on a done story, which the definition of done does not allow.

Source: design/system/server-performance.md § Where the time goes
```

The operator confirmed the orchestrator's recommendation on 2026-11-09T07:21:22Z. S-0603's criterion 3 was reworded, its title became "takes about 250 ms, not 0.8 to 2.4 s", and the 200 ms target went to an issue, "agent.status with 13 projects at once is not under 200 ms" (I-0313, on S-0603's branch), since S-0602's criterion times `threads.list`, `inbox.designer`, and `stats.get` but not `agent.status`.

So a performance story whose absolute target is out of reach because another cause in `design/system/server-performance.md` is not yet fixed, or because the host is loaded, has its criterion reworded to the relative fall it measured against main on the same host, with the before and after in its cause section. The absolute target goes to an issue when no other story measures it, rather than holding the finished story in review. This is the same answer as D-0531 (S-0602) and the same kind as D-0466: a criterion that waits on another story moves out. TH-0910 answered the same for S-0601 (D-0537): its criterion 3, a warm whole `flai check --strict` under 2 s, became half of main's on the same host, and the 2 s target went to I-0312.

Sources: TH-0911, S-0603, TH-0910, S-0601.
