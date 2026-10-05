// The planner's page (S-0259): the planner's activity document as flai answers it for
// activity.document, the /api/planner answer, and how the page reads them: the log newest first,
// each planner run with what its log entries say it cost, the run under way, and each run's outcome.

import type { PlanRun } from './activity';

/** One entry of the planner's activity log: when it ended a run, what it did, and what that cost. */
export type PlannerEntry = {
	at: string;
	summary: string;
	trigger?: string;
	items: string[];
	seconds: number;
	cost: number;
	estimated: boolean;
};

/** The planner's activity document: its totals, its last run, its file, and its log oldest first. */
export type PlannerActivity = {
	kind: 'planner';
	accrued_cost: number;
	accrued_seconds: number;
	tasks_completed: number;
	last_run: string;
	path: string;
	entries: PlannerEntry[];
};

/** The /api/planner answer: whether the plan host action is on, the planner's activity, and its runs newest first. */
export type PlannerView = {
	plan_enabled: boolean;
	activity: PlannerActivity;
	runs: PlanRun[];
};

/** What a run's log entries say it cost: their cost and seconds summed, whether any is estimated, and how many there are. */
export type RunCost = {
	cost: number;
	seconds: number;
	estimated: boolean;
	count: number;
};

/** A planner run with what its log entries say it cost. */
export type CostedRun = { run: PlanRun; cost: RunCost };

/** The planner's activity document as flai may answer it, its log or an entry's items null or absent. */
export type PlannerDocument = Omit<PlannerActivity, 'entries'> & {
	entries?: (Omit<PlannerEntry, 'items'> & { items?: string[] | null })[] | null;
};

/** The activity document with a missing log, or an entry's missing items, read as empty lists. */
export function plannerActivity(doc: PlannerDocument): PlannerActivity {
	return {
		...doc,
		entries: (doc.entries ?? []).map((e) => ({ ...e, items: e.items ?? [] }))
	};
}

const time = (s: string) => new Date(s).getTime();

/** The log's entries newest first, leaving the list given as it was. */
export function newestFirst(entries: PlannerEntry[]): PlannerEntry[] {
	return [...entries].sort((a, b) => time(b.at) - time(a.at));
}

/** The runs newest first by when they started, leaving the list given as it was. */
export function runsNewestFirst(runs: PlanRun[]): PlanRun[] {
	return [...runs].sort((a, b) => time(b.started) - time(a.started));
}

/**
 * What the log entries within a run cost: those logged from its start, truncated to the second, to
 * its end rounded up to the next second, both inclusive, since the planner logs at a run's end to the
 * second; to now while it runs. A run that could not start, with an error and no end, did no work and
 * counts none.
 */
export function runCost(run: PlanRun, entries: PlannerEntry[], now: Date = new Date()): RunCost {
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
export function pastRuns(
	runs: PlanRun[],
	entries: PlannerEntry[],
	now: Date = new Date()
): CostedRun[] {
	return runsNewestFirst(runs).map((run) => ({ run, cost: runCost(run, entries, now) }));
}

/** The newest run that has neither ended nor failed to start, or null when none is under way. */
export function currentRun(runs: PlanRun[]): PlanRun | null {
	return runsNewestFirst(runs).find((r) => !r.ended && !r.error) ?? null;
}

/**
 * How a run went, in a word: running while it has neither ended nor failed to start, could not start
 * with an error and no outcome, and otherwise its outcome, or ended when flai recorded none.
 */
export function outcomeWord(run: PlanRun): string {
	if (!run.ended && !run.error) return 'running';
	if (run.error && !run.outcome) return 'could not start';
	return run.outcome ?? 'ended';
}
