---
id: D-0381
title: flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s
kind: cause
sources: [I-0166]
topics: [flai, cli, go, code, all]
paths: [design/system/component-library.md, docs/operators/runbooks/backup.md, flai/internal/components/props.go, flai/internal/components/props.mjs, flai/internal/components/props_test.go, design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md, design/issues/summary.md]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0381 flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s

I-0166 was closed on 2026-10-23.

The cause:

```text
flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s
```

The fix:

```text
Fixed by S-0478 (T-2068): props.mjs keeps the skeletons it resolves in a cache under .flai-cache/component-props/<sub-project path>/, keyed by the script, the sub-project, its typescript, and the components' typings, and stamped by every path typescript looked at. A run with nothing changed takes them from the cache without loading typescript: flaiover's 58 components take 0.43s, against 2.7s cold. TestListTakesSkeletonsFromTheCacheUntilAFileTypescriptReadChanges reproduces it.
```

Source: I-0166.
