---
id: D-0551
title: "TestComponentsServeEndsWhenViteExits failed under load because its crash-mode fake vite exited on a 300 ms timer that raced flai's 150 ms poll; the fake now exits only after it has answered once (I-0325)"
kind: cause
sources: [I-0325, S-0650, T-2836, I-0331, I-0329, S-0652, T-2841, T-2859]
topics: [testing, flai, go, code]
paths: [flai/cmd/components_test.go, design/issues/I-0325-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md, flai/internal/serve/components_test.go, flai/internal/components/devserver_test.go, design/issues/I-0331-testadevserverthatexitsonitsownisforgottenwithitserror-fails-the-integration.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0551 TestComponentsServeEndsWhenViteExits failed under load because its crash-mode fake vite exited on a 300 ms timer that raced flai's 150 ms poll; the fake now exits only after it has answered once (I-0325)

I-0325 was closed on 2026-11-10 by S-0650 (T-2836). `TestComponentsServeEndsWhenViteExits` in `flai/cmd/components_test.go` failed under load with "exited before it answered".

The cause: the crash-mode fake vite exited 2 on a 300 ms timer. Under load that timer fired before `StartDevServer`'s 150 ms poll had seen the fake answer, so flai reported an early exit rather than a vite that answered and then died.

The fix, from the issue's close: "the crash-mode fake vite in flai/cmd/components_test.go exits 2 only after it has answered the harness's path once, instead of on a 300 ms timer that raced StartDevServer's 150 ms poll." S-0650's Decisions (`wip/archive/agents/S-0650.md`) add that it waits 100 ms more "for the answer to reach flai", and that "the product code is unchanged; it was right on both paths."

The reproduction without load: `FLAI_FAKE_COMPONENTS_VITE_DELAY`, slept before the fake listens. `TestComponentsServeEndsWhenASlowViteExits` sets it to 600 ms; with the old timer in a scratch copy it failed with I-0325's own error, and with the fix it passes.

So a fake process in a flai test exits on the event it stands for, not on a timer, as D-0491 says of tests that raced their own setup. An earlier fault in the same dev-server tests, a free port closed before the fake bound it, was I-0162 (D-0319).

S-0652 met the same race twice more. I-0331: `TestADevServerThatExitsOnItsOwnIsForgottenWithItsError` failed its close-out because "the crash-later fake vite crashes before flai's probe sees it answer". T-2859 fixed it: "the crash-later fake vites in flai/internal/serve/components_test.go and flai/internal/components/devserver_test.go now crash 300ms after they first answer, not 300ms after they start". A 500 ms listen delay in the serve test failed it on every run with the old fake. I-0329's first instance, `TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered`, gave its fake 600 ms to start and answer 404. T-2841 runs the fake with a 1 s start delay, which reproduced the missing 404 on every run, and allows 5 s. The other half of I-0329 is D-0558.

Sources: I-0325; S-0650 (`wip/archive/agents/S-0650.md`, Decisions); T-2836; I-0331, I-0329, S-0652 (`wip/archive/agents/S-0652.md`, Decisions), T-2841, T-2859.
