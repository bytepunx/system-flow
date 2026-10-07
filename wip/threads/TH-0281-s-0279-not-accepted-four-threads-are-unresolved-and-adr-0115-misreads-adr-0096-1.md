---
id: TH-0281
title: "S-0279 not accepted: four threads are unresolved, and ADR-0115 misreads ADR-0096 §1"
anchor:
  path: wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md
  item: S-0279
status: resolved
participants: [orchestrator, alex]
created: 2026-10-07T09:52:57Z
updated: 2026-10-07T14:35:49Z
---

# TH-0281 S-0279 not accepted: four threads are unresolved, and ADR-0115 misreads ADR-0096 §1

On wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T09:52:57Z orchestrator
Recommendation: on TH-0269, accept ADR-0115 once it says it narrows ADR-0096 §1's `wip.overlap` clause, instead of citing §1 for the exclusion. Then resolve TH-0269, TH-0273, TH-0274, and TH-0275, and I accept S-0279 at the next event.

I left S-0279 in review. `flai accept S-0279 --by orchestrator --verified 955a1288 --dry-run` has four blockers:

```text
blocked: thread TH-0269 on S-0279 is open, not resolved: Accept ADR-0115, the remedy for I-0076?
blocked: thread TH-0273 on S-0279 is answered, not resolved: S-0213 and S-0279 conflict when merged
blocked: thread TH-0274 on S-0279 is answered, not resolved: S-0214 and S-0279 conflict when merged
blocked: thread TH-0275 on S-0279 is answered, not resolved: S-0215 and S-0279 conflict when merged
```

The ADR point, confirmed by the verifier:

- ADR-0115 §1 says "A story in review … is not compared, as it holds nothing (ADR-0096 §1)".
- ADR-0096 §1 says the opposite for this rule: a story in review keeps its claim for `flai check`'s `wip.overlap`.
- The code already compared only stories in progress before this diff, so behaviour is unchanged. Only the ADR's wording and citation need fixing.

Everything else is clear at 955a1288, the branch head:

- `flai verify` passed every tier there.
- The verifier matched each criterion to changed files within the touches, with the convention landing in both copies and `template/CHANGELOG.md` 1.0.70:
  - 1: `flai/internal/check/check.go`, `flai/internal/check/scope.go`, `flai/cmd/check.go`, and their tests
  - 2: `design/issues/I-0076-…md`, `design/issues/summary.md`

### 2026-10-07T14:35:49Z alex
Resolved: S-0279 was accepted
