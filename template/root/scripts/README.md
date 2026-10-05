# scripts

Purpose-named shell scripts for common tasks. The `Makefile` calls these; CI calls the same ones. Each script is safe to run from any directory, prints what it does, and exits non-zero on failure. Scripts are POSIX `sh` under `set -eu`, because the host shell may be zsh.

| Script | Does |
|--------|------|
| `check.sh` | Runs `flai check --strict` on this repository; in a close-out, which exports `CLOSE_OUT_STORY`, scoped to the story with `--story` and `--record-issues`, so findings outside it are notes recorded in issues |
| `test.sh` | Behavior tests, fast, every iteration (`make test`) |
| `integration.sh` | Integration tests against real adapters (`make integration`) |
| `smoke.sh` | End-to-end smoke tests (`make smoke`) |
| `lint-md.sh` | markdownlint-cli2 over every markdown file with `.markdownlint.yaml` (`make lint-md`) |
| `close-out.sh` | Before a story goes to review, in its worktree: lint, the three test tiers, `flai check --strict` scoped to the story, the narrative's `## Current state` and `## Next steps`, then the commit with the `git commit` options given; stops at the first step that fails, and every run ends with one line naming the story, the outcome, and the step it stopped at, so nobody runs it again to learn why it stopped |
