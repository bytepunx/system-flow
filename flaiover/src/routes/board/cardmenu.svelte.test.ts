// S-0202: a right click, the context-menu key, the card's menu button, or a long press opens a card's
// menu, which offers what the item's page does (Open, Finalize on a draft, the agent, Plan since
// S-0263, Block or Unblock, Cancel) and then the lane's menu, and makes the page's writes.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { resetForTests } from '$lib/project.svelte';

const api = vi.fn();
const goto = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({ resolve: (route: string) => route }));
vi.mock('$app/navigation', () => ({ goto: (...args: unknown[]) => goto(...args) }));

import BoardPage from './+page.svelte';
import { PRESS_MS } from './longpress';

class QuietEventSource {
	addEventListener() {}
	close() {}
}

type Card = {
	id: string;
	type: string;
	title: string;
	nature: string;
	status: string;
	blocked: boolean;
	draft?: boolean;
	age_seconds: number;
};
const card = (id: string, status: string, more: Partial<Card> = {}): Card => ({
	id,
	type: 'story',
	title: `${id} title`,
	nature: 'feature',
	status,
	blocked: false,
	age_seconds: 60,
	...more
});
const answer = (body: unknown, ok = true) => ({
	ok,
	statusText: ok ? 'OK' : 'Bad',
	json: async () => body
});
const settle = async () => {
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

describe('a card’s menu (S-0202)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let posts: { url: string; body: Record<string, unknown> }[];
	let columns: Record<string, Card[]>;
	let agents: unknown;
	let posted: ReturnType<typeof answer>;

	const open = async (writable = true) => {
		api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
			if (init?.method === 'POST') {
				posts.push({ url, body: JSON.parse(init.body ?? '{}') });
				return posted;
			}
			if (url === '/api/board')
				return answer({ wip_limits: { ready: 5 }, order: [], writable, columns });
			if (url === '/api/publish') return answer({ plans: [], push_enabled: false });
			if (url === '/api/host-agent') return answer(agents);
			return answer({});
		});
		c = mount(BoardPage, { target: document.body });
		await settle();
	};
	const wrapper = (id: string) => document.querySelector<HTMLElement>(`[data-card="${id}"]`)!;
	const link = (id: string) => document.querySelector<HTMLAnchorElement>(`a[data-id="${id}"]`)!;
	const menu = () => document.querySelector<HTMLElement>('[data-testid="card-menu"]');
	const rightClick = (el: Element, init: MouseEventInit = {}) => {
		const e = new MouseEvent('contextmenu', {
			bubbles: true,
			cancelable: true,
			button: 2,
			clientX: 40,
			clientY: 50,
			...init
		});
		el.dispatchEvent(e);
		flushSync();
		return e;
	};
	const offered = () =>
		[...(menu()?.querySelectorAll<HTMLButtonElement>('[role="menuitem"]') ?? [])].map(
			(b) => b.textContent
		);
	const pick = async (action: string) => {
		menu()!.querySelector<HTMLButtonElement>(`[data-action="${action}"]`)!.click();
		await settle();
	};
	const boards = () => api.mock.calls.filter(([u]) => u === '/api/board').length;
	const notice = () => document.querySelector('[data-testid="board-notice"]')?.textContent ?? '';

	beforeEach(() => {
		resetForTests();
		vi.stubGlobal('EventSource', QuietEventSource);
		posts = [];
		agents = { enabled: false };
		posted = answer({});
		columns = {
			backlog: [card('S-0001', 'backlog', { draft: true }), card('S-0002', 'backlog')],
			ready: [card('S-0003', 'ready')],
			review: [card('S-0005', 'review')]
		};
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		goto.mockReset();
		vi.useRealTimers();
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
	});

	it('opens on a right click with the page’s actions, Finalize on a draft, then the lane’s', async () => {
		await open();
		const e = rightClick(link('S-0001'));
		expect(e.defaultPrevented).toBe(true);
		expect(document.querySelector('[data-testid="lane-menu"]')).toBeNull();
		expect(menu()!.getAttribute('aria-label')).toBe('S-0001 actions');
		expect(menu()!.style.left).toBe('40px');
		expect(offered()).toEqual([
			'Open',
			'Finalize',
			'Block…',
			'Cancel…',
			'Create item here',
			'Move stories forward to ready…'
		]);
		expect(menu()!.querySelector('[role="group"][aria-label="backlog"]')).not.toBeNull();
		await settle();
		expect(document.activeElement?.getAttribute('data-action')).toBe('open');
	});

	it('offers no Finalize on a finalized story', async () => {
		await open();
		rightClick(link('S-0002'));
		expect(offered()).toEqual([
			'Open',
			'Block…',
			'Cancel…',
			'Create item here',
			'Move stories forward to ready…'
		]);
	});

	it('finalizes a draft and shows the card the board then reads', async () => {
		await open();
		expect(wrapper('S-0001').querySelector('[data-testid="draft"]')).not.toBeNull();
		rightClick(link('S-0001'));
		const before = boards();
		columns.backlog[0] = { ...columns.backlog[0], draft: false };
		await pick('finalize');
		expect(posts).toEqual([{ url: '/api/items/S-0001/finalize', body: {} }]);
		expect(boards()).toBe(before + 1);
		expect(menu()).toBeNull();
		expect(notice()).toContain('S-0001 finalized');
		expect(wrapper('S-0001').querySelector('[data-testid="draft"]')).toBeNull();
	});

	it('opens from the card’s menu button and from the context-menu key, below the card', async () => {
		await open();
		const button = wrapper('S-0002').querySelector<HTMLButtonElement>(
			'[data-testid="card-menu-button"]'
		)!;
		expect(button.getAttribute('aria-label')).toBe('Actions for S-0002');
		expect(button.getAttribute('aria-haspopup')).toBe('menu');
		expect(button.getAttribute('aria-expanded')).toBe('false');
		button.click();
		flushSync();
		expect(menu()!.getAttribute('aria-label')).toBe('S-0002 actions');
		expect(button.getAttribute('aria-expanded')).toBe('true');

		menu()!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true }));
		await settle();
		expect(menu()).toBeNull();
		// Escape gives focus back to the card
		expect(document.activeElement).toBe(link('S-0002'));

		rightClick(link('S-0003'), { clientX: 0, clientY: 0, button: 0 });
		expect(menu()!.getAttribute('aria-label')).toBe('S-0003 actions');
	});

	it('closes when focus leaves it, and on a choice', async () => {
		await open();
		rightClick(link('S-0002'));
		await settle();
		document.querySelector<HTMLElement>('[data-testid="new-item-link"]')!.focus();
		flushSync();
		expect(menu()).toBeNull();

		rightClick(link('S-0005'));
		await pick('open');
		expect(menu()).toBeNull();
		// a story in review opens on its review, as its card does
		expect(goto).toHaveBeenCalledWith('/review/[id]');
		rightClick(link('S-0002'));
		await pick('open');
		expect(goto).toHaveBeenLastCalledWith('/items/[id]');
	});

	it('is not offered on a read-only board', async () => {
		await open(false);
		expect(rightClick(link('S-0001')).defaultPrevented).toBe(false);
		expect(menu()).toBeNull();
		expect(document.querySelector('[data-testid="card-menu-button"]')).toBeNull();
	});

	it('blocks with the reason asked for, and unblocks', async () => {
		await open();
		vi.stubGlobal(
			'prompt',
			vi.fn(() => null)
		);
		rightClick(link('S-0002'));
		await pick('block');
		expect(posts).toEqual([]);

		vi.stubGlobal(
			'prompt',
			vi.fn(() => 'waiting on X')
		);
		columns.backlog[1] = { ...columns.backlog[1], blocked: true };
		rightClick(link('S-0002'));
		await pick('block');
		expect(posts).toEqual([{ url: '/api/items/S-0002/block', body: { reason: 'waiting on X' } }]);
		expect(notice()).toContain('S-0002 blocked');
		expect(wrapper('S-0002').textContent).toContain('BLOCKED');

		rightClick(link('S-0002'));
		expect(offered()).toContain('Unblock');
		expect(offered()).not.toContain('Block…');
		await pick('unblock');
		expect(posts.at(-1)).toEqual({ url: '/api/items/S-0002/unblock', body: {} });
	});

	it('starts an agent on a ready story when the host can', async () => {
		agents = { enabled: true, state: { command: 'claude', stories: {} } };
		await open();
		rightClick(link('S-0003'));
		expect(offered()).toContain('Start agent');
		const asked = api.mock.calls.filter(([u]) => u === '/api/host-agent').length;
		await pick('agent');
		expect(posts).toEqual([{ url: '/api/items/S-0003/agent', body: { action: 'start' } }]);
		expect(notice()).toContain('S-0003: agent started');
		expect(api.mock.calls.filter(([u]) => u === '/api/host-agent').length).toBe(asked + 1);
	});

	it('starts the planner when the plan host action is on, and says where its output goes (S-0263)', async () => {
		agents = { enabled: false, plan_enabled: true };
		await open();
		rightClick(link('S-0002'));
		expect(offered()).toEqual([
			'Open',
			'Plan',
			'Block…',
			'Cancel…',
			'Create item here',
			'Move stories forward to ready…'
		]);
		const asked = api.mock.calls.filter(([u]) => u === '/api/host-agent').length;
		posted = answer({ item: 'S-0002', pid: 42, log: '/home/op/.flai/serve/plan-S-0002.log' });
		await pick('plan');
		expect(posts).toEqual([{ url: '/api/items/S-0002/plan', body: {} }]);
		expect(notice()).toContain(
			'planner started for S-0002 (pid 42); its output is in /home/op/.flai/serve/plan-S-0002.log on the host'
		);
		expect(api.mock.calls.filter(([u]) => u === '/api/host-agent').length).toBe(asked + 1);
	});

	it('says why flai refused to start the planner (S-0263)', async () => {
		agents = { enabled: false, plan_enabled: true };
		await open();
		rightClick(link('S-0002'));
		posted = answer({ error: 'the planner is already running for S-0002 (pid 42)' }, false);
		await pick('plan');
		expect(notice()).toContain('the planner is already running for S-0002 (pid 42)');
		expect(notice()).not.toContain('planner started');
	});

	it('cancels through the confirmation the board already asks', async () => {
		await open();
		rightClick(link('S-0002'));
		await pick('cancel');
		expect(menu()).toBeNull();
		expect(document.getElementById('cancel-title')!.textContent).toContain('S-0002');
	});

	it('opens on a long press and swallows the click that ends it', async () => {
		await open();
		vi.useFakeTimers();
		const touch = (type: string, x = 10, y = 10) => {
			const e = new MouseEvent(type, { bubbles: true, cancelable: true, clientX: x, clientY: y });
			Object.defineProperty(e, 'pointerType', { value: 'touch' });
			link('S-0002').dispatchEvent(e);
		};
		// a finger that moves is scrolling, not pressing
		touch('pointerdown');
		touch('pointermove', 30, 10);
		vi.advanceTimersByTime(PRESS_MS + 100);
		flushSync();
		expect(menu()).toBeNull();

		touch('pointerdown');
		vi.advanceTimersByTime(PRESS_MS + 100);
		flushSync();
		expect(menu()!.getAttribute('aria-label')).toBe('S-0002 actions');
		touch('pointerup');
		const reached = vi.fn((e: Event) => e.preventDefault());
		link('S-0002').addEventListener('click', reached);
		const click = new MouseEvent('click', { bubbles: true, cancelable: true });
		link('S-0002').dispatchEvent(click);
		expect(click.defaultPrevented).toBe(true);
		expect(reached).not.toHaveBeenCalled();

		// the next tap is a tap
		touch('pointerdown');
		touch('pointerup');
		link('S-0002').dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }));
		expect(reached).toHaveBeenCalledTimes(1);
	});
});
