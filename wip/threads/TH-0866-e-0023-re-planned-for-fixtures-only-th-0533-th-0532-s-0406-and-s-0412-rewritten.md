---
id: TH-0866
title: "E-0023 re-planned for fixtures only (TH-0533, TH-0532): S-0406 and S-0412 rewritten; proposals to bring S-0407 to S-0411 and the epic's last two criteria in line"
anchor:
  path: wip/kanban/epics/E-0023-an-opencode-adapter-runs-a-story-s-agent-on-opencode-with-the-guard-the-permission-prompt-the-roles-and-usage-measured-through-its-plugin-hooks-and-its-server-s-events-so-a-non-claude-model-works-a-story-end-to-end.md
  item: E-0023
status: resolved
participants: [planner-E-0023, orchestrator, alex]
created: 2026-11-07T09:38:00Z
updated: 2026-11-10T09:31:56Z
---

# TH-0866 E-0023 re-planned for fixtures only (TH-0533, TH-0532): S-0406 and S-0412 rewritten; proposals to bring S-0407 to S-0411 and the epic's last two criteria in line

On wip/kanban/epics/E-0023-an-opencode-adapter-runs-a-story-s-agent-on-opencode-with-the-guard-the-permission-prompt-the-roles-and-usage-measured-through-its-plugin-hooks-and-its-server-s-events-so-a-non-claude-model-works-a-story-end-to-end.md.

## Entries

### 2026-11-07T09:38:00Z planner-E-0023
E-0023 is re-planned on your answers: OpenRouter fixtures only, no live integration (TH-0533), and no OpenCode on the host (TH-0532, D-0506). I rewrote the two drafts, S-0406 and S-0412. I did not change the words of the finalized stories S-0407 to S-0411. The changes they need are proposed below.

## What changed

- **S-0406** (draft). The four open findings are now read from OpenCode's public source at one pinned release, not checked against a binary. The adapter is tested against a stub `opencode` that replays a fixture stream. The binary-gated integration test and `scripts/opencode-smoke.sh` are dropped. Forecast 1h05m (was 1h15m), value 97.01 USD a week.
- **S-0412** (draft). It was an `experiment`: a real non-Claude model through a funded OpenRouter key. It is now a `feature`: an integration test in which a stub `opencode` replays a scripted session in OpenRouter's format and drives a scratch story to review, with the guard and the permission prompt on. Its cost comes from a fake OpenRouter on `httptest`. It spends nothing. Forecast 50m (was 30m), value 74.63 USD a week.
- **`after`**. S-0409 now waits on S-0407 (was S-0408), and S-0410 on S-0406 (was S-0409). Each now names only the story whose output it reads. The old entries only set an order, which ADR-0230 leaves to the pull order.

## Stories, `after`, and task layers

| Story | After: what it reads | Layers |
|-------|----------------------|--------|
| S-0406 | none | T-1682, T-1683, T-1684, T-1685, T-1686 → T-1687 → T-1688 → T-1689 |
| S-0407 | S-0406: the adapter and the sub-agent finding | T-1690 → T-1691, T-1692 → T-1693 → T-1694 |
| S-0408 | S-0407: the plugin it extends | T-1695 → T-1696, T-1697 → T-1698 → T-1699 → T-1700 |
| S-0409 | S-0407: the guard ADR-0131 needs before roles | T-1701 → T-1702, T-1703 → T-1704 → T-1705 |
| S-0410 | S-0406: the adapter and the `step_finish` finding | T-1706, T-1707 → T-1708 → T-1709 → T-1710 → T-1711 |
| S-0411 | S-0410: its OpenCode fixtures | T-1712, T-1713 → T-1714 → T-1715 |
| S-0412 | S-0407 guard, S-0408 ask, S-0410 pricing, S-0411 stream | T-1716 → T-1717 → T-1718 → T-1719 |

## Assumptions

- "Fixtures" means hand-written files that follow OpenCode's source and OpenRouter's documented responses. Nothing is recorded from a paid call, and no test needs a binary, a network, or a key.
- A story's agent may read OpenCode's public GitHub source and docs on the web.
- A live run against a real `opencode` and OpenRouter is left for a later story outside this epic.

## Proposed changes to finalized stories (your word needed)

1. **S-0407, S-0408, S-0409: criterion 4.** Each asks for an integration test "on the known `opencode` version, skipped without it". **Recommended:** replace it with a test on fixtures: the plugin's recorded `tool.execute.before` and `permission.asked` inputs, and a stub `opencode`, with no skip. The same change applies to T-1693, T-1699, and T-1704.
2. **S-0407: criterion 3.** `Guard: true` would rest on a reading of OpenCode's source, not on a run. **Recommended:** allow it on that reading. `agent-adapters.md` and § Harnesses must say the guard was not checked against a live binary. The alternative is to keep `Guard` false until a live run, which also keeps S-0409's roles refused under ADR-0131.
3. **S-0410: criterion 5 and T-1710.** These ask for a live OpenRouter gateway test behind `OPENROUTER_API_KEY`. **Recommended:** replace them with a fake OpenRouter on `httptest` that serves fixtures, as `flai/internal/usage/openrouter_test.go` does, and asserts `priced_by: gateway`. Also, criterion 1 and T-1706 should read "fixtures written from the pinned release" in place of "recorded output from the known version".
4. **E-0023: criteria 8 and 9.** **Recommended:** reword criterion 8 to "a scripted non-Claude session works one story end to end on the opencode adapter against OpenRouter fixtures". Reword criterion 9 to "OpenCode's tests run on fixtures and a stub, with no binary". Those words are yours.
5. **Figures of S-0407 to S-0411.** I left their values as written. With S-0406 and S-0412 re-forecast, `flai cod` moves each by under 4% (for example S-0407 from 69.23 to 67.16 USD a week). That is too small to change the pull order, and rewriting them would mean rewriting finalized bodies. Say if you want them rewritten.

No story to add, split, or drop. Once you answer 1 to 3, I rewrite those criteria and tasks.

### 2026-11-07T09:39:45Z orchestrator (recommendation)
**Recommended:** accept all five of the planner's proposals as written. planner-E-0023 then rewrites S-0407 to S-0411's criteria and tasks on its next run.

Your answer on TH-0533 (D-0507) is "OpenRouter fixtures only, no direct integration in this repository for now". D-0506 adds that no OpenCode goes on the host. Any criterion that needs a live `opencode` or a funded key would be skipped on every host, so it could never be met or ticked.

1. **S-0407, S-0408, S-0409, criterion 4, and T-1693, T-1699, T-1704: yes.** Replace the binary-gated tests with tests on fixtures: recorded `tool.execute.before` and `permission.asked` inputs, and a stub `opencode`, with no skip.
2. **S-0407, criterion 3: yes, allow `Guard: true` on the reading of OpenCode's source at the pinned release.** `agent-adapters.md` and § Harnesses must say the guard was not checked against a live binary. The alternative keeps `Guard` false with no date set, which also keeps S-0409's roles refused under ADR-0131, so S-0409 could deliver nothing. A risk stays: the plugin's sub-agent sign is unproven until a live run. Even if the sign turns out wrong, nothing in this repository runs OpenCode until then.
3. **S-0410, criterion 5 and T-1710: yes.** Use a fake OpenRouter on `httptest` that asserts `priced_by: gateway`, as `flai/internal/usage/openrouter_test.go` does. Criterion 1 and T-1706 should say "fixtures written from the pinned release".
4. **E-0023, criteria 8 and 9: yes, with the planner's words.** The epic's Notes line "Spend: an OpenRouter key with a few dollars of credit" should go too, since this epic now spends nothing.
5. **Figures: leave them.** The orchestration policy is `throughput`, so the ready column is ordered by forecast. A shift of under 4% in value changes no order.

What I did not do: change any criterion. That is scope, so it is yours. S-0406 and S-0412 are already finalized and fit the fixtures-only answer. S-0406 is in progress.

Source: design/digest/D-0507-s-0412-is-a-complete-draft-but-finalizing-it-is-yours-choose-its-non-claude.md

### 2026-11-10T09:31:38Z alex
Confirmed the recommendation of 2026-11-07T09:39:45Z orchestrator.

Source: design/digest/D-0507-s-0412-is-a-complete-draft-but-finalizing-it-is-yours-choose-its-non-claude.md

### 2026-11-10T09:31:56Z alex
Resolved.
