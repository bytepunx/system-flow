---
title: Project manifest
updated: 2026-10-05
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
  ref: main                                  # branch, tag, or commit
  version: 0.1.0                             # template version applied, from template.yaml
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
flai:                                        # optional (S-0181): what the project asks of the flai that reads it
  minimum: 1.27.0                            # the oldest flai release that may read it; publishing a flai release that changes the front-matter fields raises it
```

Rules:

- `flai` refuses to run project commands in a directory tree with no `system-flow.yaml` above the current directory, except `flai new` and `flai import`.
- `layout` is the only place folder names live. Everything else resolves through it. Subfolders such as `design/conventions` are fixed names under their layout folder.
- `template.version` is the version `flai upgrade` compares against; `system-flow.lock.yaml` beside the manifest records the hash of every rendered file so upgrade can tell project edits from baseline (ADR-0015), the `topics` the template gave each marker file so upgrade keeps a project's own (S-0134), and under `vars` the value of every template variable the project was last rendered with, so upgrade renders a fork's own variables again (S-0185).
- `name`, `key`, `description`, `owner`, and `repo` are read back by `flai upgrade` as the template variables `project_name`, `project_key`, `description`, `owner`, and `repo_url`, ahead of the values the lock recorded, so an edit to them here reaches the next upgrade. `flai upgrade --var` refuses them and names the key to edit ([template.md](template.md#upgrading-a-project)).
- `projects` are the components `flai release` versions, one item at a time or, since S-0087, batched by `flai release --pending`: code kinds get `<name>/vX.Y.Z` tags, kind `template` gets its version file bumped. `flai accept` versions nothing. `tags` are aliases a story or epic tag may use to say which component it delivers to.
- `dashboard.notify_url`, when set, makes the dashboard's server POST `{ project, entry: { key, kind, title, href, at } }` to that URL for each inbox entry that appears after it started: one attempt, a short timeout, a warning in the log on failure. The token and file contents are never sent. Unset by default.
- `dashboard.autocommit: false` leaves documents saved from the dashboard uncommitted; the default commits each save on the main checkout, one path per commit (ADR-0023).
- `agent` is the project's default agent (S-0103, [ADR-0037](../adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)): a harness, a model, and `config`, a flat map of options for the harness. A story made while it is set gets a copy in its front matter, which the story may change. A change to the default reaches new stories only. `flai agent` shows it, `flai agent set --harness --model --config key=value --unset key` changes only what is given, and `flai agent clear` removes it. Each rewrites the `agent:` block alone. A harness is a lower-case name; a model may also hold `.`, `:`, `/`, and `@`; config keys are lower-case words and values are one line. `flai check` reports a default that is not (`manifest.agent`). `roles` (S-0189) maps a sub-agent role (`explore`, `verify`) to its own harness, model, and config, checked the same way; a story gets them merged role by role, and `flai agent set --role-harness role=h --role-model role=m --role-config role.key=v --unset-role role` changes them.
- `prime.budget` is the size a story's context pack fits, `flai prime --story` and the MCP `prime` tool alike ([ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md), S-0146): bytes, or a number with `KB` or `MB` (1024-based). `--budget` overrides it for one run. It is the project's, not the host's, because a pack is the same for every agent that works the project. Unset, it is 80 KB.
- `issues.story_after` is how long an issue may stay open with no open story linking it before `flai check` warns with `issues.no-story` (S-0198). It is a Go duration, such as `168h` or `24h`. Unset, it is 168h, seven days; `0` turns the warning off. A value that is not a duration, or is negative, is a `manifest.issues` error, and the warning is off until it is fixed. A story links an issue when its body names the issue's ID, as `flai issue story` writes it.
- `planning` sets what the planning data on work items is counted in ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md), S-0199; the fields are in [work-hierarchy.md](work-hierarchy.md#identifiers)). `planning.currency` is the ISO 4217 code, three capital letters, of every amount in an item's `cost_of_delay` and of `hour_rate`; unset, it is USD. `planning.hour_rate` is what an hour of work costs, a number of zero or more, which prices a cost of delay's time lost; unset means unknown, not free. `planning.cycle` is the period `time_lost_per_cycle` is counted over, a Go duration longer than zero; unset, it is 168h, a week. `planning.default_duration` is the duration `flai forecast` gives a story when fewer than three done stories with usage match it on any rung (S-0210, [strategic-agents.md](strategic-agents.md#forecast)), a Go duration longer than zero; unset, it is 1h. `flai check` reports a value that is none of these as a `manifest.planning` error. `flai show` and MCP's `item_get` give an item's amounts in the currency.
- `planning.agent` is the planner's agent (S-0208, [ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)): the same shape as `agent`, a harness, a model, `config`, and `roles`, merged over `agent` field by field, config key by config key, and role by role (`Manifest.PlanningAgent`), so it names only what the planner runs differently. Unset, the planner runs on `agent`. `flai check` reports a value that is not valid, under the name `planning.agent`, as a `manifest.planning` error. The operator sets it by hand, as the other `planning` keys; no flai command changes it, `flai agent` included. How the planner runs is in [strategic-agents.md](strategic-agents.md#its-agent).
- `planning.replan` and `planning.schedule` say when `flai serve` plans again on its own (S-0211, [ADR-0084](../adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)). Both act only while the `plan` host action is on. `planning.replan` is what it does when a story is accepted or cancelled or the pull order changes: `never` does nothing; `deterministic` plays the board out again with no agent and moves the delivery of each forecast that changed; `agent` does that and queues the planner for each story whose delivery moved. Unset, it is `deterministic`. `planning.schedule` is when it runs the planner over every ready story: a five-field cron expression in UTC, such as `0 6 * * 1-5`, or `daily`, which is 00:00 UTC. Unset, there is no schedule. `flai check` reports a value that is none of these as a `manifest.planning` error. The operator sets both by hand, as the other `planning` keys. When the planner runs again is in [strategic-agents.md](strategic-agents.md#planning-again).
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
- `flai check` reports each thing wrong with `orchestration` as a `manifest.orchestration` error on its line: a policy or release policy outside its list, a threshold with neither figure or a negative one, a theme with neither or both of `epic` and `tag`, a key under `permissions` that names no permission, an `answer_threads` that is not `off`, `recommend`, or `autonomous`, and an `agent` that is not valid, under the name `orchestration.agent`. The operator sets every `orchestration` key by hand; no flai command changes them.
- `flai.minimum` is the oldest flai release, `X.Y.Z`, that may read the project (S-0181): one that knows every front-matter field its items, threads, and issues carry. `manifest.Load` refuses the manifest for a flai below it, so every command, `flai serve` (which leaves the project unserved and says why), and `flai mcp` stop before reading any item, with `manifest.TooOldError`: the version needed, the running one, and the upgrade (`flai host upgrade`, or `flai self-upgrade` where no flai host runs). A dev build (`dev`) is never below it, and one that is not a release version is a load error. Publishing a flai release raises it to that release when `flai/internal/workitem/front-matter-fields.txt` changed since the previous `flai/v*` tag (`release.RaiseMinimum`, in the publish commit, which warns that the host's flai must be upgraded once the release's binaries are built), so a release that adds a front-matter field raises it. A fields file the previous tag did not have raises nothing. A flai older than S-0181 does not know the key and ignores it, as the manifest is decoded leniently; what keeps such a flai reading is that it is told it is behind ([flai-cli.md](flai-cli.md#versions-the-hosts-flai-and-the-tree)). Unset, any flai reads the project.
- The manifest is human-edited YAML. `flai` rewrites only the keys it owns (`template.*`, `projects`, `agent`, `flai.minimum`) and preserves comments where the YAML library allows it.
