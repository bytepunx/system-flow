import { describe, expect, it } from 'vitest';
import type { ComponentProp, ComponentPropKind } from '$lib/server/repo';
import {
	PROPS_PARAM,
	controlFor,
	format,
	fromQuery,
	functions,
	initial,
	missing,
	parse,
	toQuery,
	toTexts,
	toValues
} from './controls';

// S-0440: the props panel's model, against props as flai components list answers them.
const prop = (kind: ComponentPropKind, more: Partial<ComponentProp> = {}): ComponentProp => ({
	name: kind,
	kind,
	required: false,
	...more
});
const VARIANT = prop('union', { literals: ['"primary"', '"ghost"', '1', 'true'] });

describe('controlFor', () => {
	it('gives each kind its control, and a function none', () => {
		expect(controlFor(prop('string'))).toEqual({ kind: 'text' });
		expect(controlFor(prop('number'))).toEqual({ kind: 'number' });
		expect(controlFor(prop('boolean'))).toEqual({ kind: 'checkbox' });
		expect(controlFor(prop('array'))).toEqual({ kind: 'json' });
		expect(controlFor(prop('object'))).toEqual({ kind: 'json' });
		expect(controlFor(prop('unknown'))).toEqual({ kind: 'json' });
		expect(controlFor(prop('snippet'))).toEqual({ kind: 'snippet' });
		expect(controlFor(prop('function'))).toBeNull();
	});

	it('offers a union its literals, a string by itself and others as JSON', () => {
		expect(controlFor(VARIANT)).toEqual({
			kind: 'select',
			options: [
				{ label: 'primary', text: '"primary"' },
				{ label: 'ghost', text: '"ghost"' },
				{ label: '1', text: '1' },
				{ label: 'true', text: 'true' }
			]
		});
	});

	it('reads a literal written with single quotes', () => {
		expect(controlFor(prop('union', { literals: ["'a'"] }))).toEqual({
			kind: 'select',
			options: [{ label: 'a', text: '"a"' }]
		});
	});
});

describe('parse', () => {
	it('keeps text and snippet text as written, blank included', () => {
		expect(parse(prop('string'), 'Go ')).toEqual({ ok: true, value: 'Go ' });
		expect(parse(prop('string'), '')).toEqual({ ok: true, value: '' });
		expect(parse(prop('snippet'), '<b>hi</b>')).toEqual({ ok: true, value: '<b>hi</b>' });
	});

	it('reads a number, and refuses what is not one', () => {
		expect(parse(prop('number'), '-2.5')).toEqual({ ok: true, value: -2.5 });
		expect(parse(prop('number'), 'ten')).toEqual({ ok: false, error: 'Not a number' });
	});

	it('reads a checkbox as true or false', () => {
		expect(parse(prop('boolean'), 'true')).toEqual({ ok: true, value: true });
		expect(parse(prop('boolean'), 'false')).toEqual({ ok: true, value: false });
		expect(parse(prop('boolean'), 'yes')).toEqual({ ok: false, error: 'Not true or false' });
	});

	it('reads a union literal as its value, and refuses one not among them', () => {
		expect(parse(VARIANT, '"primary"')).toEqual({ ok: true, value: 'primary' });
		expect(parse(VARIANT, '1')).toEqual({ ok: true, value: 1 });
		expect(parse(VARIANT, 'true')).toEqual({ ok: true, value: true });
		expect(parse(VARIANT, '"other"')).toEqual({
			ok: false,
			error: 'Not one of primary, ghost, 1, true'
		});
	});

	it('reads JSON for an array, an object, or an unknown, and says why invalid JSON is not', () => {
		expect(parse(prop('array'), '[1, "a"]')).toEqual({ ok: true, value: [1, 'a'] });
		expect(parse(prop('object'), '{"a": {"b": null}}')).toEqual({
			ok: true,
			value: { a: { b: null } }
		});
		expect(parse(prop('unknown'), '"x"')).toEqual({ ok: true, value: 'x' });
		const bad = parse(prop('object'), '{a: 1}');
		expect(bad.ok).toBe(false);
		expect(!bad.ok && bad.error).toMatch(/^Invalid JSON: /);
	});

	it('gives no value for no text, blank text outside a text control, and a function', () => {
		for (const kind of ['string', 'number', 'boolean', 'union', 'array', 'snippet'] as const) {
			expect(parse(prop(kind), undefined)).toEqual({ ok: true, value: undefined });
		}
		for (const kind of ['number', 'boolean', 'union', 'array', 'object', 'unknown'] as const) {
			expect(parse(prop(kind), '  ')).toEqual({ ok: true, value: undefined });
		}
		expect(parse(prop('function'), '() => {}')).toEqual({ ok: true, value: undefined });
	});
});

describe('format', () => {
	it('writes each kind as its control shows it, which parse reads back', () => {
		const cases: [ComponentProp, unknown, string][] = [
			[prop('string'), 'Go', 'Go'],
			[prop('snippet'), 'hi', 'hi'],
			[prop('number'), 3, '3'],
			[prop('boolean'), false, 'false'],
			[VARIANT, 'ghost', '"ghost"'],
			[prop('array'), [1], '[\n  1\n]'],
			[prop('object'), { a: 1 }, '{\n  "a": 1\n}']
		];
		for (const [p, value, text] of cases) {
			expect(format(p, value)).toBe(text);
			expect(parse(p, text)).toEqual({ ok: true, value });
		}
		expect(format(prop('number'), undefined)).toBe('');
	});
});

describe('initial', () => {
	it('reads a default written as a literal', () => {
		expect(initial(prop('string', { default: "'primary'" }))).toBe('primary');
		expect(initial(prop('string', { default: '"it\\\'s \\u0041"' }))).toBe("it's A");
		expect(initial(prop('number', { default: '3' }))).toBe(3);
		expect(initial(prop('number', { default: '-1.5e2' }))).toBe(-150);
		expect(initial(prop('boolean', { default: 'true' }))).toBe(true);
		expect(initial({ ...VARIANT, default: "'ghost'" })).toBe('ghost');
		expect(initial(prop('array', { default: '[]' }))).toEqual([]);
		expect(initial(prop('array', { default: "['a', 2,]" }))).toEqual(['a', 2]);
		expect(initial(prop('object', { default: "{a: 1, 'b': [null], 2: {}, }" }))).toEqual({
			a: 1,
			b: [null],
			2: {}
		});
		expect(initial(prop('unknown', { default: 'null' }))).toBeNull();
	});

	it('gives none for no default, an expression, or a value the control cannot hold', () => {
		const none = [
			prop('string'),
			prop('string', { default: 'foo()' }),
			prop('string', { default: 'undefined' }),
			prop('string', { default: '`x`' }),
			prop('string', { default: "'a' + 'b'" }),
			prop('array', { default: '[a]' }),
			prop('object', { default: '{a}' }),
			prop('object', { default: '{...rest}' }),
			prop('string', { default: "'open" }),
			prop('number', { default: "'3'" }),
			prop('number', { default: 'null' }),
			prop('boolean', { default: '1' }),
			{ ...VARIANT, default: "'other'" },
			prop('function', { default: '() => {}' })
		];
		for (const p of none) expect(initial(p), p.default).toBeUndefined();
	});
});

describe('missing', () => {
	it('marks a required prop with no value or with text that does not parse', () => {
		const label = prop('string', { name: 'label', required: true });
		const items = prop('array', { required: true });
		expect(missing(label, undefined)).toBe(true);
		expect(missing(label, '')).toBe(false);
		expect(missing(items, '')).toBe(true);
		expect(missing(items, '[')).toBe(true);
		expect(missing(items, '[]')).toBe(false);
		expect(missing(prop('string'), undefined)).toBe(false);
		expect(missing(prop('function', { required: true }), undefined)).toBe(false);
	});
});

describe('the panel', () => {
	const props: ComponentProp[] = [
		prop('string', { name: 'label', required: true }),
		{ ...VARIANT, name: 'variant', default: "'primary'" },
		prop('number', { name: 'size', default: 'compute()' }),
		prop('array', { name: 'items', default: '[]' }),
		prop('function', { name: 'onclick' }),
		prop('function', { name: 'onhover' })
	];

	it('names the function props for the harness to log', () => {
		expect(functions(props)).toEqual(['onclick', 'onhover']);
	});

	it('resets each prop to its literal default, and leaves the rest without text', () => {
		expect(toTexts(props)).toEqual({ variant: '"primary"', items: '[]' });
	});

	it('takes values given over defaults', () => {
		expect(toTexts(props, { label: 'Go', size: 2, onclick: 1 })).toEqual({
			label: 'Go',
			variant: '"primary"',
			size: '2',
			items: '[]'
		});
	});

	it('sends what parses, and leaves out no value, invalid JSON, and functions', () => {
		const texts = { label: 'Go', variant: '"ghost"', size: '', items: '[1,', onclick: 'x' };
		expect(toValues(props, texts)).toEqual({ label: 'Go', variant: 'ghost' });
	});
});

describe('the URL', () => {
	it('holds the props under one parameter that the page does not use already', () => {
		expect(PROPS_PARAM).toBe('props');
	});

	it('round trips props as base64url JSON without padding', () => {
		const values = { label: 'Gö ✓ ~?>', n: -1.5, on: true, items: [{ a: null }], s: '' };
		const query = toQuery(values);
		expect(query).toMatch(/^[A-Za-z0-9_-]+$/);
		expect(JSON.parse(Buffer.from(query, 'base64url').toString('utf8'))).toEqual(values);
		expect(fromQuery(query)).toEqual(values);
		const url = new URL('http://127.0.0.1/design/components?p=web&c=A.svelte');
		url.searchParams.set(PROPS_PARAM, query);
		expect(fromQuery(new URL(url.href).searchParams.get(PROPS_PARAM))).toEqual(values);
	});

	it('writes nothing for no props', () => {
		expect(toQuery({})).toBe('');
		expect(toQuery({ a: undefined })).toBe('');
	});

	it('reads anything else as no props, and never throws', () => {
		const encode = (s: string) => Buffer.from(s, 'utf8').toString('base64url');
		for (const bad of [
			null,
			undefined,
			'',
			'!!!',
			'a',
			encode('{"a":'),
			encode('[1]'),
			encode('null'),
			encode('3'),
			Buffer.from([0x7b, 0xff, 0x7d]).toString('base64url')
		]) {
			expect(fromQuery(bad), String(bad)).toEqual({});
		}
	});
});
