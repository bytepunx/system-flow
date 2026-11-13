---
id: I-0372
title: flai's MCP address tests fail now and then when another test takes a port they closed or never held
class: defect
status: closed
count: 4
cost: 12m
first_reported: 2026-11-13T20:39:59Z
last_reported: 2026-11-13T21:51:02Z
updated: 2026-11-13T22:05:35Z
---

# I-0372 flai's MCP address tests fail now and then when another test takes a port they closed or never held

## Description
flai's MCP address tests fail now and then when another test takes a port they closed or never held

## Instances

### 2026-11-13T20:39:59Z
Story: S-0779.
Run 43455794271 on main (flai 1.93.1), test (ubuntu-latest, -race): TestMCPListenFailureClosesWhatItOpened failed at mcp_http_test.go:805 with `listen tcp 127.0.0.2:42545: bind: address already in use`. The test holds 127.0.0.1:<port> and counts on 127.0.0.2:<port> being free before listenMCP runs and again after it closes what it opened; nothing ensures either, and another parallel test's connection to a 127.0.0.2 server can take that port as its source. flai behaved as it should; the test's assumption broke. Cost: reading the failing job's log and running it again.

### 2026-11-13T20:40:04Z
Story: S-0779.
Run 43495183748 on main (commit 9b63bacf3, 2026-11-13), Integration tests step: TestTheHostMovesAProjectOffAPortAnotherProcessHolds failed at serve_mcp_test.go:516 with `the first start, on the remembered port: "127.0.0.1:4243" <nil>`. unusedPort asks the kernel for a port, closes it, and the test writes it as the remembered address; between the close and mcpAddrsFor's freeOnAll probe another parallel test took it, so flai rightly moved the project to the default range's 4243. The same test passed minutes before in the Behavior tests step on the same commit. Six tests use unusedPort this way.

### 2026-11-13T21:34:09Z
Story: S-0779.
S-0779's load run (T-3333), on the host: in one of ten whole-package `go test -race ./cmd/` runs beside the fifty-count `MCP|TheHost` run, TestMCPOverHTTPListensOnEveryAddress failed at mcp_http_test.go:776 with `a stop closes every listener, but 127.0.0.1:40842 is still held: listen tcp 127.0.0.1:40842: bind: address already in use`. After the server stops the test listened again on each port it had closed, counting on nobody taking it meanwhile. Fixed in the same story by T-3337: the test now asks whether anything still accepts at the address.

### 2026-11-13T21:51:02Z
Story: S-0779.
S-0779's second load run (T-3333), on the host: in the fifty-count `MCP|TheHost` run, TestHostRunsServeAndTheMCPServersAndTakesThemWithIt failed at host_integration_test.go:163 with `flai host already runs on this machine for /tmp/TestHostUpgradeToAVersion.../cfg.json`: it had closed a kernel port and given it to flai host as its address, and a host the whole-package run beside it started took it first. It does the same with the project's MCP address. TestProfilesAreOfferedOnLoopbackOnly hands servePprof a closed port the same way. Fixed in the same story by T-3338 (the scenario tries again on new ports when it loses one) and T-3339 (servePprof picks the port and the test reads it from the log). Not fixed here: closedURL in dashboard_test.go counts on a closed port staying closed, and internal/components' devserver tests take a freePort() and close it; neither is a test of flai/cmd's MCP addresses.

## Remediation
Closed 2026-11-13T22:05:35Z: S-0779: flai/cmd's tests now decide for themselves which ports are free and held. mcpAddrsFor and hostLauncher take the free-port check as a value, and the eight address tests give it a fake that holds what they say (T-3330, T-3331). The listen-failure test retries on a new port when 127.0.0.2 is lost (T-3332). The every-address test checks closure by what still accepts (T-3337). The host integration test retries on fresh ports (T-3338), and the profiles test lets servePprof pick its port (T-3339). The fifty-count `MCP|TheHost` race run then passed beside the whole package looping under race: 2050 passes in 779 s (T-3333).
