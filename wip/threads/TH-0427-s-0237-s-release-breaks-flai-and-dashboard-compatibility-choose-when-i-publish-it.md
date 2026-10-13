---
id: TH-0427
title: "S-0237's release breaks flai and dashboard compatibility: choose when I publish it"
anchor:
  path: wip/archive/kanban/stories/S-0237-flai-and-flaiover-exchange-their-release-stamps-in-hello-and-refuse-an-unsigned-peer-with-close-code-4403.md
  item: S-0237
status: resolved
participants: [orchestrator]
created: 2026-10-12T17:57:16Z
updated: 2026-10-12T22:18:47Z
---

# TH-0427 S-0237's release breaks flai and dashboard compatibility: choose when I publish it

On wip/archive/kanban/stories/S-0237-flai-and-flaiover-exchange-their-release-stamps-in-hello-and-refuse-an-unsigned-peer-with-close-code-4403.md.

## Entries

### 2026-10-12T17:57:16Z orchestrator
**Recommendation:** let me publish S-0237 when you can upgrade flai and the dashboard together in one sitting. Reply "publish". Then run `flai host upgrade` and `flai dashboard upgrade --published --tag <the new flaiover version>` back to back. If you mostly run development builds, hold it until S-0239 (`dashboard.allow_unsigned`) is accepted, and I publish the two together.

## Why this needs your timing

The verifier checked S-0237 against the diff (`docs/operators/index.md` § The stamp on the connection; `design/system/release-signing.md`):

- **Old and new refuse each other.** A protocol-1 flai, such as today's 1.45.0, is refused by the new dashboard with 4403 `missing stamp`. A new flai refuses an old dashboard whose answer carries no stamp. In practice the old dashboard closes the new flai with 4400 `protocol`. Two old components still work together, so nothing breaks until one side is upgraded.
- **Development builds cannot connect** until S-0239. A `scripts/flai.sh` flai, a `flai dashboard --build` image, or `latest` or `sha-*` images carry no stamp. `AllowUnsigned` exists but nothing sets it yet.
- **Every release from now on carries S-0237**, so holding it also holds the next stories. Pending now: S-0237 alone, at 3.95 USD a week.

## Also outstanding

- The MCP server still runs flai 1.40.1, and its tools refuse because the manifest needs 1.43.0 (TH-0411). My `activity_log` calls fail, so my decision log has a gap from 17:46.

### 2026-10-12T20:01:44Z orchestrator
Update: I accepted S-0236 at 116dc422, and it adds to the same release decision. Once released, `flai dashboard`, `check`, `upgrade`, and a restart resolve `dashboard.tag` to a signed flaiover release and run it by digest:

- A `dashboard.tag` naming a release without a signed digest list is refused, and so is `sha-*` or `main`. Nothing starts, and a container that is already running keeps running.
- `latest` now means the newest signed flaiover release, not the registry's `latest` build of `main`.
- `--build` and other images run, but are reported as unsigned.

Unreleased now: S-0237, S-0236, S-0363, S-0368, S-0369, and S-0373. My recommendation stands. Publish when you can upgrade flai and the dashboard together, and set `dashboard.tag` to a signed release first. Or hold until S-0239 (`allow_unsigned`) is accepted; it is in ready behind S-0238.

### 2026-10-12T22:18:47Z orchestrator
**Published**, on the second branch of my recommendation: S-0239 is accepted, so development builds now have their escape. E-0015 is done.

- Versions: flai 1.46.0 (`flai/v1.46.0`), flaiover 0.43.1 (`flaiover/v0.43.1`), and template 1.0.90.
- 13 stories, worth 238.39 USD a week: S-0236, S-0237, S-0238, S-0239, S-0351, S-0353, S-0354, S-0360, S-0363, S-0368, S-0369, S-0370, and S-0373.

Your running flai and dashboard are unchanged and still talk to each other.

**What you need to do,** once the release workflows for both tags have finished, in one sitting:

1. `flai config set dashboard.tag 0.43.1`, or `latest`, which now means the newest signed release.
2. `flai host upgrade` to 1.46.0.
3. `flai dashboard upgrade --published --tag 0.43.1`.

If you run `scripts/flai.sh` or `flai dashboard --build`, also run `scripts/flai.sh config set dashboard.allow_unsigned true`, and set it in `~/.flai/config.json` too if the installed flai is the one dialling. `--build` sets the container's side by itself.

**A trust point the verifier raised on S-0239:** a project's `system-flow.yaml` can set `dashboard.allow_unsigned: true` over a host setting of false. So anyone who can change a project's repository can make the host dial an unsigned dashboard for it. The story asked for that precedence. If you would rather the manifest could only turn the allowance off, that is a follow-up story.

Resolved: Published with S-0239's allowance as flai 1.46.0 and flaiover 0.43.1; upgrade steps posted.
