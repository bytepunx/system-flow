---
id: D-0563
title: "E-0023 re-planned for fixtures only (TH-0533, TH-0532): S-0406 and S-0412 rewritten; proposals to bring S-0407 to S-0411 and the epic's last two criteria in line: Recommended: accept all five of the…"
kind: answer
sources: [TH-0866]
topics: [cli, agents, template, flai, go, conventions, code, all]
paths: [flai/internal/harness, flai/internal/guard, flai/internal/usage, flai/internal/mcpserver, flai/internal/serve, flai/internal/manifest, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, template]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0563 E-0023 re-planned for fixtures only (TH-0533, TH-0532): S-0406 and S-0412 rewritten; proposals to bring S-0407 to S-0411 and the epic's last two criteria in line: Recommended: accept all five of the…

The operator, alex, answered TH-0866 at 2026-11-10T09:31:38Z.

What was asked:

```text
E-0023 is re-planned on your answers: OpenRouter fixtures only, no live integration (TH-0533), and no OpenCode on the host (TH-0532, D-0506). I rewrote the two drafts, S-0406 and S-0412. I did not change the words of the finalized stories S-0407 to S-0411. The changes they need are proposed below.
```

The answer:

```text
**Recommended:** accept all five of the planner's proposals as written. planner-E-0023 then rewrites S-0407 to S-0411's criteria and tasks on its next run.

Your answer on TH-0533 (D-0507) is "OpenRouter fixtures only, no direct integration in this repository for now". D-0506 adds that no OpenCode goes on the host. Any criterion that needs a live `opencode` or a funded key would be skipped on every host, so it could never be met or ticked.

1. **S-0407, S-0408, S-0409, criterion 4, and T-1693, T-1699, T-1704: yes.** Replace the binary-gated tests with tests on fixtures: recorded `tool.execute.before` and `permission.asked` inputs, and a stub `opencode`, with no skip.
2. **S-0407, criterion 3: yes, allow `Guard: true` on the reading of OpenCode's source at the pinned release.** `agent-adapters.md` and § Harnesses must say the guard was not checked against a live binary. The alternative keeps `Guard` false with no date set, which also keeps S-0409's roles refused under ADR-0131, so S-0409 could deliver nothing. A risk stays: the plugin's sub-agent sign is unproven until a live run. Even if the sign turns out wrong, nothing in this repository runs OpenCode until then.
3. **S-0410, criterion 5 and T-1710: yes.** Use a fake OpenRouter on `httptest` that asserts `priced_by: gateway`, as `flai/internal/usage/openrouter_test.go` does. Criterion 1 and T-1706 should say "fixtures written from the pinned release".
4. **E-0023, criteria 8 and 9: yes, with the planner's words.** The epic's Notes line "Spend: an OpenRouter key with a few dollars of credit" should go too, since this epic now spends nothing.
5. **Figures: leave them.** The orchestration policy is `throughput`, so the ready column is ordered by forecast. A shift of under 4% in value changes no order.

What I did not do: change any criterion. That is scope, so it is yours. S-0406 and S-0412 are already finalized and fit the fixtures-only answer. S-0406 is in progress.

Source: design/digest/D-0507-s-0412-is-a-complete-draft-but-finalizing-it-is-yours-choose-its-non-claude.md
```

Source: TH-0866.
