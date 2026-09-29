import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, rm } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo } from './repo';
import { flaiAsk } from './testing';
import { stats } from './stats';
import { build, hasSpend, normalise, type Report } from '$lib/viz/charts';
import { theme } from '$lib/viz/palette';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');

describe.skipIf(!existsSync(bin))('stats through flai on the fixture', () => {
	let dir: string;
	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-stats-'));
		await cp(fixture, dir, { recursive: true });
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});
	it('returns the report shape the charts consume', async () => {
		const r = await stats<Report>(new Repo(dir, flaiAsk(dir)), { since: '365d' });
		expect(r.type).toBe('story');
		expect(r.summary.completed).toBeGreaterThanOrEqual(2);
		expect(r.summary.cycle_time.p85_seconds).toBeGreaterThan(0);
		expect(Object.keys(r.burnup)).toContain('all');
		expect(r.cfd.length).toBeGreaterThan(0);
		expect(r.throughput.length).toBeGreaterThan(0);
		expect(r.items.some((i) => i.estimate_seconds)).toBe(true);
	});
	it('returns spend over time for every item type, in the bucket asked for', async () => {
		const repo = new Repo(dir, flaiAsk(dir));
		const r = normalise(await stats<Report>(repo, { since: '365d', type: 'task' }));
		expect(hasSpend(r)).toBe(true);
		expect(r.usage?.bucket).toBe('day');
		const stories = r.usage!.spend!.story;
		// S-001 and S-002 carry usage: 3 000 000 tokens for $1.50, and 1 100 000 for $0.52
		expect([stories.items, stories.tokens, stories.seconds]).toEqual([2, 4100000, 1800]);
		expect(stories.cost).toBeCloseTo(2.02);
		expect(stories.tokens_per_item).toBeCloseTo(2050000);
		expect(stories.tokens_per_minute).toBeCloseTo(4100000 / 30);
		expect(stories.buckets[0]).toMatchObject({
			at: '2026-08-03T00:00:00Z',
			items: 1,
			tokens: 3000000,
			tokens_per_item: 3000000,
			mean_tokens: 3000000
		});
		const second = stories.buckets.find((b) => b.at === '2026-08-12T00:00:00Z');
		expect(second?.estimated).toBe(true);
		expect(second?.models?.map((m) => [m.model, m.tokens])).toEqual([
			['claude-haiku-4-5', 100000],
			['claude-opus-5-5', 1000000]
		]);
		expect(second?.mean_tokens).toBeCloseTo(410000);
		expect(r.usage!.spend!.task.items).toBe(1);
		expect(r.usage!.spend!.epic.buckets).toEqual([]);
		// the builders draw it: the item types side by side, and the day's tokens by model
		const per = build('tokens-per-item', r, theme(false)) as { series: { name: string }[] };
		expect(per.series.map((s) => s.name)).toEqual(['story', 'task']);
		const weeks = await stats<Report>(repo, { since: '365d', bucket: 'week' });
		expect(weeks.usage?.bucket).toBe('week');
		expect(weeks.usage?.spend?.story.buckets[0].at).toBe('2026-08-03T00:00:00Z');
		await expect(stats<Report>(repo, { since: '365d', bucket: 'hour' })).rejects.toThrow(
			/31 days or less/
		);
	});
});
