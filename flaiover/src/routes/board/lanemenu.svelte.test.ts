// S-0167: a right click on a lane opens its menu, which offers only what the lane allows by flai's
// rules (ADR-0055): create an item starting there, move stories forward from backlog or back from
// ready, in-progress, review, and cancelled, and change the WIP limit of ready, in-progress, review.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { resetForTests } from '$lib/project.svelte';

const api = vi.fn();
const goto = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({ resolve: (route: string) => route }));
vi.mock('$app/navigation', () => ({ goto: (...args: unknown[]) => goto(...args) }));

import BoardPage from './+page.svelte';

class QuietEventSource {
	addEventListener() {}
	close() {}
}

const card = (id: string, status: string, type = 'story') => ({
	id,
	type,
	title: `${id} title`,
	nature: 'feature',
	status,
	blocked: false,
	age_seconds: 60
});
const board = (writable = true) => ({
	wip_limits: { ready: 5, 'in-progress': 2, review: 3 },
	order: [],
	writable,
	columns: {
		backlog: [
			card('S-0001', 'backlog'),
			card('S-0002', 'backlog'),
			card('E-0001', 'backlog', 'epic')
		],
		ready: [card('S-0003', 'ready')],
		'in-progress': [card('S-0004', 'in-progress')],
		review: [card('S-0005', 'review')],
		done: [card('S-0006', 'done')],
		cancelled: [card('S-0007', 'cancelled')]
	}
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

describe('the lane menu (S-0167)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let posts: { url: string; body: Record<string, unknown> }[];
	let refuse: (url: string) => string | null;

	const open = async (writable = true) => {
		api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
			if (init?.method === 'POST') {
				const body = JSON.parse(init.body ?? '{}');
				posts.push({ url, body });
				const why = refuse(url);
				if (why) return answer({ error: why }, false);
				if (url === '/api/board/limit') return answer({ column: body.column, limit: body.limit });
				return answer({ id: url.split('/')[3], status: body.to, warnings: [] });
			}
			if (url === '/api/board') return answer(board(writable));
			if (url === '/api/publish') return answer({ plans: [], push_enabled: false });
			return answer({ enabled: false });
		});
		c = mount(BoardPage, { target: document.body });
		await settle();
	};
	const lane = (name: string) => document.querySelector<HTMLElement>(`[data-lane="${name}"]`)!;
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
	const entries = () =>
		[
			...document.querySelectorAll<HTMLButtonElement>('[data-testid="lane-menu"] [role="menuitem"]')
		].map((b) => b.dataset.action);
	const pick = async (action: string) => {
		document
			.querySelector<HTMLButtonElement>(`[role="menuitem"][data-action="${action}"]`)!
			.click();
		await settle();
	};
	const ticked = () =>
		[
			...document.querySelectorAll<HTMLInputElement>('[data-testid="lane-move"] ul input:checked')
		].map((i) => i.value);

	beforeEach(() => {
		resetForTests();
		vi.stubGlobal('EventSource', QuietEventSource);
		posts = [];
		refuse = () => null;
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		goto.mockReset();
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
	});

	it('offers in each lane only what the lane allows', async () => {
		await open();
		const offered: Record<string, (string | undefined)[]> = {};
		for (const name of ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled']) {
			const e = rightClick(lane(name));
			expect(e.defaultPrevented).toBe(true);
			offered[name] = entries();
			document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
			document
				.querySelector('[data-testid="lane-menu"]')!
				.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true }));
			flushSync();
			expect(document.querySelector('[data-testid="lane-menu"]')).toBeNull();
		}
		expect(offered).toEqual({
			backlog: ['create', 'forward'],
			ready: ['create', 'back', 'limit'],
			'in-progress': ['create', 'back', 'limit'],
			review: ['create', 'back', 'limit'],
			done: ['create'],
			cancelled: ['create', 'back']
		});
	});

	it('is not offered on a read-only board, and Shift keeps the browser menu', async () => {
		await open(false);
		expect(rightClick(lane('ready')).defaultPrevented).toBe(false);
		expect(document.querySelector('[data-testid="lane-menu"]')).toBeNull();
		unmount(c!);
		document.body.innerHTML = '';
		await open();
		expect(rightClick(lane('ready'), { shiftKey: true }).defaultPrevented).toBe(false);
		expect(document.querySelector('[data-testid="lane-menu"]')).toBeNull();
	});

	it('creates an item starting in the lane it was opened on', async () => {
		await open();
		rightClick(lane('in-progress'));
		await pick('create');
		expect(goto).toHaveBeenCalledWith('/new?status=in-progress');
		rightClick(lane('review'));
		await pick('create');
		// the create screen starts anything but backlog, ready, and in-progress in backlog
		expect(goto).toHaveBeenLastCalledWith('/new?status=review');
	});

	it('moves the ticked backlog stories forward to ready, the card clicked ticked first', async () => {
		await open();
		rightClick(document.querySelector('[data-card="S-0002"]')!);
		await pick('forward');
		expect(document.querySelector('[data-testid="lane-move"] h2')!.textContent).toContain(
			'from backlog to ready'
		);
		// stories only: the epic in backlog is not offered
		expect(
			[...document.querySelectorAll<HTMLInputElement>('[data-testid="lane-move"] ul input')].map(
				(i) => i.value
			)
		).toEqual(['S-0001', 'S-0002']);
		expect(ticked()).toEqual(['S-0002']);
		document.querySelector<HTMLInputElement>('[data-testid="lane-move-all"]')!.click();
		flushSync();
		expect(ticked()).toEqual(['S-0001', 'S-0002']);
		document.querySelector<HTMLButtonElement>('[data-testid="lane-move-confirm"]')!.click();
		await settle();
		expect(posts).toEqual([
			{ url: '/api/items/S-0001/move', body: { to: 'ready' } },
			{ url: '/api/items/S-0002/move', body: { to: 'ready' } }
		]);
		expect(document.querySelector('[data-testid="lane-move"]')).toBeNull();
		expect(document.querySelector('[data-testid="board-notice"]')!.textContent).toContain(
			'S-0001, S-0002 → ready'
		);
	});

	it('moves stories back a column, asks why only from review, and names each refusal', async () => {
		await open();
		rightClick(lane('cancelled'));
		await pick('back');
		expect(ticked()).toEqual([]);
		const confirm = () =>
			document.querySelector<HTMLButtonElement>('[data-testid="lane-move-confirm"]')!;
		expect(confirm().disabled).toBe(true);
		expect(document.querySelector('[data-testid="lane-move-reason"]')).toBeNull();
		document.querySelector<HTMLInputElement>('[data-testid="lane-move"] ul input')!.click();
		flushSync();
		confirm().click();
		await settle();
		expect(posts.at(-1)).toEqual({ url: '/api/items/S-0007/move', body: { to: 'backlog' } });

		rightClick(document.querySelector('[data-card="S-0005"]')!);
		await pick('back');
		expect(ticked()).toEqual(['S-0005']);
		expect(confirm().disabled).toBe(true); // review to in-progress needs a reason
		const reason = document.querySelector<HTMLTextAreaElement>('[data-testid="lane-move-reason"]')!;
		reason.value = 'tests missing';
		reason.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		refuse = (url) => (url.includes('S-0005') ? 'rule: something' : null);
		confirm().click();
		await settle();
		expect(posts.at(-1)).toEqual({
			url: '/api/items/S-0005/move',
			body: { to: 'in-progress', reason: 'tests missing' }
		});
		const notice = document.querySelector('[data-testid="board-notice"]')!;
		expect(notice.textContent).toContain('Refused: S-0005: rule: something');
		expect(notice.parentElement!.getAttribute('role')).toBe('alert');
	});

	it("changes a lane's WIP limit, 0 for none", async () => {
		await open();
		rightClick(lane('in-progress'));
		await pick('limit');
		const input = document.querySelector<HTMLInputElement>('[data-testid="wip-limit-value"]')!;
		expect(input.value).toBe('2');
		input.value = '4';
		input.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		document.querySelector<HTMLButtonElement>('[data-testid="wip-limit-confirm"]')!.click();
		await settle();
		expect(posts).toEqual([{ url: '/api/board/limit', body: { column: 'in-progress', limit: 4 } }]);
		expect(document.querySelector('[data-testid="wip-limit"]')).toBeNull();
		expect(document.querySelector('[data-testid="board-notice"]')!.textContent).toContain(
			'WIP limit for in-progress: 4'
		);

		rightClick(lane('review'));
		await pick('limit');
		const again = document.querySelector<HTMLInputElement>('[data-testid="wip-limit-value"]')!;
		again.value = '-1';
		again.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		expect(
			document.querySelector<HTMLButtonElement>('[data-testid="wip-limit-confirm"]')!.disabled
		).toBe(true);
	});
});
