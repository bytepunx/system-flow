---
id: D-0635
title: workitem/transition.go's exported AcceptSource, CancelledState, and FirstOfKind duplicate the unexported acceptSource, cancelledState, and firstOfKind S-0739 adds to workitem/flow.go and S-0746's helpers in mcpserver, so once both stories are accepted one set can replace the other
kind: decision
sources: [S-0744, S-0739, S-0746]
topics: [flai, go, code]
paths: [flai/internal/workitem/transition.go, flai/internal/workitem/flow.go, flai/internal/mcpserver]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0635 workitem/transition.go's exported AcceptSource, CancelledState, and FirstOfKind duplicate the unexported acceptSource, cancelledState, and firstOfKind S-0739 adds to workitem/flow.go and S-0746's helpers in mcpserver, so once both stories are accepted one set can replace the other

S-0744 added exported helpers to `flai/internal/workitem/transition.go` that find a workflow's acceptance source and target, its cancelled state, and the first state of a kind: `Accepts`, `AcceptSource`, `AcceptTarget`, `CancelledState`, `FirstOfKind`, `IsKind`, and `PermitMove`. Its narrative's Decisions (`wip/archive/agents/S-0744.md`) record that they repeat what parallel stories wrote:

> The exported helpers in workitem/transition.go overlap S-0739's unexported ones in workitem/flow.go (acceptSource, cancelledState, firstOfKind) and S-0746's in mcpserver; they are distinct identifiers, so both merge, and one can replace the other once both are accepted.

S-0746 is accepted. S-0739 was in review on 2026-11-15, so `flow.go`'s unexported copies reach main with it. A later E-0034 story that needs one of these lookups calls `transition.go`'s exported helper, or `repo.Flow()` and the workflow package's lookups as D-0630 says, and adds no third copy. A story that tidies `workitem` can fold `flow.go`'s copies into `transition.go`'s once S-0739 is accepted.

Source: S-0744's narrative, `wip/archive/agents/S-0744.md`, § Decisions.
