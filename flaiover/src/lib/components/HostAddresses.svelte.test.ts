import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostAddresses from './HostAddresses.svelte';
import type { SettingsView } from '$lib/settings';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 8; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const q = <T extends HTMLElement = HTMLElement>(id: string) =>
	document.querySelector<T>(`[data-testid="${id}"]`);
const answer = (status: number, body: unknown) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});

function view(everywhere: boolean, allowed_hosts: string): SettingsView {
	return {
		here: true,
		everywhere,
		enable: 'flai serve enable settings',
		enable_everywhere: 'flai serve enable settings --all-projects',
		host: {
			actions: [],
			agent: { name: '', command: [], harnesses: {} },
			checks: { commands: [], timeout_minutes: 15 },
			import_roots: [],
			allowed_hosts
		}
	};
}

/**
 * The host as the area asks it: settings.get with the configured list, which a save changes,
 * dashboard.status with the list the container runs with, and the save's answer.
 */
function host(opts: {
	everywhere?: boolean;
	configured?: string;
	running?: string;
	dashboardEnabled?: boolean;
	save?: [number, unknown];
}) {
	let configured = opts.configured ?? '';
	api.mockImplementation(async (path: string, init?: RequestInit) => {
		if (path === '/api/dashboard') {
			return answer(200, {
				container: 'flaiover',
				running: true,
				allowed_hosts: opts.running ?? '',
				dashboard_enabled: opts.dashboardEnabled ?? true
			});
		}
		if (path === '/api/settings' && init?.method === 'POST') {
			const [status, body] = opts.save ?? [200, { warnings: [] }];
			if (status < 400) configured = JSON.parse(String(init.body)).hosts.replace(/\s/g, '');
			return answer(status, body);
		}
		return answer(200, view(opts.everywhere ?? true, configured));
	});
}

const posts = () =>
	api.mock.calls.filter(([, init]) => (init as RequestInit | undefined)?.method === 'POST');

describe('HostAddresses (S-0489)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	async function open() {
		c = mount(HostAddresses, { target: document.body });
		flushSync();
		await settle();
	}

	it('shows the list and the enabling command, and offers no Save, while settings is off host-wide', async () => {
		host({ everywhere: false, configured: 'flow.example.com,board.example.net:8443' });
		await open();
		expect(q('host-addresses')!.querySelector('h2')!.textContent).toBe('Addresses');
		expect(q('host-addresses-list')!.textContent).toContain(
			'flow.example.com, board.example.net:8443'
		);
		expect(q('host-addresses-gate')!.textContent).toContain(
			'flai serve enable settings --all-projects'
		);
		expect(q('host-addresses-field')).toBeNull();
		expect(q('host-addresses-save')).toBeNull();
		expect(document.querySelectorAll('button')).toHaveLength(0);
	});

	it('saves the list through settings.dashboard_hosts, shows it, and says a restart applies it', async () => {
		host({});
		await open();
		const field = q<HTMLInputElement>('host-addresses-field')!;
		expect(field.value).toBe('');
		field.value = 'flow.example.com';
		field.dispatchEvent(new Event('input'));
		flushSync();
		q<HTMLButtonElement>('host-addresses-save')!.click();
		await settle();
		expect(posts()).toHaveLength(1);
		const [path, init] = posts()[0] as [string, RequestInit];
		expect(path).toBe('/api/settings');
		expect(JSON.parse(String(init.body))).toEqual({
			kind: 'dashboard_hosts',
			hosts: 'flow.example.com'
		});
		expect(q<HTMLInputElement>('host-addresses-field')!.value).toBe('flow.example.com');
		const saved = q('host-addresses-saved')!.textContent!;
		expect(saved).toContain('keeps its old list until it is restarted');
		expect(saved).toContain('Restart in the Dashboard area above');
		expect(q('host-addresses-refused')).toBeNull();
	});

	it('names the dashboard host action when Restart is not offered above', async () => {
		host({ dashboardEnabled: false });
		await open();
		q<HTMLButtonElement>('host-addresses-save')!.click();
		await settle();
		expect(q('host-addresses-saved')!.textContent).toContain('flai serve enable dashboard');
	});

	it("shows flai's refusal under the field and says nothing was saved", async () => {
		const why =
			'dashboard.allowed_hosts: "https://flow.example.com" is not a host name: write no scheme, path, or spaces';
		host({ configured: 'board.example.net', save: [400, { error: why }] });
		await open();
		const field = q<HTMLInputElement>('host-addresses-field')!;
		field.value = 'https://flow.example.com';
		field.dispatchEvent(new Event('input'));
		flushSync();
		q<HTMLButtonElement>('host-addresses-save')!.click();
		await settle();
		expect(q('host-addresses-refused')!.textContent).toContain(why);
		expect(q('host-addresses-saved')).toBeNull();
		expect(q<HTMLInputElement>('host-addresses-field')!.value).toBe('https://flow.example.com');
	});

	it('says when the running dashboard was started with another list', async () => {
		host({ configured: 'flow.example.com', running: '' });
		await open();
		const differs = q('host-addresses-differs')!.textContent!;
		expect(differs).toContain('started with none');
		expect(differs).toContain('a restart applies flow.example.com');
	});

	it('says nothing of the running list when it is the configured one', async () => {
		host({ configured: 'flow.example.com', running: 'flow.example.com' });
		await open();
		expect(q('host-addresses-differs')).toBeNull();
	});
});
