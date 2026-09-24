import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostProcesses from './HostProcesses.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const settleThrough = async (ms = 20) => {
	await new Promise((r) => setTimeout(r, ms));
	flushSync();
};
const answer = (body: unknown, ok = true, status = ok ? 200 : 400) => ({
	ok,
	status,
	statusText: 'Error',
	json: async () => body
});
const host = (over: Record<string, unknown> = {}) =>
	answer({
		running: true,
		pid: 4100,
		version: '1.9.0',
		host_enabled: true,
		children: [
			{ name: 'serve', state: 'running', pid: 4101, version: '1.9.0', restarts: 0 },
			{ name: 'mcp', root: '/src/harbour', state: 'running', version: '1.9.0', restarts: 1 },
			{ name: 'mcp', root: '/src/quay', state: 'running', version: '1.9.0', restarts: 0 }
		],
		...over
	});
const q = (testid: string) =>
	document.querySelector<HTMLButtonElement>(`[data-testid="${testid}"]`);
const text = (testid: string) => q(testid)?.textContent?.replace(/\s+/g, ' ').trim();
const lastPost = () => JSON.parse(api.mock.calls.at(-1)![1].body);
const fast = { disconnectTimeoutMs: 5, reconnectPollMs: 5, reconnectGiveUpMs: 50 };

describe('HostProcesses', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	async function open(first = host(), props = {}) {
		api.mockResolvedValueOnce(first);
		c = mount(HostProcesses, { target: document.body, props });
		flushSync();
		await settle();
	}

	it('shows the host, serve, and the MCP servers with their state and version', async () => {
		await open();
		expect(text('host-processes-host')).toBe('1.9.0, pid 4100');
		expect(text('host-processes-serve-state')).toBe('running');
		expect(text('host-processes-serve-version')).toBe('1.9.0');
		expect(text('host-processes-mcp-state')).toBe('running');
		expect(text('host-processes-mcp-projects')).toBe('harbour, quay');
		expect(api).toHaveBeenCalledWith('/api/host');
	});

	it('offers start, stop, and restart per process, disabling what would do nothing', async () => {
		await open(
			host({
				children: [
					{ name: 'serve', state: 'running', version: '1.9.0', restarts: 0 },
					{ name: 'mcp', root: '/src/harbour', state: 'stopped', version: '1.9.0', restarts: 0 }
				]
			})
		);
		expect(q('host-processes-serve-start')!.disabled).toBe(true);
		expect(q('host-processes-serve-stop')!.disabled).toBe(false);
		expect(q('host-processes-serve-restart')!.disabled).toBe(false);
		expect(q('host-processes-mcp-start')!.disabled).toBe(false);
		expect(q('host-processes-mcp-stop')!.disabled).toBe(true);
		expect(text('host-processes-mcp-state')).toBe('stopped');
		expect(q('host-processes-check')).not.toBeNull();
		expect(q('host-processes-upgrade')).not.toBeNull();
	});

	it('disables the MCP controls while no project’s server is kept', async () => {
		await open(host({ children: [{ name: 'serve', state: 'running', restarts: 0 }] }));
		expect(text('host-processes-mcp-state')).toBe('none');
		for (const a of ['start', 'stop', 'restart']) {
			expect(q(`host-processes-mcp-${a}`)!.disabled).toBe(true);
		}
	});

	it('restarts the MCP servers and shows the host as it now stands', async () => {
		await open();
		api.mockResolvedValueOnce(answer({ version: '1.9.0' }));
		api.mockResolvedValueOnce(host());
		q('host-processes-mcp-restart')!.click();
		await settle();
		expect(api.mock.calls[1][0]).toBe('/api/host');
		expect(JSON.parse(api.mock.calls[1][1].body)).toEqual({ action: 'restart', process: 'mcp' });
		expect(api.mock.calls[2]).toEqual(['/api/host']);
		expect(text('host-processes-message')).toBe('MCP: restarted.');
	});

	it('shows the host action’s refusal', async () => {
		await open();
		api.mockResolvedValueOnce(answer({ error: 'the host action "host" is not enabled' }, false));
		q('host-processes-mcp-stop')!.click();
		await settle();
		expect(text('host-processes-failed')).toBe('the host action "host" is not enabled');
	});

	it('asks before stopping serve, and says how to bring it back once it is gone', async () => {
		await open(host(), fast);
		q('host-processes-serve-stop')!.click();
		flushSync();
		expect(api).toHaveBeenCalledTimes(1);
		expect(text('host-processes-confirm-stop')).toContain('flai host start serve');

		api.mockImplementationOnce(() => new Promise(() => {}));
		q('host-processes-confirm-stop-yes')!.click();
		await settleThrough();
		expect(lastPost()).toEqual({ action: 'stop', process: 'serve' });
		expect(text('host-processes-none')).toContain('flai host start serve');
		expect(q('host-processes-serve-start')).toBeNull();
	});

	it('treats a serve restart that drops the connection as proceeding, and reconnects', async () => {
		await open(host(), fast);
		api.mockImplementationOnce(() => new Promise(() => {}));
		q('host-processes-serve-restart')!.click();
		await settleThrough();
		expect(q('host-processes-reconnecting')).not.toBeNull();
		// the old serve, still answering before it went down, is not taken for the new one
		api.mockResolvedValueOnce(host());
		const restarted = host({
			children: [{ name: 'serve', state: 'running', pid: 4200, version: '1.9.0', restarts: 0 }]
		});
		api.mockResolvedValue(restarted);
		await settleThrough(40);
		expect(text('host-processes-message')).toBe('serve restarted; reconnected.');
		expect(api.mock.calls.length).toBeGreaterThan(3);
	});

	it('takes the route’s 502 for serve going away as the restart going ahead', async () => {
		await open(host(), fast);
		api.mockResolvedValueOnce(
			answer({ error: 'the host flai went away before it answered' }, false, 502)
		);
		api.mockResolvedValue(
			host({ children: [{ name: 'serve', state: 'running', pid: 4300, restarts: 1 }] })
		);
		q('host-processes-serve-restart')!.click();
		await settleThrough(40);
		expect(q('host-processes-failed')).toBeNull();
		expect(text('host-processes-message')).toBe('serve restarted; reconnected.');
	});

	it('checks for a newer flai as a read', async () => {
		await open(host({ host_enabled: false }));
		api.mockResolvedValueOnce(answer({ current: '1.9.0', latest: '1.10.0', up_to_date: false }));
		q('host-processes-check')!.click();
		await settle();
		expect(lastPost()).toEqual({ action: 'check' });
		expect(text('host-processes-check-result')).toBe('flai 1.10.0 is available (running 1.9.0).');
	});

	it('upgrades, waits for the host to answer again, and reports the new version', async () => {
		await open(host(), fast);
		api.mockImplementationOnce(() => Promise.reject(new TypeError('network error')));
		q('host-processes-upgrade')!.click();
		await settleThrough(2);
		expect(lastPost()).toEqual({ action: 'upgrade' });
		api.mockResolvedValue(host({ version: '1.10.0' }));
		await settleThrough();
		expect(text('host-processes-message')).toBe(
			'Upgraded flai 1.9.0 to 1.10.0; serve and MCP restarted.'
		);
		expect(text('host-processes-host')).toBe('1.10.0, pid 4100');
	});

	it('waits past the old host after an upgrade that answered before restarting', async () => {
		await open(host(), fast);
		api.mockResolvedValueOnce(
			answer({ upgrade: { previous: '1.9.0', installed: '1.10.0' }, restarting: true })
		);
		api.mockResolvedValueOnce(host());
		api.mockResolvedValue(host({ version: '1.10.0' }));
		q('host-processes-upgrade')!.click();
		await settleThrough(40);
		expect(text('host-processes-message')).toBe(
			'Upgraded flai 1.9.0 to 1.10.0; serve and MCP restarted.'
		);
	});

	it('says so when the upgrade finds nothing newer, without waiting to reconnect', async () => {
		await open(host(), fast);
		api.mockResolvedValueOnce(
			answer({
				upgrade: { current: '1.9.0', latest: '1.9.0', up_to_date: true },
				restarting: false
			})
		);
		q('host-processes-upgrade')!.click();
		await settle();
		expect(text('host-processes-message')).toBe('flai 1.9.0 is already the latest.');
		expect(api).toHaveBeenCalledTimes(2);
	});

	it('hides the controls but keeps the check while the host action is off', async () => {
		await open(host({ host_enabled: false }));
		expect(q('host-processes-serve-restart')).toBeNull();
		expect(q('host-processes-upgrade')).toBeNull();
		expect(q('host-processes-check')).not.toBeNull();
		expect(text('host-processes-off')).toContain('flai serve enable host');
	});

	it('says how to start a host when none answers', async () => {
		await open(answer({ running: false, reason: 'flai host is not running', host_enabled: false }));
		expect(text('host-processes-none')).toContain('flai host is not running');
		expect(text('host-processes-none')).toContain('flai host start');
		expect(q('host-processes-check')).toBeNull();
	});
});
