---
id: D-0230
title: An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look
kind: cause
sources: [I-0095]
topics: [cli, flai, go, code, all]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md, design/issues/summary.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/system/workflow.md#branches-and-collisions-adr-0019"
asked_again: []
cited: []
---

# D-0230 An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look

I-0095 was closed on 2026-10-08.

The cause:

```text
An agent that ends with a question open on a thread is recorded as `asked`, and flai serve starts it again when the thread is answered. flai serve decides this when it notices the agent has ended, at a look, up to a minute later. If the operator answers in between, the look finds an ended agent and no open question. It records the run as `failed`, and nothing starts the agent again: the story stays in progress with no agent, holding every story that waits for it.

Two things made the window wide on 2026-10-06. flai serve had been restarted since it started the agent, so it was not the agent's parent, could not see it exit, and learned of the end only at its next look. And the operator was at the dashboard and answered within thirty seconds of the agent ending.
```

The fix:

```text
S-0317: flai serve's judgeRun now records a run asked when its agent wrote on a thread of its story after the run started and that question was answered before the end was judged (askedAnswered in flai/internal/serve/agents.go), so the next look starts the agent again at once in its session instead of recording it failed. TestAnAnswerBeforeTheEndIsJudgedStartsTheAgentAgain reproduces I-0095 and fails without the change.
```

Source: I-0095.
