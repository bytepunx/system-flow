---
title: Simplifying agent-produced prose
updated: 2026-10-16
status: active
topics: [cli, template]
---

# Simplifying agent-produced prose

The finding of S-0413. The operator asked for research on making the prose flai's agents write simpler and shorter, so that the operator and other agents read less and read it faster. The operator named two means: a natural language linter, such as Vale or proselint, and a small model hosted locally and fitted with LoRA or a similar method, which rewrites verbose agent text until it passes that linter. This document says what prose flai's agents write today and what checks it, then covers prose linters, a local rewriting model, and the alternatives, and ends with a recommendation. Nothing here is built.

Every figure about this repository comes from a command run on 2026-10-16. Design documents and ADRs are linked; code is named by its path.

## Today

Measured on 2026-10-16 over the committed tree at `53eb8e40`, with a Python script kept outside the repository. A word is a run of letters and digits outside front matter, fenced code, and HTML comments. An entry is a `###` section under a narrative's `## Log` or under a thread's or message's `## Entries`. Live and archived files are counted together.

| Kind | Path | Reader | Files | Words | Median words per file | Entries | Median words per entry |
|------|------|--------|-------|-------|-----------------------|---------|------------------------|
| Story narratives | `wip/agents/S-*.md`, `wip/archive/agents/` | The agent that resumes the story; the operator | 381 | 236,110 | 556 | 1,989 log entries | 20 |
| Strategic agents' activity logs | `wip/agents/planner.md`, `orchestrator.md`, `analyzer.md` | The operator; that agent's next session | 3 | 35,031 | 4,236 | 645 | 49 |
| Threads | `wip/threads/` | The operator and agents | 468 | 173,488 | 317 | 1,012 agent entries | 105 |
| Messages between agents | `wip/messages/` | The two stories' agents; the operator on escalation | 126 | 47,026 | 308 | 554 agent entries | 64 |
| Epics | `wip/kanban/epics/`, `wip/archive/kanban/epics/` | The operator and agents | 25 | 11,338 | 304 | | |
| Stories | `wip/kanban/stories/`, `wip/archive/kanban/stories/` | The operator and agents | 413 | 219,461 | 524 | | |
| Tasks | `wip/kanban/tasks/`, `wip/archive/kanban/tasks/` | Agents; the operator | 1,719 | 266,613 | 133 | | |
| Issues | `design/issues/I-*.md` | The operator, agents, the analyzer | 146 | 43,509 | 204 | | |
| ADRs | `design/adrs/` | The operator and agents | 162 | 107,983 | 639 | | |
| Design documents | `design/system/` | The operator and agents | 26 | 206,366 | 3,490 | | |
| Analyzer reports | `design/analysis/` | The operator | 1 | 1,983 | 1,983 | | |
| Commit messages | `git log`, commits with a Claude trailer | The operator | 1,955 | 117,132 | 59 | 1,411 with a body | 58 per body |
| Final report to the operator | Not stored in the repository | The operator | | | | | |

The kinds above hold 1,466,040 words. A small part is not agent prose: the operator's thread entries, flai's own notices in threads and messages, the operator's brief, and the goals the operator writes for stories.

- A narrative's largest section is `## Decisions`: 89,314 words across the 381 narratives, against 63,672 in `## Log` and 46,242 in `## Current state`, each counted with its heading.
- Threads hold 1,558 entries. Agents wrote 1,012 of them, 151,764 words. The operator wrote 401, 2,674 words, with a median of 1 word: 161 of them are `Resolved.` alone. flai wrote 145, 9,506 words.
- Messages hold 635 entries; flai wrote 81 of them.
- The repository has 2,954 commits. 1,955 carry a Claude trailer. Of the other 999, 278 carry flaiover's trailer and 721 carry none.
- The design documents include the operator's [brief](brief.md), 865 words, and this folder's README.
- The final report goes to the operator in the harness's session, not to a file. The harness keeps it in its session log, which flai reads for usage only (`flai/internal/usage`).

### Readability

Measured over the same tree with the same script. Sentences are split at `.`, `!`, or `?` followed by a space and a capital or a digit, within paragraph and list lines. Headings, tables, and code are left out, and an inline code span counts as one word. The Flesch-Kincaid grade uses a heuristic syllable count, so it compares rows rather than standing alone.

| Prose | Sentences | Mean words per sentence | Median words per sentence | Sentences over 25 words | Flesch-Kincaid grade |
|-------|-----------|-------------------------|---------------------------|-------------------------|----------------------|
| Narrative log entries | 4,150 | 13.7 | 11 | 13% | 6.8 |
| Narratives, whole | 12,746 | 16.9 | 15 | 20% | 7.5 |
| Strategic agents' log entries | 3,363 | 9.8 | 5 | 10% | 6.2 |
| Thread entries by agents | 10,277 | 12.9 | 12 | 6% | 5.5 |
| Message entries by agents | 1,857 | 13.7 | 12 | 12% | 5.2 |
| Epics | 484 | 22.4 | 18 | 33% | 10.5 |
| Stories | 10,631 | 18.0 | 16 | 20% | 8.1 |
| Tasks | 14,200 | 16.5 | 14 | 16% | 6.6 |
| Issues | 2,178 | 18.5 | 15 | 26% | 8.5 |
| ADRs | 4,868 | 20.7 | 19 | 28% | 8.9 |
| Design documents | 4,573 | 22.2 | 19 | 32% | 9.2 |
| Analyzer reports | 69 | 21.3 | 21 | 33% | 8.0 |
| Commit messages | 11,013 | 10.6 | 11 | 1% | 5.5 |

Sentences are shortest in commits, threads, messages, and logs. They are longest in the documents agents write for later reading: design documents, epics, ADRs, the analyzer's report, and issues. There a quarter to a third of sentences run past 25 words.

### What checks this prose today

| Check | What it checks | Enforced |
|-------|----------------|----------|
| [Communication convention](../conventions/communication.md) | Reports to the operator: lead with the outcome, one idea per sentence, short sentences with verbs, no restatement of the request, no praise | No. Agents read it in every prime pack. |
| [Documentation convention](../conventions/documentation.md) | Documents: write plainly, tables for parallel facts, short sentences, no filler, no marketing, no duplicated content | No. Agents read it in every prime pack. |
| flai's markdown lint, [ADR-0061](../adrs/0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md) | 23 markdownlint rules on the work items, threads, and narratives flai writes in `wip/`: headings, lists, whitespace, bare URLs, code fences, emphasis, tables | Yes. Package `flai/internal/mdlint`: `flai check` reports a finding, and a command that writes refuses one. |
| The `markdown` test tier in [system-flow.yaml](../../system-flow.yaml) | Every markdown file, through `scripts/lint-md.sh`, which runs markdownlint-cli2 0.20.0 with `.markdownlint.yaml` | Yes, in `flai test` and the close-out. Line length (MD013) is off. |
| Title limits | A conversation's title is cut to 72 runes (`flai/internal/messages`); the host API refuses an epic or story title over 200 characters (`flai/internal/hostapi`) | Yes, on titles only. |
| The prime pack's size budget, [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) | What an agent reads at story start: 81,920 bytes by default (`flai/internal/context`) | Yes, on what is read, not what is written. |

Nothing checks the prose itself. On 2026-10-16 a search of `flai/`, `flaiover/src/`, `template/`, `scripts/`, the manifest, and the repository's configuration found no check of sentence length, wordiness, passive voice, or readability. It found no prose linter such as Vale, proselint, write-good, or textlint. The two conventions state the style; agents follow them or not, and no tool measures the result.

The prime budget is the one place flai limits what an agent reads, and conventions alone exceed it. `flai prime --story S-0413 --json`, run on 2026-10-16, reports a pack of 128,217 bytes in 1,114 lines against the 81,920-byte budget, exceeded by conventions. The 14 conventions it loads whole come to 81,868 bytes, and their README to 3,063 more. The largest are [work-management.md](../conventions/work-management.md), 21,819 bytes, and [delegation.md](../conventions/delegation.md), 13,167 bytes.

## Prose linters

Criterion 1. This section compares the natural language linters that could check what agents write, records a trial of the strongest on this repository's prose, and recommends one. The linter would also score the training data and the output of a rewriting model (the next section), so machine-readable output, per-rule alerts, and rules the project can write matter as much as what a tool checks out of the box.

| Tool | Checks | Configurable rules | Markdown aware | Runtime | Output | Licence | Last release | Fit |
|------|--------|--------------------|----------------|---------|--------|---------|--------------|-----|
| [Vale](https://github.com/vale-cli/vale) | Whatever its styles define: word lists, sentence length, readability formulas, capitalisation, consistency | Yes. Each rule is a YAML file built on one of 12 [check types](https://docs.vale.sh/checks/), turned on, off, or re-levelled in `.vale.ini` | Yes. It skips code blocks, code spans, and URLs ([markdown](https://docs.vale.sh/formats/markdown)). It lints front matter unless `BlockIgnores` excludes it | One static Go binary, 56 MB. Reads files or stdin (`--ext=.md`) | `--output=JSON`, line, CLI, or a template ([CLI](https://docs.vale.sh/cli)) | MIT | v3.24.0, 2026-10-01 | Best |
| [Harper](https://github.com/Automattic/harper) | Grammar and spelling: about 900 rules, mostly fixed-phrase corrections, with `LongSentences`, `Hedging`, and `FillerWords` | On and off with `--only` and `--ignore`. Custom rules in the [Weir](https://writewithharper.com/docs/weir) language | Yes. It skipped front matter and code in the trial | One Rust binary, 58 MB (`harper-cli`) | `--format json`, with rule, span, line, and suggestions | Apache-2.0 | v2.12.0, 2026-10-01 | Good for grammar, weak for brevity |
| [proselint](https://github.com/amperser/proselint) | Usage: clichés, redundancy, hedging, weasel words, typography | On and off only (`.proselintrc.json`). No new rules | No. It lints front matter, code spans, and fenced code | Python 3.10 or newer | `-o json` | BSD-3-Clause | 0.16.0, 2025-11-14 | Poor. Vale ships the same checks as a package |
| [write-good](https://github.com/btford/write-good) | Passive voice, weasel words, "there is", adverbs, wordiness | Flags turn checks on or off; a JavaScript checks module can be passed in | No | Node | Text | MIT | 1.0.8, 2022-06-29 | Poor. Unmaintained, and also a Vale package |
| [textlint](https://github.com/textlint/textlint) | Nothing by itself. Rules come from npm: `sentence-length` (characters, not words), `write-good`, `max-comma`, `stop-words` | Yes, as JavaScript rules | Yes, with its markdown plugin | Node 20 or newer, about 55 MB of `node_modules` | `-f json` | MIT | v15.8.0, 2026-08-01 | Fair. Slow, and every rule is another npm dependency |
| [LanguageTool](https://github.com/languagetool-org/languagetool) | Grammar and style, with a picky level | `disabledRules` and `enabledRules` per request ([API](https://languagetool.org/http-api/)) | No. The caller strips markdown | Java HTTP server, about 250 MB | JSON | LGPL-2.1 | 6.8, 2026-05-05 | Poor. A heavy server for a style job |
| [alex](https://github.com/get-alex/alex) | Insensitive wording | Allow lists | Yes | Node | Text; JSON unverified | MIT | 11.0.0, 2022-11-04 | Not about brevity |
| [coala](https://github.com/coala/coala) | A framework of plug-ins ("bears") | Yes | Through bears | Python | JSON | AGPL-3.0 | 0.11.0, 2017-05-22 | Dead |
| [textstat](https://github.com/textstat/textstat) | Readability scores for a whole text: Flesch-Kincaid and others | None. A library | No | Python | Numbers | MIT | 0.7.13, 2026-02-18 | Scores training data in Python. Not a linter |
| [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2) | Markdown structure, no prose | Yes | Yes | Node | Text, JSON | MIT | 0.23.3; `scripts/lint-md.sh` pins 0.20.0 | In use already. Not a prose tool |
| [slop-cop](https://github.com/yasyf/slop-cop) | 226 rules for text that "sounds like an LLM", with optional tiers that call an LLM | Unverified | Yes | Go binary | JSON | MIT | v0.1.73, 2026-10-07 | Watch only. Young |

Release dates and licences come from the GitHub, npm, and PyPI pages and APIs, read on 2026-10-15 and 2026-10-16.

Other finds:

- [vale-ai-tells](https://github.com/krishnasunkam/vale-ai-tells) is a Vale package of 17 rules for patterns typical of AI writing. It is new, with one star: a source of rule ideas, not a dependency.
- [vale-cli/Harper](https://github.com/vale-cli/Harper) ports 455 of Harper's rules into a Vale style.
- Vale's code sits under `internal/`, so Go cannot import it. flai would run it as a process, as the `golangci-lint` tier does.
- [neurosnap/sentences](https://github.com/neurosnap/sentences) (MIT) splits text into sentences in Go, if flai counts words per sentence itself.

### Trial on this repository

The corpus was 14 files, about 6,388 words once front matter and fenced code were removed: the narratives `wip/archive/agents/S-0386.md` to `S-0389.md` and `S-0391.md`, the threads `TH-0447` to `TH-0451`, and the issues `I-0150` to `I-0153`. The tools ran from a scratch install in the git-ignored `bin/prose-trial/`, with every configuration outside the repository.

| Run | Alerts | Per 1,000 words | Wall time |
|-----|-------:|----------------:|----------:|
| Vale with the Microsoft, Google, write-good, proselint, and Readability packages and custom rules, untuned | 1,335 | 209.0 | 0.35 s |
| Vale with the recommended set below | 113 | 17.7 | 0.23 s |
| proselint 0.16.0, defaults | 16 | 2.5 | 0.35 s |
| Harper 2.12.0, defaults | 572 | 89.5 | 1.9 s |
| Harper without `SpellCheck`, `UseTitleCase`, and `SentenceCapitalization` | 29 | 4.5 | |
| textlint 15.8.0 with write-good, sentence-length, max-comma, stop-words | 186 | 29.1 | 6.0 s |
| textstat 0.7.13, one score per document | 14 documents | Flesch-Kincaid grade 4.5 to 11.6 | 0.86 s |

This prose has almost no filler. The custom rules found no filler word, hedge, praise, or apology in the 6,388 words. Its length comes from structure:

| Pattern | Hits | Example |
|---------|-----:|---------|
| Sentence over 25 words | 25 | "Kept the planner's seven tasks and their layers as drafted (story Notes § Planning): layer 1 T-1576, T-1577, T-1578 wait for none and share no path; layer 2 T-1579 waits for T-1577 and T-1578, …" (S-0388, 70 words) |
| Semicolon joining two ideas | 33 | "accept_reviews on since 2026-10-07; the orchestrator declined …" (I-0150) |
| More than three item IDs in one sentence | 13 | "S-0216 305,070, S-0215 303,610, S-0214 301,459, S-0213 298,706, …" (I-0153) |
| Two or more parentheticals in one sentence | 9 | "S-0341 (361,859 s), S-0342 (321,146 s), and S-0343 (250,881 s)" (I-0150) |
| Title over 25 words | 6 | "# S-0388 The review page shows each path a story's branch undoes on main as a card, with Restore, Keep with a reason, and Send back, …" |
| Document over a readability limit of grade 8 | 6 of 14 | I-0150, grade 10.2 |

A made-up paragraph in a chatty style ("Great question! Let me take a look … I think the issue might basically be …") tested whether each tool sees filler at all. Vale with the custom rules raised 21 alerts on it, proselint 1, and Harper 1, a false positive.

False positives seen:

- **Vale, untuned.** It linted front matter until `BlockIgnores` excluded it. Google.Headings flagged every timestamp heading. write-good.E-Prime flagged every "is" (167 hits), and Readability.Polysyllables flagged "operator" (169).
- **Vale, sentence length.** With the default token, a path counted as many words. The letters-only token below took the long-sentence count from 41 to 25. Four alerts pointed one line early.
- **proselint.** About 14 of its 16 hits were false: "ADR-0050's" read as a decade, and number lists as repeated words.
- **Harper.** 432 spelling hits fell on IDs and timestamps such as `TH`, `ADR`, and `flai`; a user dictionary would remove most. 62 title-case hits fell on sentence-case headings.
- **textlint.** write-good raised 119 hits, mostly "only can weaken meaning". `sentence-length` counts characters, so it cannot express a word limit.

### Recommendation

Adopt Vale, with a project style of nine rules written from [communication.md](../conventions/communication.md) and a few rules from the public packages. Keep Harper in reserve as a grammar checker. The others do not fit: proselint and write-good ship inside Vale as packages, textlint is slow and counts characters, LanguageTool is a Java server, alex is not about brevity, and coala is dead.

Reasons:

- Each mechanical rule of the communication convention fits a Vale check type: short sentences, one idea per sentence, no filler, no praise, no apology, no narration, and data sets out of prose.
- It is one Go binary with no runtime, like the tools flai runs already.
- Its JSON names the rule, line, and span of every alert. That gives the per-rule rate per 1,000 words that scores a rewriting model.
- It reads markdown, and it ran in 0.23 s on the corpus.
- It was the only tool that caught filler, hedges, praise, and narration in the chatty sample, and with front matter excluded it raised nothing on IDs, paths, or code.

Leave out the Google and Microsoft styles. They raised 209 alerts per 1,000 words and conflict with the first-person voice of threads.

The starting `.vale.ini`:

```ini
StylesPath = .vale/styles
MinAlertLevel = suggestion
Packages = write-good, proselint, Readability

[*.md]
BasedOnStyles = Vale, Flow, write-good, proselint, Readability
BlockIgnores = (?s)\A---\n.*?\n---\n
Vale.Spelling = NO
write-good.E-Prime = NO
write-good.Passive = NO
write-good.TooWordy = NO
Readability.Polysyllables = NO
```

The project style, `Flow`, one YAML file per rule, shown as one listing:

```yaml
# Flow/SentenceLength.yml. The token counts letter words only, so paths, IDs, and numbers do not count.
extends: occurrence
message: "Sentence has more than 25 words. Split it."
scope: sentence
level: warning
max: 25
token: (?<![\w/.\-])[A-Za-z][A-Za-z']*(?![\w/.\-])
---
# Flow/Semicolon.yml
extends: existence
message: "A semicolon joins two ideas. Make two sentences."
level: suggestion
nonword: true
tokens: [";"]
---
# Flow/Parens.yml
extends: occurrence
message: "More than one parenthetical in a sentence. Give each aside its own sentence."
scope: sentence
level: suggestion
max: 1
token: \([^)]*\)
---
# Flow/IDDensity.yml
extends: occurrence
message: "More than three item IDs in one sentence. Use a table or a list."
scope: sentence
level: suggestion
max: 3
token: \b(?:S|T|E|I|TH|MS|ADR)-\d{4}\b
---
# Flow/Filler.yml
extends: existence
message: "Cut the filler '%s'."
ignorecase: true
level: warning
tokens: [basically, actually, really, simply, just, very, quite, in order to,
  note that, it is worth noting, it should be noted, as mentioned,
  as noted above, of course, essentially, successfully, went ahead and]
---
# Flow/Hedge.yml
extends: existence
message: "Hedge '%s'. Say what is, or say it is unverified."
ignorecase: true
level: suggestion
tokens: [might, perhaps, possibly, probably, seems to, appears to,
  I think, I believe, arguably, somewhat]
---
# Flow/Praise.yml
extends: existence
message: "Praise nothing: '%s'."
ignorecase: true
level: error
tokens: [great question, excellent, perfect, awesome, fantastic,
  you're right, good catch]
---
# Flow/Narration.yml
extends: existence
message: "Report results, not your reasoning or tool use: '%s'."
ignorecase: true
level: warning
tokens: [let me, I'll now, I will now, I'm going to, I need to, "first, I"]
---
# Flow/Apology.yml
extends: existence
message: "Skip the apology '%s'. State the correction."
ignorecase: true
level: warning
tokens: [sorry, apologies, I apologize]
```

flai could call it in three places:

| Where | How | Status |
|-------|-----|--------|
| A test tier in `system-flow.yaml` | `{name: prose, command: [scripts/lint-prose.sh, "{files}"], paths: ["wip/**/*.md", "design/**/*.md"], format: plain}`, advisory at first, with only `Flow.Praise` at error level. A `vale-json` format would give structured findings later; flai has none today | Not built |
| A check on what flai writes | flai pipes a thread entry, item body, or commit message to `vale --output=JSON --ext=.md` and returns the alerts in the write's answer, or refuses the write, as its markdown lint does (alternative 5 below) | Not built |
| The scorer of a rewriting model | Alerts per 1,000 words per check, and the readability values. Vale prints a metric only when its condition fails, so continuous scores need copies of the Readability rules whose condition always holds | Not tested |

Two choices come first: how Vale is installed (a pinned download in a script, as `scripts/lint-md.sh` does with markdownlint-cli2, recorded in `design/tech/`), and whether `Flow` ships in `template/` for every project or stays in this repository.

Not checked in the trial: LanguageTool, alex, and slop-cop were not run; Vale on commit messages; and the always-on metric rules. The false-positive counts are a reading of each hit, not a measured precision.

### Sources

- Vale: [v3.24.0](https://github.com/vale-cli/vale/releases/tag/v3.24.0), [check types](https://docs.vale.sh/checks/), [Markdown](https://docs.vale.sh/formats/markdown), [CLI](https://docs.vale.sh/cli), [packages](https://github.com/vale-cli/packages)
- Harper: [v2.12.0](https://github.com/Automattic/harper/releases/tag/v2.12.0), [Weir](https://writewithharper.com/docs/weir), [Vale port](https://github.com/vale-cli/Harper)
- proselint: [v0.16.0](https://github.com/amperser/proselint/releases/tag/v0.16.0), [PyPI](https://pypi.org/project/proselint/)
- write-good: [repository](https://github.com/btford/write-good), [npm](https://www.npmjs.com/package/write-good)
- textlint: [v15.8.0](https://github.com/textlint/textlint/releases/tag/v15.8.0), [sentence-length](https://www.npmjs.com/package/textlint-rule-sentence-length)
- LanguageTool: [changes](https://github.com/languagetool-org/languagetool/blob/master/languagetool-standalone/CHANGES.md), [HTTP server](https://dev.languagetool.org/http-server), [HTTP API](https://languagetool.org/http-api/)
- [alex](https://github.com/get-alex/alex), [coala](https://github.com/coala/coala), [textstat](https://github.com/textstat/textstat), [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2), [slop-cop](https://github.com/yasyf/slop-cop), [vale-ai-tells](https://github.com/krishnasunkam/vale-ai-tells), [neurosnap/sentences](https://github.com/neurosnap/sentences)
