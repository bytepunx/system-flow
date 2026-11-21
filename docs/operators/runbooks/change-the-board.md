---
title: "Runbook: change the board"
updated: 2026-11-18
status: active
---

# Change the board

The board's states, the moves between them, the rules each move passes, and the timers are one file, `design/workflows/board.md` ([ADR-0301](../../../design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md)). It is code you run: change it in a story, and accept that story yourself. flai reads the file of the checkout it runs in, so the story's worktree runs the new board before main does ([ADR-0307](../../../design/adrs/0307-flai-reads-the-workflow-file-of-the-checkout-it-opens-refuses-to-open-a.md)).

## Before you start

1. Open a story for the change, and name `design/workflows/board.md` in its touches. Work in its worktree.
2. Find the file. A project with no `design/workflows/board.md` runs the baseline flai carries, which is the template's file byte for byte. To change it, add the file first, in the story's worktree:
   - copy `root/design/workflows/board.md` from the template repository, or
   - let `flai upgrade` add it, which it does for a project that lacks it.
3. Edit the front matter only. It is the whole workflow, and it replaces the baseline whole: flai merges nothing into it. The body is prose for people; flai reads none of it. Every field and its allowed values are in ADR-0301 §1, and the keys that name test suites in [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md).

## Add a state

1. Add it to `states`, in board order, left to right. Give it a `kind`, which tells flai what it means:

   | Kind | The state is |
   |------|--------------|
   | `queue` | Waiting to be worked; nobody is on it |
   | `active` | Being worked; the story's agent runs |
   | `waiting` | Waiting on someone else; the agent's work is done |
   | `accepted` | Closed and delivered; nothing leaves it |
   | `cancelled` | Closed and not delivered |

   A `qa` state between `in-progress` and `review`:

   ```yaml
   states:
     - {name: backlog, kind: queue, initial: true, order: true}
     - {name: ready, kind: queue, order: true, limit: true}
     - {name: in-progress, kind: active, limit: true}
     - {name: qa, kind: waiting}
     - {name: review, kind: waiting, limit: true}
     - {name: done, kind: accepted}
     - {name: cancelled, kind: cancelled}
   ```

   `types` limits the state to some item types; it takes all three by default.

2. Add the moves into it and out of it. Every state must be reachable from the initial one, and must lead to the accepted or the cancelled state. Here the move to `review` becomes the move to `qa`:

   ```yaml
   transitions:
     # in place of the baseline's in-progress -> review
     - from: in-progress
       to: qa
       rules: [has-task, no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}]
     - {from: qa, to: review, rules: [kluster-verified]}
     - {from: qa, to: in-progress, effects: [send-back], rules: [reason-given]}
     - {from: qa, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
   ```

   `by` says who may make the move: `anyone` (the default), `operator`, `agent`, `flai`, `acceptance`, or `cascade`. `effects` says what flai does at the move: `commit`, `pull`, `accept`, `send-back`, `cancel`, or `reopen`. Stories keep exactly one `pull` move and one `accept` move. No two moves share a `from`, a `to`, and a type. The effects are in ADR-0301 §3.

3. For a WIP limit, give the state `limit: true`. The value is not in the workflow file: it stays in `wip/kanban/board.md`, which `flai board limit` and the board's lane menu set. A state with `limit: true` and no value there has no limit. `flai board limit <state> <n>` takes any state with `limit: true`, one you add or rename included, though its help still names only `ready`, `in-progress`, and `review` ([I-0403](../../../design/issues/I-0403-a-refusal-and-a-command-s-help-name-the-baseline-s-states-where-a-project-s.md)).

   An added state of the `active` kind runs the story's agent as `in-progress` does, but `flai serve` stops the story's dev server there, since it still reads the baseline's kinds for that ([I-0402](../../../design/issues/I-0402-flai-serve-stops-the-dev-server-of-a-story-in-a-state-a-project-s-workflow-adds.md)).

## Add a rule to a move

Add the rule's name to the move's `rules`. A rule is a check flai implements; the file names one and cannot define one. The rules, what each checks, the item types it applies to, and its default severity are in ADR-0301 §4. `flai check` lists the names when the file gives one it does not know.

Each rule refuses or warns by default. Set `severity` to change it for one move: `refuse` stops the move, `warn` lets it through and says why.

```yaml
  # the pull refuses a held story, where the baseline only warns
  - {from: ready, to: in-progress, effects: [pull], rules: [not-blocked, {rule: not-held, severity: refuse}, task-after]}
```

A rule applies only to the item types it names, and a move's `types` narrows it further. For a rule flai lacks, ask for a story that adds it.

The rule `verified` comes with a move's `verify` ([below](#name-the-suites-a-step-runs)): you need not list it. List it in `rules` to place it among the others; listed on a move with no `verify`, it needs every suite.

## Add a timer

A timer reads the moves every item already records, so it covers past items too. It is the time between two points, or the sum of the time spent in some states:

```yaml
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  # ... the baseline's other timers
  - {name: acceptance-wait, title: Acceptance wait, from: {first: review}, to: completed}
  - {name: qa-time, title: QA time, sum: [qa]}
```

A point is `created`, `completed`, `{first: <state>}`, `{last: <state>}`, or `{first: {kind: <kind>}}` and `{last: {kind: <kind>}}`. `until: <point>` bounds a sum. The fields are in ADR-0301 §5. `flai stats` reports the baseline's four as lead, cycle, queue, and review time, and each of the workflow's own timers after them, for past items too ([Workflow timers](../../../design/system/metrics.md#workflow-timers-s-0752)).

## Rename a state

flai never rewrites an item's `status` or its history. Give the state its new name and list the old one in its `aliases`, then rename it in every move and timer:

```yaml
states:
  - {name: doing, kind: active, limit: true, aliases: [in-progress]}
transitions:
  - {from: ready, to: doing, effects: [pull], rules: [not-blocked, not-held, task-after]}
```

flai reads an item that names `in-progress` as being in `doing`, counts its time there, and writes `doing` from then on. Keep the alias while any item, archived ones included, names it.

## Remove a state with items in it

Take the state and its moves out, and list its name in the `aliases` of the state its items go to:

```yaml
states:
  - {name: review, kind: waiting, limit: true, aliases: [qa]}
```

Its items read as being in `review`, and the time they spent in `qa` counts as time in `review`, in every figure and timer.

Changing a state's `kind` changes how flai reads all of the state's time, past time included.

## Name the suites a step runs

Each test tier in `system-flow.yaml`'s `tests` belongs to a suite: its `suite`, or its own name when it gives none. A project made from the template has `static`, `behavior`, `integration`, and `smoke` ([Run the tests for what changed](../../users/flai.md#run-the-tests-for-what-changed)). The workflow file names suites in three places ([ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md)):

- `verify` on a move: the suites a story's move needs its last `flai verify` to have passed.
- `runs` on a state: the suites, and the cluster's `deploy` and `checks`, that flai runs when a story enters the state.
- `verify_at` at the top level: the suites the merge queue's `batch`, the release's `publish`, and `ci` run.

Each takes `all`, every suite, or a list of suites. One suite is a list of one, `[static]`. A move by default refuses when its suites have not passed. To let it through with a warning, give `verify` as `{suites: ..., severity: warn}`.

The baseline gives today's behaviour: the move to review needs every suite, at `warn`; the acceptance needs every suite; review runs the cluster's deploy and checks; and the batch and CI run every suite. A project that keeps the move to review and the acceptance to its fast suites, and runs the slow ones at publish:

```yaml
states:
  # ... the baseline's other states
  - {name: review, kind: waiting, limit: true, runs: {kluster: [deploy, checks]}}
transitions:
  # ... the baseline's other moves
  - from: in-progress
    to: review
    rules: [has-task, no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}, kluster-verified]
    verify: {suites: [static, behavior], severity: warn}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked], verify: [static, behavior]}
verify_at: {batch: [static, behavior], publish: [integration, smoke], ci: all}
```

Since S-0807 a story's move and its close-out read `verify` ([The move needs the suites passed at the head](../../users/flai.md#the-move-needs-the-suites-passed-at-the-head)):

- A move by a transition with `verify` runs the rule `verified` at `verify`'s severity. It needs the story's last `flai verify` to have passed each suite `verify` lists at the head of the story's branch. On the baseline the move to review warns until every suite passed there.
- The close-out runs the tiers its paths select of the suites the moves out of the story's state need, every suite when one of them gives `all` or none gives `verify`. With the example above it runs `static` and `behavior` alone. `flai verify S-nnnn --all` runs every suite.
- The acceptance runs no `verified`.

Since S-0810 a state's `runs` decides what `flai serve` runs when a story enters the state, as [run suites and the cluster in a state](#run-suites-and-the-cluster-in-a-state) says.

Nothing else reads the keys yet but a state's `runs` since S-0810: the acceptance's gates, the batch, the publish, and CI run as they did; later releases make each read them.

The keys need a flai that knows them. The first publish from the flai release that carries them raises `flai.minimum` in `system-flow.yaml`, and from then an older flai stops before it reads the file, naming the release to install ([Keeping the host's flai current](../index.md#keeping-the-hosts-flai-current)). Before that publish, an older flai refuses a file that uses them with its parse error, an unknown field.

## Run suites and the cluster in a state

A state's `runs` is what `flai serve` runs each time a story enters the state, with no agent ([what a state's runs starts](../../../design/system/workflow.md#what-a-states-runs-starts-s-0810)). It runs one job per entry. The job rebases the story's branch onto main, runs each suite the state names, then the cluster's steps it names. A story entering a state whose `runs` names `deploy` is deployed there into its kluster namespace, and its checks run there when `runs` names `checks` ([kluster.md](../../../design/system/kluster.md#what-happens-when)). No other state deploys it. The baseline names the cluster's steps on review alone, so on the baseline a story is deployed when it enters review.

An `integration` state between `review` and `done`, which runs the `integration` suite and deploys and checks the story on the cluster before you accept it:

```yaml
states:
  # ... backlog, ready, and in-progress as in the baseline
  - {name: review, kind: waiting, limit: true}
  - {name: integration, kind: waiting, runs: {suites: [integration], kluster: [deploy, checks]}}
  # ... done and cancelled as in the baseline
transitions:
  # ... the baseline's other moves, less its review -> done for epics and stories
  - {from: review, to: integration}
  - {from: integration, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked], verify: all}
  - {from: integration, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: integration, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
```

Here `review` gives no `runs`, so a story entering it runs nothing and is not deployed. The acceptance now leaves from `integration`, so the acceptance and the merge queue take stories from there.

- **A pass** leaves the story in the state.
- **A failure** sends the story back by the move with the `send-back` effect out of the state, here `integration` to `in-progress`. A rebase that stops on conflicts or a suite that fails opens a thread `Acceptance sent back: S-nnnn's job on entering integration: ...`. A failed deploy or checks run opens the kluster's own thread. `flai serve` starts the story's agent on the thread, and runs the job again when the story is back in the state. A state with no `send-back` move out of it keeps the story where it is, and the job records why.
- **A cluster step that cannot run** is skipped, and the job records why: no `kluster:` block, no `deploy`, the `kluster` host action off, or, in a state of the `active` kind, no agent namespace free. A skip never holds the story.
- **The `kluster-verified` rule** is checked by the move whose transition names it, at the severity it gives. The baseline names it at `warn` on `in-progress` to `review`. Give it `{rule: kluster-verified, severity: refuse}` to stop a move whose namespace is not deployed and checked at the head ([Add a rule to a move](#add-a-rule-to-a-move)).

`flai serve` keeps the job's record and log under `.flai-cache/verify/jobs/`, as `S-nnnn.state.json` and `S-nnnn.state.log` ([the close-out job](../../../design/system/flai-cli.md#the-close-out-job)).

## Before the story is accepted

1. In the story's worktree, run `flai check --strict`. It checks the branch's file, and the items against it. The codes are `workflow.file`, `workflow.rule`, `workflow.timer`, `workflow.state`, `workflow.unreachable`, `workflow.dead-end`, `workflow.no-pull`, `workflow.no-accept`, `workflow.item`, `workflow.suite`, and `workflow.suite-unnamed`; what each means and what clears it is under [Check the repository](../../users/flai.md#check-the-repository).
2. Fix every error. While the file has one, but `workflow.suite`, every command run in that checkout and `flai serve` refuse to open the project, and the MCP server refuses each tool call, naming the finding:

   ```text
   open the project at <root>: its workflow file cannot be run: <path>:<line>: <code>: <message>; fix the file, or remove it to run the baseline workflow flai carries
   ```

   Only `flai check` still runs. A warning, `workflow.unreachable`, stops nothing. Neither does `workflow.item`, an item naming a state the file lacks: add that name to an `aliases`.

   The two suite findings stop nothing either; `flai check` reports them, and `--strict` fails on them. `workflow.suite`, an error, is a suite the file names that no tier in `system-flow.yaml`'s `tests` has: name a suite a tier gives, add `suite: <name>` to the tier it belongs to, or give `all`. `workflow.suite-unnamed`, a warning, is a suite with a tier that is not `manual` which nothing in the file names, so no step runs it: name it in a `verify`, a `runs`, or `verify_at`, give `all` somewhere, or mark its tiers `manual: true`. A file that gives none of `verify`, `runs`, and `verify_at` reads as the baseline's and gets no `workflow.suite-unnamed`.
3. Look at the board in the worktree: `flai board` there shows the branch's columns. A move made there follows the branch's file, but writes the item in the main checkout, where the old board still runs: an item moved into a new state is a `workflow.item` finding on main until the story is accepted. To try moves freely, copy the file into a scratch project made with `flai new`.
4. Commit, and move the story to review. Accept it yourself. A story that changes the file is the operator's to accept, never the orchestrator's or the merge queue's (ADR-0301 §7). The dashboard's saves refuse the file on main, and take it only on a story's branch.
5. Once it is accepted, main runs the new board: the main checkout's commands, `flai serve` at its next look, the MCP server at its next call (it reads the file again when it changed), and the dashboard through `flai serve`. A story's worktree runs its own branch's file until it takes main in.

## If main's file cannot be run

Every command on main refuses, and `flai serve` logs the finding once. Run `flai check` in the main checkout for the finding. Revert the commit that brought it (`git revert <commit>`), fix the field the finding names, or remove the file to run the baseline, and commit. Then make the change again in a story.

## When the template's baseline changes

`flai upgrade` adds the file only to a project that lacks it, and never writes it after. When the template's baseline changed since the last upgrade, `flai upgrade` lists the file as `keep` and says so, as `--dry-run` does. The lock then records the new baseline, so the next upgrade is silent until it changes again.

```text
the template's baseline design/workflows/board.md changed: the project's file is kept, carry the change over to it in a story
```

Compare the template's file with yours, and carry over what you want in a story, as above.

The template's file gained `verify`, `runs`, and `verify_at` in S-0806. A file you kept from before gives none of them, and flai reads it as giving the baseline's: every suite at `warn` on the move to review, every suite on the acceptance, the cluster's deploy and checks on each state an acceptance leaves from, which is review in the baseline, and every suite at the batch and in CI. It runs what it ran before, and `flai check` warns about no suite. Once you give any one of the keys, flai reads the file as it is written, and a key you leave out names nothing, so carry all three over together ([Name the suites a step runs](#name-the-suites-a-step-runs)).
