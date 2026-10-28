---
id: D-0121
title: "S-0242 and S-0198 each wait for the other on wip.overlap: may S-0242 go to review with the overlap warnings: take recommendation"
kind: answer
sources: [TH-0077]
topics: [flai, cli, go, template, conventions, code, all]
paths: [template/root/design/conventions/work-management.md, design/conventions/work-management.md, flai/cmd/release.go, docs/users/flai-reference.md, flai/cmd/release_moved_test.go, flai/internal/release/remote.go, flai/internal/release/remote_test.go, design/system/pushing-from-the-board.md, docs/users/flai.md, design/system/flai-cli.md]
date: 2026-10-03
status: active
absorbed_by: ""
asked_again: [TH-0095]
cited: []
---

# D-0121 S-0242 and S-0198 each wait for the other on wip.overlap: may S-0242 go to review with the overlap warnings: take recommendation

The operator, alex, answered TH-0077 at 2026-10-03T02:29:40Z.

What was asked:

```text
S-0242 and S-0198 are both complete. Each one's close-out passes everything except `flai check --strict`. That fails only on `wip.overlap` warnings for docs/users/flai.md, docs/users/flai-reference.md, and design/system/flai-cli.md, which both stories really change. `flai stream sync` says the two merge cleanly. On TH-0073 I said I would wait for S-0198's acceptance, and on TH-0076 S-0198's agent says it is waiting for S-0242's acceptance. So each is waiting for the other.
```

The recommendation the operator took: "let S-0242 go first. I'd run a fresh close-out verifier, which fails only on these three warnings, and move S-0242 to review with them standing, recorded in its narrative. Once you accept it, S-0198 syncs onto it and closes out clean."

The answer:

```text
take recommendation
```

So when two stories in progress each wait on the other's `wip.overlap`, one goes to review first with the warnings recorded in its narrative; the rule counts only stories in progress, so that frees the other. The same day, TH-0095 (S-0206 and S-0204) asked the same question and named this answer: "the case you settled on TH-0077 by letting one story go first". ADR-0115 later had a close-out record no `wip.overlap` at all.

Source: TH-0077.
