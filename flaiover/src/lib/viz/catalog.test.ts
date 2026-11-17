import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import {
	CHART_GROUPS,
	COD_KINDS,
	FLOW_KINDS,
	KINDS,
	PLANNING_KINDS,
	STRATEGIC_KINDS,
	isTimerKind,
	timerCharts,
	timerKind,
	timerName,
	TITLES,
	USAGE_KINDS
} from './catalog';

describe('CHART_GROUPS', () => {
	it('lists every kind exactly once', () => {
		const listed = CHART_GROUPS.flatMap((g) => g.kinds);
		expect(listed).toHaveLength(KINDS.length);
		expect([...listed].sort()).toEqual([...KINDS].sort());
	});

	it('lists the cost of delay charts under planning, after the forecasts and claims', () => {
		expect(CHART_GROUPS.map((g) => g.kinds)).toEqual([
			FLOW_KINDS,
			USAGE_KINDS,
			[...PLANNING_KINDS, ...COD_KINDS],
			STRATEGIC_KINDS
		]);
	});

	it('lists Concurrency under flow, after Throughput (S-0551)', () => {
		const flow = CHART_GROUPS.find((g) => g.key === 'flow')!.kinds;
		expect(flow).toContain('concurrency');
		expect(flow.indexOf('concurrency')).toBe(flow.indexOf('throughput') + 1);
		expect(CHART_GROUPS.filter((g) => g.kinds.includes('concurrency'))).toHaveLength(1);
		expect(TITLES.concurrency).toBe('Concurrency');
	});

	it('names the groups flow, usage, planning, strategic, in that order', () => {
		expect(CHART_GROUPS.map((g) => g.key)).toEqual(['flow', 'usage', 'planning', 'strategic']);
		expect(CHART_GROUPS.map((g) => g.label)).toEqual(['flow', 'usage', 'planning', 'strategic']);
	});
});

describe('TITLES', () => {
	it('gives every kind a title', () => {
		for (const kind of KINDS) expect(TITLES[kind], kind).toMatch(/\S/);
	});
});

// S-0754: a chart for each of the workflow's own timers, at /charts/timer-<name>
describe('the charts of the workflow timers', () => {
	it('names a timer chart by its timer, and no catalogue kind as one', () => {
		expect(timerKind('testing-wait')).toBe('timer-testing-wait');
		expect(timerName('timer-testing-wait')).toBe('testing-wait');
		expect(isTimerKind('timer-testing-wait')).toBe(true);
		expect(isTimerKind('timer-')).toBe(false);
		expect(timerName('timer-')).toBeUndefined();
		for (const kind of KINDS) {
			expect(isTimerKind(kind), kind).toBe(false);
			expect(timerName(kind), kind).toBeUndefined();
		}
	});

	it('lists the timer charts in the order given, titled by the timer, else its name', () => {
		expect(
			timerCharts([
				{ name: 'testing-wait', title: 'Testing Wait' },
				{ name: 'rework', title: '' },
				{ name: 'soak' }
			])
		).toEqual([
			{ kind: 'timer-testing-wait', title: 'Testing Wait' },
			{ kind: 'timer-rework', title: 'rework' },
			{ kind: 'timer-soak', title: 'soak' }
		]);
		expect(timerCharts([])).toEqual([]);
	});
});

describe('catalog.ts', () => {
	it('imports neither ECharts nor the chart builders', () => {
		const source = readFileSync(new URL('./catalog.ts', import.meta.url), 'utf8');
		expect(source).not.toMatch(/from\s+['"]echarts/);
		expect(source).not.toMatch(/from\s+['"]\.\/charts['"]/);
		expect(source).not.toMatch(/^\s*import\b/m);
	});
});
