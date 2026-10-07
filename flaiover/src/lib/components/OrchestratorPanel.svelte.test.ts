// S-0228: the orchestrator's panel shows the orchestrate host action on, off with the command that
// turns it on, or held after a Stop; the run under way with its stream read by role; its last
// decisions with their reasons above its whole log; its runs with their cost and outcome; and Stop
// and Start through /api/orchestrator, each saying what flai answered, both disabled while the
// action is off.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { ActivityEntry, OrchestratorRun, OrchestratorView } from '$lib/strategic';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id ?? '')
}));

import OrchestratorPanel from './OrchestratorPanel.svelte';

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

const running: OrchestratorRun = {
	story: '',
	agent: 'orchestrator',
	harness: 'claude-code',
	model: 'opus',
	command: 'claude',
	pid: 42,
	started: '2026-10-03T10:00:00Z'
};
const stopped: OrchestratorRun = {
	story: '',
	agent: 'orchestrator',
	command: 'claude',
	started: '2026-10-02T09:00:00Z',
	ended: '2026-10-02T09:30:00Z',
	outcome: 'stopped'
};
const failed: OrchestratorRun = {
	story: '',
	agent: 'orchestrator',
	command: 'claude',
	started: '2026-10-01T08:00:00Z',
	ended: '2026-10-01T08:01:00Z',
	outcome: 'failed',
	why: 'the harness exited 1'
};
/** n decisions a minute apart from 10:00 on 2026-10-02, oldest first as flai answers them. */
const decided = (n: number): ActivityEntry[] =>
	Array.from({ length: n }, (_, i) => ({
		at: `2026-10-02T09:${String(10 + i).padStart(2, '0')}:00Z`,
		summary: `decision ${i + 1}: because ${i + 1}`,
		items: i % 2 ? [] : ['S-0252'],
		seconds: 60,
		cost: 1.5,
		estimated: false
	}));
const view = (over: Partial<OrchestratorView> = {}): OrchestratorView => ({
	enabled: true,
	held: false,
	activity: {
		kind: 'orchestrator',
		accrued_cost: 1.2,
		accrued_seconds: 900,
		tasks_completed: 2,
		last_run: '2026-10-02T09:30:00Z',
		path: 'wip/agents/orchestrator.md',
		entries: [
			{
				at: '2026-10-02T09:10:00Z',
				summary: 'promoted S-0252: cost of delay 60.81 USD a week, first under policy cod',
				items: ['S-0252'],
				seconds: 300,
				cost: 0.4,
				estimated: false
			},
			{
				at: '2026-10-02T09:20:00Z',
				summary: 'answered TH-0107: the design names the split',
				items: [],
				seconds: 600,
				cost: 0.8,
				estimated: true
			}
		]
	},
	run: null,
	runs: [],
	...over
});

describe('the orchestrator panel (S-0228)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let changed = 0;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		changed = 0;
		document.body.innerHTML = '';
	});
	const show = async (v: OrchestratorView) => {
		c = mount(OrchestratorPanel, {
			target: document.body,
			props: { view: v, onchanged: () => changed++ }
		});
		await settle();
	};

	describe('its state', () => {
		it('is on, with Stop while a run goes', async () => {
			api.mockResolvedValue(answer({ entries: [], running: true, from: 0, next: 0, size: 0 }));
			await show(view({ run: running, runs: [running, stopped] }));
			expect(text(one('orchestrator-action-on'))).toContain(
				'flai serve runs the orchestrator for as long as the action stays on'
			);
			expect(one('orchestrator-action-off')).toBeNull();
			expect(one('orchestrator-action-held')).toBeNull();
			expect(button('orchestrator-stop')!.disabled).toBe(false);
			expect(one('orchestrator-start')).toBeNull();
			expect(one('orchestrator-controls-off')).toBeNull();
		});

		it('is off, with the command that turns it on, and Stop and Start disabled saying why', async () => {
			await show(view({ enabled: false, runs: [stopped] }));
			const off = one('orchestrator-action-off')!;
			expect(text(off)).toContain(
				'Off: flai serve starts no orchestrator for this project. On the host, in the project, run flai serve enable orchestrate'
			);
			expect(off.querySelector('code')!.textContent).toBe('flai serve enable orchestrate');
			expect(button('orchestrator-stop')!.disabled).toBe(true);
			expect(button('orchestrator-start')!.disabled).toBe(true);
			const why = one('orchestrator-controls-off')!;
			expect(text(why)).toContain('The orchestrate host action is off for this project');
			expect(why.querySelector('code')!.textContent).toBe('flai serve enable orchestrate');
			expect(one('orchestrator-action-on')).toBeNull();
		});

		it('is held after a Stop, with Start and no Stop', async () => {
			await show(view({ held: true, runs: [{ ...stopped, held: true }] }));
			expect(text(one('orchestrator-action-held'))).toBe(
				'Held: you stopped the orchestrator, and flai serve does not start it again until you Start it.'
			);
			expect(one('orchestrator-action-on')).toBeNull();
			expect(button('orchestrator-start')!.disabled).toBe(false);
			expect(one('orchestrator-stop')).toBeNull();
		});

		it('offers no Stop before flai has started an orchestrator', async () => {
			await show(view());
			expect(one('orchestrator-stop')).toBeNull();
			expect(one('orchestrator-start')).toBeNull();
		});
	});

	it('shows the run under way with its stream, read by the role orchestrate', async () => {
		api.mockResolvedValue(
			answer({
				story: '',
				agent: 'orchestrator',
				started: running.started,
				running: true,
				from: 0,
				next: 10,
				size: 10,
				entries: [{ kind: 'text', text: 'Ordering the ready column.' }]
			})
		);
		await show(view({ run: running, runs: [running, stopped] }));
		const box = one('orchestrator-current')!;
		expect(text(box)).toContain('orchestrator (claude-code, opus)');
		expect(text(box)).toContain('2026-10-03 10:00 UTC');
		expect(api).toHaveBeenCalledWith('/api/agent-stream/orchestrator?role=orchestrate');
		expect(text(box.querySelector('[data-testid="agent-stream"]'))).toContain(
			'Ordering the ready column.'
		);
		expect(one('orchestrator-current-none')).toBeNull();
	});

	it('says no orchestrator is running when none goes', async () => {
		await show(view({ runs: [stopped] }));
		expect(text(one('orchestrator-current-none'))).toBe('No orchestrator is running.');
		expect(api).not.toHaveBeenCalled();
	});

	it('shows its last decisions newest first, with their reasons, above the whole log', async () => {
		await show(view({ activity: { ...view().activity, entries: decided(7) } }));
		expect(all('orchestrator-decision-summary').map(text)).toEqual([
			'decision 7: because 7',
			'decision 6: because 6',
			'decision 5: because 5',
			'decision 4: because 4',
			'decision 3: because 3'
		]);
		expect(text(all('orchestrator-decision')[0])).toBe(
			'2026-10-02 09:16 UTC decision 7: because 7'
		);
		// the whole log follows, every entry with its items, duration, and cost
		const sections = [...document.querySelectorAll('section')].map((s) =>
			s.getAttribute('data-testid')
		);
		expect(sections.indexOf('orchestrator-decisions')).toBeLessThan(
			sections.indexOf('orchestrator-activity')
		);
		expect(all('orchestrator-entry')).toHaveLength(7);
		expect(text(all('orchestrator-entry-details')[0])).toBe('items S-0252; took 1m; cost $1.50');
	});

	it('says when it has decided nothing', async () => {
		await show(view({ activity: { ...view().activity, entries: [] } }));
		expect(text(one('orchestrator-decisions-none'))).toBe(
			'The orchestrator has decided nothing yet.'
		);
		expect(text(one('orchestrator-entries-none'))).toBe(
			'The orchestrator has not logged an activity yet.'
		);
	});

	it('shows the totals of its activity document', async () => {
		await show(view());
		expect(text(one('orchestrator-accrued-cost'))).toBe('$1.20');
		expect(text(one('orchestrator-accrued-time'))).toBe('15m');
		expect(text(one('orchestrator-activities'))).toBe('2');
		expect(text(one('orchestrator-last-run'))).toBe('2026-10-02 09:30 UTC');
		expect(text(document.querySelector('[data-testid="orchestrator-activity"] h2'))).toBe(
			'Activity in wip/agents/orchestrator.md'
		);
	});

	it('lists its runs newest first with their outcome, cost, and why one failed', async () => {
		api.mockResolvedValue(answer({ entries: [], running: true, from: 0, next: 0, size: 0 }));
		await show(view({ run: running, runs: [running, failed, stopped] }));
		expect(all('orchestrator-run-outcome').map(text)).toEqual(['running', 'stopped', 'failed']);
		// the stopped run cost the two entries logged within it
		expect(all('orchestrator-run-cost').map(text)).toEqual(['—', '$1.20 (estimated)', '—']);
		expect(all('orchestrator-run-why').map(text)).toEqual(['', '', 'the harness exited 1']);
		expect(text(all('orchestrator-run')[1].querySelectorAll('td')[0])).toBe('2026-10-02 09:00 UTC');
	});

	it('says flai has started no orchestrator when it has none', async () => {
		await show(view());
		expect(text(one('orchestrator-runs-none'))).toBe(
			'flai serve has started no orchestrator for this project.'
		);
	});

	describe('Stop and Start', () => {
		const posted = () =>
			api.mock.calls
				.filter(([url]) => url === '/api/orchestrator')
				.map(([, init]) => ({ method: init.method, body: JSON.parse(init.body) }));

		it('Stop asks flai to stop and hold it, and says what it answered', async () => {
			api.mockResolvedValue(answer({ entries: [], running: true, from: 0, next: 0, size: 0 }));
			await show(view({ run: running, runs: [running] }));
			api.mockResolvedValue(
				answer({ held: true, orchestrator: { ...running, held: true, ended: running.started } })
			);
			button('orchestrator-stop')!.click();
			await settle();
			expect(posted()).toEqual([{ method: 'POST', body: { action: 'stop' } }]);
			expect(text(one('orchestrator-said'))).toBe(
				'stopped the orchestrator orchestrator (pid 42) and held it stopped'
			);
			expect(one('orchestrator-said')!.getAttribute('role')).toBe('status');
			expect(changed).toBe(1);
		});

		it('Stop of a run that is not going says flai held it as it was', async () => {
			await show(view({ runs: [failed] }));
			api.mockResolvedValue(answer({ held: true, orchestrator: { ...failed, held: true } }));
			button('orchestrator-stop')!.click();
			await settle();
			expect(text(one('orchestrator-said'))).toBe('held the orchestrator orchestrator stopped');
		});

		it('Start asks flai to lift the hold, and says what it answered', async () => {
			await show(view({ held: true, runs: [{ ...stopped, held: true }] }));
			api.mockResolvedValue(answer({ held: false, orchestrator: { ...running, pid: 77 } }));
			button('orchestrator-start')!.click();
			await settle();
			expect(posted()).toEqual([{ method: 'POST', body: { action: 'start' } }]);
			expect(text(one('orchestrator-said'))).toBe(
				'lifted the hold and started claude as orchestrator (pid 77)'
			);
			expect(changed).toBe(1);
		});

		it('says why flai refused', async () => {
			await show(view({ held: true, runs: [{ ...stopped, held: true }] }));
			api.mockResolvedValue(
				answer(
					{
						error:
							'the orchestrate host action is off for this project, so flai serve runs no orchestrator to start: flai serve enable orchestrate turns it on, and starts it'
					},
					403
				)
			);
			button('orchestrator-start')!.click();
			await settle();
			expect(text(one('orchestrator-said'))).toBe(
				'refused: the orchestrate host action is off for this project, so flai serve runs no orchestrator to start: flai serve enable orchestrate turns it on, and starts it'
			);
			expect(changed).toBe(1);
			expect(button('orchestrator-start')!.disabled).toBe(false);
		});
	});
});
