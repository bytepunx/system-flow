// S-0382: the design page shows the project's design system, its token sheet, its guidance, and
// the live preview; tells a project with no token file so, with the command that writes one; and
// reads it all again when a document changes, as the other pages follow theirs.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { DesignSystem, DesignToken } from '$lib/server/repo';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const enhance = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$lib/markdown', async (actual) => ({
	...(await actual<typeof import('$lib/markdown')>()),
	enhance
}));
const followed = vi.hoisted(() => ({
	kinds: [] as string[][],
	f: null as null | (() => void)
}));
vi.mock('$lib/events', () => ({
	listen: () => () => {},
	follow: (kinds: string[], f: () => void) => {
		followed.kinds.push(kinds);
		followed.f = f;
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

describe('the design page (S-0382)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		followed.kinds = [];
		followed.f = null;
		document.body.innerHTML = '';
	});

	function serve(design: () => unknown, guidance: () => unknown = () => answer(guidanceDoc)) {
		api.mockImplementation(async (url: string) => {
			if (url === '/api/design') return design();
			if (url === guidanceUrl) return guidance();
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
		expect(followed.kinds).toEqual([['document']]);
		expect(q('no-design-system')).not.toBeNull();

		current = withSystem;
		followed.f!();
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
