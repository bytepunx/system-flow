// A proposal for the design system against the main branch's (S-0384): the tokens the two differ
// in, per theme and token path, each changed, added in the proposal, or removed from it, with both
// values. Two tokens differ when they render differently: their css and their custom properties,
// what the gallery's preview reads, so that an alias made literal at the same value, or a type
// renamed to one that renders alike, is not a difference.
import type { DesignSystem, DesignTheme, DesignToken } from '$lib/server/repo';

/** How a token of the proposal differs from the main branch's. */
export type TokenChange = 'changed' | 'added' | 'removed';

/** One token that differs: main is absent when it was added, proposal when it was removed. */
export type TokenDifference = {
	path: string;
	change: TokenChange;
	main?: DesignToken;
	proposal?: DesignToken;
};

/** One theme's differing tokens: main's file order, then the added ones in the proposal's. */
export type ThemeDifferences = { theme: DesignTheme['name']; differences: TokenDifference[] };

/** The two design systems' themes compared, as far as diffDesign reads them. */
export type Compared = Pick<DesignSystem, 'themes'>;

const THEMES: DesignTheme['name'][] = ['light', 'dark'];

/** Whether two tokens render alike: the same css and the same custom properties, in order. */
export function sameValue(a: DesignToken, b: DesignToken): boolean {
	return (
		a.css === b.css &&
		a.variables.length === b.variables.length &&
		a.variables.every((v, i) => v.name === b.variables[i].name && v.value === b.variables[i].value)
	);
}

/**
 * A token's rendered value as text: its css, then each further custom property it renders as, by
 * the part of its name after the first's, such as "1rem; line-height 1.5" for a type step.
 */
export function valueText(token: DesignToken): string {
	const [first, ...rest] = token.variables;
	const named = rest.map((v) => {
		const prefix = first ? `${first.name}--` : '';
		const name = prefix && v.name.startsWith(prefix) ? v.name.slice(prefix.length) : v.name;
		return `${name} ${v.value}`;
	});
	return [token.css, ...named].join('; ');
}

/** The tokens of one theme that differ between main and the proposal. */
function diffTheme(main: readonly DesignToken[], proposal: readonly DesignToken[]) {
	const proposed = new Map(proposal.map((t) => [t.path, t]));
	const known = new Set(main.map((t) => t.path));
	const out: TokenDifference[] = [];
	for (const m of main) {
		const p = proposed.get(m.path);
		if (!p) out.push({ path: m.path, change: 'removed', main: m });
		else if (!sameValue(m, p)) out.push({ path: m.path, change: 'changed', main: m, proposal: p });
	}
	for (const p of proposal)
		if (!known.has(p.path)) out.push({ path: p.path, change: 'added', proposal: p });
	return out;
}

/** The themes in which the proposal's tokens differ from main's, light then dark; [] when none. */
export function diffDesign(main: Compared, proposal: Compared): ThemeDifferences[] {
	const tokens = (ds: Compared, theme: DesignTheme['name']) =>
		ds.themes.find((t) => t.name === theme)?.tokens ?? [];
	return THEMES.map((theme) => ({
		theme,
		differences: diffTheme(tokens(main, theme), tokens(proposal, theme))
	})).filter((t) => t.differences.length > 0);
}
