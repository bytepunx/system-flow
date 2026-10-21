// S-0383: the token editor changes each kind of token in place, bound to the edit model, and marks
// a value flai would not read and each failing contrast pair beside the token, in its theme.
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import TokenEditor from './TokenEditor.svelte';
import type { DesignContrastPair, DesignTheme } from '$lib/server/repo';
import { editedThemes, editKey, type Edits } from '$lib/design/tokens';
import { pairsOf, rolePairs } from '$lib/design/preview';

// flai's design.get of its test token file, as `flai design show --json` prints it.
const GOLDEN = resolve('../flai/internal/designsystem/testdata/show.json.golden');
const SAVED = (JSON.parse(readFileSync(GOLDEN, 'utf8')) as { themes: DesignTheme[] }).themes;

// A plain token file: no theme context overrides a token, so a light edit carries into dark.
const PLAIN = '{}';

// Both themes' role pairs on the themes given, as the page computes them for the preview.
const pairsOn = (themes: readonly DesignTheme[]): DesignContrastPair[] =>
	themes.flatMap((t) => rolePairs(t.name, t.tokens));

describe('TokenEditor (S-0383)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	// Mounts the editor on namespace with the page's part played here: edits bound, and edited and
	// pairs recomputed from them on each edit.
	function render(namespace: string) {
		const props = $state({
			namespace,
			saved: SAVED,
			edited: SAVED,
			edits: {} as Edits,
			pairs: pairsOn(SAVED),
			onedit: (edits: Edits) => {
				props.edited = editedThemes(SAVED, edits, PLAIN);
				props.pairs = pairsOn(props.edited);
			}
		});
		c = mount(TokenEditor, { target: document.body, props });
		flushSync();
		return props;
	}

	const el = <T extends HTMLElement = HTMLInputElement>(label: string) => {
		const found = document.querySelector<T>(`[aria-label="${label}"]`);
		if (!found) throw new Error(`no field labelled ${label}`);
		return found;
	};
	function type(label: string, text: string) {
		const input = el(label);
		input.value = text;
		input.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
	}
	function commit(label: string) {
		el(label).dispatchEvent(new Event('change', { bubbles: true }));
		flushSync();
	}
	function choose(label: string, value: string) {
		const select = el<HTMLSelectElement>(label);
		select.value = value;
		select.dispatchEvent(new Event('change', { bubbles: true }));
		flushSync();
	}
	const row = (path: string) =>
		document.querySelector<HTMLElement>(`[data-testid="edit-${path}"]`)!;
	const cell = (path: string, theme: string) =>
		row(path).querySelector<HTMLElement>(`[data-cell-theme="${theme}"]`)!;
	const words = (e: Element) => e.textContent!.replace(/\s+/g, ' ').trim();
	const fails = (path: string, theme: string) =>
		[...cell(path, theme).querySelectorAll('[data-testid="contrast-fail"]')].map(words);
	const note = (path: string) => words(row(path).querySelector('[data-testid="readonly"]')!);
	const alertFor = (label: string) => {
		const id = el(label).getAttribute('aria-describedby')!.split(' ')[0];
		return document.getElementById(id)!;
	};

	it('changes a colour by its text field, as #rrggbb or oklch, in the theme it names', () => {
		const props = render('color');
		// An oklch colour shows its text, and its picker its sRGB hex.
		expect(el('color.primary, light').value).toBe('oklch(52% 0.18 255)');
		expect(el('color.primary, light, picker').value).toMatch(/^#[0-9a-f]{6}$/);
		expect(el('color.background, light').value).toBe('#ffffff');

		type('color.primary, light', '#336699');
		expect(props.edits).toEqual({
			[editKey('light', 'color.primary')]: {
				colorSpace: 'srgb',
				components: [0.2, 0.4, 0.6],
				hex: '#336699'
			}
		});
		type('color.primary, dark', 'oklch(60% 0.1 200)');
		expect(props.edits[editKey('dark', 'color.primary')]).toEqual({
			colorSpace: 'oklch',
			components: [0.6, 0.1, 200]
		});
		expect(row('color.primary').querySelector('[data-testid="changed"]')).not.toBeNull();
		// The swatch follows the edited value.
		const swatch = row('color.primary').querySelector<HTMLElement>('[data-testid="swatch"]')!;
		expect(swatch.style.backgroundColor).toBe('rgb(51, 102, 153)');
		// Once the field is left, it shows the value as the edit model writes it.
		commit('color.primary, dark');
		expect(el('color.primary, dark').value).toBe('oklch(60% 0.1 200)');
	});

	it('changes a colour by its picker into srgb, keeping its alpha', () => {
		const props = render('color');
		type('color.primary, light, picker', '#ff0000');
		expect(props.edits[editKey('light', 'color.primary')]).toEqual({
			colorSpace: 'srgb',
			components: [1, 0, 0],
			hex: '#ff0000'
		});
		expect(el('color.primary, light').value).toBe('#ff0000');

		type('color.scrim, dark, picker', '#102030');
		expect(props.edits[editKey('dark', 'color.scrim')]).toMatchObject({
			colorSpace: 'srgb',
			hex: '#102030',
			alpha: 0.4
		});
	});

	it('changes a length by its number and unit, one field for both themes where they agree', () => {
		const props = render('radius');
		expect(el('radius.md, light and dark').value).toBe('0.375');
		type('radius.md, light and dark', '0.5');
		expect(props.edits).toEqual({ [editKey('light', 'radius.md')]: { unit: 'rem', value: 0.5 } });
		choose('radius.md, light and dark, unit', 'px');
		expect(props.edits).toEqual({ [editKey('light', 'radius.md')]: { unit: 'px', value: 0.5 } });
		// The dark theme took the light edit, so one field still serves both.
		expect(cell('radius.md', 'light-and-dark')).not.toBeNull();
		expect(el<HTMLSelectElement>('radius.md, light and dark, unit').value).toBe('px');

		unmount(c!);
		const spacing = render('spacing');
		type('spacing, light and dark', '0.3');
		expect(spacing.edits).toEqual({ [editKey('light', 'spacing')]: { unit: 'rem', value: 0.3 } });
	});

	it("changes a type step's font size, its unit, and its line height", () => {
		const props = render('text');
		const key = editKey('light', 'text.2xl');
		type('text.2xl, light and dark, font size', '2');
		expect(props.edits[key]).toEqual({ fontSize: { unit: 'rem', value: 2 }, lineHeight: 1.333 });
		choose('text.2xl, light and dark, font size unit', 'px');
		type('text.2xl, light and dark, line height', '1.25');
		expect(props.edits[key]).toEqual({ fontSize: { unit: 'px', value: 2 }, lineHeight: 1.25 });

		// A line height given as a length has its own unit; a step that is a plain length, its number.
		choose('text.sm, light and dark, line height unit', 'px');
		expect(props.edits[editKey('light', 'text.sm')]).toEqual({
			fontSize: { unit: 'rem', value: 0.875 },
			lineHeight: { unit: 'px', value: 1.25 }
		});
		type('text.lg, light and dark', '1.2');
		expect(props.edits[editKey('light', 'text.lg')]).toEqual({ unit: 'rem', value: 1.2 });
	});

	it('changes a font family from its stack, quotes stripped and blanks dropped', () => {
		const props = render('font');
		expect(el('font.mono, light and dark').value).toBe('JetBrains Mono, ui-monospace, monospace');
		type('font.mono, light and dark', `"Fira Code", , 'IBM Plex Mono',monospace`);
		expect(props.edits).toEqual({
			[editKey('light', 'font.mono')]: ['Fira Code', 'IBM Plex Mono', 'monospace']
		});
	});

	it('keeps a value flai would not read in its field, marked, and makes no edit of it', () => {
		const props = render('color');
		type('color.primary, light', 'blue');
		const field = el('color.primary, light');
		expect(field.value).toBe('blue');
		expect(field.getAttribute('aria-invalid')).toBe('true');
		expect(alertFor('color.primary, light').getAttribute('role')).toBe('alert');
		expect(words(alertFor('color.primary, light'))).toContain('write a colour as #rrggbb');
		expect(props.edits).toEqual({});
		// It stays when the field is left.
		commit('color.primary, light');
		expect(el('color.primary, light').value).toBe('blue');

		// A colour that parses but that flai refuses is marked in flai's words.
		type('color.primary, light', 'oklch(150% 0.1 200)');
		expect(words(alertFor('color.primary, light'))).toContain('oklch lightness is 1.5');
		expect(props.edits).toEqual({});
		// A good value clears the mark.
		type('color.primary, light', '#123456');
		expect(el('color.primary, light').hasAttribute('aria-invalid')).toBe(false);
		expect(Object.keys(props.edits)).toEqual([editKey('light', 'color.primary')]);
	});

	it('marks a bad number, a line height of zero, and an empty font stack', () => {
		const props = render('text');
		type('text.xs, light and dark, font size', 'big');
		expect(words(alertFor('text.xs, light and dark, font size'))).toContain('write a number');
		type('text.xs, light and dark, line height', '0');
		expect(words(alertFor('text.xs, light and dark, line height'))).toContain('lineHeight is 0');
		expect(props.edits).toEqual({});

		unmount(c!);
		const font = render('font');
		type('font.sans, light and dark', ' , ');
		expect(words(alertFor('font.sans, light and dark'))).toContain('a font family is a name');
		expect(font.edits).toEqual({});
	});

	it('marks each failing contrast pair a colour is in, with its ratio, in its theme', () => {
		const props = render('color');
		// The saved border fails 3:1 on the ground in light: both tokens say so, in light only.
		const border = fails('color.border', 'light');
		expect(border).toEqual(['Fails 3:1 on color.background in light: 1.59:1']);
		const ground = fails('color.background', 'light');
		expect(ground).toEqual(['Fails 3:1 under color.border in light: 1.59:1']);
		const dark = fails('color.border', 'dark');
		expect(dark).toEqual(['Fails 3:1 on color.background in dark: 1.67:1']);
		expect(fails('color.text-muted', 'light')).toEqual([]);

		// Muted text made too light fails 4.5:1 on the ground and the surface, in light and not dark.
		type('color.text-muted, light', '#bbbbbb');
		const failing = pairsOf('color.text-muted', props.pairs).filter((p) => p.theme === 'light');
		expect(failing.map((p) => p.background)).toEqual(['color.background', 'color.surface']);
		expect(fails('color.text-muted', 'light')).toEqual(
			failing.map((p) => `Fails 4.5:1 on ${p.background} in light: ${p.ratio.toFixed(2)}:1`)
		);
		expect(fails('color.text-muted', 'dark')).toEqual([]);
		// The field is described by its marks.
		const described = el('color.text-muted, light').getAttribute('aria-describedby')!;
		expect(words(document.getElementById(described)!)).toContain('Fails 4.5:1');
	});

	it('leaves no edit of a value set back to its saved one, and Reset drops the edits', () => {
		const props = render('color');
		type('color.primary, light', '#000000');
		type('color.primary, light', 'oklch(52% 0.18 255)');
		expect(props.edits).toEqual({});
		expect(row('color.primary').querySelector('[data-testid="changed"]')).toBeNull();

		type('color.primary, light', '#111111');
		type('color.primary, dark', '#222222');
		type('color.text, light', '#333333');
		type('color.primary, dark', 'not a colour');
		el<HTMLButtonElement>('Reset color.primary').click();
		flushSync();
		expect(props.edits).toEqual({ [editKey('light', 'color.text')]: expect.anything() });
		expect(row('color.primary').querySelector('[data-testid="changed"]')).toBeNull();
		expect(el('color.primary, light').value).toBe('oklch(52% 0.18 255)');
		expect(el('color.primary, dark').value).toBe('oklch(72% 0.15 255)');
		expect(el('color.primary, dark').hasAttribute('aria-invalid')).toBe(false);
		expect(row('color.text').querySelector('[data-testid="changed"]')).not.toBeNull();
	});

	it('shows shadows and aliases read-only, saying where they are changed', () => {
		render('shadow');
		for (const path of ['shadow.sm', 'shadow.md', 'shadow.lg', 'shadow.inner']) {
			expect(note(path)).toBe('A shadow token is edited in the token file.');
			expect(row(path).querySelector('input, select')).toBeNull();
		}
		expect(words(row('shadow.md'))).toContain('oklch(0% 0 none / 0.5)');

		unmount(c!);
		render('color');
		expect(note('color.link')).toContain('Follows color.primary');
		expect(row('color.link').querySelector('input')).toBeNull();
	});

	it('labels every field with its token path and theme', () => {
		render('color');
		const unlabelled = [...document.querySelectorAll('input, select')].filter((f) => {
			const label = f.getAttribute('aria-label') ?? '';
			return !/^color\.[a-z-]+, (light|dark)(, picker)?$/.test(label);
		});
		expect(unlabelled).toEqual([]);
	});

	// I-0159: layout.css scopes the dark palette and Tailwind's dark variant to any element with
	// data-theme="dark", so a cell carrying it drew the other theme's fields unreadable.
	it('sets data-theme on no element, so the page theme styles every cell (I-0159)', () => {
		for (const namespace of ['color', 'font', 'radius', 'shadow', 'text']) {
			render(namespace);
			const editor = document.querySelector('[data-testid="token-editor"]')!;
			expect(editor.querySelectorAll('[data-theme]'), namespace).toHaveLength(0);
			if (c) unmount(c);
			c = undefined;
			document.body.innerHTML = '';
		}
		render('color');
		expect(cell('color.primary', 'dark').querySelector('input')).not.toBeNull();
	});
});
