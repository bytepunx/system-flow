---
id: TH-0318
title: "Plan: eight draft stories give agents messages of their own, route conflicts to them, and let a holding agent share paths"
anchor:
  path: wip/kanban/epics/E-0018-improve-inter-agent-coordination-via-inboxes.md
  item: E-0018
status: open
participants: [planner-E-0018]
created: 2026-10-07T20:22:53Z
updated: 2026-10-07T20:22:53Z
---

# TH-0318 Plan: eight draft stories give agents messages of their own, route conflicts to them, and let a holding agent share paths

On wip/kanban/epics/E-0018-improve-inter-agent-coordination-via-inboxes.md.

## Entries

### 2026-10-07T20:22:53Z planner-E-0018
E-0018 had no stories. I drafted eight, S-0330 to S-0337, with 36 tasks. All are drafts in the backlog, with touches, forecasts, and cost of delay values. `flai check --strict` passes.

## Stories and their order

| Story | What it delivers | after | Forecast | CoD / week |
|-------|------------------|-------|----------|------------|
| S-0330 | `flai message`: a store of conversations between two open stories, apart from threads | none | 48m | 124.03 |
| S-0331 | MCP `message_*` tools, `inbox.messages`, `wait_for_events` wakes on one, guard, convention | S-0330 | 44m | 113.70 |
| S-0332 | Trial-merge conflicts and grown overlaps become a conversation; `message escalate` asks the operator | S-0331 | 55m | 142.12 |
| S-0333 | `flai task done` messages each open story whose claim covers what the task changed | S-0331 | 44m | 113.70 |
| S-0334 | flai serve asks a holding agent about a held story; `message share` lifts the hold | S-0332 | 58m | 149.87 |
| S-0335 | An agent waiting on another agent's reply ends, and flai serve restarts it on the reply | S-0331 | 36m | 93.02 |
| S-0336 | Dashboard: Messages view, conversations on a story, held card says it was asked | S-0334 | 48m | 124.03 |
| S-0337 | `flai stats` coordination values: conflicts found early against late, shares, hold time saved | S-0332, S-0333, S-0334 | 54m | 139.53 |

The total is 6h27m. Each value is the story's share of the epic's 1000 USD a week by forecast duration (`flai cod`). Many stories touch the same design documents and the messages package, so the overlap hold will mostly run them one at a time even where `after` allows more.

## Tasks and layers

Layers are separated by `→`. Tasks inside brackets run together.

- S-0330: T-1185 ADR → T-1186 package → [T-1187 command, T-1188 close at accept or archive, and check] → T-1189 docs
- S-0331: [T-1190 tools and instructions, T-1191 guard] → T-1192 inbox and wait → T-1193 convention and docs
- S-0332: T-1194 ADR → [T-1195 sync, T-1196 grown claim, T-1197 escalate] → T-1198 docs
- S-0333: T-1199 move the coverage test into itemedit → T-1200 notice at task close → T-1201 docs
- S-0334: T-1202 ADR → T-1203 share and ask → [T-1204 hold honours shares, T-1205 serve asks and passes the split on] → T-1206 docs
- S-0335: [T-1207 end on a reply, T-1208 serve restarts] → T-1209 dashboard → T-1210 docs
- S-0336: T-1211 host reads → T-1212 repo and route → T-1213 view → T-1214 story page and card → T-1215 docs
- S-0337: T-1216 ADR and metrics.md → T-1217 conflict log → T-1218 stats → T-1219 chart → T-1220 docs

## Assumptions

- Messages are addressed story to story, not agent to agent. Whichever agent works a story reads its messages, so restarts lose nothing (I-0037).
- Each conversation is one markdown file under `wip/messages/`, like threads. T-1185's ADR decides this and may choose otherwise. Messages never count as awaiting the operator.
- flai writes the conversations it starts (conflicts, task notices, hold requests) under its own author name, as the conflict threads are today.
- Each conversation is between two open stories. The one exception is S-0334's hold request, where one side is the held story, which is still in ready.
- No new manifest setting and no new host action.

## Decision for you when you finalize S-0334

S-0334 changes the hold policy you chose on TH-0019 (ADR-0046). A holding agent's share would lift a hold without you. I recommend it: the share is recorded and shown on the card (S-0336), and the trial merge still catches a split that did not hold (S-0332). If you would rather confirm each share yourself, say so and I will change S-0334 so that a share waits for your confirmation.

## Proposals

- **Split S-0336:** the reads, the view, and the story page need only S-0330. The held-card part needs S-0334. Split it, and the operator sees messages three stories sooner. I did not split it, because one dashboard story touches `flaiover` once.
- **Add later, not drafted:** a report of the paths most often held on or in conflict, proposing seams that would separate concerns, as the epic's notes suggest. It needs S-0337's data first. The analyzer could do it, or a research story once S-0337 has run a few weeks.
- **Outside E-0018:** S-0306 (ADR README conflicts, and a thread for each sync that finds one) changes the code S-0332 reroutes. Plan S-0306 after S-0332, or fold its thread part into S-0332.

Nothing else to merge or drop.
