---
id: TH-0878
title: "S-0410 cannot be accepted: criteria 1 and 5 are unticked, and both need a live opencode or a funded key"
anchor:
  path: wip/kanban/stories/S-0410-a-story-on-opencode-has-its-usage-read-from-opencode-s-json-events-priced-by-the-harness-or-its-gateway-s-spend-log-and-reaches-its-provider-through-an-entry-derived-from-the-host-s.md
  item: S-0410
status: open
participants: [orchestrator]
created: 2026-11-07T21:45:47Z
updated: 2026-11-07T21:45:47Z
---

# TH-0878 S-0410 cannot be accepted: criteria 1 and 5 are unticked, and both need a live opencode or a funded key

On wip/kanban/stories/S-0410-a-story-on-opencode-has-its-usage-read-from-opencode-s-json-events-priced-by-the-harness-or-its-gateway-s-spend-log-and-reaches-its-provider-through-an-entry-derived-from-the-host-s.md.

## Entries

### 2026-11-07T21:45:47Z orchestrator
**Recommended:** answer TH-0866 point 3. Then reword S-0410's criteria 1 and 5 to the fixtures-only form and tick them. I accept it on my next run.

- Criterion 1: "tested on fixtures written from the pinned release".
- Criterion 5: "a gateway test on a fake spend log asserts `priced_by: gateway` and a cost above zero".

**Blockers:** `flai accept S-0410 --by orchestrator --verified 6940235ff6f --dry-run` lists only these:

```text
blocked: acceptance criterion 1 of S-0410 is unticked
blocked: acceptance criterion 5 of S-0410 is unticked
```

**Verified:** `flai verify S-0410 --last` passed every step at 6940235ff6f, the head of `story/S-0410`. That includes integration and smoke.

**Why the criteria are unticked:** the story's Notes say why.

- Criterion 1: the fixtures under `flai/internal/usage/testdata/opencode/` follow opencode 1.19.5's source. They were not recorded from a live run. This host has no `opencode`.
- Criterion 5: `TestGatewayOpenCode` exists and its skip is tested. Its paid branch never ran, because the host has neither `OPENROUTER_API_KEY` nor `opencode`.

Your answers D-0506 and D-0507 say fixtures only, with no OpenCode on the host. Under them, both criteria as written can never be ticked.

**What the branch already has for the reworded criteria:**

- `TestAGatewayPricesAnOpenCodeRunThroughAProvider` in `flai/internal/serve/usage_test.go`. It prices two OpenCode runs through a fake spend log. It asserts `priced_by: gateway`, a cost of 1.0, and OpenCode's tokens kept.
- The reader tests in `flai/internal/usage/opencode_test.go`, on the two fixtures.

**Alternative:** accept it yourself as it stands:

```bash
flai accept S-0410 --by alex
```

Rewording criteria is a change of scope, so it is your call. I have changed nothing on S-0410. S-0411 is held in ready until S-0410 is done.
