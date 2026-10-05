// Behaviour tests for the planner page's reading of the planner's log and runs (S-0259).

import { describe, expect, it } from 'vitest';
import type { PlanRun } from './activity';
import {
	currentRun,
	newestFirst,
	outcomeWord,
	pastRuns,
	plannerActivity,
	runCost,
	type PlannerEntry
} from './planner';

const entry = (at: string, cost: number, seconds: number, estimated = false): PlannerEntry => ({
	at,
	summary: `planned at ${at}`,
	items: [],
	seconds,
	cost,
	estimated
});

const run = (item: string, started: string, more: Partial<PlanRun> = {}): PlanRun => ({
	item,
	agent: 'planner',
	command: 'claude',
	started,
	...more
});

describe('planner (S-0259)', () => {
	it('reads a missing log, or an entry with no items, as empty lists', () => {
		const doc = {
			kind: 'planner' as const,
			accrued_cost: 0,
			accrued_seconds: 0,
			tasks_completed: 0,
			last_run: '',
			path: 'wip/agents/planner.md'
		};
		expect(plannerActivity({ ...doc, entries: null }).entries).toEqual([]);
		expect(plannerActivity(doc).entries).toEqual([]);
		const [e] = plannerActivity({
			...doc,
			entries: [{ ...entry('2026-10-03T10:00:00Z', 1, 60), items: null }]
		}).entries;
		expect(e.items).toEqual([]);
	});

	it('lists the log newest first and leaves the log given as it was', () => {
		const log = [
			entry('2026-10-01T09:00:00Z', 1, 10),
			entry('2026-10-03T09:00:00Z', 3, 30),
			entry('2026-10-02T09:00:00Z', 2, 20)
		];
		expect(newestFirst(log).map((e) => e.cost)).toEqual([3, 2, 1]);
		expect(log.map((e) => e.cost)).toEqual([1, 3, 2]);
	});

	it('costs a run by the entries logged from its start to its end, its last second included', () => {
		const r = run('S-0252', '2026-10-03T10:00:00.700Z', { ended: '2026-10-03T10:30:12.400Z' });
		const log = [
			entry('2026-10-03T09:59:59Z', 100, 100),
			entry('2026-10-03T10:00:00Z', 1, 10),
			entry('2026-10-03T10:30:13Z', 2.5, 20, true),
			entry('2026-10-03T10:30:14Z', 100, 100)
		];
		expect(runCost(r, log)).toEqual({ cost: 3.5, seconds: 30, estimated: true, count: 2 });
		expect(runCost(r, [])).toEqual({ cost: 0, seconds: 0, estimated: false, count: 0 });
	});

	it('costs a run under way up to now, and a run that could not start as nothing', () => {
		const now = new Date('2026-10-03T11:00:00Z');
		const log = [entry('2026-10-03T10:10:00Z', 1, 10), entry('2026-10-03T11:00:01Z', 9, 90)];
		expect(runCost(run('E-0016', '2026-10-03T10:00:00Z'), log, now)).toMatchObject({
			cost: 1,
			count: 1
		});
		const failed = run('E-0016', '2026-10-03T10:00:00Z', { error: 'claude: not found' });
		expect(runCost(failed, log, now)).toEqual({
			cost: 0,
			seconds: 0,
			estimated: false,
			count: 0
		});
	});

	it('lists every run newest first, each with its cost, and leaves the runs given as they were', () => {
		const runs = [
			run('E-0016', '2026-10-01T10:00:00Z', { ended: '2026-10-01T10:05:00Z', outcome: 'worked' }),
			run('S-0252', '2026-10-03T10:00:00Z', { ended: '2026-10-03T10:05:00Z', outcome: 'asked' })
		];
		const log = [entry('2026-10-01T10:05:00Z', 1, 300), entry('2026-10-03T10:05:00Z', 2, 300)];
		const past = pastRuns(runs, log);
		expect(past.map((p) => [p.run.item, p.cost.cost])).toEqual([
			['S-0252', 2],
			['E-0016', 1]
		]);
		expect(runs.map((r) => r.item)).toEqual(['E-0016', 'S-0252']);
		expect(pastRuns([], log)).toEqual([]);
	});

	it('finds the newest run under way, passing over ended runs and runs that could not start', () => {
		const ended = run('E-0016', '2026-10-01T10:00:00Z', { ended: '2026-10-01T10:05:00Z' });
		const older = run('S-0250', '2026-10-02T10:00:00Z');
		const newer = run('S-0252', '2026-10-03T10:00:00Z');
		const failed = run('S-0259', '2026-10-04T10:00:00Z', { error: 'claude: not found' });
		expect(currentRun([ended, older, failed, newer])).toBe(newer);
		expect(currentRun([ended, failed])).toBeNull();
		expect(currentRun([])).toBeNull();
	});

	it('says how a run went in a word', () => {
		const at = '2026-10-03T10:00:00Z';
		const done = { ended: '2026-10-03T10:05:00Z' };
		expect(outcomeWord(run('S-0252', at))).toBe('running');
		expect(outcomeWord(run('S-0252', at, { error: 'claude: not found' }))).toBe('could not start');
		for (const outcome of ['worked', 'failed', 'asked', 'stopped'] as const)
			expect(outcomeWord(run('S-0252', at, { ...done, outcome }))).toBe(outcome);
		expect(outcomeWord(run('S-0252', at, done))).toBe('ended');
	});
});
