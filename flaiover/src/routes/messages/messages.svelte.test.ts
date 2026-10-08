// S-0336: the messages page lists every open conversation between stories, the closed ones on
// request, and says when there is none.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
vi.mock('$lib/events', () => ({ follow: () => () => {} }));

import MessagesPage from './+page.svelte';

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const settle = async () => {
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const conversation = (id: string, closed = false) => ({
	id,
	title: `about ${id}`,
	from: 'S-0001',
	to: 'S-0002',
	about: [],
	status: closed ? 'closed' : 'open',
	closed,
	closed_reason: closed ? 'S-0001 was accepted' : '',
	awaiting: closed ? '' : 'S-0002',
	participants: ['agent-S-0001'],
	created: '2026-10-08T06:00:00Z',
	updated: '2026-10-08T07:00:00Z',
	path: `wip/messages/${id}-x.md`,
	entries: [{ at: '2026-10-08T07:00:00Z', author: 'agent-S-0001', story: 'S-0001', text: 'hi' }]
});
const listed = () =>
	[...document.querySelectorAll('[data-conversation]')].map(
		(m) => (m as HTMLElement).dataset.conversation
	);

describe('the messages page (S-0336)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('lists every open conversation, and the closed ones when asked', async () => {
		api.mockImplementation(async (url: string) =>
			answer(
				url === '/api/messages?all=1'
					? [conversation('MS-0001'), conversation('MS-0002', true)]
					: [conversation('MS-0001')]
			)
		);
		c = mount(MessagesPage, { target: document.body });
		await settle();
		expect(document.querySelector('h1')!.textContent).toBe('Messages');
		expect(document.querySelector('[data-messages]')!.getAttribute('data-messages')).toBe('all');
		expect(listed()).toEqual(['MS-0001']);

		(document.querySelector('input[type="checkbox"]') as HTMLInputElement).click();
		await settle();
		expect(listed()).toEqual(['MS-0001', 'MS-0002']);
		expect(document.querySelector('[data-conversation="MS-0002"] [data-closed]')!.textContent).toBe(
			'closed: S-0001 was accepted'
		);
	});

	it('says so when the project has no conversation', async () => {
		api.mockResolvedValue(answer([]));
		c = mount(MessagesPage, { target: document.body });
		await settle();
		expect(listed()).toEqual([]);
		expect(document.querySelector('[data-empty]')!.textContent!.trim()).toBe(
			'No open conversations here.'
		);
	});
});
