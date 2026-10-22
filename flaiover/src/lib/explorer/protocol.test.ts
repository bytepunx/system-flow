import { describe, expect, it, vi } from 'vitest';
import {
	connect,
	harnessMessage,
	harnessUrl,
	type FrameWindow,
	type Handlers,
	type PageWindow
} from './protocol';

// S-0439: the page's side of the harness protocol, against a fake frame and page window.
const SERVER = 'http://127.0.0.1:5173/__flai/components/';
const ORIGIN = 'http://127.0.0.1:5173';
const READY = { source: 'flai-components', version: 1, type: 'ready' };

function fakeFrame() {
	const posted: { message: Record<string, unknown>; targetOrigin: string }[] = [];
	const win: FrameWindow = {
		postMessage: (message, targetOrigin) =>
			void posted.push({ message: message as Record<string, unknown>, targetOrigin })
	};
	return { win, posted, types: () => posted.map((p) => p.message.type) };
}

function fakePage() {
	const listeners = new Set<(e: MessageEvent) => void>();
	const page: PageWindow = {
		addEventListener: (_type, l) => void listeners.add(l),
		removeEventListener: (_type, l) => void listeners.delete(l)
	};
	function deliver(data: unknown, origin: string, source: unknown) {
		const e = { data, origin, source } as MessageEvent;
		for (const l of [...listeners]) l(e);
	}
	return { page, listeners, deliver };
}

function setup(handlers: Handlers = {}) {
	const frame = fakeFrame();
	const page = fakePage();
	const client = connect(() => frame.win, SERVER, handlers, page.page);
	const answer = (m: Record<string, unknown>) =>
		page.deliver({ source: 'flai-components', version: 1, ...m }, ORIGIN, frame.win);
	return { frame, page, client, answer };
}

describe('connect', () => {
	it('queues what is sent before ready and sends it on ready, to the dev server origin', () => {
		const { frame, client, answer } = setup();
		client.select('src/lib/Button.svelte', { label: 'Go' }, ['onclick']);
		client.theme('dark');
		expect(frame.posted).toEqual([]);
		answer({ type: 'ready' });
		expect(frame.posted).toEqual([
			{
				message: { source: 'flai-components', version: 1, type: 'theme', theme: 'dark' },
				targetOrigin: ORIGIN
			},
			{
				message: {
					source: 'flai-components',
					version: 1,
					type: 'select',
					component: 'src/lib/Button.svelte',
					props: { label: 'Go' },
					functions: ['onclick']
				},
				targetOrigin: ORIGIN
			}
		]);
	});

	it('names the functions to log and the snippets to render on a select and on props (S-0440)', () => {
		const { frame, client, answer } = setup();
		answer({ type: 'ready' });
		client.select('A.svelte', { label: 'Go', children: 'Hi' }, ['onclick'], ['children']);
		client.props({ label: 'Stop', children: 'Bye' }, ['onclick'], ['children']);
		expect(frame.posted.map((p) => p.message)).toEqual([
			{
				source: 'flai-components',
				version: 1,
				type: 'select',
				component: 'A.svelte',
				props: { label: 'Go', children: 'Hi' },
				functions: ['onclick'],
				snippets: ['children']
			},
			{
				source: 'flai-components',
				version: 1,
				type: 'props',
				props: { label: 'Stop', children: 'Bye' },
				functions: ['onclick'],
				snippets: ['children']
			}
		]);
	});

	it('carries the answers on a select and on props, each whole (S-0476)', () => {
		const { frame, client, answer } = setup();
		answer({ type: 'ready' });
		const threads = [{ path: '/api/threads', body: [{ id: 'TH-1' }] }];
		const events = [{ path: '/api/events', events: [{ event: 'change', data: { kind: 'item' } }] }];
		client.select('A.svelte', {}, [], [], threads);
		client.props({ n: 1 }, [], [], events);
		client.props({ n: 1 });
		expect(frame.posted.map((p) => p.message)).toEqual([
			{
				source: 'flai-components',
				version: 1,
				type: 'select',
				component: 'A.svelte',
				props: {},
				functions: [],
				snippets: [],
				answers: threads
			},
			{
				source: 'flai-components',
				version: 1,
				type: 'props',
				props: { n: 1 },
				functions: [],
				snippets: [],
				answers: events
			},
			{ source: 'flai-components', version: 1, type: 'props', props: { n: 1 } }
		]);
		expect(frame.posted[2].message.answers).toBeUndefined();
	});

	it('keeps only the last of each kind before ready, in theme, tokens, select, props order', () => {
		const { frame, client, answer } = setup();
		client.props({ n: 1 });
		client.theme('light');
		client.select('A.svelte');
		client.props({ n: 2 });
		client.props({ n: 3 });
		client.tokens({ accent: 'red' }, {});
		client.theme('dark');
		answer({ type: 'ready' });
		expect(frame.posted.map((p) => p.message)).toMatchObject([
			{ type: 'theme', theme: 'dark' },
			{ type: 'tokens', light: { accent: 'red' }, dark: {} },
			{ type: 'select', component: 'A.svelte' },
			{ type: 'props', props: { n: 3 } }
		]);
	});

	it('drops the props given before a select, which carries its own', () => {
		const { frame, client, answer } = setup();
		client.select('A.svelte');
		client.props({ n: 1 });
		client.select('B.svelte', { m: 2 });
		answer({ type: 'ready' });
		expect(frame.posted.map((p) => p.message)).toMatchObject([
			{ type: 'select', component: 'B.svelte', props: { m: 2 } }
		]);
	});

	it('sends at once after ready', () => {
		const { frame, client, answer } = setup();
		answer({ type: 'ready' });
		client.select('A.svelte');
		client.props({ n: 1 });
		expect(frame.types()).toEqual(['select', 'props']);
		expect(frame.posted.every((p) => p.targetOrigin === ORIGIN)).toBe(true);
	});

	it('sends the last of each kind again on a later ready, as after a reload', () => {
		const ready = vi.fn();
		const { frame, client, answer } = setup({ ready });
		answer({ type: 'ready' });
		client.select('A.svelte');
		client.props({ n: 1 });
		client.theme('dark');
		client.props({ n: 2 });
		frame.posted.length = 0;
		answer({ type: 'ready', component: 'A.svelte' });
		expect(frame.posted.map((p) => p.message)).toMatchObject([
			{ type: 'theme', theme: 'dark' },
			{ type: 'select', component: 'A.svelte' },
			{ type: 'props', props: { n: 2 } }
		]);
		expect(ready).toHaveBeenCalledTimes(2);
		expect(ready).toHaveBeenLastCalledWith(
			expect.objectContaining({ type: 'ready', component: 'A.svelte' })
		);
	});

	it('ignores a message from another origin', () => {
		const ready = vi.fn();
		const { frame, page, client } = setup({ ready });
		client.theme('dark');
		page.deliver(READY, 'http://127.0.0.1:5174', frame.win);
		expect(ready).not.toHaveBeenCalled();
		expect(frame.posted).toEqual([]);
	});

	it('ignores a message from another window with the right origin', () => {
		const ready = vi.fn();
		const { frame, page, client } = setup({ ready });
		client.theme('dark');
		page.deliver(READY, ORIGIN, fakeFrame().win);
		page.deliver(READY, ORIGIN, null);
		expect(ready).not.toHaveBeenCalled();
		expect(frame.posted).toEqual([]);
	});

	it('ignores everything while the frame has no window', () => {
		const ready = vi.fn();
		const page = fakePage();
		connect(() => null, SERVER, { ready }, page.page);
		page.deliver(READY, ORIGIN, null);
		expect(ready).not.toHaveBeenCalled();
	});

	it.each([
		['another source', { source: 'other', version: 1, type: 'ready' }],
		['another version', { source: 'flai-components', version: 2, type: 'ready' }],
		['no version', { source: 'flai-components', type: 'ready' }],
		['a string', 'ready'],
		['null', null]
	])('ignores %s', (_name, data) => {
		const ready = vi.fn();
		const { frame, page } = setup({ ready });
		page.deliver(data, ORIGIN, frame.win);
		expect(ready).not.toHaveBeenCalled();
	});

	it('gives each answer to its handler', () => {
		const handlers = {
			ready: vi.fn(),
			rendered: vi.fn(),
			error: vi.fn(),
			event: vi.fn(),
			size: vi.fn()
		};
		const { answer } = setup(handlers);
		const answers = {
			ready: { type: 'ready', component: 'A.svelte' },
			rendered: { type: 'rendered', component: 'A.svelte' },
			error: {
				type: 'error',
				component: 'A.svelte',
				message: 'boom',
				file: 'src/A.svelte',
				line: 3,
				column: 7
			},
			event: { type: 'event', component: 'A.svelte', name: 'onclick', args: [{ type: 'click' }] },
			size: { type: 'size', width: 120, height: 40 }
		};
		for (const m of Object.values(answers)) answer(m);
		for (const [type, m] of Object.entries(answers)) {
			const h = handlers[type as keyof typeof handlers];
			expect(h).toHaveBeenCalledTimes(1);
			expect(h).toHaveBeenCalledWith(expect.objectContaining(m));
		}
	});

	it('ignores an answer without the fields its type needs, or of an unknown type', () => {
		const handlers = { rendered: vi.fn(), error: vi.fn(), event: vi.fn(), size: vi.fn() };
		const { answer } = setup(handlers);
		answer({ type: 'rendered' });
		answer({ type: 'error', component: 'A.svelte' });
		answer({ type: 'event', component: 'A.svelte', name: 'onclick' });
		answer({ type: 'size', width: '1px', height: 2 });
		answer({ type: 'select', component: 'A.svelte' });
		for (const h of Object.values(handlers)) expect(h).not.toHaveBeenCalled();
	});

	it('stops listening and sending once closed', () => {
		const ready = vi.fn();
		const { frame, page, client, answer } = setup({ ready });
		client.theme('dark');
		expect(page.listeners.size).toBe(1);
		client.close();
		expect(page.listeners.size).toBe(0);
		answer({ type: 'ready' });
		client.select('A.svelte');
		expect(ready).not.toHaveBeenCalled();
		expect(frame.posted).toEqual([]);
	});

	it('takes the origin of the URL it is given', () => {
		const frame = fakeFrame();
		const page = fakePage();
		const client = connect(() => frame.win, ORIGIN, {}, page.page);
		page.deliver(READY, ORIGIN, frame.win);
		client.theme('light');
		expect(frame.posted).toMatchObject([{ targetOrigin: ORIGIN }]);
	});
});

describe('harnessUrl', () => {
	it("tells the harness page the page's origin and the component to render first", () => {
		const u = new URL(harnessUrl(SERVER, 'http://127.0.0.1:4242', 'src/lib/A B.svelte'));
		expect(u.origin + u.pathname).toBe(SERVER);
		expect(u.searchParams.get('origin')).toBe('http://127.0.0.1:4242');
		expect(u.searchParams.get('c')).toBe('src/lib/A B.svelte');
	});

	it('names no component when given none, and the origin alone of a page URL', () => {
		const u = new URL(harnessUrl(SERVER, 'http://127.0.0.1:4242/components?x=1'));
		expect(u.searchParams.get('origin')).toBe('http://127.0.0.1:4242');
		expect(u.searchParams.has('c')).toBe(false);
	});
});

describe('harnessMessage', () => {
	it('reads a harness message, and nothing else', () => {
		expect(harnessMessage(READY)).toEqual(READY);
		expect(harnessMessage({ ...READY, type: 'toString' })).toBeNull();
		expect(harnessMessage(undefined)).toBeNull();
	});
});
