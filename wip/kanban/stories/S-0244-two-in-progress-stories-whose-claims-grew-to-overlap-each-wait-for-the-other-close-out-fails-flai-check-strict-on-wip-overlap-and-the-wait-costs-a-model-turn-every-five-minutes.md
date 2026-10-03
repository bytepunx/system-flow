---
id: S-0244
type: story
nature: remediation
title: "Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes"
status: backlog
owner: alex
created: 2026-10-03T03:31:01Z
updated: 2026-10-03T03:31:01Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0244 Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes

## Goal

This story remediates [I-0059](../../../design/issues/I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md), "Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes". The issue recommends this solution:

Possible remediations, cheapest first:

1. **The close-out does not fail on another story's claim.** `wip.overlap` between two in-progress stories is the pull-time hold's business and advisory afterwards (ADR-0019, ADR-0046); `scripts/close-out.sh` (here and in the template) should not fail on it: pass an allow list to `flai check --strict` (`--allow wip.overlap`), or give `flai check` a `--story S-nnnn` scope that reports only what the story can clear. The delegation and work-management conventions should say that an overlap with another open story is not the story's to clear: note it in the narrative and go to review.
2. **Report a claim that grows into another story's.** When `flai task new`, `flai edit --touches`, or the MCP writes give an open story's task touches that overlap another in-progress story's claim, say so in the command's output and send both stories an inbox event (`overlapped`, as acceptance already does, S-0132), so the agents coordinate while the work is small instead of discovering it at close-out.
3. **Let a long poll be long.** `wait_for_events` and `wait_for_work` should honour the requested `timeout_seconds` up to what the transport allows (the agents asked for 900 and 1800 s and got 300), so an idle agent costs one turn per answer, not one per five minutes. Measure: at 150k cached tokens a turn, a five-minute tick costs about 0.05 USD; an hour of waiting about 0.60 USD per agent plus the turns it spends deciding to wait again.
4. **Do not wait for sub-agents with flai.** The delegation convention and `harness.Prompt` should tell the story's agent to wait for a background sub-agent through the harness's own completion notice, not by polling `wait_for_events`, which cannot see a sub-agent end.
5. **Detect the mutual wait.** `flai serve` or the designer's inbox could report two in-progress stories each blocked only on the other's claim (both agents idle in waits, both close-outs failing on `wip.overlap` against each other) and open a thread to the operator naming the pair and what clears it, rather than leaving two agents polling for an hour.

Related: S-0197 (stream sync and conflict reporting), the agent-waiting chart in E-0016 (S-0215) would have made this visible; ADR-0046 for the hold.

## Acceptance criteria
- [ ] The cause I-0059 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0059 is closed with `flai issue close I-0059 --reason` saying what fixed it

## Tasks

## Notes
