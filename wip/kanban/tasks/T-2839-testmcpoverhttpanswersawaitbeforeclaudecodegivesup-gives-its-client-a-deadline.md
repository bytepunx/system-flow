---
id: T-2839
type: task
nature: remediation
title: TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp gives its client a deadline far above the server's bound
status: backlog
parent: S-0651
owner: alex
created: 2026-11-10T07:00:52Z
updated: 2026-11-10T07:00:52Z
transitions: []
stream: S-0651
tags: [cli, testing]
touches: [flai/cmd/mcp_http_test.go]
---
# T-2839 TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp gives its client a deadline far above the server's bound

## Work

The test scales the server's bound (`mcpserver.HTTPLongestWait`, 45 s) and the client's deadline (`claudeCodeHTTPTimeout`, 60 s) down by the same 20, leaving 2.25 s against 3 s. Under a load average of 160 to 220 on 24 cores the 0.75 s margin is gone, and `wait_for_work` reads as outliving the client (I-0309, instances of 2026-11-08T00:27:14Z and 2026-11-09T22:38:42Z).

- Keep the static check that `opt.MaxWait` is `mcpserver.HTTPLongestWait`, above zero and below `claudeCodeHTTPTimeout`: that is what keeps a real wait under Claude Code's timeout.
- Test the behaviour apart from the two real figures: set `opt.MaxWait` to a short bound, such as 200 ms, and give the client a deadline far above it and far below the 600 s the calls ask for, such as 30 s. Both waits must answer `timed_out`; a server that held them for the 600 s asked still fails at the client's deadline, as I-0228's did.
- Reproduce the failure before the change: run the old test under a CPU load, for example several `go test -race ./flai/...` runs at once, or with a delay put in the client's transport, and record in the task's Notes what failed and how.

It waits for no task: it changes only this file, which no other task of the story touches.

## Done when

- `go test -race -run '^TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp$' -count=20 ./flai/cmd/` passes under the load that made the old test fail, recorded in the task's Notes.
- With `opt.MaxWait` made longer than the client's deadline, the test fails, so it still catches I-0228.
- `flai test flai/cmd/mcp_http_test.go` passes.

## Notes
