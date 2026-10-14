---
id: T-1487
type: task
nature: improvement
title: Design, convention, template, and user docs describe the settled reply, the blocked hold on wakes, and the notice without generated paths
status: done
parent: S-0372
owner: alex
created: 2026-10-12T17:23:37Z
updated: 2026-10-14T21:40:22Z
transitions:
  - to: ready
    at: 2026-10-14T21:36:47Z
    by: agent-S-0372
  - to: in-progress
    at: 2026-10-14T21:36:48Z
    by: agent-S-0372
  - to: done
    at: 2026-10-14T21:40:22Z
    by: agent-S-0372
stream: S-0372
tags: [messages, docs, conventions]
touches: [design/system/agent-coordination.md, design/system/workflow.md, design/system/agent-narrative.md, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1483, T-1485, T-1486]
usage:
  source: log
  seconds: 214
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 22428
      cache_read: 3106003
      cache_write: 143879
      cost: 1.8792
---
# T-1487 Design, convention, template, and user docs describe the settled reply, the blocked hold on wakes, and the notice without generated paths

## Work

Write down what T-1483, T-1485, and T-1486 built, linking the ADR of T-1482:

- `design/system/agent-coordination.md` › Messages between stories: the settled entry, and that a conversation awaits neither story until the next message.
- `design/system/workflow.md` › Branches and collisions (ADR-0019): the task-done notice leaves out the paths sync regenerates. Its "Asked" and "An answered agent" rules gain the settled reply and the blocked hold.
- `design/system/agent-narrative.md` › What an agent is told through MCP: `message_reply`'s `settled`, a settled conversation's `awaiting` in `inbox`, and the no-reply text of its `message` event.
- `design/system/flai-cli.md`: the `flai message` row gains `reply --settled`, and how `list`, `show`, and `--json` show a settled entry.
- `design/conventions/work-management.md` › Rules, in the baseline, and its copy `template/root/design/conventions/work-management.md`:
  - The rule on answering a message: when your answer asks nothing, reply settled.
  - The rule on a task-done notice: "Reply on the conversation either way" becomes "reply settled when the change does not break your work". Keep the two copies the same.
- `template/CHANGELOG.md`: an entry for the convention change.
- `docs/users/flai.md` › Work items › Messages between stories: the settled reply, and when to use it.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference`, so that it carries `--settled`.

It waits for T-1483, T-1485, and T-1486, because it describes what they built.

## Done when

- Each document above says what the code does, and links the ADR where it states a rule.
- The two copies of the work-management convention say the same in the changed rules.
- `docs/users/flai-reference.md` matches `make flai-reference`'s output.
- `flai check --strict` and `flai test` pass on the changed files.

## Notes

The baseline of `work-management.md` changes here. That is done in this repository, the template's source, in the same story as the design, as CLAUDE.md says. The plan thread on S-0372 asks the operator to agree to it.
