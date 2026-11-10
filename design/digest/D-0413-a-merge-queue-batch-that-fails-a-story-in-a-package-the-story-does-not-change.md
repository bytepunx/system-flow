---
id: D-0413
title: A merge queue batch that fails a story in a package the story does not change is fixed inside the story, not re-enqueued unchanged, since the next loaded batch fails the same way
kind: decision
sources: [S-0507, TH-0723, TH-0724, I-0225, I-0226, S-0538, TH-0732, I-0203, S-0529, TH-0767, I-0221, S-0548, S-0652, I-0331]
topics: [merge-queue, testing, go, code]
paths: [scripts/integration.sh, flai/internal/mcpserver, flai/internal/serve/agents_test.go, flai/internal/serve/forgemerge_test.go, flai/internal/guard/roles_test.go, flai/internal/mcpserver/verify_test.go, flai/internal/serve/components_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0413 A merge queue batch that fails a story in a package the story does not change is fixed inside the story, not re-enqueued unchanged, since the next loaded batch fails the same way

The merge queue sent S-0507 back twice. Both times the failure was in a package S-0507 does not change. Each time its agent fixed the cause on the story's branch instead of enqueueing the same commit again. Its reason: "the next batch would likely fail the same way."

- TH-0723: batch 20261028-231632 failed the integration tier because `flai/internal/mcpserver` hit Go's 10-minute default test timeout (600.17 s). That day the package took 470 to 591 s under `-race` in every batch. T-2316 gave `scripts/integration.sh` `-timeout 20m`, and the script's comment says why. S-0548 closed I-0225, the suite's slowness, on 2026-10-31: the tests run in parallel, 9m14s to 136 to 162 s, and the go-test tier passes `-timeout 20m` too (D-0416).
- TH-0724: the next batch passed mcpserver at 679 s, then failed two `flai/internal/serve` tests that pass when run alone. Each checked for something before it had happened. T-2317 made one test check for one start and the other wait for the log line. I-0226 records it. S-0548's T-2362 closed I-0226. The second instance "was a race in the test, not a short deadline (2 failures in 100 runs reproduced, 0 in 200 after)": the stub could end before the same look restarted what ended, so the test now holds the stub. `waitFor` and the merge test's waits allow 30 s.
- TH-0732: batch 20261029-032429 sent S-0538 back at integration, because main itself failed `TestADefinedRoleIsHeldToItsFile` in `flai/internal/guard/roles_test.go`. S-0416 had been accepted by hand, outside the queue, after S-0535 reworded the to-backlog condition, and its test still expected the old reason (I-0203, bumped). S-0538's agent added T-2324 to update the expectation, byte-identical to S-0519's T-2323, so the two stories rebase cleanly over each other. It settled that with both in their conversations.
- TH-0767: batch 20261030-182125 sent S-0529 back at integration. S-0529 had changed `flai verify`'s reuse rule on purpose. `TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch`, in `flai/internal/mcpserver/verify_test.go`, a package S-0529 does not change, still expected the old rule. Its close-out had passed because the `go-test` tier runs only the changed packages (I-0221, bumped). T-2345 updated the test, and "a search found no other test outside the changed packages that pins the old rule."

A later story sent back for a failure outside its own change should do the same:

1. Find the cause first, in the batch log and in the package's times across batches.
2. If the cause is timing or load, fix the timeout or the test in the story. A test that fails only under load may race its own stub rather than wait too briefly, as I-0226's did. Reproduce it with `GOMAXPROCS=2 go test -race -count=100` before you lengthen a wait.
3. If main itself fails, fix main's expectation in the story with the same edit any other story making it uses.
4. If the story's own deliberate change broke a test in a package it does not change, update that test in the story. Search the other packages for tests that pin the old behaviour.
5. Record the underlying cost as an issue for a story of its own.
6. Enqueue again only after that.

A close-out stopped by a test outside the story's change is handled the same way. S-0652's first close-out stopped at integration on `TestADevServerThatExitsOnItsOwnIsForgottenWithItsError`, outside its issue but of the same cause. From its Decisions (`wip/archive/agents/S-0652.md`): "I recorded it as I-0331 and fixed it here as T-2859, a third layer after the first two, rather than running the close-out again until it passed (code-quality: no retries)." (D-0551)

Source: S-0507 (wip/archive/agents/S-0507.md, Decisions), TH-0723, TH-0724; S-0538 (wip/archive/agents/S-0538.md, Decisions), TH-0732; S-0529 (wip/archive/agents/S-0529.md, Decisions), TH-0767; S-0548 (wip/archive/agents/S-0548.md, Decisions), I-0225, I-0226; S-0652 (wip/archive/agents/S-0652.md, Decisions), I-0331.
