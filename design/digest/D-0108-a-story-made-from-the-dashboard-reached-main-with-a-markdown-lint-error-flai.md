---
id: D-0108
title: A story made from the dashboard reached main with a markdown lint error flai task new would have refused
kind: cause
sources: [I-0055]
topics: [dashboard, markdown, flai, cli, go, flaiover, sveltekit, code, all]
paths: [flai/internal/mdlint, flai/internal/itemnew, flai/internal/hostapi, flaiover/src/lib/server, flaiover/src/lib/components, design/issues, design/system/flai-cli.md, docs/users/flaiover.md]
date: 2026-10-02
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0108 A story made from the dashboard reached main with a markdown lint error flai task new would have refused

I-0055 was closed on 2026-10-02.

The cause:

```text
A story made from the dashboard reached main with a markdown lint error flai task new would have refused
```

The fix:

```text
S-0240: the dashboard's item.new and doc.save already ran flai's wip lint; the lint had no MD007, so flai story new took S-0231's body too. mdlint now implements MD007 and refuses it with the rule and the line.
```

Source: I-0055.
