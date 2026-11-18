---
id: D-0565
title: "Opening a read form to the agents flai guard holds to its reads takes no new ADR: S-0428, S-0491, S-0612, and S-0680 each added reads without one"
kind: decision
sources: [S-0680, I-0318, S-0829, I-0355]
topics: [guard, flai, go, kluster]
paths: [flai/internal/guard/guard_tables.go, design/adrs, design/system/kluster.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0565 Opening a read form to the agents flai guard holds to its reads takes no new ADR: S-0428, S-0491, S-0612, and S-0680 each added reads without one

S-0680 (I-0318) opened the read forms of `flai serve` that ADR-0220 lists, such as `flai serve status`, to every agent `flai guard` holds to its reads. It wrote no ADR. Its narrative gives the reason: "No new ADR: earlier read additions (S-0428, S-0491, S-0612) took none." (`wip/archive/agents/S-0680.md`, § Decisions.)

S-0829 (I-0355) did the same for `flai kluster status` and `flai kluster deployments`, which `design/system/kluster.md` already called reads any agent may run. It added both to the guard's `cliReads` and wrote no ADR. I-0355's fix: "flai guard's cliReads, the table of flai commands a sub-agent may run, now lists kluster status and deployments, so every sub-agent and strategic role reads the cluster as kluster.md says."

So a story that adds a form the guard should pass as a read, one an ADR or the design already names as a read, adds it to the guard's tables in `flai/internal/guard/guard_tables.go` and to the design and user docs, without an ADR of its own.

The operator's answer on the scope of S-0680 is D-0562 (TH-0991).

Sources: S-0680; I-0318; S-0829; I-0355.
