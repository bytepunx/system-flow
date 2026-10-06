---
id: T-0905
type: task
nature: feature
title: The orchestrator's prompt, agent file, and convention say how it reviews and accepts a story
status: done
parent: S-0221
owner: alex
created: 2026-10-05T04:47:45Z
updated: 2026-10-06T11:24:16Z
transitions:
  - to: ready
    at: 2026-10-06T11:18:33Z
    by: agent-S-0221
  - to: in-progress
    at: 2026-10-06T11:18:33Z
    by: agent-S-0221
  - to: done
    at: 2026-10-06T11:24:16Z
    by: agent-S-0221
stream: S-0221
tags: [flai, template]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml]
after: [T-0898]
usage:
  source: log
  seconds: 343
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 65
      output: 306
      cache_read: 2190811
      cache_write: 82943
      cost: 0.9782
---
# T-0905 The orchestrator's prompt, agent file, and convention say how it reviews and accepts a story

## Work

S-0218 writes the orchestrator's prompt in `flai/internal/harness` and its agent file. Add what it does with a story in review when `accept_reviews` is on:

1. Hand the story's worktree to its verifier sub-agent. The verifier runs the tests, the lint, and `flai check --strict`, and checks the diff against each acceptance criterion.
2. Run `flai accept <S-nnnn> --dry-run` with the verified commit, and read the blockers.
3. With no blocker, and every criterion matched to changed files, run `flai accept` with `--by orchestrator`, the verified commit, and the evidence. Log the acceptance with `activity_log`, naming the story and the evidence.
4. Otherwise, leave the story in review. Open a thread on it saying what is missing: each blocker, and each criterion it could not check against the diff. Log that decision.

Write the same rule into the convention's "As the orchestrator", template first:

- `template/root/design/conventions/strategic-agents.md`, then `design/conventions/strategic-agents.md`.
- The agent file, `template/root/.claude/agents/orchestrator.md`, then `.claude/agents/orchestrator.md`.
- Bump `template/template.yaml` and add a `template/CHANGELOG.md` entry.

A write under `.claude/` asks the operator through the permission prompt; answer its thread.

Test the prompt in `harness_test.go`: it names the verifier step, the dry run, the evidence, and the thread when something is missing.

This task waits for T-0898, which fixes the conditions and the evidence.

## Done when

- The orchestrator's prompt, its agent file, and the convention's "As the orchestrator" all give the same review steps, in the template and in this project alike.
- The template's version and changelog record the change.
- The harness test, `scripts/flai-test.sh`, `scripts/lint-md.sh`, and `flai check --strict` pass.

## Notes
