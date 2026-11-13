---
id: T-3280
type: task
nature: feature
title: kluster.md, flai-cli.md, the user guide, and the settings index say how a story's namespace gets its own kubeconfig and how its agent runs exec and log
status: backlog
parent: S-0763
owner: alex
created: 2026-11-13T09:03:41Z
updated: 2026-11-13T09:03:59Z
transitions: []
stream: S-0763
tags: [cli, docs]
touches: [design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/settings.md]
after: [T-3276, T-3277, T-3278, T-3279]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 15
          cache_read: 81708
          cache_write: 2328
          cost: 0.0242
---
# T-3280 kluster.md, flai-cli.md, the user guide, and the settings index say how a story's namespace gets its own kubeconfig and how its agent runs exec and log

## Work

Criterion 6, written from what the four tasks before built, which is why it waits for them all.

- `design/system/kluster.md`: in § What the agent runs against its namespace (S-0660, S-0661, S-0763), replace the bullet on its own commands, not built, with what S-0763 built: the kubeconfig step after the deploy step, its ServiceAccount, RoleBinding, token, file, and mode, that its failure fails the deploy at `deploy`, and that the teardown removes it; `flai kluster exec` and `flai kluster log` and their refusals; and the guard's rules, removing them from the list of what is not built in § What `flai guard` refuses everyone else and § The `kluster` host action. Add the agents' kubeconfig to the table of what flai keeps under `.flai-cache/kluster/`, and say in § The steps and their variables that `exec` gives the variables now and `env` stays S-0661's. Add `exec` and `log` to the table in § Commands (S-0577, S-0578), and say in § The host methods that `kluster.log` reads through the kluster package's reader, which `flai kluster log` shares. Say that an agent in a container reaches neither the kubeconfig nor the logs this way, since `.flai-cache/kluster` is not mounted.
- `design/system/flai-cli.md`: the two commands in § Commands, and the guard's kluster rules where it lists them.
- `docs/users/flai.md`: how a story's agent runs its own commands and reads a step's log against its namespace, and what the guard refuses it.
- `docs/operators/settings.md`: name `flai kluster exec` in the row of the kluster variables, with `KUBECONFIG` the namespace's kubeconfig there and no `FLAI_KLUSTER_IMAGES`. The flag index is the generator's, written by the tasks before.
- Run `flai test` on the four files, and `scripts/flai-reference.sh` once more if the help changed.

## Done when

- `design/system/kluster.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` say how a deploy gives the namespace its kubeconfig, how the teardown removes it, how `flai kluster exec` and `flai kluster log` work, and what the guard lets through and refuses, with nothing of S-0763's left marked not built.
- `docs/operators/settings.md` names `flai kluster exec` among the commands that give the kluster variables.
- `flai test` passes on the changed files, the markdown lint among its tiers.

## Notes
