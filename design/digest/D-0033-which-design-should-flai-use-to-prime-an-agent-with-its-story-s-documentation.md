---
id: D-0033
title: "Which design should flai use to prime an agent with its story's documentation: right now, stories already use tags to indicate how version tags will be calculated for the components. let's…"
kind: answer
sources: [TH-0020]
topics: [flai, cli, go, code, all]
paths: [flai/cmd]
date: 2026-09-26
status: absorbed
absorbed_by: "design/adrs/0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md#topics-on-documents"
asked_again: []
cited: []
---

# D-0033 Which design should flai use to prime an agent with its story's documentation: right now, stories already use tags to indicate how version tags will be calculated for the components. let's…

The operator, alex, answered TH-0020 at 2026-09-26T17:37:57Z.

What was asked:

```text
The findings are in `design/system/agent-context.md` on `story/S-0125`. Today every agent reads all 12 conventions (47 KB) and the whole issues table, then has to find the relevant ADRs and design itself.
```

The answer:

```text
right now, stories already use tags to indicate how version tags will be calculated for the components. let's disambiguate this by calling them "topics" instead of "tags" with the same intended functionality, just without the overloading of terminology
```

Source: TH-0020.
