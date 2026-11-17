// Pure helpers for the review page (S-0041).

import { localTime } from './localtime';

export type Criterion = { text: string; checked: boolean };

/** The checkboxes under "## Acceptance criteria" in a story body. */
export function criteriaOf(body: string): Criterion[] {
	const section = sectionOf(body, 'Acceptance criteria');
	const out: Criterion[] = [];
	for (const line of section.split('\n')) {
		const m = /^\s*[-*]\s+\[( |x|X)\]\s+(.*\S)\s*$/.exec(line);
		if (m) out.push({ checked: m[1] !== ' ', text: m[2] });
	}
	return out;
}

/** The text under a level-two heading, up to the next one; empty when the heading is missing. */
export function sectionOf(markdown: string, heading: string): string {
	const lines = markdown.split('\n');
	const start = lines.findIndex((l) => l.trim().toLowerCase() === `## ${heading}`.toLowerCase());
	if (start < 0) return '';
	const rest = lines.slice(start + 1);
	const end = rest.findIndex((l) => /^##\s/.test(l));
	return (end < 0 ? rest : rest.slice(0, end)).join('\n').trim();
}

export type PatchKind = 'hunk' | 'add' | 'del' | 'context' | 'note';
/** One line of a patch: `sign` is its + or -, empty for the rest, and `text` is the line without it. */
export type PatchLine = { kind: PatchKind; sign: '+' | '-' | ''; text: string };

/**
 * A unified patch as lines with their kind, the sign of an added or removed line apart from its
 * text, so a view can put it in a margin. A hunk header and a "no newline" note keep their whole
 * text; the empty line the patch's last newline leaves is not a line.
 */
export function patchLines(patch: string): PatchLine[] {
	if (!patch) return [];
	const lines = patch.split('\n');
	if (lines[lines.length - 1] === '') lines.pop();
	return lines.map((raw): PatchLine => {
		if (raw.startsWith('@@')) return { kind: 'hunk', sign: '', text: raw };
		if (raw.startsWith('\\')) return { kind: 'note', sign: '', text: raw };
		if (raw.startsWith('+')) return { kind: 'add', sign: '+', text: raw.slice(1) };
		if (raw.startsWith('-')) return { kind: 'del', sign: '-', text: raw.slice(1) };
		// a context line begins with the space that stands where a sign would
		return { kind: 'context', sign: '', text: raw.startsWith(' ') ? raw.slice(1) : raw };
	});
}

/** Consecutive lines of one kind: a view outlines a run of added or of removed lines as one block. */
export type PatchRun = { kind: PatchKind; lines: PatchLine[] };

/** A unified patch as runs: every line beside its neighbours of the same kind, in the patch's order. */
export function patchRuns(patch: string): PatchRun[] {
	const runs: PatchRun[] = [];
	for (const line of patchLines(patch)) {
		const last = runs[runs.length - 1];
		if (last && last.kind === line.kind) last.lines.push(line);
		else runs.push({ kind: line.kind, lines: [line] });
	}
	return runs;
}

/**
 * Read a newline-delimited JSON response as it arrives, calling `each` for every complete line.
 * Lines that are not JSON are skipped; a last line without a newline is still delivered.
 */
export async function readNdjson(
	response: Response,
	each: (line: Record<string, unknown>) => void
): Promise<void> {
	const emit = (raw: string) => {
		const text = raw.trim();
		if (!text) return;
		try {
			each(JSON.parse(text));
		} catch {
			// not a JSON line
		}
	};
	if (!response.body) {
		for (const l of (await response.text()).split('\n')) emit(l);
		return;
	}
	const reader = response.body.getReader();
	const decoder = new TextDecoder();
	let pending = '';
	for (;;) {
		const { value, done } = await reader.read();
		if (done) break;
		pending += decoder.decode(value, { stream: true });
		const lines = pending.split('\n');
		pending = lines.pop() ?? '';
		lines.forEach(emit);
	}
	emit(pending + decoder.decode());
}

/** Where a step of flai verify got to (S-0270); a reused tier passed in an earlier run of the same head (S-0341), a tier another selected tier covers is skipped (S-0342), and a tier left to the merge queue's batch is queued (ADR-0174). */
export type VerifyState = 'passed' | 'failed' | 'not-reached' | 'reused' | 'skipped' | 'queued';
/** What a failed step found: a test, a check, or a linter's complaint, where it is, and what it says. */
export type VerifyFinding = { name?: string; path?: string; line?: number; message?: string };
/** One step of flai verify: one of its own, or a test or lint tier when `tier` is set. */
export type VerifyStep = {
	name: string;
	tier?: boolean;
	state: VerifyState;
	/** The tier that covers a skipped one. */
	covered_by?: string;
	duration_ms: number;
	duration?: string;
	findings?: VerifyFinding[];
	omitted?: number;
	command?: string[];
	dir?: string;
	exit_code?: number;
	/** When the run that ran a reused tier began, in UTC. */
	reused_from?: string;
};
/** A finding of flai check outside the story: it does not fail the story (S-0249). */
export type VerifyNote = {
	rule: string;
	level: string;
	path?: string;
	line?: number;
	message: string;
};
/** How a record's commit stands to the story branch's head (ADR-0292 §3), as flai judges it. */
export type KlusterHeadJudgement = {
	commit: string;
	head: string;
	/** The commit counts as the head: it is the head, or only design, docs, and wip changed since. */
	at_head: boolean;
	/** The commit is the head itself. */
	exact?: boolean;
	/** Paths outside design, docs, and wip that changed since, when that is why it is not at the head. */
	outside?: string[];
	/** Why it is not at the head; absent when it is. */
	why?: string;
};
/** One step of a namespace's last job, with its log in the main checkout. */
export type KlusterEvidenceStep = {
	name: string;
	started_at?: string;
	ended_at?: string;
	outcome: string;
	exit_code?: number;
	log?: string;
};
/** A story's namespace's last deploy, or teardown, as its record keeps it. */
export type KlusterDeployEvidence = {
	commit: string;
	/** running, passed, or failed. */
	outcome: string;
	failed_step?: string;
	exit_code?: number;
	error?: string;
	/** deploy or teardown. */
	job: string;
	/** The outcome in words: "passed", "failed at the build step", "is still running". */
	said: string;
	/** Who asked for it: a story's agent or flai serve; absent for one run in a shell. */
	asked_by?: string;
	started_at?: string;
	ended_at?: string;
	steps?: KlusterEvidenceStep[];
	head: KlusterHeadJudgement;
	passed_at_head: boolean;
};
/** The last checks run against a story's namespace since its deploy. */
export type KlusterChecksEvidence = {
	commit: string;
	/** passed, failed, timed-out, or cancelled. */
	outcome: string;
	log?: string;
	head: KlusterHeadJudgement;
	passed_at_head: boolean;
};
/**
 * What a story's kluster namespace shows of its change (ADR-0292 §6): its last deploy and checks
 * run, and whether each passed at the head of its branch. verify.status reads it on the host and
 * leaves it out while none applies: the action off, or a block with no deploy.
 */
export type KlusterEvidence = {
	applies: boolean;
	why?: string;
	story: string;
	namespace?: string;
	/** The story branch's head; absent when the branch cannot be resolved. */
	head?: string;
	/** The checks host action is on, so the checks must pass at the head too. */
	checks_on: boolean;
	deploy?: KlusterDeployEvidence;
	checks?: KlusterChecksEvidence;
	passed: boolean;
	commands?: { deploy: string; deploy_tool: string; checks: string; checks_tool: string };
};
/** flai verify's report for a story, as it stores the last one and verify.status answers it. */
export type VerifyReport = {
	story: string;
	commit: string;
	base: string;
	ran_at: string;
	duration_ms: number;
	duration?: string;
	passed: boolean;
	stopped_at?: string;
	paths?: string[];
	steps: VerifyStep[];
	notes?: VerifyNote[];
	/** The story's kluster namespace, read on the host as the page asks (S-0662). */
	kluster?: KlusterEvidence;
	/** The report's commit is not the head of the story's branch: the head moved since it ran (I-0384). */
	moved?: boolean;
	/** The close-out job running at another commit than the report's (I-0384). */
	running?: CloseoutProgress;
};
/** A close-out job as flai answers it (closeout.Progress): its ID, when it started, the head it runs at, and its step. */
export type CloseoutProgress = { id: string; started: string; head?: string; step?: string };
/** What verify.status answers while a close-out job runs for a story that has no report yet (I-0384). */
export type VerifyRunning = { story: string; running: CloseoutProgress };
/** What verify.status answers for a story with a report or a running close-out job: the report, or the job alone. */
export type VerifyStatus = VerifyReport | VerifyRunning;

/** A finding or a note as a page shows it: where, as path:line and a name, and what it says. */
export type VerifyLine = { where: string; message: string };
/** A step as a row: its state in words, its duration or a reused tier's origin, and the findings when it failed. */
export type VerifyRow = {
	name: string;
	tier: boolean;
	state: VerifyState;
	label: string;
	duration: string;
	findings: VerifyLine[];
	/** How many findings flai left out past its cap. */
	omitted: number;
};
/** A report as a page shows it: the outcome, the commit and when, a row per step, the notes apart. */
export type VerifyView = {
	/** Whether there is a report; false while a close-out runs before the story's first, and then only running is said. */
	reported: boolean;
	passed: boolean;
	outcome: string;
	commit: string;
	base: string;
	ranAt: string;
	duration: string;
	rows: VerifyRow[];
	notes: VerifyLine[];
	/** That the report is from before the branch's head moved; empty when it is at the head. */
	moved: string;
	/** The close-out job running at another commit, its head, step, and start; empty when none runs. */
	running: string;
};

const STATE_LABEL: Record<VerifyState, string> = {
	passed: 'passed',
	failed: 'failed',
	'not-reached': 'not reached',
	reused: 'reused',
	skipped: 'skipped',
	queued: 'left to the merge queue'
};

/** A step's state in words: a skipped step names the tier that covers it, when the report does. */
function label(s: VerifyStep): string {
	if (s.state === 'skipped' && s.covered_by) return `skipped: covered by ${s.covered_by}`;
	return STATE_LABEL[s.state] ?? s.state;
}

/** A step not reached, skipped, or left to the merge queue ran nothing, so it has no duration. */
const RAN_NOTHING: VerifyState[] = ['not-reached', 'skipped', 'queued'];

/** A duration in milliseconds as flai writes one when the report carries none: 340ms, 5.1s. */
function took(ms: number): string {
	return ms < 1000 ? `${Math.round(ms)}ms` : `${(ms / 1000).toFixed(1)}s`;
}

/** Where a finding or note is, as flai verify prints it: path:line, then its name. */
function where(path?: string, line?: number, name?: string): string {
	const loc = path ? (line ? `${path}:${line}` : path) : '';
	return [loc, name ?? ''].filter((s) => s).join(' ');
}

/** A step's duration as its row shows it: none when not reached, skipped, or queued, and for a reused tier the run it comes from. */
function rowDuration(s: VerifyStep): string {
	if (RAN_NOTHING.includes(s.state)) return '';
	if (s.state === 'reused' && s.reused_from) return `from ${localTime(s.reused_from)}`;
	return s.duration ?? took(s.duration_ms);
}

/**
 * A report as rows for the review page (S-0270): a row per step in the order flai ran them, with
 * the findings under the step that failed only; a step not reached or skipped has no duration, and
 * a reused tier shows, in its place, when the run it comes from began (S-0341). The notes from
 * outside the story are apart, since they do not fail it. A report from before the head moved
 * says so, and a close-out running at another commit is named with its job, step, and start; while
 * one runs before the story's first report, the view says only that (I-0384).
 */
export function verifyView(status: VerifyStatus): VerifyView {
	const running = status.running ? runningNote(status.running) : '';
	if (!('steps' in status)) {
		return {
			reported: false,
			passed: false,
			outcome: 'No flai verify result yet',
			commit: '',
			base: '',
			ranAt: '',
			duration: '',
			rows: [],
			notes: [],
			moved: '',
			running
		};
	}
	const report = status;
	const head = report.running?.head;
	const rows = report.steps.map((s): VerifyRow => ({
		name: s.name,
		tier: s.tier === true,
		state: s.state,
		label: label(s),
		duration: rowDuration(s),
		findings:
			s.state === 'failed'
				? (s.findings ?? []).map((f) => ({
						where: where(f.path, f.line, f.name),
						message: f.message ?? ''
					}))
				: [],
		omitted: s.state === 'failed' ? (s.omitted ?? 0) : 0
	}));
	return {
		reported: true,
		passed: report.passed,
		outcome: report.passed
			? 'Passed every step'
			: report.stopped_at
				? `Stopped at ${report.stopped_at}`
				: 'Failed',
		commit: report.commit.slice(0, 12),
		base: report.base,
		ranAt: report.ran_at,
		duration: report.duration ?? took(report.duration_ms),
		rows,
		notes: (report.notes ?? []).map((n) => ({
			where: where(n.path, n.line, n.rule),
			message: `${n.level}: ${n.message}`
		})),
		moved: report.moved
			? `This report is from before the head moved${head ? ` to ${head.slice(0, 12)}` : ''}.`
			: '',
		running
	};
}

/**
 * A running close-out job in words (I-0384): the head it runs at, short, its ID, its step, and when
 * it started; each part left out when flai gives none.
 */
export function runningNote(job: CloseoutProgress): string {
	const at = job.head ? ` at ${job.head.slice(0, 12)}` : '';
	const step = job.step ? `, at ${job.step}` : '';
	const since = job.started ? `, since ${localTime(job.started)}` : '';
	return `Close-out running${at}: job ${job.id}${step}${since}.`;
}

/**
 * What accept.resolve answers once flai serve took the request (S-0804): the story, the job's ID,
 * whether this request started it or one of the story ran already, where its output is, and what
 * it does, in a sentence.
 */
export type ResolveStarted = {
	id: string;
	job: string;
	started: boolean;
	running: boolean;
	log: string;
	reason: string;
};
/** A step or tier of a job as its log's progress lines show it (closeout.StepProgress). */
export type JobStep = {
	name: string;
	kind: 'step' | 'tier';
	/** Who ran it: flai verify, or the resolve job itself. */
	source: 'verify' | 'resolve';
	/** Absent for a tier that ended without running, as one reused or skipped. */
	started?: string;
	/** Absent, with outcome, while it runs. */
	ended?: string;
	outcome?: string;
};
/**
 * A story's last close-out or Rebase and resolve job as resolve.status answers it (closeout.JobStatus
 * and whether it runs): when it started and ended, its verdict and the step it stopped at, how a
 * resolve job ended and why, the step and tier it is on, its steps, and where its log is.
 */
export type ResolveJob = {
	id: string;
	story: string;
	/** resolve for a Rebase and resolve; absent for a close-out. */
	kind?: string;
	by: string;
	started: string;
	head?: string;
	ended?: string;
	exit_code?: number;
	/** passed, stopped, failed, or interrupted, when flai serve's restart cut the job short. */
	verdict?: string;
	stopped_at?: string;
	commit?: string;
	duration_ms?: number;
	/** stayed in review or sent back, once a resolve job ended; absent for one interrupted. */
	outcome?: string;
	thread?: string;
	conflicts?: string[];
	reason?: string;
	step?: string;
	tier?: string;
	steps?: JobStep[];
	log: string;
	running: boolean;
};
/** What resolve.status answers: the story, and its last job, absent when it has had none. */
export type ResolveStatus = { id: string; job?: ResolveJob };

/** How a job ended, as the review page says it. */
export type JobEndKind = 'stayed' | 'sent-back' | 'interrupted' | 'closed-out';
/** A job's end as the review page shows it: what came of it, why, and what it names. */
export type JobEnd = {
	kind: JobEndKind;
	headline: string;
	reason: string;
	/** The send-back thread, when the job opened one. */
	thread: string;
	/** The paths the rebase stopped on, which the job aborted. */
	conflicts: string[];
	/** The step a close-out stopped at; empty when none. */
	stoppedAt: string;
};
/** A job as the review page shows it while it runs and once it ended. */
export type JobView = {
	/** Rebase and resolve, or Close-out. */
	label: string;
	running: boolean;
	/** The steps it is on, outermost first, such as close-out › flai verify › check; empty when its log cannot say. */
	step: string;
	/** The tiers running, comma-separated; empty when none does. */
	tiers: string;
	/** Null while it runs. */
	end: JobEnd | null;
};

/** A span of milliseconds as the review page counts it up: 7s, 1m 05s, 1h 02m. */
export function elapsed(ms: number): string {
	const s = Math.max(0, Math.floor(ms / 1000));
	const m = Math.floor(s / 60);
	const h = Math.floor(m / 60);
	const two = (n: number) => String(n).padStart(2, '0');
	if (h) return `${h}h ${two(m % 60)}m`;
	return m ? `${m}m ${two(s % 60)}s` : `${s}s`;
}

/**
 * How long a job has run at now, in milliseconds since the epoch, while it runs; once it ended,
 * how long it ran, from its duration or its start and end. Empty when its start cannot be read.
 */
export function jobElapsed(job: ResolveJob, now: number): string {
	if (!job.running && job.duration_ms !== undefined) return elapsed(job.duration_ms);
	const from = Date.parse(job.started);
	const to = !job.running && job.ended ? Date.parse(job.ended) : now;
	return Number.isFinite(from) && Number.isFinite(to) ? elapsed(to - from) : '';
}

/** The send back's words for a close-out that did not pass: the step it stopped at, or failed. */
function stopped(job: ResolveJob): string {
	return job.stopped_at ? `stopped at ${job.stopped_at}` : 'failed';
}

/** How an ended job came out, from its outcome, else its verdict. */
function endOf(job: ResolveJob): JobEnd {
	const base = {
		reason: job.reason ?? '',
		thread: job.thread ?? '',
		conflicts: job.conflicts ?? [],
		stoppedAt: job.stopped_at ?? ''
	};
	if (job.verdict === 'interrupted')
		return {
			...base,
			kind: 'interrupted',
			headline: `Interrupted: flai serve restarted while job ${job.id} ran; nothing was sent back`
		};
	if (job.outcome === 'stayed in review')
		return {
			...base,
			kind: 'stayed',
			headline: `Stayed in review: ${job.story} is ready to accept`
		};
	if (job.outcome === 'sent back')
		return {
			...base,
			kind: 'sent-back',
			headline: base.conflicts.length
				? 'Sent back to in-progress: the rebase stopped on conflicts'
				: job.stopped_at === 'rebase'
					? 'Sent back to in-progress: the rebase did not end clean'
					: `Sent back to in-progress: the close-out ${stopped(job)}`
		};
	return {
		...base,
		kind: 'closed-out',
		headline:
			job.verdict === 'passed'
				? 'The close-out passed; no Rebase and resolve ran'
				: `The close-out ${stopped(job)}; no Rebase and resolve ran`
	};
}

/**
 * A close-out or Rebase and resolve job as the review page shows it (S-0804). While it runs, the
 * steps it is on: the resolve job's own (the rebase, or the close-out), then the close-out's step
 * as flai serve names it, then flai verify's, and the tiers running. Once it ended, how: a resolve
 * job that stayed in review, one sent back with its thread and the paths in conflict or the step,
 * one interrupted, or a close-out's verdict.
 */
export function jobView(job: ResolveJob): JobView {
	const label = job.kind === 'resolve' ? 'Rebase and resolve' : 'Close-out';
	if (!job.running) return { label, running: false, step: '', tiers: '', end: endOf(job) };
	const on = (job.steps ?? []).filter((s) => s.started && !s.ended);
	const names = on.filter((s) => s.kind === 'step' && s.source === 'resolve').map((s) => s.name);
	if (job.step && !names.includes(job.step)) names.push(job.step);
	names.push(...on.filter((s) => s.kind === 'step' && s.source === 'verify').map((s) => s.name));
	const tiers = on.filter((s) => s.kind === 'tier').map((s) => s.name);
	if (job.tier && !tiers.includes(job.tier)) tiers.push(job.tier);
	return { label, running: true, step: names.join(' › '), tiers: tiers.join(', '), end: null };
}

/** A deploy or a checks run of the story's namespace as the review page shows it (S-0662). */
export type KlusterRunView = {
	/** A deploy's or teardown's, or a checks run's: their logs are shown in different places. */
	kind: 'deploy' | 'checks';
	/** Deploy, Teardown, or Checks. */
	label: string;
	/** Its commit, short. */
	commit: string;
	/** Its outcome in words. */
	outcome: string;
	/** Its outcome's colour: a run that failed or timed out failed. */
	tone: 'passed' | 'failed' | 'other';
	/** Who asked for it, in words. */
	asked: string;
	/** When it ended, in UTC; empty when it has not or the record does not say. */
	endedAt: string;
	atHead: boolean;
	/** How it stands to the head, in words: at the head, or why not. */
	head: string;
	/** Where its log is shown; empty when it has none. */
	logHref: string;
	/** What the log link opens. */
	logLabel: string;
};
/** A story's kluster evidence as the review page shows it beside the criteria (S-0662). */
export type KlusterView = {
	namespace: string;
	/** The story branch's head, short; empty when the branch cannot be resolved. */
	head: string;
	passed: boolean;
	/** Whether the namespace passed at the head, in a sentence. */
	headline: string;
	deploy: KlusterRunView | null;
	checks: KlusterRunView | null;
	/** What the page says in place of a deploy the record lacks. */
	noDeploy: string;
	/** What the page says in place of a checks run the record lacks. */
	noChecks: string;
};

/** Where the review page shows a checks run's log: its Checks section. */
export const CHECKS_LOG_HREF = '#review-checks';

const short = (commit: string | undefined) => (commit ?? '').slice(0, 12);

/** How a record's commit stands to the head, in words. */
function headWords(j: KlusterHeadJudgement): string {
	if (!j.at_head) return `not at the head: ${j.why ?? 'its commit is not the head'}`;
	if (j.exact) return 'at the head';
	return `at the head: from ${short(j.commit)}, followed only by changes under design, docs, and wip`;
}

const asked = (by: string | undefined) => (by ? `asked by ${by}` : 'who asked is not recorded');

/** A deploy's view: its log is the Deployments page's Logs pane, deploymentsHref#logs. */
function deployRun(d: KlusterDeployEvidence, deploymentsHref: string): KlusterRunView {
	return {
		kind: 'deploy',
		label: d.job === 'teardown' ? 'Teardown' : 'Deploy',
		commit: short(d.commit),
		outcome: d.said || d.outcome,
		tone: d.outcome === 'passed' ? 'passed' : d.outcome === 'failed' ? 'failed' : 'other',
		asked: asked(d.asked_by),
		endedAt: d.ended_at ?? '',
		atHead: d.head.at_head,
		head: headWords(d.head),
		logHref: d.steps?.length ? `${deploymentsHref}#logs` : '',
		logLabel: d.failed_step ? `the ${d.failed_step} step's log` : "its steps' logs"
	};
}

/** A checks run's view: its log is where the review page shows checks logs. */
function checksRun(c: KlusterChecksEvidence): KlusterRunView {
	return {
		kind: 'checks',
		label: 'Checks',
		commit: short(c.commit),
		outcome: c.outcome === 'timed-out' ? 'timed out' : c.outcome,
		tone:
			c.outcome === 'passed'
				? 'passed'
				: c.outcome === 'failed' || c.outcome === 'timed-out'
					? 'failed'
					: 'other',
		asked: asked(undefined),
		endedAt: '',
		atHead: c.head.at_head,
		head: headWords(c.head),
		logHref: c.log ? CHECKS_LOG_HREF : '',
		logLabel: 'its log'
	};
}

/**
 * A story's kluster evidence as the review page shows it beside the criteria (S-0662): the
 * namespace's last deploy and checks run, each with its commit, outcome, who asked, when it ended,
 * whether it is at the head and why not, and where its log is. A deploy's log opens the Deployments
 * page, at deploymentsHref, on its Logs pane; a checks run's opens the page's Checks section. Null
 * when no evidence applies, and the page shows nothing.
 */
export function klusterView(
	e: KlusterEvidence | undefined | null,
	deploymentsHref: string
): KlusterView | null {
	if (!e?.applies) return null;
	const head = short(e.head);
	const what = e.checks_on ? 'deployed and checked' : 'deployed';
	let headline: string;
	if (!head) headline = `${e.story} has no branch story/${e.story}, so nothing is at its head`;
	else if (e.passed) headline = `Its namespace is ${what} at the head ${head}`;
	else headline = `Its namespace is not ${what} at the head ${head}`;
	return {
		namespace: e.namespace ?? '',
		head,
		passed: e.passed,
		headline,
		deploy: e.deploy ? deployRun(e.deploy, deploymentsHref) : null,
		checks: e.checks ? checksRun(e.checks) : null,
		noDeploy: 'No deploy of the namespace is recorded.',
		noChecks: e.checks_on
			? 'No checks run against it has ended since its deploy.'
			: 'The checks host action is off, so no checks run is needed.'
	};
}

/**
 * The body with its nth acceptance criterion ticked or unticked (S-0085). Only lines of the
 * "Acceptance criteria" section count, in order, as criteriaOf lists them; null when there is no
 * such criterion. Everything else in the body is left exactly as it was.
 */
export function toggleCriterion(body: string, index: number): string | null {
	const lines = body.split('\n');
	const start = lines.findIndex((l) => l.trim().toLowerCase() === '## acceptance criteria');
	if (start < 0 || index < 0) return null;
	let n = -1;
	for (let i = start + 1; i < lines.length && !/^##\s/.test(lines[i]); i++) {
		const m = /^(\s*[-*]\s+\[)( |x|X)(\]\s+.*\S\s*)$/.exec(lines[i]);
		if (!m) continue;
		if (++n === index) {
			lines[i] = m[1] + (m[2] === ' ' ? 'x' : ' ') + m[3];
			return lines.join('\n');
		}
	}
	return null;
}
