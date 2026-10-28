---
id: D-0184
title: flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out
kind: cause
sources: [I-0081]
topics: [cli, flai, go, code, all]
paths: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, design/adrs, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0081-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md, design/issues/summary.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-06
status: absorbed
absorbed_by: "design/adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md#decision"
asked_again: []
cited: []
---

# D-0184 flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out

I-0081 was closed on 2026-10-06.

The cause:

```text
flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out
```

The fix:

```text
S-0284: permission_prompt takes an answer from the story's owner or the project's owner, the manifest's owner, and never the asking agent's own entry (ADR-0097, refining ADR-0086); TestPermissionPromptTakesTheProjectOwnersAnswer reproduces S-0218's case, story owner arobson and the operator answering allow as alex
```

Source: I-0081.
