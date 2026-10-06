// S-0154: a story page follows the project: its status, history, children, and body change on the
// page as they change in the project, without a reload.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/state', () => ({
	page: { params: { id: 'S-0154' }, url: new URL('http://localhost/items/S-0154') }
}));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
const enhance = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$lib/markdown', async (actual) => ({
	...(await actual<typeof import('$lib/markdown')>()),
	enhance
}));
const events = vi.hoisted(() => ({ heard: [] as { change?: (path: string) => void }[] }));
vi.mock('$lib/events', () => ({
	listen: (l: { change?: (path: string) => void }) => {
		events.heard.push(l);
		return () => events.heard.splice(events.heard.indexOf(l), 1);
	},
	follow: (_kinds: string[], f: () => void) => {
		const l = { change: () => f() };
		events.heard.push(l);
		return () => events.heard.splice(events.heard.indexOf(l), 1);
	}
}));
const changed = (path: string) => events.heard.forEach((l) => l.change?.(path));

import ItemPage from './+page.svelte';
import { inboxState } from '$lib/inbox.svelte';

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const story = {
	id: 'S-0154',
	type: 'story',
	nature: 'feature',
	title: 'Story pages receive live updates',
	status: 'in-progress',
	created: '2026-09-29T05:55:11Z',
	updated: '2026-09-29T07:06:54Z',
	transitions: [{ to: 'in-progress', at: '2026-09-29T07:06:54Z', by: 'agent-S-0154' }],
	path: 'wip/kanban/stories/S-0154-story-pages-receive-live-updates.md',
	archived: false,
	body: '# S-0154\n\n## Goal\n\nPages follow the project.\n'
};

describe('the item page (S-0154)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		enhance.mockClear();
		document.body.innerHTML = '';
	});

	const serve = (item: typeof story, children: unknown[] = []) =>
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item, children });
			if (url === '/api/board') return answer({ writable: false });
			if (url.startsWith('/api/threads')) return answer([]);
			return answer({ enabled: false });
		});
	const asked = (url: string) => api.mock.calls.filter(([u]) => u === url).length;
	const history = () => document.querySelector('aside ol')!.textContent!.replace(/\s+/g, ' ');

	it('shows a move, a new task, and a changed body when the project changes', async () => {
		serve(story);
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('h1')!.textContent).toContain('Story pages receive live updates');
		expect(enhance).toHaveBeenCalledTimes(1);

		// the agent writes a task and moves the story to review
		const review = {
			...story,
			status: 'review',
			transitions: [
				...story.transitions,
				{ to: 'review', at: '2026-09-29T08:00:00Z', by: 'agent-S-0154' }
			]
		};
		serve(review, [{ ...story, id: 'T-0553', type: 'task', title: 'A task', status: 'done' }]);
		changed('wip/kanban/stories/S-0154-story-pages-receive-live-updates.md');
		await settle();
		expect(document.querySelector('[data-testid="item-line"]')!.textContent).toContain('review');
		expect(history()).toContain('2026-09-29T08:00:00Z review by agent-S-0154');
		expect(document.body.textContent).toContain('T-0553');
		expect(document.body.textContent).toContain('Review this story');
		// the body did not change: it is not drawn again, and writable is not asked again
		expect(enhance).toHaveBeenCalledTimes(1);
		expect(asked('/api/board')).toBe(1);

		serve({ ...review, body: '# S-0154\n\n## Goal\n\nPages follow the project, live.\n' });
		changed('wip/kanban/stories/S-0154-story-pages-receive-live-updates.md');
		await settle();
		expect(document.querySelector('article')!.textContent).toContain(
			'Pages follow the project, live.'
		);
		expect(enhance).toHaveBeenCalledTimes(2);
	});

	it("says in the story's thread that its agent is working on the operator's reply", async () => {
		const reply = {
			id: 'TH-0040',
			title: 'Which animation?',
			anchor: { path: story.path, item: 'S-0154' },
			status: 'answered',
			participants: ['agent-S-0154', 'alex'],
			updated: '2026-09-29T08:05:00Z',
			entries: [
				{ at: '2026-09-29T08:00:00Z', author: 'agent-S-0154', text: 'Dots or a bar?' },
				{ at: '2026-09-29T08:05:00Z', author: 'alex', text: 'Dots.', operator: true }
			]
		};
		const run = {
			story: 'S-0154',
			command: 'claude',
			agent: 'agent-S-0154',
			started: '2026-09-29T07:06:00Z'
		};
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item: story, children: [] });
			if (url === '/api/board') return answer({ writable: false });
			if (url.startsWith('/api/threads')) return answer([reply]);
			if (url === '/api/host-agent')
				return answer({
					enabled: true,
					state: { command: 'claude', stories: { 'S-0154': { state: 'working', run } } }
				});
			return answer({});
		});
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="agent-working"]')!.textContent).toContain(
			'agent-S-0154 is working on your reply'
		);
	});

	it('dismisses the notice a refused edit leaves by its X (S-0151)', async () => {
		const ticked = {
			...story,
			body: '# S-0154\n\n## Goal\n\nPages follow the project.\n\n## Acceptance criteria\n\n- [ ] Live\n'
		};
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item: ticked, children: [] });
			if (url === '/api/board') return answer({ writable: true });
			if (url === '/api/items/S-0154/edit')
				return { ok: false, json: async () => ({ error: 'the story is being edited' }) };
			if (url.startsWith('/api/threads')) return answer([]);
			return answer({ enabled: false });
		});
		c = mount(ItemPage, { target: document.body });
		await settle();
		const box = document.querySelector<HTMLInputElement>('input[aria-label="criterion 1"]')!;
		box.click();
		await settle();
		const notice = () => document.querySelector('[data-testid="item-notice"]');
		expect(notice()!.textContent).toBe('refused: the story is being edited');
		document.querySelector<HTMLButtonElement>('[data-testid="dismiss"]')!.click();
		flushSync();
		expect(notice()).toBeNull();
	});

	describe('the New link at the top (S-0171)', () => {
		const show = async (item: Record<string, unknown>, writable = true) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			return document.querySelector<HTMLAnchorElement>('[data-testid="new-same-type"]');
		};

		it('names a story and opens the form for a story on its own, even under an epic (S-0192)', async () => {
			const link = await show({ ...story, parent: 'E-0013' });
			expect(link!.textContent).toBe('New story');
			expect(link!.getAttribute('href')).toBe('/new?type=story');
		});

		it('names a story with no epic and opens the form for a story alone', async () => {
			const link = await show(story);
			expect(link!.getAttribute('href')).toBe('/new?type=story');
		});

		it('names an epic and opens the form for an epic', async () => {
			const link = await show({ ...story, id: 'E-0013', type: 'epic', parent: undefined });
			expect(link!.textContent).toBe('New epic');
			expect(link!.getAttribute('href')).toBe('/new?type=epic');
		});

		it("is not offered on a task, which is the agent's to write, or when the dashboard cannot write", async () => {
			expect(await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0171' })).toBeNull();
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
			expect(await show(story, false)).toBeNull();
		});
	});

	describe('the Create story link on an epic (S-0191)', () => {
		const epic = { ...story, id: 'E-0013', type: 'epic', parent: undefined, status: 'in-progress' };
		const show = async (item: Record<string, unknown>, writable = true) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			return document.querySelector<HTMLAnchorElement>('[data-testid="new-child-story"]');
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};

		it('comes after New epic and opens the form for a story under the epic', async () => {
			const link = await show(epic);
			expect(link!.textContent).toBe('Create story');
			expect(link!.getAttribute('href')).toBe('/new?type=story&parent=E-0013');
			const links = [...link!.parentElement!.querySelectorAll('a')].map((a) => a.textContent);
			expect(links).toEqual(['New epic', 'Create story']);
		});

		it('is not offered on a closed or archived epic, whose stories the form would not take', async () => {
			expect(await show({ ...epic, status: 'done' })).toBeNull();
			again();
			expect(await show({ ...epic, status: 'cancelled' })).toBeNull();
			again();
			expect(await show({ ...epic, archived: true })).toBeNull();
		});

		it('is not offered on a story or a task, or when the dashboard cannot write', async () => {
			expect(await show({ ...story, parent: 'E-0013' })).toBeNull();
			again();
			expect(await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0171' })).toBeNull();
			again();
			expect(await show(epic, false)).toBeNull();
		});
	});

	describe('the New sibling link on a story (S-0192)', () => {
		const sibling = { ...story, parent: 'E-0013' };
		const epic = { ...story, id: 'E-0013', type: 'epic', parent: undefined, status: 'in-progress' };
		const show = async (
			item: Record<string, unknown>,
			{ writable = true, parent = (): unknown => answer({ item: epic, children: [] }) } = {}
		) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/items/E-0013') return parent();
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			return document.querySelector<HTMLAnchorElement>('[data-testid="new-sibling-story"]');
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};
		const withEpic = (e: Record<string, unknown>) => ({
			parent: () => answer({ item: { ...epic, ...e }, children: [] })
		});

		it('comes after New story and opens the form for a story under the same epic', async () => {
			const link = await show(sibling);
			expect(link!.textContent).toBe('New sibling');
			expect(link!.getAttribute('href')).toBe('/new?type=story&parent=E-0013');
			const links = [...link!.parentElement!.querySelectorAll('a')].map((a) => a.textContent);
			expect(links).toEqual(['New story', 'New sibling']);
		});

		it('asks after the epic once, not at every change in the project', async () => {
			await show(sibling);
			changed(story.path);
			await settle();
			expect(api.mock.calls.filter(([u]) => u === '/api/items/E-0013').length).toBe(1);
			expect(document.querySelector('[data-testid="new-sibling-story"]')).not.toBeNull();
		});

		it('is not offered on a story with no epic', async () => {
			expect(await show(story)).toBeNull();
			expect(api.mock.calls.some(([u]) => u === '/api/items/E-0013')).toBe(false);
		});

		it('is not offered under a closed or archived epic, which the form would not take', async () => {
			expect(await show(sibling, withEpic({ status: 'done' }))).toBeNull();
			again();
			expect(await show(sibling, withEpic({ status: 'cancelled' }))).toBeNull();
			again();
			expect(await show(sibling, withEpic({ archived: true }))).toBeNull();
		});

		it('is not offered when the epic cannot be read', async () => {
			const failed = () => ({ ok: false, json: async () => ({ error: 'no such item' }) });
			expect(await show(sibling, { parent: failed })).toBeNull();
			again();
			const offline = () => {
				throw new Error('offline');
			};
			expect(await show(sibling, { parent: offline })).toBeNull();
		});

		it('is not offered when the dashboard cannot write, or on an epic', async () => {
			expect(await show(sibling, { writable: false })).toBeNull();
			again();
			expect(await show(epic)).toBeNull();
		});
	});

	// S-0173: the inbox leads an open question here, so the story's page shows its open questions.
	it("shows the story's open questions from the inbox, and a task's page none", async () => {
		inboxState.data = {
			total: 1,
			counts: { thread: 0, question: 1, review: 0, blocked: 0, overlap: 0 },
			notes: [],
			entries: [
				{
					key: 'question:S-0154:abc',
					kind: 'question',
					title: 'Which port should it use?',
					href: '/items/S-0154?question=question%3AS-0154%3Aabc',
					item: 'S-0154'
				}
			]
		};
		serve(story);
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="open-questions"]')?.textContent).toContain(
			'Which port should it use?'
		);
		unmount(c);
		serve({ ...story, type: 'task' });
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="open-questions"]')).toBeNull();
		inboxState.data = null;
	});

	// S-0199: a story's page shows its cost of delay and forecast, read only, each only when set;
	// S-0204: the inputs and the planner's value each with who set them and when (ADR-0080)
	describe('its cost of delay and forecast', () => {
		const text = (testid: string) =>
			[...(document.querySelector(`[data-testid="${testid}"]`)?.children ?? [])].map((d) =>
				d.textContent!.trim()
			);
		const stale = () => document.querySelector('[data-testid="item-cost-of-delay-stale"]');
		const costOfDelay = (inputsAt: string) => ({
			inputs: {
				revenue_per_week: 1200,
				time_lost_per_cycle: '6h',
				by: 'alex',
				at: inputsAt
			},
			value: 1650,
			by: 'planner',
			at: '2026-10-03T09:20:00Z'
		});

		it("shows a story's cost of delay and forecast when set, and neither when not", async () => {
			serve({
				...story,
				cost_of_delay: costOfDelay('2026-10-03T09:00:00Z'),
				forecast: { delivery: '2026-10-07T17:00:00Z', by: 'planner', at: '2026-10-03T09:30:00Z' }
			} as typeof story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(text('item-cost-of-delay')).toEqual([
				'cost of delay: 1,650 per week, set by planner at 2026-10-03T09:20:00Z',
				'inputs: revenue 1,200 per week · 6h lost per cycle, set by alex at 2026-10-03T09:00:00Z'
			]);
			expect(stale()).toBeNull();
			expect(text('item-forecast')).toEqual([
				'forecast: delivery 2026-10-07T17:00:00Z',
				'set by planner at 2026-10-03T09:30:00Z'
			]);
			unmount(c);
			serve(story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(document.querySelector('[data-testid="item-cost-of-delay"]')).toBeNull();
			expect(document.querySelector('[data-testid="item-forecast"]')).toBeNull();
		});

		it('notes that the value is out of date when the inputs were set after it', async () => {
			serve({ ...story, cost_of_delay: costOfDelay('2026-10-03T10:00:00Z') } as typeof story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(stale()?.textContent?.trim()).toBe(
				"The inputs changed after the value was set, so the planner's value is out of date."
			);
			expect(text('item-cost-of-delay')[1]).toBe(
				'inputs: revenue 1,200 per week · 6h lost per cycle, set by alex at 2026-10-03T10:00:00Z'
			);
		});

		it('says an epic has no value yet when only its inputs are set, with no note', async () => {
			const { inputs } = costOfDelay('2026-10-03T10:00:00Z');
			serve({ ...story, type: 'epic', cost_of_delay: { inputs } } as typeof story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(text('item-cost-of-delay')).toEqual([
				'cost of delay: no value yet',
				'inputs: revenue 1,200 per week · 6h lost per cycle, set by alex at 2026-10-03T10:00:00Z'
			]);
			expect(stale()).toBeNull();
		});

		// S-0225, ADR-0083: what the planner expects it to cost, and what planning it cost, apart
		it("shows a planned story's expected cost and the planner's usage before any agent works it", async () => {
			const planned = {
				...story,
				status: 'backlog',
				forecast: { duration: '3h', by: 'planner', at: '2026-10-03T09:30:00Z' },
				usage: {
					source: 'sum',
					seconds: 0,
					models: [],
					strategic: [
						{
							kind: 'planner',
							seconds: 412,
							estimated: true,
							models: [
								{
									model: 'claude-opus-5-5',
									input: 12000,
									output: 0,
									cache_read: 800000,
									cache_write: 0,
									cost: 0.81
								}
							]
						}
					]
				}
			};
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154')
					return answer({
						item: planned,
						children: [],
						expected_cost: { cost: 12, from: 'forecast', estimated: true }
					});
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(
				document.querySelector('[data-testid="item-expected-cost"]')?.textContent?.trim()
			).toBe('expected cost: $12.00 (estimated, from the forecast)');
			const usage = document.querySelector('[data-testid="item-usage"]')!.textContent!;
			expect(usage).toContain('planner, strategic: 812.0K tokens · $0.810 (estimated) · 6m');
			// no agent has worked it, so there is no agents' line of nothing
			expect(usage).not.toContain('of agent work');
			unmount(c);
			serve(story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(document.querySelector('[data-testid="item-expected-cost"]')).toBeNull();
			expect(document.querySelector('[data-testid="item-usage"]')).toBeNull();
		});
	});

	// S-0201: a draft story says so beside its title and is finalized from its page without a reload
	describe('a draft story (S-0201)', () => {
		const draftStory = { ...story, status: 'backlog', draft: true };
		const epic = { ...story, id: 'E-0016', type: 'epic', draft: undefined };
		const indicator = () => document.querySelector('[data-testid="item-draft"]');
		const button = () => document.querySelector<HTMLButtonElement>('[data-testid="item-finalize"]');
		const show = async (
			item: Record<string, unknown>,
			children: unknown[] = [],
			writable = true
		) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};

		it('shows [Draft] beside the title and a Finalize button', async () => {
			await show(draftStory);
			expect(indicator()!.textContent).toBe('[Draft]');
			expect(indicator()!.closest('h1')).not.toBeNull();
			expect(button()!.textContent).toBe('Finalize');
		});

		it('shows neither on a story that is not a draft, nor on an epic', async () => {
			await show(story);
			expect(indicator()).toBeNull();
			expect(button()).toBeNull();
			again();
			await show({ ...epic, draft: true });
			expect(indicator()).toBeNull();
			expect(button()).toBeNull();
		});

		it('offers Finalize only where the dashboard may write, and not on an archived story', async () => {
			await show(draftStory, [], false);
			expect(indicator()).not.toBeNull();
			expect(button()).toBeNull();
			again();
			await show({ ...draftStory, archived: true });
			expect(button()).toBeNull();
		});

		it('finalizes through the API and, once reloaded, shows neither', async () => {
			let current: Record<string, unknown> = draftStory;
			api.mockImplementation(async (url: string, init?: RequestInit) => {
				if (url === '/api/items/S-0154') return answer({ item: current, children: [] });
				if (url === '/api/items/S-0154/finalize' && init?.method === 'POST') {
					current = { ...draftStory, draft: false };
					return answer({ id: 'S-0154', draft: false, warnings: [] });
				}
				if (url === '/api/board') return answer({ writable: true });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			button()!.click();
			await settle();
			expect(asked('/api/items/S-0154/finalize')).toBe(1);
			expect(indicator()).toBeNull();
			expect(button()).toBeNull();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toBe('done');
		});

		it('says why when flai refuses, and keeps the draft', async () => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item: draftStory, children: [] });
				if (url === '/api/items/S-0154/finalize')
					return { ok: false, json: async () => ({ error: 'S-0154 is not a draft' }) };
				if (url === '/api/board') return answer({ writable: true });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			button()!.click();
			await settle();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toBe(
				'refused: S-0154 is not a draft'
			);
			expect(indicator()).not.toBeNull();
		});

		it("marks an epic's draft story among its children", async () => {
			await show(epic, [
				{ ...story, id: 'S-0201', title: 'Drafted', draft: true },
				{ ...story, id: 'S-0200', title: 'Finalized' }
			]);
			const marks = [...document.querySelectorAll('[data-testid="child-draft"]')];
			expect(marks).toHaveLength(1);
			expect(marks[0].textContent).toBe('Draft');
			expect(marks[0].closest('li')!.textContent).toContain('S-0201');
		});
	});

	// S-0208: an open epic's or story's page offers Plan while the plan host action is on
	describe('the Plan action (S-0208)', () => {
		const epic = { ...story, id: 'E-0016', type: 'epic' };
		const button = () => document.querySelector<HTMLButtonElement>('[data-testid="item-plan"]');
		const show = async (item: Record<string, unknown>, writable = true, enabled = true) => {
			api.mockImplementation(async (url: string, init?: RequestInit) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				if (url === `/api/items/${item.id}/plan` && init?.method === 'POST')
					return {
						ok: false,
						json: async () => ({
							error: `the planner is already running for ${item.id} (pid 42, started 2026-10-03T10:00:00Z); one item has one planner at a time`
						})
					};
				if (url === `/api/items/${item.id}/plan`)
					return answer({ plan_enabled: enabled, run: null });
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};

		it("is among an epic's and a story's actions while the host action is on", async () => {
			await show(epic);
			expect(button()!.textContent).toBe('Plan');
			again();
			await show(story);
			expect(button()!.textContent).toBe('Plan');
			again();
			await show(story, true, false);
			expect(button()).toBeNull();
		});

		it('is not offered on a task, a closed or archived item, or when the dashboard cannot write', async () => {
			for (const item of [
				{ ...story, id: 'T-0788', type: 'task' },
				{ ...story, status: 'done' },
				{ ...epic, status: 'cancelled' },
				{ ...story, archived: true }
			]) {
				await show(item);
				expect(button()).toBeNull();
				again();
			}
			await show(story, false);
			expect(button()).toBeNull();
		});

		it("says in the page's notice why flai refused", async () => {
			await show(epic);
			button()!.click();
			await settle();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toContain(
				'refused: the planner is already running for E-0016'
			);
		});
	});

	// S-0176: a story's page shows its task plan, and follows it as tasks move
	it("shows a story's task plan when flai sends one, and follows it as the tasks move", async () => {
		const plan = {
			tasks: [
				{ id: 'T-0001', state: 'waiting', after: ['T-0002'], waiting_for: ['T-0002'] },
				{ id: 'T-0002', state: 'in-progress' },
				{ id: 'T-0003', state: 'ready' }
			],
			layers: [['T-0002', 'T-0003'], ['T-0001']]
		};
		const children = [{ ...story, id: 'T-0002', type: 'task', title: 'The flai side' }];
		const show = (p?: unknown) =>
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item: story, children, plan: p });
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
		const shown = () => document.querySelector('[data-testid="task-plan"]');
		const tasks = () =>
			[...document.querySelectorAll('[data-testid="plan-task"]')].map((t) =>
				t.textContent!.replace(/\s+/g, ' ').trim()
			);
		show(plan);
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(shown()!.closest('aside')).not.toBeNull();
		expect(tasks()).toEqual([
			'T-0001 waiting for T-0002',
			'T-0002 in progress',
			'T-0003 ready to start'
		]);
		expect(document.querySelectorAll('[data-testid="plan-layer"]')).toHaveLength(2);
		expect(shown()!.querySelector('a[href="/items/T-0002"]')!.getAttribute('title')).toBe(
			'The flai side'
		);

		// T-0002 is done: T-0001 can start
		show({
			...plan,
			tasks: [
				{ id: 'T-0001', state: 'ready', after: ['T-0002'] },
				{ id: 'T-0002', state: 'done' },
				{ id: 'T-0003', state: 'ready' }
			]
		});
		changed('wip/kanban/tasks/T-0002-the-flai-side.md');
		await settle();
		expect(tasks()).toEqual(['T-0001 ready to start', 'T-0002 done', 'T-0003 ready to start']);

		// a story without tasks has no plan, and shows none
		show(undefined);
		changed('wip/kanban/tasks/T-0002-the-flai-side.md');
		await settle();
		expect(shown()).toBeNull();
	});

	// ADR-0093: a done story says who accepted it and on which day, from its last transition to
	// done; the orchestrator's acceptance links to the evidence it wrote in the story's Notes.
	describe('who accepted a done story (ADR-0093)', () => {
		const done = (by: string, notes = '') => ({
			...story,
			status: 'done',
			transitions: [
				...story.transitions,
				{ to: 'review', at: '2026-10-05T22:00:00Z', by: 'agent-S-0154' },
				{ to: 'done', at: '2026-10-06T09:30:00Z', by }
			],
			body: story.body + notes
		});
		const line = () =>
			document
				.querySelector('[data-testid="item-accepted"]')
				?.textContent!.replace(/\s+/g, ' ')
				.trim();
		const link = () =>
			document.querySelector<HTMLAnchorElement>('[data-testid="item-accepted-evidence"]');
		const notes =
			'\n## Notes\n\n### Accepted by the orchestrator\n\n- Verified: 0123456789abcdef0123456789abcdef01234567\n- At: 2026-10-06T09:30:00Z\n\nVerdict: pass, tests and lint clean\n- 1: `flai/cmd/accept.go`\n';

		it('names the orchestrator and links to its evidence in the Notes', async () => {
			serve(done('orchestrator', notes));
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(line()).toBe(
				'Accepted by orchestrator on 2026-10-06, with its evidence in the Notes.'
			);
			expect(link()!.getAttribute('href')).toBe('#accepted-by-the-orchestrator');
			// the link lands on the heading the body renders
			expect(document.getElementById('accepted-by-the-orchestrator')!.textContent).toContain(
				'Accepted by the orchestrator'
			);
		});

		it('names the operator with no evidence link, and says nothing of a story that is not done', async () => {
			serve(done('alex'));
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(line()).toBe('Accepted by alex on 2026-10-06.');
			expect(link()).toBeNull();
			unmount(c);
			document.body.innerHTML = '';
			serve(story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(line()).toBeUndefined();
		});
	});

	it('stops following the project once it is left', async () => {
		serve(story);
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(events.heard.length).toBeGreaterThan(0);
		unmount(c);
		c = undefined;
		expect(events.heard).toHaveLength(0);
	});
});
