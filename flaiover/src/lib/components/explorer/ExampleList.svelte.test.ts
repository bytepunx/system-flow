// S-0443: the example list shows the picked component's examples, hands up the one chosen, and
// saves the panel as a named example, showing a conflict or a refusal with its findings and keeping
// the name to save again.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount, type ComponentProps } from 'svelte';
import type { ComponentExample } from '$lib/server/repo';
import ExampleList, { exampleFrom, type SaveFailure } from './ExampleList.svelte';

const PATH = 'design/components/flaiover/src/lib/components/Card.examples.json';
const LONG: ComponentExample = {
	name: 'Long title',
	props: { title: 'A long one' },
	theme: 'dark',
	width: 375
};
const EMPTY: ComponentExample = { name: 'Empty', props: {} };

describe('exampleFrom', () => {
	it('keeps the theme when one is shown and the width when one is picked', () => {
		expect(exampleFrom(' Long ', { title: 'x' }, ['dark'], 375)).toEqual({
			name: 'Long',
			props: { title: 'x' },
			theme: 'dark',
			width: 375
		});
		expect(exampleFrom('Both', {}, ['light', 'dark'], null)).toEqual({ name: 'Both', props: {} });
	});
});

describe('ExampleList (S-0443)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	const q = (testid: string) => document.querySelector(`[data-testid="${testid}"]`);
	const settle = async () => {
		for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
		flushSync();
	};
	const click = async (testid: string) => {
		(q(testid) as HTMLElement).click();
		await settle();
	};
	const type = (value: string) => {
		const input = q('example-name') as HTMLInputElement;
		input.value = value;
		input.dispatchEvent(new Event('input'));
		flushSync();
	};
	const submit = async () => {
		(q('example-save-form') as HTMLFormElement).requestSubmit();
		await settle();
	};

	function show(
		more: Partial<ComponentProps<typeof ExampleList>> = {},
		answer: SaveFailure | null = null
	) {
		const onchoose = vi.fn();
		const onsave = vi.fn(async () => answer);
		c = mount(ExampleList, {
			target: document.body,
			props: { examples: [EMPTY, LONG], path: PATH, onchoose, onsave, ...more }
		});
		flushSync();
		return { onchoose, onsave };
	}

	it('says there are none yet', () => {
		show({ examples: [] });
		expect(q('example-list-none')).not.toBeNull();
	});

	it('lists the examples in order and hands up the one chosen, marked current', async () => {
		const { onchoose } = show();
		const names = [...document.querySelectorAll('li button')].map((b) => b.textContent);
		expect(names).toEqual(['Empty', 'Long title']);
		await click('example-Long title');
		expect(onchoose).toHaveBeenCalledWith(LONG);
		expect(q('example-Long title')!.getAttribute('aria-current')).toBe('true');
		expect(q('example-Empty')!.getAttribute('aria-current')).toBeNull();
	});

	it('asks for a name, saves under it, and says where it goes', async () => {
		const { onsave } = show();
		expect(q('example-save-form')).toBeNull();
		await click('example-save-as');
		expect(q('example-save-form')!.textContent).toContain(PATH);
		expect(q('example-save-form')!.textContent).toContain('on main');
		expect((q('example-save') as HTMLButtonElement).disabled).toBe(true);
		type('  Short  ');
		await submit();
		expect(onsave).toHaveBeenCalledWith('Short');
		expect(q('example-save-form')).toBeNull();
		expect(q('example-saved')!.textContent).toContain('Short');
	});

	it("offers the chosen example's name, says the save replaces it, and names the story", async () => {
		show({ story: 'S-0443' });
		await click('example-Empty');
		await click('example-save-as');
		expect((q('example-name') as HTMLInputElement).value).toBe('Empty');
		expect(q('example-replaces')).not.toBeNull();
		expect(q('example-save-form')!.textContent).toContain('story/S-0443');
		type('New');
		expect(q('example-replaces')).toBeNull();
	});

	it('shows a conflict and keeps the name to save again', async () => {
		show({}, { error: 'Card.examples.json changed after it was loaded', conflict: true });
		await click('example-save-as');
		type('Mine');
		await submit();
		const alert = q('example-save-error')!;
		expect(alert.textContent).toContain('changed after it was loaded');
		expect(alert.textContent).toContain('save again');
		expect((q('example-name') as HTMLInputElement).value).toBe('Mine');
		expect(q('example-saved')).toBeNull();
	});

	it('shows a refusal with its findings', async () => {
		const findings = [
			{ level: 'warning', rule: 'components.examples', path: PATH, line: 6, message: 'no prop x' }
		];
		show({}, { error: 'flai check refused the save', findings });
		await click('example-save-as');
		type('Bad');
		await submit();
		const alert = q('example-save-error')!;
		expect(alert.textContent).toContain('flai check refused the save');
		expect(alert.querySelector('li')!.textContent).toBe(
			`${PATH}:6: warning: components.examples: no prop x`
		);
		expect(q('example-save-form')).not.toBeNull();
	});

	it('offers no save without a path to save to', () => {
		show({ path: undefined });
		expect(q('example-save-as')).toBeNull();
	});
});
