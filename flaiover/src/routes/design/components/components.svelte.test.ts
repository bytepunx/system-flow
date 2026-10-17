// S-0439: the Components page lists the catalogue by sub-project and folder with a filter, frames
// the picked component from its dev server's harness in light, dark, or both at a chosen width,
// sends it the theme and the component, sizes the frame from the harness, shows a harness error in
// its place, starts a stopped dev server and shows a failed start's output, tells the operator to
// turn the host action on or upgrade flai, and asks again every minute. S-0440: beside the frame,
// the props panel sends the component its props, a change as a props message, keeps them in the
// URL, and shows a required prop without a value in the frame's place; the event log lists the
// calls of its function props.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type {
	Component,
	ComponentProp,
	ComponentsCatalogue,
	ComponentsStatus
} from '$lib/server/repo';
import { fromQuery, toQuery } from '$lib/explorer/controls';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const goto = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$app/navigation', () => ({ goto }));
vi.mock('$app/paths', () => ({ resolve: (route: string) => route }));
const where = vi.hoisted(() => ({ url: new URL('http://localhost/design/components') }));
vi.mock('$app/state', () => ({ page: where }));
const followed = vi.hoisted(() => ({ kinds: [] as string[][], f: () => {} }));
vi.mock('$lib/events', () => ({
	listen: () => () => {},
	follow: (kinds: string[], f: () => void) => {
		followed.kinds.push(kinds);
		followed.f = f;
		return () => {};
	}
}));

import ComponentsPage from './+page.svelte';

const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

const HARNESS = 'http://127.0.0.1:5199/__flai/components/';
const DEV = 'http://127.0.0.1:5199';
const ENABLE = 'flai serve enable components';
const BUTTON = 'src/lib/components/Button.svelte';
const SWATCH = 'src/lib/components/design/Swatch.svelte';
const BROKEN = 'src/lib/components/Broken.svelte';
const VITE_LINES = 'vite did not answer in 60s; its last lines:\nerror: port 5199 in use';

const component = (path: string, props: ComponentProp[] = []): Component => ({
	project: 'flaiover',
	path,
	name: path.slice(path.lastIndexOf('/') + 1, -'.svelte'.length),
	props
});
const SWATCH_PROPS: ComponentProp[] = [
	{ name: 'label', kind: 'string', type: 'string', required: true },
	{ name: 'count', kind: 'number', type: 'number', default: '3', required: false },
	{ name: 'items', kind: 'array', type: 'string[]', required: false },
	{ name: 'children', kind: 'snippet', type: 'Snippet', required: false },
	{ name: 'onpick', kind: 'function', type: '(n: number) => void', required: false }
];
const LOGGED = { functions: ['onpick'], snippets: ['children'] };
const catalogue: ComponentsCatalogue = {
	projects: [
		{
			name: 'flaiover',
			path: 'flaiover',
			kind: 'sveltekit',
			supported: true,
			components: [
				{ ...component(BROKEN), error: 'Unexpected token' },
				component(BUTTON),
				component(SWATCH, SWATCH_PROPS),
				component('src/routes/Chart.svelte')
			]
		},
		{ name: 'flai', path: 'flai', kind: 'go', supported: false, components: [] }
	]
};
const running: ComponentsStatus = {
	enabled: true,
	enable: ENABLE,
	idle_seconds: 600,
	servers: [{ project: 'flaiover', running: true, url: HARNESS, port: 5199, pid: 42 }]
};
const stopped: ComponentsStatus = {
	...running,
	servers: [{ project: 'flaiover', running: false }]
};

type Fake = { postMessage: ReturnType<typeof vi.fn> };
type Posted = Record<string, unknown>;

describe('the components page (S-0439)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		where.url = new URL('http://localhost/design/components');
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		goto.mockClear();
		followed.kinds = [];
		followed.f = () => {};
		vi.useRealTimers();
		document.body.innerHTML = '';
	});

	const serve = (status: ComponentsStatus = running) =>
		api.mockImplementation(async (url: string) =>
			url === '/api/components'
				? answer({ catalogue, status })
				: answer({ error: `unexpected ${url}` }, 404)
		);
	const q = (testid: string) => document.querySelector(`[data-testid="${testid}"]`);
	const gets = () => api.mock.calls.filter(([url]) => url === '/api/components').length;
	const show = async () => {
		c = mount(ComponentsPage, { target: document.body });
		await settle();
	};
	const click = async (el: Element | null) => {
		(el as HTMLElement).click();
		await settle();
	};
	const text = (e: Element) => e.textContent!.replace(/\s+/g, ' ').trim();
	const named = (b: Element, name: string) =>
		text(b) === name || b.getAttribute('aria-label') === name;
	const buttonNamed = (name: string) =>
		[...document.querySelectorAll('button')].find((b) => named(b, name))!;
	const frames = () =>
		[...document.querySelectorAll('[data-testid="component-frame"]')] as HTMLIFrameElement[];

	/** Gives each frame a fake window, which records what the page posts to it. */
	function fakeWindows(): Fake[] {
		return frames().map((f) => {
			const w: Fake = { postMessage: vi.fn() };
			Object.defineProperty(f, 'contentWindow', { value: w, configurable: true });
			return w;
		});
	}
	/** The harness in a frame answers, from the dev server's origin. */
	const harness = async (from: Fake, m: Record<string, unknown>) => {
		const e = new MessageEvent('message', {
			data: { source: 'flai-components', version: 1, ...m },
			origin: DEV
		});
		Object.defineProperty(e, 'source', { value: from });
		window.dispatchEvent(e);
		await settle();
	};
	const posted = (w: Fake) => w.postMessage.mock.calls.map(([m]) => m as Posted);

	it('lists each sub-project by folder, an unsupported one with its kind', async () => {
		serve();
		await show();
		expect(document.querySelector('h1')!.textContent).toBe('Components');
		const flaiover = q('component-project-flaiover')!;
		const folder = (f: Element) => ({
			dir: f.querySelector('h3')!.textContent,
			names: [...f.querySelectorAll('button')].map(text)
		});
		const folders = [...flaiover.querySelectorAll('[data-testid="component-folder"]')].map(folder);
		expect(folders).toEqual([
			{ dir: 'src/lib/components', names: ['Broken props unreadable', 'Button'] },
			{ dir: 'src/lib/components/design', names: ['Swatch'] },
			{ dir: 'src/routes', names: ['Chart'] }
		]);
		const flai = q('component-project-flai')!;
		expect(flai.textContent).toContain('go');
		expect(flai.querySelector('[data-testid="component-project-unsupported"]')).not.toBeNull();
		expect(flai.querySelectorAll('button')).toHaveLength(0);
		expect(q('components-frames')).toBeNull();
	});

	it('filters on name and on path, case-insensitively', async () => {
		serve();
		await show();
		const input = document.querySelector<HTMLInputElement>(
			'input[aria-label="Filter components by name or path"]'
		)!;
		const testid = (e: Element) => e.getAttribute('data-testid');
		const listed = () =>
			[...document.querySelectorAll('[data-testid="component-list"] li button')].map(testid);
		input.value = 'BUTT';
		input.dispatchEvent(new Event('input'));
		flushSync();
		expect(listed()).toEqual([`component-${BUTTON}`]);
		// the unsupported sub-project has no match, so it is left out under a filter
		expect(q('component-project-flai')).toBeNull();
		input.value = '/design/';
		input.dispatchEvent(new Event('input'));
		flushSync();
		expect(listed()).toEqual([`component-${SWATCH}`]);
		input.value = 'nothing-like-it';
		input.dispatchEvent(new Event('input'));
		flushSync();
		expect(q('component-list-none')!.textContent).toContain('nothing-like-it');
	});

	it('frames the picked component at the harness URL with the page origin, and keeps it in the URL', async () => {
		serve();
		await show();
		await click(q(`component-${BUTTON}`));
		expect(q(`component-${BUTTON}`)!.getAttribute('aria-current')).toBe('true');
		const [frame] = frames();
		const src = new URL(frame.src);
		expect(src.origin + src.pathname).toBe(HARNESS);
		expect(src.searchParams.get('origin')).toBe(location.origin);
		expect(src.searchParams.get('c')).toBe(BUTTON);
		expect(goto).toHaveBeenCalledWith(
			`/design/components?p=flaiover&c=${encodeURIComponent(BUTTON)}`,
			expect.objectContaining({ replaceState: true })
		);
	});

	it('frames the component the URL names on a reload', async () => {
		where.url = new URL(`http://localhost/design/components?p=flaiover&c=${SWATCH}`);
		serve();
		await show();
		expect(new URL(frames()[0].src).searchParams.get('c')).toBe(SWATCH);
		expect(q(`component-${SWATCH}`)!.getAttribute('aria-current')).toBe('true');
	});

	it('lists a component whose props are unreadable and still frames it', async () => {
		serve();
		await show();
		await click(q(`component-${BROKEN}`));
		expect(q('component-props-error')!.textContent).toContain('Unexpected token');
		expect(new URL(frames()[0].src).searchParams.get('c')).toBe(BROKEN);
	});

	it('posts the theme and the select to the dev server once the harness is ready, and a new pick as a select', async () => {
		serve();
		await show();
		await click(q(`component-${BUTTON}`));
		const [w] = fakeWindows();
		const src = frames()[0].src;
		expect(w.postMessage).not.toHaveBeenCalled();
		await harness(w, { type: 'ready' });
		expect(posted(w).map((m) => m.type)).toEqual(['theme', 'select']);
		expect(posted(w)[0]).toMatchObject({ source: 'flai-components', version: 1, theme: 'light' });
		expect(posted(w)[1]).toMatchObject({ component: BUTTON });
		expect(w.postMessage.mock.calls.every(([, origin]) => origin === DEV)).toBe(true);

		await click(q(`component-${SWATCH}`));
		expect(posted(w).at(-1)).toMatchObject({ type: 'select', component: SWATCH });
		// the frame is not reloaded for a new pick
		expect(frames()[0].src).toBe(src);
	});

	it('shows a harness error in place of the frame, with its file and line, until it renders again', async () => {
		serve();
		await show();
		await click(q(`component-${BUTTON}`));
		const [w] = fakeWindows();
		await harness(w, { type: 'ready' });
		await harness(w, {
			type: 'error',
			component: BUTTON,
			message: 'label is not defined',
			file: BUTTON,
			line: 12,
			column: 4
		});
		const shown = q('component-error')!;
		expect(shown.textContent).toContain('label is not defined');
		expect(q('component-error-place')!.textContent).toBe(`at ${BUTTON}:12:4`);
		expect((q('component-frame-box') as HTMLElement).hidden).toBe(true);
		// still mounted, so the next hot update's render comes back
		expect(frames()).toHaveLength(1);
		await harness(w, { type: 'rendered', component: BUTTON });
		expect(q('component-error')).toBeNull();
		expect((q('component-frame-box') as HTMLElement).hidden).toBe(false);
	});

	it('takes a message only from its own frame', async () => {
		serve();
		await show();
		await click(q(`component-${BUTTON}`));
		const [w] = fakeWindows();
		await harness({ postMessage: vi.fn() }, { type: 'size', width: 300, height: 999 });
		await harness(w, { type: 'size', width: 300, height: 412.4 });
		expect(frames()[0].style.height).toBe('413px');
	});

	it('frames both themes side by side, each sent its own theme', async () => {
		serve();
		await show();
		await click(q(`component-${BUTTON}`));
		await click(buttonNamed('Both'));
		expect(buttonNamed('Both').getAttribute('aria-pressed')).toBe('true');
		expect(q('component-frame-light')).not.toBeNull();
		expect(q('component-frame-dark')).not.toBeNull();
		const [light, dark] = fakeWindows();
		await harness(light, { type: 'ready' });
		await harness(dark, { type: 'ready' });
		expect(posted(light).find((m) => m.type === 'theme')).toMatchObject({ theme: 'light' });
		expect(posted(dark).find((m) => m.type === 'theme')).toMatchObject({ theme: 'dark' });

		await click(buttonNamed('Dark'));
		expect(frames()).toHaveLength(1);
		expect(q('component-frame-dark')).not.toBeNull();
		expect(posted(light).at(-1)).toMatchObject({ type: 'theme', theme: 'dark' });
	});

	it('sizes the frame to the width chosen, or the panel', async () => {
		serve();
		await show();
		await click(q(`component-${BUTTON}`));
		const box = () => q('component-frame-box') as HTMLElement;
		expect(buttonNamed("The panel's width").getAttribute('aria-pressed')).toBe('true');
		expect(box().style.width).toBe('');
		await click(buttonNamed('768 pixels wide'));
		// the frame is 768 pixels inside the box's 1-pixel border
		expect(box().style.width).toBe('770px');
		expect(q('component-frame-light')!.textContent).toContain('768px');
		await click(buttonNamed('375 pixels wide'));
		expect(box().style.width).toBe('377px');
		await click(buttonNamed("The panel's width"));
		expect(box().style.width).toBe('');
	});

	it('says to turn the host action on when the host refuses with 403', async () => {
		const off = { error: 'the action is off', action: 'components', enable: ENABLE };
		api.mockResolvedValue(answer(off, 403));
		await show();
		expect(q('components-off')!.textContent).toContain(ENABLE);
	});

	it('says to turn the host action on when the status says it is off, and offers no Start', async () => {
		serve({ ...stopped, enabled: false });
		await show();
		expect(q('components-off')!.textContent).toContain(ENABLE);
		expect(buttonNamed('Start')).toBeUndefined();
		// the catalogue is still listed
		expect(q(`component-${BUTTON}`)).not.toBeNull();
	});

	it('says to upgrade flai on the host when it answers 501', async () => {
		api.mockResolvedValue(
			answer({ error: 'This flai cannot list the components; upgrade flai on the host.' }, 501)
		);
		await show();
		expect(q('components-too-old')!.textContent).toMatch(/upgrade flai on the host/);
		expect(q('component-list')).toBeNull();
	});

	it('starts a stopped dev server for the picked sub-project and frames the component', async () => {
		where.url = new URL(`http://localhost/design/components?p=flaiover&c=${BUTTON}`);
		let release: () => void = () => {};
		api.mockImplementation(async (url: string, init?: RequestInit) => {
			if (url === '/api/components') return answer({ catalogue, status: stopped });
			if (url === '/api/components/start') {
				expect(init?.method).toBe('POST');
				expect(JSON.parse(String(init?.body))).toEqual({ project: 'flaiover' });
				await new Promise<void>((r) => (release = r));
				return answer(running);
			}
			return answer({}, 404);
		});
		await show();
		expect(frames()).toHaveLength(0);
		await click(buttonNamed('Start'));
		expect(q('components-starting')).not.toBeNull();
		release();
		await settle();
		expect(q('components-starting')).toBeNull();
		expect(frames()).toHaveLength(1);
		expect(new URL(frames()[0].src).searchParams.get('c')).toBe(BUTTON);
	});

	it('shows a server that is starting as waiting', async () => {
		serve({ ...running, servers: [{ project: 'flaiover', running: false, starting: true }] });
		await show();
		expect(q('components-starting')).not.toBeNull();
		expect(buttonNamed('Start')).toBeUndefined();
	});

	it("shows a failed start's output whole, and asks the status again for the last error", async () => {
		api.mockImplementation(async (url: string) => {
			if (url === '/api/components')
				return answer({
					catalogue,
					status: { ...stopped, servers: [{ ...stopped.servers[0], last_error: VITE_LINES }] }
				});
			if (url === '/api/components/start') return answer({ error: VITE_LINES }, 500);
			return answer({}, 404);
		});
		await show();
		expect(q('components-last-error')!.textContent).toBe(VITE_LINES);
		expect(gets()).toBe(1);
		await click(buttonNamed('Start'));
		expect(q('components-start-error')!.textContent).toBe(VITE_LINES);
		expect(q('components-start-error')!.tagName).toBe('PRE');
		expect(gets()).toBe(2);
		expect(buttonNamed('Start')).toBeDefined();
	});

	describe('the props panel and the event log (S-0440)', () => {
		const field = <T extends HTMLElement = HTMLInputElement>(name: string) =>
			document.querySelector<T>(`#props-panel-${name}`)!;
		const type = async (el: HTMLInputElement | HTMLTextAreaElement, text: string) => {
			el.value = text;
			el.dispatchEvent(new Event('input', { bubbles: true }));
			await settle();
		};
		// The props the URL was last given, as the page writes it.
		const urlProps = () => {
			const [href] = goto.mock.lastCall as unknown[];
			const url = new URL(String(href), 'http://localhost');
			return { c: url.searchParams.get('c'), props: fromQuery(url.searchParams.get('props')) };
		};
		const box = () => q('component-frame-box') as HTMLElement;
		const entries = () => [...document.querySelectorAll('[data-testid="event-log-entry"]')];

		it("sends the picked component its defaults, logs and snippets, and shows a required prop's want in the frame's place", async () => {
			serve();
			await show();
			await click(q(`component-${SWATCH}`));
			const [w] = fakeWindows();
			await harness(w, { type: 'ready' });
			expect(posted(w).map((m) => m.type)).toEqual(['theme', 'select']);
			expect(posted(w)[1]).toMatchObject({ component: SWATCH, props: { count: 3 }, ...LOGGED });
			expect(urlProps()).toEqual({ c: SWATCH, props: { count: 3 } });
			expect(q('prop-label')!.querySelector('[data-testid="prop-missing"]')).not.toBeNull();
			expect(q('component-missing')!.textContent).toContain('needs a value for label');
			expect(box().hidden).toBe(true);
			// the harness rendering it does not lift the want
			await harness(w, { type: 'rendered', component: SWATCH });
			expect(q('component-missing')).not.toBeNull();
		});

		it('sends each change as props without remounting, keeps it in the URL, and holds back invalid JSON', async () => {
			serve();
			await show();
			await click(q(`component-${SWATCH}`));
			const [w] = fakeWindows();
			await harness(w, { type: 'ready' });
			await type(field('label'), 'Go');
			expect(posted(w).at(-1)).toEqual({
				source: 'flai-components',
				version: 1,
				type: 'props',
				props: { count: 3, label: 'Go' },
				...LOGGED
			});
			expect(urlProps()).toEqual({ c: SWATCH, props: { count: 3, label: 'Go' } });
			expect(q('component-missing')).toBeNull();
			expect(box().hidden).toBe(false);

			const sent = posted(w).length;
			const written = goto.mock.calls.length;
			await type(field<HTMLTextAreaElement>('items'), '["a",');
			expect(q('prop-items')!.querySelector('[data-testid="prop-error"]')).not.toBeNull();
			expect(field<HTMLTextAreaElement>('items').value).toBe('["a",');
			expect(posted(w)).toHaveLength(sent);
			expect(goto.mock.calls).toHaveLength(written);
			// another prop's change is held back too, until the JSON parses
			await type(field('count'), '4');
			expect(posted(w)).toHaveLength(sent);
			await type(field<HTMLTextAreaElement>('items'), '["a", "b"]');
			expect(posted(w).at(-1)).toMatchObject({
				type: 'props',
				props: { count: 4, label: 'Go', items: ['a', 'b'] }
			});
			await type(field<HTMLTextAreaElement>('children'), 'Hello');
			expect(posted(w).at(-1)).toMatchObject({ props: { children: 'Hello' }, ...LOGGED });

			(q('props-reset') as HTMLElement).click();
			await settle();
			expect(posted(w).at(-1)).toMatchObject({ type: 'props', props: { count: 3 } });
			expect(Object.keys(posted(w).at(-1)!.props as object)).toEqual(['count']);
			expect(urlProps().props).toEqual({ count: 3 });
			expect(q('component-missing')).not.toBeNull();
			// one select: every change went to the component mounted
			expect(posted(w).filter((m) => m.type === 'select')).toHaveLength(1);
		});

		it('renders the props a link gives, and starts another component from its defaults', async () => {
			const props = toQuery({ label: 'Hi', count: 5, items: ['x'] });
			where.url = new URL(
				`http://localhost/design/components?p=flaiover&c=${SWATCH}&props=${props}`
			);
			serve();
			await show();
			expect(field('label').value).toBe('Hi');
			expect(field('count').value).toBe('5');
			expect(JSON.parse(field<HTMLTextAreaElement>('items').value)).toEqual(['x']);
			const [w] = fakeWindows();
			await harness(w, { type: 'ready' });
			expect(posted(w).find((m) => m.type === 'select')).toMatchObject({
				component: SWATCH,
				props: { label: 'Hi', count: 5, items: ['x'] },
				...LOGGED
			});
			expect(q('component-missing')).toBeNull();

			await click(q(`component-${BUTTON}`));
			expect(urlProps()).toEqual({ c: BUTTON, props: {} });
			await click(q(`component-${SWATCH}`));
			expect(field('label').value).toBe('');
			expect(field('count').value).toBe('3');
		});

		it("lists the calls of the component's function props, and clears them on another pick", async () => {
			serve();
			await show();
			await click(q(`component-${SWATCH}`));
			const [w] = fakeWindows();
			await harness(w, { type: 'ready' });
			expect(q('event-log-none')).not.toBeNull();
			await harness(w, { type: 'event', component: SWATCH, name: 'onpick', args: [2] });
			await harness(w, { type: 'event', component: SWATCH, name: 'onpick', args: [{ n: 3 }] });
			// a late call of the component picked before is not its
			await harness(w, { type: 'event', component: BUTTON, name: 'onclick', args: [] });
			expect(entries().map((e) => e.querySelector('code')!.textContent)).toEqual([
				'[{"n":3}]',
				'[2]'
			]);
			expect(entries()[0].textContent).toContain('onpick');
			expect(entries()[0].querySelector('time')).not.toBeNull();
			for (let i = 0; i < 55; i++)
				await harness(w, { type: 'event', component: SWATCH, name: 'onpick', args: [i] });
			expect(entries()).toHaveLength(50);
			expect(entries()[0].querySelector('code')!.textContent).toBe('[54]');

			await click(q(`component-${BUTTON}`));
			expect(entries()).toHaveLength(0);
		});
	});

	it('asks the status again every minute, and on a change to the project', async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		serve();
		await show();
		expect(gets()).toBe(1);
		await vi.advanceTimersByTimeAsync(59_000);
		await settle();
		expect(gets()).toBe(1);
		await vi.advanceTimersByTimeAsync(1_000);
		await settle();
		expect(gets()).toBe(2);
		expect(followed.kinds).toEqual([['document']]);
		followed.f();
		await settle();
		expect(gets()).toBe(3);
		unmount(c!);
		c = undefined;
		await vi.advanceTimersByTimeAsync(120_000);
		expect(gets()).toBe(3);
	});
});
