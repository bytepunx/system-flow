---
id: TH-0479
title: Which of S-0413's nine prose steps become stories?
anchor:
  path: design/system/agent-prose.md
status: open
participants: [orchestrator, alex]
created: 2026-10-16T06:07:54Z
updated: 2026-10-31T07:55:23Z
---

# TH-0479 Which of S-0413's nine prose steps become stories?

On design/system/agent-prose.md.

## Entries

### 2026-10-16T06:07:54Z orchestrator
Recommendation: draft stories for steps 1 to 3 now (measure in flai stats, a word budget per artefact in the prompt and the communication convention, and flai refusing writes over budget), under a new epic for plainer agent prose. Leave steps 4 to 9 until step 1's figures are in. Steps 1 and 3 need nothing, and step 2 needs a template release. The document puts the rewriter (step 9) last and makes it depend on those figures.

Why you: I accepted S-0413 under accept_reviews at 55882ab. Its last section, "Recommendation", says "The operator chooses at review which of these become stories. None is created by this story." Because I accepted it, you never saw that choice at review. Which stories to create is a scope decision, so it is yours.

What I need: name the steps to turn into stories, or say "none". Once there is an epic, I will start its planner under plan_backlog_epics.

### 2026-10-31T07:48:27Z alex
provide the steps in question - referring to them by number only means nothing to me without the reference

### 2026-10-31T07:55:23Z orchestrator
These are the nine steps from `design/system/agent-prose.md` § Recommendation, which S-0413 wrote. None has become a story yet.

| Step | What the story would do | Needs |
|------|-------------------------|-------|
| 1. Measure | `flai stats` reports words per artefact, words per sentence, and prose alerts per 1,000 words | nothing |
| 2. Prompts and conventions | The harness prompt and the communication convention give each artefact a word budget, with one good and one bad example | a template release |
| 3. Budgets | flai refuses a thread entry, log entry, Current state, or title over its word budget, and names the budget | nothing |
| 4. The linter | `flai test` runs Vale with the Flow style on what a story changes, as an advisory prose tier | Vale, pinned |
| 5. Self-revision | flai refuses a `wip` write with prose-lint errors, and returns each finding with its line | step 4 |
| 6. Reader summaries | The review and activity pages show each entry's first sentence and fold the rest | nothing |
| 7. Read less | The inbox in a tool's answer lists only the threads on the agent's own story, and the conventions drop what they repeat | a template release |
| 8. Structured threads | `thread_open` takes the question, the recommendation, and the options as fields | an ADR |
| 9. Rewriter experiment | An experiment: a hosted small model rewrites `wip` text that fails the prose lint, keeping every ID and path | steps 1 and 4 |

The document runs step 9 only if steps 1 to 5 leave the alerts above target, or retries cost more than a hosted rewrite would.

My suggestion stays as it was: draft stories for steps 1, 2, and 3 under a new epic for plainer agent prose, and leave steps 4 to 9 until step 1's figures are in.

What I need: the step numbers to turn into stories, or "none". Once there is an epic, I start its planner under plan_backlog_epics.
