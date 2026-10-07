import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import DismissibleNotice from './DismissibleNotice.svelte';

describe('DismissibleNotice (S-0151)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	it('shows its text, and a right-aligned X labelled Dismiss that calls ondismiss', () => {
		const ondismiss = vi.fn();
		c = mount(DismissibleNotice, {
			target: document.body,
			props: { text: 'S-0151 → review', ondismiss, class: 'border-good', testid: 'said' }
		});
		flushSync();
		const banner = document.querySelector<HTMLElement>('[role="status"]')!;
		expect(banner.className).toContain('border-good');
		expect(document.querySelector('[data-testid="said"]')!.textContent).toBe('S-0151 → review');
		const x = document.querySelector<HTMLButtonElement>('[data-testid="dismiss"]')!;
		expect(banner.lastElementChild).toBe(x);
		expect(x.getAttribute('aria-label')).toBe('Dismiss');
		expect(ondismiss).not.toHaveBeenCalled();
		x.click();
		expect(ondismiss).toHaveBeenCalledOnce();
	});

	it('takes the alert role for a refusal', () => {
		c = mount(DismissibleNotice, {
			target: document.body,
			props: { text: 'refused: no', ondismiss: () => {}, role: 'alert' }
		});
		flushSync();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('refused: no');
	});

	it('shows a time flai names in its text in the local zone (S-0329)', () => {
		c = mount(DismissibleNotice, {
			target: document.body,
			props: {
				text: 'the planner is already running for E-0016 (pid 42, started 2026-10-03T10:00:00Z)',
				ondismiss: () => {},
				testid: 'said'
			}
		});
		flushSync();
		expect(document.querySelector('[data-testid="said"]')!.textContent).toBe(
			'the planner is already running for E-0016 (pid 42, started 2026-10-03 06:00 EDT)'
		);
	});
});
