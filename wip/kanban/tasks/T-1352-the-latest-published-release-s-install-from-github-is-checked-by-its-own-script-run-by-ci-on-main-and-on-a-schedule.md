---
id: T-1352
type: task
nature: remediation
title: The latest published release's install from GitHub is checked by its own script, run by CI on main and on a schedule
status: backlog
parent: S-0340
owner: alex
created: 2026-10-08T08:05:51Z
updated: 2026-10-08T08:06:10Z
transitions: []
stream: S-0340
tags: [cli]
touches: [scripts/install-published-test.sh, ".github/workflows/install-published.yml", Makefile]
---
# T-1352 The latest published release's install from GitHub is checked by its own script, run by CI on main and on a schedule

## Work

Keep today's GitHub check as `scripts/install-published-test.sh`. It does what `scripts/install-test.sh` does at the story's base commit: it installs the latest published release with `install.sh` and with `flai self-upgrade`, with a token from `GITHUB_TOKEN`, `GH_TOKEN`, or `gh`. Start from `git show <base>:scripts/install-test.sh`, since T-1354 rewrites that file. No close-out runs it: it is in no tier of `system-flow.yaml`'s `tests`, and `scripts/smoke.sh` does not call it.

- When it fails, it records the failure through `flai issue`. It bumps the open issue titled "The latest published flai release does not install from GitHub", or opens it with `flai issue new`, naming the failing step and its output. Then it exits non-zero.
- Add `.github/workflows/install-published.yml`, triggered on push to `main`, on a daily `schedule`, and on `workflow_dispatch`. It builds flai and runs the script with `GITHUB_TOKEN`. How an issue recorded in CI reaches the repository is the open question in the plan's thread on S-0340: follow its answer.
- Add a `make install-published-test` target beside `install-test`. Its help says it reaches GitHub.

Waits for no task: it starts from the script as committed, not from T-1354's rewrite.

## Done when

- `scripts/install-published-test.sh` passes on the host with a token.
- Run in a scratch clone with `FLAI_API` at an address that does not answer, it records the issue through `flai issue` there and exits non-zero.
- `install-published.yml` names the three triggers, and `actionlint` passes on it where it is installed.
- No tier in `system-flow.yaml` and no line of `scripts/smoke.sh` runs the script.

## Notes

Drafted by the planner from S-0340's criterion 4. The script name and the workflow name are the planner's; the story's agent may rename them, and then updates the touches.
