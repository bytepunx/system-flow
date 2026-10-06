---
id: T-1064
type: task
nature: remediation
title: flai upgrade with no --ref takes the newest template tag and asks which version to apply when the manifest, the lock, and the newest tag disagree
status: done
parent: S-0301
owner: alex
created: 2026-10-06T22:50:26Z
updated: 2026-10-06T23:11:45Z
transitions:
  - to: ready
    at: 2026-10-06T23:00:24Z
    by: agent-S-0301
  - to: in-progress
    at: 2026-10-06T23:00:24Z
    by: agent-S-0301
  - to: done
    at: 2026-10-06T23:11:45Z
    by: agent-S-0301
stream: S-0301
tags: [cli, template]
touches: [flai/cmd/upgrade.go, flai/cmd/upgrade_test.go]
after: [T-1057, T-1060, T-1062]
usage:
  source: log
  seconds: 681
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 86
      output: 36690
      cache_read: 4788882
      cache_write: 127338
      cost: 2.4557
---
# T-1064 flai upgrade with no --ref takes the newest template tag and asks which version to apply when the manifest, the lock, and the newest tag disagree

## Work

Criteria 3 and 4, as the ADR from T-1062 decides.

- In `flai/cmd/upgrade.go`, when `--ref` is empty and the template is a git repository, find the newest version tag with T-1060's `Latest`. Compare three versions: the manifest's `template.version`, the version `system-flow.yaml` names; the version the lock recorded; and the newest tag.
  - When the manifest's ref follows releases (T-1060's `FollowsReleases`) and the operator has not edited it, upgrade to the newest tag with no question, and record it. A recorded tag equal to the lock's is not a pin. A ref naming another branch or a commit is used as given, fetched fresh. A repository with no version tags keeps today's behaviour: the manifest's ref, fetched fresh.
  - The operator edited `system-flow.yaml` when a lock exists and the manifest's `template.ref` or `template.version` differs from the lock's. A changed ref names that ref. A changed version names the tag `v<version>`, refused with the tags listed when there is none. When that intent equals the newest tag, apply it with no question.
  - When it differs from the newest tag, ask with `a.prompts().Select`. The choices are the version the manifest names, the newest tag, and changing nothing, each labelled with its version. The answer is the ref the upgrade applies, and T-1057 records it in `template.ref`. It never falls back to the lock's version.
  - Without a terminal (`prompt.ErrNoAnswer`), or with `--yes`, change nothing. Exit non-zero with a message that names the versions and the `flai upgrade --ref <tag>` to run.
  - The version the project is at is the lock's `template.version` when there is a lock, else the manifest's. "Already at template" compares with it, and `upgrade.Compute` is given it as the version upgraded from. Today a `template.version` hand-set to the newest makes upgrade report "already at" and apply nothing, which is the reported revert.
- `--dry-run` says which version it would apply and whether it would ask.
- Update the command's `Long` help and the `--ref` flag's text (default: the newest version tag). T-1066 regenerates the reference from them.
- This task waits for T-1057, which changes the same lines of `upgrade.go`. It also waits for T-1060, whose `Latest` it calls, and for T-1062, whose ADR it implements. It runs alongside T-1065, which touches only `new.go`, `import.go`, and config.

## Done when

- [ ] Tests in `flai/cmd/upgrade_test.go` use a local bare template repository with tags, and cover each of these:
  - newest tag taken with no question;
  - a manifest version edited by hand while a terminal is attached, where the prompt is shown and each choice is applied and recorded;
  - the same case without a terminal, which refuses and changes nothing;
  - the reported case: lock at `main` 1.0.0, manifest version hand-set to the newest, which upgrades to it without asking and does not report "already at";
  - a manifest whose recorded tag equals the lock's while a newer tag exists, which upgrades without asking;
  - a repository with no tags, which keeps today's behaviour.
- [ ] No path through `flai upgrade` renders the lock's version after the operator chose another.
- [ ] `scripts/flai-test.sh` passes for `flai/cmd`.

## Notes

Drafted by the planner. The prompt package is `flai/internal/prompt` (ADR-0052). The upgrade already uses `Select` for conflicts in `askConflict`.
