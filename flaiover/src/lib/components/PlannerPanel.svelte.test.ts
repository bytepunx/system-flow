// S-0259: the planner's page shows the plan host action, the run under way with its stream, the
// activity document newest first with its totals, the past runs with their outcome and cost, and a
// form that asks flai serve to plan an epic or a story.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { PlanRun } from '$lib/activity';
import type { PlannerView } from '$lib/planner';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id ?? '')
}));

import PlannerPanel from './PlannerPanel.svelte';

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

const running: PlanRun = {
	item: 'E-0016',
	agent: 'planner',
	harness: 'claude-code',
	model: 'opus',
	command: 'claude',
	pid: 42,
	log: '/home/op/.flai/serve/plan-E-0016.log',
	started: '2026-10-03T10:00:00Z',
	trigger: 'asked'
};
const worked: PlanRun = {
	item: 'S-0259',
	agent: 'planner',
	command: 'claude',
	started: '2026-10-02T09:00:00Z',
	ended: '2026-10-02T09:10:00Z',
	outcome: 'worked'
};
const failed: PlanRun = {
	item: 'S-0001',
	agent: 'planner',
	command: 'claude',
	started: '2026-10-01T08:00:00Z',
	ended: '2026-10-01T08:01:00Z',
	outcome: 'failed',
	why: 'the harness exited 1'
};
const view = (over: Partial<PlannerView> = {}): PlannerView => ({
	plan_enabled: true,
	activity: {
		kind: 'planner',
		accrued_cost: 3.5,
		accrued_seconds: 5400,
		tasks_completed: 2,
		last_run: '2026-10-02T09:10:00Z',
		path: 'wip/agents/planner.md',
		entries: [
			{
				at: '2026-10-01T08:00:30Z',
				summary: 'drafted tasks',
				items: [],
				seconds: 600,
				cost: 1.25,
				estimated: false
			},
			{
				at: '2026-10-02T09:10:00Z',
				summary: 'enriched the story',
				trigger: 'asked',
				items: ['S-0259', 'T-0838'],
				seconds: 1200,
				cost: 2.25,
				estimated: true
			}
		]
	},
	runs: [],
	...over
});

describe('the planner panel (S-0259)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let planned = 0;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		planned = 0;
		document.body.innerHTML = '';
	});
	const show = async (v: PlannerView) => {
		c = mount(PlannerPanel, {
			target: document.body,
			props: { view: v, onplanned: () => planned++ }
		});
		await settle();
	};

	it('says the plan host action is on, and the form can be used', async () => {
		await show(view());
		expect(text(one('planner-action-on'))).toContain('flai serve starts the planner when you ask');
		expect(one('planner-action-off')).toBeNull();
		expect(one('planner-id')).not.toBeNull();
		expect((one('planner-id') as HTMLInputElement).disabled).toBe(false);
		expect((one('planner-plan') as HTMLButtonElement).disabled).toBe(false);
		expect(one('planner-form-off')).toBeNull();
	});

	it('says the action is off with the command that turns it on, and disables the form saying why', async () => {
		await show(view({ plan_enabled: false }));
		const off = one('planner-action-off')!;
		expect(text(off)).toContain('On the host, in the project, run flai serve enable plan');
		expect(off.querySelector('code')!.textContent).toBe('flai serve enable plan');
		expect((one('planner-id') as HTMLInputElement).disabled).toBe(true);
		expect((one('planner-plan') as HTMLButtonElement).disabled).toBe(true);
		const why = one('planner-form-off')!;
		expect(text(why)).toContain('The plan host action is off for this project');
		expect(why.querySelector('code')!.textContent).toBe('flai serve enable plan');
	});

	it('shows the run under way, linked to its item, with its stream asked with ?plan', async () => {
		api.mockResolvedValue(
			answer({
				story: '',
				item: 'E-0016',
				agent: 'planner',
				started: running.started,
				running: true,
				from: 0,
				next: 10,
				size: 10,
				entries: [{ kind: 'text', text: 'Reading the epic.' }]
			})
		);
		await show(view({ runs: [running, worked] }));
		const box = one('planner-current')!;
		expect(one('planner-current-item')!.querySelector('a')!.getAttribute('href')).toBe(
			'/items/E-0016'
		);
		expect(text(box)).toContain('planner (claude-code, opus)');
		expect(text(box)).toContain('2026-10-03 10:00 UTC');
		expect(text(one('planner-current-trigger'))).toBe('asked');
		expect(api).toHaveBeenCalledWith('/api/agent-stream/E-0016?plan');
		expect(text(box.querySelector('[data-testid="agent-stream"]'))).toContain('Reading the epic.');
		expect(one('planner-current-none')).toBeNull();
	});

	it('says no planner is running when none is under way', async () => {
		await show(view({ runs: [worked, failed] }));
		expect(text(one('planner-current-none'))).toBe('No planner is running.');
		expect(document.querySelector('[data-testid="agent-stream"]')).toBeNull();
		expect(api).not.toHaveBeenCalled();
	});

	it('shows the totals and the entries newest first, with items, duration, and cost', async () => {
		await show(view());
		expect(text(one('planner-accrued-cost'))).toBe('$3.50');
		expect(text(one('planner-accrued-time'))).toBe('1h30m');
		expect(text(one('planner-activities'))).toBe('2');
		expect(text(one('planner-last-run'))).toBe('2026-10-02 09:10 UTC');
		expect(all('planner-entry-summary').map(text)).toEqual(['enriched the story', 'drafted tasks']);
		const [newest, oldest] = all('planner-entry');
		expect(text(newest.querySelector('[data-testid="planner-entry-trigger"]'))).toBe(
			'trigger asked'
		);
		expect(
			[...newest.querySelectorAll('[data-testid="planner-entry-items"] a')].map((a) =>
				a.getAttribute('href')
			)
		).toEqual(['/items/S-0259', '/items/T-0838']);
		expect(text(newest.querySelector('[data-testid="planner-entry-items"]'))).toBe(
			'S-0259, T-0838'
		);
		expect(text(newest.querySelector('[data-testid="planner-entry-duration"]'))).toBe('20m');
		expect(text(newest.querySelector('[data-testid="planner-entry-cost"]'))).toBe(
			'$2.25 (estimated)'
		);
		expect(text(newest.querySelector('[data-testid="planner-entry-details"]'))).toBe(
			'trigger asked; items S-0259, T-0838; took 20m; cost $2.25 (estimated)'
		);
		expect(text(oldest.querySelector('[data-testid="planner-entry-details"]'))).toBe(
			'items none; took 10m; cost $1.25'
		);
		expect(oldest.querySelector('[data-testid="planner-entry-trigger"]')).toBeNull();
		expect(text(oldest.querySelector('[data-testid="planner-entry-items"]'))).toBe('none');
		expect(text(oldest.querySelector('[data-testid="planner-entry-cost"]'))).toBe('$1.25');
	});

	it('says when the planner has logged nothing', async () => {
		await show(view({ activity: { ...view().activity, entries: [] } }));
		expect(text(one('planner-entries-none'))).toBe('The planner has not logged an activity yet.');
		expect(one('planner-entries')).toBeNull();
	});

	it('lists the runs newest first with their outcome, cost, and why one failed', async () => {
		api.mockResolvedValue(answer({ entries: [], running: true, from: 0, next: 0, size: 0 }));
		const unstarted: PlanRun = {
			item: 'S-0002',
			agent: 'planner',
			command: 'claude',
			started: '2026-09-30T08:00:00Z',
			error: 'claude: not found'
		};
		await show(view({ runs: [failed, unstarted, worked, running] }));
		const rows = all('planner-run');
		expect(rows.map((r) => r.querySelector('a')!.textContent)).toEqual([
			'E-0016',
			'S-0259',
			'S-0001',
			'S-0002'
		]);
		expect(all('planner-run-outcome').map(text)).toEqual([
			'running',
			'worked',
			'failed',
			'could not start'
		]);
		// each ended run costs the entry logged within it; the run under way has logged none since it
		// started, and the run that could not start did no work
		expect(all('planner-run-cost').map(text)).toEqual(['—', '$2.25 (estimated)', '$1.25', '—']);
		expect(all('planner-run-why').map(text)).toEqual([
			'',
			'',
			'the harness exited 1',
			'claude: not found'
		]);
		expect(text(rows[1].querySelectorAll('td')[2])).toBe('2026-10-02 09:10 UTC');
		expect(text(rows[0].querySelectorAll('td')[2])).toBe('—');
	});

	it('asks flai serve to plan the ID given and says what it answered', async () => {
		await show(view());
		api.mockResolvedValue(answer(running));
		const input = one('planner-id') as HTMLInputElement;
		input.value = '  e-0016 ';
		input.dispatchEvent(new Event('input'));
		one('planner-plan')!.click();
		await settle();
		expect(api).toHaveBeenCalledWith('/api/items/E-0016/plan', { method: 'POST' });
		expect(text(one('planner-said'))).toBe(
			'planner started for E-0016 (pid 42); its output is in /home/op/.flai/serve/plan-E-0016.log on the host'
		);
		expect(one('planner-said')!.getAttribute('role')).toBe('status');
		expect(planned).toBe(1);
	});

	it('says why flai refused', async () => {
		await show(view());
		api.mockResolvedValue(answer({ error: 'S-0001 is done; flai plans only open items' }, 400));
		const input = one('planner-id') as HTMLInputElement;
		input.value = 'S-0001';
		input.dispatchEvent(new Event('input'));
		one('planner-plan')!.click();
		await settle();
		expect(text(one('planner-said'))).toBe('refused: S-0001 is done; flai plans only open items');
		expect(planned).toBe(1);
		expect((one('planner-plan') as HTMLButtonElement).disabled).toBe(false);
	});

	it('refuses an ID that is not an epic’s or a story’s without asking flai', async () => {
		await show(view());
		const input = one('planner-id') as HTMLInputElement;
		input.value = 'T-0838';
		input.dispatchEvent(new Event('input'));
		one('planner-plan')!.click();
		await settle();
		expect(api).not.toHaveBeenCalled();
		expect(text(one('planner-said'))).toBe(
			"T-0838 is not an epic's or a story's ID; give one such as E-0001 or S-0001"
		);
		expect(planned).toBe(0);
	});
});
