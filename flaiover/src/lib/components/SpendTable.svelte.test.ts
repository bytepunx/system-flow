import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import SpendTable from './SpendTable.svelte';
import { normalise, spendRows, type Bucket, type Report } from '$lib/viz/charts';

const day = (at: string, items: number, tokens: number, cost: number, rest: Partial<Bucket>) =>
	({
		at,
		items,
		tokens,
		cost,
		seconds: 0,
		mean_tokens: 0,
		mean_cost: 0,
		...rest
	}) as Bucket;
const opus = {
	model: 'claude-opus-5-5',
	items: 2,
	tokens: 2500000,
	cost: 1.5,
	seconds: 5400,
	tokens_per_item: 1250000,
	cost_per_item: 0.75,
	tokens_per_minute: 27777.8,
	tokens_per_dollar: 1666666.7
};
const haiku = {
	model: 'claude-haiku-4-5',
	items: 1,
	tokens: 1000000,
	cost: 0.25,
	seconds: 3600,
	tokens_per_item: 1000000,
	cost_per_item: 0.25,
	tokens_per_minute: 16666.7,
	tokens_per_dollar: 4000000
};
const report = normalise({
	type: 'story',
	items: [],
	usage: {
		bucket: 'day',
		spend: {
			epic: { items: 0, tokens: 0, cost: 0, seconds: 0, models: [], buckets: [] },
			story: {
				items: 3,
				tokens: 4500000,
				cost: 2.25,
				seconds: 5400,
				models: [],
				buckets: [
					day('2026-08-03T00:00:00Z', 2, 3500000, 1.75, {
						seconds: 5400,
						estimated: true,
						tokens_per_item: 1750000,
						cost_per_item: 0.875,
						tokens_per_minute: 38888.9,
						tokens_per_dollar: 2000000,
						mean_tokens: 3500000,
						mean_cost: 1.75,
						models: [haiku, opus]
					}),
					day('2026-08-04T00:00:00Z', 0, 0, 0, { mean_tokens: 1750000, mean_cost: 0.875 }),
					day('2026-08-05T00:00:00Z', 1, 1000000, 0.5, {
						tokens_per_item: 1000000,
						cost_per_item: 0.5,
						tokens_per_dollar: 2000000,
						mean_tokens: 1500000,
						mean_cost: 0.75,
						models: [{ ...opus, items: 1, tokens: 1000000, cost: 0.5, seconds: 0 }]
					})
				]
			},
			task: {
				items: 4,
				tokens: 2000000,
				cost: 1,
				seconds: 1200,
				models: [],
				buckets: [
					day('2026-08-03T00:00:00Z', 4, 2000000, 1, {
						seconds: 1200,
						tokens_per_item: 500000,
						cost_per_item: 0.25
					})
				]
			}
		}
	}
} as unknown as Report);

describe('SpendTable', () => {
	let component: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});
	const cells = () =>
		[...document.querySelectorAll('[data-testid="spend-table"] tbody tr')].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent?.trim())
		);

	it('lists each bucket with items, newest first, whole and per model', () => {
		const rows = spendRows(report, 'tokens-spent');
		expect(rows.map((r) => [r.at.slice(0, 10), r.of])).toEqual([
			['2026-08-05', 'story'],
			['2026-08-05', 'claude-opus-5-5'],
			['2026-08-03', 'all models'],
			['2026-08-03', 'claude-haiku-4-5'],
			['2026-08-03', 'claude-opus-5-5']
		]);
		component = mount(SpendTable, { target: document.body, props: { rows, bucket: 'day' } });
		flushSync();
		const heads = [...document.querySelectorAll('th')].map((th) => th.textContent?.trim());
		expect(heads).toContain('mean tokens per day');
		const all = cells()[2];
		// 3 August: two stories, 3.5M tokens for $1.75, estimated in part, over 90 agent minutes
		expect(all).toEqual([
			'2026-08-03',
			'all models',
			'2',
			'3.5M',
			'$1.75*',
			'90.0',
			'1.8M',
			'$0.875',
			'38.9K',
			'2.0M',
			'3.5M',
			'$1.75'
		]);
		// a model's row has no running mean, and a bucket with no agent time no rate
		expect(cells()[3].slice(-2)).toEqual(['-', '-']);
		expect(cells()[0][8]).toBe('-');
	});

	it('lists the item types side by side for a chart per item', () => {
		const rows = spendRows(report, 'cost-per-item');
		expect(rows.map((r) => [r.at.slice(0, 10), r.of, r.cost_per_item])).toEqual([
			['2026-08-05', 'story', 0.5],
			['2026-08-03', 'story', 0.875],
			['2026-08-03', 'task', 0.25]
		]);
		expect(spendRows(report, 'cost-per-model').map((r) => r.of)).toContain('claude-haiku-4-5');
		component = mount(SpendTable, { target: document.body, props: { rows, bucket: 'week' } });
		flushSync();
		expect(cells()[2].slice(0, 3)).toEqual(['week of 2026-08-03', 'task', '4']);
	});
});
