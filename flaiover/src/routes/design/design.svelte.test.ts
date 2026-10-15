// S-0382: the design page shows the project's design system, its token sheet, its guidance, and
// the live preview; tells a project with no token file so, with the command that writes one; and
// reads it all again when a document changes, as the other pages follow theirs.
// S-0383: in Edit mode the tokens change in place with the preview following before saving; Save
// writes the edits into the token file through /api/design/tokens, and a refusal, a conflict, or a
// flai too old to save is shown with the edits kept; Discard goes back to the saved tokens; leaving
// with unsaved changes asks first.
// S-0384: the page lists the open stories that propose a change to the design system, compares
// the one chosen with the main branch's tokens, closes it when it leaves the list, and offers the
// form that requests a change under an open epic.
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { DesignSystem, DesignToken } from '$lib/server/repo';
import { composeTokenFile, editKey } from '$lib/design/tokens';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const navigation = vi.hoisted(() => ({
	before: null as null | ((nav: { cancel: () => void }) => void)
}));
vi.mock('$app/navigation', () => ({
	beforeNavigate: (f: (nav: { cancel: () => void }) => void) => {
		navigation.before = f;
	}
}));
vi.mock('$app/paths', () => ({
	resolve: (route: string, p?: { id: string }) => route.replace('[id]', p?.id ?? '')
}));
const enhance = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$lib/markdown', async (actual) => ({
	...(await actual<typeof import('$lib/markdown')>()),
	enhance
}));
// Each follow's callback, by the kinds it follows joined with commas.
const followed = vi.hoisted(() => ({
	kinds: [] as string[][],
	of: {} as Record<string, () => void>
}));
vi.mock('$lib/events', () => ({
	listen: () => () => {},
	follow: (kinds: string[], f: () => void) => {
		followed.kinds.push(kinds);
		followed.of[kinds.join()] = f;
		return () => {};
	}
}));

import DesignPage from './+page.svelte';

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

const colour = (name: string, hex: string): DesignToken => ({
	path: `color.${name}`,
	type: 'color',
	value: { colorSpace: 'srgb', hex },
	css: hex,
	variables: [{ name: `--color-${name}`, value: hex }],
	line: 1
});
const GUIDANCE = 'design/system/design-system.md';
const TOKENS = 'design/system/design-system.tokens.json';
const NO_FILE = `${TOKENS} does not exist: write it with flai design init.`;

const withSystem: DesignSystem = {
	path: TOKENS,
	exists: true,
	guidance: GUIDANCE,
	headings: ['Design system', 'Colour'],
	resolver: true,
	themes: [
		{ name: 'light', tokens: [colour('background', '#ffffff'), colour('primary', '#2563eb')] },
		{ name: 'dark', tokens: [colour('background', '#0f172a'), colour('primary', '#60a5fa')] }
	],
	problems: [{ path: 'color.brand', line: 30, message: 'color.brand has no $value' }],
	css: '',
	tailwind: '',
	contrast: []
};
const without: DesignSystem = {
	path: TOKENS,
	exists: false,
	no_token_file: NO_FILE,
	guidance: '',
	headings: [],
	resolver: false,
	themes: [],
	problems: [],
	css: '',
	tailwind: '',
	contrast: []
};
const guidanceDoc = {
	path: GUIDANCE,
	frontMatter: { title: 'Design system' },
	body: '\n# Design system\n\nUse the primary colour for one action per view.\n'
};
const guidanceUrl = `/api/docs/file?path=${encodeURIComponent(GUIDANCE)}`;
const PROPOSALS = '/api/design/proposals';
const EPICS = '/api/items?type=epic&archived=false';

describe('the design page (S-0382)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		followed.kinds = [];
		followed.of = {};
		document.body.innerHTML = '';
	});

	function serve(design: () => unknown, guidance: () => unknown = () => answer(guidanceDoc)) {
		api.mockImplementation(async (url: string) => {
			if (url === '/api/design') return design();
			if (url === guidanceUrl) return guidance();
			if (url === PROPOSALS || url === EPICS) return answer([]);
			return answer({ error: `unexpected ${url}` }, 404);
		});
	}
	const q = (testid: string) => document.querySelector(`[data-testid="${testid}"]`);
	const show = async () => {
		c = mount(DesignPage, { target: document.body });
		await settle();
	};

	it('shows the token sheet, the problems, the guidance rendered, and the preview in both themes', async () => {
		serve(() => answer(withSystem));
		await show();
		expect(document.querySelector('h1')!.textContent).toBe('Design system');
		expect(q('token-sheet')).not.toBeNull();
		expect(q('swatch-color.primary')).not.toBeNull();
		expect(q('design-problems')!.textContent).toContain('line 30');
		expect(q('design-problems')!.textContent).toContain('color.brand has no $value');
		const article = q('design-guidance')!.querySelector('article')!;
		expect(article.querySelector('h1')!.textContent).toBe('Design system');
		expect(article.textContent).toContain('Use the primary colour for one action per view.');
		expect(enhance).toHaveBeenCalled();
		const frame = q('design-preview') as HTMLIFrameElement;
		expect(frame.tagName).toBe('IFRAME');
		const panels = [...frame.contentDocument!.querySelectorAll<HTMLElement>('[data-ds-theme]')];
		expect(panels.map((p) => p.dataset.dsTheme)).toEqual(['light', 'dark']);
		expect(frame.contentDocument!.querySelector('style[data-ds-tokens]')!.textContent).toContain(
			'--ds-color-primary: #60a5fa;'
		);
	});

	it('says so when the guidance document is missing, and still shows the tokens', async () => {
		serve(
			() => answer(withSystem),
			() => answer({ error: 'not found' }, 404)
		);
		await show();
		expect(q('design-guidance')!.textContent).toContain(
			`The guidance document ${GUIDANCE} is missing.`
		);
		expect(q('design-guidance')!.querySelector('article')).toBeNull();
		expect(q('token-sheet')).not.toBeNull();
	});

	it('tells a project with no token file so, with the command that writes one, and nothing else', async () => {
		serve(() => answer(without));
		await show();
		const none = q('no-design-system')!;
		expect(none.textContent).toContain(NO_FILE);
		expect([...none.querySelectorAll('code')].map((e) => e.textContent)).toContain(
			'flai design init'
		);
		expect(q('token-sheet')).toBeNull();
		expect(q('design-preview')).toBeNull();
		expect(q('design-guidance')).toBeNull();
		expect(api).not.toHaveBeenCalledWith(guidanceUrl);
	});

	it('reads the design system again when a document changes', async () => {
		let current: DesignSystem = without;
		serve(() => answer(current));
		await show();
		expect(followed.kinds).toEqual([['document'], ['item']]);
		expect(q('no-design-system')).not.toBeNull();

		current = withSystem;
		followed.of.document();
		await settle();
		expect(q('no-design-system')).toBeNull();
		expect(q('token-sheet')).not.toBeNull();
		expect(q('design-guidance')!.textContent).toContain('one action per view');
		expect(q('design-preview')).not.toBeNull();
	});

	it('shows the error of a host flai older than design.get, as the other pages show theirs', async () => {
		const older = 'flai on the host is older than this dashboard: it does not offer design.get';
		serve(() => answer({ error: older }, 501));
		await show();
		expect(document.querySelector('[role="alert"]')!.textContent!.trim()).toBe(older);
		expect(q('token-sheet')).toBeNull();
	});
});

// flai's design.get of its test token file, and that file's raw text, as flai's tests hold them.
const TESTDATA = resolve('../flai/internal/designsystem/testdata');
const SHOWN = JSON.parse(readFileSync(`${TESTDATA}/show.json.golden`, 'utf8')) as DesignSystem;
const RAW = readFileSync(`${TESTDATA}/themes.tokens.json`, 'utf8');
// The token file as someone else saved it meanwhile: without font.mono.
const CURRENT = RAW.replace(/,\s*"mono": \{[^}]*\}/, '');

const PRIMARY = { colorSpace: 'srgb', components: [0.2, 0.4, 0.6], hex: '#336699' };
const MONO = ['Fira Code', 'monospace'];
const CANNOT_SAVE = 'This flai cannot save the token file yet; upgrade flai on the host.';

describe('editing the design system on the design page (S-0383)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let design: DesignSystem;
	const readRaw = () => answer({ path: TOKENS, content: RAW, hash: 'h1' });
	const saveRaw = () => answer({ path: TOKENS, hash: 'h2', log: '' });
	let get: () => unknown = readRaw;
	let put: (body: Record<string, unknown>) => unknown = saveRaw;
	let gets: number;
	let puts: Record<string, unknown>[];
	const confirm = vi.fn(() => false);

	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		confirm.mockClear();
		vi.unstubAllGlobals();
		get = readRaw;
		put = saveRaw;
		followed.of = {};
		navigation.before = null;
		document.body.innerHTML = '';
	});

	// Serves design.get, the guidance, and the token file's GET and PUT.
	function serve() {
		gets = 0;
		puts = [];
		api.mockImplementation(async (url: string, init?: RequestInit) => {
			if (url === '/api/design') return answer(design);
			if (url === guidanceUrl) return answer(guidanceDoc);
			if (url === PROPOSALS || url === EPICS) return answer([]);
			if (url === '/api/design/tokens' && init?.method === 'PUT') {
				const body = JSON.parse(String(init.body)) as Record<string, unknown>;
				puts.push(body);
				return put(body);
			}
			if (url === '/api/design/tokens') {
				gets++;
				return get();
			}
			return answer({ error: `unexpected ${url}` }, 404);
		});
	}

	// Mounts the page on the test token file, and enters Edit mode.
	async function edit() {
		vi.stubGlobal('confirm', confirm);
		design = structuredClone(SHOWN);
		serve();
		c = mount(DesignPage, { target: document.body });
		await settle();
		button('Edit tokens').click();
		await settle();
	}

	const q = (testid: string) => document.querySelector<HTMLElement>(`[data-testid="${testid}"]`);
	const words = (e: Element | null) => e!.textContent!.replace(/\s+/g, ' ').trim();
	function button(text: string) {
		const found = [...document.querySelectorAll('button')].find(
			(b) => b.textContent!.trim() === text
		);
		if (!found) throw new Error(`no button ${text}`);
		return found;
	}
	const field = (label: string) => {
		const found = document.querySelector<HTMLInputElement>(`[aria-label="${label}"]`);
		if (!found) throw new Error(`no field labelled ${label}`);
		return found;
	};
	function type(label: string, text: string) {
		const input = field(label);
		input.value = text;
		input.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
	}
	const changedMarks = () =>
		[...document.querySelectorAll('[data-testid="changed"]')].map(
			(m) => m.closest<HTMLElement>('[data-testid^="edit-"]')!.dataset.testid
		);
	// The preview's token stylesheet, split at its dark block.
	function preview() {
		const frame = q('design-preview') as HTMLIFrameElement;
		const css = frame.contentDocument!.querySelector('style[data-ds-tokens]')!.textContent!;
		const [light, dark] = css.split("[data-ds-theme='dark']");
		return { light, dark };
	}
	// Answers the save of the hash first read with a conflict, and any other save with success.
	function conflictFirst(diff?: string) {
		const conflict = { error: 'the token file changed', current: CURRENT, hash: 'h2', diff };
		return (body: Record<string, unknown>) =>
			body.hash === 'h1' ? answer(conflict, 409) : saveRaw();
	}
	const lightPrimary = editKey('light', 'color.primary');

	it('enters Edit mode with an editor per namespace, in the token sheet order, and reads the raw file', async () => {
		await edit();
		expect(q('design-edit')).not.toBeNull();
		expect(gets).toBe(1);
		const editors = [...document.querySelectorAll<HTMLElement>('[data-testid="token-editor"]')];
		expect(editors.map((e) => e.dataset.namespace)).toEqual([
			'color',
			'font',
			'text',
			'spacing',
			'radius',
			'shadow'
		]);
		expect(button('Save').disabled).toBe(true);
		expect(button('Discard changes').disabled).toBe(true);
	});

	it('moves the preview with each change before saving, a light edit reaching dark where dark does not override it', async () => {
		await edit();
		expect(preview().light).not.toContain('#336699');
		type('color.primary, light', '#336699');
		await settle();
		expect(preview().light).toContain('--ds-color-primary: #336699;');
		// dark overrides color.primary, so its own value stays
		expect(preview().dark).not.toContain('#336699');
		// nothing in dark overrides radius.md, so the light edit reaches it
		type('radius.md, light and dark', '0.5');
		await settle();
		expect(preview().light).toContain('--ds-radius-md: 0.5rem;');
		expect(preview().dark).toContain('--ds-radius-md: 0.5rem;');
		expect(words(q('design-edit'))).toContain('unsaved changes');
		expect(puts).toEqual([]);
	});

	it('flags a failing contrast beside the token and in the token sheet after an edit', async () => {
		await edit();
		const muted = () => q('edit-color.text-muted')!.querySelector('[data-theme="light"]')!;
		expect(muted().querySelector('[data-testid="contrast-fail"]')).toBeNull();
		type('color.text-muted, light', '#bbbbbb');
		await settle();
		expect(words(muted().querySelector('[data-testid="contrast-fail"]'))).toContain(
			'Fails 4.5:1 on color.background in light'
		);
		const pair = q('pair-color.text-muted-on-color.background')!;
		const verdict = (theme: string) =>
			words(pair.querySelector(`[data-sheet-theme="${theme}"] [data-testid="verdict"]`));
		expect(verdict('light')).toBe('fails 4.5');
		expect(verdict('dark')).toBe('passes 4.5');
	});

	it('saves the edits composed into the raw token file with its hash, then takes the saved values as the base', async () => {
		await edit();
		type('color.primary, light', '#336699');
		type('What changed', 'A bluer primary');
		await settle();
		const reads = api.mock.calls.filter(([url]) => url === '/api/design').length;
		button('Save').click();
		await settle();
		expect(puts).toEqual([
			{
				path: TOKENS,
				content: composeTokenFile(RAW, { [lightPrimary]: PRIMARY }),
				hash: 'h1',
				message: 'A bluer primary'
			}
		]);
		expect(puts[0].content).toContain(
			'"$value": {"colorSpace":"srgb","components":[0.2,0.4,0.6],"hex":"#336699"}'
		);
		expect(api.mock.calls.filter(([url]) => url === '/api/design').length).toBe(reads + 1);
		expect(gets).toBe(2);
		expect(words(q('design-notice'))).toBe(`Saved ${TOKENS}.`);
		expect(changedMarks()).toEqual([]);
		expect(field('What changed').value).toBe('');
		expect(button('Save').disabled).toBe(true);
	});

	it("shows flai's refusal and each finding, keeping the edits", async () => {
		const refusal = {
			error: 'flai refused the token file',
			findings: [
				{
					level: 'error',
					rule: 'design.tokens',
					path: TOKENS,
					line: 9,
					message: 'color.text-muted fails 4.5:1'
				}
			]
		};
		put = () => answer(refusal, 422);
		await edit();
		type('color.primary, light', '#336699');
		button('Save').click();
		await settle();
		const alert = q('design-save-error')!;
		expect(alert.getAttribute('role')).toBe('alert');
		expect(words(alert)).toContain('flai refused the token file');
		const finding = `${TOKENS}:9: error: design.tokens: color.text-muted fails 4.5:1`;
		expect(words(alert)).toContain(finding);
		expect(changedMarks()).toEqual(['edit-color.primary']);
		expect(field('color.primary, light').value).toBe('#336699');
		expect(preview().light).toContain('--ds-color-primary: #336699;');
		expect(button('Save').disabled).toBe(false);
	});

	it('shows a conflict with its diff, and saves mine over it with the hash the conflict gives', async () => {
		put = conflictFirst('-a\n+b');
		await edit();
		type('color.primary, light', '#336699');
		button('Save').click();
		await settle();
		const alert = q('design-save-error')!;
		expect(words(alert)).toContain('the token file changed');
		expect(alert.querySelector('pre')!.textContent).toBe('-a\n+b');
		expect(changedMarks()).toEqual(['edit-color.primary']);

		button('Save mine over it').click();
		await settle();
		expect(puts.map((p) => p.hash)).toEqual(['h1', 'h2']);
		expect(puts[1].content).toBe(puts[0].content);
		expect(q('design-save-error')).toBeNull();
		expect(changedMarks()).toEqual([]);
	});

	it('loads the current version on a conflict, keeping the edits that still have a place in it', async () => {
		put = conflictFirst();
		await edit();
		type('color.primary, light', '#336699');
		type('font.mono, light and dark', MONO.join(', '));
		button('Save').click();
		await settle();
		button('Load the current version').click();
		await settle();
		expect(words(q('design-notice'))).toBe(
			'Loaded the current version. Your edits are kept where its tokens still take them; 1 had no place in it and was discarded.'
		);
		expect(changedMarks()).toEqual(['edit-color.primary']);
		expect(q('design-save-error')).toBeNull();

		button('Save').click();
		await settle();
		expect(puts[1]).toMatchObject({
			hash: 'h2',
			content: composeTokenFile(CURRENT, { [lightPrimary]: PRIMARY })
		});
	});

	it('edits with the preview following when flai cannot read the token file, saying it cannot save', async () => {
		get = () => answer({ error: 'flai on the host is older than this dashboard' }, 501);
		await edit();
		expect(words(q('design-cannot-save'))).toBe(CANNOT_SAVE);
		type('color.primary, light', '#336699');
		await settle();
		expect(preview().light).toContain('--ds-color-primary: #336699;');
		expect(changedMarks()).toEqual(['edit-color.primary']);
		expect(button('Save').disabled).toBe(true);
		expect(puts).toEqual([]);
	});

	it('keeps the edits and says it cannot save when the save answers 501', async () => {
		put = () => answer({ error: 'flai on the host does not offer file.save' }, 501);
		await edit();
		expect(q('design-cannot-save')).toBeNull();
		type('color.primary, light', '#336699');
		button('Save').click();
		await settle();
		expect(words(q('design-cannot-save'))).toBe(CANNOT_SAVE);
		expect(changedMarks()).toEqual(['edit-color.primary']);
		expect(field('color.primary, light').value).toBe('#336699');
		expect(button('Save').disabled).toBe(true);
	});

	it('discards the edits, and text that made none, back to the saved tokens', async () => {
		await edit();
		const saved = field('color.primary, light').value;
		type('color.primary, light', '#336699');
		type('color.text, light', 'not a colour');
		await settle();
		button('Discard changes').click();
		await settle();
		expect(changedMarks()).toEqual([]);
		expect(field('color.primary, light').value).toBe(saved);
		expect(field('color.text, light').hasAttribute('aria-invalid')).toBe(false);
		expect(preview().light).not.toContain('#336699');
		expect(words(q('design-edit'))).not.toContain('unsaved changes');
		expect(puts).toEqual([]);
	});

	it('keeps the edits when a document change reloads the design system, refreshing the saved tokens only', async () => {
		await edit();
		type('color.primary, light', '#336699');
		const surface = design.themes[0].tokens.find((t) => t.path === 'color.surface')!;
		surface.value = { colorSpace: 'srgb', components: [0.933, 0.933, 0.933], hex: '#eeeeee' };
		surface.css = '#eeeeee';
		surface.variables = [{ name: '--color-surface', value: '#eeeeee' }];
		followed.of.document();
		await settle();
		expect(field('color.surface, light').value).toBe('#eeeeee');
		expect(changedMarks()).toEqual(['edit-color.primary']);
		expect(preview().light).toContain('--ds-color-primary: #336699;');
		expect(gets).toBe(1);
	});

	it('asks before leaving with unsaved changes, in the app and out of it', async () => {
		await edit();
		const cancel = vi.fn();
		const unload = () => {
			const e = new Event('beforeunload', { cancelable: true });
			window.dispatchEvent(e);
			return e.defaultPrevented;
		};
		navigation.before!({ cancel });
		expect(confirm).not.toHaveBeenCalled();
		expect(unload()).toBe(false);

		type('color.primary, light', '#336699');
		navigation.before!({ cancel });
		expect(confirm).toHaveBeenCalledWith('Leave without saving your changes?');
		expect(cancel).toHaveBeenCalledTimes(1);
		expect(unload()).toBe(true);

		confirm.mockReturnValueOnce(true);
		navigation.before!({ cancel });
		expect(cancel).toHaveBeenCalledTimes(1);
	});
});

const STYLESHEET = 'flaiover/src/routes/layout.css';
const BLUER = { id: 'S-0390', title: 'A bluer primary', status: 'in-progress', paths: [TOKENS] };
const ACCENT = {
	id: 'S-0391',
	title: 'Say when to use the accent',
	status: 'ready',
	paths: [TOKENS, GUIDANCE]
};
const EPIC_LIST = [
	{ id: 'E-0001', title: 'Dashboard', status: 'in-progress' },
	{ id: 'E-0002', title: 'Shipped', status: 'done' },
	{ id: 'E-0003', title: 'Dropped', status: 'cancelled' },
	{ id: 'E-0004', title: 'Design system', status: 'backlog' }
];
const STORY = '/api/design?story=';

describe('the proposals and the request for a change on the design page (S-0384)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let main: DesignSystem;
	let list: () => unknown;
	let story: (id: string) => unknown;

	// S-0390's design system: the main branch's with a bluer light primary.
	const proposed = () => {
		const d = structuredClone(withSystem);
		d.problems = [];
		d.themes[0].tokens[1] = colour('primary', '#336699');
		return d;
	};
	const withStylesheet = () => ({ ...structuredClone(withSystem), stylesheet: STYLESHEET });

	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		followed.kinds = [];
		followed.of = {};
		document.body.innerHTML = '';
	});

	// Serves design.get, the guidance, the proposals, the epics, and a proposal's design.get.
	async function show(proposals: unknown[] = [BLUER, ACCENT], design = withStylesheet()) {
		main = design;
		list = () => answer(proposals);
		story = (id) => (id === 'S-0390' ? answer(proposed()) : answer({ error: `no ${id}` }, 404));
		api.mockImplementation(async (url: string) => {
			if (url === '/api/design') return answer(main);
			if (url === guidanceUrl) return answer(guidanceDoc);
			if (url === PROPOSALS) return list();
			if (url === EPICS) return answer(EPIC_LIST);
			if (url.startsWith(STORY)) return story(decodeURIComponent(url.slice(STORY.length)));
			return answer({ error: `unexpected ${url}` }, 404);
		});
		c = mount(DesignPage, { target: document.body });
		await settle();
	}

	const q = (testid: string) => document.querySelector<HTMLElement>(`[data-testid="${testid}"]`);
	const words = (e: Element | null) => e!.textContent!.replace(/\s+/g, ' ').trim();
	const compare = (id: string) => q(`proposal-${id}`)!.querySelector('button')!;
	const pathsOf = (row: Element) =>
		[...row.querySelectorAll('[data-testid="proposal-path"]')].map(words);
	const storyReads = () =>
		api.mock.calls.map(([url]) => String(url)).filter((url) => url.startsWith(STORY));
	function button(text: string) {
		const found = [...document.querySelectorAll('button')].find(
			(b) => b.textContent!.trim() === text
		);
		if (!found) throw new Error(`no button ${text}`);
		return found;
	}

	it('lists each proposal with its link, title, state, and which of the token file and guidance it changes', async () => {
		await show();
		const bluer = q('proposal-S-0390')!;
		expect(bluer.querySelector('a')!.getAttribute('href')).toBe('/items/S-0390');
		expect(words(bluer)).toContain('A bluer primary');
		expect(words(bluer)).toContain('in-progress');
		expect(pathsOf(bluer)).toEqual([`token file ${TOKENS}`]);
		const accent = q('proposal-S-0391')!;
		expect(accent.querySelector('a')!.getAttribute('href')).toBe('/items/S-0391');
		expect(words(accent)).toContain('ready');
		expect(pathsOf(accent)).toEqual([`token file ${TOKENS}`, `guidance ${GUIDANCE}`]);
		expect(q('proposals-none')).toBeNull();
		expect(q('proposal-chosen')).toBeNull();
	});

	it('says no open story proposes a change, pointing to the form, when the list is empty', async () => {
		await show([]);
		expect(words(q('proposals-none'))).toBe(
			'No open story proposes a change to the design system. Request one below.'
		);
		expect(q('propose-change')).not.toBeNull();
	});

	it("says the host's flai is too old to list the proposals when it answers 501", async () => {
		await show();
		list = () => answer({ error: 'flai on the host does not offer design.proposals' }, 501);
		followed.of.item();
		await settle();
		expect(words(q('proposals-cannot-list'))).toBe(
			'This flai cannot list the proposals yet; upgrade flai on the host.'
		);
		expect(q('proposal-S-0390')).toBeNull();
	});

	it("compares the chosen proposal's tokens with the main branch's, and closes it when chosen again", async () => {
		await show();
		compare('S-0390').click();
		await settle();
		expect(storyReads()).toEqual([`${STORY}S-0390`]);
		expect(compare('S-0390').getAttribute('aria-pressed')).toBe('true');
		const heading = q('proposal-compare')!.querySelector('h3');
		expect(words(heading)).toBe('Proposal S-0390: A bluer primary');
		const row = q('difference-light-color.primary')!;
		expect(row.dataset.change).toBe('changed');
		expect(words(row.querySelector('[data-side="main"]'))).toBe('#2563eb');
		expect(words(row.querySelector('[data-side="proposal"]'))).toBe('#336699');
		expect(q('difference-dark-color.primary')).toBeNull();

		compare('S-0390').click();
		await settle();
		expect(q('proposal-chosen')).toBeNull();
		expect(compare('S-0390').getAttribute('aria-pressed')).toBe('false');
	});

	it('closes the comparison with its close control', async () => {
		await show();
		compare('S-0390').click();
		await settle();
		button('Close comparison').click();
		await settle();
		expect(q('proposal-chosen')).toBeNull();
	});

	it("shows flai's message when it refuses to read a proposal's design system", async () => {
		await show();
		story = () => answer({ error: 'S-0391 has no worktree: start it again' }, 400);
		compare('S-0391').click();
		await settle();
		expect(storyReads()).toEqual([`${STORY}S-0391`]);
		const refused = q('proposal-error')!;
		expect(refused.getAttribute('role')).toBe('alert');
		expect(words(refused)).toBe(
			"Could not read S-0391's design system: S-0391 has no worktree: start it again"
		);
		expect(q('proposal-compare')).toBeNull();
	});

	it('closes the comparison when its proposal leaves the list, and keeps it while it stays', async () => {
		await show();
		compare('S-0390').click();
		await settle();
		list = () => answer([BLUER]);
		followed.of.item();
		await settle();
		expect(q('proposal-compare')).not.toBeNull();
		expect(q('proposal-S-0391')).toBeNull();

		list = () => answer([ACCENT]);
		followed.of.item();
		await settle();
		expect(q('proposal-S-0390')).toBeNull();
		expect(q('proposal-chosen')).toBeNull();
	});

	it('reads the proposals and the chosen proposal again when a document changes', async () => {
		await show();
		compare('S-0390').click();
		await settle();
		const lists = () => api.mock.calls.filter(([url]) => url === PROPOSALS).length;
		const before = lists();
		followed.of.document();
		await settle();
		expect(lists()).toBe(before + 1);
		expect(storyReads()).toEqual([`${STORY}S-0390`, `${STORY}S-0390`]);
		expect(q('proposal-compare')).not.toBeNull();
	});

	it('offers the request form with the open epics, the token file, and the stylesheet', async () => {
		await show();
		const request = q('design-request')!;
		expect(words(request.querySelector('h2'))).toBe('Request a change');
		expect(request.querySelector('[data-testid="propose-change"]')).not.toBeNull();
		const offered = [...q('propose-epic')!.querySelectorAll('option')].map((o) => o.value);
		expect(offered).toEqual(['', 'E-0001', 'E-0004']);
		expect(words(q('propose-about'))).toContain(`${TOKENS} and ${STYLESHEET}`);
	});

	it('closes the comparison and holds compare back while the tokens are edited', async () => {
		await show();
		compare('S-0390').click();
		await settle();
		button('Edit tokens').click();
		await settle();
		expect(q('design-edit')).not.toBeNull();
		expect(q('proposal-chosen')).toBeNull();
		expect(compare('S-0390').disabled).toBe(true);
	});

	it('offers neither the proposals nor the request form to a project with no token file', async () => {
		await show([BLUER], without);
		expect(q('no-design-system')).not.toBeNull();
		expect(q('design-proposals')).toBeNull();
		expect(q('design-request')).toBeNull();
		expect(api).not.toHaveBeenCalledWith(PROPOSALS);
		expect(api).not.toHaveBeenCalledWith(EPICS);
	});
});
