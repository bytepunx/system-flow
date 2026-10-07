// S-0228: the analyzer's panel shows the analyze host action on, or off with the command that turns
// it on; the run under way with its focus and its stream read by role; its log newest first, each
// entry with the report it wrote linked on the documents page; its runs with their focus, report,
// cost, and outcome; and Run with a focus or none through /api/analyzer, saying what flai answered,
// disabled, saying why, while the action is off or a run goes.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { AnalyzerRun, AnalyzerView } from '$lib/strategic';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[id]', params.id ?? '').replace('[...path]', params.path ?? '')
}));

import AnalyzerPanel from './AnalyzerPanel.svelte';

const settle = async () => {
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const one = (id: string) => document.querySelector<HTMLElement>(`[data-testid="${id}"]`);
const all = (id: string) => [...document.querySelectorAll<HTMLElement>(`[data-testid="${id}"]`)];
const text = (el: Element | null) => el!.textContent!.replace(/\s+/g, ' ').trim();
const button = (id: string) => one(id) as HTMLButtonElement | null;
const href = (el: Element | null) => el!.querySelector('a')!.getAttribute('href');
const choose = (value: string) => {
	const select = one('analyzer-focus') as HTMLSelectElement;
	select.value = value;
	select.dispatchEvent(new Event('change', { bubbles: true }));
	flushSync();
};
const streamed = answer({ entries: [], running: true, from: 0, next: 0, size: 0 });

const running: AnalyzerRun = {
	story: '',
	agent: 'analyzer',
	harness: 'claude-code',
	model: 'sonnet',
	command: 'claude',
	pid: 51,
	started: '2026-10-05T10:00:00Z',
	focus: 'risk',
	trigger: 'asked'
};
const worked: AnalyzerRun = {
	story: '',
	agent: 'analyzer',
	command: 'claude',
	started: '2026-10-04T09:00:00Z',
	ended: '2026-10-04T09:20:00Z',
	outcome: 'worked',
	focus: 'all',
	trigger: 'asked',
	report: 'design/analysis/2026-10-04-all.md'
};
const failed: AnalyzerRun = {
	story: '',
	agent: 'analyzer',
	command: 'claude',
	started: '2026-10-03T08:00:00Z',
	ended: '2026-10-03T08:05:00Z',
	outcome: 'failed',
	why: 'ended (exit 0) without writing a report under design/analysis',
	focus: 'bottlenecks',
	trigger: '0 6 * * 1'
};
const view = (over: Partial<AnalyzerView> = {}): AnalyzerView => ({
	enabled: true,
	activity: {
		kind: 'analyzer',
		accrued_cost: 3.5,
		accrued_seconds: 1500,
		tasks_completed: 3,
		last_run: '2026-10-04T09:20:00Z',
		path: 'wip/agents/analyzer.md',
		entries: [
			{
				at: '2026-10-02T07:10:00Z',
				summary: 'Wrote design/analysis/2026-10-02-intent.md: two gaps between the design and flai',
				trigger: 'asked',
				items: ['I-0070'],
				seconds: 300,
				cost: 0.5,
				estimated: false
			},
			{
				at: '2026-10-03T08:05:00Z',
				summary: 'Found nothing to report (no report written under design/analysis)',
				trigger: '0 6 * * 1',
				items: [],
				seconds: 300,
				cost: 1,
				estimated: true
			},
			{
				at: '2026-10-04T09:20:00Z',
				summary: 'Three bottlenecks and one risk (report design/analysis/2026-10-04-all.md)',
				trigger: 'asked',
				items: ['I-0071', 'I-0072'],
				seconds: 900,
				cost: 2,
				estimated: false
			}
		]
	},
	run: null,
	runs: [],
	...over
});

describe('the analyzer panel (S-0228)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let changed = 0;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		changed = 0;
		document.body.innerHTML = '';
	});
	const show = async (v: AnalyzerView) => {
		c = mount(AnalyzerPanel, {
			target: document.body,
			props: { view: v, onchanged: () => changed++ }
		});
		await settle();
	};

	describe('its state', () => {
		it('is on, with Run and a focus to choose', async () => {
			await show(view({ runs: [worked] }));
			expect(text(one('analyzer-action-on'))).toBe(
				'On: flai serve starts the analyzer when you Run it, and when analysis.schedule comes round.'
			);
			expect(one('analyzer-action-off')).toBeNull();
			expect(button('analyzer-run')!.disabled).toBe(false);
			const select = one('analyzer-focus') as HTMLSelectElement;
			expect(select.disabled).toBe(false);
			expect([...select.options].map((o) => [o.value, o.textContent])).toEqual([
				['', 'all three'],
				['bottlenecks', 'bottlenecks'],
				['intent', 'intent'],
				['risk', 'risk']
			]);
			expect(one('analyzer-run-off')).toBeNull();
			expect(one('analyzer-run-busy')).toBeNull();
		});

		it('is off, with the command that turns it on, and Run disabled saying why', async () => {
			await show(view({ enabled: false, runs: [worked] }));
			const off = one('analyzer-action-off')!;
			expect(text(off)).toBe(
				'Off: flai serve starts no analyzer for this project. On the host, in the project, run flai serve enable analyze.'
			);
			expect(off.querySelector('code')!.textContent).toBe('flai serve enable analyze');
			expect(button('analyzer-run')!.disabled).toBe(true);
			expect((one('analyzer-focus') as HTMLSelectElement).disabled).toBe(true);
			const why = one('analyzer-run-off')!;
			expect(text(why)).toContain('The analyze host action is off for this project');
			expect(why.querySelector('code')!.textContent).toBe('flai serve enable analyze');
			expect(one('analyzer-action-on')).toBeNull();
		});

		it('disables Run, saying why, while a run goes', async () => {
			api.mockResolvedValue(streamed);
			await show(view({ run: running, runs: [running, worked] }));
			expect(button('analyzer-run')!.disabled).toBe(true);
			expect((one('analyzer-focus') as HTMLSelectElement).disabled).toBe(true);
			expect(text(one('analyzer-run-busy'))).toBe(
				'An analyzer is running, focus risk; one runs at a time, so Run waits for it to end.'
			);
		});
	});

	it('shows the run under way with its focus and its stream, read by the role analyze', async () => {
		api.mockResolvedValue(
			answer({
				story: '',
				agent: 'analyzer',
				started: running.started,
				running: true,
				from: 0,
				next: 10,
				size: 10,
				entries: [{ kind: 'text', text: 'Reading the cycle times.' }]
			})
		);
		await show(view({ run: running, runs: [running, worked] }));
		const box = one('analyzer-current')!;
		expect(text(one('analyzer-current-focus'))).toBe('risk');
		expect(text(box)).toContain('analyzer (claude-code, sonnet)');
		// S-0329: started at 10:00 UTC, shown in the local zone, New York's under test
		expect(text(box)).toContain('2026-10-05 06:00 EDT');
		expect(text(box)).not.toContain('UTC');
		expect(text(one('analyzer-current-trigger'))).toBe('asked');
		expect(api).toHaveBeenCalledWith('/api/agent-stream/analyzer?role=analyze');
		expect(text(box.querySelector('[data-testid="agent-stream"]'))).toContain(
			'Reading the cycle times.'
		);
	});

	it('says no analyzer is running when none goes', async () => {
		await show(view({ runs: [worked] }));
		expect(text(one('analyzer-current-none'))).toBe('No analyzer is running.');
		expect(api).not.toHaveBeenCalled();
	});

	it('shows the totals of its activity document', async () => {
		await show(view());
		expect(text(one('analyzer-accrued-cost'))).toBe('$3.50');
		expect(text(one('analyzer-accrued-time'))).toBe('25m');
		expect(text(one('analyzer-activities'))).toBe('3');
		expect(text(one('analyzer-last-run'))).toBe('2026-10-04 05:20 EDT');
		expect(text(document.querySelector('[data-testid="analyzer-activity"] h2'))).toBe(
			'Activity in wip/agents/analyzer.md'
		);
	});

	it('shows its entries newest first, each with its items, duration, cost, and the report it wrote', async () => {
		await show(view({ runs: [worked, failed] }));
		expect(all('analyzer-entry-summary').map(text)).toEqual([
			'Three bottlenecks and one risk (report design/analysis/2026-10-04-all.md)',
			'Found nothing to report (no report written under design/analysis)',
			'Wrote design/analysis/2026-10-02-intent.md: two gaps between the design and flai'
		]);
		expect(all('analyzer-entry-details').map(text)).toEqual([
			'trigger asked; items I-0071, I-0072; took 15m; cost $2.00',
			'trigger 0 6 * * 1; items none; took 5m; cost $1.00 (estimated)',
			'trigger asked; items I-0070; took 5m; cost $0.500'
		]);
		expect(all('analyzer-entry').map((e) => text(e).slice(0, 20))).toEqual([
			'2026-10-04 05:20 EDT',
			'2026-10-03 04:05 EDT',
			'2026-10-02 03:10 EDT'
		]);
		expect(text(one('analyzer-activity'))).not.toContain('UTC');
		const reports = all('analyzer-entry-report');
		expect(reports.map(text)).toEqual([
			'report design/analysis/2026-10-04-all.md',
			'no report',
			// a report of a run flai no longer keeps is read from the summary
			'report design/analysis/2026-10-02-intent.md'
		]);
		expect(href(reports[0])).toBe('/docs/design/analysis/2026-10-04-all.md');
		expect(reports[1].querySelector('a')).toBeNull();
		expect(href(reports[2])).toBe('/docs/design/analysis/2026-10-02-intent.md');
	});

	it('says when it has logged nothing', async () => {
		await show(view({ activity: { ...view().activity, entries: [] } }));
		expect(text(one('analyzer-entries-none'))).toBe('The analyzer has not logged an activity yet.');
	});

	it('lists its runs newest first with their focus, report, outcome, cost, and why one failed', async () => {
		api.mockResolvedValue(streamed);
		await show(view({ run: running, runs: [running, worked, failed] }));
		expect(all('analyzer-run-focus').map(text)).toEqual(['risk', 'all three', 'bottlenecks']);
		const reports = all('analyzer-run-report');
		expect(reports.map(text)).toEqual(['—', 'design/analysis/2026-10-04-all.md', '—']);
		expect(href(reports[1])).toBe('/docs/design/analysis/2026-10-04-all.md');
		expect(all('analyzer-run-outcome').map(text)).toEqual(['running', 'worked', 'failed']);
		expect(all('analyzer-run-cost').map(text)).toEqual(['—', '$2.00', '$1.00 (estimated)']);
		expect(all('analyzer-run-why').map(text)).toEqual([
			'',
			'',
			'ended (exit 0) without writing a report under design/analysis'
		]);
		// the Run button shares the rows' test ID
		const rows = document.querySelectorAll('tr[data-testid="analyzer-run"]');
		const cells = rows[1].querySelectorAll('td');
		expect(text(cells[2])).toBe('2026-10-04 05:00 EDT');
		expect(text(cells[3])).toBe('2026-10-04 05:20 EDT');
		expect(text(one('analyzer-runs'))).not.toContain('UTC');
	});

	it('says flai has started no analyzer when it has none', async () => {
		await show(view());
		expect(text(one('analyzer-runs-none'))).toBe(
			'flai serve has started no analyzer for this project.'
		);
	});

	describe('Run', () => {
		const posted = () =>
			api.mock.calls
				.filter(([url]) => url === '/api/analyzer')
				.map(([, init]) => ({ method: init.method, body: JSON.parse(init.body) }));
		const started = (focus: string) =>
			answer({
				focus,
				agent: 'analyzer',
				harness: 'claude-code',
				command: 'claude',
				pid: 88,
				log: '.flai-cache/serve/sf-analyzer-1.log',
				started: '2026-10-06T10:00:00Z',
				trigger: 'asked'
			});

		it('with no focus asks flai for all of them, and says what it answered', async () => {
			await show(view());
			api.mockResolvedValue(started('all'));
			button('analyzer-run')!.click();
			await settle();
			expect(posted()).toEqual([{ method: 'POST', body: {} }]);
			expect(text(one('analyzer-said'))).toBe(
				'started claude to analyze, focus all, as analyzer (pid 88); its output is in .flai-cache/serve/sf-analyzer-1.log on the host'
			);
			expect(one('analyzer-said')!.getAttribute('role')).toBe('status');
			expect(changed).toBe(1);
		});

		it('with a focus asks flai for it', async () => {
			await show(view());
			choose('intent');
			api.mockResolvedValue(started('intent'));
			button('analyzer-run')!.click();
			await settle();
			expect(posted()).toEqual([{ method: 'POST', body: { focus: 'intent' } }]);
			expect(text(one('analyzer-said'))).toContain('to analyze, focus intent, as analyzer');
			expect(changed).toBe(1);
		});

		it('says why flai refused', async () => {
			await show(view());
			choose('risk');
			api.mockResolvedValue(
				answer(
					{
						error:
							'the analyze host action is off for this project: on the host, run flai serve enable analyze'
					},
					403
				)
			);
			button('analyzer-run')!.click();
			await settle();
			expect(posted()).toEqual([{ method: 'POST', body: { focus: 'risk' } }]);
			expect(text(one('analyzer-said'))).toBe(
				'refused: the analyze host action is off for this project: on the host, run flai serve enable analyze'
			);
			expect(changed).toBe(1);
			expect(button('analyzer-run')!.disabled).toBe(false);
		});
	});
});
