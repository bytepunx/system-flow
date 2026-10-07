---
title: Project manifest
updated: 2026-10-07
status: active
topics: [cli, template]
---

# Project manifest

A conforming repo has `system-flow.yaml` at its root. It is how `flai` and `flaiover` recognise a project, find its folders, and know which template it came from.

```yaml
# system-flow.yaml
version: 1                                   # manifest schema version
name: system-flow                            # project name, used in titles
key: sf                                      # short key, used in generated IDs when a project opts into prefixed IDs
description: Agentic lean project management system
template:
  repo: https://github.com/bytepunx/system-flow-template
  ref: main                                  # the ref last applied: branch, tag, or commit; main or a version tag follows releases (ADR-0103)
  version: 0.1.0                             # template version last applied, from template.yaml
  applied: 2026-09-15T16:00:00Z
layout:                                      # folder names, defaults shown, renameable at import
  design: design
  docs: docs
  wip: wip
projects:                                    # releasable components at the repo root
  - name: flai
    path: flai
    kind: go
    tags: [cli]                              # story tags that mean "delivers to flai"
  - name: flaiover
    path: flaiover
    kind: sveltekit
    tags: [dashboard]
  - name: template
    path: template
    kind: template                           # released by bumping template.yaml and CHANGELOG.md, not by tag
    tags: [template, conventions]
dashboard:
  image: ghcr.io/bytepunx/flaiover
  tag: latest
  port: 4242
  autocommit: true                           # optional, default true: commit documents saved from the dashboard (ADR-0023)
  notify_url: ""                             # optional, default unset: POST new inbox entries here as JSON (S-0042)
tests:                                       # optional (S-0273): the test tiers, cheapest first, that flai test runs; default one plain tier, scripts/test.sh
  - name: unit                               # unique among the tiers
    command: [go, test, -json, "{packages}"] # an argument list; {packages} or {files}, an argument of its own, stands for what the paths selected
    dir: flai                                # optional: the folder it runs in, from the root; default the root
    paths: ["flai/**", "!flai/testdata/**"]  # globs from the root that select it, as claims.shared; ! takes paths out; required unless all_only
    format: go-test-json                     # go-test-json, vitest-json, golangci-json, gofmt-list, or plain; default plain
    all_command: [go, test, -json, ./...]    # optional: what runs instead under --all, with no placeholder
  - name: smoke
    command: [scripts/smoke.sh]
    all_only: true                           # runs only under --all
agent:                                       # optional (S-0103): the agent every new story gets; flai agent set writes it
  harness: claude-code                       # what runs the agent
  model: claude-opus-5-5
  config:                                    # optional: options for the harness, passed on as they are
    effort: high
  roles:                                     # optional (S-0189): the sub-agents' agents, by role
    explore:
      model: haiku
    verify:
      model: sonnet
prime:                                       # optional (S-0146): how flai prime --story builds context packs
  budget: 80KB                               # the size a story's pack fits; default 80KB
issues:                                      # optional (S-0198): how flai check treats open issues
  story_after: 168h                          # how long an issue may stay open with no open story linking it; default 168h, 0 turns the warning off
planning:                                    # optional (S-0199): what items' planning data is counted in
  currency: USD                              # ISO 4217 code of every amount; default USD
  hour_rate: 95                              # what an hour of work costs, in currency; default unset, meaning unknown
  cycle: 168h                                # the period time lost is counted over; default 168h
  default_duration: 1h                       # optional (S-0210): a story's forecast duration with too little history; default 1h
  agent:                                     # optional (S-0208): the planner's agent, over agent above
    model: claude-sonnet-5
    config:
      effort: medium
  replan: deterministic                      # optional (S-0211): what flai serve does when work ahead completes or the order changes: never, deterministic, or agent; default deterministic
  schedule: daily                            # optional (S-0211): when flai serve plans the ready column, a cron expression in UTC or daily; default none
orchestration:                               # optional (S-0217): how the ready column is ordered and when a release is due
  policy: wsjf                               # cod, wsjf, throughput, or fifo, as flai order --by; default fifo
  release:                                   # when accepted stories not yet released are due a release
    policy: threshold                        # judgement, threshold, or theme; default judgement
    value: 500                               # threshold: the unreleased cost of delay per week, in planning.currency
    count: 5                                 # threshold: the accepted stories not yet released
    whole_epics: true                        # optional (S-0222): hold a batch back while a story's epic is in neither review nor done; default off
  permissions:                               # optional (S-0218): what the orchestrator may do without the operator; each off when unset
    promote_to_ready: true                   # also plan_backlog_epics, finalize_drafts, order_ready, accept_reviews, publish
    answer_threads: recommend                # off, recommend, or autonomous; default off
  agent:                                     # optional (S-0218): the orchestrator's agent, over agent above
    model: claude-sonnet-5
analysis:                                    # optional (S-0223): the analyzer's agent and schedule
  agent:                                     # the analyzer's agent, over agent above
    model: claude-sonnet-5
  schedule: "0 6 * * 1"                      # when flai serve runs the analyzer, a cron expression in UTC or daily; default none
claims:                                      # optional (S-0295): how stories' touches claim paths
  shared:                                    # glob patterns of paths whose overlaps hold no story; flai shared edits them
    - docs/users/flai.md
    - design/adrs
    - design/issues
flai:                                        # optional (S-0181): what the project asks of the flai that reads it
  minimum: 1.27.0                            # the oldest flai release that may read it; publishing a flai release that changes the front-matter fields raises it
```

Rules:

- `flai` refuses to run project commands in a directory tree with no `system-flow.yaml` above the current directory, except `flai new` and `flai import`.
- `layout` is the only place folder names live. Everything else resolves through it. Subfolders such as `design/conventions` are fixed names under their layout folder.
- `template.ref` and `template.version` are the ref and version last applied ([ADR-0103](../adrs/0103-flai-new-import-and-upgrade-follow-the-template-s-releases-the-newest-version.md), [template.md](template.md#which-version-is-applied)). A `ref` that is empty, `main`, or a version tag follows releases: with no `--ref`, `flai upgrade` applies the template's newest version tag. Another branch or a commit is used as given. `flai upgrade --ref` writes its ref here. Editing `ref` or `version` here, so that it differs from the lock and from the newest tag, makes the next `flai upgrade` ask which version to apply: the one named here, the newest, or none; a changed `version` names the tag `v<version>`. The version a project is at is the lock's, not this one, unless there is no lock.
- `system-flow.lock.yaml` beside the manifest records the hash of every rendered file so upgrade can tell project edits from baseline (ADR-0015), the `topics` the template gave each marker file so upgrade keeps a project's own (S-0134), and under `vars` the value of every template variable the project was last rendered with, so upgrade renders a fork's own variables again (S-0185).
- `name`, `key`, `description`, `owner`, and `repo` are read back by `flai upgrade` as the template variables `project_name`, `project_key`, `description`, `owner`, and `repo_url`, ahead of the values the lock recorded, so an edit to them here reaches the next upgrade. `flai upgrade --var` refuses them and names the key to edit ([template.md](template.md#upgrading-a-project)).
- `projects` are the components `flai release` versions, one item at a time or, since S-0087, batched by `flai release --pending`: code kinds get `<name>/vX.Y.Z` tags, kind `template` gets its version file bumped. `flai accept` versions nothing. `tags` are aliases a story or epic tag may use to say which component it delivers to.
- `dashboard.notify_url`, when set, makes the dashboard's server POST `{ project, entry: { key, kind, title, href, at } }` to that URL for each inbox entry that appears after it started: one attempt, a short timeout, a warning in the log on failure. The token and file contents are never sent. Unset by default.
- `dashboard.autocommit: false` leaves documents saved from the dashboard uncommitted; the default commits each save on the main checkout, one path per commit (ADR-0023).
- `tests` lists the project's test tiers, cheapest first, which `flai test` runs for the paths each selects (S-0273), and `flai verify` runs for what a story's branch changed against the main branch (S-0270). A tier has:
  - `name`, unique among the tiers.
  - `command`, an argument list run as it stands, never through a shell. `{packages}` and `{files}`, each as an argument of its own, stand for the packages and the files the tier's paths selected, and `flai test` fills them. Any other argument that is a word in braces, such as `{story}`, and a placeholder inside a longer argument are refused.
  - `dir`, the folder the command runs in, relative to the root and inside it. Unset, it is the root.
  - `paths`, glob patterns relative to the root with the syntax of `claims.shared` (below); one beginning with `!` takes paths out again. Required, with at least one pattern without `!`, unless `all_only` is set.
  - `format`, how the command's output becomes findings: `go-test-json`, `vitest-json`, `golangci-json`, `gofmt-list`, or `plain`, which reads only the exit status. Unset, it is `plain`. The list is `manifest.TestFormats`.
  - `all_only`, set for a tier that runs only under `flai test --all`, such as integration and smoke tests, and `all_command`, an argument list with no placeholder that runs instead of `command` under `--all`.

  `flai test` without `--all` never runs an `all_only` tier, whatever its `paths`. `flai verify S-nnnn`, which the close-out runs, runs the tiers `flai test` selects for the paths the story's branch changed, and then each `all_only` tier whose `paths` select one of those paths, or that has no `paths`, as it runs under `--all` (`verify.SelectStory`). So `paths` on an `all_only` tier says which changes need it before review: this repository gives `integration` and `smoke` `flai/**`, `template` `template/**`, and `flaiover` `flaiover/**`, and a story that changes only documents runs none of them. An `all_only` tier with no `paths` runs for every story.

  Unset, the project has one `plain` tier named `test` that runs `scripts/test.sh` for every path when that file exists, and none when it does not (`Manifest.TestTiers`), so a project made before the key still has a tier. `tests: []` means none. A project made from the template starts with `scripts/lint-md.sh` as `markdown`, for its markdown files (S-0270), its `scripts/test.sh` as `test`, and `scripts/integration.sh` and `scripts/smoke.sh` as `integration` and `smoke`, both `all_only` with no `paths`, so `flai verify` runs them for every story; all four are `plain`. `flai check` reports each thing wrong with a tier as a `manifest.tests` error, by its index, its field, and its name, such as `tests[1].format (tier "unit") "junit" is not a format flai reads; write …` (`Manifest.TestErrors`); `TestTiers` refuses such a list, so `flai test` runs none of it. `flai manifest set tests=<JSON>` writes the list whole (below).
- `agent` is the project's default agent (S-0103, [ADR-0037](../adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)): a harness, a model, and `config`, a flat map of options for the harness. A story made while it is set gets a copy in its front matter, which the story may change. A change to the default reaches new stories only. `flai agent` shows it, `flai agent set --harness --model --config key=value --unset key` changes only what is given, and `flai agent clear` removes it. Each rewrites the `agent:` block alone. A harness is a lower-case name; a model may also hold `.`, `:`, `/`, and `@`; config keys are lower-case words and values are one line. `flai check` reports a default that is not (`manifest.agent`). `roles` (S-0189) maps a sub-agent role (`explore`, `verify`) to its own harness, model, and config, checked the same way; a story gets them merged role by role, and `flai agent set --role-harness role=h --role-model role=m --role-config role.key=v --unset-role role` changes them.
- `prime.budget` is the size a story's context pack fits, `flai prime --story` and the MCP `prime` tool alike ([ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md), S-0146): bytes, or a number with `KB` or `MB` (1024-based). `--budget` overrides it for one run. It is the project's, not the host's, because a pack is the same for every agent that works the project. Unset, it is 80 KB.
- `issues.story_after` is how long an issue may stay open with no open story linking it before `flai check` warns with `issues.no-story` (S-0198). It is a Go duration, such as `168h` or `24h`. Unset, it is 168h, seven days; `0` turns the warning off. A value that is not a duration, or is negative, is a `manifest.issues` error, and the warning is off until it is fixed. A story links an issue when its body names the issue's ID, as `flai issue story` writes it.
- `planning` sets what the planning data on work items is counted in ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md), S-0199; the fields are in [work-hierarchy.md](work-hierarchy.md#identifiers)). `planning.currency` is the ISO 4217 code, three capital letters, of every amount in an item's `cost_of_delay` and of `hour_rate`; unset, it is USD. `planning.hour_rate` is what an hour of work costs, a number of zero or more, which prices a cost of delay's time lost; unset means unknown, not free. `planning.cycle` is the period `time_lost_per_cycle` is counted over, a Go duration longer than zero; unset, it is 168h, a week. `planning.default_duration` is the duration `flai forecast` gives a story when fewer than three done stories with usage match it on any rung (S-0210, [strategic-agents.md](strategic-agents.md#forecast)), a Go duration longer than zero; unset, it is 1h. `flai check` reports a value that is none of these as a `manifest.planning` error. `flai show` and MCP's `item_get` give an item's amounts in the currency.
- `planning.agent` is the planner's agent (S-0208, [ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)): the same shape as `agent`, a harness, a model, `config`, and `roles`, merged over `agent` field by field, config key by config key, and role by role (`Manifest.PlanningAgent`), so it names only what the planner runs differently. Unset, the planner runs on `agent`. `flai check` reports a value that is not valid, under the name `planning.agent`, as a `manifest.planning` error. `flai manifest set` writes it whole, and the Planner page's settings panel through it (S-0229, below); `flai agent` does not change it. How the planner runs is in [strategic-agents.md](strategic-agents.md#its-agent).
- `planning.replan` and `planning.schedule` say when `flai serve` plans again on its own (S-0211, [ADR-0084](../adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)). Both act only while the `plan` host action is on. `planning.replan` is what it does when a story is accepted or cancelled or the pull order changes: `never` does nothing; `deterministic` plays the board out again with no agent and moves the delivery of each forecast that changed; `agent` does that and queues the planner for each story whose delivery moved. Unset, it is `deterministic`. `planning.schedule` is when it runs the planner over every ready story: a five-field cron expression in UTC, such as `0 6 * * 1-5`, or `daily`, which is 00:00 UTC. Unset, there is no schedule. `flai check` reports a value that is none of these as a `manifest.planning` error. The operator sets both with `flai manifest set`, on the Planner page, or by hand, as the other `planning` keys but `currency`. When the planner runs again is in [strategic-agents.md](strategic-agents.md#planning-again).
- `orchestration.policy` orders the ready column (S-0217). Its values are the names `flai order --by` takes: `cod` orders by cost of delay, the largest first; `wsjf` by cost of delay over forecast duration, the largest first; `throughput` by forecast duration, the shortest first; `fifo` leaves the operator's order alone. Unset, it is `fifo`. The list is `manifest.OrderPolicies`, and `Orchestration.PolicyOrDefault` gives the policy in effect.
- `orchestration.release` says when accepted stories not yet released are due a release (S-0217). `orchestration.release.policy` is `judgement`, `threshold`, or `theme`. Unset, it is `judgement`, which leaves the release to the operator and is never met by itself. `threshold` takes `value`, the unreleased cost of delay per week in `planning.currency`, and `count`, the number of accepted stories not yet released; either or both, each zero or more. `theme` takes `epic`, an epic ID such as `E-0001`, or `tag`, a tag; one of them, not both. The list is `manifest.ReleasePolicies`, and `Release.PolicyOrDefault` gives the policy in effect. `orchestration.release.whole_epics` is off when unset (S-0222). Set, it holds a batch back under every policy while a story accepted and not yet released belongs to an epic in neither `review` nor `done`. Such a story is in `flai release --evaluate`'s `held_by_epic`, with its epic and the epic's status, and while one is there the policy is not met. Under `judgement` the list is what the orchestrator must not publish. A story with no epic is never held; one whose epic is not found is.
- `orchestration.permissions` is what the orchestrator may do without the operator (S-0218, [ADR-0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md)). Each permission is off when unset. `flai guard` reads them from the manifest at each of the orchestrator's calls, so a change applies at its next call ([strategic-agents.md](strategic-agents.md#its-permissions-and-the-guard)). The list is `manifest.PermissionNames`, and `Permissions.Allows` says whether one is on.

  | Key | Values | Default | Lets the orchestrator |
  |-----|--------|---------|-----------------------|
  | `plan_backlog_epics` | `true`, `false` | `false` | ask for the planner on an epic in the backlog |
  | `finalize_drafts` | `true`, `false` | `false` | finalize a draft story, `flai edit --no-draft` |
  | `promote_to_ready` | `true`, `false` | `false` | move a story to `ready` |
  | `order_ready` | `true`, `false` | `false` | write the order of the ready column, `flai order` |
  | `answer_threads` | `off`, `recommend`, `autonomous` | `off` | reply on threads: `recommend` with a recommendation for the operator, `autonomous` with an answer of its own; the guard lets it reply with either |
  | `accept_reviews` | `true`, `false` | `false` | accept a story in review, `flai accept` |
  | `publish` | `true`, `false` | `false` | release and push accepted work, `flai release --pending` and `flai push` |

- `orchestration.agent` is the orchestrator's agent (S-0218): the same shape as `agent`, merged over it as `planning.agent` is (`Manifest.OrchestrationAgent`), so it names only what the orchestrator runs differently. Unset, the orchestrator runs on `agent`. How it runs is in [strategic-agents.md](strategic-agents.md#the-orchestrator).
- `flai check` reports each thing wrong with `orchestration` as a `manifest.orchestration` error on its line: a policy or release policy outside its list, a threshold with neither figure or a negative one, a theme with neither or both of `epic` and `tag`, a key under `permissions` that names no permission, an `answer_threads` that is not `off`, `recommend`, or `autonomous`, and an `agent` that is not valid, under the name `orchestration.agent`. `flai manifest set` writes every `orchestration` key but `agent`, and the Orchestrator page's settings panel through it (S-0229, below); `orchestration.agent` is set by hand.
- `analysis.agent` is the analyzer's agent (S-0223, [ADR-0099](../adrs/0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md)): the same shape as `agent`, merged over it as `planning.agent` is (`Manifest.AnalysisAgent`), so it names only what the analyzer runs differently. Unset, the analyzer runs on `agent`. How it runs is in [strategic-agents.md](strategic-agents.md#the-analyzer).
- `analysis.schedule` is when `flai serve` runs the analyzer on its own, with no focus, while the `analyze` host action is on: a five-field cron expression in UTC, such as `0 6 * * 1`, or `daily`, which is 00:00 UTC, parsed as `planning.schedule` is. Unset, there is no schedule, and the analyzer runs only when the operator asks for it ([The analysis schedule](strategic-agents.md#the-analysis-schedule)).
- `flai check` reports each thing wrong with `analysis` as a `manifest.analysis` error on its line: a `schedule` it cannot parse, or one that never comes round, and an `agent` that is not valid, under the name `analysis.agent`. `flai manifest set` writes both keys, and the Analyzer page's settings panel through it (S-0229, below).
- `flai manifest set <key>=<value>... [--unset <key>]... [--autocommit] [--trailer ...]` writes the strategic agents' settings (S-0229, [ADR-0101](../adrs/0101-with-the-settings-host-action-on-the-dashboard-edits-the-strategic-agents.md)) and the project's test tiers (S-0273): the keys of its catalog, `manifest.Settings`, below. The dashboard's settings panels on the Orchestrator, Planner, and Analyzer pages run it through the host API's `settings.manifest`, and `settings.get` gives each key's value with its kind, values, default, meaning, and a permission's risk from the same catalog ([flaiover-dashboard.md](flaiover-dashboard.md#strategic-agents-settings-s-0229)). The meaning and the risk below are the sentences the panels show.

  | Key | Kind | Default | Meaning |
  |-----|------|---------|---------|
  | `orchestration.permissions.plan_backlog_epics` | boolean | `false` | Lets the orchestrator ask the planner to draft stories for an epic in the backlog. Risk: the planner runs, and spends, on an epic you may not mean to start yet, and its drafts fill the backlog |
  | `orchestration.permissions.finalize_drafts` | boolean | `false` | Lets the orchestrator finalize a draft story, as `flai edit --no-draft` does. Risk: a story the planner drafted becomes one that can be promoted without your having read it |
  | `orchestration.permissions.promote_to_ready` | boolean | `false` | Lets the orchestrator move a story to ready. Risk: agents may pull, work, and spend on a story you have not chosen to start |
  | `orchestration.permissions.order_ready` | boolean | `false` | Lets the orchestrator write the order of the ready column by `orchestration.policy`. Risk: it replaces the order you set on the board, so the next story pulled is the policy's choice, not yours |
  | `orchestration.permissions.answer_threads` | choice: `off`, `recommend`, `autonomous` | `off` | How the orchestrator answers threads: off leaves them to you, recommend replies with a recommendation for you, and autonomous answers them itself. Risk: with autonomous, agents act on answers you did not give; with recommend, the decision stays yours |
  | `orchestration.permissions.accept_reviews` | boolean | `false` | Lets the orchestrator accept a story in review, as `flai accept` does. Risk: work is merged to the main branch without your review, and a story accepted wrongly has to be undone by hand |
  | `orchestration.permissions.publish` | boolean | `false` | Lets the orchestrator release and push accepted work, as `flai release --pending` and `flai push` do. Risk: a release reaches everyone who pulls or installs the project, and a published release cannot be taken back |
  | `orchestration.policy` | choice: `cod`, `wsjf`, `throughput`, `fifo` | `fifo` | How the ready column is ordered: cod by cost of delay, wsjf by cost of delay over forecast duration, throughput shortest first, and fifo leaves your order alone |
  | `orchestration.release.policy` | choice: `judgement`, `threshold`, `theme` | `judgement` | When accepted stories not yet released are due a release: judgement leaves it to you, threshold at a value or a count of unreleased work, and theme when every story of an epic or a tag is accepted |
  | `orchestration.release.value` | number | unset | Under threshold, the unreleased cost of delay per week, in the project's currency, at which a release is due |
  | `orchestration.release.count` | number, whole | unset | Under threshold, the number of accepted stories not yet released at which a release is due |
  | `orchestration.release.epic` | text | unset | Under theme, the epic, such as E-0001, whose stories, every one accepted, make a release due |
  | `orchestration.release.tag` | text | unset | Under theme, the tag whose stories, every one accepted, make a release due |
  | `orchestration.release.whole_epics` | boolean | `false` | Under every policy, holds a release back while a story accepted and not yet released belongs to an epic in neither review nor done |
  | `planning.agent` | agent | the project's `agent` | The planner's agent, over the project's: what it sets wins, and what it leaves out is the project's agent's |
  | `planning.replan` | choice: `never`, `deterministic`, `agent` | `deterministic` | What flai serve does when a story is accepted or cancelled or the pull order changes: never nothing, deterministic plays the board out again and moves forecast deliveries, and agent does that and queues the planner for each story whose delivery moved |
  | `planning.schedule` | cron | unset | When flai serve runs the planner over the ready column, a five-field cron expression in UTC or daily; unset, there is no schedule |
  | `planning.hour_rate` | number | unset | What an hour of work costs, in the project's currency, which prices a cost of delay's time lost; unset means unknown, not free |
  | `planning.cycle` | duration | `168h` | The period a cost of delay's time lost is counted over |
  | `planning.default_duration` | duration | `1h` | The work a story is forecast to take when there is too little history to forecast it from |
  | `analysis.agent` | agent | the project's `agent` | The analyzer's agent, over the project's: what it sets wins, and what it leaves out is the project's agent's |
  | `analysis.schedule` | cron | unset | When flai serve runs the analyzer, a five-field cron expression in UTC or daily; unset, it runs only when you ask |
  | `tests` | tests | one `plain` tier running `scripts/test.sh`, when it exists | The project's test tiers, cheapest first, that flai test runs for the paths each selects: each a name, a command, the paths that select it, and the format of its output; unset, one plain tier runs scripts/test.sh when it exists, and [] means none |

  - **Values.** A boolean is `true` or `false`; a number is written in digits, a count a whole one; a duration is a Go duration such as `168h`; a cron expression has five fields in UTC, or is `daily`; a choice is one of its values; a text is one line; an agent is JSON, such as `{"harness":"claude-code","model":"claude-sonnet-5","config":{"effort":"medium"}}`, written whole as a block, and `{}` unsets it; a list of test tiers is JSON too, such as `[{"name":"test","command":["scripts/test.sh"],"paths":["**"]}]`, written whole as a block, one tier an item, and `[]` writes a list of none. An empty value, a value on more than one line, and a key given twice are refused. `--unset <key>` removes the key, so that its default applies.
  - **Not among them.** `planning.currency` is refused with its reason: changing it re-denominates every amount on the items, which flai does not convert, so it stays a hand edit. So are `orchestration.agent` and every key outside the strategic blocks but `tests`, named in the refusal with the keys it writes.
  - **How it writes.** Each key's line, or an agent's block, is replaced where it is, or added at the end of its block, adding the blocks above it that are missing; every other key, the order, and the comments are kept, a replaced line's comment included. A block written on one line, such as `release: {policy: theme}`, is refused, naming it, rather than rewritten. The result is read back, and a key that does not read back as given is refused, leaving the file as it was.
  - **How it checks.** The manifest is checked whole, as it would be after the change, with the validation `flai check` uses: the project's `agent`, `issues.story_after`, `planning`, `orchestration`, `analysis`, `tests`, and `claims`. Any problem is refused with its field and its reason, such as `orchestration.release.value: -1 is not an amount of zero or more; write …`, and nothing is written; the exit status is 1, as for a workflow rule. So a problem elsewhere in the manifest blocks every change until it is fixed by hand. With `--json` a refusal is `{"refused": [{"field": ..., "reason": ...}]}` and a change `{"set": [{"key": ..., "value": ...}], "unset": [...]}`, with `commit` when it was committed.
  - **Committing.** With `--autocommit`, and unless `dashboard.autocommit` is false, it commits `system-flow.yaml` alone as `chore: change the manifest's settings: <keys>`, with each `--trailer`. `settings.manifest` runs it so, with the dashboard's trailer.
  - **While the orchestrator runs.** A change does not restart it. `flai guard` reads the permissions at each of its calls, so a change holds from its next call, and its prompt tells it to read its permissions, policy, and release policy again before each decision ([strategic-agents.md](strategic-agents.md#its-permissions-and-the-guard)).
- `claims.shared` lists the paths that many stories change in separate sections or new files, as glob patterns (S-0295, [ADR-0096](../adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)). An overlap whose narrower entry lies wholly inside a pattern holds no ready story, is not a `wip.overlap`, and is not told as a grown claim; the trial merge at sync and the notice at acceptance still report it ([workflow.md](workflow.md#branches-and-collisions-adr-0019)).
  - Patterns are relative to the repository root and separated by `/`. `*` matches any characters within one segment, `**` zero or more whole segments, and `?` one character within a segment. A pattern with no glob character matches the path and everything below it, so `design/adrs` and `design/adrs/**` mean the same. A folder entry lies inside a pattern only when the whole folder does: `docs/users` is inside `docs/users/**` but not inside `docs/users/*.md`.
  - `flai check` reports an empty pattern, an absolute path, a `..` segment, or a malformed glob as an error on its line. A pattern that is not valid frees nothing.
  - Unset, the list is empty and every overlap holds. A project made from the template starts with `adrs` and `issues` under its design layout folder, `design/adrs` and `design/issues` by default. This project's list adds `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, and `template/CHANGELOG.md`.
  - `flai shared list`, `check`, `add`, and `remove` read and change it ([flai-cli.md](flai-cli.md#commands)), as do the host API's `settings.get`, `settings.shared_check`, and `settings.shared`, the dashboard's project settings, and the MCP tools `shared_paths` and `shared_paths_edit`. Add and remove rewrite only this key, keeping the file's other keys and comments. Only the operator's own session may change it: `flai guard` refuses every session flai serve starts.
  - A flai older than S-0295 ignores the key, as the manifest is decoded leniently, and holds every overlap. The key is not a front-matter field, so it raises no `flai.minimum`.
- `flai.minimum` is the oldest flai release, `X.Y.Z`, that may read the project (S-0181): one that knows every front-matter field its items, threads, and issues carry. `manifest.Load` refuses the manifest for a flai below it, so every command, `flai serve` (which leaves the project unserved and says why), and `flai mcp` stop before reading any item, with `manifest.TooOldError`: the version needed, the running one, and the upgrade (`flai host upgrade`, or `flai self-upgrade` where no flai host runs). A dev build (`dev`) is never below it, and one that is not a release version is a load error. Publishing a flai release raises it to that release when `flai/internal/workitem/front-matter-fields.txt` changed since the previous `flai/v*` tag (`release.RaiseMinimum`, in the publish commit, which warns that the host's flai must be upgraded once the release's binaries are built), so a release that adds a front-matter field raises it. A fields file the previous tag did not have raises nothing. A flai older than S-0181 does not know the key and ignores it, as the manifest is decoded leniently; what keeps such a flai reading is that it is told it is behind ([flai-cli.md](flai-cli.md#versions-the-hosts-flai-and-the-tree)). Unset, any flai reads the project.
- The manifest is human-edited YAML. `flai` rewrites only the keys it owns (`template.*`, `projects`, `agent` with `flai agent set`, `claims.shared`, `flai.minimum`) and the strategic agents' settings and `tests` with `flai manifest set` (S-0229, S-0273), and preserves comments where the YAML library allows it.
