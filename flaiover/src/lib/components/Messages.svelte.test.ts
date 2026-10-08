// S-0336: the conversations between two stories' agents, open before closed, read-only, with the
// stories, the paths, whom each awaits, its age, its entries, and the thread an escalation opened.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Messages from './Messages.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

// The page's live events, as follow() would deliver them, with no wait for changes to settle.
const events = vi.hoisted(() => ({ heard: [] as { kinds: string[]; f: () => void }[] }));
vi.mock('$lib/events', () => ({
	follow: (kinds: string[], f: () => void) => {
		const l = { kinds, f };
		events.heard.push(l);
		return () => events.heard.splice(events.heard.indexOf(l), 1);
	}
}));

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, json: async () => body });

const NOW = Date.parse('2026-10-08T10:00:00Z');
const conversation = (id: string, over: Record<string, unknown> = {}) => ({
	id,
	title: `about ${id}`,
	from: 'S-0001',
	to: 'S-0002',
	about: ['flaiover/src/lib/sitemenu.ts'],
	status: 'open',
	closed: false,
	closed_reason: '',
	awaiting: 'S-0002',
	participants: ['agent-S-0001'],
	created: '2026-10-08T06:00:00Z',
	updated: '2026-10-08T07:00:00Z',
	path: `wip/messages/${id}-x.md`,
	entries: [
		{
			at: '2026-10-08T07:00:00Z',
			author: 'agent-S-0001',
			story: 'S-0001',
			text: 'Will you **keep** `locate`?'
		}
	],
	...over
});
const shown = () => [...document.querySelectorAll('[data-conversation]')] as HTMLElement[];

describe('Messages (S-0336)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => vi.spyOn(Date, 'now').mockReturnValue(NOW));
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		vi.restoreAllMocks();
		document.body.innerHTML = '';
	});

	it('lists each open conversation with its two stories, its paths, whom it awaits, its age, and its entries', async () => {
		api.mockResolvedValue(answer([conversation('MS-0001')]));
		c = mount(Messages, { target: document.body });
		await settle();

		expect(api).toHaveBeenCalledWith('/api/messages');
		const [m] = shown();
		expect(m.dataset.conversation).toBe('MS-0001');
		expect(m.textContent).toContain('about MS-0001');
		const between = [...m.querySelectorAll('[data-between] a')];
		expect(between.map((a) => a.getAttribute('href'))).toEqual(['/items/S-0001', '/items/S-0002']);
		expect(m.querySelector('[data-awaiting]')!.textContent).toBe('awaits S-0002');
		expect(m.querySelector('[data-about]')!.textContent).toBe('flaiover/src/lib/sitemenu.ts');
		expect(m.querySelector('[data-age]')!.textContent).toBe('updated 3h ago');
		// the tests run in New York
		expect(m.querySelector('[data-age]')!.getAttribute('title')).toBe('2026-10-08 03:00 EDT');
		expect(m.querySelector('details')!.open).toBe(true);
		const entry = m.querySelector('[data-entry]')!;
		expect(entry.textContent).toContain('2026-10-08 03:00 EDT');
		expect(entry.textContent).toContain('agent-S-0001');
		expect(entry.querySelector('strong')?.textContent).toBe('keep');
		expect(entry.querySelector('code')?.textContent).toBe('locate');
		expect(m.querySelector('form, textarea, input:not([type="checkbox"])')).toBeNull();
	});

	it('says a closed conversation is closed and why, with its entries on request', async () => {
		api.mockResolvedValue(
			answer([
				conversation('MS-0002', {
					status: 'closed',
					closed: true,
					closed_reason: 'S-0002 was accepted',
					awaiting: ''
				})
			])
		);
		c = mount(Messages, { target: document.body });
		await settle();

		const [m] = shown();
		expect(m.querySelector('header')!.textContent).toContain('closed');
		expect(m.querySelector('[data-closed]')!.textContent).toBe('closed: S-0002 was accepted');
		expect(m.querySelector('[data-awaiting]')).toBeNull();
		expect(m.querySelector('details')!.open).toBe(false);
	});

	it('shows the closed ones too when asked', async () => {
		api.mockResolvedValue(answer([conversation('MS-0001')]));
		c = mount(Messages, { target: document.body });
		await settle();

		api.mockResolvedValue(
			answer([conversation('MS-0001'), conversation('MS-0002', { status: 'closed', closed: true })])
		);
		const box = document.querySelector('input[type="checkbox"]') as HTMLInputElement;
		box.click();
		await settle();
		expect(api).toHaveBeenLastCalledWith('/api/messages?all=1');
		expect(shown().map((m) => m.dataset.conversation)).toEqual(['MS-0001', 'MS-0002']);
	});

	it("on a story's page reads that story's conversations and names the other story", async () => {
		api.mockResolvedValue(
			answer([
				conversation('MS-0001', { awaiting: 'S-0001' }),
				conversation('MS-0003', { from: 'S-0003', to: 'S-0001', awaiting: 'S-0003' })
			])
		);
		c = mount(Messages, { target: document.body, props: { story: 'S-0001' } });
		await settle();

		expect(api).toHaveBeenCalledWith('/api/messages?story=S-0001');
		const [to, from] = shown();
		expect(to.querySelector('[data-with]')!.textContent).toBe('to S-0002');
		expect(to.querySelector('[data-awaiting]')!.textContent).toBe('awaits this story');
		expect(from.querySelector('[data-with]')!.textContent).toBe('from S-0003');
		expect(from.querySelector('[data-with] a')!.getAttribute('href')).toBe('/items/S-0003');
		expect(from.querySelector('[data-awaiting]')!.textContent).toBe('awaits S-0003');
		expect(document.querySelector('[data-between]')).toBeNull();
	});

	it('links an escalated conversation to the thread its story opened for the operator', async () => {
		const escalated = conversation('MS-0004', {
			awaiting: 'S-0001',
			entries: [
				conversation('MS-0004').entries[0],
				{
					at: '2026-10-08T08:00:00Z',
					author: 'agent-S-0002',
					story: 'S-0002',
					text: 'Asked the operator on TH-0042, `wip/threads/TH-0042-x.md`, since we do not agree. What we could not agree:\n\n> who keeps it\n\nThe operator answers on the thread. This conversation stays open.'
				}
			]
		});
		api.mockResolvedValue(answer([escalated, conversation('MS-0001')]));
		c = mount(Messages, { target: document.body });
		await settle();

		const [m, plain] = shown();
		const link = m.querySelector('[data-escalated]')!;
		expect(link.textContent).toBe('escalated on TH-0042');
		expect(link.getAttribute('href')).toBe('/items/S-0002?thread=TH-0042');
		expect(m.querySelector('[data-closed]')).toBeNull();
		expect(m.querySelector('blockquote')?.textContent).toContain('who keeps it');
		expect(plain.querySelector('[data-escalated]')).toBeNull();
	});

	it('says when there is no open conversation, and when there is none at all', async () => {
		api.mockResolvedValue(answer([]));
		c = mount(Messages, { target: document.body });
		await settle();
		expect(document.querySelector('[data-empty]')!.textContent!.trim()).toBe(
			'No open conversations here.'
		);

		(document.querySelector('input[type="checkbox"]') as HTMLInputElement).click();
		await settle();
		expect(document.querySelector('[data-empty]')!.textContent!.trim()).toBe(
			'No conversations here.'
		);
	});

	it('says what went wrong when the conversations cannot be read', async () => {
		api.mockResolvedValue({
			ok: false,
			statusText: 'Service Unavailable',
			json: async () => ({ error: 'no host flai is connected' })
		});
		c = mount(Messages, { target: document.body });
		await settle();
		expect(document.querySelector('.text-danger')!.textContent).toBe(
			'could not read the conversations: no host flai is connected'
		);
		expect(document.querySelector('[data-empty]')).toBeNull();
	});

	it('reads them again when a document or an item changes, and stops once gone', async () => {
		api.mockResolvedValue(answer([]));
		c = mount(Messages, { target: document.body });
		await settle();
		expect(events.heard.map((l) => l.kinds)).toEqual([['document', 'item']]);

		api.mockResolvedValue(answer([conversation('MS-0001')]));
		events.heard.forEach((l) => l.f());
		await settle();
		expect(shown()).toHaveLength(1);

		unmount(c);
		c = undefined;
		expect(events.heard).toHaveLength(0);
	});
});
