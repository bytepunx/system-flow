// S-0176: a story's page shows its task plan: each task's state, what a waiting task waits for, and
// the layers of tasks that can run at once.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount, type ComponentProps } from 'svelte';
import TaskPlan from './TaskPlan.svelte';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id)
}));

// in the contract's shape: item.get's plan beside the item and its children
const plan: ComponentProps<typeof TaskPlan>['plan'] = {
	tasks: [
		{ id: 'T-0001', state: 'done' },
		{ id: 'T-0002', state: 'in-progress', after: ['T-0001'] },
		{ id: 'T-0003', state: 'ready', after: ['T-0001'] },
		{
			id: 'T-0004',
			state: 'waiting',
			after: ['T-0001', 'T-0002', 'T-0003'],
			waiting_for: ['T-0002', 'T-0003']
		},
		{ id: 'T-0005', state: 'cancelled' }
	],
	layers: [['T-0001'], ['T-0002', 'T-0003'], ['T-0004']]
};

// the space after a comma between linked tasks is the row's gap, not text: read it as a space
const text = (el: Element) =>
	el
		.textContent!.replace(/\s+/g, ' ')
		.replace(/,(?=\S)/g, ', ')
		.trim();
const tasks = () => [...document.querySelectorAll<HTMLElement>('[data-testid="plan-task"]')];
const layers = () => [...document.querySelectorAll<HTMLElement>('[data-testid="plan-layer"]')];

describe('TaskPlan', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});
	function render(props: ComponentProps<typeof TaskPlan>) {
		c = mount(TaskPlan, { target: document.body, props });
		flushSync();
	}

	it('says each task’s state, and what a waiting task waits for, linked', () => {
		render({ plan, titles: { 'T-0002': 'The flai side' } });
		expect(tasks().map(text)).toEqual([
			'T-0001 done',
			'T-0002 in progress',
			'T-0003 ready to start',
			'T-0004 waiting for T-0002, T-0003',
			'T-0005 cancelled'
		]);
		expect(tasks().map((t) => t.dataset.state)).toEqual([
			'done',
			'in-progress',
			'ready',
			'waiting',
			'cancelled'
		]);
		const waiting = tasks()[3];
		expect([...waiting.querySelectorAll('a')].map((a) => a.getAttribute('href'))).toEqual([
			'/items/T-0004',
			'/items/T-0002',
			'/items/T-0003'
		]);
		// a task's title, from the story's children, is its link's tooltip
		expect(tasks()[1].querySelector('a')!.getAttribute('title')).toBe('The flai side');
	});

	it('lists the layers first to last, each with its tasks linked', () => {
		render({ plan });
		expect(document.querySelector('h3')!.textContent).toBe('Layers');
		expect(layers().map(text)).toEqual([
			'layer 1: T-0001',
			'layer 2: T-0002, T-0003',
			'layer 3: T-0004'
		]);
		expect([...layers()[1].querySelectorAll('a')].map((a) => a.getAttribute('href'))).toEqual([
			'/items/T-0002',
			'/items/T-0003'
		]);
		expect(document.querySelector('[data-testid="plan-cycle"]')).toBeNull();
	});

	it('names the open tasks a cycle leaves out of every layer', () => {
		render({
			plan: {
				tasks: [
					{ id: 'T-0001', state: 'ready' },
					{ id: 'T-0002', state: 'waiting', after: ['T-0003'], waiting_for: ['T-0003'] },
					{ id: 'T-0003', state: 'waiting', after: ['T-0002'], waiting_for: ['T-0002'] }
				],
				layers: [['T-0001']]
			}
		});
		expect(text(document.querySelector('[data-testid="plan-cycle"]')!)).toBe(
			'in no layer, as what they wait for runs in a cycle: T-0002, T-0003'
		);
	});

	it('shows no layers when there are none', () => {
		render({ plan: { tasks: [{ id: 'T-0001', state: 'done' }], layers: [] } });
		expect(tasks().map(text)).toEqual(['T-0001 done']);
		expect(document.querySelector('h3')).toBeNull();
		expect(document.querySelector('[data-testid="plan-cycle"]')).toBeNull();
	});
});
