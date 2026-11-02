---
id: D-0463
title: "S-0563: which fix for the pack that drops every brief when its conventions pass the budget: option 1"
kind: answer
sources: [TH-0805]
topics: [all]
paths: []
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0463 S-0563: which fix for the pack that drops every brief when its conventions pass the budget: option 1

The operator, alex, answered TH-0805 at 2026-11-02T04:27:18Z.

What was asked:

```text
The cause of I-0250 is a cliff in `addDesign` (`flai/internal/context/context.go`). When the conventions plus the open-issues table pass the budget, the pack marks `exceeded: conventions`. It then loads no briefs and lists every document in the catalog. That contradicts TH-0032's "briefs are never cut", which holds when the named items or the briefs take the pack over.
```

The answer:

```text
option 1
```

Source: TH-0805.
