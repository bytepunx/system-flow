---
id: D-0155
title: item_edit refuses a touch that starts with a dot, though flai touches accepts it
kind: cause
sources: [I-0071]
topics: [flai, cli, go, code, all]
paths: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/itemedit/itemedit.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, design/system/work-hierarchy.md, design/issues, docs/users/flai.md]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0155 item_edit refuses a touch that starts with a dot, though flai touches accepts it

I-0071 was closed on 2026-10-05.

The cause:

```text
item_edit refuses a touch that starts with a dot, though flai touches accepts it
```

The fix:

```text
S-0260: every writer of touches (flai touches, flai edit, story and task new, MCP item_new and item_edit, the dashboard) applies one rule, workitem.CleanTouches, which takes a path starting with a dot and refuses a comma, a leading dash, an absolute path, a .. segment, or a control character
```

Source: I-0071.
