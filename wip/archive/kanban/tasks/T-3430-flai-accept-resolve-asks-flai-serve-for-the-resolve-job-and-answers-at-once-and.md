---
id: T-3430
type: task
nature: feature
title: flai accept --resolve asks flai serve for the resolve job and answers at once, and the host method resolve.status gives the job to the dashboard
status: done
parent: S-0804
owner: alex
created: 2026-11-16T23:19:37Z
updated: 2026-11-17T01:18:53Z
transitions:
  - to: ready
    at: 2026-11-16T23:20:38Z
    by: agent-S-0804
  - to: in-progress
    at: 2026-11-17T00:25:41Z
    by: agent-S-0804
  - to: done
    at: 2026-11-17T00:43:35Z
    by: agent-S-0804
stream: S-0804
tags: []
touches: [flai/cmd/accept_resolve.go, flai/cmd/accept_resolve_test.go, flai/cmd/accept.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/resolve_status.go, flai/internal/hostapi/resolve_status_test.go, docs/users/flai-reference.md, flai/internal/guard/guard_tables.go, flai/internal/hostapi/guardreads_test.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, docs/operators/settings.md]
after: [T-3429]
usage:
  source: log
  seconds: 1074
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 89
      output: 37404
      cache_read: 6748122
      cache_write: 253567
      cost: 3.5128
---
# T-3430 flai accept --resolve asks flai serve for the resolve job and answers at once, and the host method resolve.status gives the job to the dashboard

## Work

Waits for T-3429: it asks for the job that task runs, and changes `flai/cmd/accept.go`, which that task changes too.

1. `flai accept S-nnnn --resolve` keeps its refusals of ADR-0302 (the acceptance's flags, `--yes`, the orchestrator, a story not in review, a story nothing holds), then asks flai serve for the resolve job with `closeout.Ask` and a `Request` of kind `resolve`, `By` the operator, and answers at once. It moves nothing and opens no thread. With `--json` it prints `{id, job, started, running, log, reason}`; `reason` says that the job rebases the branch and closes it out in flai serve, and that the story stays in review unless the rebase stops on conflicts or a step fails. A request flai serve does not take (`closeout.ErrServeNotTaking`) or refuses is refused with exit 4 and `{"refused": {"id", "reason"}}`, naming flai serve. Remove `sendBackBlocked`, `resolveText`, and whatever else only they used.
2. Update `flai accept`'s long help and the `--resolve` flag's text for the job, and regenerate `docs/users/flai-reference.md` and `docs/operators/settings.md` with `make flai-reference`.
3. `accept.resolve` keeps running `flai accept <id> --resolve --by=<owner> --json`; its answer is the new object.
4. A host method `resolve.status` with `{id}`, a read gated by no host action, answers the story's last close-out or resolve job from `closeout.Status`: `{id, job}` with the job's id, kind, by, started, ended, verdict, stopped_at, outcome, thread, conflicts, reason, running, step, tier, steps, and log, or `{id}` with no job. Register it as `verify.status` is registered, and add it to every list of methods a test checks.
5. Tests: the request written and the answer printed (with a fake flai serve answering the request file, as the verify `--job` tests do), each refusal, and `resolve.status` over a fixture job and log.

## Done when

- `flai test` passes on the changed paths, the settings-index and reference tiers included.
- `flai accept --resolve` changes no work item and no thread.

## Notes
