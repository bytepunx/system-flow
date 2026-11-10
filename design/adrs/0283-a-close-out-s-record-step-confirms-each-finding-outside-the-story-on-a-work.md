---
id: ADR-0283
title: "A close-out's record step confirms each finding outside the story on a work item under the checked-write lock and records only what still holds, and a story's containers mount the lock file"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0272]
topics: [cli, check, container]
---

# ADR-0283 A close-out's record step confirms each finding outside the story on a work item under the checked-write lock and records only what still holds, and a story's containers mount the lock file

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) has a close-out record each finding outside its story in an issue, one per rule. [ADR-0272](0272-a-checked-creation-holds-a-checked-write-lock-from-its-write-until-it-keeps-or.md) and [ADR-0275](0275-a-checked-edit-of-a-work-item-holds-the-checked-write-lock-from-its-first-save.md) give a checked creation or edit of a work item the checked-write lock, `.flai-cache/items-checked.lock`, from its first write until it keeps the item or puts it back. Only flai serve's look reads under that lock. Other readers read without it, and the scoped check of a close-out is one of them.

I-0316 shows what that costs. S-0599's close-out at 2026-11-09T07:36:46Z recorded S-0636's body as missing `## Tasks`. S-0636 was created at 07:36:47Z, and its first commit, `1a0ee3c7b`, has every section. The close-out's check read a creation that its check then refused and removed. The next creation wrote the item whole. The issue it recorded names a finding that never stayed on disk.

Inside a story's container the lock file is not mounted ([ADR-0253](0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md)). There a lock taken by the close-out's verify ([ADR-0274](0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)) is a file of its own that no checked write on the host holds.

## Decision

A close-out's record step confirms each finding outside the story on a work item under the checked-write lock and records only what still holds, and a story's containers mount the lock file.

1. **The record step confirms.** `flai check --story --record-issues` and `flai verify --record-issues` share it (`recordOutside` in `flai/cmd/check.go`). Before it records, it takes `(*workitem.Repo).LockCheckedWrites` and runs the checked write's rules, `check.RunItemWrite`, again on each active epic, story, or task file a finding outside the story names. A finding is kept when that run still gives one of its rule on its file. The rest are passed over, with a warning that counts and names them. Findings on other files are recorded as before.
2. **A lock not taken records as before.** When the lock cannot be taken within its wait, or the run fails, the step warns and records every finding unconfirmed. A close-out never stops on it.
3. **The lock is held only for the check run again.** It is let go before the issues are written.
4. **A story's containers mount the lock file.** The close-out's and the operator's checks' container (`ContainerCommand`), the story's agent's container (`Start`), and `scripts/agent-container.sh` mount the main checkout's `.flai-cache/items-checked.lock` read-write at its host path. Each makes it on the host as an empty 0600 file when it is absent, since Docker would make a missing bind-mount source a folder owned by root.

## Consequences

- A refused creation or edit no longer leaves an issue that names what never stayed. A finding that still holds is recorded as before.
- A close-out with findings on work items outside its story waits for any checked write in flight, a few seconds at most.
- An `flock` taken in the container excludes one taken on the host across the bind mount, both ways. This was measured on Docker Desktop 4.88.1 with engine 29.7.2 on WSL2, on tmpfs and ext4.
- A story's container mounts one more path of the main checkout. The lock file holds nothing.
- `RunItemWrite` runs every rule that reports on a work item's file, so a finding it no longer gives is gone, not unchecked.

## Alternatives considered

- **Hold the lock across the whole scoped check.** It blocks every checked write for the check's seconds, and the check outside the work items needs no lock.
- **Leave findings on work items out of the record.** That hides the real ones on items no write is changing.
- **Run the check again later.** A later read can land in another creation's window. It narrows the window and does not close it.
