// S-0173, TH-0041: the threads page shows every open thread and opens on the one a link names, so a
// thread on no item is answered there.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/state', () => ({
	page: { url: new URL('http://localhost/threads?thread=TH-0017') }
}));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
vi.mock('$lib/events', () => ({ follow: () => () => {} }));

import ThreadsPage from './+page.svelte';

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const settle = async () => {
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const thread = (id: string, anchor: Record<string, string>) => ({
	id,
	title: `about ${id}`,
	anchor,
	status: 'open',
	participants: ['agent'],
	updated: '2026-09-29T05:23:48Z',
	entries: [{ at: '2026-09-29T05:23:48Z', author: 'agent', text: 'may I?' }]
});

describe('the threads page (S-0173)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('opens on the thread the link names, answerable when the dashboard can write', async () => {
		Element.prototype.scrollIntoView = () => {};
		api.mockImplementation(async (url: string) => {
			if (url === '/api/board') return answer({ writable: true });
			if (url === '/api/threads')
				return answer([
					thread('TH-0026', { path: 'wip/kanban/stories/S-0139-a.md', item: 'S-0139' }),
					thread('TH-0017', { path: 'design/issues/I-0027-a.md' })
				]);
			return answer({});
		});
		c = mount(ThreadsPage, { target: document.body });
		await settle();
		expect(document.querySelector('h1')!.textContent).toBe('Threads');
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0017');
		expect(document.querySelector('article form input[placeholder="reply"]')).not.toBeNull();
	});
});
