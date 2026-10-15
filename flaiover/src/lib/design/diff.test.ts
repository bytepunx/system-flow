// S-0384: a proposal's tokens against the main branch's, per theme and token path.
import { describe, expect, it } from 'vitest';
import type { DesignTheme, DesignToken } from '$lib/server/repo';
import { diffDesign, sameValue, valueText } from './diff';

const colour = (name: string, hex: string, line = 1): DesignToken => ({
	path: `color.${name}`,
	type: 'color',
	value: { colorSpace: 'srgb', hex },
	css: hex,
	variables: [{ name: `--color-${name}`, value: hex }],
	line
});
const text = (step: string, size: string, lineHeight: string): DesignToken => ({
	path: `text.${step}`,
	type: 'typography',
	value: {},
	css: size,
	variables: [
		{ name: `--text-${step}`, value: size },
		{ name: `--text-${step}--line-height`, value: lineHeight }
	],
	line: 1
});
const system = (light: DesignToken[], dark: DesignToken[] = light) => ({
	themes: [
		{ name: 'light', tokens: light },
		{ name: 'dark', tokens: dark }
	] as DesignTheme[]
});

const MAIN = system(
	[colour('background', '#ffffff'), colour('primary', '#2563eb'), text('base', '1rem', '1.5')],
	[colour('background', '#0f172a'), colour('primary', '#60a5fa'), text('base', '1rem', '1.5')]
);

describe('diffDesign (S-0384)', () => {
	it('finds no difference between a design system and itself', () => {
		expect(diffDesign(MAIN, MAIN)).toEqual([]);
	});

	it('marks a token whose value changed, with both values', () => {
		const proposal = system(
			[colour('background', '#ffffff'), colour('primary', '#ff0000'), text('base', '1rem', '1.5')],
			MAIN.themes[1].tokens
		);
		expect(diffDesign(MAIN, proposal)).toEqual([
			{
				theme: 'light',
				differences: [
					{
						path: 'color.primary',
						change: 'changed',
						main: colour('primary', '#2563eb'),
						proposal: colour('primary', '#ff0000')
					}
				]
			}
		]);
	});

	it('marks a token added in the proposal and one removed from it', () => {
		const light = MAIN.themes[0].tokens;
		const proposal = system(
			[light[0], light[2], colour('accent', '#9333ea')],
			MAIN.themes[1].tokens
		);
		const [only] = diffDesign(MAIN, proposal);
		expect(only.theme).toBe('light');
		expect(only.differences.map((d) => [d.change, d.path])).toEqual([
			['removed', 'color.primary'],
			['added', 'color.accent']
		]);
		expect(only.differences[0].main?.css).toBe('#2563eb');
		expect(only.differences[0].proposal).toBeUndefined();
		expect(only.differences[1].main).toBeUndefined();
		expect(only.differences[1].proposal?.css).toBe('#9333ea');
	});

	it('compares each theme on its own, light then dark', () => {
		const proposal = system(
			[colour('background', '#fafafa'), ...MAIN.themes[0].tokens.slice(1)],
			[colour('background', '#000000'), ...MAIN.themes[1].tokens.slice(1)]
		);
		const themes = diffDesign(MAIN, proposal);
		expect(themes.map((t) => t.theme)).toEqual(['light', 'dark']);
		expect(themes[0].differences[0].proposal?.css).toBe('#fafafa');
		expect(themes[1].differences[0].main?.css).toBe('#0f172a');
		expect(themes[1].differences[0].proposal?.css).toBe('#000000');
		// A change in dark alone leaves light out.
		const darkOnly = system(MAIN.themes[0].tokens, proposal.themes[1].tokens);
		expect(diffDesign(MAIN, darkOnly).map((t) => t.theme)).toEqual(['dark']);
	});

	it('takes every token of a theme the proposal lacks as removed', () => {
		const proposal = { themes: [MAIN.themes[0]] };
		const [dark] = diffDesign(MAIN, proposal);
		expect(dark.theme).toBe('dark');
		expect(dark.differences.every((d) => d.change === 'removed')).toBe(true);
		expect(dark.differences.length).toBe(3);
	});

	it('compares rendered values: a further property counts, the line and the alias do not', () => {
		const proposal = system(
			[
				{ ...colour('background', '#ffffff', 9), ref: 'color.white', value: null },
				colour('primary', '#2563eb'),
				text('base', '1rem', '1.75')
			],
			MAIN.themes[1].tokens
		);
		const [light] = diffDesign(MAIN, proposal);
		expect(light.differences.map((d) => [d.change, d.path])).toEqual([['changed', 'text.base']]);
	});
});

describe('sameValue and valueText (S-0384)', () => {
	it('treats tokens alike when their css and custom properties are', () => {
		expect(sameValue(colour('a', '#000000'), colour('a', '#000000', 7))).toBe(true);
		expect(sameValue(colour('a', '#000000'), colour('a', '#000001'))).toBe(false);
		const base = text('base', '1rem', '1.5');
		expect(sameValue(base, { ...base, variables: base.variables.slice(0, 1) })).toBe(false);
	});

	it('writes a value as its css, then each further property by the rest of its name', () => {
		expect(valueText(colour('primary', '#2563eb'))).toBe('#2563eb');
		expect(valueText(text('base', '1rem', '1.5'))).toBe('1rem; line-height 1.5');
		const none = { ...colour('none', ''), css: '', variables: [] };
		expect(valueText(none)).toBe('');
	});
});
