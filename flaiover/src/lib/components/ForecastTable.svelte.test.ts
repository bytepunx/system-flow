// S-0212: the table view under the planning charts lists each story done in the window with its
// forecast, its cycle time, and its errors, newest first, with a dash for what it lacks.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ForecastTable from './ForecastTable.svelte';
import { forecastRows, type ItemMetrics, type Report } from '$lib/viz/charts';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id)
}));

const story = (id: string, completed: string, rest: Partial<ItemMetrics>): ItemMetrics => ({
	id,
	type: 'story',
	nature: 'feature',
	title: `Story ${id}`,
	status: 'done',
	created: '2026-09-01T00:00:00Z',
	started: '2026-09-01T00:00:00Z',
	completed,
	blocked_seconds: 0,
	time_in_state_seconds: {},
	...rest
});
const report = {
	generated_at: '2026-09-29T21:00:00Z',
	type: 'story',
	window_days: 30,
	window_start: '2026-08-30T21:00:00Z',
	items: [
		story('S-0001', '2026-09-10T12:00:00Z', {
			model: 'claude-opus-5-5',
			estimate_error_seconds: -3600
		}),
		story('S-0003', '2026-09-28T12:00:00Z', {
			model: 'claude-opus-5-5',
			forecast_seconds: 21600,
			cycle_time_seconds: 28800,
			forecast_error_seconds: 7200,
			delivery_error_seconds: -86400,
			estimate_error_seconds: 1800
		}),
		story('S-0002', '2026-09-20T12:00:00Z', {
			nature: 'remediation',
			forecast_seconds: 14400,
			cycle_time_seconds: 3600,
			forecast_error_seconds: -10800
		}),
		story('S-0004', '2026-09-25T12:00:00Z', { cycle_time_seconds: 3600 })
	]
} as unknown as Report;

describe('ForecastTable', () => {
	let component: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});
	const show = (rows: ReturnType<typeof forecastRows>) => {
		component = mount(ForecastTable, { target: document.body, props: { rows } });
		flushSync();
	};
	const cells = () =>
		[...document.querySelectorAll('[data-testid="forecast-table"] tbody tr')].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent?.replace(/\s+/g, ' ').trim())
		);

	it('lists the stories with an error, newest first, each linked, with a dash for what it lacks', () => {
		show(forecastRows(report));
		expect([...document.querySelectorAll('th')].map((th) => th.textContent?.trim())).toEqual([
			'story',
			'completed',
			'nature',
			'model',
			'forecast',
			'actual',
			'forecast error',
			'delivery error',
			'estimate error'
		]);
		expect(cells()).toEqual([
			[
				'S-0003 Story S-0003',
				'2026-09-28',
				'feature',
				'claude-opus-5-5',
				'6h',
				'8h',
				'+2h',
				'-1d',
				'+30m'
			],
			['S-0002 Story S-0002', '2026-09-20', 'remediation', '(none)', '4h', '1h', '-3h', '-', '-'],
			['S-0001 Story S-0001', '2026-09-10', 'feature', 'claude-opus-5-5', '-', '-', '-', '-', '-1h']
		]);
		expect(
			[...document.querySelectorAll('[data-testid="forecast-table"] a')].map((a) =>
				a.getAttribute('href')
			)
		).toEqual(['/items/S-0003', '/items/S-0002', '/items/S-0001']);
	});

	it('lists only the stories of the nature and the model chosen', () => {
		show(forecastRows(report, { nature: 'feature', model: 'claude-opus-5-5' }));
		expect(cells().map((r) => r[0])).toEqual(['S-0003 Story S-0003', 'S-0001 Story S-0001']);
	});

	it('lists no row when no story in the window has an error', () => {
		show([]);
		expect(cells()).toEqual([]);
		expect(document.querySelectorAll('[data-testid="forecast-table"] th').length).toBe(9);
	});
});
