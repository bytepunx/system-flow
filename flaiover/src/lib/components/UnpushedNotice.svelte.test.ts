import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import UnpushedNotice from './UnpushedNotice.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, json: async () => body });
const pending = {
	upstream: 'origin/main',
	commits: 3,
	acceptances: ['S-0053'],
	tags: ['flai/v1.2.11', 'flaiover/v0.17.1'],
	command: 'flai push --pending'
};
const notice = () => document.querySelector<HTMLElement>('[data-testid="unpushed"]');

describe('UnpushedNotice', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('says what was accepted and not pushed, the tags, and the command, and offers no button', async () => {
		api.mockResolvedValue(answer({ unpushed: pending }));
		c = mount(UnpushedNotice, { target: document.body, props: {} });
		await settle();
		expect(api).toHaveBeenCalledWith('/api/unpushed');
		const text = notice()!.textContent!.replace(/\s+/g, ' ');
		expect(text).toContain('Accepted, not pushed: S-0053');
		expect(text).toContain('3 commits ahead of origin/main, tags flai/v1.2.11, flaiover/v0.17.1');
		expect(text).toContain('flai push --pending');
		expect(text).toContain('A push made from another clone is not seen here until someone fetches');
		expect(document.querySelectorAll('button')).toHaveLength(0);
	});

	it('shows nothing when nothing is pending', async () => {
		api.mockResolvedValue(answer({ unpushed: null }));
		c = mount(UnpushedNotice, { target: document.body, props: {} });
		await settle();
		expect(notice()).toBeNull();
	});

	it('asks again when refresh changes, and clears by itself once pushed', async () => {
		api.mockResolvedValueOnce(answer({ unpushed: pending }));
		const props = $state({ refresh: 0 });
		c = mount(UnpushedNotice, { target: document.body, props });
		await settle();
		expect(notice()).not.toBeNull();
		api.mockResolvedValue(answer({ unpushed: null }));
		props.refresh = 1;
		await settle();
		expect(api).toHaveBeenCalledTimes(2);
		expect(notice()).toBeNull();
	});

	it('asks again when the window regains focus', async () => {
		api.mockResolvedValueOnce(answer({ unpushed: pending }));
		c = mount(UnpushedNotice, { target: document.body, props: {} });
		await settle();
		api.mockResolvedValue(answer({ unpushed: null }));
		window.dispatchEvent(new Event('focus'));
		await settle();
		expect(notice()).toBeNull();
	});

	it('on an item page shows only that item’s acceptance', async () => {
		api.mockResolvedValue(answer({ unpushed: pending }));
		c = mount(UnpushedNotice, { target: document.body, props: { item: 'S-0001' } });
		await settle();
		expect(notice()).toBeNull();
		unmount(c);
		c = mount(UnpushedNotice, { target: document.body, props: { item: 'S-0053' } });
		await settle();
		expect(notice()).not.toBeNull();
	});

	it('shows a diverged clone as a problem on the board, and keeps quiet when the question fails', async () => {
		api.mockResolvedValue(
			answer({ unpushed: null, problem: 'main and origin/main have diverged' })
		);
		c = mount(UnpushedNotice, { target: document.body, props: {} });
		await settle();
		expect(document.querySelector('[data-testid="unpushed-problem"]')!.textContent).toContain(
			'have diverged'
		);
		unmount(c);
		document.body.innerHTML = '';
		api.mockRejectedValue(new Error('offline'));
		c = mount(UnpushedNotice, { target: document.body, props: {} });
		await settle();
		expect(document.body.textContent!.trim()).toBe('');
	});
});
