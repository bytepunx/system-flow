---
id: ADR-0110
title: "A story's agent runs the close-out, which runs flai verify, itself before review, and hands the verifier only the review of the diff against the acceptance criteria and the conventions"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0059]
---

# ADR-0110 A story's agent runs the close-out, which runs flai verify, itself before review, and hands the verifier only the review of the diff against the acceptance criteria and the conventions

## Context

[ADR-0059](0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md) gave the story's agent a verifier sub-agent that runs the project's tests, lint, and `flai check`, and checks the diff against the story's acceptance criteria and the conventions. Its reason was context: a test or lint log read in the story's agent's context stays there and is read again on every later turn. Before review, the conventions and the harness prompt had one fresh verifier run the close-out, and one more after each fix.

That run cost more than it saved. Verifiers spent 88 minutes over 22 stories running `scripts/close-out.sh` and reading its log, S-0248 ran the same two-minute script three times, and S-0266 had to teach the verifier to run it once, without a pipe, and to read its last line. Each fix took a fresh verifier to confirm. The verdict the story's agent needed was a few lines: which step stopped, and what it found.

S-0273 gave flai the project's test tiers and `flai test`, which answers pass or the first findings for the paths given. S-0270 adds `flai verify S-nnnn`, the MCP tool `verify`, and the host methods `verify.run` and `verify.status`. In the story's worktree they run, cheapest first: no rebase left unfinished, the branch contains the main branch, the narrative's `## Current state` and `## Next steps` are written, `flai check --strict` scoped to the story, then each test and lint tier the branch's diff selects. They answer each step's state, duration, and findings, the commit verified, and a last line naming the outcome and the step it stopped at. Each run is stored; `flai verify S-nnnn --last` prints it again, and the review page shows it. A finding outside the story is a note, as [ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) made it, recorded as an issue with `--record-issues`. `scripts/close-out.sh` now runs `flai verify "$story" --record-issues` for its checks, and keeps only the commit and the sync check after it.

The run's answer is now as short as a verifier's report, so handing it to a sub-agent no longer keeps anything out of the story's agent's context.

## Decision

A story's agent runs the close-out, which runs `flai verify`, itself before review, and hands the verifier only the review of the diff against the acceptance criteria and the conventions. This refines ADR-0059's hand-off.

1. **The run is the agent's.** Before review, after it commits and syncs, the story's agent runs the project's close-out script once, in the worktree, as one command with the longest timeout the harness allows and its exit status echoed after it, never piped or redirected into a file. Where the project has no close-out script, it runs `flai verify S-nnnn`, or the MCP tool `verify`. It reads the last line and the findings of the step that stopped, fixes what they name, commits, and runs it again. It never hands the run, or a run after a fix, to a sub-agent. A passing run is the story's run before review, and `flai verify S-nnnn --last` prints it again.
2. **The close-out runs `flai verify`.** `work-management.md` names `flai verify S-nnnn --record-issues` as what the close-out runs for its checks, before its commit and its sync check. The close-out and `flai verify` never disagree, because the close-out has no checks of its own.
3. **The verifier reviews.** The verifier checks the diff against the story's acceptance criteria and the conventions, and names the criteria it meets by number. The story's agent hands it that review when it pays, as when the diff is too large to read in its own context, and does it itself otherwise. The verifier reads the last result with `flai verify S-nnnn --last` and does not run the tests, the lint, or the check to learn whether they pass. It runs the close-out, or `flai verify`, only when its prompt asks.
4. **The orchestrator verifies at the branch head.** Under `accept_reviews`, the orchestrator reads the story's last result with `flai verify S-nnnn --last --json`, and runs `flai verify` when that result did not pass or its commit is not the head of the story's branch. A run that did not pass is a blocker. It then hands the worktree and the result to its verifier, which checks the diff against each criterion and names the commit it verified, as [ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md) requires. ADR-0093's conditions are unchanged.
5. **Where it is said.** `delegation.md`, `work-management.md`, the story's agent's and the orchestrator's prompts in `flai/internal/harness`, and the verifier's definition, `.claude/agents/verifier.md`, with their template copies, say the same.

## Consequences

- No sub-agent is started to run the close-out, and none to confirm a fix. The story's agent reads a few lines per run, not a log.
- The story's agent waits for the close-out in its own turn, up to its tool's timeout. A close-out longer than that is cut off, as a verifier's was.
- A verifier is started only for the review, and only when the review pays. A small story is reviewed by its own agent, against the criteria it has been ticking task by task.
- `flai guard` lets a sub-agent and the orchestrator run `flai verify`, `--last` included, and call the MCP tool `verify`, which store only flai's cache; `flai verify --record-issues`, which writes issues, stays the story's agent's.
- ADR-0085's notes and issues are unchanged: `flai verify --record-issues`, run by the close-out, records them, and the close-out commits them.
- A project made from the template needs the flai release that ships `flai verify`, since its close-out calls it.

## Alternatives considered

- **Keep the verifier for the run, now that the run is short.** It keeps the cost of starting a sub-agent and of a second one to confirm each fix, for an answer the agent can read itself.
- **The agent runs `flai verify`, and the close-out keeps its own checks.** The two would disagree, as the close-out and the verifier's steps did before S-0270.
- **Drop the verifier.** A review of a large diff against the criteria and the conventions still fills the story's agent's context; the verifier keeps that out when it pays.
- **Always have a verifier review before review.** For a small story the agent has already checked each task against its criteria; a second reading adds a sub-agent's start and little else.
