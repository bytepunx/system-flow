---
id: T-0869
type: task
nature: remediation
title: The story agent's prompt and the conventions say to wait for sub-agents through the harness's notice, and what to do on a grown-claim overlapped change
status: in-progress
parent: S-0244
owner: alex
created: 2026-10-05T03:16:11Z
updated: 2026-10-05T05:03:51Z
transitions:
  - to: ready
    at: 2026-10-05T05:01:49Z
    by: agent-S-0244
  - to: in-progress
    at: 2026-10-05T05:01:49Z
    by: agent-S-0244
stream: S-0244
tags: [flai, conventions, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, template/template.yaml]
after: [T-0867, T-0868]
---
# T-0869 The story agent's prompt and the conventions say to wait for sub-agents through the harness's notice, and what to do on a grown-claim overlapped change

## Work

This covers I-0059's remediation 4, and the agent-facing half of remediations 2 and 3. While its task sub-agents ran in the background, S-0198 polled `wait_for_events`, which cannot see a sub-agent end. It noticed their completion only at the next five-minute tick.

- In `harness.Prompt`'s delegation text (`flai/internal/harness/harness.go`, around line 344, "run a layer's tasks at once ... and then wait for all of them"), tell the story's agent to wait for background sub-agents through the harness's own completion notice. It must not poll `wait_for_events`, which reports work items and threads, not sub-agents. Say the same in `delegation.md` under Tasks worked by sub-agents, here and in `template/root`.
- In `work-management.md`, here and in `template/root`, extend the `overlapped` rule to the change T-0867 adds. When `cause` is a story in progress, not an accepted one, the two claims grew to overlap. Coordinate with that story's agent on a thread while the work is small, and narrow your touches if you can. Otherwise note the overlap in the narrative and go on: it does not stop the close-out (S-0249).
- If T-0868 found that the claude-code MCP client cuts a tool call shorter than the new cap, set the client's limit, such as `MCP_TOOL_TIMEOUT`, in the environment `flai serve` gives a claude-code agent, here in the harness.
- Add a `template/CHANGELOG.md` entry for the convention changes.

It waits for T-0867, because it describes the change T-0867 adds, and for T-0868, because the limit it may set depends on the cap T-0868 settles.

## Done when

- [ ] `harness_test.go` asserts the prompt tells the agent to wait for a background sub-agent through the harness's completion notice, not `wait_for_events`
- [ ] `delegation.md` and `work-management.md` say the same in `design/conventions` and `template/root/design/conventions`, below or in the baseline as each file's rules allow
- [ ] When T-0868 called for it, a test shows the claude-code agent's environment carries the MCP tool-call limit
- [ ] The markdown lint and `scripts/flai-test.sh` pass

## Notes
