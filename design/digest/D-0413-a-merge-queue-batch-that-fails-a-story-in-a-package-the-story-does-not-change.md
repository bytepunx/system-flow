---
id: D-0413
title: A merge queue batch that fails a story in a package the story does not change is fixed inside the story, not re-enqueued unchanged, since the next loaded batch fails the same way
kind: decision
sources: [S-0507, TH-0723, TH-0724, I-0225, I-0226, S-0538, TH-0732, I-0203, S-0529, TH-0767, I-0221, S-0548, S-0652, I-0331]
topics: [merge-queue, testing, go, code]
paths: [scripts/integration.sh, flai/internal/mcpserver, flai/internal/serve/agents_test.go, flai/internal/serve/forgemerge_test.go, flai/internal/guard/roles_test.go, flai/internal/mcpserver/verify_test.go, flai/internal/serve/components_test.go]
date: 2026-10-29
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0413 A merge queue batch that fails a story in a package the story does not change is fixed inside the story, not re-enqueued unchanged, since the next loaded batch fails the same way
