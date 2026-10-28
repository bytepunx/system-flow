---
id: T-2215
type: task
nature: remediation
title: "flai serve starts nothing for a story whose agent runs elsewhere: look names the reason in waiting, and the automatic restart, the resume, and the send backs leave it alone, with tests that reproduce I-0214"
status: backlog
parent: S-0516
owner: alex
created: 2026-10-28T13:13:47Z
updated: 2026-10-28T13:14:11Z
transitions: []
stream: S-0516
tags: [cli, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2214]
---
# T-2215 flai serve starts nothing for a story whose agent runs elsewhere: look names the reason in waiting, and the automatic restart, the resume, and the send backs leave it alone, with tests that reproduce I-0214

## Work

In `flai/internal/serve/agents.go`, check `harness.Elsewhere(agent)` at each automatic start path beside today's "no harness and no command" check, before anything reaches `harness.For`: in `look`, where a ready story whose agent runs elsewhere joins a list of its own and is named in `waiting` with one reason per story, such as `S-0513's agent runs elsewhere (harness none): flai serve starts nothing for it`, logged as `agent not started` as the other reasons are, and counted as no failure; in `restartEnded`, where such a story in progress is left alone, with no automatic restart and no thread at the limit; and in `resume`, `startSentBack`, `startQueueSentBack`, and `startAcceptSentBack`, which skip it the same way. Write the tests that reproduce I-0214's two instances in `agents_test.go`: a story with `harness: none` that enters ready while the agent action is on, a command is set on the host, and the limit has room gets no agent, and `AgentState.Waiting` names it with the reason; and a story in progress with `harness: none` whose newest run ended early is not restarted, however `agent.auto_restarts` stands, and gets no limit thread. Show each test fails on the main branch's code as `git.md` says. Waits for T-2214, which gives `harness.Elsewhere`.

## Done when

- A ready story with `harness: none` is never started by `look`, whatever the host's command, and `waiting` names it with the reason, which changes the `agent not started` log line once
- `restartEnded`, `resume`, `startSentBack`, `startQueueSentBack`, and `startAcceptSentBack` each skip a story whose agent runs elsewhere, and a test covers the restart and at least one of the send backs
- The two tests that reproduce I-0214's instances fail without the change and pass with it, and the serve package's tests pass

## Notes

Layer 3, beside T-2216's refusals, since the two touch no path in common. The strategic agents (`plan.go`, `analyze.go`, `curate.go`, `orchestrate.go`) are not changed: a strategic agent with `harness: none` fails its start at `harness.For` with the refusal T-2214 words, which is enough for a value no role is expected to carry.
