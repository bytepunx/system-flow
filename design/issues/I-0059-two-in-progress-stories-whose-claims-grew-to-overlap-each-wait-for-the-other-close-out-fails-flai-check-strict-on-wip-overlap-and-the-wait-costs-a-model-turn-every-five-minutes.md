---
id: I-0059
title: "Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes"
class: defect
status: closed
count: 4
cost: 21m
first_reported: 2026-10-03T02:27:09Z
last_reported: 2026-10-03T19:31:58Z
updated: 2026-10-05T05:06:36Z
---

# I-0059 Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes

## Description

Two stories pulled within a minute of each other, S-0242 (01:22:36Z) and S-0198 (01:23:30Z), passed the pull-time hold (ADR-0046): their `touches` did not overlap. The tasks their agents then wrote widened each story's claim (a claim is the story's `touches` and its open tasks') into the other's: both ended up claiming `docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/system/flai-cli.md`. Nothing reported the new overlap when the tasks were created.

From then on neither story could finish. `scripts/close-out.sh` runs `flai check --strict`, which turns the advisory `wip.overlap` warning (ADR-0019: "a warning plus a visible badge keeps humans in charge") into a failure, and it checks the whole repository, so each story's close-out fails on the other story being in progress. S-0242 finished its work at 01:51Z and chose to wait for S-0198's acceptance to clear the overlap ("The only thing left before the final close-out is S-0198's acceptance"). S-0198 finished its tasks at 02:11Z; its verifier's close-out ran ten minutes and failed on the same three warnings against S-0242 at 02:22Z. Each now waits for the other, and only the operator can break it.

The waiting is not free. Both agents wait with the MCP `wait_for_events`, asking for 600 to 1800 seconds; the server caps every wait at five minutes (`maxWait` in `flai/internal/mcpserver/server.go` and `folder.go`), so each return is a model turn that re-reads the whole context: S-0198 made 22 waits (2776 s) and 227 top-level turns reading 33.8M cached tokens, 12.68 USD in an hour, most of it after its work was done; S-0242 made 44 waits (3170 s) for 3.14 USD. S-0198 also used `wait_for_events` as a sleep while its task sub-agents ran in the background (01:52Z to 02:08Z), so it noticed their completion only at the next five-minute tick. Its TH-0074 was answered at 01:46Z and acted on at 02:11Z.

## Instances

### 2026-10-03T02:27:09Z
2026-10-03, S-0198 and S-0242: both in progress for over an hour; S-0242 waiting since 01:51Z for S-0198's acceptance to clear the overlap, S-0198's verifier close-out failed at 02:22Z on the same three wip.overlap warnings against S-0242. 66 wait_for_events calls between them, each capped at 300 s, each return a full model turn; S-0198 cost 12.68 USD, 2776 s of it in waits.

### 2026-10-03T02:36:28Z
2026-10-03 02:29:45Z to 02:34:47Z, S-0242: after the designer allowed it to go first (TH-0077), the story's agent handed the close-out to a verifier and waited for it with wait_for_events (900 s asked, 300 s given). The verifier's events are in the log but the agent's feed showed nothing for five minutes, which read as stuck to the operator. It moved to review at 02:34:57Z and ended at 02:35:07Z, outcome worked.

### 2026-10-03T08:15:16Z
Story: S-0201.
S-0201's touches grew to flai/internal/workitem, flai/internal/itemedit, flai/cmd/edit.go, and docs/users while S-0200 was in progress on them; flai check --strict warns wip.overlap four times and stops S-0201's close-out.

### 2026-10-03T19:31:58Z
Story: S-0206.
S-0206's close-out stopped on flai check --strict's 14 wip.overlap warnings against S-0204 (in progress), after S-0206 widened its touches to the paths its tasks changed (flai/cmd files, hostapi, mcpserver, docs/users/flai.md, design/issues/summary.md); every test tier and the lint passed. S-0206 waited for S-0204 to leave in-progress and ran the close-out again.

## Remediation

Possible remediations, cheapest first:

1. **The close-out does not fail on another story's claim.** `wip.overlap` between two in-progress stories is the pull-time hold's business and advisory afterwards (ADR-0019, ADR-0046); `scripts/close-out.sh` (here and in the template) should not fail on it: pass an allow list to `flai check --strict` (`--allow wip.overlap`), or give `flai check` a `--story S-nnnn` scope that reports only what the story can clear. The delegation and work-management conventions should say that an overlap with another open story is not the story's to clear: note it in the narrative and go to review.
2. **Report a claim that grows into another story's.** When `flai task new`, `flai edit --touches`, or the MCP writes give an open story's task touches that overlap another in-progress story's claim, say so in the command's output and send both stories an inbox event (`overlapped`, as acceptance already does, S-0132), so the agents coordinate while the work is small instead of discovering it at close-out.
3. **Let a long poll be long.** `wait_for_events` and `wait_for_work` should honour the requested `timeout_seconds` up to what the transport allows (the agents asked for 900 and 1800 s and got 300), so an idle agent costs one turn per answer, not one per five minutes. Measure: at 150k cached tokens a turn, a five-minute tick costs about 0.05 USD; an hour of waiting about 0.60 USD per agent plus the turns it spends deciding to wait again.
4. **Do not wait for sub-agents with flai.** The delegation convention and `harness.Prompt` should tell the story's agent to wait for a background sub-agent through the harness's own completion notice, not by polling `wait_for_events`, which cannot see a sub-agent end.
5. **Detect the mutual wait.** `flai serve` or the designer's inbox could report two in-progress stories each blocked only on the other's claim (both agents idle in waits, both close-outs failing on `wip.overlap` against each other) and open a thread to the operator naming the pair and what clears it, rather than leaving two agents polling for an hour.

Related: S-0197 (stream sync and conflict reporting), the agent-waiting chart in E-0016 (S-0215) would have made this visible; ADR-0046 for the hold.
Closed 2026-10-05T05:06:36Z: S-0244. Remediation 1 was S-0249's close-out, scoped to the story, so a wip.overlap with another story no longer stops it. 2 by T-0867: a write that grows a story's claim into another in-progress story's says so and tells both stories as an overlapped change. 3 by T-0868: wait_for_events and wait_for_work hold for the timeout asked up to 30 minutes, with a progress notification every minute that keeps Claude Code's idle limit from dropping the call. 4 by T-0869: the prompt and the delegation convention say to wait for sub-agents through the harness's notice, not wait_for_events. 5 is not needed: once the close-out no longer stops on wip.overlap, two stories can no longer each wait for the other.
