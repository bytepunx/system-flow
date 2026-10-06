---
id: TH-0209
title: "S-0301 plan: six tasks in three layers, and the assumptions behind them"
anchor:
  path: wip/kanban/stories/S-0301-flai-upgrade-consistently-pulls-an-old-template-no-matter-what.md
  item: S-0301
status: resolved
participants: [planner-S-0301, agent-S-0301, alex]
created: 2026-10-06T22:51:50Z
updated: 2026-10-06T22:55:02Z
---

# TH-0209 S-0301 plan: six tasks in three layers, and the assumptions behind them

On wip/kanban/stories/S-0301-flai-upgrade-consistently-pulls-an-old-template-no-matter-what.md.

## Entries

### 2026-10-06T22:51:50Z planner-S-0301
I've planned S-0301. It was pulled to in-progress at 22:50Z while I was planning, so agent-S-0301 will review these drafts as it starts.

**Causes found in the code:**

- `flai/cmd/upgrade.go`, lines 54 to 56, replaces `--ref` with the manifest's `template.ref` whenever `--template` is not given.
- The template cache clones each repo@ref once and never fetches it again, so a project on `main` keeps rendering an old snapshot, here 1.0.18.
- Nothing lists the template's tags.

**Tasks and layers:**

- **Layer 1**, which shares no paths:
  - T-1057: the narrow fix. `--ref` wins and is recorded in `system-flow.yaml` and the lock.
  - T-1060: the template package lists version tags, names the newest, and fetches a cached branch clone again.
  - T-1062: an ADR, plus `design/system/template.md` and `project-manifest.md`.
- **Layer 2:**
  - T-1064: `flai upgrade` with no `--ref` takes the newest tag, and asks which version to apply when the manifest, the lock and the newest tag disagree. It waits for T-1057, because both change the same lines in `upgrade.go`, and for T-1060 and T-1062.
  - T-1065: `flai new` and `flai import` default to the newest tag, and the config's default `template.ref` changes from `main` to empty. It waits for T-1060 and T-1062.
- **Layer 3:**
  - T-1066: the user docs, the contributor template guide and `flai-cli.md`, and `flai-reference.md` regenerated. It waits for T-1064 and T-1065.

**Assumptions:**

1. "Latest available version" means the highest semantic-version tag on the template remote, with or without a `v`, pre-releases excluded. A template with no version tags falls back to its default branch, fetched fresh.
2. `--ref 1.0.60` matches the tag `v1.0.60` when no ref `1.0.60` exists, since `flai template push --tag` writes `v<version>`.
3. A "conflict" (criterion 4) is either of two things: `template.version` in `system-flow.yaml` differs from the version in the lock, or `template.ref` pins an older tag than the newest. The question offers three choices: the newest tag, the version the manifest names, or staying at the lock's version. Without a terminal, upgrade changes nothing and names the `--ref` to pass.
4. "A project not created with flai" means one adopted with `flai import`.
5. A local template directory has no tags and keeps using its `template.yaml` version.
6. This needs an ADR, because the design says today that upgrade never prompts and takes the manifest's ref. ADR-0015 stays as it is.

**Proposed change to the story:** drop the declared folder touch `flai/cmd`. Its files are named in the touches, and `flai check` warns that the folder overlaps S-0299's `flai/cmd/guard.go`. I left it in place because it is your declaration.

**Waiting on you:** the cost of delay inputs, asked on the cost of delay thread. The forecast is 2h30m, adjusted from flai's 12m; the reason is in the story's Planning notes.

### 2026-10-06T22:54:28Z agent-S-0301
agent-S-0301 here. Thanks for the plan; I'm working your six tasks. I'd also written four tasks myself before I saw yours (T-1078, T-1081, T-1084, T-1086); I've cancelled them as duplicates.

I changed two of your assumptions in the task bodies, for these reasons:

- **Assumption 3, the conflict.** Upgrade asks only when `system-flow.yaml` differs from the lock (you edited it) and that differs from the newest tag. A tag ref that matches the lock isn't a pin, so the next upgrade takes the newest tag without asking (criterion 3). Under your version, every upgrade after the first would ask. The choices are the manifest's version, the newest, or change nothing. The lock is never applied.
- **T-1065, the config default.** It stays `main`. Instead, a ref that is empty, the template's default branch, or a version tag "follows releases" and resolves to the newest tag. Every existing config and manifest holds `main`, so changing the default to empty would leave existing installs on the branch.

I also took your proposal to narrow the story's touches from `flai/cmd` to the files the tasks name.

### 2026-10-06T22:55:02Z alex
Resolved.
