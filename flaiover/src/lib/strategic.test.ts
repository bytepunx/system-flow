// Behaviour tests for how the strategic agents' pages read an activity log and runs (S-0228), on
// the orchestrator's and the analyzer's runs; planner.test.ts covers the planner's.

import { describe, expect, it } from 'vitest';
import {
	activityDocument,
	currentRun,
	newestFirst,
	outcomeWord,
	pastRuns,
	runCost,
	runsNewestFirst,
	type ActivityEntry,
	type AnalyzerRun,
	type OrchestratorRun
} from './strategic';

const entry = (at: string, cost: number, seconds: number, estimated = false): ActivityEntry => ({
	at,
	summary: `promoted S-0252 at ${at}: cost of delay 60.81 USD a week, first under policy cod`,
	items: ['S-0252'],
	seconds,
	cost,
	estimated
});

const orchestrator = (started: string, more: Partial<OrchestratorRun> = {}): OrchestratorRun => ({
	story: '',
	agent: 'orchestrator',
	command: 'claude',
	started,
	...more
});

const analyzer = (started: string, more: Partial<AnalyzerRun> = {}): AnalyzerRun => ({
	story: '',
	agent: 'analyzer',
	command: 'claude',
	focus: 'risk',
	started,
	...more
});

describe('strategic (S-0228)', () => {
	it("reads the orchestrator's document with a missing log, or an entry with no items, as empty lists", () => {
		const doc = {
			kind: 'orchestrator' as const,
			accrued_cost: 0,
			accrued_seconds: 0,
			tasks_completed: 0,
			last_run: '',
			path: 'wip/agents/orchestrator.md',
			refusals: [{ at: '2026-10-03T10:00:00Z', call: 'flai push', needs: '' }]
		};
		expect(activityDocument({ ...doc, entries: null }).entries).toEqual([]);
		expect(activityDocument(doc).entries).toEqual([]);
		const read = activityDocument({
			...doc,
			entries: [{ ...entry('2026-10-03T10:00:00Z', 1, 60), items: null }]
		});
		expect(read.entries[0].items).toEqual([]);
		expect(read.entries[0].summary).toContain('first under policy cod');
		expect(read.refusals).toEqual(doc.refusals);
		expect(read.kind).toBe('orchestrator');
	});

	it('lists the log and the runs newest first, leaving what it was given as it was', () => {
		const log = [
			entry('2026-10-01T09:00:00Z', 1, 10),
			entry('2026-10-03T09:00:00Z', 3, 30),
			entry('2026-10-02T09:00:00Z', 2, 20)
		];
		expect(newestFirst(log).map((e) => e.cost)).toEqual([3, 2, 1]);
		expect(log.map((e) => e.cost)).toEqual([1, 3, 2]);
		const runs = [
			analyzer('2026-10-01T09:00:00Z', { focus: 'intent' }),
			analyzer('2026-10-03T09:00:00Z', { focus: 'risk' }),
			analyzer('2026-10-02T09:00:00Z', { focus: 'all' })
		];
		expect(runsNewestFirst(runs).map((r) => r.focus)).toEqual(['risk', 'all', 'intent']);
		expect(runs.map((r) => r.focus)).toEqual(['intent', 'risk', 'all']);
	});

	it('costs a run by the entries that ended while it ran, its first and last seconds included', () => {
		const r = analyzer('2026-10-03T10:00:00.700Z', {
			ended: '2026-10-03T10:30:12.400Z',
			outcome: 'worked',
			report: 'design/analysis/2026-10-03-risk.md'
		});
		const log = [
			entry('2026-10-03T09:59:59Z', 100, 100),
			entry('2026-10-03T10:00:00Z', 1, 10),
			entry('2026-10-03T10:30:13Z', 2.5, 20, true),
			entry('2026-10-03T10:30:14Z', 100, 100)
		];
		expect(runCost(r, log)).toEqual({ cost: 3.5, seconds: 30, estimated: true, count: 2 });
		expect(runCost(r, [])).toEqual({ cost: 0, seconds: 0, estimated: false, count: 0 });
	});

	it("costs the orchestrator's run under way up to now, and a run that could not start as nothing", () => {
		const now = new Date('2026-10-03T11:00:00Z');
		const log = [
			entry('2026-10-03T10:10:00Z', 1, 10),
			entry('2026-10-03T10:40:00Z', 2, 20),
			entry('2026-10-03T11:00:01Z', 9, 90)
		];
		expect(runCost(orchestrator('2026-10-03T10:00:00Z'), log, now)).toEqual({
			cost: 3,
			seconds: 30,
			estimated: false,
			count: 2
		});
		const failed = orchestrator('2026-10-03T10:00:00Z', { error: 'claude: not found' });
		expect(runCost(failed, log, now).count).toBe(0);
	});

	it('lists every run newest first, each with its own cost', () => {
		const runs = [
			orchestrator('2026-10-01T10:00:00Z', { ended: '2026-10-01T12:00:00Z', outcome: 'stopped' }),
			orchestrator('2026-10-03T10:00:00Z', { ended: '2026-10-03T12:00:00Z', outcome: 'failed' })
		];
		const log = [entry('2026-10-01T11:00:00Z', 1, 300), entry('2026-10-03T11:00:00Z', 2, 300)];
		const past = pastRuns(runs, log);
		expect(past.map((p) => [p.run.started, p.cost.cost])).toEqual([
			['2026-10-03T10:00:00Z', 2],
			['2026-10-01T10:00:00Z', 1]
		]);
		expect(pastRuns([], log)).toEqual([]);
	});

	it('finds the newest run under way, passing over ended runs and runs that could not start', () => {
		const ended = orchestrator('2026-10-01T10:00:00Z', {
			ended: '2026-10-01T10:05:00Z',
			held: true
		});
		const older = orchestrator('2026-10-02T10:00:00Z');
		const newer = orchestrator('2026-10-03T10:00:00Z');
		const failed = orchestrator('2026-10-04T10:00:00Z', { error: 'claude: not found' });
		expect(currentRun([ended, older, failed, newer])).toBe(newer);
		expect(currentRun([ended, failed])).toBeNull();
		expect(currentRun<AnalyzerRun>([])).toBeNull();
	});

	it('says how a run went in a word', () => {
		const at = '2026-10-03T10:00:00Z';
		const done = { ended: '2026-10-03T10:05:00Z' };
		expect(outcomeWord(analyzer(at))).toBe('running');
		expect(outcomeWord(analyzer(at, { error: 'claude: not found' }))).toBe('could not start');
		expect(outcomeWord(orchestrator(at, { ...done, outcome: 'stopped', held: true }))).toBe(
			'stopped'
		);
		expect(outcomeWord(analyzer(at, { ...done, outcome: 'worked' }))).toBe('worked');
		expect(outcomeWord(analyzer(at, done))).toBe('ended');
	});
});
