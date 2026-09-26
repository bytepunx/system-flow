import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import SettingsPanel from './SettingsPanel.svelte';
import type { SettingsView } from '$lib/settings';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

// The switcher's list, which Serve asks again until the new project is in it (S-0122).
const switcher = vi.hoisted(() => ({
	list: [] as { key: string; name: string; connected: boolean; candidate?: boolean }[],
	refresh: async () => {}
}));
vi.mock('$lib/project.svelte', () => ({ projectState: switcher }));

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

function view(here: boolean, everywhere: boolean): SettingsView {
	return {
		here,
		everywhere,
		enable: 'flai serve enable settings',
		enable_everywhere: 'flai serve enable settings --all-projects',
		host: {
			actions: [
				{ name: 'push', means: 'push accepted work', here: false, everywhere: false },
				{ name: 'settings', means: 'change settings', here, everywhere }
			],
			default_agent: { harness: 'claude-code', model: 'claude-haiku-4-5' },
			agent: {
				name: 'builder',
				command: ['run-agent', '{story}'],
				harnesses: {
					'claude-code': {
						program: 'claude',
						args: ['--permission-mode', 'acceptEdits'],
						set: false
					}
				}
			},
			checks: { commands: [], timeout_minutes: 15 },
			manifest_checks: [{ name: 'unit', command: ['make', 'test'] }],
			import_roots: ['/home/me/git'],
			mcp: { running: true, url: 'http://127.0.0.1:4243/mcp' }
		}
	};
}

type Sent = { kind: string } & Record<string, unknown>;
function backend(v: SettingsView, answers: Record<string, [number, unknown]> = {}) {
	const sent: Sent[] = [];
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		if (init?.method === 'POST') {
			const body = JSON.parse(init.body!) as Sent;
			sent.push(body);
			const [status, res] = answers[body.kind] ?? [200, { ok: true }];
			return answer(status, res);
		}
		return answer(200, v);
	});
	return sent;
}

describe('SettingsPanel (S-0105)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});
	const show = async () => {
		c = mount(SettingsPanel, { target: document.body });
		await settle();
	};

	it('shows everything read-only, with the command that allows changing it, while settings is off', async () => {
		backend(view(false, false));
		await show();
		expect(q('settings-state')!.textContent).toContain('flai serve enable settings');
		expect(q<HTMLInputElement>('action-push')!.disabled).toBe(true);
		expect(q<HTMLButtonElement>('save-default-agent')!.disabled).toBe(true);
		expect(q<HTMLButtonElement>('save-agent')!.disabled).toBe(true);
		expect(q<HTMLTextAreaElement>('agent-command')!.value).toBe('run-agent\n{story}');
		expect(q('section-checks')!.textContent).toContain("the manifest's are used (unit)");
		expect(q('section-import')!.textContent).toContain('/home/me/git');
		expect(q('section-tokens')!.textContent).toContain('http://127.0.0.1:4243/mcp');
	});

	it('changes this project’s own settings when settings is on here, and not the host-wide ones', async () => {
		const sent = backend(view(true, false));
		await show();
		expect(q<HTMLInputElement>('action-settings')!.disabled).toBe(true);
		expect(q<HTMLButtonElement>('save-agent')!.disabled).toBe(true);
		expect(q('section-agent')!.textContent).toContain('flai serve enable settings --all-projects');
		const push = q<HTMLInputElement>('action-push')!;
		push.checked = true;
		push.dispatchEvent(new Event('change', { bubbles: true }));
		await settle();
		expect(sent[0]).toEqual({ kind: 'action', action: 'push', on: true });
		q<HTMLInputElement>('agent-model')!.value = 'claude-sonnet-5';
		q<HTMLInputElement>('agent-model')!.dispatchEvent(new Event('input', { bubbles: true }));
		q<HTMLButtonElement>('save-default-agent')!.click();
		await settle();
		expect(sent[1]).toEqual({
			kind: 'default_agent',
			agent: { harness: 'claude-code', model: 'claude-sonnet-5' }
		});
		expect(q('said-default')!.textContent).toContain('saved');
	});

	it('changes the host-wide settings when settings is on everywhere, and says a refusal', async () => {
		const sent = backend(view(true, true), {
			import: [400, { error: '"relative" is not an absolute, clean path' }],
			dashboard_token: [200, { login_url: 'http://x/login#token=abc' }]
		});
		await show();
		const cmd = q<HTMLTextAreaElement>('agent-command')!;
		cmd.value = 'claude\n-p\nwork on {story}; echo $HOME\n';
		cmd.dispatchEvent(new Event('input', { bubbles: true }));
		q<HTMLButtonElement>('save-agent')!.click();
		await settle();
		// one argument a line: nothing split, nothing quoted
		expect(sent[0]).toEqual({
			kind: 'agent',
			name: 'builder',
			command: ['claude', '-p', 'work on {story}; echo $HOME']
		});

		const args = q('harness-claude-code')!.querySelector<HTMLTextAreaElement>(
			'[data-testid="harness-args"]'
		)!;
		args.value = '--permission-mode\nbypassPermissions';
		args.dispatchEvent(new Event('input', { bubbles: true }));
		q('harness-claude-code')!
			.querySelector<HTMLButtonElement>('[data-testid="save-harness"]')!
			.click();
		await settle();
		expect(sent[1]).toEqual({
			kind: 'harness',
			name: 'claude-code',
			program: 'claude',
			args: ['--permission-mode', 'bypassPermissions']
		});

		q<HTMLInputElement>('new-check-name')!.value = 'unit';
		q<HTMLInputElement>('new-check-name')!.dispatchEvent(new Event('input', { bubbles: true }));
		q<HTMLTextAreaElement>('new-check-command')!.value = 'go\ntest\n./...';
		q<HTMLTextAreaElement>('new-check-command')!.dispatchEvent(
			new Event('input', { bubbles: true })
		);
		await settle();
		q<HTMLButtonElement>('add-check')!.click();
		await settle();
		expect(sent[2]).toEqual({ kind: 'check', name: 'unit', command: ['go', 'test', './...'] });

		q<HTMLInputElement>('new-folder')!.value = 'relative';
		q<HTMLInputElement>('new-folder')!.dispatchEvent(new Event('input', { bubbles: true }));
		await settle();
		q<HTMLButtonElement>('add-folder')!.click();
		await settle();
		expect(sent[3]).toEqual({ kind: 'import', folder: 'relative', add: true });
		expect(q('said-import')!.textContent).toContain('is not an absolute, clean path');
		expect(q<HTMLInputElement>('new-folder')!.value).toBe('relative');

		q<HTMLButtonElement>('rotate-dashboard')!.click();
		await settle();
		expect(q('once-dashboard')!.textContent).toContain('http://x/login#token=abc');
	});

	describe('projects (S-0122)', () => {
		const projects = (): SettingsView => {
			const v = view(true, false);
			v.host!.projects = {
				running: true,
				served: [
					{
						key: 'sf',
						name: 'system-flow',
						root: '/home/me/git/sf',
						from: 'registry',
						state: 'connected',
						since: '2026-09-26T06:00:00Z',
						settings: true
					},
					{
						key: 'blog',
						name: 'Blog',
						root: '/home/me/git/blog',
						from: 'registry',
						state: 'not-connected',
						last_error: 'dial: connection refused',
						settings: false
					},
					{
						key: 'notes',
						name: 'Notes',
						root: '/home/me/git/notes',
						from: 'import',
						below: '/home/me/git',
						state: 'connected',
						settings: true
					}
				],
				unserved: [
					{
						key: 'shop',
						name: 'Shop',
						root: '/home/me/git/shop',
						reason: 'no dashboard is known to serve it on',
						settings: true
					},
					{
						key: '',
						name: '',
						root: '/home/me/git/odd',
						reason: 'its system-flow.yaml does not load',
						settings: false
					},
					{
						key: 'diary',
						name: 'Diary',
						root: '/home/me/git/diary',
						reason: 'removed from the dashboard; flai serve project add /home/me/git/diary',
						removed: true,
						settings: true
					}
				]
			};
			return v;
		};

		it('lists each served project with its health, and the ones not served with why', async () => {
			backend(projects());
			await show();
			expect(q('served-sf')!.textContent).toContain('/home/me/git/sf');
			expect(q('served-sf')!.textContent).toContain('connected since 2026-09-26T06:00:00Z');
			expect(q('served-blog')!.textContent).toContain('not connected: dial: connection refused');
			expect(
				q('served-blog')!.querySelector<HTMLButtonElement>('[data-testid="remove-project"]')!
					.disabled
			).toBe(true);
			expect(q('served-blog')!.querySelector('[data-testid="gate"]')!.textContent).toContain(
				'flai serve enable settings'
			);
			// served from below an import folder: where it comes from, and Remove (S-0123)
			expect(q('served-notes')!.textContent).toContain('below /home/me/git');
			expect(
				q('served-notes')!.querySelector<HTMLButtonElement>('[data-testid="remove-project"]')!
					.disabled
			).toBe(false);
			// removed from below a folder: said so, and Serve brings it back
			expect(q('unserved-diary')!.textContent).toContain(
				'Removed from the dashboard; Serve brings it back.'
			);
			expect(
				q('unserved-diary')!.querySelector<HTMLButtonElement>('[data-testid="serve-project"]')!
					.disabled
			).toBe(false);
			expect(q('unserved-shop')!.textContent).toContain('no dashboard is known to serve it on');
			expect(
				q('unserved-shop')!.querySelector<HTMLButtonElement>('[data-testid="serve-project"]')!
					.disabled
			).toBe(false);
			expect(
				q('unserved-/home/me/git/odd')!.querySelector<HTMLButtonElement>(
					'[data-testid="serve-project"]'
				)!.disabled
			).toBe(true);
		});

		it('removes a project only once it is confirmed, and says how to add it back', async () => {
			const sent = backend(projects());
			const refresh = vi.spyOn(switcher, 'refresh');
			await show();
			q('served-sf')!.querySelector<HTMLButtonElement>('[data-testid="remove-project"]')!.click();
			flushSync();
			expect(q('confirm-remove')!.textContent).toContain('None of its files is touched');
			q<HTMLButtonElement>('confirm-remove-no')!.click();
			flushSync();
			expect(q('confirm-remove')).toBeNull();
			expect(sent).toEqual([]);

			q('served-sf')!.querySelector<HTMLButtonElement>('[data-testid="remove-project"]')!.click();
			flushSync();
			q<HTMLButtonElement>('confirm-remove-yes')!.click();
			await settle();
			expect(sent).toEqual([{ kind: 'unserve', root: '/home/me/git/sf', key: 'sf' }]);
			expect(q('said-projects')!.textContent).toContain('flai serve project add /home/me/git/sf');
			expect(refresh).toHaveBeenCalled();
			refresh.mockRestore();
		});

		it('removes a project served from below a folder, waits for the switcher to drop it, and says Serve brings it back (S-0123)', async () => {
			const sent = backend(projects(), {
				unserve: [200, { removed: { key: 'notes' }, listed_below: '/home/me/git' }]
			});
			switcher.list = [{ key: 'notes', name: 'Notes', connected: true }];
			let asked = 0;
			const refresh = vi.spyOn(switcher, 'refresh').mockImplementation(async () => {
				if (++asked === 2) switcher.list = [];
			});
			await show();
			q('served-notes')!
				.querySelector<HTMLButtonElement>('[data-testid="remove-project"]')!
				.click();
			flushSync();
			expect(q('confirm-remove')!.textContent!.replace(/\s+/g, ' ')).toContain(
				'It is below /home/me/git, so it is listed below as removed, and Serve there brings it back'
			);
			q<HTMLButtonElement>('confirm-remove-yes')!.click();
			for (let i = 0; i < 20 && asked < 2; i++) {
				await new Promise((r) => setTimeout(r, 100));
				flushSync();
			}
			await settle();
			expect(sent).toEqual([{ kind: 'unserve', root: '/home/me/git/notes', key: 'notes' }]);
			expect(q('said-projects')!.textContent).toContain(
				'It is listed below as removed, and Serve there brings it back.'
			);
			expect(asked).toBe(2);
			refresh.mockRestore();
			switcher.list = [];
		});

		it('serves a removed project again and waits for the switcher to have it (S-0123)', async () => {
			const sent = backend(projects(), {
				serve: [200, { project: { key: 'diary' }, served_below: '/home/me/git', restored: true }]
			});
			const refresh = vi.spyOn(switcher, 'refresh').mockImplementation(async () => {
				switcher.list = [{ key: 'diary', name: 'Diary', connected: true }];
			});
			await show();
			q('unserved-diary')!
				.querySelector<HTMLButtonElement>('[data-testid="serve-project"]')!
				.click();
			for (let i = 0; i < 20 && !q('said-projects')?.textContent?.includes('switcher'); i++) {
				await new Promise((r) => setTimeout(r, 100));
				flushSync();
			}
			expect(sent).toEqual([{ kind: 'serve', root: '/home/me/git/diary', key: 'diary' }]);
			expect(q('said-projects')!.textContent).toContain('diary is served, and in the switcher');
			refresh.mockRestore();
			switcher.list = [];
		});

		it('serves a project and waits for the switcher to have it', async () => {
			const sent = backend(projects());
			let asked = 0;
			const refresh = vi.spyOn(switcher, 'refresh').mockImplementation(async () => {
				if (++asked === 2) switcher.list = [{ key: 'shop', name: 'Shop', connected: true }];
			});
			await show();
			q('unserved-shop')!
				.querySelector<HTMLButtonElement>('[data-testid="serve-project"]')!
				.click();
			for (let i = 0; i < 20 && !q('said-projects')?.textContent?.includes('switcher'); i++) {
				await new Promise((r) => setTimeout(r, 100));
				flushSync();
			}
			expect(sent).toEqual([{ kind: 'serve', root: '/home/me/git/shop', key: 'shop' }]);
			expect(q('said-projects')!.textContent).toContain('shop is served, and in the switcher');
			expect(asked).toBe(2);
			refresh.mockRestore();
			switcher.list = [];
		});

		it('says flai’s refusal', async () => {
			backend(projects(), {
				unserve: [403, { error: 'the host action "settings" is not enabled for /home/me/git/sf' }]
			});
			await show();
			q('served-sf')!.querySelector<HTMLButtonElement>('[data-testid="remove-project"]')!.click();
			flushSync();
			q<HTMLButtonElement>('confirm-remove-yes')!.click();
			await settle();
			expect(q('said-projects')!.textContent).toContain('not enabled for /home/me/git/sf');
		});
	});
});
