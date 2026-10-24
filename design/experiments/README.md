# experiments

What experiments found: one document per `experiment` story, named for the story, `<S-nnnn>-<slug>.md`. Each says the hypothesis, the success measure, what was done, the results, and a recommendation to adopt, adapt, or drop what was tried. Start one from `template.md`.

An experiment story is accepted only with its results document committed on its branch; it releases nothing, whatever it touched. `flai check` validates every document here. A recommendation to adopt or adapt becomes a story or an ADR; this folder records what was learned, and `../system/` says how the system is. Convention: `../conventions/git.md` and `../conventions/documentation.md`.

## Results documents

| Story | Document | Recommendation |
|-------|----------|----------------|
| S-0176 | [Planning a story's tasks in layers and working each layer with sub-agents at once](S-0176-the-story-s-agent-plans-which-tasks-can-run-in-parallel-and-works-them-with-sub-agents.md) | Adapt: keep `after` and the sub-agents, make parallel layers optional |
| S-0427 | [A SQLite or Postgres read model loaded from the files, for the chart data's broad reads](S-0427-sqlite-and-postgres-read-models.md) | Drop, for the charts |
| S-0435 | [Priming with the digest, against the threads opened to the operator per completed story](S-0435-digest-priming-and-threads-to-the-operator.md) | Adopt: keep the digest, turn the curator on, and measure again from 2026-11-09 (I-0207) |
