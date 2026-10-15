---
id: ADR-0161
title: "flai serve's upgrade thread names a flai release only once it is published, and asks again on the next minute's look until it is"
status: proposed
date: 2026-10-15
supersedes: []
superseded_by: []
refines: [ADR-0158]
---

# ADR-0161 flai serve's upgrade thread names a flai release only once it is published, and asks again on the next minute's look until it is

## Context

[ADR-0158](0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md) point 7 has `flai serve` open one thread per project and release asking the operator to run `flai host upgrade --version X.Y.Z`, naming the newest `flai/vX.Y.Z` tag in the project's history. It opens that thread with `auto-upgrade` off, with it on while `flai serve` runs under no host, and after an hour of refused asks.

A release is tagged minutes before its binaries are published: `flai/v1.52.0` was tagged at 05:02Z and published about 05:08Z. The host refuses a release that is not published, so a thread opened in that window names a command that fails. The record under the serve directory then holds that release, so no later look asks again, and the operator is never told when the command would work. [I-0149](../issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md) records the host left behind.

## Decision

`flai serve` opens the upgrade thread for a release only once that release is published.

1. **The check.** Before it opens the thread of ADR-0158 point 7 with `auto-upgrade` off, or on while `flai serve` runs under no host, `flai serve` asks whether the newest release is among the published flai releases, as `flai self-upgrade --list` reads them: the GitHub releases tagged `flai/vX.Y.Z` that are neither drafts nor prereleases, from `FLAI_RELEASES_API` when set, else api.github.com. One look is bounded at 30 seconds.
2. **Not published.** It logs at info that the thread waits for the release to be published, and neither opens the thread nor records it. The next minute's look asks again.
3. **Published.** It opens and records the thread as ADR-0158 says, and remembers the release as published for as long as it runs, so it does not look again for it.
4. **Cannot tell.** When the releases cannot be listed, it logs at warn and opens the thread as before: a thread naming a release not yet published beats no thread.
5. **Stopping.** While `flai serve` stops, it opens nothing.

## Consequences

- The thread's `flai host upgrade --version X.Y.Z` names a release the host can install, and the one thread per release ADR-0158 records is not spent on a command that fails.
- In the minutes between tag and publication, each look lists the published releases once per project behind, one GitHub API call each.
- Without a network the thread opens as before, and may name a release the host refuses.
- The auto-upgrade path is unchanged: its asks are refused as not published and retried every 15 minutes, as ADR-0158 point 6 says. The thread it opens after an hour of refusals is not held: it names the host's last refusal, which says why the release is not installed.
- A release whose publication never completes gets no thread with `auto-upgrade` off. `flai serve` logs that it waits on each look, and the inbox's `flai_outdated` still names the newer release.

## Alternatives considered

- **Name the newest published release instead.** It may be a minor version behind the project's history, and ADR-0142 still refuses `task_done` under it.
- **Tell the operator to wait for the release in the thread.** The record means no second thread, so the operator is never told when the command works.
