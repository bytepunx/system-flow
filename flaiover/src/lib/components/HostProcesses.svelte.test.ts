import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostProcesses from './HostProcesses.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const settleThrough = async (ms: number) => {
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
// Real timers, cut short. Under machine load a 5 ms timer can take far longer, so the give-up is
// seconds away and the tests wait for their outcome rather than a fixed time (I-0053).
const fast = { disconnectTimeoutMs: 5, reconnectPollMs: 5, reconnectGiveUpMs: 3000 };
const eventually = (assert: () => void) =>
	vi.waitFor(
		() => {
			flushSync();
			assert();
		},
		{ timeout: 4000, interval: 5 }
	);
// an answer that takes ms to come, as a poll does on a loaded machine
const slowly = (ms: number, value: unknown) =>
	new Promise((resolve) => setTimeout(() => resolve(value), ms));

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

	it('keeps Process narrow and gives State, Version, and Restarts equal, padded columns', async () => {
		for (const enabled of [true, false]) {
			await open(host({ host_enabled: enabled }));
			expect(q('host-processes-table')!.classList).toContain('table-fixed');
			const [process, ...rest] = [...document.querySelectorAll('thead th')];
			expect(process.textContent).toBe('Process');
			expect(process.classList).toContain('w-28');
			const [state, version, restarts, actions] = rest;
			for (const th of [state, version, restarts]) {
				// no width of their own: a fixed table shares what is left between them equally
				expect([...th.classList].filter((c) => /^w-/.test(c))).toEqual([]);
				expect(th.classList).toContain('px-4');
			}
			expect(actions?.classList.contains('w-52') ?? false).toBe(enabled);
			if (c) unmount(c);
			c = undefined;
			document.body.innerHTML = '';
		}
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
		await eventually(() => expect(text('host-processes-none')).toContain('flai host start serve'));
		expect(lastPost()).toEqual({ action: 'stop', process: 'serve' });
		expect(q('host-processes-serve-start')).toBeNull();
	});

	it('treats a serve restart that drops the connection as proceeding, and reconnects', async () => {
		await open(host(), fast);
		api.mockImplementationOnce(() => new Promise(() => {}));
		q('host-processes-serve-restart')!.click();
		await eventually(() => expect(q('host-processes-reconnecting')).not.toBeNull());
		// the old serve, still answering before it went down, is not taken for the new one
		api.mockResolvedValueOnce(host());
		const restarted = host({
			children: [{ name: 'serve', state: 'running', pid: 4200, version: '1.9.0', restarts: 0 }]
		});
		api.mockResolvedValue(restarted);
		await eventually(() =>
			expect(text('host-processes-message')).toBe('serve restarted; reconnected.')
		);
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
		await eventually(() =>
			expect(text('host-processes-message')).toBe('serve restarted; reconnected.')
		);
		expect(q('host-processes-failed')).toBeNull();
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
		await eventually(() =>
			expect(text('host-processes-message')).toBe(
				'Upgraded flai 1.9.0 to 1.10.0; serve and MCP restarted.'
			)
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
		await eventually(() =>
			expect(text('host-processes-message')).toBe(
				'Upgraded flai 1.9.0 to 1.10.0; serve and MCP restarted.'
			)
		);
	});

	it('waits past an old host whose answer comes slowly, as on a loaded machine', async () => {
		// I-0053: the old host's answer took longer than the fixed 40 ms the test waited and the
		// 50 ms give-up, so the test read the page before the new host answered
		await open(host(), fast);
		api.mockResolvedValueOnce(
			answer({ upgrade: { previous: '1.9.0', installed: '1.10.0' }, restarting: true })
		);
		api.mockImplementationOnce(() => slowly(80, host()));
		api.mockResolvedValue(host({ version: '1.10.0' }));
		q('host-processes-upgrade')!.click();
		await eventually(() =>
			expect(text('host-processes-message')).toBe(
				'Upgraded flai 1.9.0 to 1.10.0; serve and MCP restarted.'
			)
		);
		expect(q('host-processes-failed')).toBeNull();
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

	const releases = () =>
		answer([
			{
				version: '1.10.0',
				tag: 'flai/v1.10.0',
				// just after midnight in UTC, the day before in New York, where the tests run (S-0329)
				published: '2026-10-01T02:00:00Z',
				installed: false,
				latest: true
			},
			{ version: '1.9.0', tag: 'flai/v1.9.0', installed: true, latest: false },
			{
				version: '1.8.0',
				tag: 'flai/v1.8.0',
				installed: false,
				latest: false,
				below_minimum: [{ project: 'harbour', minimum: '1.9.0' }]
			}
		]);
	const rows = () =>
		[...document.querySelectorAll('ul[aria-label="Published flai releases"] li')].map((li) =>
			li.textContent!.replace(/\s+/g, ' ').trim()
		);
	const deploy = (version: string) =>
		document.querySelector<HTMLButtonElement>(`button[aria-label="Deploy ${version}"]`);

	async function listed(first = host(), props = {}) {
		await open(first, props);
		api.mockResolvedValueOnce(releases());
		q('host-processes-versions-list')!.click();
		await settle();
	}

	it('lists the published releases newest first, marking the installed, newest, and below a minimum', async () => {
		await listed();
		expect(lastPost()).toEqual({ action: 'versions' });
		expect(rows()).toEqual([
			'1.10.0 2026-09-30 newest Deploy',
			'1.9.0 installed',
			"1.8.0 below harbour's minimum 1.9.0 Deploy"
		]);
		expect(deploy('1.9.0')).toBeNull();
	});

	it('lists the releases while the host action is off, offering no deploy', async () => {
		await listed(host({ host_enabled: false }));
		expect(rows()).toHaveLength(3);
		expect(deploy('1.10.0')).toBeNull();
		expect(deploy('1.8.0')).toBeNull();
	});

	it('deploys a chosen release once confirmed, waits for the host, and reports the version', async () => {
		await listed(host(), fast);
		deploy('1.10.0')!.click();
		flushSync();
		expect(api).toHaveBeenCalledTimes(2);
		expect(text('host-processes-confirm-deploy')).toContain('Install flai 1.10.0 on the host?');
		expect(text('host-processes-confirm-deploy')).not.toContain('flai.minimum');

		api.mockImplementationOnce(() => Promise.reject(new TypeError('network error')));
		q('host-processes-confirm-deploy-yes')!.click();
		await settleThrough(2);
		expect(lastPost()).toEqual({ action: 'upgrade', version: '1.10.0' });
		api.mockResolvedValue(host({ version: '1.10.0' }));
		await eventually(() =>
			expect(text('host-processes-message')).toBe(
				'Deployed flai 1.10.0 over 1.9.0; serve and MCP restarted.'
			)
		);
		expect(q('host-processes-confirm-deploy')).toBeNull();
	});

	it('warns before going back below a project’s flai.minimum, and deploys it when confirmed', async () => {
		await listed(host(), fast);
		deploy('1.8.0')!.click();
		flushSync();
		const said = text('host-processes-confirm-deploy')!;
		expect(said).toContain("below harbour's flai.minimum, 1.9.0");
		expect(said).toContain('flai serve leaves harbour unserved');

		api.mockResolvedValueOnce(
			answer({ upgrade: { previous: '1.9.0', installed: '1.8.0' }, restarting: true })
		);
		api.mockResolvedValue(host({ version: '1.8.0' }));
		q('host-processes-confirm-deploy-yes')!.click();
		await eventually(() =>
			expect(text('host-processes-message')).toBe(
				'Deployed flai 1.8.0 over 1.9.0; serve and MCP restarted.'
			)
		);
		expect(JSON.parse(api.mock.calls[2][1].body)).toEqual({ action: 'upgrade', version: '1.8.0' });
	});

	it('shows flai’s refusal of a release it does not list as published', async () => {
		await listed();
		deploy('1.10.0')!.click();
		flushSync();
		api.mockResolvedValueOnce(
			answer({ error: 'flai 1.10.0 is not published; published: 1.9.0, 1.8.0' }, false)
		);
		q('host-processes-confirm-deploy-yes')!.click();
		await settle();
		expect(text('host-processes-failed')).toBe(
			'flai 1.10.0 is not published; published: 1.9.0, 1.8.0'
		);
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
