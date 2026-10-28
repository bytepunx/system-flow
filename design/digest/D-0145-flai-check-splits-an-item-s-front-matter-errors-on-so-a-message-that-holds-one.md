---
id: D-0145
title: "flai check splits an item's front-matter errors on \"; \", so a message that holds one becomes two findings, the second on line 1"
kind: cause
sources: [I-0061]
topics: [flai, cli, go, code, all]
paths: [design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/I-0061-flai-check-splits-an-item-s-front-matter-errors-on-so-a-message-that-holds-one-becomes-two-findings-the-second-on-line-1.md, design/issues/summary.md, design/system/flai-cli.md, flai/internal/check/check.go, flai/internal/check/planning_test.go, flai/internal/manifest/agent.go, flai/internal/workitem/item.go]
date: 2026-10-04
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0145 flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1

I-0061 was closed on 2026-10-04.

The cause:

```text
flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1
```

The fix:

```text
S-0248: Item.Validate returns workitem.Problems, one phrase per problem, and flai check ranges over it instead of splitting the joined message on "; "; TestAFrontMatterProblemHoldingASemicolonIsOneFinding reproduces the split
```

Source: I-0061.
