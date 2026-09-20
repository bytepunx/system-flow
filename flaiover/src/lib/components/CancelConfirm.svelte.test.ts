import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import CancelConfirm from './CancelConfirm.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const json = (body: unknown, ok = true) => ({
	ok,
	statusText: ok ? 'OK' : 'Conflict',
	json: async () => body
});
const button = (text: string) =>
	[...document.querySelectorAll('button')].find((b) => b.textContent?.includes(text))!;

describe('CancelConfirm', () => {
	afterEach(() => {
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('asks for a dry run, lists what goes with the item, and acts only on confirm with a reason', async () => {
		api.mockResolvedValue(
			json({
				id: 'E-0007',
				dry_run: true,
				cancelled: [
					{ id: 'S-0065', type: 'story', title: 'Registry', from: 'ready' },
					{ id: 'T-0301', type: 'task', title: 'Write it', from: 'backlog' },
					{ id: 'S-0068', type: 'story', title: 'Embedded', from: 'review' }
				],
				left_behind: [{ id: 'S-0068', narrative: 'wip/agents/S-0068.md', branch: 'story/S-0068' }]
			})
		);
		const onconfirm = vi.fn();
		const oncancel = vi.fn();
		const c = mount(CancelConfirm, {
			target: document.body,
			props: { id: 'E-0007', onconfirm, oncancel }
		});
		await settle();
		expect(api).toHaveBeenCalledWith(
			'/api/items/E-0007/move',
			expect.objectContaining({
				method: 'POST',
				body: JSON.stringify({ to: 'cancelled', dry_run: true })
			})
		);
		const text = document.body.textContent ?? '';
		expect(text).toContain('also cancels 3');
		expect(document.querySelectorAll('[data-cascade] li').length).toBe(3);
		expect(text).toContain('S-0068');
		expect(text).toMatch(/S-0068\s+is in review/);
		expect(text).toContain('branch story/S-0068');

		const go = button('Cancel E-0007 and 3 more');
		expect(go.disabled).toBe(true);
		const reason = document.querySelector('textarea')!;
		reason.value = '  a different route ';
		reason.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		expect(go.disabled).toBe(false);
		expect(onconfirm).not.toHaveBeenCalled();
		go.click();
		await settle();
		expect(onconfirm).toHaveBeenCalledWith('a different route');

		button('Keep it').click();
		expect(oncancel).toHaveBeenCalled();
		unmount(c);
	});

	it('says so when nothing is open under the item', async () => {
		api.mockResolvedValue(json({ id: 'T-0001', cancelled: [], left_behind: [] }));
		const c = mount(CancelConfirm, {
			target: document.body,
			props: { id: 'T-0001', onconfirm: vi.fn(), oncancel: vi.fn() }
		});
		await settle();
		expect(document.body.textContent).toContain('Nothing open under it');
		expect(button('Cancel T-0001')).toBeTruthy();
		unmount(c);
	});

	it('shows a refusal and offers no confirm', async () => {
		api.mockResolvedValue(
			json({ error: 'rule: S-0004 cannot go from review to cancelled' }, false)
		);
		const c = mount(CancelConfirm, {
			target: document.body,
			props: { id: 'S-0004', onconfirm: vi.fn(), oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[role="alert"]')?.textContent).toContain(
			'cannot go from review'
		);
		expect(button('Cancel S-0004').disabled).toBe(true);
		unmount(c);
	});
});
