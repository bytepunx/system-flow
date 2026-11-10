---
id: I-0268
title: "ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built"
class: impression
status: open
count: 2
cost: 10m
first_reported: 2026-11-07T03:22:56Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-10T08:54:25Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
---

# I-0268 ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built

## Description
ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built

## Instances

### 2026-11-07T03:22:56Z
Story: S-0579.
S-0579 built flai serve's kluster pass to its criteria: a failed deploy, or a failed checks run against the story's namespace, opens one thread on the story and moves nothing. ADR-0254 §5a also says flai serve sends such a story back to in-progress and starts its agent. It also says every acceptance is refused under the `kluster` blocker unless the namespace passed at the story's head, with `flai accept --without-kluster` as the operator's override. No criterion of S-0579 named either, and no other story under E-0032 builds them. So a story whose namespace failed can be accepted, and design/system/kluster.md § A failed deploy holds the story says so. A story would build the send-back in flai/internal/serve/kluster.go, and the blocker and the flag in flai accept, the orchestrator's acceptance, and the merge queue.

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: the living design is honest about the gap, the accepted ADR reads as the behaviour.

## Impact
ADR-0254 § 5a and its Consequences (S-0579 adds the kluster blocker) are accepted; flai/internal/serve/kluster.go:43 and :626 move nothing on a failure; no --without-kluster under flai/cmd or flai/internal; design/system/kluster.md line 119 records the gap

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0694 remediates this issue, created from it at 2026-11-10T08:54:25Z.
