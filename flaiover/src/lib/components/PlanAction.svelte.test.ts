import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import PlanAction from './PlanAction.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const events = vi.hoisted(() => ({
	heard: [] as (() => void)[],
	agents: [] as ((id: string) => void)[]
}));
vi.mock('$lib/events', () => ({
	follow: (_kinds: string[], f: () => void) => {
		events.heard.push(f);
		return () => events.heard.splice(events.heard.indexOf(f), 1);
	},
	listen: (l: { agent?: (id: string) => void }) => {
		const f = l.agent ?? (() => {});
		events.agents.push(f);
		return () => events.agents.splice(events.agents.indexOf(f), 1);
	}
}));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const button = () => document.querySelector<HTMLButtonElement>('[data-testid="item-plan"]');
const running = () => document.querySelector('[data-testid="item-plan-running"]');
const run = {
	item: 'E-0016',
	agent: 'planner',
	harness: 'claude-code',
	command: 'claude',
	pid: 42,
	log: '/home/op/.flai/serve/plan-E-0016.log',
	started: '2026-10-03T10:00:00Z'
};

describe('PlanAction (S-0208)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let said: string[] = [];
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		events.heard = [];
		events.agents = [];
		said = [];
		document.body.innerHTML = '';
	});
	const show = async (status: unknown) => {
		api.mockResolvedValue(answer(status));
		c = mount(PlanAction, {
			target: document.body,
			props: { id: 'E-0016', onresult: (t: string) => said.push(t) }
		});
		await settle();
	};

	it('shows Plan only while the plan host action is on, and follows it as the project changes', async () => {
		await show({ plan_enabled: false, run: null });
		expect(api).toHaveBeenCalledWith('/api/items/E-0016/plan');
		expect(button()).toBeNull();

		api.mockResolvedValue(answer({ plan_enabled: true, run: null }));
		events.heard.forEach((f) => f());
		await settle();
		expect(button()!.textContent).toBe('Plan');
		expect(running()).toBeNull();

		api.mockResolvedValue(answer({ plan_enabled: false, run: null }));
		events.heard.forEach((f) => f());
		await settle();
		expect(button()).toBeNull();
	});

	it('says when the item’s planner has not ended', async () => {
		await show({ plan_enabled: true, run });
		// S-0329: started at 10:00 UTC, shown in the local zone, New York's under test
		expect(running()!.textContent).toContain('planner running since 2026-10-03 06:00 EDT');
		expect(running()!.textContent).not.toContain('UTC');
		unmount(c!);
		c = undefined;
		await show({ plan_enabled: true, run: { ...run, ended: '2026-10-03T10:05:00Z' } });
		expect(button()).not.toBeNull();
		expect(running()).toBeNull();
	});

	it('asks again when flai serve says the item’s planner started or ended, and only then', async () => {
		await show({ plan_enabled: true, run });
		expect(running()).not.toBeNull();
		api.mockResolvedValue(
			answer({ plan_enabled: true, run: { ...run, ended: '2026-10-03T10:05:00Z' } })
		);
		const asked = api.mock.calls.length;
		events.agents.forEach((f) => f('S-0001'));
		await settle();
		expect(api.mock.calls.length).toBe(asked);
		events.agents.forEach((f) => f('E-0016'));
		await settle();
		expect(running()).toBeNull();
	});

	it('starts the planner and says where its output goes', async () => {
		await show({ plan_enabled: true, run: null });
		api.mockImplementation(async (_url: string, init?: RequestInit) =>
			init?.method === 'POST' ? answer(run) : answer({ plan_enabled: true, run })
		);
		button()!.click();
		await settle();
		expect(api).toHaveBeenCalledWith('/api/items/E-0016/plan', { method: 'POST' });
		expect(said).toEqual([
			'planner started for E-0016 (pid 42); its output is in /home/op/.flai/serve/plan-E-0016.log on the host'
		]);
		expect(running()).not.toBeNull();
	});

	it('says why flai refused, as the page’s other actions do', async () => {
		await show({ plan_enabled: true, run: null });
		api.mockImplementation(async (_url: string, init?: RequestInit) =>
			init?.method === 'POST'
				? answer(
						{
							error:
								'the planner is already running for E-0016 (pid 42, started 2026-10-03T10:00:00Z); one item has one planner at a time'
						},
						400
					)
				: answer({ plan_enabled: true, run })
		);
		button()!.click();
		await settle();
		expect(said).toEqual([
			'refused: the planner is already running for E-0016 (pid 42, started 2026-10-03T10:00:00Z); one item has one planner at a time'
		]);
		expect(button()!.disabled).toBe(false);
	});
});
