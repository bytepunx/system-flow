---
id: T-2217
type: task
nature: remediation
title: The users guide, the operators guide, and agent-adapters.md describe harness none, and I-0214 is closed saying what fixed it
status: backlog
parent: S-0516
owner: alex
created: 2026-10-28T13:14:28Z
updated: 2026-10-28T13:14:28Z
transitions: []
stream: S-0516
tags: [cli, serve, docs]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/issues/I-0214-flai-serve-starts-a-host-agent-for-a-story-meant-to-be-worked-in-a-container-since-nothing-on-a-story-or-the-host-says-its-agent-runs-elsewhere.md, design/issues/summary.md]
after: [T-2215, T-2216]
---
# T-2217 The users guide, the operators guide, and agent-adapters.md describe harness none, and I-0214 is closed saying what fixed it

## Work

Describe the reserved name where the harnesses are described, in the same words T-2213's ADR uses. `docs/users/flai.md` § Who works a story: its agent: `harness: none` says the story's agent runs elsewhere, by hand or in a container, and `flai serve` starts nothing for it, so start it with `flai story start`. `docs/operators/index.md` § Starting an agent when a story becomes ready: the **Harnesses** bullet names `none` beside `claude-code`; the **Why a story waits** bullet adds the reason; the **Restart**, **Start**, and **Commit** bullets say each refuses such a story; the **Automatic restart** bullet says it leaves one alone. `design/system/agent-adapters.md` § Today › Starting and resuming an agent: the `Adapters` row says `For` refuses `none` and why, and the paragraph on the manifest's shape check says where the name is decided, linking the ADR. Bump each document's `updated`. Then close the issue from the story's worktree with `flai issue close I-0214 --reason "..."`, the reason naming `harness: none`, the ADR, and that `flai serve` starts nothing for such a story and does not restart it, which rewrites `design/issues/summary.md`. Waits for T-2215 and T-2216, so that the docs describe what was built.

## Done when

- The three documents describe `harness: none`, what `flai serve` does and refuses for it, and how a story with it is started, with `updated` bumped, and the markdown lint passes on them
- I-0214 is closed with a reason that says what fixed it, `design/issues/summary.md` no longer lists it as open, and `flai check --strict` is clean

## Notes

Layer 4, the last. `docs/users/flai-reference.md` and `docs/operators/settings.md` are not touched: no flag and no setting changes, since the name goes on the story's agent with `flai edit --harness none` or `flai agent set --harness none`, which take any name of the right shape today.
