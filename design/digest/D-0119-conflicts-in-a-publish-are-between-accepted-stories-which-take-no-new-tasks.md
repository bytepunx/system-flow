---
id: D-0119
title: "Conflicts in a publish are between accepted stories, which take no new tasks: where do the conflict tasks go: I would like a different convention that 1 - I don't want new stories created for a story…"
kind: answer
sources: [TH-0073]
topics: [flai, cli, go, template, conventions, code, all]
paths: [template/root/design/conventions/work-management.md, design/conventions/work-management.md, flai/cmd/release.go, docs/users/flai-reference.md, flai/cmd/release_moved_test.go, flai/internal/release/remote.go, flai/internal/release/remote_test.go, design/system/pushing-from-the-board.md, docs/users/flai.md, design/system/flai-cli.md]
date: 2026-10-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0119 Conflicts in a publish are between accepted stories, which take no new tasks: where do the conflict tasks go: I would like a different convention that 1 - I don't want new stories created for a story…

The operator, alex, answered TH-0073 at 2026-10-03T01:50:36Z.

What was asked:

```text
TH-0070 says merge conflicts in a publish are worked "through tasks on the stories in question, in threads with the operator". But the stories whose changes conflict when a clone takes its remote are accepted: done and archived, on both sides. `flai task new` refuses a task under a done or archived story (workitem.Create: "S-nnnn is done; reopen or pick another parent"), and nothing moves a story back out of done.
```

The answer:

```text
I would like a different convention that 1 - I don't want new stories created for a story stuck in review since the overlaps would likely cause unwanted blocking behavior. Instead I prefer option 3 - use threads to resolve conflicts.
```

Source: TH-0073.
