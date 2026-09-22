import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostPanel from './HostPanel.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
// For a flow that itself waits on real (if short) timers: real ticks, not zero-delay ones.
const settleThrough = async (ms = 100) => {
	await new Promise((r) => setTimeout(r, ms));
	flushSync();
};
const answer = (body: unknown, ok = true) => ({
	ok,
	statusText: 'Error',
	json: async () => body
});
const status = (over: Record<string, unknown> = {}) =>
	answer({
		container: 'flaiover',
		running: true,
		image: 'ghcr.io/bytepunx/flaiover:0.22.6',
		serves: ['harbour'],
		dashboard_enabled: true,
		...over
	});

describe('HostPanel', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('shows what is running and hides the buttons when the action is off', async () => {
		api.mockResolvedValueOnce(status({ dashboard_enabled: false }));
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		const text = document.querySelector('[data-testid="host-panel-status"]')!.textContent!;
		expect(text).toContain('flaiover');
		expect(text).toContain('ghcr.io/bytepunx/flaiover:0.22.6');
		expect(text).toContain('harbour');
		expect(document.querySelectorAll('button')).toHaveLength(0);
		expect(document.body.textContent).toContain('flai serve enable dashboard');
	});

	it('offers the four actions once enabled', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		for (const testid of [
			'host-panel-check',
			'host-panel-restart',
			'host-panel-upgrade',
			'host-panel-stop'
		]) {
			expect(document.querySelector(`[data-testid="${testid}"]`)).not.toBeNull();
		}
	});

	it('checks for an update without touching the running container', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		api.mockResolvedValueOnce(answer({ running: true, upgrade_available: true }));
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-check"]')!.click();
		await settle();
		expect(api).toHaveBeenLastCalledWith(
			'/api/dashboard',
			expect.objectContaining({ method: 'POST', body: JSON.stringify({ action: 'check' }) })
		);
		expect(document.querySelector('[data-testid="host-panel-check-result"]')!.textContent).toBe(
			'An update is available.'
		);
	});

	it('shows a plain refusal when a check fails outright', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		api.mockResolvedValueOnce(answer({ error: 'docker is required' }, false));
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-check"]')!.click();
		await settle();
		expect(document.querySelector('[data-testid="host-panel-failed"]')!.textContent).toBe(
			'docker is required'
		);
	});

	// Real timing, cut to a few milliseconds via props, rather than vi.useFakeTimers(): the
	// component races a fetch against setTimeout inside its own async flow, which fake timers do
	// not advance until the microtasks they unblock are themselves awaited, and getting that
	// interleaving right is far more fragile than just making the real waits short.
	const fast = { disconnectTimeoutMs: 5, reconnectPollMs: 5, reconnectGiveUpMs: 50 };

	it('treats a restart that drops the connection as proceeding, then reports what reconnected', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body, props: fast });
		flushSync();
		await settle();

		// the POST never resolves: the container answering it was stopped mid-request
		api.mockImplementationOnce(() => new Promise(() => {}));
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-restart"]')!.click();
		await settleThrough(20);
		expect(document.querySelector('[data-testid="host-panel-reconnecting"]')).not.toBeNull();

		// the new container answers once it is back
		api.mockResolvedValue(status({ image: 'ghcr.io/bytepunx/flaiover:0.22.6' }));
		await settleThrough(20);
		expect(document.querySelector('[data-testid="host-panel-message"]')!.textContent).toBe(
			'Reconnected.'
		);
		expect(document.querySelector('[data-testid="host-panel-reconnecting"]')).toBeNull();
	});

	it('says what changed when an upgrade reconnects to a different version', async () => {
		api.mockResolvedValueOnce(status({ image: 'ghcr.io/bytepunx/flaiover:0.22.6' }));
		c = mount(HostPanel, { target: document.body, props: fast });
		flushSync();
		await settle();

		api.mockImplementationOnce(() => new Promise(() => {}));
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-upgrade"]')!.click();
		await settleThrough(20);

		api.mockResolvedValue(status({ image: 'ghcr.io/bytepunx/flaiover:0.22.7' }));
		await settleThrough(20);
		expect(document.querySelector('[data-testid="host-panel-message"]')!.textContent).toBe(
			'Reconnected — now running ghcr.io/bytepunx/flaiover:0.22.7.'
		);
	});

	it('reports an upgrade that never became healthy as a plain failure, the previous container untouched', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		api.mockResolvedValueOnce(
			answer(
				{
					error:
						'upgrade to ghcr.io/bytepunx/flaiover:latest failed: the new image did not answer healthy in time; flaiover keeps running unchanged'
				},
				false
			)
		);
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-upgrade"]')!.click();
		await settle();
		expect(document.querySelector('[data-testid="host-panel-failed"]')!.textContent).toContain(
			'keeps running unchanged'
		);
		expect(document.querySelector('[data-testid="host-panel-reconnecting"]')).toBeNull();
	});

	it('stops and shows whether the container itself stopped or just unregistered', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		api.mockResolvedValueOnce(
			answer({ container: 'flaiover', state: 'still-running', serves: ['quay'] })
		);
		api.mockResolvedValueOnce(status({ serves: ['quay'] }));
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-stop"]')!.click();
		await settle();
		expect(document.querySelector('[data-testid="host-panel-message"]')!.textContent).toBe(
			'Unregistered; still serving quay.'
		);
	});

	// Found live (S-0081, T-0325): stopping the last project stops the very container answering
	// the request, the same self-termination a restart causes; the first version of this
	// component only raced that fetch for restart and upgrade, so a real "last project" stop
	// threw an unhandled rejection instead of reporting anything.
	it('treats a stop that drops the connection as the container genuinely gone, not a failure', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body, props: fast });
		flushSync();
		await settle();
		api.mockImplementationOnce(() => new Promise(() => {}));
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-stop"]')!.click();
		await settleThrough(20);
		expect(document.querySelector('[data-testid="host-panel-message"]')!.textContent).toContain(
			'last project'
		);
		expect(document.querySelector('[data-testid="host-panel-reconnecting"]')).toBeNull();
		expect(document.querySelector('[data-testid="host-panel-failed"]')).toBeNull();
	});

	// Found live: an upgrade that changed nothing (already current) still made the page say
	// "Reconnected", implying a disruption that never happened, because it always waited to
	// reconnect after any clean response.
	it('does not claim a reconnect for an upgrade that touched nothing', async () => {
		api.mockResolvedValueOnce(status());
		c = mount(HostPanel, { target: document.body });
		flushSync();
		await settle();
		api.mockResolvedValueOnce(
			answer({ container: 'flaiover-s0325-scratch', outcome: 'up-to-date', to: 'a' })
		);
		document.querySelector<HTMLButtonElement>('[data-testid="host-panel-upgrade"]')!.click();
		await settle();
		expect(document.querySelector('[data-testid="host-panel-message"]')!.textContent).toContain(
			'already running the latest'
		);
		expect(document.querySelector('[data-testid="host-panel-reconnecting"]')).toBeNull();
	});
});
