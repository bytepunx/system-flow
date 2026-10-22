import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import {
	CHART_GROUPS,
	COD_KINDS,
	FLOW_KINDS,
	KINDS,
	PLANNING_KINDS,
	STRATEGIC_KINDS,
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

describe('catalog.ts', () => {
	it('imports neither ECharts nor the chart builders', () => {
		const source = readFileSync(new URL('./catalog.ts', import.meta.url), 'utf8');
		expect(source).not.toMatch(/from\s+['"]echarts/);
		expect(source).not.toMatch(/from\s+['"]\.\/charts['"]/);
		expect(source).not.toMatch(/^\s*import\b/m);
	});
});
