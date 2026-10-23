---
id: D-0161
title: "Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes"
kind: cause
sources: [I-0059]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/workitem/hold.go, flai/internal/itemedit, flai/internal/itemnew, flai/cmd/items.go, flai/cmd/edit.go, flai/cmd/touches.go, flai/internal/mcpserver, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/harness, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/agent-narrative.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md, design/issues/summary.md, flai/cmd/touches_overlap_test.go, template/template.yaml, docs/operators/index.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0161 Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes

I-0059 was closed on 2026-10-05.

The cause:

```text
Two stories pulled within a minute of each other, S-0242 (01:22:36Z) and S-0198 (01:23:30Z), passed the pull-time hold (ADR-0046): their `touches` did not overlap. The tasks their agents then wrote widened each story's claim (a claim is the story's `touches` and its open tasks') into the other's: both ended up claiming `docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/system/flai-cli.md`. Nothing reported the new overlap when the tasks were created.

From then on neither story could finish. `scripts/close-out.sh` runs `flai check --strict`, which turns the advisory `wip.overlap` warning (ADR-0019: "a warning plus a visible badge keeps humans in charge") into a failure, and it checks the whole repository, so each story's close-out fails on the other story being in progress. S-0242 finished its work at 01:51Z and chose to wait for S-0198's acceptance to clear the overlap ("The only thing left before the final close-out is S-0198's acceptance"). S-0198 finished its tasks at 02:11Z; its verifier's close-out ran ten minutes and failed on the same three warnings against S-0242 at 02:22Z. Each now waits for the other, and only the operator can break it.

The waiting is not free. Both agents wait with the MCP `wait_for_events`, asking for 600 to 1800 seconds; the server caps every wait at five minutes (`maxWait` in `flai/internal/mcpserver/server.go` and `folder.go`), so each return is a model turn that re-reads the whole context: S-0198 made 22 waits (2776 s) and 227 top-level turns reading 33.8M cached tokens, 12.68 USD in an hour, most of it after its work was done; S-0242 made 44 waits (3170 s) for 3.14 USD. S-0198 also used `wait_for_events` as a sleep while its task sub-agents ran in the background (01:52Z to 02:08Z), so it noticed their completion only at the next five-minute tick. Its TH-0074 was answered at 01:46Z and acted on at 02:11Z.
```

The fix:

```text
S-0244. Remediation 1 was S-0249's close-out, scoped to the story, so a wip.overlap with another story no longer stops it. 2 by T-0867: a write that grows a story's claim into another in-progress story's says so and tells both stories as an overlapped change. 3 by T-0868: wait_for_events and wait_for_work hold for the timeout asked up to 30 minutes, with a progress notification every minute that keeps Claude Code's idle limit from dropping the call. 4 by T-0869: the prompt and the delegation convention say to wait for sub-agents through the harness's notice, not wait_for_events. 5 is not needed: once the close-out no longer stops on wip.overlap, two stories can no longer each wait for the other.
```

Source: I-0059.
