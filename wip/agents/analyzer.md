---
kind: analyzer
accrued_cost: 46.3509
accrued_seconds: 2989
tasks_completed: 4
last_run: 2026-11-27T22:48:56Z
---

# Analyzer activity

flai writes this document, one entry per activity, newest last; do not edit it by hand.

## Log

### 2026-10-15T22:49:08Z

- Summary: Wrote `design/analysis/2026-10-15-bottlenecks.md` (focus bottlenecks, window 2026-09-15 to 2026-10-15, active, indexed in `design/analysis/README.md`) with four findings and filed I-0150 (review waits of 259 h when the orchestrator's acceptance stops), I-0151 (pending recommendations waiting 25 to 113 h on the operator, S-0391 waiting now on TH-0461), I-0152 (folder-wide touches and `after` hubs holding up to 10 ready stories with 3 of 5 lanes free), and I-0153 (impression: ready stories unheld and unpulled 75 h each); `flai check` and the markdown lint pass.
- Trigger: asked
- Items: I-0150, I-0151, I-0152, I-0153
- Seconds: 454
- Cost: 5.0040 USD, estimated

### 2026-11-07T08:18:10Z

- Summary: Wrote `design/analysis/2026-11-07-all.md` (focus all, window 2026-10-08 to 2026-11-07, active, indexed in `design/analysis/README.md`) with 18 findings, filing I-0284 to I-0298 (confirmed cost-of-delay inputs waiting on `set_cost_of_delay` off, `after` chains holding ready stories 56 to 196 h a week, touches exact in 16 percent of stories, one-story merge-queue batches at 2.95 verifications per landing, three idle days on an all-draft backlog, guard refusing the orchestrator's `git tag` reads, a NUL byte in the serve journal dropping the Close-outs section, `publish` inert under judgement, no analysis schedule, stale design sentences, raw-HTML rendering in flaiover, the close-out job running worktree scripts on the host, a usage-limit refusal recorded as a failed `permission_prompt` check, tag-pinned CI inputs and two lockfiles, five files over 2,000 lines) and bumping I-0230, I-0246, I-0257, and I-0268; `flai check` and the markdown lint pass.
- Trigger: formula analysis-cadence S-0584
- Items: I-0230, I-0246, I-0257, I-0268, I-0284, I-0285, I-0286, I-0287, I-0288, I-0289, I-0290, I-0291, I-0292, I-0293, I-0294, I-0295, I-0296, I-0297, I-0298
- Seconds: 793
- Cost: 15.4957 USD, estimated

### 2026-11-16T06:17:09Z

- Summary: All links resolve, `flai check` reports nothing on the analysis folder, and the markdown lint passes with 0 errors. The report is `active`, indexed, and its nine issues are filed and linked. (report design/analysis/2026-11-16-all.md)
- Trigger: formula analysis-cadence S-0803
- Items: I-0386, I-0387, I-0388, I-0389, I-0390, I-0391, I-0392, I-0393, I-0394
- Seconds: 889
- Cost: 11.3293 USD, estimated

### 2026-11-27T22:48:56Z

- Summary: Report written: `design/analysis/2026-11-27-all.md` (focus all, window 2026-10-28 to 2026-11-27), indexed in `design/analysis/README.md`; it filed I-0411 (the curator's unlinted digest entries on main refused every acceptance for three days), I-0412 (666 orchestrator restarts in 14 h on an expired login, no hold or thread), I-0413 (a passed close-out or resolve job starts no orchestrator, so S-0839 and S-0843 waited six days), I-0414 and I-0415 (S-0811's three send-backs on one shared sentence and the `--keep-removal` only the operator can give; folded with S-0811's I-0408 and I-0410 at its acceptance), impressions I-0416 (MCP HTTP agent identity is a header), I-0417 (host config written 0644 with a notify URL that may carry a secret), and I-0418 (goreleaser hooks run tidy and tests with the signing key in the environment), and bumped I-0388 (13 close-outs interrupted by auto-upgrade since 2026-11-16) and I-0390 (nine more design sentences the code does not match); lint and `flai check` are clean, and the 26 uncommitted digest files on main, one failing MD028, are the operator's first item.
- Trigger: formula analysis-cadence S-0850
- Items: I-0388, I-0390, I-0411, I-0412, I-0413, I-0414, I-0415, I-0416, I-0417, I-0418
- Seconds: 853
- Cost: 14.5219 USD, estimated
