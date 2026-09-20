import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostAgentNotice from './HostAgentNotice.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, json: async () => body });
const notice = () => document.querySelector<HTMLElement>('[data-testid="host-agent"]');
const text = () => notice()!.textContent!.replace(/\s+/g, ' ');
const run = {
	story: 'S-0079',
	command: 'claude',
	agent: 'builder',
	pid: 4242,
	started: '2026-09-20T16:00:00Z'
};

describe('HostAgentNotice (S-0079)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});
	const show = async (body: unknown) => {
		api.mockResolvedValue(answer(body));
		c = mount(HostAgentNotice, { target: document.body, props: {} });
		await settle();
	};

	it('says nothing while the operator has not enabled it, or nothing has happened', async () => {
		await show({ enabled: false, state: { command: 'claude', running: run } });
		expect(api).toHaveBeenCalledWith('/api/host-agent');
		expect(notice()).toBeNull();
		unmount(c!);
		c = undefined;
		await show({ enabled: true, state: { command: 'claude' } });
		expect(notice()).toBeNull();
	});

	it('says that an agent was started, for which story, by which command, as whom, and when', async () => {
		await show({ enabled: true, state: { command: 'claude', running: run } });
		expect(text()).toContain('Agent started for S-0079 by claude as builder, 2026-09-20 16:00 UTC');
		expect(text()).toContain('Stopping an agent is done on the host; the dashboard cannot');
		expect(document.querySelectorAll('button')).toHaveLength(0);
	});

	it('says when a command could not be started, and why a ready story waits', async () => {
		await show({
			enabled: true,
			state: {
				command: 'claude',
				last: { ...run, error: 'claude: executable file not found in $PATH' },
				waiting: 'the in-progress limit leaves no room to pull a story'
			}
		});
		expect(document.querySelector('[data-testid="host-agent-failed"]')!.textContent).toContain(
			'executable file not found'
		);
		expect(text()).toContain('flai serve agent set');
		expect(document.querySelector('[data-testid="host-agent-waiting"]')!.textContent).toContain(
			'leaves no room'
		);
	});

	it('says when the last one ended', async () => {
		await show({
			enabled: true,
			state: { command: 'claude', last: { ...run, ended: '2026-09-20T16:30:00Z', exit: 2 } }
		});
		expect(text()).toContain('ended 2026-09-20 16:30 UTC with exit code 2');
	});
});
