---
id: T-1062
type: task
nature: remediation
title: An ADR and the template design say which template ref flai new, flai import, and flai upgrade take, and when upgrade asks
status: done
parent: S-0301
owner: alex
created: 2026-10-06T22:50:14Z
updated: 2026-10-06T23:00:10Z
transitions:
  - to: ready
    at: 2026-10-06T22:54:35Z
    by: agent-S-0301
  - to: in-progress
    at: 2026-10-06T22:54:36Z
    by: agent-S-0301
  - to: done
    at: 2026-10-06T23:00:10Z
    by: agent-S-0301
stream: S-0301
tags: [design, template]
touches: [design/adrs, design/adrs/README.md, design/system/template.md, design/system/project-manifest.md]
usage:
  source: log
  seconds: 333
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 267
      cache_read: 1683293
      cache_write: 70014
      cost: 0.7722
---
# T-1062 An ADR and the template design say which template ref flai new, flai import, and flai upgrade take, and when upgrade asks

## Work

The story changes what `flai upgrade` defaults to and makes it ask a question. Today `design/system/template.md` under "Upgrading a project" says that upgrade takes the manifest's `template.ref` and never prompts. Record the decision before the code tasks build on it.

- Write an ADR with `flai adr new`. Its number is not known yet, which is why the story keeps the folder touch `design/adrs`. It records the following:
  - With no `--ref`, `flai new`, `flai import` and `flai upgrade` take the newest version tag of a git template when the configured or recorded ref follows releases: it is empty, the template's default branch (`main`, which every config and manifest so far holds), or a version tag. Another branch or a commit is used as given. A template with no version tags falls back to the ref, fetched fresh. The config's default `template.ref` stays `main`.
  - `--ref` always wins and is written to `template.ref` in `system-flow.yaml` and to the lock.
  - The lock is the record of what was applied, never a target. The version a project is at is the lock's.
  - `flai upgrade` with no `--ref` asks the operator which version to apply when `system-flow.yaml` names a different ref or version than the lock recorded, because the operator edited it, and that differs from the newest tag. The choices are the version the manifest names, the newest tag, and changing nothing. A recorded tag equal to the lock's is not a pin: the next upgrade takes the newest tag without asking (criterion 3).
  - Without a terminal it changes nothing. It names the choices and the `--ref` to pass.
  - A branch clone in the cache is fetched again before it is used.
- Add the ADR to `design/adrs/README.md`.
- Rewrite "Upgrading a project" in `design/system/template.md`, the `ref` line of the manifest example in `design/system/project-manifest.md`, and the description beside it.
- This task waits for nothing and shares no path with the two code tasks of this layer. The second-layer tasks wait for it.

## Done when

- [ ] The ADR is proposed. Its decision covers the default ref, `--ref` precedence, the prompt and its choices, the run without a terminal, and the cache refetch.
- [ ] `design/system/template.md` and `design/system/project-manifest.md` say the same as the ADR and link it.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner. ADR-0015 is accepted and stays as it is; the new ADR refines what it says about the ref.
