import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, rm } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo } from './repo';
import { resetFlaiBinary } from './flai';
import { clearStatsCache, stats, validate } from './stats';
import type { Report } from '$lib/viz/charts';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');

describe('stats query validation', () => {
	it('builds flai arguments and rejects bad values', () => {
		expect(validate({})).toEqual(['stats']);
		expect(validate({ since: '90d', type: 'task', by: 'nature' })).toEqual([
			'stats',
			'--since',
			'90d',
			'--type',
			'task',
			'--by',
			'nature'
		]);
		expect(() => validate({ since: 'soon' })).toThrow(/since/);
		expect(() => validate({ type: 'bug' })).toThrow(/type/);
		expect(() => validate({ by: 'owner' })).toThrow(/by/);
	});
});

describe.skipIf(!existsSync(bin))('stats through flai on the fixture', () => {
	let dir: string;
	beforeAll(async () => {
		process.env.FLAI_BIN = bin;
		resetFlaiBinary();
		clearStatsCache();
		dir = await mkdtemp(join(tmpdir(), 'flaiover-stats-'));
		await cp(fixture, dir, { recursive: true });
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});
	it('returns the report shape the charts consume', async () => {
		const r = await stats<Report>(new Repo(dir), { since: '365d' });
		expect(r.type).toBe('story');
		expect(r.summary.completed).toBeGreaterThanOrEqual(2);
		expect(r.summary.cycle_time.p85_seconds).toBeGreaterThan(0);
		expect(Object.keys(r.burnup)).toContain('all');
		expect(r.cfd.length).toBeGreaterThan(0);
		expect(r.throughput.length).toBeGreaterThan(0);
		expect(r.items.some((i) => i.estimate_seconds)).toBe(true);
	});
});
