import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import CostOfDelayPanel from './CostOfDelayPanel.svelte';
import type { CostFields } from '$lib/costofdelay';

const q = <T extends HTMLElement>(id: string) => document.querySelector<T>(`[data-testid="${id}"]`);
const toggle = (open: boolean) => {
	const box = q<HTMLDetailsElement>('cod-panel')!;
	box.open = open;
	box.dispatchEvent(new Event('toggle'));
	flushSync();
};

describe('CostOfDelayPanel (S-0204)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	it('starts closed with no summary when nothing is set, and labels the amounts with the currency', () => {
		const fields = $state<CostFields>({
			revenue_per_week: '',
			penalty_per_week: '',
			time_lost_per_cycle: ''
		});
		c = mount(CostOfDelayPanel, { target: document.body, props: { fields, currency: 'EUR' } });
		flushSync();
		expect(q<HTMLDetailsElement>('cod-panel')!.open).toBe(false);
		expect(q<HTMLDetailsElement>('cod-panel')!.textContent).toContain('Cost of delay');
		expect(q('cod-summary')).toBeNull();
		const text = document.body.textContent!;
		expect(text).toContain('revenue per week, EUR');
		expect(text).toContain('penalty per week, EUR');
		expect(text).toContain('time lost per cycle, a duration');
	});

	it('says what is set beside its title while closed, and binds what is typed', () => {
		const fields = $state<CostFields>({
			revenue_per_week: '1200',
			penalty_per_week: '',
			time_lost_per_cycle: '6h'
		});
		c = mount(CostOfDelayPanel, { target: document.body, props: { fields, currency: 'USD' } });
		flushSync();
		expect(q('cod-summary')!.textContent).toBe('revenue 1,200 USD/week · 6h lost per cycle');
		toggle(true);
		expect(q('cod-summary')).toBeNull();
		const penalty = q<HTMLInputElement>('cod-penalty')!;
		penalty.value = '300';
		penalty.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		expect(fields.penalty_per_week).toBe('300');
		toggle(false);
		expect(q('cod-summary')!.textContent).toBe(
			'revenue 1,200 USD/week · penalty 300 USD/week · 6h lost per cycle'
		);
	});
});
