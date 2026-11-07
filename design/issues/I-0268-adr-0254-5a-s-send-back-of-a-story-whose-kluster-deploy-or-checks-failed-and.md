---
id: I-0268
title: "ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built"
class: impression
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T03:22:56Z
last_reported: 2026-11-07T03:22:56Z
updated: 2026-11-07T03:22:56Z
---

# I-0268 ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built

## Description
ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built

## Instances

### 2026-11-07T03:22:56Z
Story: S-0579.
S-0579 built flai serve's kluster pass to its criteria: a failed deploy, or a failed checks run against the story's namespace, opens one thread on the story and moves nothing. ADR-0254 §5a also says flai serve sends such a story back to in-progress and starts its agent. It also says every acceptance is refused under the `kluster` blocker unless the namespace passed at the story's head, with `flai accept --without-kluster` as the operator's override. No criterion of S-0579 named either, and no other story under E-0032 builds them. So a story whose namespace failed can be accepted, and design/system/kluster.md § A failed deploy holds the story says so. A story would build the send-back in flai/internal/serve/kluster.go, and the blocker and the flag in flai accept, the orchestrator's acceptance, and the merge queue.

## Remediation
