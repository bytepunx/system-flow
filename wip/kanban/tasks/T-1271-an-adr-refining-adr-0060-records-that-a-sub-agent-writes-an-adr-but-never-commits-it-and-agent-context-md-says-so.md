---
id: T-1271
type: task
nature: remediation
title: An ADR refining ADR-0060 records that a sub-agent writes an ADR but never commits it, and agent-context.md says so
status: in-progress
parent: S-0287
owner: alex
created: 2026-10-07T23:19:16Z
updated: 2026-10-08T07:14:09Z
transitions:
  - to: ready
    at: 2026-10-08T07:14:08Z
    by: agent-S-0287
  - to: in-progress
    at: 2026-10-08T07:14:09Z
    by: agent-S-0287
stream: S-0287
tags: [adr, design, guard]
touches: [design/system/agent-context.md]
---
# T-1271 An ADR refining ADR-0060 records that a sub-agent writes an ADR but never commits it, and agent-context.md says so

## Work

The story's agent records the decision with `flai adr new "<the decision, as a sentence>" --refines ADR-0060 --status accepted --body-stdin --commit` in the story's worktree, which numbers it, writes its row in `design/adrs/README.md`, commits both, and adds them to the story's touches. Its decision: `flai guard` lets a sub-agent run `flai adr new`, `topics`, and `accept`, and call `adr_new`, without a commit, because an ADR is a document ADR-0060 does not guard; `--commit`, `--autocommit`, and `adr_new`'s `commit` stay the story's agent's. Its context is I-0062's three instances (S-0207, S-0249, S-0220); its alternatives name refusing all of them. Set its topics with `flai adr topics`.

In `design/system/agent-context.md` § Sub-agents, in the paragraph on ADR-0060, say what a sub-agent may now do with an ADR, linking the new ADR.

It waits for no task: it records the rule the story's plan chose, which T-1270 builds at the same time on other paths.

## Done when

- The ADR is accepted, refines ADR-0060, has topics, and is committed on the story's branch with its index row.
- `design/system/agent-context.md` § Sub-agents names the rule and links the ADR, and its `updated` is bumped.

## Notes
