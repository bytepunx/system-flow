// The token values the Components page lays over a project's generated design stylesheet (S-0441,
// ADR-0159): the custom properties flai design css --write writes for the given themes, light
// under :root and dark under [data-theme='dark'], as flai's Tailwind rendering lays them out
// (flai/internal/designsystem/css.go). Tailwind compiles the stylesheet's @theme rule to :root,
// so the light values go there; the dark rule carries only the properties whose dark value
// differs, and, set after the light rule, wins where the attribute is.

import type { DesignTheme, DesignToken } from '$lib/server/repo';
import { unsafeVariable, type RefusedVariable } from '$lib/design/preview';

/** The selector the generated stylesheet's light values take effect under. */
export const LIGHT_SELECTOR = ':root';

/** The selector the generated stylesheet's dark values sit under (ADR-0159). */
export const DARK_SELECTOR = "[data-theme='dark']";

/** Custom properties by name, with their leading --, as the harness's tokens message takes them. */
export type Declarations = Record<string, string>;

/** A stylesheet's light and dark custom properties, and those refused as unsafe. */
export type StylesheetDeclarations = {
	light: Declarations;
	dark: Declarations;
	refused: RefusedVariable[];
};

// The groups in the order flai's Groups gives, which orders a theme's properties.
const GROUPS = ['color', 'font', 'text', 'spacing', 'radius', 'shadow'];

// A token's place in GROUPS; tokens outside them render nothing and go last.
const rank = (t: DesignToken) => {
	const i = GROUPS.indexOf(t.path.split('.')[0]);
	return i < 0 ? GROUPS.length : i;
};

// A theme's properties, by group then in file order, as flai's themeDeclarations gives them.
function themeVariables(theme: DesignTheme | undefined): { name: string; value: string }[] {
	const tokens = [...(theme?.tokens ?? [])].sort((a, b) => rank(a) - rank(b));
	return tokens.flatMap((t) => t.variables);
}

/**
 * The custom properties the generated design stylesheet carries for themes, light then dark as
 * design.get gives them: every light one, and each dark one whose value differs from its light
 * value; a property whose name or value preview's unsafeVariable refuses is left out and listed.
 */
export function stylesheetDeclarations(themes: readonly DesignTheme[]): StylesheetDeclarations {
	const out: StylesheetDeclarations = { light: {}, dark: {}, refused: [] };
	const light = themeVariables(themes.find((th) => th.name === 'light'));
	const lightValues = new Map(light.map((v) => [v.name, v.value]));
	const dark = themeVariables(themes.find((th) => th.name === 'dark')).filter(
		(v) => lightValues.get(v.name) !== v.value
	);
	const put = (theme: 'light' | 'dark', { name, value }: { name: string; value: string }) => {
		const reason = unsafeVariable(name, value);
		if (reason) out.refused.push({ theme, name, value, reason });
		else out[theme][name] = value;
	};
	for (const v of light) put('light', v);
	for (const v of dark) put('dark', v);
	return out;
}
