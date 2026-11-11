---
id: D-0596
title: This host's ephemeral port range is 40525 to 44620, only 4096 ports, which every port 0 listener and outgoing connection of concurrent close-outs draws from, so a test or server that frees a port and binds it again by number picks outside that range
kind: decision
sources: [S-0714, T-3114, I-0342]
topics: [testing, go, code]
paths: [flai/internal/components/devserver.go, design/system/component-library.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0596 This host's ephemeral port range is 40525 to 44620, only 4096 ports, which every port 0 listener and outgoing connection of concurrent close-outs draws from, so a test or server that frees a port and binds it again by number picks outside that range

S-0714 read the host's range while remediating I-0342: "This host's ephemeral range is 40525-44620 (/proc/sys/net/ipv4/ip_local_port_range), only 4096 ports, which every httptest server and outgoing connection of concurrent close-outs draws from." That made the planner's proposed cause for I-0342 (D-0572) plausible: a Go listener given Vite's port once Vite stopped listening (wip/archive/agents/S-0714.md, Decisions).

So T-3114 made the components dev server pick its port "from the ports outside the ephemeral range, from 10000 up to 65535, so a host whose range starts low still has a band above it. Below 10000 is left out to keep off well-known and common service ports." `design/system/component-library.md` § "The dev server" describes the pick in general terms; this entry records how small the range is on this host.

A narrow range means a freed port in it is handed out again soon under the load of parallel close-outs. A test that frees a port and binds it again by number, as D-0319's fake Vite did, picks it outside the range.

Sources: S-0714 (T-3114), I-0342.
