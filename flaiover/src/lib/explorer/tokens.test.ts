import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import type { DesignSystem, DesignTheme } from '$lib/server/repo';
import { editedThemes, editKey, parseColorText } from '$lib/design/tokens';
import {
	DARK_SELECTOR,
	LIGHT_SELECTOR,
	stylesheetDeclarations,
	type StylesheetDeclarations
} from './tokens';

// flai's own design system fixtures, so that the port cannot drift from flai design css.
const DIR = new URL('../../../../flai/internal/designsystem/testdata/', import.meta.url);
const read = (name: string) => readFileSync(new URL(name, DIR), 'utf8');
const RAW = read('themes.tokens.json');
// What flai design show --json answers for themes.tokens.json: its themes, as design.get has them.
const SAVED: DesignTheme[] = (JSON.parse(read('show.json.golden')) as DesignSystem).themes;

// Rules by selector, each its custom properties by name.
type Rules = Map<string, Record<string, string>>;

// Each rule of a stylesheet.
function rules(css: string): Rules {
	const out: Rules = new Map();
	for (const [, selector, body] of css.matchAll(/^([^\s{][^{\n]*?)\s*\{\n([\s\S]*?)^\}/gm)) {
		const props: Record<string, string> = {};
		for (const [, name, value] of body.matchAll(/(--[a-z0-9-]+):\s*(.*?);\s*$/gm))
			props[name] = value;
		out.set(selector, props);
	}
	return out;
}

// The rules a result lays over the stylesheet; a rule with nothing is left out.
function overlay(d: StylesheetDeclarations): Rules {
	const out: Rules = new Map();
	if (Object.keys(d.light).length) out.set(LIGHT_SELECTOR, d.light);
	if (Object.keys(d.dark).length) out.set(DARK_SELECTOR, d.dark);
	return out;
}

// SAVED with the dark color.primary rendering as value.
function withDarkPrimary(value: string): DesignTheme[] {
	return SAVED.map((th) => {
		if (th.name !== 'dark') return th;
		const tokens = th.tokens.map((t) =>
			t.path === 'color.primary' ? { ...t, variables: [{ name: '--color-primary', value }] } : t
		);
		return { ...th, tokens };
	});
}

const GREEN = parseColorText('oklch(60% 0.2 140)');

describe('stylesheetDeclarations', () => {
	it("gives the properties of flai's Tailwind golden, @theme's under :root", () => {
		const golden = rules(read('themes.tailwind.css.golden'));
		expect([...golden.keys()]).toEqual(['@theme', DARK_SELECTOR]);
		golden.set(LIGHT_SELECTOR, golden.get('@theme')!);
		golden.delete('@theme');
		const got = stylesheetDeclarations(SAVED);
		expect(got.refused).toEqual([]);
		expect(overlay(got)).toEqual(golden);
		expect(Object.keys(got.light)).toEqual(Object.keys(golden.get(LIGHT_SELECTOR)!));
		expect(Object.keys(got.dark)).toEqual(Object.keys(golden.get(DARK_SELECTOR)!));
	});

	it("gives the Design page's tried values, a light edit reaching dark where dark inherits", () => {
		const tried = editedThemes(SAVED, { [editKey('light', 'color.secondary')]: GREEN }, RAW);
		const got = stylesheetDeclarations(tried);
		expect(got.light['--color-secondary']).toBe('oklch(60% 0.2 140)');
		expect(got.dark).not.toHaveProperty('--color-secondary');
		const dark = editedThemes(SAVED, { [editKey('dark', 'color.secondary')]: GREEN }, RAW);
		expect(stylesheetDeclarations(dark).dark['--color-secondary']).toBe('oklch(60% 0.2 140)');
	});

	it('leaves out and lists a value that would break out of its rule, keeping the rest', () => {
		for (const value of ['url(https://example.com/x.png)', 'red; } body { display: none']) {
			const got = stylesheetDeclarations(withDarkPrimary(value));
			expect(got.dark, value).not.toHaveProperty('--color-primary');
			expect(got.dark['--color-surface']).toBe('oklch(23% 0.02 265)');
			expect(got.refused, value).toEqual([
				{ theme: 'dark', name: '--color-primary', value, reason: expect.any(String) }
			]);
		}
	});

	it('gives nothing for no themes', () => {
		expect(stylesheetDeclarations([])).toEqual({ light: {}, dark: {}, refused: [] });
	});
});
