// Editing the design system's tokens on the Design page (S-0383): a token's value rendered as the
// CSS custom properties flai renders it as (flai/internal/designsystem/css.go), checked as flai
// checks it (tokens.go), the saved themes with the edits applied, and the token file with the
// edited values written into it, every other byte kept.
import type { DesignTheme, DesignToken } from '$lib/server/repo';

/** A theme the token file gives. */
export type Theme = 'light' | 'dark';

/** An edited token: its theme and its dotted path. */
export type EditTarget = { theme: Theme; path: string };

/** Edited values, each a DTCG $value, keyed by editKey(theme, path). */
export type Edits = Record<string, unknown>;

/** A JSON object, such as a colour value. */
export type JSONObject = Record<string, unknown>;

/** A CSS custom property and its value. */
export type Variable = { name: string; value: string };

/** A token's CSS custom properties, and the first one's value; '' when it renders none. */
export type Rendered = { css: string; variables: Variable[] };

/** The key of an edit to the token at path in theme. */
export function editKey(theme: Theme, path: string): string {
	return `${theme}:${path}`;
}

/** The theme and path of an edit's key; null when it is not one. */
export function parseEditKey(key: string): EditTarget | null {
	const at = key.indexOf(':');
	const theme = key.slice(0, at);
	if (at < 1 || at === key.length - 1) return null;
	if (theme !== 'light' && theme !== 'dark') return null;
	return { theme, path: key.slice(at + 1) };
}

// ---- Reading a value as flai reads it (tokens.go) ----

type Color = {
	space: 'srgb' | 'oklch';
	components: [number, number, number];
	noHue: boolean;
	alpha: number;
	hex: string;
};

type Dimension = { value: number; unit: string };

type Typography = {
	fontSize: Dimension;
	lineHeight: number | Dimension;
	fontWeight: number;
	letterSpacing: Dimension | null;
};

function isNumber(v: unknown): v is number {
	return typeof v === 'number' && Number.isFinite(v);
}

function isObject(v: unknown): v is JSONObject {
	return typeof v === 'object' && v !== null && !Array.isArray(v);
}

// A value's kind as flai's messages name it.
function kind(v: unknown): string {
	if (v === null || v === undefined) return 'null';
	if (Array.isArray(v)) return 'an array';
	if (typeof v === 'object') return 'an object';
	return `a ${typeof v}`;
}

const HEX = /^#[0-9a-fA-F]{6}$/;

// A 0 to 1 channel as the nearest 8-bit value.
const byte8 = (v: number) => Math.round(Math.min(Math.max(v, 0), 1) * 255);

// A 0 to 1 channel as two hex digits.
const hexByte = (v: number) => byte8(v).toString(16).padStart(2, '0');

// The #rrggbb an srgb colour's components round to.
function hexOf(components: readonly number[]): string {
	return '#' + components.map(hexByte).join('');
}

// flai's DecodeColor: the colour, or what is wrong with it.
function decodeColor(v: unknown): Color | string {
	if (!isObject(v))
		return `a colour is an object such as {"colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff"}, not ${kind(v)}`;
	const space = v.colorSpace;
	if (typeof space !== 'string') return 'a colour needs its colorSpace, srgb or oklch';
	if (space !== 'srgb' && space !== 'oklch')
		return `flai does not read the colour space ${JSON.stringify(space)}; write the colour in srgb, with its hex, or in oklch`;
	const comps = v.components;
	if (!Array.isArray(comps) || comps.length !== 3)
		return `an ${space} colour needs components, an array of three numbers`;
	const noHue = space === 'oklch' && comps[2] === 'none';
	const c: Color = { space, components: [0, 0, 0], noHue, alpha: 1, hex: '' };
	for (let i = 0; i < 3; i++) {
		if (i === 2 && noHue) continue;
		if (!isNumber(comps[i])) return `component ${i + 1} is ${kind(comps[i])}; write a number`;
		c.components[i] = comps[i];
	}
	const range = rangeProblem(c);
	if (range) return range;
	if (v.alpha !== undefined) {
		if (!isNumber(v.alpha) || v.alpha < 0 || v.alpha > 1)
			return 'alpha is a number from 0 to 1; correct it or remove it';
		c.alpha = v.alpha;
	}
	if (v.hex === undefined && space === 'srgb')
		return `an srgb colour needs its hex, the six-digit ${hexOf(c.components)}`;
	if (v.hex === undefined) return c;
	if (typeof v.hex !== 'string' || !HEX.test(v.hex))
		return 'hex is a six-digit #rrggbb such as #1a2b3c; correct it';
	c.hex = v.hex.toLowerCase();
	if (space === 'srgb' && !hexAgrees(c))
		return `its hex ${c.hex} and its components, which give ${hexOf(c.components)}, are not the same colour; make them agree`;
	return c;
}

// Which component is outside its colour space's range, or ''.
function rangeProblem(c: Color): string {
	const [l, ch, h] = c.components;
	if (c.space === 'srgb') {
		const i = c.components.findIndex((v) => v < 0 || v > 1);
		if (i < 0) return '';
		return `srgb component ${i + 1} is ${c.components[i]}; red, green, and blue are 0 to 1`;
	}
	if (l < 0 || l > 1) return `oklch lightness is ${l}; it is 0 to 1`;
	if (ch < 0) return `oklch chroma is ${ch}; it is 0 or more`;
	if (h < 0 || h > 360) return `oklch hue is ${h}; it is 0 to 360 degrees, or "none"`;
	return '';
}

// Whether an srgb colour's hex and components are one colour, each channel within 1/255.
function hexAgrees(c: Color): boolean {
	return c.components.every((v, i) => {
		const b = parseInt(c.hex.slice(1 + 2 * i, 3 + 2 * i), 16);
		return Math.abs(v - b / 255) <= 1 / 255 + 1e-9;
	});
}

// flai's DecodeDimension.
function decodeDimension(v: unknown): Dimension | string {
	if (!isObject(v))
		return `a dimension is an object such as {"value": 0.25, "unit": "rem"}, not ${kind(v)}`;
	if (!isNumber(v.value)) return 'a dimension needs its value, a number';
	if (v.unit !== 'rem' && v.unit !== 'px')
		return "a dimension's unit is rem or px; flai reads no other";
	return { value: v.value, unit: v.unit };
}

// flai's DecodeFontFamily.
function decodeFontFamily(v: unknown): string[] | string {
	const want =
		'a font family is a name, or a non-empty array of names such as ["Inter", "sans-serif"]';
	if (typeof v === 'string' && v !== '') return [v];
	if (!Array.isArray(v) || v.length === 0) return want;
	if (!v.every((f) => typeof f === 'string' && f !== '')) return want;
	return v as string[];
}

// DTCG's named font weights.
const FONT_WEIGHTS: Record<string, number> = {
	thin: 100,
	hairline: 100,
	'extra-light': 200,
	'ultra-light': 200,
	light: 300,
	normal: 400,
	regular: 400,
	book: 400,
	medium: 500,
	'semi-bold': 600,
	'demi-bold': 600,
	bold: 700,
	'extra-bold': 800,
	'ultra-bold': 800,
	black: 900,
	heavy: 900,
	'extra-black': 950,
	'ultra-black': 950
};

// A typography value's fontWeight as a number, 0 when it gives none, or what is wrong with it.
function decodeWeight(w: unknown): number | string {
	if (w === undefined) return 0;
	if (isNumber(w) && w >= 1 && w <= 1000) return w;
	if (typeof w === 'string' && Object.hasOwn(FONT_WEIGHTS, w)) return FONT_WEIGHTS[w];
	return 'fontWeight is a number from 1 to 1000 or a DTCG weight name such as bold';
}

// flai's DecodeTypography.
function decodeTypography(v: unknown): Typography | string {
	if (!isObject(v))
		return `a typography value is an object with fontSize and lineHeight, not ${kind(v)}`;
	if (v.fontSize === undefined) return 'a typography value needs its fontSize, a dimension';
	const fontSize = decodeDimension(v.fontSize);
	if (typeof fontSize === 'string') return `fontSize: ${fontSize}`;
	if (v.lineHeight === undefined)
		return 'a typography value needs its lineHeight, a number or a dimension';
	const lh = isNumber(v.lineHeight) ? v.lineHeight : decodeDimension(v.lineHeight);
	if (typeof lh === 'string') return `lineHeight is a number or a dimension: ${lh}`;
	const fontWeight = decodeWeight(v.fontWeight);
	if (typeof fontWeight === 'string') return fontWeight;
	const family = v.fontFamily === undefined ? [] : decodeFontFamily(v.fontFamily);
	if (typeof family === 'string') return `fontFamily: ${family}`;
	const spacing = v.letterSpacing === undefined ? null : decodeDimension(v.letterSpacing);
	if (typeof spacing === 'string') return `letterSpacing: ${spacing}`;
	return { fontSize, lineHeight: lh, fontWeight, letterSpacing: spacing };
}

// The decoders of the types the page edits.
const DECODE: Record<string, (v: unknown) => object | string> = {
	color: decodeColor,
	dimension: decodeDimension,
	fontFamily: decodeFontFamily,
	typography: decodeTypography
};

// The colour roles, which are opaque (ADR-0148).
const COLOR_ROLES = [
	'background',
	'surface',
	'text',
	'text-muted',
	'border',
	'primary',
	'on-primary',
	'secondary',
	'on-secondary',
	'focus',
	'success',
	'success-soft',
	'warning',
	'warning-soft',
	'danger',
	'danger-soft',
	'info',
	'info-soft'
];

/**
 * What is wrong with a value of type, in flai's words, or null when flai reads it. Given the
 * token's path, a colour role with alpha is wrong, as flai says; a typography step's line height
 * must be positive. The page edits colours, dimensions, font families, and typography steps, and
 * refuses any other type.
 */
export function valueProblem(type: string, value: unknown, path?: string): string | null {
	if (!Object.hasOwn(DECODE, type))
		return `the Design page does not edit ${type || 'untyped'} tokens; change it in the token file`;
	const decoded = DECODE[type](value);
	if (typeof decoded === 'string') return decoded;
	const name = path?.startsWith('color.') ? path.slice(6) : '';
	const alpha = type === 'color' ? (decoded as Color).alpha : 1;
	if (alpha !== 1 && COLOR_ROLES.includes(name))
		return `the colour role ${name} has alpha ${alpha}; a role is opaque, alpha is only for shadows and extra tokens: remove alpha`;
	const lh = type === 'typography' ? (decoded as Typography).lineHeight : 1;
	if (typeof lh === 'number' && lh <= 0)
		return `lineHeight is ${lh}; write a positive number, a multiple of the font size`;
	return null;
}

// ---- Rendering a value as flai renders it (css.go) ----

// The groups flai renders, and the types each holds.
const GROUP_TYPES: Record<string, string[]> = {
	color: ['color'],
	font: ['fontFamily'],
	text: ['typography', 'dimension'],
	spacing: ['dimension'],
	radius: ['dimension'],
	shadow: ['shadow']
};

const CSS_NAME = /^[a-z0-9][a-z0-9-]*$/;

// Whether a token sits where flai renders it: spacing is the group key itself, and every other
// token one name under its group.
function placed(parts: string[]): boolean {
	if (!Object.hasOwn(GROUP_TYPES, parts[0])) return false;
	return parts.length === (parts[0] === 'spacing' ? 1 : 2);
}

/** A number as flai writes it in CSS: to six decimals at most, with no trailing zeros. */
export function cssNumber(v: number): string {
	// Go's math.Round, which rounds half away from zero
	const r = (Math.sign(v) * Math.round(Math.abs(v) * 1e6)) / 1e6;
	if (r === 0) return '0';
	return Math.abs(r) >= 1e21 ? BigInt(r).toString() : String(r);
}

const dimensionString = (d: Dimension) => cssNumber(d.value) + d.unit;

// An oklch colour's lightness as a percentage, its chroma, and its hue, as flai writes them.
function oklchText(c: Color): string {
	const [l, ch, h] = c.components;
	return `${cssNumber(l * 100)}% ${cssNumber(ch)} ${c.noHue ? 'none' : cssNumber(h)}`;
}

// A colour as flai writes it: srgb as #rrggbb, #rrggbbaa below alpha 1; oklch as
// oklch(63.7% 0.237 25.331), with / alpha below alpha 1.
function colorString(c: Color): string {
	if (c.space === 'srgb' && c.alpha < 1) return hexOf(c.components) + hexByte(c.alpha);
	if (c.space === 'srgb') return hexOf(c.components);
	const alpha = c.alpha < 1 ? ' / ' + cssNumber(c.alpha) : '';
	return `oklch(${oklchText(c)}${alpha})`;
}

// CSS's generic font families, written bare.
const GENERIC_FAMILIES = [
	'serif',
	'sans-serif',
	'monospace',
	'cursive',
	'fantasy',
	'system-ui',
	'ui-serif',
	'ui-sans-serif',
	'ui-monospace',
	'ui-rounded',
	'math',
	'emoji',
	'fangsong'
];

// The keywords that may not be a family name unquoted.
const CSS_WIDE_KEYWORDS = ['inherit', 'initial', 'unset', 'revert', 'revert-layer', 'default'];

const IDENT = /^-?[A-Za-z_][A-Za-z0-9_-]*$/;

// One font family as flai writes it: bare when generic or a plain identifier, else quoted.
function cssFamily(f: string): string {
	if (GENERIC_FAMILIES.includes(f)) return f;
	if (IDENT.test(f) && !CSS_WIDE_KEYWORDS.includes(f.toLowerCase())) return f;
	return '"' + f.replaceAll('\\', '\\\\').replaceAll('"', '\\"') + '"';
}

// A value's properties, the first named name, as flai's Declarations gives them; null when the
// value does not decode.
function declarations(name: string, type: string, value: unknown): Variable[] | null {
	if (!Object.hasOwn(DECODE, type)) return null;
	const v = DECODE[type](value);
	if (typeof v === 'string') return null;
	if (type === 'color') return [{ name, value: colorString(v as Color) }];
	if (type === 'dimension') return [{ name, value: dimensionString(v as Dimension) }];
	if (type === 'fontFamily') return [{ name, value: (v as string[]).map(cssFamily).join(', ') }];
	const t = v as Typography;
	const lh = t.lineHeight;
	const out: Variable[] = [];
	const add = (suffix: string, value: string) => out.push({ name: name + suffix, value });
	add('', dimensionString(t.fontSize));
	add('--line-height', typeof lh === 'number' ? cssNumber(lh) : dimensionString(lh));
	if (t.fontWeight !== 0) add('--font-weight', cssNumber(t.fontWeight));
	if (t.letterSpacing) add('--letter-spacing', dimensionString(t.letterSpacing));
	return out;
}

/**
 * The CSS custom properties a token renders as, and the first one's value, as flai's Declarations
 * and CSSValue give them, for a colour, a dimension, a font family, or a typography step; none,
 * and css '', when it renders nothing.
 */
export function renderToken(path: string, type: string, value: unknown): Rendered {
	const parts = path.split('.');
	const [group, name] = [parts[0], parts[parts.length - 1]];
	const none: Rendered = { css: '', variables: [] };
	if (!placed(parts) || !CSS_NAME.test(name) || !GROUP_TYPES[group].includes(type)) return none;
	const property = group === 'spacing' ? '--spacing' : `--${group}-${name}`;
	const variables = declarations(property, type, value);
	return variables ? { css: variables[0].value, variables } : none;
}

// ---- The colour field ----

/**
 * An srgb colour of a #rrggbb hex, its components to three decimals as flai's starter writes
 * them, keeping previous's alpha; null when hex is not #rrggbb.
 */
export function colorFromHex(hex: string, previous?: unknown): JSONObject | null {
	if (!HEX.test(hex)) return null;
	const n = parseInt(hex.slice(1), 16);
	const components = [n >> 16, (n >> 8) & 255, n & 255].map(
		(b) => Math.round((b / 255) * 1000) / 1000
	);
	return withAlpha({ colorSpace: 'srgb', components, hex: hex.toLowerCase() }, previous);
}

// The colour with previous's alpha, when previous is a colour that has one.
function withAlpha(colour: JSONObject, previous: unknown): JSONObject {
	if (!isObject(previous) || previous.alpha === undefined) return colour;
	return { ...colour, alpha: previous.alpha };
}

// A number the colour field gives, a percentage divided by 100; NaN when it is neither.
function fieldNumber(s: string): number {
	const m = /^([+-]?(?:\d+\.?\d*|\.\d+)(?:e[+-]?\d+)?)(%?)$/i.exec(s.trim());
	if (!m) return NaN;
	return m[2] ? Number((Number(m[1]) / 100).toPrecision(12)) : Number(m[1]);
}

/**
 * The colour a colour field's text gives: #rrggbb, or oklch(L C H) with L from 0 to 1 or a
 * percentage, H a number or none, and an optional / alpha; previous's alpha is kept unless the
 * text gives one. Null when the text is neither; valueProblem says what is out of range.
 */
export function parseColorText(text: string, previous?: unknown): JSONObject | null {
	const t = text.trim();
	if (HEX.test(t)) return colorFromHex(t, previous);
	const m = /^oklch\((.*)\)$/i.exec(t);
	if (!m) return null;
	const [lch, alpha, ...more] = m[1].split('/');
	const parts = lch.trim().split(/\s+/);
	if (more.length > 0 || parts.length !== 3 || /%/.test(parts[1] + parts[2])) return null;
	const [l, c] = [fieldNumber(parts[0]), fieldNumber(parts[1])];
	const h = parts[2].toLowerCase() === 'none' ? 'none' : fieldNumber(parts[2]);
	if (Number.isNaN(l) || Number.isNaN(c) || Number.isNaN(h)) return null;
	const colour = { colorSpace: 'oklch', components: [l, c, h] };
	if (alpha === undefined) return withAlpha(colour, previous);
	const a = fieldNumber(alpha);
	return Number.isNaN(a) ? null : { ...colour, alpha: a };
}

/** The text a colour field shows for a colour: its #rrggbb, or oklch(L% C H); '' for none. */
export function colorText(value: unknown): string {
	const c = decodeColor(value);
	if (typeof c === 'string') return '';
	return c.space === 'srgb' ? hexOf(c.components) : `oklch(${oklchText(c)})`;
}

// ---- Edits against the saved themes ----

// Whether two JSON values are equal, member by member.
function jsonEqual(a: unknown, b: unknown): boolean {
	if (a === b) return true;
	if (Array.isArray(a) && Array.isArray(b))
		return a.length === b.length && a.every((x, i) => jsonEqual(x, b[i]));
	if (!isObject(a) || !isObject(b)) return false;
	const keys = Object.keys(a);
	if (keys.length !== Object.keys(b).length) return false;
	return keys.every((k) => Object.hasOwn(b, k) && jsonEqual(a[k], b[k]));
}

// The saved token at path in theme.
function savedToken(saved: readonly DesignTheme[], theme: Theme, path: string) {
	return saved.find((th) => th.name === theme)?.tokens.find((t) => t.path === path);
}

// Whether an edit differs from the value saved for its token: an edit of no saved token does, and
// a key that is no editKey never does.
function differs(saved: readonly DesignTheme[], key: string, value: unknown): boolean {
	const k = parseEditKey(key);
	if (!k) return false;
	const token = savedToken(saved, k.theme, k.path);
	return !token || !jsonEqual(token.value, value);
}

/** The edits that differ from the saved values, in the order the edits were made. */
export function changed(saved: readonly DesignTheme[], edits: Edits): EditTarget[] {
	return Object.keys(edits).flatMap((key) => {
		const k = parseEditKey(key);
		return k && differs(saved, key, edits[key]) ? [k] : [];
	});
}

/** The edits without those equal to the saved value, or whose key is no editKey. */
export function pruneEdits(saved: readonly DesignTheme[], edits: Edits): Edits {
	return Object.fromEntries(Object.entries(edits).filter(([k, v]) => differs(saved, k, v)));
}

/**
 * The saved themes with the edits applied: each edited token's value, css, and variables, an
 * alias following the token it names in its theme, and every other token as it was. Given the
 * token file's raw text, a light edit applies to dark too where neither theme context overrides
 * the token and dark has no edit of its own. saved is left as it is.
 */
export function editedThemes(
	saved: readonly DesignTheme[],
	edits: Edits,
	raw?: string | null
): DesignTheme[] {
	if (Object.keys(edits).length === 0) return [...saved];
	const overridden = raw ? contextPaths(raw) : null;
	return saved.map((theme) => editedTheme(theme, edits, overridden));
}

type Edited = { value: unknown };

// One theme with the edits applied; overridden is the paths the theme contexts override, null
// when they are not known.
function editedTheme(
	theme: DesignTheme,
	edits: Edits,
	overridden: ReadonlySet<string> | null
): DesignTheme {
	// The value edited for path in this theme, a light edit carried into dark where no context
	// overrides the token.
	const edit = (path: string): Edited | undefined => {
		const own = editKey(theme.name, path);
		const light = editKey('light', path);
		if (Object.hasOwn(edits, own)) return { value: edits[own] };
		const inherits = theme.name === 'dark' && overridden !== null && !overridden.has(path);
		if (inherits && Object.hasOwn(edits, light)) return { value: edits[light] };
		return undefined;
	};
	const byPath = new Map(theme.tokens.map((t) => [t.path, t] as const));
	// The edit of the first edited token along path's chain of aliases.
	const followed = (path: string): Edited | undefined => {
		const seen = new Set<string>();
		let t = byPath.get(path);
		while (t && !seen.has(t.path)) {
			seen.add(t.path);
			const e = edit(t.path);
			if (e || !t.ref) return e;
			t = byPath.get(t.ref);
		}
		return undefined;
	};
	const tokens = theme.tokens.map((t): DesignToken => {
		const e = followed(t.path);
		if (!e) return t;
		const next: DesignToken = { ...t, value: e.value, ...renderToken(t.path, t.type, e.value) };
		if (edit(t.path)) delete next.ref;
		return next;
	});
	return { ...theme, tokens };
}

// The token paths either theme context overrides: none for a plain token file, whose one value
// serves both themes; null when raw does not parse.
function contextPaths(raw: string): Set<string> | null {
	let doc: Sources;
	try {
		doc = sources(raw, scan(raw));
	} catch (e) {
		if (e instanceof ScanError) return null;
		throw e;
	}
	const out = new Set<string>();
	for (const src of [...doc.lightContext, ...doc.dark]) tokenPaths(src, [], out);
	return out;
}

// Adds every token path in a token tree to out.
function tokenPaths(node: JObject, prefix: string[], out: Set<string>): void {
	for (const m of node.members) {
		if (m.key.startsWith('$') || m.value.kind !== 'object') continue;
		const path = [...prefix, m.key];
		if (member(m.value, '$value')) out.add(path.join('.'));
		else tokenPaths(m.value, path, out);
	}
}

// ---- Writing the edits into the token file ----

type JObject = { kind: 'object'; start: number; end: number; members: JMember[] };
type JArray = { kind: 'array'; start: number; end: number; items: JNode[] };
type JScalar = { kind: 'scalar'; start: number; end: number };
type JNode = JObject | JArray | JScalar;
type JMember = { key: string; keyStart: number; value: JNode };

class ScanError extends Error {}

const LITERAL = /-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?|true|false|null/y;

// Reads raw as one JSON value, keeping where each value and each key starts and ends.
function scan(raw: string): JNode {
	let i = 0;
	const fail = (what: string): never => {
		const line = raw.slice(0, i).split('\n').length;
		throw new ScanError(
			`the token file does not parse at line ${line}: ${what}; fix its JSON, then reload the page`
		);
	};
	const space = () => {
		while (i < raw.length && ' \t\r\n'.includes(raw[i])) i++;
	};
	const string = (): string => {
		const start = i++;
		while (i < raw.length && raw[i] !== '"') i += raw[i] === '\\' ? 2 : 1;
		if (i >= raw.length) fail('a string is not closed');
		i++;
		try {
			return JSON.parse(raw.slice(start, i)) as string;
		} catch {
			return fail('a string holds a character or an escape JSON does not allow');
		}
	};
	function object(start: number): JObject {
		const members: JMember[] = [];
		i++;
		space();
		if (raw[i] === '}') return { kind: 'object', start, end: ++i, members };
		for (;;) {
			space();
			if (raw[i] !== '"') fail('an object key is not a string');
			const keyStart = i;
			const key = string();
			space();
			if (raw[i] !== ':') fail(`the key ${JSON.stringify(key)} has no colon after it`);
			i++;
			members.push({ key, keyStart, value: value() });
			space();
			if (raw[i] === '}') return { kind: 'object', start, end: ++i, members };
			if (raw[i] !== ',') fail('an object member is not followed by a comma or a closing brace');
			i++;
		}
	}
	function array(start: number): JArray {
		const items: JNode[] = [];
		i++;
		space();
		if (raw[i] === ']') return { kind: 'array', start, end: ++i, items };
		for (;;) {
			items.push(value());
			space();
			if (raw[i] === ']') return { kind: 'array', start, end: ++i, items };
			if (raw[i] !== ',') fail('an array item is not followed by a comma or a closing bracket');
			i++;
		}
	}
	function value(): JNode {
		space();
		const start = i;
		if (raw[i] === '{') return object(start);
		if (raw[i] === '[') return array(start);
		if (raw[i] === '"') string();
		else {
			LITERAL.lastIndex = i;
			if (!LITERAL.test(raw)) fail('a value is expected');
			i = LITERAL.lastIndex;
		}
		return { kind: 'scalar', start, end: i };
	}
	const root = value();
	space();
	if (i < raw.length) fail("there is more after the document's closing brace");
	return root;
}

// The value of an object's member key; undefined when n is no object or has no such member.
function member(n: JNode | undefined, key: string): JNode | undefined {
	return n?.kind === 'object' ? n.members.find((m) => m.key === key)?.value : undefined;
}

// The inline sources of an array of sources: its objects that are no $ref.
function inline(n: JNode | undefined): JObject[] {
	if (n?.kind !== 'array') return [];
	return n.items.filter((s): s is JObject => s.kind === 'object' && !member(s, '$ref'));
}

// Where a token file's values stand.
type Sources = {
	// the sources the light theme merges, in resolution order: the sets', then the light context's
	light: JObject[];
	lightContext: JObject[];
	dark: JObject[];
	// the dark context's array of sources; null for a plain token file, which has none
	darkArray: JArray | null;
	// the file's indent unit, its top-level members' indent, else two spaces
	unit: string;
};

// The token file's sources, a plain token file being its one light source.
function sources(raw: string, root: JNode): Sources {
	if (root.kind !== 'object')
		throw new ScanError('the token file is not an object of token groups; fix it and reload');
	const first = root.members.at(0);
	const unit = (first && indentAt(raw, first.keyStart)) || '  ';
	const order = member(root, 'resolutionOrder');
	if (!order) return { light: [root], lightContext: [], dark: [], darkArray: null, unit };
	const contexts = member(member(member(root, 'modifiers'), 'theme'), 'contexts');
	const lightContext = inline(member(contexts, 'light'));
	const dark = member(contexts, 'dark');
	const darkArray = dark?.kind === 'array' ? dark : null;
	const light: JObject[] = [];
	for (const step of order.kind === 'array' ? order.items : []) {
		const ref = member(step, '$ref');
		if (ref?.kind !== 'scalar') continue;
		const text: unknown = JSON.parse(raw.slice(ref.start, ref.end));
		if (text === '#/modifiers/theme') light.push(...lightContext);
		if (typeof text !== 'string' || !text.startsWith('#/sets/')) continue;
		const name = text.slice(7).replaceAll('~1', '/').replaceAll('~0', '~');
		light.push(...inline(member(member(member(root, 'sets'), name), 'sources')));
	}
	return { light, lightContext, dark: inline(dark), darkArray, unit };
}

// The $value of the token at path in the last of srcs that gives it.
function lastValue(srcs: JObject[], path: string[]): JNode | undefined {
	let found: JNode | undefined;
	for (const src of srcs) {
		const token = path.reduce<JNode | undefined>((n, key) => member(n, key), src);
		found = member(token, '$value') ?? found;
	}
	return found;
}

// The leading white space of the line offset stands on.
function indentAt(raw: string, offset: number): string {
	const start = raw.lastIndexOf('\n', offset - 1) + 1;
	let end = start;
	while (raw[end] === ' ' || raw[end] === '\t') end++;
	return raw.slice(start, end);
}

// New members to write into one object or array: group names down to each value's JSON.
type Tree = Map<string, Tree | string>;

type Change = { start: number; end: number; text: string };

const replace = (n: JNode, text: string): Change => ({ start: n.start, end: n.end, text });

const insert = (at: number, text: string): Change => ({ start: at, end: at, text });

// Tree's members as text: a token on one line, and a group on as many lines as its members when
// multiline, indent being the indent of the line it starts on.
function memberLines(tree: Tree, indent: string, unit: string, multiline: boolean): string[] {
	return [...tree].map(([key, node]) => {
		const k = JSON.stringify(key);
		if (typeof node === 'string') return `${k}: { "$value": ${node} }`;
		const inner = memberLines(node, indent + unit, unit, multiline);
		if (!multiline) return `${k}: { ${inner.join(', ')} }`;
		return `${k}: {\n${inner.map((m) => indent + unit + m).join(',\n')}\n${indent}}`;
	});
}

// A new source holding tree's members, its closing brace at indent at.
function source(tree: Tree, at: string, unit: string): string {
	const lines = memberLines(tree, at + unit, unit, true).map((l) => at + unit + l);
	return `{\n${lines.join(',\n')}\n${at}}`;
}

/**
 * The token file raw with each edit's $value written in place as compact JSON, and every other
 * byte as it was: a light edit where the light theme takes the token from, a dark edit where the
 * dark context overrides the token, or, where it does not yet, added to the context's first
 * source with the groups its path needs. raw itself when there are no edits. Throws an Error,
 * naming the token and what to do, when no set defines an edited token, or a plain token file is
 * given a dark value.
 */
export function composeTokenFile(raw: string, edits: Edits): string {
	const keys = Object.keys(edits);
	if (keys.length === 0) return raw;
	const doc = sources(raw, scan(raw));
	const changes: Change[] = [];
	const plan = new Map<JObject | JArray, Tree>();
	for (const key of keys) {
		const k = parseEditKey(key);
		if (!k) throw new Error(`${JSON.stringify(key)} is no edit key, theme:path; discard it`);
		const text = JSON.stringify(edits[key]);
		if (text === undefined) throw new Error(`the edit of ${key} is no JSON value; discard it`);
		const path = k.path.split('.');
		if (!lastValue(doc.light, path))
			throw new Error(
				`no set of the token file defines the token ${k.path}, so its ${k.theme} value has nowhere to go; reload the page to read the token file again, or discard the edit`
			);
		const target = lastValue(k.theme === 'light' ? doc.light : doc.dark, path);
		if (target) changes.push(replace(target, text));
		else planOverride(plan, doc, k.path, text);
	}
	for (const [container, tree] of plan) changes.push(insertion(raw, container, tree, doc.unit));
	changes.sort((a, b) => b.start - a.start);
	return changes.reduce((out, c) => out.slice(0, c.start) + c.text + out.slice(c.end), raw);
}

// Plans a dark value of a token the dark context does not override yet: into the context's first
// source, under the groups its path needs, or into a new source when it has none.
function planOverride(
	plan: Map<JObject | JArray, Tree>,
	doc: Sources,
	dotted: string,
	text: string
): void {
	if (!doc.darkArray)
		throw new Error(
			`the token file is a plain token file, whose one value of ${dotted} serves both themes; edit its light value, or make the file a resolver document with a dark context`
		);
	const path = dotted.split('.');
	let container: JObject | JArray = doc.dark.at(0) ?? doc.darkArray;
	let depth = 0;
	for (; container.kind === 'object' && depth < path.length; depth++) {
		const m = container.members.find((mm) => mm.key === path[depth]);
		if (!m) break;
		if (m.value.kind !== 'object' || member(m.value, '$value') || depth === path.length - 1)
			throw new Error(
				`the dark context's ${path.slice(0, depth + 1).join('.')} is not a group of tokens, so ${dotted} cannot be added under it; correct the token file`
			);
		container = m.value;
	}
	let tree: Tree = plan.get(container) ?? new Map();
	plan.set(container, tree);
	for (const name of path.slice(depth, -1)) {
		const sub = tree.get(name);
		const next: Tree = typeof sub === 'object' ? sub : new Map();
		tree.set(name, next);
		tree = next;
	}
	tree.set(path[path.length - 1], text);
}

// The change that adds tree's members to container, formatted as its neighbours are.
function insertion(raw: string, container: JObject | JArray, tree: Tree, unit: string): Change {
	const indent = indentAt(raw, container.start);
	if (container.kind === 'array') {
		const item = container.items.at(-1);
		const inner = indent + unit;
		if (!item) return replace(container, `[\n${inner}${source(tree, inner, unit)}\n${indent}]`);
		const at = indentAt(raw, item.start);
		return insert(item.end, `,\n${at}${source(tree, at, unit)}`);
	}
	const last = container.members.at(-1);
	if (!last) return replace(container, source(tree, indent, unit));
	const first = container.members[0];
	if (!raw.slice(container.start + 1, first.keyStart).includes('\n')) {
		const lines = memberLines(tree, '', '', false).map((l) => ', ' + l);
		return insert(last.value.end, lines.join(''));
	}
	const inner = indentAt(raw, first.keyStart);
	const nested = inner.startsWith(indent) ? inner.slice(indent.length) : '';
	const lines = memberLines(tree, inner, nested || unit, true).map((l) => `,\n${inner}${l}`);
	return insert(last.value.end, lines.join(''));
}
