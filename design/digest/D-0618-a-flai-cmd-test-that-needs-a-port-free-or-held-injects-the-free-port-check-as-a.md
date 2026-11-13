---
id: D-0618
title: A flai/cmd test that needs a port free or held injects the free-port check as a value, and one that must use real ports retries on new ones, rather than count on a port it closed or never held staying free (I-0372)
kind: cause
sources: [I-0372, S-0779, T-3330, T-3332, T-3337, T-3338, T-3339]
topics: [testing, flai, go, code]
paths: [flai/cmd/serve_mcp.go, flai/cmd/host.go, flai/cmd/serve_mcp_test.go, flai/cmd/mcp_http_test.go, flai/cmd/host_integration_test.go, flai/cmd/serve_pprof_test.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0618 A flai/cmd test that needs a port free or held injects the free-port check as a value, and one that must use real ports retries on new ones, rather than count on a port it closed or never held staying free (I-0372)

I-0372 was closed on 2026-11-13 by S-0779. Its cause: "flai's MCP address tests fail now and then when another test takes a port they closed or never held". On this host every port 0 listener of concurrent close-outs draws from a narrow ephemeral range (D-0596), so a freed port is soon handed out again.

S-0779's fix, from I-0372 and the story's Decisions (`wip/archive/agents/S-0779.md`):

- Where the code under test only asks whether a port is free, make the check a value the test fakes. T-3330: "the free-port check is a value (`free func(addr string) bool`) that `mcpAddrsFor` and `freeOnAll` take and `hostLauncher` carries, `freePort` when nil. A test then says which addresses are held, with no real port raced for."
- Where the test must listen for real, retry on a new port when a try visibly loses one. T-3332 keeps `listenMCP` on real addresses and "retries on a new port, at most 20 times, when 127.0.0.2:<port> is taken before or after the call. A listener `listenMCP` failed to close would fail every try, so the test still catches that bug." T-3338: the host integration test "hands ports to other processes, so it cannot hold them: it tries again on new ports, at most five tries".
- Check that a listener closed by what still accepts, not by re-listening. T-3337: the every-address test "counts an address gone when it can be listened on or refuses a connection, retrying for up to 3 s while it accepts, so a listener left open still fails it."
- Best of all, let the server pick its port. T-3339: "servePprof listens on 127.0.0.1:0 and logs its address", so the profiles test reads it from the log and needs no retry.

The fifty-count `MCP|TheHost` race run then passed beside the whole package looping under race: 2050 passes in 779 s (T-3333, run with `go test` by hand as D-0492 says).

Left for later and recorded in I-0372: `dashboard_test.go`'s `closedURL` relies on a closed port staying closed, not free, and `internal/components`' devserver tests are outside `flai/cmd` (D-0319, D-0596 cover the dev server's own port pick).

Sources: I-0372; S-0779 (`wip/archive/agents/S-0779.md`, Decisions), T-3330, T-3332, T-3337, T-3338, T-3339.
