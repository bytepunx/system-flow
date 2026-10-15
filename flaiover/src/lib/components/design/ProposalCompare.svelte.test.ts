// S-0384: a proposal for the design system previewed on the gallery beside the main branch's, with
// each token that differs marked and its two values shown.
import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ProposalCompare from './ProposalCompare.svelte';
import type { DesignSystem, DesignToken } from '$lib/server/repo';

const colour = (name: string, hex: string): DesignToken => ({
	path: `color.${name}`,
	type: 'color',
	value: { colorSpace: 'srgb', hex },
	css: hex,
	variables: [{ name: `--color-${name}`, value: hex }],
	line: 1
});
const radius: DesignToken = {
	path: 'radius.md',
	type: 'dimension',
	value: { unit: 'rem', value: 0.375 },
	css: '0.375rem',
	variables: [{ name: '--radius-md', value: '0.375rem' }],
	line: 1
};

const system = (light: DesignToken[], dark: DesignToken[]): DesignSystem => ({
	path: 'design/system/design-system.tokens.json',
	exists: true,
	guidance: 'design/system/design-system.md',
	headings: [],
	resolver: false,
	themes: [
		{ name: 'light', tokens: light },
		{ name: 'dark', tokens: dark }
	],
	problems: [],
	css: '',
	tailwind: '',
	contrast: []
});

const MAIN = system(
	[colour('background', '#ffffff'), colour('primary', '#2563eb'), radius],
	[colour('background', '#0f172a'), colour('primary', '#60a5fa'), radius]
);
// Light's primary changed, its radius removed, an accent added; dark unchanged.
const PROPOSAL = system(
	[colour('background', '#ffffff'), colour('primary', '#dc2626'), colour('accent', '#9333ea')],
	MAIN.themes[1].tokens
);

describe('ProposalCompare (S-0384)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	function render(proposal = PROPOSAL) {
		c = mount(ProposalCompare, {
			target: document.body,
			props: { main: MAIN, proposal, story: 'S-0400', title: 'A warmer primary' }
		});
		flushSync();
	}
	const q = <T extends Element = HTMLElement>(sel: string) => document.querySelector<T>(sel)!;
	const words = (el: Element) => el.textContent!.replace(/\s+/g, ' ').trim();
	const row = (theme: string, path: string) => q(`[data-testid="difference-${theme}-${path}"]`);
	const side = (r: Element, s: string) => r.querySelector(`[data-side="${s}"]`)!;

	it("previews the gallery with main's tokens and the proposal's, side by side, each labelled", () => {
		render();
		expect(words(q('h2'))).toBe('Proposal S-0400: A warmer primary');
		const main = q('[data-testid="compare-main"]');
		const proposal = q('[data-testid="compare-proposal"]');
		expect(words(main.querySelector('h3')!)).toBe('Main branch');
		expect(words(proposal.querySelector('h3')!)).toBe('Proposal S-0400');
		const frame = (el: Element) => el.querySelector<HTMLIFrameElement>('iframe')!;
		expect(frame(main).title).toBe("The main branch's design system preview");
		expect(frame(proposal).title).toBe("S-0400's proposed design system preview");
		const tokens = (el: Element) =>
			frame(el).contentDocument!.querySelector('style[data-ds-tokens]')!.textContent!;
		expect(tokens(main)).toContain('--ds-color-primary: #2563eb;');
		expect(tokens(proposal)).toContain('--ds-color-primary: #dc2626;');
		expect(tokens(proposal)).not.toContain('#2563eb');
	});

	it("marks each token that differs, per theme, with main's value and the proposal's", () => {
		render();
		const differences = q('[data-testid="proposal-differences"]');
		expect([...differences.querySelectorAll('h4')].map(words)).toEqual(['Light']);
		const rows = [...differences.querySelectorAll<HTMLElement>('tbody tr')];
		expect(rows.map((r) => [r.dataset.change, r.querySelector('td')!.textContent])).toEqual([
			['changed', 'color.primary'],
			['removed', 'radius.md'],
			['added', 'color.accent']
		]);

		const primary = row('light', 'color.primary');
		expect(words(primary.querySelector('[data-testid="change"]')!)).toBe('changed');
		expect(words(side(primary, 'main'))).toBe('#2563eb');
		expect(words(side(primary, 'proposal'))).toBe('#dc2626');
		const swatch = (r: Element, s: string) =>
			side(r, s).querySelector<HTMLElement>('[data-testid="swatch"]')!;
		expect(swatch(primary, 'main').style.backgroundColor).not.toBe('');
		expect(swatch(primary, 'proposal').style.backgroundColor).not.toBe(
			swatch(primary, 'main').style.backgroundColor
		);

		const removed = row('light', 'radius.md');
		expect(words(side(removed, 'main'))).toBe('0.375rem');
		expect(side(removed, 'main').querySelector('[data-testid="swatch"]')).toBeNull();
		expect(words(side(removed, 'proposal'))).toBe('not in the proposal');
		const added = row('light', 'color.accent');
		expect(words(side(added, 'main'))).toBe('not on the main branch');
		expect(words(side(added, 'proposal'))).toBe('#9333ea');
		// The table is named by its theme's heading, and its columns by their headers.
		const table = primary.closest('table')!;
		expect(document.getElementById(table.getAttribute('aria-labelledby')!)!.textContent).toContain(
			'Light'
		);
		expect([...table.querySelectorAll('th')].map(words)).toEqual([
			'Token',
			'Change',
			'Main branch',
			'Proposal'
		]);
		expect(document.querySelector('[data-testid="proposal-no-difference"]')).toBeNull();
	});

	it('says a proposal that changes no token does so, its guidance aside', () => {
		render(system(MAIN.themes[0].tokens, MAIN.themes[1].tokens));
		expect(words(q('[data-testid="proposal-no-difference"]'))).toBe(
			'S-0400 changes no token; its guidance may still differ.'
		);
		expect(document.querySelector('[data-testid="proposal-differences"] table')).toBeNull();
		expect(document.querySelectorAll('iframe').length).toBe(2);
	});

	it("shows the proposal's problems, and lists nothing removed when it has no tokens to compare", () => {
		render({
			...system([], []),
			themes: [],
			problems: [{ line: 3, message: 'unexpected end of JSON input' }]
		});
		expect(words(q('[data-testid="proposal-problems"]'))).toContain(
			'line 3: unexpected end of JSON input'
		);
		expect(words(q('[data-testid="proposal-no-tokens"]'))).toBe(
			"S-0400's token file does not parse, so it has no tokens to compare."
		);
		expect(document.querySelector('[data-testid="proposal-differences"] tbody')).toBeNull();
	});
});
