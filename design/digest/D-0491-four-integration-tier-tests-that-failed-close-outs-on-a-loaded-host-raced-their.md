---
id: D-0491
title: Four integration-tier tests that failed close-outs on a loaded host raced their own setup, so each now waits for the event it depends on rather than sleeping longer (I-0247)
kind: cause
sources: [I-0247, S-0569, T-2522, T-2525, T-2528, T-2529, I-0252, S-0617, T-2723, T-2724]
topics: [testing, flai, go, code]
paths: [flai/internal/serve/outdated_test.go, flai/internal/serve/serve_test.go, flai/internal/mcpserver/folder_test.go, flai/cmd/host_integration_test.go, flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, design/issues/I-0271-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md, design/issues/I-0273-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md, design/issues/I-0252-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0543]
---

# D-0491 Four integration-tier tests that failed close-outs on a loaded host raced their own setup, so each now waits for the event it depends on rather than sleeping longer (I-0247)

I-0247 recorded close-outs of S-0518 and S-0578 failing the integration tier on timing tests in packages neither story changes, at load averages up to 76 on 24 CPUs. S-0569 closed it on 2026-11-07. The issue's remediation says: "Each test raced its own setup; the load only widened the window. Each fix waits for the event it depends on."

- `TestWaitForEventsAcrossAFolder` and its two wait_for_work siblings slept 150 ms, then changed the project, before the wait had fingerprinted it. They now wait for the wait's first progress heartbeat and hold up to 30 s (T-2522).
- `TestServeSaysTheFlaiIsOlderThanTheProject` read the watcher's change notification for serve's upgrade thread as project.info's reply, "because channeltest's Ask returns the next message". It now asks until an answer names the project (T-2525).
- `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` moved the story before serve's watcher took its first snapshot. It now waits for the watcher to report a probe file (T-2528).
- `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` lost the first free MCP port to concurrent runs. It now seeds its own kernel-picked port (T-2529).

Two defects in flai's own code were found and left as issues, since fixing them changes non-test code beyond the story's touches: I-0271, the host restarting a project's MCP server on the same port, and I-0273, serve's watcher taking its first snapshot after its first looks. "The operator chooses at acceptance which of them become stories" (wip/archive/agents/S-0569.md, Decisions).

I-0252, filed separately for `TestServeSaysTheFlaiIsOlderThanTheProject` (project.info answered with no flai_outdated under load), was the same race. S-0617 closed it on 2026-11-09: "the one instance (S-0506, 2026-11-02, project.info: <nil>) was I-0247's symptom, a watcher notification read as the reply, which S-0569 fixed on 2026-11-06". T-2724 ran the three outdated serve tests 100 times each on 2 loaded CPUs since that fix, and none failed. T-2723 closed the one way left for flai_outdated to read nil while a newer flai is tagged: `release.FlaiOutdated` no longer keeps a failed git look for a minute, so one failed look is not served as not outdated (`TestFlaiOutdatedKeepsNoAnswerFromAGitThatFailed`). So an issue filed for one instance of a load failure is checked first against the race fixes that came after it, with a stress run, before new code is written.

A later test of serve, the host, or the MCP waits should wait for the event it needs (a heartbeat, a probe file the watcher reports, an answer that names its subject, a port of its own), not sleep. D-0413 and D-0427 record the same lesson from I-0226 and I-0233.

Sources: I-0247 and its `## Remediation`; S-0569 (wip/archive/agents/S-0569.md); I-0252 and its fix; S-0617 (wip/archive/agents/S-0617.md), T-2723, T-2724.
