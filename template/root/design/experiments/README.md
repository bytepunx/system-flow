# experiments

What experiments found: one document per `experiment` story, named for the story, `<S-nnnn>-<slug>.md`. Each says the hypothesis, the success measure, what was done, the results, and a recommendation to adopt, adapt, or drop what was tried. Start one from `template.md`.

An experiment story is accepted only with its results document committed on its branch; it releases nothing, whatever it touched. `flai check` validates every document here. A recommendation to adopt or adapt becomes a story or an ADR; this folder records what was learned, and `../system/` says how the system is. Convention: `../conventions/git.md` and `../conventions/documentation.md`.
