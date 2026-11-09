---
id: ADR-0270
title: "A permission_prompt check whose claude -p ends on the harness's usage limit records no outcome and opens no thread, and flai serve checks that version again after the limit resets, or an hour after the refusal when no reset time is named"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
refines: [ADR-0106]
topics: [permissions, usage-limit, flai]
---

# ADR-0270 A permission_prompt check whose claude -p ends on the harness's usage limit records no outcome and opens no thread, and flai serve checks that version again after the limit resets, or an hour after the refusal when no reset time is named

## Context

[ADR-0106](0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md) § 2 has `flai serve` check each new Claude Code version once against `permission_prompt`. A headless `claude -p` is asked to Write a file in a scratch story's worktree. It passes when the file holds what was asked, and fails otherwise. flai records the outcome in `claude-checks.json` beside its configuration. § 2 then says: "A version with a record is not checked again, whatever its outcome."

[I-0296](../issues/I-0296-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md) records what that rule does to a refusal for the account's usage limit. On 2026-10-26T21:55:54Z the check of Claude Code 2.1.311 was refused with "You've hit your weekly limit". The prompt was never reached, so the file was not written, and flai recorded `failed`. It opened a thread on `system-flow.yaml` saying protected-path writes would not go through, and it never checked 2.1.311 again. Protected-path writes under 2.1.311 were refused until Claude Code moved on. The only remedy the thread named was a hand edit of the host's `claude-checks.json`. TH-0669 and TH-0819 are still open on `system-flow.yaml`.

A usage-limit refusal says nothing about the prompt. Claude Code's stream-json says so as data, which the usage reader of [ADR-0214](0214-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-with.md) already reads: a `rate_limit_event` whose `rate_limit_info` has `status: rejected` and `resetsAt`, then a `result` with `is_error: true`, `api_error: usage_limit_reached`, and `api_error_status: 429`. The reader turns these into a neutral `Limit`: the error, its status, and the reset time when one is given.

## Decision

**A `permission_prompt` check whose `claude -p` ends on the harness's usage limit records no outcome and opens no thread, and `flai serve` checks that version again after the limit resets, or an hour after the refusal when no reset time is named.** This refines ADR-0106 § 2; the rest of it stands.

1. **What counts as a limit.** The check's run ended on the usage limit when the harness's usage reader, reading the check's `--output-format stream-json` output, gives its last run a `Limit`: a rejecting `rate_limit_event`, or a `result` naming the usage-limit error or status 429, as ADR-0214 reads it. flai does not match Claude Code's text.
2. **No outcome.** flai writes no record for that version and provider in `claude-checks.json`, and opens no thread. It logs one warning, `claude code check refused on the usage limit`, with the version, the provider the check ran through, the image when the check ran in one, and the reset time when Claude Code gave one.
3. **Checked again after the reset.** The first look after the reset checks the version again through that provider. When the refusal names no reset time, the version is checked again an hour after the refusal, so a limit with no reset does not start a refused `claude -p` at every look.
4. **Host and image alike.** The check of the host's `claude` and the check inside an agent image (S-0618) follow the same rule, each held by its own version, provider, and image.
5. **A restart checks again.** The hold until the reset is kept in `flai serve`'s process. A `flai serve` started again has no record of the version, and so checks it at its first look. A check refused again is held again.
6. **Every other failure stays final.** A check that reached Claude Code and did not write the file, for any cause but the usage limit, is recorded `failed` and told on a thread, and is not checked again, as ADR-0106 § 2 decided.
7. **Records made before.** A `failed` record written for a usage limit before this change stays until the operator removes its entry from `claude-checks.json`. The host's state is the operator's.

## Consequences

- In I-0296's case, Claude Code 2.1.311 would have had no record and no thread, and would have been checked again at the first look after the weekly limit reset.
- A limit costs a refused `claude -p` at most once per reset, or once an hour with no reset time, for each version and provider.
- A thread about a failed check now always means the prompt itself failed, or the check could not write for a cause other than the limit.
- `design/system/agent-adapters.md` § Permissions and `design/system/flai-cli.md` describe the rule, and the user and operator docs say a usage-limit refusal is checked again.
- The check's refusal does not feed the host's budget state of [ADR-0209](0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md) § 5. The check is not an agent run.

## Alternatives considered

- **Record a `limited` outcome with its reset time in `claude-checks.json`.** It would survive a restart and show the operator why the version is unchecked. But every reader of the file today treats any record as done, and a restart that checks again costs one refused run at most. The log line says why.
- **Match Claude Code's text, such as "You've hit your weekly limit".** It depends on Claude Code's wording, locale, and time zone. The stream-json fields need neither, and ADR-0214's reader already reads them.
- **Feed the refusal into the host's budget state (ADR-0209 § 5).** It would hold every start on the provider. The check runs on its own model and account settings, apart from any story, and holding every agent on it is a wider change than this issue.
- **Check again at every look.** A limit with no reset time would start a refused `claude -p` every minute until it lifts.
