// S-0757: the Agent images page lists the project's agent images, each with its ref, where its
// name came from, whether it is on the host, and its state with why it is stale, and says why flai
// could not judge them; it lists the running job, those that wait, and the last ended, each an
// ImageJob card whose output it reads through /api/images?job=&offset=; it asks /api/images again
// every 3 s while a job runs or waits and stops once none does; a failed job's Retry, or Make the
// images with nothing failed and an image missing, posts /api/images/retry and the page shows the
// job it started; a refusal is shown word for word, the agent host action's off with the command
// that turns it on and Retry disabled; a 501 says to upgrade flai; and a link to a job's card
// scrolls to it.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { AgentImage, ImageJob, ImagesState } from '$lib/images';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

import ImagesPage from './+page.svelte';

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

const ENABLE = 'flai serve enable agent';
const OFF = [
	'the host action "agent" is not enabled for this project.',
	`On the host, in the project, run: ${ENABLE}`
].join(' ');
const RUNNING = '20261113T093000Z-001';
const FAILED = '20261113T080000Z-001';
const NEXT = '20261113T100000Z-002';

const image = (over: Partial<AgentImage> = {}): AgentImage => ({
	kind: 'harness',
	inner: 'claude-code',
	ref: 'flai-agent:1.94.0',
	source: 'default',
	present: true,
	labels: { flai: '1.94.0', recipe: 'r1', base: 'b1' },
	state: 'current',
	want: { flai: '1.94.0', recipe: 'r1', base: 'b1' },
	...over
});
const layer = (over: Partial<AgentImage> = {}) =>
	image({
		kind: 'layer',
		ref: 'acme-agent:1.94.0',
		source: 'layer',
		dockerfile: 'agent/Dockerfile',
		present: false,
		labels: {},
		state: 'missing',
		...over
	});
const job = (id: string, over: Partial<ImageJob> = {}): ImageJob => ({
	id,
	images: ['flai-agent:1.94.0', 'acme-agent:1.94.0'],
	state: 'running',
	ways: [],
	requested: '2026-11-13T09:30:00Z',
	started: '2026-11-13T09:30:02Z',
	step: 'agent image: building acme-agent:1.94.0',
	log_size: 12,
	...over
});
const waiting = (id: string) =>
	job(id, { state: 'queued', started: undefined, step: undefined, log_size: 0 });
const failed = (over: Partial<ImageJob> = {}) =>
	job(FAILED, {
		state: 'ended',
		outcome: 'failed',
		step: 'agent image: building acme-agent:1.94.0',
		cause: 'docker build exited with status 1',
		ended: '2026-11-13T08:02:00Z',
		duration_ms: 120000,
		ways: [{ image: 'flai-agent:1.94.0', way: 'pulled' }],
		...over
	});
const state = (over: Partial<ImagesState> = {}): ImagesState => ({
	flai: '1.94.0',
	images: [image(), layer()],
	running: null,
	queued: [],
	ended: [],
	...over
});

const all = (id: string) => [...document.querySelectorAll(`[data-testid="${id}"]`)];
const q = (id: string) => document.querySelector(`[data-testid="${id}"]`);
const card = (id: string) => document.querySelector<HTMLElement>(`[data-job="${id}"]`);
const retryOf = (id: string) =>
	card(id)!.querySelector<HTMLButtonElement>('[data-testid="image-job-retry"]');
const retryError = (id: string) =>
	card(id)!.querySelector('[data-testid="image-job-retry-error"]')?.textContent?.trim();
const posts = () =>
	api.mock.calls.filter(([, init]) => (init as RequestInit | undefined)?.method === 'POST');
const asked = () => api.mock.calls.filter(([p]) => p === '/api/images').length;

describe('Agent images page (S-0757)', () => {
	let c: ReturnType<typeof mount> | undefined;
	// the answers GET /api/images gives, in turn; the last is kept
	let reads: ReturnType<typeof answer>[] = [];
	let retry: ReturnType<typeof answer> = answer({});

	const show = async () => {
		c = mount(ImagesPage, { target: document.body });
		await settle();
	};

	beforeEach(() => {
		history.replaceState(null, '', '/ops/images');
		reads = [answer(state())];
		retry = answer({});
		api.mockReset();
		api.mockImplementation(async (path: string, init?: RequestInit) => {
			if (path === '/api/images') return reads.length > 1 ? reads.shift()! : reads[0];
			if (path.startsWith('/api/images?job='))
				return answer({
					id: RUNNING,
					state: 'running',
					exists: true,
					offset: 0,
					next: 12,
					size: 12,
					text: 'pulling one\n'
				});
			if (path === '/api/images/retry' && init?.method === 'POST') return retry;
			throw new Error(`unexpected ${path}`);
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		vi.useRealTimers();
	});

	it('lists the images with their ref, source, presence, state, and reasons, and the jobs by state', async () => {
		reads = [
			answer(
				state({
					images: [
						image({
							state: 'stale',
							reasons: [{ kind: 'flai', message: 'it holds flai 1.93.0, not 1.94.0' }]
						}),
						layer()
					],
					running: job(RUNNING),
					queued: [waiting(NEXT)],
					ended: [failed()]
				})
			)
		];
		await show();
		expect(document.querySelector('h1')!.textContent).toBe('Agent images');
		expect(q('images-flai')!.textContent).toContain('1.94.0');
		const rows = all('image-row');
		expect(rows.map((r) => r.querySelector('[data-testid="image-ref"]')!.textContent)).toEqual([
			'flai-agent:1.94.0',
			'acme-agent:1.94.0'
		]);
		expect(rows[0].textContent).toContain('harness for claude-code');
		expect(rows[0].textContent).toContain("flai's default");
		expect(rows[1].textContent).toContain('project layer for claude-code');
		expect(rows[1].textContent).toContain("the project's layer");
		expect(all('image-present').map((p) => p.textContent)).toEqual(['yes', 'no']);
		expect(all('image-state').map((s) => s.textContent)).toEqual(['stale', 'missing']);
		expect(q('image-reasons')!.textContent).toContain('it holds flai 1.93.0, not 1.94.0');
		expect(q('images-running')!.querySelector('[data-job]')!.getAttribute('data-job')).toBe(
			RUNNING
		);
		expect(q('images-queued')!.querySelector('[data-job]')!.getAttribute('data-job')).toBe(NEXT);
		expect(q('images-ended')!.querySelector('[data-job]')!.getAttribute('data-job')).toBe(FAILED);
		// a failed job shows its cause and step, with Retry; the page's own button gives way to it
		expect(card(FAILED)!.querySelector('[data-testid="image-job-cause"]')!.textContent).toBe(
			'docker build exited with status 1'
		);
		expect(card(FAILED)!.querySelector('[data-testid="image-job-failed-step"]')!.textContent).toBe(
			'agent image: building acme-agent:1.94.0'
		);
		expect(retryOf(FAILED)).not.toBeNull();
		expect(retryOf(FAILED)!.disabled).toBe(false);
		expect(q('images-make')).toBeNull();
		expect(q('images-no-jobs')).toBeNull();
	});

	it('says why flai could not judge the images, and that no job has run', async () => {
		reads = [
			answer(state({ flai: '', images: [], unjudged: "docker is not on flai serve's PATH" }))
		];
		await show();
		expect(q('images-unjudged')!.textContent!.trim()).toBe("docker is not on flai serve's PATH");
		expect(q('images-flai')).toBeNull();
		expect(q('images-none')).not.toBeNull();
		expect(q('images-no-jobs')).not.toBeNull();
		expect(q('images-make')).toBeNull();
	});

	it("follows a running job's output, and asks again every 3 s until no job runs or waits", async () => {
		vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval'] });
		const passed = job(RUNNING, {
			state: 'ended',
			outcome: 'passed',
			ended: '2026-11-13T09:34:00Z',
			duration_ms: 240000,
			ways: [
				{ image: 'flai-agent:1.94.0', way: 'pulled' },
				{ image: 'acme-agent:1.94.0', way: 'layered' }
			]
		});
		reads = [
			answer(state({ running: job(RUNNING) })),
			answer(state({ running: job(RUNNING) })),
			answer(
				state({ images: [image(), layer({ present: true, state: 'current' })], ended: [passed] })
			)
		];
		await show();
		const logs = api.mock.calls.filter(([p]) => (p as string).startsWith('/api/images?job='));
		expect(logs.map(([p]) => p)).toContain(`/api/images?job=${RUNNING}&offset=0`);
		expect(card(RUNNING)!.querySelector('[data-testid="log-feed-text"]')!.textContent).toBe(
			'pulling one\n'
		);
		const before = asked();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 1);
		expect(q('images-running')).not.toBeNull();
		vi.advanceTimersByTime(3000);
		await settle();
		expect(asked()).toBe(before + 2);
		expect(q('images-running')).toBeNull();
		expect(card(RUNNING)!.getAttribute('data-outcome')).toBe('passed');
		// with nothing running or waiting, the page does not ask again
		vi.advanceTimersByTime(30000);
		await settle();
		expect(asked()).toBe(before + 2);
	});

	it('retries a failed job through /api/images/retry and shows the job it started', async () => {
		reads = [answer(state({ ended: [failed()] }))];
		const started = waiting(NEXT);
		retry = answer({
			...state({ queued: [started], ended: [failed()] }),
			started: { ...started, ahead: 0 }
		});
		await show();
		retryOf(FAILED)!.click();
		await settle();
		expect(posts().map(([p]) => p)).toEqual(['/api/images/retry']);
		expect(q('images-queued')!.querySelector('[data-job]')!.getAttribute('data-job')).toBe(NEXT);
		expect(retryError(FAILED)).toBeUndefined();
	});

	it('makes the images when one is missing and no job has failed', async () => {
		await show();
		retry = answer({ ...state({ running: job(NEXT) }), started: { ...job(NEXT), ahead: 0 } });
		(q('images-make') as HTMLButtonElement).click();
		await settle();
		expect(posts().map(([p]) => p)).toEqual(['/api/images/retry']);
		expect(q('images-running')!.querySelector('[data-job]')!.getAttribute('data-job')).toBe(NEXT);
		expect(q('images-make')).toBeNull();
	});

	it('shows the agent action refusal word for word, says how to turn it on, and disables Retry', async () => {
		reads = [answer(state({ ended: [failed()] }))];
		retry = answer({ error: OFF, action: 'agent', enable: ENABLE, project: 'acme' }, 403);
		await show();
		retryOf(FAILED)!.click();
		await settle();
		expect(retryError(FAILED)).toBe(OFF);
		const off = q('images-off')!;
		expect(off.textContent).toContain('The agent host action is off');
		expect(off.querySelector('code')!.textContent).toBe(ENABLE);
		expect(retryOf(FAILED)!.disabled).toBe(true);
	});

	it('shows a refusal while a job makes the images word for word, and reads the images again to show it', async () => {
		reads = [
			answer(state({ ended: [failed()] })),
			answer(state({ running: job(RUNNING), ended: [failed()] }))
		];
		const busy = `the image job ${RUNNING} already makes flai-agent:1.94.0, acme-agent:1.94.0`;
		retry = answer({ error: busy, job: job(RUNNING), project: 'acme' }, 400);
		await show();
		const before = asked();
		retryOf(FAILED)!.click();
		await settle();
		expect(retryError(FAILED)).toBe(busy);
		expect(asked()).toBe(before + 1);
		expect(q('images-running')!.querySelector('[data-job]')!.getAttribute('data-job')).toBe(
			RUNNING
		);
		expect(q('images-off')).toBeNull();
		expect(retryOf(FAILED)!.disabled).toBe(false);
	});

	it('says to upgrade flai on the host when it is too old, in place of the page', async () => {
		const old =
			'the host flai is too old to show the agent images (unknown method images.get): upgrade flai on the host with flai self-upgrade';
		reads = [answer({ error: old, project: 'acme' }, 501)];
		await show();
		expect(q('images-too-old')!.textContent).toContain('upgrade flai on the host');
		expect(q('images-too-old')!.textContent).toContain(old);
		expect(q('images')).toBeNull();
	});

	it("scrolls to the job a story's agent links to once the jobs are drawn", async () => {
		const scrolled: string[] = [];
		const original = Element.prototype.scrollIntoView;
		Element.prototype.scrollIntoView = function (this: Element) {
			scrolled.push(this.id);
		};
		try {
			history.replaceState(null, '', `/ops/images#image-job-${FAILED}`);
			reads = [answer(state({ ended: [failed()] }))];
			await show();
			expect(scrolled).toEqual([`image-job-${FAILED}`]);
		} finally {
			Element.prototype.scrollIntoView = original;
		}
	});
});
