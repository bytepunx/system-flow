import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import type { DesignSystem, DesignTheme } from '$lib/server/repo';
import {
	changed,
	colorFromHex,
	colorText,
	composeTokenFile,
	cssNumber,
	editedThemes,
	editKey,
	parseColorText,
	parseEditKey,
	pruneEdits,
	renderToken,
	valueProblem,
	type Edits
} from './tokens';

// flai's own design system fixtures, so that the port cannot drift from flai.
const DIR = new URL('../../../../flai/internal/designsystem/', import.meta.url);
const read = (name: string) => readFileSync(new URL(name, DIR), 'utf8');
const THEMES = read('testdata/themes.tokens.json');
const STARTER = read('starter.tokens.json');
// What flai design show --json answers for THEMES: the saved themes, as design.get gives them.
const SHOWN = JSON.parse(read('testdata/show.json.golden')) as DesignSystem;
const SAVED: DesignTheme[] = SHOWN.themes;

// The custom properties of one rule of the CSS flai renders THEMES as.
function golden(selector: string): Map<string, string> {
	const css = read('testdata/themes.css.golden');
	const start = css.indexOf(selector + ' {');
	const out = new Map<string, string>();
	const rule = css.slice(start, css.indexOf('}', start));
	for (const m of rule.matchAll(/(--[a-z0-9-]+): (.*);/g)) out.set(m[1], m[2]);
	return out;
}
const LIGHT = golden(':root');
const DARK = golden("[data-theme='dark']");

// The token at path in one of themes.
function token(themes: DesignTheme[], theme: string, path: string) {
	const t = themes.find((th) => th.name === theme)?.tokens.find((tk) => tk.path === path);
	if (!t) throw new Error(`no ${theme} token ${path}`);
	return t;
}

// The value at a path of keys in a parsed JSON document.
function at(doc: unknown, ...path: (string | number)[]): unknown {
	return path.reduce<unknown>((v, k) => (v as Record<string, unknown>)[k], doc);
}
const SET = ['sets', 'tokens', 'sources', 0];
const DARK_SOURCE = ['modifiers', 'theme', 'contexts', 'dark', 0];

const css = (path: string, type: string, value: unknown) => renderToken(path, type, value).css;

const GREEN = parseColorText('oklch(60% 0.2 140)');
const GREEN_JSON = JSON.stringify(GREEN);

describe('renderToken', () => {
	it('renders every token as flai does, but shadows, which the page does not edit', () => {
		const tokens = SAVED.flatMap((th) => th.tokens).filter((t) => t.type !== 'shadow');
		for (const t of tokens) {
			const want = { css: t.css, variables: t.variables };
			expect(renderToken(t.path, t.type, t.value), t.path).toEqual(want);
		}
		expect(tokens.length).toBeGreaterThan(60);
	});

	it('renders an edited colour, from a hex or oklch, as the golden CSS writes it', () => {
		const scrim = token(SAVED, 'light', 'color.scrim').value;
		const cases: [string, unknown, string | undefined][] = [
			['color.border', colorFromHex('#CCCCD1'), LIGHT.get('--color-border')],
			['color.scrim', colorFromHex('#000000', scrim), LIGHT.get('--color-scrim')],
			['color.primary', parseColorText('oklch(72% 0.15 255)'), DARK.get('--color-primary')],
			['color.shade', parseColorText('oklch(0 0 none / 0.1)'), 'oklch(0% 0 none / 0.1)']
		];
		for (const [path, value, want] of cases) expect(css(path, 'color', value), path).toBe(want);
		expect(LIGHT.get('--shadow-md')).toContain('oklch(0% 0 none / 0.1)');
		const border = renderToken('color.border', 'color', colorFromHex('#ccccd1'));
		expect(border.variables).toEqual([{ name: '--color-border', value: '#ccccd1' }]);
	});

	it('renders a dimension, spacing being the group key itself', () => {
		const px8 = { value: 8, unit: 'px' };
		const lg = [{ name: '--radius-lg', value: LIGHT.get('--radius-lg') }];
		expect(renderToken('radius.lg', 'dimension', px8).variables).toEqual(lg);
		const spacing = [{ name: '--spacing', value: LIGHT.get('--spacing') }];
		const rem = { value: 0.25, unit: 'rem' };
		expect(renderToken('spacing', 'dimension', rem).variables).toEqual(spacing);
		expect(css('radius.md', 'dimension', { value: 0.1 + 0.2, unit: 'rem' })).toBe('0.3rem');
	});

	it('renders a font family, quoting what is no plain identifier', () => {
		const sans = ['Inter Variable', 'Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'];
		const stack = [...sans, 'Apple Color Emoji'];
		expect(css('font.sans', 'fontFamily', stack)).toBe(LIGHT.get('--font-sans'));
		const odd = ['inherit', 'Say "hi"', 'a\\b'];
		expect(css('font.mono', 'fontFamily', odd)).toBe('"inherit", "Say \\"hi\\"", "a\\\\b"');
		expect(css('font.mono', 'fontFamily', 'Fira Code')).toBe('"Fira Code"');
	});

	it('renders a typography step with its line height, weight, and letter spacing', () => {
		const base = token(SAVED, 'light', 'text.base').value as object;
		const step = { ...base, fontSize: { value: 1.125, unit: 'rem' } };
		expect(renderToken('text.base', 'typography', step).variables).toEqual([
			{ name: '--text-base', value: '1.125rem' },
			{ name: '--text-base--line-height', value: LIGHT.get('--text-base--line-height') },
			{ name: '--text-base--font-weight', value: LIGHT.get('--text-base--font-weight') },
			{ name: '--text-base--letter-spacing', value: LIGHT.get('--text-base--letter-spacing') }
		]);
	});

	it('renders nothing for a value flai does not read, or a token outside the groups', () => {
		const white = colorFromHex('#ffffff');
		const px = { value: 1, unit: 'px' };
		const none = { css: '', variables: [] };
		expect(renderToken('color.primary', 'color', { ...white, hex: undefined })).toEqual(none);
		expect(renderToken('color.Primary', 'color', white)).toEqual(none);
		expect(renderToken('color', 'color', white)).toEqual(none);
		expect(renderToken('spacing.sm', 'dimension', px)).toEqual(none);
		expect(renderToken('radius.sm', 'color', white)).toEqual(none);
		expect(renderToken('brand.sm', 'dimension', px)).toEqual(none);
	});

	it('writes numbers as flai does: six decimals at most, no trailing zeros, no -0', () => {
		expect(cssNumber(1 / 3)).toBe('0.333333');
		expect(cssNumber(2.5)).toBe('2.5');
		expect(cssNumber(-0.0000001)).toBe('0');
		expect(cssNumber(0.21 * 100)).toBe('21');
	});
});

describe('the colour field', () => {
	it('makes an srgb colour of a hex as the starter writes it, keeping the alpha', () => {
		const blue = { colorSpace: 'srgb', components: [0.145, 0.388, 0.922], hex: '#2563eb' };
		expect(colorFromHex('#2563EB')).toEqual(blue);
		expect(colorFromHex('#000000', { alpha: 0.4 })).toMatchObject({ alpha: 0.4 });
		expect(colorFromHex('#00000')).toBeNull();
		expect(valueProblem('color', colorFromHex('#848e9f'))).toBeNull();
	});

	it('reads oklch with a percentage or a number, and refuses what is neither', () => {
		const green = { colorSpace: 'oklch', components: [0.52, 0.18, 255] };
		expect(parseColorText(' oklch(52% 0.18 255) ')).toEqual(green);
		const grey = { colorSpace: 'oklch', components: [0.7, 0.3, 'none'], alpha: 0.5 };
		expect(parseColorText('oklch(0.7 0.3 none)', { alpha: 0.5 })).toEqual(grey);
		expect(parseColorText('oklch(70% 0.1 20 / 50%)')).toMatchObject({ alpha: 0.5 });
		expect(parseColorText('#abcdef')).toMatchObject({ hex: '#abcdef' });
		for (const text of ['red', 'oklch(1 2)', 'oklch(a b c)', 'oklch(1 2% 3)', '#abc'])
			expect(parseColorText(text), text).toBeNull();
	});

	it('shows a colour as text its field reads back', () => {
		for (const t of SAVED[0].tokens.filter((t) => t.type === 'color')) {
			const text = colorText(t.value);
			expect(colorText(parseColorText(text)), t.path).toBe(text);
		}
		expect(colorText(token(SAVED, 'light', 'color.text').value)).toBe('oklch(21% 0.02 265)');
		expect(colorText(token(SAVED, 'light', 'color.surface').value)).toBe('#f5f5f7');
		expect(colorText('nope')).toBe('');
	});
});

const WHITE = { colorSpace: 'srgb', components: [1, 1, 1], hex: '#ffffff' };
const REM = { value: 1, unit: 'rem' };
const oklch = (...components: unknown[]) => ({ colorSpace: 'oklch', components });
const step = (o: object) => ({ fontSize: REM, lineHeight: 1.5, ...o });

// A type, a value of it, and what flai says is wrong with the value.
const PROBLEMS: [string, unknown, RegExp][] = [
	['color', '#fff', /a colour is an object .*, not a string/],
	['color', { components: [1, 1, 1] }, /needs its colorSpace/],
	['color', { colorSpace: 'hsl' }, /does not read the colour space "hsl"/],
	['color', oklch(0.5, 0.1), /an oklch colour needs components/],
	['color', oklch(0.5, '1', 2), /component 2 is a string/],
	['color', { ...WHITE, components: [1.2, 1, 1] }, /srgb component 1 is 1.2/],
	['color', oklch(52, 0.1, 2), /oklch lightness is 52; it is 0 to 1/],
	['color', oklch(0.5, -0.1, 2), /oklch chroma is -0.1/],
	['color', oklch(0.5, 0.1, 400), /oklch hue is 400/],
	['color', { ...WHITE, alpha: 2 }, /alpha is a number from 0 to 1/],
	['color', { colorSpace: 'srgb', components: [1, 1, 1] }, /the six-digit #ffffff/],
	['color', { ...WHITE, hex: '#fff' }, /hex is a six-digit #rrggbb/],
	['color', { ...WHITE, hex: '#000000' }, /are not the same colour/],
	['dimension', 4, /a dimension is an object .*, not a number/],
	['dimension', { unit: 'px' }, /needs its value, a number/],
	['dimension', { value: NaN, unit: 'px' }, /needs its value, a number/],
	['dimension', { value: 1, unit: 'em' }, /unit is rem or px/],
	['fontFamily', '', /a font family is a name/],
	['fontFamily', [], /a font family is a name/],
	['fontFamily', ['Inter', ''], /a font family is a name/],
	['typography', { lineHeight: 1.5 }, /needs its fontSize/],
	['typography', step({ fontSize: { value: 1, unit: 'em' } }), /^fontSize: .*rem or px/],
	['typography', { fontSize: REM }, /needs its lineHeight/],
	['typography', step({ lineHeight: 0 }), /lineHeight is 0; write a positive number/],
	['typography', step({ lineHeight: 'tall' }), /^lineHeight is a number or a dimension: /],
	['typography', step({ fontWeight: 'huge' }), /fontWeight is a number from 1 to 1000/],
	['typography', step({ fontFamily: [] }), /^fontFamily: a font family is a name/],
	['shadow', {}, /the Design page does not edit shadow tokens/]
];

describe('valueProblem', () => {
	it('says what is wrong with each value as flai does', () => {
		for (const [type, value, message] of PROBLEMS)
			expect(valueProblem(type, value), `${type} ${JSON.stringify(value)}`).toMatch(message);
	});

	it('refuses alpha on a colour role, but not on an extra colour', () => {
		const translucent = colorFromHex('#000000', { alpha: 0.4 });
		const role = /the colour role primary has alpha 0.4; a role is opaque/;
		expect(valueProblem('color', translucent, 'color.primary')).toMatch(role);
		expect(valueProblem('color', translucent, 'color.scrim')).toBeNull();
	});

	it('accepts every value the saved themes hold', () => {
		for (const t of SAVED[0].tokens.filter((t) => t.type !== 'shadow'))
			expect(valueProblem(t.type, t.value, t.path), t.path).toBeNull();
		const lineHeight = { value: 1.5, unit: 'rem' };
		expect(valueProblem('typography', step({ lineHeight }))).toBeNull();
		expect(valueProblem('fontFamily', 'Inter')).toBeNull();
	});
});

describe('editedThemes', () => {
	it('gives an edited token its value, css, and variables, and leaves the rest alone', () => {
		const before = structuredClone(SAVED);
		const radius = { value: 12, unit: 'px' };
		const out = editedThemes(SAVED, { [editKey('light', 'radius.lg')]: radius });
		const variables = [{ name: '--radius-lg', value: '12px' }];
		const lg = { value: radius, css: '12px', variables };
		expect(token(out, 'light', 'radius.lg')).toMatchObject(lg);
		expect(token(out, 'light', 'radius.md')).toBe(token(SAVED, 'light', 'radius.md'));
		expect(token(out, 'dark', 'radius.lg')).toBe(token(SAVED, 'dark', 'radius.lg'));
		expect(SAVED).toEqual(before);
	});

	it('changes a typography step and a font family', () => {
		const base = token(SAVED, 'light', 'text.base').value as object;
		const edits = {
			[editKey('light', 'text.base')]: { ...base, fontSize: { value: 1.125, unit: 'rem' } },
			[editKey('light', 'font.mono')]: ['Fira Code', 'monospace']
		};
		const out = editedThemes(SAVED, edits);
		expect(token(out, 'light', 'text.base').css).toBe('1.125rem');
		expect(token(out, 'light', 'font.mono').css).toBe('"Fira Code", monospace');
	});

	it('makes an alias follow the token it names, in that theme only', () => {
		const out = editedThemes(SAVED, { [editKey('light', 'color.primary')]: GREEN }, THEMES);
		const link = token(out, 'light', 'color.link');
		const green = 'oklch(60% 0.2 140)';
		expect(link).toMatchObject({ ref: 'color.primary', value: GREEN, css: green });
		expect(link.variables).toEqual([{ name: '--color-link', value: green }]);
		// dark overrides color.primary, so dark's link keeps dark's primary
		expect(token(out, 'dark', 'color.link')).toBe(token(SAVED, 'dark', 'color.link'));
		const dark = editedThemes(SAVED, { [editKey('dark', 'color.primary')]: GREEN });
		expect(token(dark, 'dark', 'color.link').css).toBe('oklch(60% 0.2 140)');
		expect(token(dark, 'light', 'color.link')).toBe(token(SAVED, 'light', 'color.link'));
	});

	it('carries a light edit into dark where dark does not override the token', () => {
		const dark = (themes: DesignTheme[], path: string) => token(themes, 'dark', path);
		const secondary = { [editKey('light', 'color.secondary')]: GREEN };
		const carried = editedThemes(SAVED, secondary, THEMES);
		expect(dark(carried, 'color.secondary').value).toEqual(GREEN);
		const unknown = editedThemes(SAVED, secondary);
		expect(dark(unknown, 'color.secondary')).toBe(dark(SAVED, 'color.secondary'));
		// dark overrides color.background
		const background = { [editKey('light', 'color.background')]: GREEN };
		const kept = editedThemes(SAVED, background, THEMES);
		expect(dark(kept, 'color.background')).toBe(dark(SAVED, 'color.background'));
		// dark's own edit wins over the light one
		const own = colorFromHex('#101010');
		const both = { ...secondary, [editKey('dark', 'color.secondary')]: own };
		const mine = editedThemes(SAVED, both, THEMES);
		expect(dark(mine, 'color.secondary').value).toEqual(own);
	});

	it('returns the saved themes as they were with no edits', () => {
		expect(editedThemes(SAVED, {}, THEMES)).toEqual(SAVED);
	});
});

describe('changed and pruneEdits', () => {
	const same = token(SAVED, 'light', 'color.text').value;
	const edits: Edits = {
		[editKey('light', 'color.text')]: structuredClone(same),
		[editKey('dark', 'color.text')]: colorFromHex('#ffffff'),
		[editKey('light', 'radius.sm')]: { value: 0.25, unit: 'rem' },
		[editKey('light', 'radius.md')]: { unit: 'rem', value: 0.5 },
		'dusk:radius.md': { value: 1, unit: 'rem' }
	};

	it('lists the edits that differ from the saved value, by deep equality', () => {
		expect(changed(SAVED, edits)).toEqual([
			{ theme: 'dark', path: 'color.text' },
			{ theme: 'light', path: 'radius.md' }
		]);
	});

	it('drops the edits equal to the saved value', () => {
		const kept = [editKey('dark', 'color.text'), editKey('light', 'radius.md')];
		expect(Object.keys(pruneEdits(SAVED, edits))).toEqual(kept);
		const sm = { [editKey('light', 'radius.sm')]: { value: 0.25, unit: 'rem' } };
		expect(pruneEdits(SAVED, sm)).toEqual({});
	});

	it('reads an edit key back', () => {
		const key = editKey('dark', 'color.text');
		expect(parseEditKey(key)).toEqual({ theme: 'dark', path: 'color.text' });
		expect(parseEditKey('dusk:color.text')).toBeNull();
		expect(parseEditKey('color.text')).toBeNull();
	});
});

describe('composeTokenFile', () => {
	it('returns the token file byte for byte with no edits', () => {
		expect(composeTokenFile(THEMES, {})).toBe(THEMES);
		expect(composeTokenFile(STARTER, {})).toBe(STARTER);
	});

	it('replaces a light and a dark value where each stands, and nothing else', () => {
		const black = colorFromHex('#000000');
		const out = composeTokenFile(THEMES, {
			[editKey('light', 'color.primary')]: GREEN,
			[editKey('dark', 'color.background')]: black
		});
		const primary = '{ "colorSpace": "oklch", "components": [0.52, 0.18, 255] }';
		const background = '{ "colorSpace": "oklch", "components": [0.18, 0.02, 265] }';
		const light = THEMES.replace(primary, GREEN_JSON);
		expect(out).toBe(light.replace(background, JSON.stringify(black)));
		const doc = JSON.parse(out);
		expect(at(doc, ...SET, 'color', 'primary', '$value')).toEqual(GREEN);
		expect(at(doc, ...DARK_SOURCE, 'color', 'background', '$value')).toEqual(black);
	});

	it('writes a typography step and a font family into the starter', () => {
		const base = { fontSize: { value: 1.125, unit: 'rem' }, lineHeight: 1.5 };
		const mono = ['Fira Code', 'monospace'];
		const out = composeTokenFile(STARTER, {
			[editKey('light', 'text.base')]: base,
			[editKey('light', 'font.mono')]: mono
		});
		const doc = JSON.parse(out);
		expect(at(doc, ...SET, 'text', 'base', '$value')).toEqual(base);
		expect(at(doc, ...SET, 'font', 'mono', '$value')).toEqual(mono);
		expect(out.split('\n').length).toBe(STARTER.split('\n').length);
	});

	it('adds a dark override to the dark context beside its neighbours', () => {
		const out = composeTokenFile(THEMES, { [editKey('dark', 'color.secondary')]: GREEN });
		const soft = '{ "colorSpace": "oklch", "components": [0.28, 0.06, 25] }';
		const anchor = `"danger-soft": { "$value": ${soft} }`;
		const added = `,\n              "secondary": { "$value": ${GREEN_JSON} }`;
		expect(out).toBe(THEMES.replace(anchor, anchor + added));
		expect(at(JSON.parse(out), ...DARK_SOURCE, 'color', 'secondary', '$value')).toEqual(GREEN);
	});

	it('adds the groups a dark override needs, once for every override under them', () => {
		const sm = { value: 2, unit: 'px' };
		const md = { value: 4, unit: 'px' };
		const edits = { [editKey('dark', 'radius.sm')]: sm, [editKey('dark', 'radius.md')]: md };
		const out = composeTokenFile(THEMES, edits);
		const block = [
			',',
			'            "radius": {',
			`              "sm": { "$value": ${JSON.stringify(sm)} },`,
			`              "md": { "$value": ${JSON.stringify(md)} }`,
			'            }'
		];
		expect(out.replace(block.join('\n'), '')).toBe(THEMES);
		const radius = { sm: { $value: sm }, md: { $value: md } };
		expect(at(JSON.parse(out), ...DARK_SOURCE, 'radius')).toEqual(radius);
	});

	it('makes a source of an empty dark context, keeping the order of every key', () => {
		const raw = [
			'{',
			'  "version": "2025.10",',
			'  "$extensions": { "z": 1, "50": 2 },',
			'  "sets": { "tokens": { "sources": [ { "color": { "$type": "color",',
			'    "primary": { "$value": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000" } }',
			'  } } ] } },',
			'  "modifiers": { "theme": { "contexts": { "light": [], "dark": [] } } },',
			'  "resolutionOrder": [{ "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" }]',
			'}',
			''
		].join('\n');
		const white = colorFromHex('#ffffff');
		const out = composeTokenFile(raw, { [editKey('dark', 'color.primary')]: white });
		const source = [
			'[',
			'    {',
			'      "color": {',
			`        "primary": { "$value": ${JSON.stringify(white)} }`,
			'      }',
			'    }',
			'  ]'
		];
		expect(out).toBe(raw.replace('"dark": []', `"dark": ${source.join('\n')}`));
		const doc = JSON.parse(out);
		expect(at(doc, ...DARK_SOURCE, 'color', 'primary', '$value')).toEqual(white);
	});

	it('refuses a token the file does not define, naming it and what to do', () => {
		const nope = /no set of the token file defines the token color\.nope.* or discard the edit/;
		const light = { [editKey('light', 'color.nope')]: GREEN };
		const dark = { [editKey('dark', 'color.nope')]: GREEN };
		expect(() => composeTokenFile(THEMES, light)).toThrow(nope);
		expect(() => composeTokenFile(THEMES, dark)).toThrow(nope);
	});

	it('writes a plain token file but no dark value in it, and refuses broken JSON', () => {
		const sm = '"sm": { "$value": { "value": 1, "unit": "px" } }';
		const plain = `{ "radius": { "$type": "dimension", ${sm} } }`;
		const two = { value: 2, unit: 'px' };
		const want = plain.replace('{ "value": 1, "unit": "px" }', JSON.stringify(two));
		expect(composeTokenFile(plain, { [editKey('light', 'radius.sm')]: two })).toBe(want);
		const dark = { [editKey('dark', 'radius.sm')]: two };
		expect(() => composeTokenFile(plain, dark)).toThrow(/is a plain token file/);
		const broken = { [editKey('light', 'a')]: 1 };
		expect(() => composeTokenFile('{ "a": [1, }', broken)).toThrow(/does not parse at line 1/);
	});
});
