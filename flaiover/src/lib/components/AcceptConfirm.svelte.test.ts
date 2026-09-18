import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import AcceptConfirm from './AcceptConfirm.svelte';

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

describe('AcceptConfirm', () => {
	afterEach(() => {
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('shows the release plan and the branch, and only acts on confirm', async () => {
		api.mockResolvedValue(
			json({
				id: 'S-0046',
				branch: 'story/S-0046',
				plan: {
					level: 'minor',
					commits: ['a', 'b'],
					steps: [
						{
							component: { name: 'flai' },
							delivered: true,
							level: 'minor',
							from: { Major: 1, Minor: 1, Patch: 6 },
							to: { Major: 1, Minor: 2, Patch: 0 },
							files: ['x', 'y']
						},
						{
							component: { name: 'flaiover' },
							delivered: false,
							level: 'patch',
							from: '0.10.2',
							to: '0.10.3',
							files: ['z']
						}
					]
				}
			})
		);
		const onconfirm = vi.fn();
		const oncancel = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0046', onconfirm, oncancel }
		});
		await settle();
		expect(api).toHaveBeenCalledWith('/api/items/S-0046/acceptance');
		const text = document.body.textContent ?? '';
		expect(text).toContain('story/S-0046');
		expect(text).toContain('minor release from 2');
		expect(text).toMatch(/flai\s+1\.1\.6 → 1\.2\.0/);
		expect(text).toMatch(/flaiover\s+0\.10\.2 → 0\.10\.3/);
		const [cancel, accept] = [...document.querySelectorAll('button')];
		cancel.click();
		expect(oncancel).toHaveBeenCalledOnce();
		expect(onconfirm).not.toHaveBeenCalled();
		accept.click();
		await settle();
		expect(onconfirm).toHaveBeenCalledOnce();
		unmount(c);
	});

	it('lists blockers from the preview and keeps accept disabled', async () => {
		api.mockResolvedValue(
			json({ id: 'S-0046', blockers: ['git has no committer identity here'], plan: null })
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0046', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[role=alert]')?.textContent).toContain('committer identity');
		expect(([...document.querySelectorAll('button')][1] as HTMLButtonElement).disabled).toBe(true);
		unmount(c);
	});

	it('shows flai’s refusal verbatim and keeps accept disabled', async () => {
		api.mockResolvedValue(
			json({ error: 'S-0046 is in-progress; a story is accepted from review' }, false)
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0046', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[role=alert]')?.textContent).toContain(
			'a story is accepted from review'
		);
		const accept = [...document.querySelectorAll('button')][1] as HTMLButtonElement;
		expect(accept.disabled).toBe(true);
		accept.click();
		expect(onconfirm).not.toHaveBeenCalled();
		unmount(c);
	});
});
