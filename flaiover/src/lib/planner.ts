// The planner's page (S-0259): the planner's activity document as flai answers it for
// activity.document, the /api/planner answer, and how the page reads them, as strategic.ts reads
// any strategic agent's (S-0228): the log newest first, each planner run with what its log entries
// say it cost, the run under way, and each run's outcome.

import type { PlanRun } from './activity';
import {
	activityDocument,
	type ActivityDocument,
	type ActivityEntry,
	type CostedRun as Costed,
	type RawActivityDocument
} from './strategic';

export {
	currentRun,
	newestFirst,
	outcomeWord,
	pastRuns,
	runCost,
	runsNewestFirst,
	type RunCost
} from './strategic';

/** One entry of the planner's activity log: when it ended a run, what it did, and what that cost. */
export type PlannerEntry = ActivityEntry;

/** The planner's activity document: its totals, its last run, its file, and its log oldest first. */
export type PlannerActivity = ActivityDocument<'planner'>;

/** The /api/planner answer: whether the plan host action is on, the planner's activity, and its runs newest first. */
export type PlannerView = {
	plan_enabled: boolean;
	activity: PlannerActivity;
	runs: PlanRun[];
};

/** A planner run with what its log entries say it cost. */
export type CostedRun = Costed<PlanRun>;

/** The planner's activity document as flai may answer it, its log or an entry's items null or absent. */
export type PlannerDocument = RawActivityDocument<'planner'>;

/** The activity document with a missing log, or an entry's missing items, read as empty lists. */
export function plannerActivity(doc: PlannerDocument): PlannerActivity {
	return activityDocument(doc);
}
