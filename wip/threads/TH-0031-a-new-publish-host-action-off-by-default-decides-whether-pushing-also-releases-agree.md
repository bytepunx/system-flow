---
id: TH-0031
title: A new publish host action, off by default, decides whether pushing also releases. Agree?
anchor:
  path: wip/kanban/stories/S-0144-automatic-publish-push-should-be-a-setting.md
  item: S-0144
status: resolved
participants: [agent-S-0144, alex]
created: 2026-09-29T01:14:05Z
updated: 2026-09-29T02:06:11Z
---

# TH-0031 A new publish host action, off by default, decides whether pushing also releases. Agree?

On wip/kanban/stories/S-0144-automatic-publish-push-should-be-a-setting.md.

## Entries

### 2026-09-29T01:14:05Z agent-S-0144
**Cause.** Since S-0094, `flai push --pending` tags a release of everything accumulated before it pushes. Agents are told to run it whenever `inbox` reports an unpushed acceptance, and the board's Push now button runs it too. So every acceptance ends up published, one release per story.

**What I am doing (recommendation A):**
- `flai push --pending` pushes the merged commits and any tags already made. It tags nothing new unless a new host action, `publish`, is enabled for the project (`flai serve enable publish`, or its toggle in the dashboard's settings panel). It is off by default. Enabling `push` no longer implies any release.
- Publishing stays a deliberate step: the board's Publish button (still gated by `push`, as ADR-0032 says) or `flai release --pending`.
- Agents keep pushing plain merged commits after an acceptance. That only puts what you accepted on the remote; it cuts no release, so batching is kept.

**B:** also stop agents pushing after an acceptance, so nothing reaches the remote until you push or publish. Say B if "pushed" in the first criterion means that. It changes the work-management convention and the MCP instructions, which S-0138 is also editing.

**Name:** `publish`, unless you prefer something else, such as `auto-publish`.

I will record this in an ADR that refines ADR-0032 and supersedes S-0094's release-at-push.

### 2026-09-29T02:01:31Z alex
A: Call it `auto-publish`. Will the tagging be affected such that previously pushed commits end up missing version changes as a result of this change?

### 2026-09-29T02:01:48Z agent-S-0144
A, named `auto-publish`: renaming it now.

**No version change goes missing.** What a release covers is worked out from local git history: every "accept and archive" commit since each component's last release tag. Whether those commits were pushed plays no part. With `auto-publish` off, a pushed acceptance stays pending. The next Publish or `flai release --pending` bumps and tags it with everything else accumulated, at the highest bump among them, and pushes the version commit and tags.

A test already checks that the release is still pending after such a push. I am extending it to publish afterwards and check that the tag and version commit reach the remote.

### 2026-09-29T02:04:40Z alex
Proceed as planned.

### 2026-09-29T02:06:11Z agent-S-0144
Resolved: A, named auto-publish (ADR-0048, accepted). Answering the question turned up a bug: a publish made after its acceptances were pushed tagged locally and never pushed the tag. It is fixed in 385b443: such tags are now pushed on their own.
