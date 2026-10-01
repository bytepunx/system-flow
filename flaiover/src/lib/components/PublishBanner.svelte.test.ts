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
		// S-0151: the result is dismissed by its X, and the plan stays
		document.querySelector<HTMLButtonElement>('[data-testid="dismiss"]')!.click();
		flushSync();
		expect(document.querySelector('[data-testid="published"]')).toBeNull();
		expect(pending()).not.toBeNull();
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
	// S-0174: a clone missing release tags its remote has is told so, with no plan and no button.
	it('says which release tags the clone is missing and how to fetch them, offering nothing', () => {
		const remote = {
			remote: 'origin',
			behind: [
				{ component: 'flai', local: 'flai/v1.18.0', remote: 'flai/v1.26.3' },
				{ component: 'flaiover', remote: 'flaiover/v0.32.1' }
			],
			fix: 'git fetch --tags origin',
			message: 'missing'
		};
		c = mount(PublishBanner, {
			target: document.body,
			props: { plans: [], enabled: true, remote }
		});
		flushSync();
		const text = document
			.querySelector('[data-testid="publish-missing-tags"]')!
			.textContent!.replace(/\s+/g, ' ');
		expect(text).toContain('missing release tags origin has');
		expect(text).toContain('flai/v1.26.3 (here flai/v1.18.0)');
		expect(text).toContain('flaiover/v0.32.1 (here none)');
		expect(text).toContain('git fetch --tags origin');
		expect(pending()).toBeNull();
		expect(document.querySelectorAll('button')).toHaveLength(0);
	});

	// S-0174: a remote that could not be asked leaves the plan in view, warned, with no button.
	it('warns over the plan when the remote could not be asked, and offers no Publish', () => {
		const remote = {
			remote: 'origin',
			unchecked: 'unable to access the remote',
			fix: 'git fetch --tags origin',
			message: 'could not ask'
		};
		c = mount(PublishBanner, { target: document.body, props: { plans, enabled: true, remote } });
		flushSync();
		expect(pending()!.textContent).toContain('S-0101, S-0102');
		const warn = document.querySelector('[data-testid="publish-unchecked"]')!.textContent!;
		expect(warn).toContain('Not checked against origin');
		expect(warn).toContain('unable to access the remote');
		expect(document.querySelector('[data-testid="publish-now"]')).toBeNull();
		expect(document.querySelector('[data-testid="publish-missing-tags"]')).toBeNull();
	});

	// I-0024: what no plan covers is named with why, beside the plan.
	it('names the accepted items no plan covers, with why', () => {
		const unplanned = [
			{
				id: 'S-0103',
				title: 'Both sides',
				reason: 'S-0103 touches cli, web but no tag says which it delivers to'
			}
		];
		c = mount(PublishBanner, { target: document.body, props: { plans, enabled: true, unplanned } });
		flushSync();
		const text = document.querySelector('[data-testid="publish-unplanned"]')!.textContent!;
		expect(text).toContain('S-0103');
		expect(text).toContain('Both sides');
		expect(text).toContain('no tag says which it delivers to');
		expect(pending()).not.toBeNull();
	});
});
