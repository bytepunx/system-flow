// The strategic agents' pages (S-0259, S-0228): a strategic agent's activity document as flai
// answers it for activity.document, the planner's, the orchestrator's, and the analyzer's runs as
// agent.status reports them, the /api/orchestrator and /api/analyzer answers, and how the pages read
// them: the log newest first, each run with what its log entries say it cost, the run under way,
// and each run's outcome.

import type { AgentRun } from './activity';

/** The strategic agents that keep an activity document, wip/agents/<kind>.md. */
export type StrategicKind = 'planner' | 'orchestrator' | 'analyzer';

/** One entry of a strategic agent's activity log: when it ended, what it did, and what that cost. */
export type ActivityEntry = {
	at: string;
	summary: string;
	trigger?: string;
	items: string[];
	seconds: number;
	cost: number;
	estimated: boolean;
};

/** A call flai guard refused a strategic agent, with the permission that would allow it (S-0218). */
export type ActivityRefusal = { at: string; call: string; needs?: string };

/** A strategic agent's activity document: its totals, its last run, its file, and its log oldest first. */
export type ActivityDocument<K extends StrategicKind = StrategicKind> = {
	kind: K;
	accrued_cost: number;
	accrued_seconds: number;
	tasks_completed: number;
	last_run: string;
	path: string;
	entries: ActivityEntry[];
	/** The calls flai guard refused it, oldest first; the orchestrator's alone has any. */
	refusals?: ActivityRefusal[];
};

/** An activity document as flai may answer it, its log or an entry's items null or absent. */
export type RawActivityDocument<K extends StrategicKind = StrategicKind> = Omit<
	ActivityDocument<K>,
	'entries'
> & {
	entries?: (Omit<ActivityEntry, 'items'> & { items?: string[] | null })[] | null;
};

/** The activity document with a missing log, or an entry's missing items, read as empty lists. */
export function activityDocument<K extends StrategicKind>(
	doc: RawActivityDocument<K>
): ActivityDocument<K> {
	return {
		...doc,
		entries: (doc.entries ?? []).map((e) => ({ ...e, items: e.items ?? [] }))
	};
}

/** What the run functions read of a run: when it started and ended, and how it went. */
export type RunSpan = Pick<AgentRun, 'started' | 'ended' | 'error' | 'outcome'>;

/**
 * The project's orchestrator run (S-0218), no story's: `held` while the operator has stopped it and
 * flai serve does not start it again for as long as the orchestrate action stays on (S-0228).
 */
export type OrchestratorRun = AgentRun & { held?: boolean };

/**
 * The project's analyzer run (S-0223), no story's: what it looks for (bottlenecks, intent, risk, or
 * all), what its starter said, and the report it wrote once it has ended, if any.
 */
export type AnalyzerRun = AgentRun & { focus: string; trigger?: string; report?: string };

/** What the operator may have flai do to the orchestrator from its page. */
export const ORCHESTRATOR_ACTIONS = ['stop', 'start'] as const;
export type OrchestratorAction = (typeof ORCHESTRATOR_ACTIONS)[number];

/** What the operator may ask the analyzer to look for; none asks for all of them. */
export const ANALYZER_FOCUSES = ['bottlenecks', 'intent', 'risk'] as const;
export type AnalyzerFocus = (typeof ANALYZER_FOCUSES)[number];

/** The strategic roles whose newest run agent.stream reads by role. */
export const STREAM_ROLES = ['orchestrate', 'analyze'] as const;
export type StreamRole = (typeof STREAM_ROLES)[number];

/**
 * The /api/orchestrator answer: whether the orchestrate host action is on, whether the operator
 * holds the orchestrator stopped, its activity (its decisions, each summary with its reason), the
 * run under way or null, and every run flai knows of, newest first.
 */
export type OrchestratorView = {
	enabled: boolean;
	held: boolean;
	activity: ActivityDocument<'orchestrator'>;
	run: OrchestratorRun | null;
	runs: OrchestratorRun[];
};

/**
 * The /api/analyzer answer: whether the analyze host action is on, the analyzer's activity, the run
 * under way or null, and every run flai knows of, newest first.
 */
export type AnalyzerView = {
	enabled: boolean;
	activity: ActivityDocument<'analyzer'>;
	run: AnalyzerRun | null;
	runs: AnalyzerRun[];
};

/** What flai's analyze.run answers: the analyzer it started now. */
export type AnalyzerStarted = {
	focus: string;
	agent: string;
	harness?: string;
	command: string;
	pid: number;
	log?: string;
	session?: string;
	started: string;
	trigger?: string;
};

/** What a run's log entries say it cost: their cost and seconds summed, whether any is estimated, and how many there are. */
export type RunCost = {
	cost: number;
	seconds: number;
	estimated: boolean;
	count: number;
};

/** A run with what its log entries say it cost. */
export type CostedRun<R extends RunSpan = RunSpan> = { run: R; cost: RunCost };

const time = (s: string) => new Date(s).getTime();

/** The log's entries newest first, leaving the list given as it was. */
export function newestFirst(entries: ActivityEntry[]): ActivityEntry[] {
	return [...entries].sort((a, b) => time(b.at) - time(a.at));
}

/** The runs newest first by when they started, leaving the list given as it was. */
export function runsNewestFirst<R extends RunSpan>(runs: R[]): R[] {
	return [...runs].sort((a, b) => time(b.started) - time(a.started));
}

/**
 * What the log entries within a run cost: those logged from its start, truncated to the second, to
 * its end rounded up to the next second, both inclusive, since an agent logs at an activity's end to
 * the second; to now while it runs. A run that could not start, with an error and no end, did no
 * work and counts none.
 */
export function runCost(run: RunSpan, entries: ActivityEntry[], now: Date = new Date()): RunCost {
	const none: RunCost = { cost: 0, seconds: 0, estimated: false, count: 0 };
	if (run.error && !run.ended) return none;
	const from = Math.floor(time(run.started) / 1000) * 1000;
	const to = run.ended ? Math.ceil(time(run.ended) / 1000) * 1000 : now.getTime();
	return entries
		.filter((e) => time(e.at) >= from && time(e.at) <= to)
		.reduce(
			(c, e) => ({
				cost: c.cost + e.cost,
				seconds: c.seconds + e.seconds,
				estimated: c.estimated || e.estimated,
				count: c.count + 1
			}),
			none
		);
}

/** Every run given, newest first by when it started, each with its cost from the log. */
export function pastRuns<R extends RunSpan>(
	runs: R[],
	entries: ActivityEntry[],
	now: Date = new Date()
): CostedRun<R>[] {
	return runsNewestFirst(runs).map((run) => ({ run, cost: runCost(run, entries, now) }));
}

/** The newest run that has neither ended nor failed to start, or null when none is under way. */
export function currentRun<R extends RunSpan>(runs: R[]): R | null {
	return runsNewestFirst(runs).find((r) => !r.ended && !r.error) ?? null;
}

/**
 * How a run went, in a word: running while it has neither ended nor failed to start, could not start
 * with an error and no outcome, and otherwise its outcome, or ended when flai recorded none.
 */
export function outcomeWord(run: RunSpan): string {
	if (!run.ended && !run.error) return 'running';
	if (run.error && !run.outcome) return 'could not start';
	return run.outcome ?? 'ended';
}
