---
id: D-0040
title: "A new publish host action, off by default, decides whether pushing also releases. Agree: A: Call it `auto-publish`"
kind: answer
sources: [TH-0031]
topics: [flai, cli, go, code, all]
paths: [flai/cmd/push.go, flai/cmd/push_test.go, flai/cmd/release.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/pending, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/conventions/git.md, design/adrs, design/issues, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
date: 2026-09-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0040 A new publish host action, off by default, decides whether pushing also releases. Agree: A: Call it `auto-publish`

The operator, alex, answered TH-0031 at 2026-09-29T02:01:31Z.

What was asked:

```text
**Cause.** Since S-0094, `flai push --pending` tags a release of everything accumulated before it pushes. Agents are told to run it whenever `inbox` reports an unpushed acceptance, and the board's Push now button runs it too. So every acceptance ends up published, one release per story.
```

The answer:

```text
A: Call it `auto-publish`. Will the tagging be affected such that previously pushed commits end up missing version changes as a result of this change?
```

Source: TH-0031.
