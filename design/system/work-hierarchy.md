---
title: Work item hierarchy and schema
updated: 2026-10-02
status: active
topics: [all]
---

# Work item hierarchy and schema

Work is specified top-down and delivered bottom-up.

```mermaid
flowchart TD
    E[Epic<br/>deliverable spanning several stories] --> S1[Story<br/>one incremental deliverable]
    E --> S2[Story]
    S1 --> T1[Task<br/>a piece of work]
    S1 --> T2[Task]
    S2 --> T3[Task]
```

| Level | Represents | Sized so that | Owned by |
|-------|------------|---------------|----------|
| Epic | A deliverable that spans multiple stories | It closes when its stories close | Human, sets direction |
| Story | One incremental, demonstrable deliverable | An agent can complete it in one to a few sessions | Agent or human |
| Task | A distinct piece of work required by its story | It is done or not done, no partial state | Agent |

Every item, at every level, has a **nature** that says what kind of deliverable it is to the system:

| Nature | Meaning |
|--------|---------|
| `feature` | New capability that did not exist |
| `improvement` | Existing capability made better, faster, or clearer |
| `remediation` | Fixing something that is wrong, including defects and tech debt |
| `research` | Producing knowledge, a spike or investigation, with a written finding as the deliverable |
| `experiment` | Testing a hypothesis with a defined success measure, may be thrown away |

How a nature is accepted and released is in [workflow.md](workflow.md#natures-that-do-not-release-adr-0025-adr-0066): `research` and `experiment` are accepted like any other nature and give no component a bump, and an `experiment` story is accepted only with its results document, `design/experiments/<S-nnnn>-<slug>.md`, whose code and findings may still be dropped ([ADR-0066](../adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md)).

Natures are a closed list. Adding one requires an ADR because the dashboard groups metrics by nature.

## Identifiers

IDs are a type letter, a dash, and a sequence zero-padded to four digits: `E-0001`, `S-0012`, `T-0134`. Sequences are per type and per repo and never reused. `flai` allocates the next ID by scanning `kanban/` and `archive/`; the next number is one more than the highest existing number whatever its width. Repositories created before ADR-0017 carry three-digit IDs; those stay valid side by side with four-digit ones, sorting is numeric, and commands accept an ID with any zero padding or none, so a two, three, or four digit spelling of the same number names the same item. Issues under `design/issues` follow the same width (`I-0012`). `flai migrate ids` widens an existing repository in one step: it renames every item, narrative, and issue with `git mv` and rewrites every reference under `design/`, `docs/`, `wip/`, the root markdown and yaml files, and each project's root markdown files. Run it with `--dry-run` first.

File name: `<ID>-<slug>.md` in `wip/kanban/<epics|stories|tasks>/`.

Stories and tasks may carry `touches`, a list of repository paths or component names the work is expected to change (ADR-0019, S-0037; set with `--touches` on `new` or `flai touches`). It is advisory: `flai check` warns on overlap between in-progress items, and the dashboard shows it. Since S-0128 it is also a claim that holds a ready story ([workflow.md](workflow.md#branches-and-collisions-adr-0019)).

A story may carry `after`, a list of the stories that must be done before it starts, for a dependency that is not about files ([ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md), S-0130; set with `flai edit --after`, cleared with `--clear-after`). While any story it names is not done, a ready story is held (`after`) as an overlapping one is. A cancelled story it names keeps the hold, and so does one that does not exist. `flai check` reports (`story.after`, an error) an entry that names no story, the story itself, or anything but a story, and every cycle, once, on its lowest story. The field is for stories only, in canonical IDs, and is written only when it is not empty. A flai older than S-0130 refused an item that carried it, and with it the listing that read the item (board, `flai serve`, MCP), because front matter was parsed strictly; since S-0181 such a flai reads past the field with a warning, though it does not hold the story (see [Fields an older flai does not know](#fields-an-older-flai-does-not-know)).

Stories and epics may carry `topics`, a list of words naming what the work is about beyond the components it reaches, such as `logging` or `release` ([ADR-0047](../adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md), S-0135; set with `--topics` on `flai story new` and `flai epic new`, `flai edit --topics`, cleared with `--clear-topics`, and from the dashboard). Each is one word of letters, digits, dot, dash, or underscore. `tags` keep their one meaning, the component a story delivers to; `topics` choose what an agent is primed with. A story's topics are the union of:

- its own `topics` and its epic's;
- the name, tags, and kind in `system-flow.yaml` of every sub-project that one of its `tags` names, or that its claim reaches: its `touches` and those of its open tasks, an entry reaching a sub-project when it is the sub-project's name or one of its tags, or when it and the sub-project's path overlap (ADR-0046's rule);
- `code`, when one of those sub-projects is not the template;
- `all`.

A touch outside every sub-project, such as `design/system`, and a tag that names none add nothing. `flai show S-nnnn` prints the story's topics with where each came from, and `--json` returns every source (`own`, `epic`, `tag`, `claim`, `code`, `all`, with the item, the sub-project, and the tag or touch). The topics stories and epics declare are words that `flai check` accepts on a document's topics. A task carries none. As with `after`, a flai older than S-0135 refuses an item that carries `topics`: upgrade the flai on the host first.

Epics, stories, and tasks may carry `usage`, what agents spent on them: the tokens each model read and wrote, what they cost in US dollars, and the seconds of agent work ([ADR-0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md), S-0143). No command takes it as an argument; flai writes it:

- `flai serve` measures a story from the logs it keeps of the agents it started for the story, and each of its tasks from the part of those logs when the task was in progress, when an agent ends and when a task of a running agent enters done; `flai serve agent usage --write` does the same by hand. Such usage has `source: log`. How the logs are read is in [flai-cli.md](flai-cli.md).
- Whenever an item enters done, by `flai move`, MCP's `item_move`, or `flai accept`, each item above it whose usage was not measured is given the sum of its children's, archived ones and cancelled ones included, up to its epic, with `source: sum`. A story measured from its log keeps its measurement, which already holds its tasks'; its epic sums it with its other stories.
- A task's usage is its story's session totals in the share of their input and cache tokens its window holds, so it is `estimated`; so is any cost the harness did not report, priced at the rate the logs report for the model.

Only agents flai serve starts with the `claude-code` harness are measured: their stream-json logs are the only record flai reads. As with `after`, a flai older than S-0143 refuses an item that carries `usage`: upgrade the flai on the host before any item is measured.

## States

One state machine for all three levels:

```mermaid
stateDiagram-v2
    [*] --> backlog
    backlog --> ready
    ready --> in_progress
    in_progress --> review
    in_progress --> done: tasks only
    review --> in_progress: changes requested
    review --> done
    backlog --> cancelled
    ready --> cancelled
    in_progress --> cancelled
    review --> cancelled: only with its parent
    done --> [*]
    cancelled --> [*]
```

| State | Meaning | Clock |
|-------|---------|-------|
| `backlog` | Captured, not committed, may be one paragraph | Lead time starts at `created` |
| `ready` | Refined enough to start: goal and acceptance criteria present. Tasks are written by the agent that starts the story | Queue time |
| `in-progress` | Being worked | Cycle time starts on first entry |
| `review` | Work complete, awaiting verification or human acceptance | Review time |
| `done` | Accepted | Cycle and lead time end |
| `cancelled` | Will not be done, reason recorded | Terminal |

Tasks may go straight from `in-progress` to `done`: review is the human acceptance step and happens at story level, so a task-level review column would only add ceremony. Stories and epics must pass through `review`.

Blocked is not a state. It is a flag with a timestamped interval so the item keeps its column and blocked time is measured separately. See the `blocked` key below.

## Front matter schema

```yaml
---
id: S-0004
type: story                      # epic | story | task
nature: feature                  # feature | improvement | remediation | research | experiment
title: CLI scaffold and config
status: ready                    # backlog | ready | in-progress | review | done | cancelled
parent: E-0002                    # required for task; optional for story (S-0092), absent for epic
owner: agent                     # free text: agent, a person's handle, or team
created: 2026-09-15T16:10:00Z
updated: 2026-09-15T16:40:00Z
transitions:                     # append-only, one entry per state change
  - to: ready
    at: 2026-09-15T16:40:00Z
    by: agent
blocked:                         # optional, append-only, open interval has no `until`
  - from: 2026-09-16T09:00:00Z
    until: 2026-09-16T11:30:00Z
    reason: waiting on template repo access
estimate: 4h                     # optional, Go duration, used for estimate-vs-actual
tags: [cli, config]              # optional, free text
topics: [logging]                # stories and epics, optional (S-0135): what it is about beyond its components
stream: S-0004                    # tasks only: the narrative file in wip/agents they report to
after: [S-0002, S-0003]          # stories only, optional (S-0130): stories that must be done before it starts
agent:                           # stories only, optional (S-0103): who works it
  harness: claude-code
  model: claude-opus-5-5
  config:                        # optional: options for the harness
    effort: high
  roles:                         # optional (S-0189): the agents of the sub-agents its agent hands work to
    verify:
      model: sonnet
usage:                           # optional, written by flai (S-0143): what agents spent on it
  source: log                    # log: measured from the agents' logs | sum: summed from its children
  seconds: 1083                  # agent work
  estimated: true                # optional: some cost was apportioned or estimated, not reported
  models:                        # one entry per model, by name
    - model: claude-opus-5-5
      input: 256                 # tokens
      output: 89342
      cache_read: 19723140
      cache_write: 327605
      cost: 8.1258               # US dollars
---
```

Rules:

- `created`, `updated`, and every `at` are UTC ISO 8601 timestamps.
- Titles and other free text that contain a colon followed by a space, or start with a YAML-special character, are double-quoted. `flai` writes them that way; hand edits must too.
- `transitions` is the source of truth for state. `status` must equal the `to` of the last transition; `flai check` enforces this. Creation implies `backlog` and is not recorded as a transition; a later move back to `backlog` is ([ADR-0055](../adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)).
- `started` and `completed` are not stored. They are derived as the first `in-progress` transition and the last transition while the item is `done` or `cancelled`. See [metrics.md](metrics.md).
- An epic cannot be `done` while any child story is not `done` or `cancelled`. A story cannot be `done` while any child task is not `done` or `cancelled`.
- A cancelled epic has no open story and a cancelled story has no open task: cancelling a parent cancels what is open under it, including an item in `review`, which can be cancelled in no other way ([ADR-0028](../adrs/0028-cancelling-an-item-cancels-everything-open-under-it.md)). `flai check` reports a tree where this does not hold.
- Who changes what after an item is made (S-0085). `title`, `nature`, `tags`, `touches`, `parent`, a story's or epic's `topics`, a story's `after` and `agent`, and the body below the heading are the item's own words and change with `flai edit`, or from a story's or an epic's page in the dashboard, which runs it. A title also lives in the heading, the file's name, the parent's list, and a story's narrative; `flai edit` keeps them in step, a hand edit does not. `id`, `type`, `status`, `transitions`, `blocked`, `owner`, `created`, `updated`, and `usage` are the item's state and change only through the commands that own them (`flai move`, `flai block`, `flai accept`, and for `usage` `flai serve`). No key is added to the front matter to record an edit: a key an older flai does not know is one it cannot act on.
- A story's `agent` names the harness, the model, and the harness's options that work it (S-0103, [ADR-0037](../adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)). A story made while the project has a default agent (`agent` in `system-flow.yaml`) gets a copy, with what `flai story new --harness --model --agent-config` or the dashboard gives laid over it. A story made with neither has no `agent` key, so a project that does not use agents stays readable by an older flai. Only a story carries one. `flai edit --harness/--model/--agent-config/--clear-agent`, MCP's `item_edit`, and the story's page change it. Its `roles` (S-0189) map a sub-agent role, `explore` or `verify`, to a harness, model, and config of that role's own; a story gets the default's, merged role by role, and `--role-harness`, `--role-model`, `--role-config role.key=value`, and `--unset-role` set them. `flai serve` starts `claude-code` with each role's model over the project's sub-agent definition ([agent-context.md](agent-context.md#sub-agents)). The dashboard shows roles and keeps them on a save; flai sets them.
- A story cannot be `ready` without an acceptance criteria section with at least one checkbox. It can be `ready` and `in-progress` with no tasks, and cannot be `review` without at least one ([ADR-0021](../adrs/0021-story-ready-without-tasks.md)).

### Fields an older flai does not know

The flai installed on the host serves MCP and starts agents, and it lags the tree whenever a release is accepted and not installed (S-0181). Until S-0181, work items, threads, and issues were decoded strictly, so one item carrying a field that flai did not know made it refuse every listing that read the item (I-0051). Now:

- The paths that read work items, threads, and issues to list or serve them (`workitem.ParseItem`, `threads.Parse`, `issues.Parse`) decode leniently. Each top-level key the struct does not name is kept on the value (`Unknown`, a `workitem.Field` with the key's lines as written), and `workitem.WarnUnknown` logs `front matter has fields this flai does not know` with the file and the fields, once per process for each file and set of fields.
- A write gives those keys back unchanged: `Marshal` appends them, as written, after the fields it knows. No write drops a field it did not understand.
- `flai check` stays strict: each unknown key is an error, `item.unknown-field`, `threads.unknown-field`, or `issues.unknown-field`, on the key's line.
- Reading past a field is not acting on it: a flai that does not know `after` does not hold the story. What keeps the host's flai current is the version check and the minimum: `flai serve`, the MCP `inbox`, and the dashboard say when the host's flai is older than the newest flai release in the project's history, and `flai.minimum` in `system-flow.yaml` stops an older flai before it reads any item ([project-manifest.md](project-manifest.md), [flai-cli.md](flai-cli.md#versions-the-hosts-flai-and-the-tree)).
- The fields flai reads are listed in `flai/internal/workitem/front-matter-fields.txt`, which a test keeps equal to the code, with the keys of a story's `agent` block and of each of its roles too (S-0189): the lenient read ignores a key it does not know inside `agent`, and a write would drop it. A story that adds or removes a front-matter field changes it, and publishing the flai release that carries the change raises `flai.minimum` to that release.

## Body structure

The body below the front matter has fixed headings so agents and the dashboard can find things.

Epic:

```markdown
# E-0002 flai CLI

## Outcome
One paragraph: what is true when this epic is done.

## Stories
Maintained by flai. A list of child IDs and titles.

## Notes
Anything else.
```

Story:

```markdown
# S-0004 CLI scaffold and config

## Goal
What this increment delivers and for whom.

## Acceptance criteria
- [ ] Checkable statements. Required before `ready`.

## Tasks
Maintained by flai. A list of child IDs and titles.

## Notes
Decisions, links to ADRs, anything discovered on the way.
```

Task:

```markdown
# T-0021 Read and write ~/.flai/config.json

## Work
What to do, concretely.

## Done when
A single sentence.

## Notes
```
