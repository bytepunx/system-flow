// The analyzer page's read and its Run (S-0228): whether the analyze host action is on, the
// analyzer's activity document, and its runs, each asked of flai on the host.

import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { StrategicStatus } from '$lib/server/agent';
import {
	activityDocument,
	ANALYZER_FOCUSES,
	currentRun,
	type AnalyzerFocus,
	type AnalyzerStarted,
	type AnalyzerView,
	type RawActivityDocument
} from '$lib/strategic';
import type { RequestHandler } from './$types';

/** An answer of flai's, or null when flai cannot be asked: the page then reads it as off and empty. */
const orNull = <T>(asked: Promise<T>): Promise<T | null> => asked.catch(() => null);

/**
 * GET: whether the analyze host action is on (project.info's host_actions.analyze), the analyzer's
 * activity document, and its runs from agent.status, newest first. A flai that cannot be asked reads as
 * the action off and no runs; the activity document is the page's subject, so failing to read it
 * is the answer's error.
 */
export const GET: RequestHandler = () =>
	respond(async (): Promise<AnalyzerView> => {
		const [info, status, activity] = await Promise.all([
			orNull(repo().ask<{ host_actions?: Record<string, boolean> }>('project.info')),
			orNull(repo().ask<StrategicStatus>('agent.status')),
			repo()
				.ask<RawActivityDocument<'analyzer'>>('activity.document', { kind: 'analyzer' })
				.then(activityDocument)
		]);
		const newest = status?.state?.analyzer ?? null;
		const runs = [...(newest ? [newest] : []), ...(status?.state?.past_analyzers ?? [])];
		return {
			enabled: info?.host_actions?.analyze === true,
			activity,
			run: currentRun(runs),
			runs
		};
	});

/**
 * POST { focus }: flai's analyze.run, the analyzer started now looking for the focus, or for all of
 * them when none is given, which answers the run it started. flai judges whether it may and says
 * why not (one already running), and it is 403 with what enables it while the operator has not
 * (flai serve enable analyze); a focus that is not one is 400.
 */
export const POST: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => null)) as { focus?: unknown } | null;
		const focus = body?.focus ?? '';
		if (focus !== '' && !ANALYZER_FOCUSES.includes(focus as AnalyzerFocus))
			throw new RepoError(
				400,
				`focus must be one of ${ANALYZER_FOCUSES.join(', ')}, or none for all of them, not ${JSON.stringify(focus)}`
			);
		const { data } = await repo().write<AnalyzerStarted>('analyze.run', focus ? { focus } : {});
		return data;
	});
