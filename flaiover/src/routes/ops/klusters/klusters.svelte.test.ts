// S-0580: the Klusters page shows the project's cluster, its name, profile, provider, whether it
// runs, and its age; brings it up, and takes it down once confirmed, through /api/kluster; says
// while a job runs and asks again until it ends, showing the end of its log; shows flai's refusal
// word for word; and in place of the buttons says when the kluster host action is off, linking the
// Features page, or when the manifest has no kluster: block, naming it and linking the design.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { KlusterJob, KlusterStatus } from '$lib/server/repo';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params?: Record<string, string>) =>
		params ? route.replace(/\[(?:\.\.\.)?(\w+)\]/g, (_, k: string) => params[k]) : route
}));

import KlustersPage from './+page.svelte';

const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

const ENABLE = 'flai serve enable kluster';
const NO_BLOCK =
	"this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it; add one";

const job = (action: 'up' | 'down', over: Partial<KlusterJob> = {}): KlusterJob => ({
	action,
	cluster: 'authstar',
	started_at: '2026-11-07T01:00:00Z',
	log: `.flai-cache/kluster/${action}.log`,
	...over
});
const status = (over: Partial<KlusterStatus> = {}): KlusterStatus => ({
	enabled: true,
	enable: ENABLE,
	block: true,
	cluster: {
		name: 'authstar',
		profile: 'integration',
		provider: 'k3d',
		context: 'k3d-authstar',
		exists: true,
		running: true,
		age: '2h50m39s',
		age_seconds: 10239
	},
	...over
});

const q = (id: string) => document.querySelector(`[data-testid="${id}"]`);
const buttonNamed = (name: string) =>
	[...document.querySelectorAll<HTMLButtonElement>('button')].find(
		(b) => b.textContent!.trim() === name
	);
const posts = () =>
	api.mock.calls.filter(([, init]) => (init as RequestInit | undefined)?.method === 'POST');

describe('Klusters page (S-0580)', () => {
	let c: ReturnType<typeof mount> | undefined;
	// the answers GET /api/kluster gives, in turn; the last is kept
	let reads: ReturnType<typeof answer>[] = [];
	let writes: Record<string, ReturnType<typeof answer>> = {};

	const show = async () => {
		c = mount(KlustersPage, { target: document.body });
		await settle();
	};

	beforeEach(() => {
		reads = [answer(status())];
		writes = {};
		api.mockReset();
		api.mockImplementation(async (path: string) => {
			if (path === '/api/kluster') return reads.length > 1 ? reads.shift()! : reads[0];
			const w = writes[path];
			if (!w) throw new Error(`unexpected ${path}`);
			return w;
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		vi.useRealTimers();
	});

	it("shows the cluster's name, profile, provider, state, and age", async () => {
		await show();
		const text = q('kluster-cluster')!.textContent!;
		for (const s of ['authstar', 'integration', 'k3d', 'Running', '2h50m39s']) {
			expect(text).toContain(s);
		}
		expect(buttonNamed('Up')).toBeDefined();
		expect(buttonNamed('Down')).toBeDefined();
	});

	it('says a cluster is stopped or not created', async () => {
		reads = [answer(status({ cluster: { ...status().cluster!, running: false, age: undefined } }))];
		await show();
		expect(q('kluster-state')!.textContent).toBe('Stopped');
		unmount(c!);
		reads = [answer(status({ cluster: { ...status().cluster!, running: false, exists: false } }))];
		c = mount(KlustersPage, { target: document.body });
		await settle();
		expect(q('kluster-state')!.textContent).toBe('Not created');
	});

	it('brings the cluster up and shows the job running, the end of its log, and asks again until it ends', async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		await show();
		const running = job('up', { log_tail: 'creating cluster authstar\n' });
		writes['/api/kluster/up'] = answer(status({ running }));
		buttonNamed('Up')!.click();
		await settle();
		expect(posts().map(([p]) => p)).toEqual(['/api/kluster/up']);
		expect(q('kluster-running')!.textContent).toContain('kluster up of the cluster authstar runs');
		expect(q('kluster-log')!.textContent).toBe('creating cluster authstar\n');
		expect(buttonNamed('Up')!.disabled).toBe(true);
		expect(buttonNamed('Down')!.disabled).toBe(true);

		const done = job('up', {
			ended_at: '2026-11-07T01:03:00Z',
			outcome: 'done',
			log_tail: 'cluster authstar is up\n'
		});
		reads = [answer(status({ running })), answer(status({ last: done }))];
		const asked = () => api.mock.calls.filter(([p]) => p === '/api/kluster').length;
		const before = asked();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 1);
		expect(q('kluster-running')).not.toBeNull();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 2);
		expect(q('kluster-running')).toBeNull();
		expect(q('kluster-last')!.textContent).toContain('kluster up of the cluster authstar');
		expect(q('kluster-last')!.textContent).toContain('done');
		expect(q('kluster-log')!.textContent).toBe('cluster authstar is up\n');
		// with nothing running, the page does not ask again
		vi.advanceTimersByTime(30000);
		await settle();
		expect(asked()).toBe(before + 2);
	});

	it('takes the cluster down only once the operator confirms', async () => {
		await show();
		writes['/api/kluster/down'] = answer(status({ running: job('down') }));
		buttonNamed('Down')!.click();
		flushSync();
		expect(q('kluster-confirm-down')!.textContent).toContain('Take the cluster authstar down?');
		expect(posts()).toEqual([]);
		buttonNamed('Keep it')!.click();
		flushSync();
		expect(q('kluster-confirm-down')).toBeNull();
		expect(posts()).toEqual([]);

		buttonNamed('Down')!.click();
		flushSync();
		buttonNamed('Take it down')!.click();
		await settle();
		expect(posts().map(([p]) => p)).toEqual(['/api/kluster/down']);
		expect(q('kluster-confirm-down')).toBeNull();
		expect(q('kluster-running')!.textContent).toContain(
			'kluster down of the cluster authstar runs'
		);
	});

	it("shows a failed job's error as flai gives it", async () => {
		reads = [
			answer(
				status({
					last: job('up', {
						ended_at: '2026-11-07T01:01:00Z',
						outcome: 'failed',
						error: "kluster is not on the host's PATH: install it, as design/tech/kluster.md says"
					})
				})
			)
		];
		await show();
		expect(q('kluster-last')!.textContent).toContain('failed');
		expect(q('kluster-last-error')!.textContent).toBe(
			"kluster is not on the host's PATH: install it, as design/tech/kluster.md says"
		);
	});

	it("shows flai's refusal word for word", async () => {
		await show();
		const refusal =
			'kluster up of the cluster authstar runs since 2026-11-07T01:00:00Z; wait for it to end, as kluster.status says, then ask again';
		writes['/api/kluster/up'] = answer({ error: refusal }, 400);
		buttonNamed('Up')!.click();
		await settle();
		expect(q('kluster-action-error')!.textContent!.trim()).toBe(refusal);
		expect(buttonNamed('Up')!.disabled).toBe(false);
	});

	it('shows the status error kluster gave, word for word', async () => {
		const why = "kluster is not on the host's PATH";
		reads = [answer(status({ status_error: why }))];
		await show();
		expect(q('kluster-status-error')!.textContent!.trim()).toBe(why);
	});

	it('says the host action is off, with the command and the Features page, in place of the buttons', async () => {
		reads = [answer(status({ enabled: false }))];
		await show();
		const off = q('kluster-off')!;
		expect(off.textContent).toContain('The kluster host action is off');
		expect(off.textContent).toContain(ENABLE);
		expect(off.querySelector('a')!.getAttribute('href')).toBe('/settings/features');
		expect(buttonNamed('Up')).toBeUndefined();
		expect(buttonNamed('Down')).toBeUndefined();
		// the cluster is still shown
		expect(q('kluster-cluster')!.textContent).toContain('authstar');
	});

	it('turns to the off notice when a write is refused because the action went off', async () => {
		await show();
		const message = `the host action "kluster" is not enabled for this project. On the host, in the project, run: ${ENABLE}`;
		writes['/api/kluster/up'] = answer({ error: message, action: 'kluster', enable: ENABLE }, 403);
		buttonNamed('Up')!.click();
		await settle();
		expect(q('kluster-off')).not.toBeNull();
		expect(q('kluster-action-error')!.textContent!.trim()).toBe(message);
		expect(buttonNamed('Up')).toBeUndefined();
	});

	it('names the block to add and links the design, in place of the buttons, with no block', async () => {
		reads = [answer({ enabled: true, enable: ENABLE, block: false, no_block: NO_BLOCK })];
		await show();
		const none = q('kluster-no-block')!;
		expect(none.textContent).toContain('no kluster: block');
		expect(none.textContent).toContain(NO_BLOCK);
		expect(none.querySelector('a')!.getAttribute('href')).toBe('/docs/design/system/kluster.md');
		expect(buttonNamed('Up')).toBeUndefined();
		expect(q('kluster-cluster')).toBeNull();
	});

	it('tells to upgrade flai on the host when it is too old', async () => {
		reads = [answer({ error: 'the host flai is too old to show the kluster' }, 501)];
		await show();
		expect(q('kluster-too-old')!.textContent).toContain('upgrade flai on the host');
		expect(buttonNamed('Up')).toBeUndefined();
	});
});
