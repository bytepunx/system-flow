import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import PublishBanner from './PublishBanner.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, json: async () => body });
const plans = [
	{
		component: { name: 'cli' },
		level: 'minor',
		from: '1.0.0',
		to: '1.1.0',
		items: [
			{ id: 'S-0101', title: 'Fix one', level: 'patch' },
			{ id: 'S-0102', title: 'Add a thing', level: 'minor' }
		]
	}
];
const pending = () => document.querySelector<HTMLElement>('[data-testid="publish-pending"]');

describe('PublishBanner', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('shows nothing when nothing is pending', () => {
		c = mount(PublishBanner, { target: document.body, props: { plans: [], enabled: true } });
		flushSync();
		expect(document.body.textContent!.trim()).toBe('');
	});

	it('shows every pending component, its bump, and the story IDs it bundles, with no button when the action is off', () => {
		c = mount(PublishBanner, { target: document.body, props: { plans, enabled: false } });
		flushSync();
		const text = pending()!.textContent!.replace(/\s+/g, ' ');
		expect(text).toContain('cli');
		expect(text).toContain('1.0.0 → 1.1.0 (minor)');
		expect(text).toContain('S-0101, S-0102');
		expect(document.querySelectorAll('button')).toHaveLength(0);
		expect(text).toContain('Publishing from the board is off');
		expect(text).toContain('flai serve enable push');
	});

	it('offers a Publish button when the action is enabled, and reports success', async () => {
		let release: (v: unknown) => void = () => {};
		api.mockImplementation(() => new Promise((r) => (release = r)));
		const onpublished = vi.fn();
		c = mount(PublishBanner, {
			target: document.body,
			props: { plans, enabled: true, onpublished }
		});
		flushSync();
		const button = document.querySelector<HTMLButtonElement>('[data-testid="publish-now"]')!;
		expect(button.textContent).toContain('Publish');
		button.click();
		await settle();
		expect(api).toHaveBeenCalledWith('/api/publish', { method: 'POST' });
		expect(button.disabled).toBe(true);
		expect(button.textContent).toContain('Publishing…');
		release(answer({ pushed: true, tags: ['cli/v1.1.0'], published: [] }));
		await settle();
		expect(onpublished).toHaveBeenCalledOnce();
		const done = document.querySelector('[data-testid="published"]')!.textContent!;
		expect(done).toContain('Published cli/v1.1.0, pushed.');
	});

	it('shows why a publish was refused and keeps the pending plan visible', async () => {
		api.mockResolvedValue({
			ok: false,
			statusText: 'Conflict',
			json: async () => ({ error: 'main and origin/main have diverged: fetch and merge first' })
		});
		c = mount(PublishBanner, { target: document.body, props: { plans, enabled: true } });
		flushSync();
		document.querySelector<HTMLButtonElement>('[data-testid="publish-now"]')!.click();
		await settle();
		const refused = document.querySelector('[data-testid="publish-refused"]')!.textContent!;
		expect(refused).toContain('Not published:');
		expect(refused).toContain('have diverged');
		expect(pending()).not.toBeNull();
	});
});
