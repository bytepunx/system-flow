# scripts

Purpose-named shell scripts for common tasks. The `Makefile` calls these; CI calls the same ones. Each script is safe to run from any directory, prints what it does, and exits non-zero on failure. Scripts are POSIX `sh` under `set -eu`, because the host shell may be zsh.

| Script | Does |
|--------|------|
| `check.sh` | Runs `flai check --strict` on this repository; with `CLOSE_OUT_STORY` set, as `flai verify` sets it for the tiers it runs, scoped to the story with `--story` and `--record-issues`, so findings outside it are notes recorded in issues |
| `test.sh` | Behavior tests, fast, every iteration (`make test`); the `test` tier of `system-flow.yaml`'s `tests` |
| `integration.sh` | Integration tests against real adapters (`make integration`); the `integration` tier |
| `smoke.sh` | End-to-end smoke tests (`make smoke`); the `smoke` tier |
| `lint-md.sh` | markdownlint-cli2 over every markdown file with `.markdownlint.yaml` (`make lint-md`); the `markdown` tier |
| `close-out.sh` | Before a story goes to review, in its worktree, on its branch: `flai verify <story> --record-issues`, which runs every check (no rebase left unfinished, the branch contains the main branch, the narrative's `## Current state` and `## Next steps`, `flai check --strict` scoped to the story, then the tiers of `system-flow.yaml`'s `tests` that the branch's changes select, `integration` and `smoke` on every story), so the close-out and `flai verify` never disagree; then the commit with the `git commit` options given, which also commits the issues `--record-issues` recorded, the check with `flai verify <story> --sync-only` that the branch still contains the main branch, which passes over commits on it that change only `wip/` paths the branch does not change, and the clean worktree check. It stops at the first step that fails, and every run ends with one line naming the story, the outcome, and the step it stopped at, so nobody runs it again to learn why it stopped. Add a check as a tier in `tests`, not as a step here |
