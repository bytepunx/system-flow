// The design system preview's arithmetic (S-0382, ADR-0149): a theme's tokens, as flai's
// design.get answers them, as the --ds-* custom properties the gallery's stylesheet reads, and
// each colour's WCAG 2 contrast against the theme's ground, computed as flai computes it
// (flai/internal/designsystem/color.go).

/** A token as design.get answers it, only as far as the preview reads it. */
export interface PreviewToken {
	path: string;
	type: string;
	value: unknown;
	css: string;
	variables: { name: string; value: string }[];
}

/** A custom property left out of the preview, and why. */
export interface RefusedVariable {
	theme: string;
	name: string;
	value: string;
	reason: string;
}

/** The attribute that names a gallery panel's theme; its --ds-* values are set on it. */
export const THEME_ATTRIBUTE = 'data-ds-theme';

/** The themes the preview shows, in the order its panels stand. */
export const PREVIEW_THEMES = ['light', 'dark'] as const;
export type PreviewTheme = (typeof PREVIEW_THEMES)[number];

/** The token whose colour every swatch's contrast is measured against. */
export const GROUND = 'color.background';

// A variable's name: -- and a name as flai's cssNameRe allows it.
const NAME = /^--[a-z0-9][a-z0-9-]*$/;

// What a value must not carry into a stylesheet (ADR-0149 § 4): flai's parse already keeps these
// out, and the preview refuses them again rather than trust it. An unclosed comment would swallow
// the declarations after it.
const UNSAFE: [RegExp, string][] = [
	[/url\(/i, 'url('],
	[/[{}]/, 'a brace'],
	[/;/, 'a semicolon'],
	[/</, '<'],
	[/@/, '@'],
	[/\\/, 'a backslash'],
	[/\/\*/, 'a comment']
];

/** Why a custom property cannot go into the preview's stylesheet, or null when it can. */
export function unsafeVariable(name: string, value: string): string | null {
	if (!NAME.test(name)) return `its name ${JSON.stringify(name)} is not a custom property`;
	for (const [re, what] of UNSAFE) if (re.test(value)) return `its value contains ${what}`;
	return null;
}

/** A theme's tokens as --ds-* declarations, each variable renamed from --x to --ds-x, and those refused. */
export function dsDeclarations(
	theme: string,
	tokens: readonly PreviewToken[]
): { declarations: string[]; refused: RefusedVariable[] } {
	const declarations: string[] = [];
	const refused: RefusedVariable[] = [];
	for (const token of tokens) {
		for (const { name, value } of token.variables) {
			const reason = unsafeVariable(name, value);
			if (reason) refused.push({ theme, name, value, reason });
			else declarations.push(`--ds-${name.slice(2)}: ${value};`);
		}
	}
	return { declarations, refused };
}

/** The preview's one token stylesheet: each theme's --ds-* values on its panel, and the values refused. */
export function previewStyle(themes: Record<PreviewTheme, readonly PreviewToken[]>): {
	css: string;
	refused: RefusedVariable[];
} {
	const refused: RefusedVariable[] = [];
	const blocks = PREVIEW_THEMES.map((theme) => {
		const d = dsDeclarations(theme, themes[theme]);
		refused.push(...d.refused);
		const body = [`color-scheme: ${theme};`, ...d.declarations].map((l) => `\t${l}`).join('\n');
		return `[${THEME_ATTRIBUTE}='${theme}'] {\n${body}\n}`;
	});
	return { css: blocks.join('\n\n') + '\n', refused };
}

/** A gamma-encoded sRGB colour, each channel 0 to 1, and whether it had to be clipped to sRGB. */
export interface SRGB {
	rgb: [number, number, number];
	hex: string;
	clipped: boolean;
}

// How far outside 0 to 1 a converted channel may fall and still count as in gamut: half an 8-bit
// step, as flai's clipTolerance.
const CLIP_TOLERANCE = 0.5 / 255;

/** A DTCG colour value (srgb with its hex or components, or oklch) in sRGB, or null when it is not one. */
export function toSRGB(value: unknown): SRGB | null {
	if (typeof value !== 'object' || value === null) return null;
	const v = value as { colorSpace?: unknown; components?: unknown; hex?: unknown };
	const comps = Array.isArray(v.components) ? v.components : null;
	if (v.colorSpace === 'srgb') {
		if (typeof v.hex === 'string' && /^#[0-9a-f]{6}$/i.test(v.hex)) {
			const n = parseInt(v.hex.slice(1), 16);
			return encoded([(n >> 16) / 255, ((n >> 8) & 255) / 255, (n & 255) / 255], false);
		}
		if (comps?.length === 3 && comps.every((c) => typeof c === 'number'))
			return encoded(comps as [number, number, number], false);
		return null;
	}
	if (v.colorSpace !== 'oklch' || comps?.length !== 3) return null;
	const [l, c] = comps;
	const h = comps[2] === 'none' ? 0 : comps[2];
	if (typeof l !== 'number' || typeof c !== 'number' || typeof h !== 'number') return null;
	const rad = (h * Math.PI) / 180;
	const linear = oklabToLinearSRGB(l, c * Math.cos(rad), c * Math.sin(rad));
	const gamma = linear.map(encode) as [number, number, number];
	const clipped = gamma.some((x) => x < -CLIP_TOLERANCE || x > 1 + CLIP_TOLERANCE);
	return encoded(gamma, clipped);
}

// Björn Ottosson's OKLab to linear sRGB, through LMS, with the matrices CSS Color 4 and flai use.
function oklabToLinearSRGB(l: number, a: number, b: number): [number, number, number] {
	const l1 = l + 0.3963377774 * a + 0.2158037573 * b;
	const m1 = l - 0.1055613458 * a - 0.0638541728 * b;
	const s1 = l - 0.0894841775 * a - 1.291485548 * b;
	const lc = l1 * l1 * l1;
	const mc = m1 * m1 * m1;
	const sc = s1 * s1 * s1;
	return [
		4.0767416621 * lc - 3.3077115913 * mc + 0.2309699292 * sc,
		-1.2684380046 * lc + 2.6097574011 * mc - 0.3413193965 * sc,
		-0.0041960863 * lc - 0.7034186147 * mc + 1.707614701 * sc
	];
}

// The sRGB transfer function, extended to negative values by symmetry as CSS Color 4 extends it.
function encode(v: number): number {
	const x = Math.abs(v);
	const e = x > 0.0031308 ? 1.055 * Math.pow(x, 1 / 2.4) - 0.055 : 12.92 * x;
	return Math.sign(v) * e;
}

const byte8 = (v: number) => Math.round(Math.min(Math.max(v, 0), 1) * 255);

// The colour clamped to sRGB and rounded to the 8-bit channels its hex gives, as flai's contrast reads it.
function encoded(rgb: [number, number, number], clipped: boolean): SRGB {
	const bytes = rgb.map(byte8);
	return {
		rgb: bytes.map((b) => b / 255) as [number, number, number],
		hex: '#' + bytes.map((b) => b.toString(16).padStart(2, '0')).join(''),
		clipped
	};
}

/** The colour's WCAG 2 relative luminance, 0 for black to 1 for white. */
export function luminance([r, g, b]: readonly [number, number, number]): number {
	const lin = (v: number) => (v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4));
	return 0.2126 * lin(r) + 0.7152 * lin(g) + 0.0722 * lin(b);
}

/** One colour's contrast on another, as flai reports a pair. */
export interface Contrast {
	foregroundHex: string;
	backgroundHex: string;
	/** The WCAG 2 ratio, truncated to two decimals so that it never reads as passing when it fails. */
	ratio: number;
	/** The ratio before truncation, for a pass or fail against a threshold. */
	exact: number;
	/** Either colour is an oklch colour outside sRGB, clipped to it. */
	clipped: boolean;
}

/** The WCAG 2 contrast of two DTCG colour values, or null when either is not a colour the preview reads. */
export function contrast(foreground: unknown, background: unknown): Contrast | null {
	const fg = toSRGB(foreground);
	const bg = toSRGB(background);
	if (!fg || !bg) return null;
	const [hi, lo] = [luminance(fg.rgb), luminance(bg.rgb)].sort((a, b) => b - a);
	const exact = (hi + 0.05) / (lo + 0.05);
	return {
		foregroundHex: fg.hex,
		backgroundHex: bg.hex,
		ratio: Math.floor(exact * 100) / 100,
		exact,
		clipped: fg.clipped || bg.clipped
	};
}

/** Each colour token of a theme with its contrast on the theme's ground, in the tokens' order; empty without a ground. */
export function contrastOnGround(tokens: readonly PreviewToken[]): ({ path: string } & Contrast)[] {
	const ground = tokens.find((t) => t.path === GROUND);
	if (!ground) return [];
	return tokens.flatMap((t) => {
		if (t.type !== 'color') return [];
		const c = contrast(t.value, ground.value);
		return c ? [{ path: t.path, ...c }] : [];
	});
}
