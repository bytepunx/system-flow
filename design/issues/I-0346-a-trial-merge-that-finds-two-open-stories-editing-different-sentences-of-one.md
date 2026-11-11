---
id: I-0346
title: A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance
class: efficiency
status: closed
count: 1
cost: 15m
first_reported: 2026-11-11T01:43:39Z
last_reported: 2026-11-11T01:43:39Z
updated: 2026-11-11T14:53:12Z
---

# I-0346 A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance

## Description
A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance

## Instances

### 2026-11-11T01:43:39Z
Story: S-0688.
S-0688 and S-0666 each changed a different sentence of one paragraph in design/system/agent-adapters.md § A container harness of flai's own and in docs/operators/index.md § Strategic agents in a container. flai stream sync's trial merge reported the conflict on MS-0474, and the two agents agreed that the second to land keeps both sides. But nothing could act on that until acceptance: the section merge (ADR-0175, ADR-0276) does not settle a paragraph both sides changed, so the operator's acceptance of S-0688 after S-0666 aborted, sent S-0688 back to in-progress on TH-1061, and S-0688 needed a hand rebase, a second close-out, and a second acceptance. The resolution was mechanical: main's paragraph with S-0688's sentence applied. A merge by sentence within a paragraph, or the second story rebasing as soon as the first is accepted (the overlapped event) rather than at its own acceptance, would avoid the round trip.

## Remediation

Story S-0719 remediates this issue, created from it at 2026-11-11T02:29:04Z.
Closed 2026-11-11T14:53:12Z: S-0719: the section merge now merges by sentences the prose lines both sides change (ADR-0287, mergeLines in flai/internal/storygit/mdmerge.go), so a paragraph two stories change in different sentences no longer stops flai stream sync, task done's sync, flai accept, or the merge queue's batch, and the second story is not sent back; both changing one sentence still stops. testSyncMergesBySentencesAParagraphMainChangedAnotherSentenceOf and testAcceptLandsAStoryThatChangedAnotherSentenceOfAParagraphMainChanged reproduce I-0346 in real git and fail without the fix. The trial merge still reports the pair, as ADR-0276 leaves it.
