---
id: TH-0823
title: "S-0510 cannot be published: under judgement the publish formula stops at step 2, and release_publish refuses (TH-0811)"
anchor:
  path: wip/archive/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md
  item: S-0510
status: resolved
participants: [orchestrator]
created: 2026-11-03T02:41:14Z
updated: 2026-11-03T17:57:40Z
---

# TH-0823 S-0510 cannot be published: under judgement the publish formula stops at step 2, and release_publish refuses (TH-0811)

On wip/archive/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md.

## Entries

### 2026-11-03T02:41:14Z orchestrator
**Recommended:** answer TH-0811. Either set `policy: threshold` with `count: 1` under `orchestration.release` in `system-flow.yaml`, or publish S-0510 by hand with `flai release --pending`.

**Where it stands:**

- I judge S-0510 a coherent, complete release: one accepted feature, all five tasks done, 27.78 USD/week.
- `release_evaluate` gives policy `judgement`, `met: false`.
- `release_publish` refused:

```text
publishing in this project is the publish formula, design/workflows/formulas/publish.md, so release_publish publishes nothing: start it with flai formula run publish --ready --autocommit, as your role file says under publish, when release_evaluate is met (ADR-0240)
```

- I did not run the publish formula. Its step 2, `flai release --evaluate --require-met`, always exits 2 under judgement, so the instance would release nothing. flai 1.82.0 adds no flag that changes this.

I will try again at the next acceptance, or once you answer here or on TH-0811.

### 2026-11-03T17:57:40Z orchestrator
Resolved: S-0510 is released: the operator published flai 1.83.0, flaiover 0.56.2, template 1.1.22 in 3b80afc2e, and release_evaluate now lists nothing pending. The judgement-policy question itself stays open on TH-0811.
