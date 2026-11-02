---
id: T-2423
type: task
nature: improvement
title: The tests that read the strategic agents' .claude/agents copies read their generated definitions instead
status: done
parent: S-0553
owner: alex
created: 2026-11-02T00:52:08Z
updated: 2026-11-02T00:53:22Z
transitions:
  - to: ready
    at: 2026-11-02T00:52:45Z
    by: agent-S-0553
  - to: in-progress
    at: 2026-11-02T00:52:50Z
    by: agent-S-0553
  - to: done
    at: 2026-11-02T00:53:22Z
    by: agent-S-0553
stream: S-0553
tags: []
touches: [flai/internal/harness/harness_test.go, flai/internal/pilot/run_test.go, flai/internal/serve/plan_test.go]
after: [T-2347]
usage:
  source: log
  seconds: 32
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1442
      cache_read: 355109
      cache_write: 7534
      cost: 0.1489
---
# T-2423 The tests that read the strategic agents' .claude/agents copies read their generated definitions instead

## Work

`TestTheTemplatesDefinitionsRead`, `TestTheOrchestratorsDefinitionSaysWhatEachPermissionDoes`, and `TestTheTemplatesAnalyzerIsHeldToItsReport` in `flai/internal/harness/harness_test.go` read the planner's, orchestrator's, analyzer's, and curator's `.claude/agents/<name>.md`, here and in the template, which T-2352 replaces with links and removes. They read the definition flai generates from the role files instead (ADR-0236), which is what runs since T-2347. Three fixtures that still write a `.claude/agents/planner.md` no start reads any more go: in `TestClaudeCodeStartsAStrategicAgentOnItsProvider`, `flai/internal/pilot/run_test.go`, and `flai/internal/serve/plan_test.go`.

It waits for T-2347, whose generation it reads. Written while T-2352 waits on the operator, so the copies can go without breaking the tests.

## Done when

The three tests pass reading the generated definitions, with or without the copies, and no test writes a strategic agent's `.claude/agents` fixture that nothing reads.

## Notes

Added by S-0553's agent on 2026-11-02, from what T-2347's sub-agent found was left.
