// The props panel's model (S-0440): the control each prop of the picked component takes, by the
// kind flai read for it (S-0437), what the text in that control means, and the props in the
// page's URL. The panel keeps each prop's control as text, so that what the operator types stays
// until it parses; only what parses is sent, and a function prop gets the harness's logger.

import type { ComponentProp } from '$lib/server/repo';
import type { Props } from './protocol';

/** One of a union's literals as its select offers it: what it shows, and the value's JSON. */
export type Option = { label: string; text: string };

/** The control a prop takes. A snippet's text is rendered as the snippet. */
export type Control =
	| { kind: 'text' }
	| { kind: 'number' }
	| { kind: 'checkbox' }
	| { kind: 'select'; options: Option[] }
	| { kind: 'json' }
	| { kind: 'snippet' };

/** What a control's text means: its value, undefined for none, or why it is not one. */
export type Parsed = { ok: true; value: unknown } | { ok: false; error: string };

/** The page's URL parameter that holds the picked component's props. */
export const PROPS_PARAM = 'props';

/** The control for prop, by its kind; null for a function, which the harness logs instead. */
export function controlFor(prop: ComponentProp): Control | null {
	switch (prop.kind) {
		case 'string':
			return { kind: 'text' };
		case 'number':
			return { kind: 'number' };
		case 'boolean':
			return { kind: 'checkbox' };
		case 'union':
			return { kind: 'select', options: options(prop) };
		case 'snippet':
			return { kind: 'snippet' };
		case 'array':
		case 'object':
		case 'unknown':
			return { kind: 'json' };
		case 'function':
			return null;
	}
}

// A union's literals, each a value as JSON or as the source writes it; one that is not is left out.
function options(prop: ComponentProp): Option[] {
	return (prop.literals ?? []).flatMap((l) => {
		const v = literal(l);
		if (!v) return [];
		const text = JSON.stringify(v.value);
		return [{ label: typeof v.value === 'string' ? v.value : text, text }];
	});
}

const ok = (value: unknown): Parsed => ({ ok: true, value });
const fail = (error: string): Parsed => ({ ok: false, error });

/** What text in prop's control means. No text is no value, and so is blank text, except in a
 * text or snippet control, where it is the empty string. */
export function parse(prop: ComponentProp, text: string | undefined): Parsed {
	const control = controlFor(prop);
	if (text === undefined || !control) return ok(undefined);
	if (control.kind === 'text' || control.kind === 'snippet') return ok(text);
	if (text.trim() === '') return ok(undefined);
	switch (control.kind) {
		case 'number': {
			const n = Number(text);
			return Number.isFinite(n) ? ok(n) : fail('Not a number');
		}
		case 'checkbox':
			if (text === 'true' || text === 'false') return ok(text === 'true');
			return fail('Not true or false');
		case 'select':
			if (control.options.some((o) => o.text === text)) return ok(JSON.parse(text));
			return fail(`Not one of ${control.options.map((o) => o.label).join(', ')}`);
		case 'json':
			try {
				return ok(JSON.parse(text));
			} catch (e) {
				return fail(`Invalid JSON: ${e instanceof Error ? e.message : String(e)}`);
			}
	}
}

/** The text prop's control shows for value, which parse reads back; '' for undefined. */
export function format(prop: ComponentProp, value: unknown): string {
	if (value === undefined) return '';
	switch (prop.kind) {
		case 'string':
		case 'snippet':
			return typeof value === 'string' ? value : JSON.stringify(value);
		case 'number':
		case 'boolean':
			return String(value);
		case 'union':
			return JSON.stringify(value);
		case 'array':
		case 'object':
		case 'unknown':
			return JSON.stringify(value, null, 2);
		case 'function':
			return '';
	}
}

/** prop's default as a value, where the source writes it as a literal its control can hold;
 * undefined where it has none, writes an expression, or writes a value of another kind. */
export function initial(prop: ComponentProp): unknown {
	const v = prop.default ? literal(prop.default) : null;
	return v && fits(prop, v.value) ? v.value : undefined;
}

function fits(prop: ComponentProp, value: unknown): boolean {
	switch (prop.kind) {
		case 'string':
		case 'snippet':
			return typeof value === 'string';
		case 'number':
			return Number.isFinite(value);
		case 'boolean':
			return typeof value === 'boolean';
		case 'union':
			return options(prop).some((o) => o.text === JSON.stringify(value));
		case 'array':
		case 'object':
		case 'unknown':
			return true;
		case 'function':
			return false;
	}
}

/** Whether prop is required and text gives it no value, so that the panel marks it. A function
 * never is: the harness gives it a logger. */
export function missing(prop: ComponentProp, text: string | undefined): boolean {
	if (!prop.required || prop.kind === 'function') return false;
	const p = parse(prop, text);
	return !p.ok || p.value === undefined;
}

/** The names of props' function props, for the harness to log. */
export function functions(props: ComponentProp[]): string[] {
	return props.filter((p) => p.kind === 'function').map((p) => p.name);
}

/** Each prop's text: its value in values where it has one there, else its default where initial
 * reads one. A prop with neither, and a function, has no text. With no values, the reset. */
export function toTexts(props: ComponentProp[], values: Props = {}): Record<string, string> {
	const texts: Record<string, string> = {};
	for (const prop of props) {
		if (prop.kind === 'function') continue;
		const value = Object.hasOwn(values, prop.name) ? values[prop.name] : initial(prop);
		if (value !== undefined) texts[prop.name] = format(prop, value);
	}
	return texts;
}

/** The props to send for texts: each prop whose text parses to a value. A prop with no value, or
 * with text that does not parse, is left out, and so is a function. */
export function toValues(props: ComponentProp[], texts: Record<string, string | undefined>): Props {
	const values: Props = {};
	for (const prop of props) {
		const p = parse(prop, texts[prop.name]);
		if (p.ok && p.value !== undefined) values[prop.name] = p.value;
	}
	return values;
}

/** values as PROPS_PARAM's value: their JSON, base64url-encoded without padding; '' for none. */
export function toQuery(values: Props): string {
	const json = JSON.stringify(values);
	if (json === '{}') return '';
	let binary = '';
	for (const b of new TextEncoder().encode(json)) binary += String.fromCharCode(b);
	return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

/** The props in PROPS_PARAM's value; {} when it is absent or is not an object toQuery wrote. */
export function fromQuery(param: string | null | undefined): Props {
	if (!param) return {};
	try {
		const binary = atob(param.replace(/-/g, '+').replace(/_/g, '/'));
		const bytes = Uint8Array.from(binary, (c) => c.charCodeAt(0));
		const v: unknown = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
		return typeof v === 'object' && v !== null && !Array.isArray(v) ? (v as Props) : {};
	} catch {
		return {};
	}
}

const NUMBER = /[+-]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?/y;
const WORD = /[A-Za-z_$][\w$]*/y;
const ESCAPES: Record<string, string> = {
	n: '\n',
	t: '\t',
	r: '\r',
	b: '\b',
	f: '\f',
	v: '\v',
	0: '\0'
};

// source as a value where it is a literal as JavaScript writes one: a string, a number, true,
// false, or null, or an array or object of literals, with either quote, bare keys, and trailing
// commas; null for any other expression, such as a call, a name, or a template literal.
function literal(source: string): { value: unknown } | null {
	let i = 0;
	const stop = (): never => {
		throw new SyntaxError(source);
	};
	const skip = () => {
		while (/\s/.test(source[i] ?? '')) i++;
	};
	const eat = (c: string) => {
		skip();
		if (source[i] !== c) return false;
		i++;
		return true;
	};
	const match = (re: RegExp) => {
		skip();
		re.lastIndex = i;
		const m = re.exec(source);
		if (m) i = re.lastIndex;
		return m?.[0];
	};

	function string(): string {
		const quote = source[i++];
		let s = '';
		while (source[i] !== quote) {
			const c = source[i++];
			if (c === undefined || c === '\n') stop();
			if (c !== '\\') {
				s += c;
				continue;
			}
			const e = source[i++];
			if (e === undefined) stop();
			if (e !== 'u' && e !== 'x') {
				s += ESCAPES[e] ?? e;
				continue;
			}
			const size = e === 'u' ? 4 : 2;
			const hex = source.slice(i, i + size);
			if (hex.length !== size || !/^[\da-fA-F]+$/.test(hex)) stop();
			s += String.fromCharCode(parseInt(hex, 16));
			i += size;
		}
		i++;
		return s;
	}

	function list<T>(close: string, item: () => T): T[] {
		const items: T[] = [];
		while (!eat(close)) {
			items.push(item());
			if (eat(',')) continue;
			if (!eat(close)) stop();
			break;
		}
		return items;
	}

	function key(): string {
		skip();
		if (source[i] === '"' || source[i] === "'") return string();
		return match(WORD) ?? match(NUMBER) ?? stop();
	}

	function value(): unknown {
		skip();
		if (source[i] === '"' || source[i] === "'") return string();
		if (eat('[')) return list(']', value);
		if (eat('{')) {
			const entry = (): [string, unknown] => {
				const k = key();
				if (!eat(':')) stop();
				return [k, value()];
			};
			return Object.fromEntries(list('}', entry));
		}
		const n = match(NUMBER);
		if (n !== undefined) return Number(n);
		const w = match(WORD);
		if (w === 'true') return true;
		if (w === 'false') return false;
		if (w === 'null') return null;
		return stop();
	}

	try {
		const v = value();
		skip();
		return i === source.length ? { value: v } : null;
	} catch {
		return null;
	}
}
