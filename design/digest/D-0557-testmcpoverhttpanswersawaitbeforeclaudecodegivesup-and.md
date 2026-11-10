---
id: D-0557
title: "TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp and TestBacksOffWhenTheDashboardSaysTheProjectIsHeld failed under load on margins under the bound they waited on; each now allows about four times it (I-0309)"
kind: cause
sources: [I-0309, S-0651, T-2839, T-2842]
topics: [testing, flai, go, code]
paths: [flai/cmd/mcp_http_test.go, flai/internal/channel/channel_test.go, design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0557 TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp and TestBacksOffWhenTheDashboardSaysTheProjectIsHeld failed under load on margins under the bound they waited on; each now allows about four times it (I-0309)

I-0309 recorded `flai/cmd`'s host and MCP-over-HTTP timing tests failing close-outs at host loads of about 100 to 200 on 24 cores, in code the failing stories did not touch. S-0651 closed it on 2026-11-10. The two import tests had been fixed already by S-0606's T-2709, which stamps the fake host state an hour ahead of `host.StaleAfter`.

The two that were left, from the issue's close:

- `TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp` (`flai/cmd/mcp_http_test.go`) "gives its client a 12 s deadline against the server's 2.25 s scaled bound, in place of 3 s, so only a wait held for the 600 s asked fails (T-2839)". S-0651's Decisions add: "A fatal check keeps the deadline at least four times the bound. The static check that HTTPLongestWait is under Claude Code's 60 s stays."
- `TestBacksOffWhenTheDashboardSaysTheProjectIsHeld` (`flai/internal/channel/channel_test.go`) "uses a 4 s held back-off, at least 2 s after jitter, and allows an ordinary close 1 s, in place of 300 ms against 250 ms (T-2842)". That is "about four times the 252 ms seen under load", and the test takes about 3 s longer.

Pinned to one CPU beside 40 busy loops, main's channel test failed 5 of 10 runs and its MCP test 1 of 5, and S-0651's versions passed (D-0492 says how).

So when a timing test fails under load because its allowance sits close to the bound it measures, the fix keeps the test's meaning and widens the margin to about four times the bound or the time seen under load, with a check that keeps that ratio if the bound moves. D-0551 and the I-0329 fixes are the same kind of cause in the dev-server tests.

Sources: I-0309; S-0651 (`wip/archive/agents/S-0651.md`, Decisions); T-2839; T-2842.
