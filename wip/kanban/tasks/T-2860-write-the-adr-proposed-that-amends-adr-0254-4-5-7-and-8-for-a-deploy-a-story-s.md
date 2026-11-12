---
id: T-2860
type: task
nature: feature
title: "Write the ADR, proposed, that amends ADR-0254 §4, §5, §7, and §8 for a deploy a story's agent asks for while its story is in progress"
status: in-progress
parent: S-0659
owner: alex
created: 2026-11-10T08:20:31Z
updated: 2026-11-12T16:03:58Z
transitions:
  - to: ready
    at: 2026-11-12T16:03:58Z
    by: agent-S-0659
  - to: in-progress
    at: 2026-11-12T16:03:58Z
    by: agent-S-0659
stream: S-0659
tags: [cli]
touches: [design/adrs/README.md]
---
# T-2860 Write the ADR, proposed, that amends ADR-0254 §4, §5, §7, and §8 for a deploy a story's agent asks for while its story is in progress

## Work

Record the decision with `flai adr new "<the decision, as a sentence>" --refines 0254 --body-stdin --commit` from the story's worktree, status proposed. `--commit` writes the next number, the file, and the index, and adds the new file to the story's touches, so no task names it now. Use `--supersedes 0254` instead only if the decision replaces §1 to §3 and §6 as well, which this plan does not expect.

Weigh first the hand-off `flai verify --job` makes to flai serve ([ADR-0229](../../../design/adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)): the agent asks with an MCP tool or a `--job` flag, gets a job id at once, and waits for an event. Read ADR-0254 §4 to §8 and `design/system/kluster.md` § What happens when, § A failed deploy holds the story, § The `kluster` host action, § flai serve (S-0579), and § The jobs before writing.

Start from the design as built now, not as ADR-0254 first wrote it. S-0694 built §5a: a failed deploy or checks run sends the story back to in-progress and keeps its namespace there, and `kluster.md`'s table has rows for the send-back. [ADR-0259](../../../design/adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md) makes `deploy` optional, and [ADR-0263](../../../design/adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md) adds components and a record of every step. Say what an agent's request gets on a project with no `deploy`, and how an agent's deploy in progress fits with a namespace a send-back kept.

The body answers each of the story's criteria 1 to 6 in its own numbered section:

1. How a story's agent asks for a deploy and a teardown of its own story's namespace only, that flai serve runs the job, and that the `kluster` host action gates it.
2. What the agent may run against the namespace: the project's checks, its own commands with §4's variables, and what `flai guard` still refuses, another story's namespace, the default one, `up`, `down`, and `promote` among them.
3. When a namespace is made, deployed again, and torn down across in-progress, review, a send back, acceptance, and cancellation, as a table that replaces ADR-0254 §5's and the send-back rows `kluster.md` has today.
4. How many story namespaces may live at once, where the project sets it (the manifest's `kluster:` block, by this plan's assumption), its default, and what a request over the limit gets.
5. How the agent waits for a job behind the queue and for how long, what it is told on a failure, and when a failure reaches the operator, against §5a's thread as S-0694 built it.
6. What is recorded of a deploy and a checks run before review, and where the move to review, the close-out, and an acceptance read it, against §5a's `kluster` blocker and the record's `checks_commit` and `checks_outcome`.

It also names the sections of ADR-0254 it amends, as the story's last criterion asks. It waits for no task: it is the first layer.

## Done when

- A proposed ADR with sections answering criteria 1 to 6 is committed on `story/S-0659`, refining ADR-0254 and naming the sections it amends.
- `design/adrs/README.md` lists it, as `flai adr new` writes it.

## Notes
