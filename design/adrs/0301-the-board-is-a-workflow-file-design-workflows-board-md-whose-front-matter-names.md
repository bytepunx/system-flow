---
id: ADR-0301
title: "The board is a workflow file, design/workflows/board.md, whose front matter names its states with a kind each, its transitions with who may make them, the rules they pass from a vocabulary flai implements, and the effects flai gives them, and its timers; a project with no file runs the baseline flai carries, which is today's board"
status: proposed
date: 2026-11-14
supersedes: []
superseded_by: []
refines: [ADR-0003, ADR-0004, ADR-0055, ADR-0073, ADR-0076, ADR-0096]
topics: [cli, template, conventions, dashboard]
---

# ADR-0301 The board is a workflow file, design/workflows/board.md, whose front matter names its states with a kind each, its transitions with who may make them, the rules they pass from a vocabulary flai implements, and the effects flai gives them, and its timers; a project with no file runs the baseline flai carries, which is today's board

## Context

The board is fixed in flai's code. `flai/internal/workitem/item.go` holds the six states as a closed list (`States`). `flai/internal/workitem/rules.go` holds the transition map (`allowed`, with `Follows` adding a task's move from `in-progress` to `done`) and every rule a move passes (`Repo.Move`). About 320 references in 20 packages compare a state's name to decide what to do: `flai serve` starts an agent for a story in `ready`, kluster, the forge pass, and the merge queue act on a story in `review`, acceptance is the move to `done`. The metrics of `design/system/metrics.md` read the names too. The dashboard keeps three lists of its own. [ADR-0004](0004-workflow-states-and-transitions.md) made the states one closed list shared by every item type; [ADR-0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md) set the moves back. Changing the board needs a flai release.

E-0034 makes the board data a project holds, as [ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) made a strategic role. Before any code reads a workflow file, this ADR fixes what the file is, so that the epic's other stories build to one format.

Today's board, read from the code:

| From | To | Types | Refused unless | Warns when | Code |
|------|----|-------|----------------|------------|------|
| backlog | ready | all | a story has criteria and is not a draft (or the move finalizes it) | a story names in `after` a story not done | `Repo.Move`, `case Ready` and `HoldAfter` |
| ready | in-progress | all | | the item is blocked; a story is held; a task's `after` is not done | `Repo.Move` |
| in-progress | review | all | a story has a task, no open question in its narrative, nothing uncommitted in its worktree | a story has an unticked criterion; its worktree could not be read; kluster is not deployed or checked at its head | `Repo.Move`, `case Review`; `mcpserver.ReviewWarnings` |
| review | done | all | every child is closed; a story's criteria are all ticked; a story or an epic only by acceptance | | `Repo.Move`, `case Done`; `mcpserver` `itemMove`; `cmd/move.go` runs `flai accept` |
| in-progress | done | task | every child is closed | | `Follows` |
| review | in-progress | all | a reason is given | the item is blocked | `Repo.Move` |
| in-progress | ready | all | a story has criteria and is not a draft | a story names in `after` a story not done | `Repo.Move`, `case Ready` |
| ready | backlog | all | | | `allowed` |
| backlog, ready, in-progress | cancelled | all | a reason is given | | `Repo.Move` |
| review | cancelled | all | only as a parent's cancellation takes it along (`Cascade`, [ADR-0028](0028-cancelling-an-item-cancels-everything-open-under-it.md)) | | `Repo.Move` |
| cancelled | backlog | all | the parent is not cancelled | | `Repo.Move` |

Every move into a state with a WIP limit warns when it takes the column of stories over it. A move into `done` or `cancelled` closes an open blocked interval. A story entering `ready` goes last in the pull order, and leaving it, or moving anywhere else, takes it out of the order. Creation is `backlog` and is not a transition.

### Choices for the operator

Each is decided below as recommended. The operator's answer may change it before this ADR is accepted.

1. **Where the file lives.** Recommended: `design/workflows/board.md`, beside `roles/` and `formulas/`. Alternative: `design/workflows/workflow.md`, which does not share a name with `wip/kanban/board.md`.
2. **Where the WIP limits' values live.** Recommended: they stay in `wip/kanban/board.md`, where `flai board limit` and the board's lane menu set them without a story; the workflow file says only which states take one, and a state whose limit `wip/kanban/board.md` does not set has none, as today. Alternative: the values move into the workflow file, and a limit changes only through a story.
3. **How a renamed or removed state's items and history read.** Recommended: a state lists its old names in `aliases`; flai reads every status and transition through them and rewrites nothing, and a removed state's time counts in the state it is read as. Alternatives: a migration command rewrites each item's `status` and `transitions` when the story that changes the file is accepted; or a removed state's time keeps its own name in time in state and counts in no timer.
4. **ADR-0003 and ADR-0004.** Recommended: this ADR refines both, and supersedes neither. Each keeps what it decided but the closed list of states. Alternative: supersede ADR-0004, restating here that transitions are recorded in front matter and blocked is a flag.
5. **A rule's severity.** Recommended: every rule has a default, refuse or warn, and a transition may set it. The baseline needs it: an unticked criterion warns at review and refuses at done. Alternative: a fixed severity per rule, with two rules for the criteria.
6. **How many accepted and cancelled states.** Recommended: exactly one of each, and for stories exactly one pull and one acceptance. Alternative: several of each, which every reader of `done` would have to choose between.
7. **Who may make a transition.** Recommended: a short vocabulary, `anyone`, `operator`, `agent`, `flai`, `acceptance`, and `cascade`, which a role file and its permissions narrow further for a strategic role. Alternative: the workflow names roles and permissions itself, which ADR-0207's role files already hold.
8. **How flai upgrade treats the file.** Recommended: it adds the baseline file when the project lacks it and never changes it afterwards, saying when the template's baseline moved. Alternative: merge it as a baseline role file is merged, the template's front matter over the project's, which would undo the project's board.
9. **Who accepts a story that changes the file.** Recommended: the operator alone, never the orchestrator or the merge queue, as for a protected path (ADR-0106), and the dashboard's saves take the file only on a story's branch. Alternative: any acceptance, the file being reviewed as any other design file is.

## Decision

**The board is one markdown file, `design/workflows/board.md`. Its front matter names the states, each with a kind that tells flai what the state means; the transitions between them, each with the item types it takes, who may make it, the rules it passes from a vocabulary flai implements, and the effects flai gives it; and the timers measured over the transitions every item records. Code reads a state's kind and a transition's effect, never a state's name. A project with no file runs the baseline flai carries, which is today's board exactly.** This refines [ADR-0003](0003-work-item-hierarchy.md), [ADR-0004](0004-workflow-states-and-transitions.md), [ADR-0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md), [ADR-0073](0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md), [ADR-0076](0076-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review.md), and [ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md); the rest of each stands.

### 1. Where it lives and what it holds

- The file is `design/workflows/board.md`, under the design folder the manifest's `layout.design` names. `design/workflows/README.md` lists it with `roles/` and `formulas/`.
- It is markdown with YAML front matter, as every design file is. The front matter is the workflow. The body is prose for the operator: why the board is shaped as it is. flai reads nothing from the body, so the file has no baseline marker.
- One file per project. Epics, stories, and tasks follow it, each where the states and transitions name its type. There is no workflow per epic and no editor in the dashboard.

Every field is plain YAML. `flai check` reports `workflow.file` for a file that does not parse, misses a required field, or gives a value outside the allowed ones.

| Field | Type | Meaning | Allowed values |
|-------|------|---------|----------------|
| `title` | string | The document's title. | Any. Required. |
| `updated` | date | The day of the last meaningful edit. | `YYYY-MM-DD`. Required. |
| `status` | string | The documentation standard's status. It gates nothing. | `draft`, `active`, `deprecated`. Required. |
| `states` | list of mappings | The board's columns, in board order, left to right (§2). | One entry or more. Required. |
| `transitions` | list of mappings | The moves between states (§3). | One entry or more. Required. |
| `timers` | list of mappings | The durations measured over every item's transitions (§5). | Default `[]`; flai's figures then lack the timers. |

A state is a mapping:

| Field | Type | Meaning | Allowed values |
|-------|------|---------|----------------|
| `name` | string | The state's name, as `status` and `transitions[].to` hold it. | Kebab-case; unique among the names and aliases. Required. |
| `kind` | string | What the state means to flai (§2). | `queue`, `active`, `waiting`, `accepted`, `cancelled`. Required. |
| `title` | string | The column's heading on the board. | Any. Default the name. |
| `types` | list of strings | The item types the state takes. | `epic`, `story`, `task`. Default all three. |
| `initial` | boolean | Items are created in this state. | Exactly one state, of kind `queue`. Default `false`. |
| `order` | boolean | The state's stories read in the board's pull order. | Only a `queue` state. Default `false`. |
| `limit` | boolean | The state takes a WIP limit for stories. Its value is `wip/kanban/board.md`'s; with none set there, the state has no limit, as today. | Default `false`. |
| `aliases` | list of strings | Names the state had, which flai reads as this state (§6). | Kebab-case; never a current state's name. Default `[]`. |

A transition is a mapping:

| Field | Type | Meaning | Allowed values |
|-------|------|---------|----------------|
| `from` | string or list | The state or states it leaves. | States of the file. Required. |
| `to` | string | The state it enters. | A state of the file. Required. |
| `types` | list of strings | The item types that may make it. | `epic`, `story`, `task`. Default all three. |
| `by` | list of strings | Who may make it (§3). | `anyone`, `operator`, `agent`, `flai`, `acceptance`, `cascade`. Default `[anyone]`. |
| `rules` | list of rule entries | What the move must pass (§4). | A rule's name, or `{rule: <name>, severity: refuse \| warn}`. Default `[]`. |
| `effects` | list of strings | What flai does with the move besides recording it (§3). | `commit`, `pull`, `accept`, `send-back`, `cancel`, `reopen`. Default `[]`. |

No two transitions share a `from`, a `to`, and a type. Any move may carry a reason, which flai appends to the item's Notes as today; the rule `reason-given` makes one required.

### 2. State kinds

A state's kind is what lets flai, the check, and the charts work on a board they have never seen. A state's name is only a name.

| Kind | Means | What flai does with an item in it |
|------|-------|-----------------------------------|
| `queue` | Waiting to be worked; nobody is on it. | Holds no claim. A story in the pull's source state waits for the pull: `flai serve` starts its agent and `wait_for_work` offers it while the pull's target has room (ADR-0043), held by an overlap or `after` (ADR-0230). The state with `order` reads in the pull order. |
| `active` | Being worked. | The story's agent runs, and `flai serve` restarts it when it ends with work left (ADR-0108). Its claim holds ready stories. `stream_state`, `task_done`, and `flai stream sync` take it. |
| `waiting` | Waiting on someone else. | The agent's run ended with its work done. Its claim holds nothing (ADR-0096). Its time is the review wait `flai stats` reports. The waiting state the acceptance leaves is where kluster deploys and checks the story, the forge pass opens its pull request, and the merge queue batches it. |
| `accepted` | Closed and delivered. | Final: no transition leaves it (ADR-0055). The item is completed, archived by `flai archive`, and counts toward a release by its nature. |
| `cancelled` | Closed and not delivered. | The item is completed until it reopens. Its stories leave lead and cycle time aggregates and count in the cancellation rate. |

A move into an `accepted` or `cancelled` state closes an open blocked interval (I-0347). A full `active` or `waiting` state with a limit holds the pull, as a full review does today (ADR-0073). The order of `states` is the board's: the dashboard draws the columns in it, and an epic walks it.

An epic follows its stories by kind (ADR-0076): to the accepted state when every story is accepted; to the acceptance's source when every story is there or accepted; to the pull's target when any story has started (is in an `active`, `waiting`, or `accepted` state); to the pull's source when any story is there; else to the initial state. It walks the board's order one transition at a time, as `Follow` does.

`flai check` reports `workflow.state` for a workflow without exactly one initial state, one `accepted` state, and one `cancelled` state; for a state no transition reaches but the initial one; for a state with no way to the accepted or the cancelled state; and for a transition out of the accepted state.

Today's six states:

| State | Kind | Types | WIP limit | Order | Initial |
|-------|------|-------|-----------|-------|---------|
| `backlog` | `queue` | all | none | yes, placed by hand | yes |
| `ready` | `queue` | all | takes one (the template's `board.md` sets 5) | yes, last on entry | no |
| `in-progress` | `active` | all | takes one (2) | no | no |
| `review` | `waiting` | all | takes one (3) | no | no |
| `done` | `accepted` | all | none | no | no |
| `cancelled` | `cancelled` | all | none | no | no |

A story entering the pull's source state goes last in its order; leaving it, or entering any other state, takes it out of the order; the operator places stories by hand in either ordered state (S-0057).

### 3. Transitions, who makes them, and their effects

**Who.** `by` bounds who may make a move. A role file and the operator's permissions narrow it further for a strategic role, as ADR-0207 §5 says, and the guard still refuses every sub-agent's move (ADR-0060).

| `by` | Who may make the move |
|------|------------------------|
| `anyone` | Every author flai lets write items: a person, a story's agent on its story and tasks, `flai serve`, and a strategic role whose file allows it. |
| `operator` | A person, never an agent `flai serve` started. |
| `agent` | The story's own agent, on its story and its tasks. |
| `flai` | `flai serve` and the merge queue, as `sendBackBy`. |
| `acceptance` | Only an acceptance: `flai accept`, `flai move` of a story from its waiting state to done, the dashboard's accept, the orchestrator's `accept` under `accept_reviews`, and the merge queue's landing. `item_move` never makes it. |
| `cascade` | Only as a parent's cancellation takes the item along. |

**Effects.** An effect is what flai does once, at the move. What flai does with an item while it sits in a state is the state's kind (§2).

| Effect | flai does |
|--------|-----------|
| `commit` | The move that commits an item to be worked. Its rules are the definition of ready; the orchestrator's `promote_to_ready` and `flai promote --candidates` make it; a planner never does. Its first time is `committed`. |
| `pull` | The move that starts a story's work: `flai story start` and `story_start` make it, `flai serve` starts the agent for it, holds and claims apply. Exactly one for stories, from a `queue` state to an `active` one. |
| `accept` | Acceptance: the acceptance gate (`flai-cli.md` § The acceptance gate), the merge, the move, the archive, and the commit; the usage roll-up; the epic's acceptance with its last story; the release accounting; the curator's and the replanner's triggers. Exactly one for stories, from a `waiting` state to the `accepted` one, with `by: [acceptance]`. |
| `send-back` | Work comes back from a wait: `flai serve` starts the agent again, a thread flai opened on the send back is resolved once the story returns (ADR-0248), kluster, forge, and the merge queue drop the story. |
| `cancel` | Everything open under the item is cancelled with it (ADR-0028), the agent is stopped, its kluster namespace is torn down. |
| `reopen` | The item comes back from cancelled and is no longer completed (ADR-0055). |

A move back one column that is no send back, such as `in-progress` to `ready`, has no effect: the state it enters says what follows, as a story entering `ready` may be pulled again.

Today's transitions, as the baseline names them (§7 gives the file):

| From | To | Types | By | Effects | Rules |
|------|----|-------|----|---------|-------|
| backlog | ready | all | anyone | commit | criteria-present, not-draft, after-done |
| ready | in-progress | all | anyone | pull | not-blocked, not-held, task-after |
| in-progress | review | all | anyone | | has-task, no-open-question, clean-worktree, criteria-ticked (warn), kluster-verified |
| review | done | epic, story | acceptance | accept | children-closed, criteria-ticked |
| review | done | task | anyone | | children-closed |
| in-progress | done | task | anyone | | children-closed |
| review | in-progress | all | anyone | send-back | reason-given, not-blocked, task-after |
| in-progress | ready | all | anyone | | criteria-present, not-draft, after-done |
| ready | backlog | all | anyone | | |
| backlog, ready, in-progress | cancelled | all | anyone | cancel | reason-given |
| review | cancelled | all | cascade | cancel | reason-given |
| cancelled | backlog | all | anyone | reopen | parent-open |

### 4. Rules

A rule is a check flai implements in Go, in `flai/internal/workflow/rules.go`, named in kebab-case, as a role file's conditions are (ADR-0207 §5). A file names one; it cannot define one. A rule applies only to the item types it names below; a transition's `types` narrows it further. Each has a default severity, and a transition's entry may set `severity`: `refuse` stops the move and changes nothing, `warn` lets it through and says why. An unknown rule is `workflow.rule`. A rule the vocabulary lacks is a story that adds it.

| Rule | Checks | Types | Default |
|------|--------|-------|---------|
| `criteria-present` | The body has `## Acceptance criteria` with at least one checkbox (ADR-0021). | story | refuse |
| `not-draft` | The story is not a draft, or the move finalizes it (`--yes`, `finalize`). | story | refuse |
| `has-task` | The story has at least one task. | story | refuse |
| `no-open-question` | The story's narrative has no open question. | story | refuse |
| `clean-worktree` | Nothing is uncommitted in the story's worktree. A worktree git cannot read warns and lets the move through (I-0017). | story | refuse |
| `criteria-ticked` | Every acceptance criterion is ticked. | story | refuse |
| `kluster-verified` | The story's namespace is deployed at its branch's head and its checks passed there, where the manifest's `kluster:` block names a deploy and the `kluster` host action is on (ADR-0292). | story | warn |
| `children-closed` | Every child is in an `accepted` or `cancelled` state. | epic, story | refuse |
| `reason-given` | The move gives a reason. A cascade gives its parent's. | all | refuse |
| `parent-open` | The item's parent is not in a `cancelled` state. | all | refuse |
| `not-blocked` | The item has no open blocked interval. | all | warn |
| `not-held` | The story is not held by an overlap or by `after` (ADR-0046, ADR-0230). | story | warn |
| `after-done` | Every story the story names in `after` is accepted (ADR-0279). | story | warn |
| `task-after` | Every task the task names in `after` is closed (S-0176). | task | warn |

The WIP limit is no rule a transition lists: every move into a state with `limit` warns when it takes the column's stories over the limit `wip/kanban/board.md` sets for it. The acceptance gate stays the `accept` effect's: what `flai accept` checks beyond the transition's rules, the verify record, the merge queue's tiers, the pull request's checks, and kluster at the head, is unchanged by this ADR.

### 5. Timers

Every state change is already recorded in an item's `transitions` (ADR-0004), so a timer is a reading of that history. Adding or changing one needs no new data and applies to past items too.

A timer has one of two forms:

- **Between two points**, `from` and `to`: the time from one to the other. It has no value until both exist.
- **A sum**, `sum`: the time the item spent in the states it names, an open interval counted up to now. With `until`, only the intervals that end at or before that point count. It has no value for an item never in its states, or, with `until`, before the point exists.

A point is one of these:

| Point | Is |
|-------|----|
| `created` | The item's front matter `created`. |
| `completed` | The `at` of the item's last transition while it is in an `accepted` or `cancelled` state; none while it is open (ADR-0055). |
| `{first: <state>}` | The `at` of the first transition into the state. |
| `{last: <state>}` | The `at` of the last transition into the state. |
| `{first: {kind: <kind>}}`, `{last: {kind: <kind>}}` | The same, into any state of the kind. |

A timer is a mapping:

| Field | Type | Meaning | Allowed values |
|-------|------|---------|----------------|
| `name` | string | The timer's name in `flai stats` and the charts. | Kebab-case; unique. Required. |
| `title` | string | Its heading on the charts. | Any. Required. |
| `from`, `to` | point | The two points, for the first form. | Points above. Both, or neither. |
| `sum` | list of strings | The states summed, for the second form. | States of the file. Required without `from` and `to`. |
| `until` | point | Bounds a sum. | A point. Only with `sum`. |
| `types` | list of strings | The item types it is measured for. | `epic`, `story`, `task`. Default all three. |
| `counts` | string | Which completed items its aggregates take. | `accepted`, items closed in the accepted state; `closed`, cancelled ones too. Default `accepted`. |

A timer naming a state the file lacks is `workflow.timer`. Time in state, the sum per state of every state, and blocked time stay figures flai gives for every state, not timers.

The baseline's timers are today's metrics:

| Timer | Form | Today's definition in `metrics.md` |
|-------|------|------------------------------------|
| `lead` | `from: created`, `to: completed` | Lead time, `completed - created` |
| `cycle` | `from: {first: in-progress}`, `to: completed` | Cycle time, `completed - started` |
| `queue` | `sum: [ready]`, `until: {first: in-progress}` | Queue time, the `ready` intervals that end at or before `started` (ADR-0190); its fallback, `started - committed`, is the same sum for an item that went from `ready` to `in-progress` alone |
| `review` | `sum: [review]` | `wait_review_seconds`, the seconds in `review`, the open interval up to now |

`counts: accepted` keeps cancelled items out of lead and cycle time aggregates, as `metrics.md` says. An operator's own timer, such as the time from the first move to review to the acceptance, is `{name: acceptance-wait, title: Acceptance wait, from: {first: review}, to: completed}`.

This ADR does not change `metrics.md`. S-0750 computes the four figures from the baseline's timers, byte for byte as today, and writes the ADR that changes the metrics contract; S-0752 reports a workflow's own timers.

### 6. Changing a workflow with items on the board

The transitions an item records are its history, reviewable in git (ADR-0004). flai never rewrites them, nor an item's `status`, when the workflow changes.

- **Adding a state** needs nothing else: no item is in it yet.
- **Renaming a state** keeps its old name in the state's `aliases`. flai reads a `status` or a transition's `to` that names an alias as the state that lists it. Every write names the current state. The history reads under the new name, and the timers count it there.
- **Removing a state** puts its name in the `aliases` of the state its items go to. Its items read as being in that state, and the time they spent in it counts as time in that state, in time in state and in every timer.
- **Changing a state's kind** changes how flai reads all of it, past intervals included: a kind is not recorded per transition.
- A name stays an alias for as long as any item, archived ones included, names it. It is never a current state's name, and never two states' alias.

`flai check` reports `workflow.item` for an item whose `status` or a transition names a state the file neither has nor lists as an alias, saying to add it as an alias of the state its items go to. It reports `board.wip-state` for a limit in `wip/kanban/board.md` on a state that takes none.

### 7. The baseline, and how the file changes

- **No file.** A project with no `design/workflows/board.md` runs the baseline flai carries in its binary (`flai/internal/workflow/baseline/board.md`), which is today's board exactly: the same six states, the same transitions, the same rules, the same figures.
- **A file.** It replaces the baseline whole. flai merges nothing from the baseline into it.
- **The template** ships the baseline as `design/workflows/board.md`. `flai upgrade` adds it to a project that lacks it, and never changes a project's file afterwards; when the template's baseline changed, it says so, and the operator carries the change over in a story.
- **It is code the operator runs.** Like every file under `design/workflows/`, it changes only through a story the operator accepts. No role's `files` may name the folder (ADR-0207 §2), so no role edits it. A story whose branch changes it is accepted by the operator alone, never by the orchestrator or the merge queue. The dashboard's document and file saves refuse it on main, and take it only on a story's branch (ADR-0145).

The baseline file:

```yaml
---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: in-progress, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - from: in-progress
    to: review
    rules: [has-task, no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}, kluster-verified]
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked]}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: in-progress, to: done, types: [task], rules: [children-closed]}
  - {from: review, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: in-progress, to: ready, rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  - {name: cycle, title: Cycle time, from: {first: in-progress}, to: completed}
  - {name: queue, title: Queue time, sum: [ready], until: {first: in-progress}}
  - {name: review, title: Review time, sum: [review]}
---
```

### 8. ADR-0003 and ADR-0004

This ADR refines both and supersedes neither. ADR-0003's three levels and five natures stay closed lists; the types a state or a transition takes are those levels. ADR-0004's transitions recorded in front matter, `status` as the last transition's copy, blocked as a flag, and derived timestamps computed, never stored, all stand. Only its one closed list of states shared by every type becomes the workflow's states, each with the types it takes. ADR-0055's moves back become the baseline's transitions; ADR-0073's full review and ADR-0096's review that holds nothing, and ADR-0076's epic following, read kinds.

### 9. Which story builds what

Until each story below lands, flai keeps deciding from its Go, and the baseline is what that Go does.

| Section | Stories |
|---------|---------|
| §1, the baseline file flai carries, the states and transition map | S-0737 |
| §4, the rules | S-0738 |
| §2, kinds read in place of names | S-0739 (workitem), S-0742 and S-0743 (`flai serve`), S-0744 (`flai move`, `flai accept`), S-0746 (MCP), S-0747 (host API), S-0749 (the rest), S-0753 and S-0758 (the dashboard) |
| §2, §6, the check | S-0741 |
| A role file's `to-review` and `to-backlog` as one condition naming any state | S-0748 |
| §5, the timers | S-0750 (the baseline's, with the metrics ADR), S-0752 (a workflow's own), S-0754 (the charts) |
| §1, §6, §7, a project's own file, the template, upgrade, the guides | S-0755 |
| The epic's outcome, end to end | S-0756 |

## Consequences

- An operator changes the board by editing one file in a story, with no flai release. The board's rules are reviewed as data, and each is tested Go.
- Code that compared a state's name reads its kind or a transition's effect. A state's name appears only in the baseline file and in the items.
- A project that runs the baseline sees no change: the same moves, refusals, warnings, and figures.
- Every rule is one flai implements. A project that needs another asks for it, which is slower than writing it and keeps the moves' rules in tested code.
- History is never rewritten. A renamed or removed state's name lives on as an alias for as long as an item names it.
- A timer added today reports past items too, since it reads the transitions they already carry.
- `wip/kanban/board.md` keeps the limits' values and the pull order. The workflow file says which states take them.
- A story that changes `design/workflows/board.md` needs the operator's acceptance, as one that changes a protected path does (ADR-0106).
- `design/system/workflow.md` and `design/system/work-hierarchy.md` describe the states as the baseline's, not as a closed list.

## Alternatives considered

- **Rules as expressions in the file**, such as `when: criteria.unticked == 0`. The moves would need a small language to evaluate, and a file could write a rule no test covers. Named rules keep the logic in tested Go, as ADR-0207 chose for conditions.
- **Effects read from a state's name**, as today. A board with a state between `in-progress` and `review` would then need code to know what it means. Kinds and effects make a name only a name.
- **Every effect as a kind**, with no transition effects. Acceptance and the send back happen once at a move, not while an item sits in a state; a kind cannot say which move out of `review` is the acceptance.
- **A `submit` effect on the move to review.** The waiting kind already ends the agent's work and holds nothing, and the acceptance's source is where kluster, the forge, and the merge queue act; a fourth effect would say it twice.
- **A workflow per epic, or per item type.** The epic leaves both out. One file whose states and transitions name their types covers today's one difference, a task's move from `in-progress` to `done`.
- **A YAML file, `board.yaml`.** The documentation standard keeps design files markdown with front matter; a role file and a formula are too.
- **Rewriting the items when a state is renamed or removed.** It would rewrite history ADR-0004 keeps reviewable, in many files at once on main. Aliases change one file.
- **The baseline merged under a project's file**, a project's file naming only what it changes. A reader could not see the whole board in one place, and removing a baseline state would need a way to say so.
- **The marker below which a project adds its own.** The whole file is the project's, and upgrade never merges it, so a marker would divide nothing.
