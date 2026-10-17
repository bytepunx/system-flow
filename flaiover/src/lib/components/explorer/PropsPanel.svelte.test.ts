// S-0440: the props panel lists each prop with its type, its default, and a control by its kind,
// shows text that does not parse with why, marks a required prop without a value, hands each
// change up whole, and resets every prop to its default.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { ComponentProp, ComponentPropKind } from '$lib/server/repo';
import PropsPanel from './PropsPanel.svelte';

const prop = (name: string, kind: ComponentPropKind, more: Partial<ComponentProp> = {}) =>
	({ name, kind, required: false, ...more }) as ComponentProp;

const PROPS: ComponentProp[] = [
	prop('label', 'string', { type: 'string', required: true }),
	prop('count', 'number', { type: 'number', default: '3' }),
	prop('open', 'boolean', { type: 'boolean', default: 'false' }),
	prop('variant', 'union', {
		type: "'primary' | 'ghost'",
		default: "'primary'",
		literals: ['"primary"', '"ghost"']
	}),
	prop('items', 'array', { type: 'string[]', default: "['a']" }),
	prop('options', 'object', { type: '{ dense: boolean }' }),
	prop('data', 'unknown'),
	prop('children', 'snippet', { type: 'Snippet' }),
	prop('onclick', 'function', { type: '() => void' })
];

describe('PropsPanel (S-0440)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	const q = <T extends Element = HTMLElement>(sel: string) => document.querySelector<T>(sel);
	const row = (name: string) => q(`[data-testid="prop-${name}"]`)!;
	const control = <T extends Element>(name: string) => q<T>(`#props-panel-${name}`)!;
	const inRow = (name: string, testid: string) =>
		row(name).querySelector(`[data-testid="${testid}"]`);

	/** Mounts the panel over texts that take each change, as the page does; onchange records them. */
	function show(texts: Record<string, string> = {}, props = PROPS) {
		const onchange = vi.fn();
		const state = $state({ texts });
		c = mount(PropsPanel, {
			target: document.body,
			props: {
				props,
				get texts() {
					return state.texts;
				},
				onchange: (next: Record<string, string>) => {
					onchange(next);
					state.texts = next;
				}
			}
		});
		flushSync();
		return { onchange, state };
	}
	const type = (el: HTMLInputElement | HTMLTextAreaElement, text: string) => {
		el.value = text;
		el.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
	};
	const choose = (el: HTMLSelectElement | HTMLInputElement) => {
		el.dispatchEvent(new Event('change', { bubbles: true }));
		flushSync();
	};

	it('lists each prop with its type and its default as written, and a control by its kind', () => {
		show();
		const types = PROPS.map((p) => inRow(p.name, 'prop-type')!.textContent);
		expect(types).toEqual([
			'string',
			'number',
			'boolean',
			"'primary' | 'ghost'",
			'string[]',
			'{ dense: boolean }',
			'unknown',
			'Snippet',
			'() => void'
		]);
		expect(inRow('count', 'prop-default')!.textContent).toBe('= 3');
		expect(inRow('variant', 'prop-default')!.textContent).toBe("= 'primary'");
		expect(inRow('label', 'prop-default')).toBeNull();

		const kind = (name: string) => {
			const el = control<HTMLElement>(name);
			return el instanceof HTMLInputElement ? `input:${el.type}` : el.tagName.toLowerCase();
		};
		expect(Object.fromEntries(PROPS.map((p) => [p.name, kind(p.name)]))).toEqual({
			label: 'input:text',
			count: 'input:number',
			open: 'input:checkbox',
			variant: 'select',
			items: 'textarea',
			options: 'textarea',
			data: 'textarea',
			children: 'textarea',
			onclick: 'p'
		});
		// each control is labelled with its prop's name
		expect(q(`label[for="props-panel-count"]`)!.textContent).toBe('count');
		const offered = [...control<HTMLSelectElement>('variant').options].map((o) => [
			o.value,
			o.textContent
		]);
		expect(offered).toEqual([
			['', '(not set)'],
			['"primary"', 'primary'],
			['"ghost"', 'ghost']
		]);
		expect(inRow('onclick', 'prop-logged')!.textContent).toContain('Logged');
	});

	it('shows the texts it is given in their controls', () => {
		show({ label: 'Go', count: '3', open: 'true', variant: '"ghost"', items: '["a"]' });
		expect(control<HTMLInputElement>('label').value).toBe('Go');
		expect(control<HTMLInputElement>('count').value).toBe('3');
		expect(control<HTMLInputElement>('open').checked).toBe(true);
		expect(control<HTMLSelectElement>('variant').value).toBe('"ghost"');
		expect(control<HTMLTextAreaElement>('items').value).toBe('["a"]');
	});

	it('hands each change up with every text, from each kind of control', () => {
		const { onchange } = show({ count: '3' });
		type(control('label'), 'Go');
		expect(onchange).toHaveBeenLastCalledWith({ count: '3', label: 'Go' });
		type(control('count'), '7');
		expect(onchange).toHaveBeenLastCalledWith({ count: '7', label: 'Go' });
		const open = control<HTMLInputElement>('open');
		open.checked = true;
		choose(open);
		expect(onchange.mock.lastCall![0]).toMatchObject({ open: 'true' });
		const variant = control<HTMLSelectElement>('variant');
		variant.value = '"ghost"';
		choose(variant);
		expect(onchange.mock.lastCall![0]).toMatchObject({ variant: '"ghost"' });
		type(control('options'), '{"dense": true}');
		expect(onchange.mock.lastCall![0]).toMatchObject({ options: '{"dense": true}' });
		type(control('children'), 'Hello');
		expect(onchange.mock.lastCall![0]).toMatchObject({ children: 'Hello' });
	});

	it('shows invalid JSON in the editor with why, as not sent', () => {
		show({ label: 'Go' });
		type(control('items'), '["a",');
		expect(control<HTMLTextAreaElement>('items').value).toBe('["a",');
		expect(control('items').getAttribute('aria-invalid')).toBe('true');
		const error = inRow('items', 'prop-error')!;
		expect(error.textContent).toMatch(/Invalid JSON.*not sent/s);
		expect(control('items').getAttribute('aria-describedby')).toBe(error.id);
		type(control('items'), '["a", "b"]');
		expect(inRow('items', 'prop-error')).toBeNull();
		expect(control('items').getAttribute('aria-invalid')).toBe('false');
	});

	it('marks a required prop without a value until it has one', () => {
		show();
		expect(inRow('label', 'prop-required')).not.toBeNull();
		expect(inRow('count', 'prop-required')).toBeNull();
		expect(inRow('label', 'prop-missing')).not.toBeNull();
		type(control('label'), 'Go');
		expect(inRow('label', 'prop-missing')).toBeNull();
	});

	it('never marks a required function as without a value, as the harness logs it', () => {
		show({}, [prop('onclick', 'function', { required: true })]);
		expect(q('[data-testid="prop-missing"]')).toBeNull();
	});

	it('resets every prop to its default', () => {
		const { onchange } = show({ label: 'Go', count: '9', variant: '"ghost"', items: '[' });
		q<HTMLButtonElement>('[data-testid="props-reset"]')!.click();
		flushSync();
		expect(onchange).toHaveBeenLastCalledWith({
			count: '3',
			open: 'false',
			variant: '"primary"',
			items: '[\n  "a"\n]'
		});
		expect(control<HTMLInputElement>('label').value).toBe('');
		expect(control<HTMLInputElement>('count').value).toBe('3');
		expect(control<HTMLSelectElement>('variant').value).toBe('"primary"');
	});

	it('says when the component takes no props, and offers no reset', () => {
		show({}, []);
		expect(q('[data-testid="props-none"]')).not.toBeNull();
		expect(q('[data-testid="props-reset"]')).toBeNull();
	});
});
