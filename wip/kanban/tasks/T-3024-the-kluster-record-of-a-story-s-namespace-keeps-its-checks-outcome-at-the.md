---
id: T-3024
type: task
nature: remediation
title: The kluster record of a story's namespace keeps its checks outcome at the deployed commit, with a gate acceptance can read
status: done
parent: S-0694
owner: alex
created: 2026-11-10T10:32:34Z
updated: 2026-11-11T09:13:35Z
transitions:
  - to: ready
    at: 2026-11-11T08:57:52Z
    by: agent-S-0694
  - to: in-progress
    at: 2026-11-11T08:57:52Z
    by: agent-S-0694
  - to: done
    at: 2026-11-11T09:13:35Z
    by: agent-S-0694
stream: S-0694
tags: [flai, kluster, serve]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go]
usage:
  source: log
  seconds: 943
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 58
      output: 271
      cache_read: 2580346
      cache_write: 124432
      cost: 1.0011
---
# T-3024 The kluster record of a story's namespace keeps its checks outcome at the deployed commit, with a gate acceptance can read

## Work

ADR-0254 §5a refuses an acceptance unless the story's namespace was deployed at the story's head and its checks run against it passed there. Today `kluster.Record` (`flai/internal/kluster/deploy.go`) keeps the deploy's commit and outcome. `ChecksRun` in package serve keeps no commit. The pass-only `klusterTried` is unexported, so `flai accept` cannot read either.

- Give `kluster.Record` the outcome and the commit of the last checks run against the namespace. Have `checks()` in `flai/internal/serve/kluster.go` write both when the run ends, passed or failed.
- Export one gate in package kluster, for example `kluster.AcceptGate(root, story, head, checksOn)`. It answers whether acceptance may go on. When it may not, it names the last job and its outcome.
- The gate passes when the deploy passed at `head` and either `checksOn` is false or the checks passed at `head`. Otherwise it refuses: no record, a deploy at another commit, a failed deploy, or checks failed or not yet run at `head`.
- `head` is the story branch's head, `refs/heads/story/<ID>`, which `KlusterEnv` deploys.

This is the first layer: it waits for no task. The send-back, the blocker, and the queue's reason all read this record.

## Done when

- A checks run against a story's namespace leaves its outcome and commit in the namespace's record.
- Table tests in `deploy_test.go` cover the gate: deploy passed with checks passed, checks off, deploy at another commit, deploy failed, checks failed, checks not run, and no record. Each refusal names the last job and its outcome.
- `flai test` passes on the four files.

## Notes
