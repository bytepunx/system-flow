---
id: T-2216
type: task
nature: remediation
title: flai serve agent start, restart, and commit refuse a story whose agent runs elsewhere, saying flai serve starts nothing for it
status: backlog
parent: S-0516
owner: alex
created: 2026-10-28T13:13:55Z
updated: 2026-10-28T13:14:16Z
transitions: []
stream: S-0516
tags: [cli, serve]
touches: [flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/internal/serve/commit.go, flai/internal/serve/commit_test.go]
after: [T-2214]
---
# T-2216 flai serve agent start, restart, and commit refuse a story whose agent runs elsewhere, saying flai serve starts nothing for it

## Work

In `flai/internal/serve/start.go` (`Start`), `restart.go` (`Restart`), and `commit.go`, refuse a story whose agent `harness.Elsewhere` reports, beside today's "names no harness, and no command is set on the host" refusal and before anything reaches `harness.For`: the refusal says the story's agent runs elsewhere (`harness: none`), that `flai serve` starts nothing for it, and what to do instead (`flai story start <story>` for a ready story; name a harness flai serve can start to have it started here). The dashboard's **Start agent**, **Retry**, and **Have an agent commit them** buttons reach these through the `agent.start`, `agent.restart`, and `agent.commit` writes, so the same words reach the story's page and the journal. Add a test in each of `start_test.go`, `restart_test.go`, and `commit_test.go` beside the no-harness refusal's. Waits for T-2214, which gives `harness.Elsewhere`.

## Done when

- `flai serve agent start`, `restart`, and `commit` each refuse a story with `harness: none`, whether or not a command is set on the host, with a message that names `flai story start` and the harness, and the refusal is journalled as the no-harness one is
- A test in each of the three test files checks the refusal, and the serve package's tests pass

## Notes

Layer 3, beside T-2215: it touches no path T-2215 touches, so the two can run together.
