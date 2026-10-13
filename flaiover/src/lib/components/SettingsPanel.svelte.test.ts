import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import SettingsPanel from './SettingsPanel.svelte';
import type { PlanningTriggers, SettingsView } from '$lib/settings';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({ resolve: (route: string) => route }));

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
			args: ['--permission-mode', 'bypassPermissions'],
			guard: '',
			deny_protected: false
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

	describe('a harness’s guard and deny_protected (S-0355, ADR-0131)', () => {
		const lifted = (): SettingsView => {
			const v = view(true, true);
			const harnesses = {
				'claude-code': {
					program: 'claude',
					args: ['--permission-mode', 'acceptEdits'],
					set: false,
					guard: '',
					deny_protected: false,
					guard_capable: true,
					ask_capable: true
				},
				command: {
					program: 'run-agent',
					args: ['{story}'],
					set: false,
					guard: 'none',
					deny_protected: true,
					guard_capable: false,
					ask_capable: false
				}
			};
			v.host!.agent.harnesses = harnesses;
			return v;
		};
		const within = <T extends HTMLElement>(box: string, id: string) =>
			q(box)!.querySelector<T>(`[data-testid="${id}"]`)!;

		it('shows each harness’s guard and deny_protected, with what each gives up', async () => {
			backend(lifted());
			await show();
			expect(within<HTMLSelectElement>('harness-claude-code', 'harness-guard').value).toBe('');
			const ccDeny = within<HTMLInputElement>('harness-claude-code', 'harness-deny-protected');
			expect(ccDeny.checked).toBe(false);
			expect(within('harness-claude-code', 'harness-guard-gives-up').textContent).toContain(
				'claude-code runs flai guard, so this changes nothing here'
			);
			expect(within('harness-claude-code', 'harness-deny-gives-up').textContent).toContain(
				'claude-code asks before a protected write'
			);
			expect(q('harness-claude-code')!.textContent).toContain('(defaults)');
			const ccReset = within<HTMLButtonElement>('harness-claude-code', 'reset-harness');
			expect(ccReset.disabled).toBe(true);

			expect(q('harness-command')!.textContent).toContain('harness command: the command above');
			expect(q('harness-command')!.querySelector('[data-testid="harness-program"]')).toBeNull();
			expect(q('harness-command')!.querySelector('[data-testid="harness-args"]')).toBeNull();
			expect(within<HTMLSelectElement>('harness-command', 'harness-guard').value).toBe('none');
			const cmdDeny = within<HTMLInputElement>('harness-command', 'harness-deny-protected');
			expect(cmdDeny.checked).toBe(true);
			const means = within('harness-command', 'harness-guard-gives-up').textContent!;
			expect(means).toContain(
				"a story's sub-agents run without flai guard, so one could move items, write threads, or commit"
			);
			expect(means).not.toContain('changes nothing here');
			expect(within('harness-command', 'harness-deny-gives-up').textContent).toContain(
				'you say its arguments deny the protected paths; flai cannot check them'
			);
			expect(within<HTMLButtonElement>('harness-command', 'reset-harness').disabled).toBe(false);
		});

		it('saves the command’s guard and deny_protected alone, and a harness’s beside its program', async () => {
			const sent = backend(lifted());
			await show();
			const guard = within<HTMLSelectElement>('harness-command', 'harness-guard');
			guard.value = '';
			guard.dispatchEvent(new Event('change', { bubbles: true }));
			const deny = within<HTMLInputElement>('harness-command', 'harness-deny-protected');
			deny.checked = false;
			deny.dispatchEvent(new Event('change', { bubbles: true }));
			within<HTMLButtonElement>('harness-command', 'save-harness').click();
			await settle();
			expect(sent[0]).toEqual({
				kind: 'harness',
				name: 'command',
				guard: '',
				deny_protected: false
			});

			const cc = within<HTMLSelectElement>('harness-claude-code', 'harness-guard');
			cc.value = 'none';
			cc.dispatchEvent(new Event('change', { bubbles: true }));
			within<HTMLButtonElement>('harness-claude-code', 'save-harness').click();
			await settle();
			expect(sent[1]).toEqual({
				kind: 'harness',
				name: 'claude-code',
				program: 'claude',
				args: ['--permission-mode', 'acceptEdits'],
				guard: 'none',
				deny_protected: false
			});

			within<HTMLButtonElement>('harness-command', 'reset-harness').click();
			await settle();
			expect(sent[2]).toEqual({ kind: 'harness', name: 'command', reset: true });
		});
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
			expect(q('served-sf')!.textContent).toContain('connected since 2026-09-26 02:00 EDT');
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

	// S-0359, ADR-0129: the host's providers, set under the settings action on for every project; a
	// key is never shown or asked for, only its variable's name.
	describe('providers', () => {
		const providers = (here: boolean, everywhere: boolean) => {
			const v = view(here, everywhere);
			v.host!.default_agent = { harness: 'claude-code', model: 'opus', provider: 'litellm' };
			v.host!.agent.providers = {
				litellm: {
					api: 'anthropic-messages',
					base_url: 'http://127.0.0.1:4000',
					key_env: 'LITELLM_API_KEY',
					models: {
						opus: 'anthropic/claude-opus-5-5',
						haiku: 'anthropic/claude-haiku-4-5'
					},
					key_set: true
				},
				openrouter: {
					api: 'openai-chat',
					base_url: 'https://openrouter.ai/api/v1',
					key_env: 'OPENROUTER_API_KEY',
					models: {},
					key_set: false
				}
			};
			return v;
		};
		const box = (name: string) => q(`provider-${name}`)!;
		const field = <T extends HTMLElement>(within: HTMLElement, id: string) =>
			within.querySelector<T>(`[data-testid="${id}"]`)!;
		const type = (
			el: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement,
			value: string
		) => {
			el.value = value;
			el.dispatchEvent(
				new Event(el instanceof HTMLSelectElement ? 'change' : 'input', { bubbles: true })
			);
			flushSync();
		};

		it('shows each provider and whether its key’s variable is set, read-only until settings is on everywhere', async () => {
			const sent = backend(providers(true, false));
			await show();
			const litellm = box('litellm');
			expect(field<HTMLSelectElement>(litellm, 'provider-api').value).toBe('anthropic-messages');
			const apis = field<HTMLSelectElement>(litellm, 'provider-api').options;
			expect([...apis].filter((o) => !o.disabled).map((o) => o.value)).toEqual([
				'anthropic-messages',
				'openai-chat',
				'openai-responses'
			]);
			expect(field<HTMLInputElement>(litellm, 'provider-base-url').value).toBe(
				'http://127.0.0.1:4000'
			);
			expect(field<HTMLInputElement>(litellm, 'provider-key-env').value).toBe('LITELLM_API_KEY');
			expect(field<HTMLTextAreaElement>(litellm, 'provider-models').value).toBe(
				'haiku=anthropic/claude-haiku-4-5\nopus=anthropic/claude-opus-5-5'
			);
			expect(field(litellm, 'provider-key').textContent!.trim()).toBe('key set');
			expect(field(box('openrouter'), 'provider-key').textContent).toContain(
				"OPENROUTER_API_KEY is not set in flai serve's environment"
			);
			expect(q('section-providers')!.querySelector('[data-testid="gate"]')!.textContent).toContain(
				'flai serve enable settings --all-projects'
			);
			expect(field<HTMLButtonElement>(litellm, 'save-provider').disabled).toBe(true);
			expect(field<HTMLButtonElement>(litellm, 'reset-provider').disabled).toBe(true);
			type(q<HTMLInputElement>('new-provider-name')!, 'openai');
			expect(q<HTMLButtonElement>('add-provider')!.disabled).toBe(true);
			expect(sent).toEqual([]);
		});

		it('says when there are none', async () => {
			backend(view(true, true));
			await show();
			expect(q('providers-none')!.textContent).toContain('None');
		});

		it('saves a provider’s fields, its models one alias=name a line', async () => {
			const sent = backend(providers(true, true));
			await show();
			const litellm = box('litellm');
			expect(q('section-providers')!.querySelector('[data-testid="gate"]')).toBeNull();
			type(field<HTMLSelectElement>(litellm, 'provider-api'), 'openai-chat');
			type(field<HTMLInputElement>(litellm, 'provider-base-url'), ' https://llm.example.com ');
			type(field<HTMLTextAreaElement>(litellm, 'provider-models'), 'opus = openai/gpt-6\n\n');
			field<HTMLButtonElement>(litellm, 'save-provider').click();
			await settle();
			expect(sent).toEqual([
				{
					kind: 'provider',
					name: 'litellm',
					api: 'openai-chat',
					base_url: 'https://llm.example.com',
					key_env: 'LITELLM_API_KEY',
					models: { opus: 'openai/gpt-6' }
				}
			]);
			expect(q('said-provider-litellm')!.textContent).toContain('saved');
		});

		it('sends no models when they are emptied, and says a line that is not alias=name without sending it', async () => {
			const sent = backend(providers(true, true));
			await show();
			const litellm = box('litellm');
			type(field<HTMLTextAreaElement>(litellm, 'provider-models'), 'opus');
			field<HTMLButtonElement>(litellm, 'save-provider').click();
			await settle();
			expect(sent).toEqual([]);
			expect(q('said-provider-litellm')!.getAttribute('role')).toBe('alert');
			expect(q('said-provider-litellm')!.textContent).toContain('models: "opus" is not key=value');

			type(field<HTMLTextAreaElement>(litellm, 'provider-models'), '');
			field<HTMLButtonElement>(litellm, 'save-provider').click();
			await settle();
			expect(sent).toEqual([
				{
					kind: 'provider',
					name: 'litellm',
					api: 'anthropic-messages',
					base_url: 'http://127.0.0.1:4000',
					key_env: 'LITELLM_API_KEY'
				}
			]);
		});

		it('removes a provider whole', async () => {
			const sent = backend(providers(true, true));
			await show();
			field<HTMLButtonElement>(box('openrouter'), 'reset-provider').click();
			await settle();
			expect(sent).toEqual([{ kind: 'provider', name: 'openrouter', reset: true }]);
		});

		it('adds a provider by name and clears the form, and keeps it when flai refuses one', async () => {
			const sent = backend(providers(true, true), {
				provider: [
					400,
					{
						error:
							'agent.providers.groq is not complete: give api, base_url, and key_env for a new provider'
					}
				]
			});
			await show();
			const add = q('new-provider')!;
			type(q<HTMLInputElement>('new-provider-name')!, 'groq');
			type(field<HTMLInputElement>(add, 'provider-base-url'), 'https://api.groq.com/openai/v1');
			q<HTMLButtonElement>('add-provider')!.click();
			await settle();
			expect(sent).toEqual([
				{ kind: 'provider', name: 'groq', base_url: 'https://api.groq.com/openai/v1' }
			]);
			expect(q('said-providers')!.textContent).toContain('is not complete');
			expect(q<HTMLInputElement>('new-provider-name')!.value).toBe('groq');

			const again = backend(providers(true, true));
			type(field<HTMLSelectElement>(add, 'provider-api'), 'openai-chat');
			type(field<HTMLInputElement>(add, 'provider-key-env'), 'GROQ_API_KEY');
			type(field<HTMLTextAreaElement>(add, 'provider-models'), 'sonnet=llama-4');
			q<HTMLButtonElement>('add-provider')!.click();
			await settle();
			expect(again).toEqual([
				{
					kind: 'provider',
					name: 'groq',
					api: 'openai-chat',
					base_url: 'https://api.groq.com/openai/v1',
					key_env: 'GROQ_API_KEY',
					models: { sonnet: 'llama-4' }
				}
			]);
			expect(q<HTMLInputElement>('new-provider-name')!.value).toBe('');
			expect(field<HTMLInputElement>(add, 'provider-base-url').value).toBe('');
		});

		it('shows the default agent’s provider and sends the one typed', async () => {
			const sent = backend(providers(true, false));
			await show();
			const provider = q<HTMLInputElement>('agent-provider')!;
			expect(provider.value).toBe('litellm');
			type(provider, ' openrouter ');
			q<HTMLButtonElement>('save-default-agent')!.click();
			await settle();
			expect(sent).toEqual([
				{
					kind: 'default_agent',
					agent: { harness: 'claude-code', model: 'opus', provider: 'openrouter' }
				}
			]);
		});
	});

	// S-0295, ADR-0096: the shared paths, changed under the settings action and checked by anyone.
	describe('shared paths', () => {
		const shared = (here: boolean, patterns = ['docs/users/*.md', 'design/adrs']) => {
			const v = view(here, false);
			v.host!.shared = patterns;
			return v;
		};
		const type = (id: string, value: string) => {
			q<HTMLInputElement>(id)!.value = value;
			q<HTMLInputElement>(id)!.dispatchEvent(new Event('input', { bubbles: true }));
			flushSync();
		};
		const removeOf = (pattern: string) =>
			q(`shared-${pattern}`)!.querySelector<HTMLButtonElement>('[data-testid="remove-shared"]');

		it('lists the patterns read-only while settings is off, with no add field or remove button', async () => {
			backend(shared(false));
			await show();
			expect(q('shared-docs/users/*.md')!.textContent).toContain('docs/users/*.md');
			expect(q('shared-design/adrs')!.textContent).toContain('design/adrs');
			expect(removeOf('design/adrs')).toBeNull();
			expect(q('new-shared')).toBeNull();
			expect(q('add-shared')).toBeNull();
			expect(q('section-shared')!.querySelector('[data-testid="gate"]')!.textContent).toContain(
				'flai serve enable settings'
			);
			expect(q('section-shared')!.textContent).toContain(
				'An overlap inside a shared path holds no ready story'
			);
			// the check changes nothing, so it stays
			expect(q('shared-check-text')).not.toBeNull();
		});

		it('says when there are none, and leaves the section out when flai does not say', async () => {
			backend(shared(true, []));
			await show();
			expect(q('shared-none')!.textContent).toContain('None');
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
			backend(view(true, false));
			await show();
			expect(q('section-shared')).toBeNull();
		});

		it('adds a pattern, explaining the glob dialect beside the field, and clears the field', async () => {
			const sent = backend(shared(true), {
				shared: [200, { added: ['flai/cmd/**/*.go'], shared: [], warnings: [] }]
			});
			await show();
			expect(q('section-shared')!.textContent!.replace(/\s+/g, ' ')).toContain(
				'* within a folder, ** across folders, and a plain path covers everything under it'
			);
			type('new-shared', ' flai/cmd/**/*.go ');
			q<HTMLButtonElement>('add-shared')!.click();
			await settle();
			expect(sent).toEqual([{ kind: 'shared', action: 'add', pattern: 'flai/cmd/**/*.go' }]);
			expect(q('said-shared')!.textContent).toContain('added flai/cmd/**/*.go');
			expect(q<HTMLInputElement>('new-shared')!.value).toBe('');
		});

		it('removes a pattern', async () => {
			const sent = backend(shared(true), {
				shared: [200, { removed: ['design/adrs'], shared: ['docs/users/*.md'], warnings: [] }]
			});
			await show();
			expect(removeOf('design/adrs')!.getAttribute('aria-label')).toBe('Remove design/adrs');
			removeOf('design/adrs')!.click();
			await settle();
			expect(sent).toEqual([{ kind: 'shared', action: 'remove', pattern: 'design/adrs' }]);
			expect(q('said-shared')!.textContent).toContain('removed design/adrs');
		});

		it('shows flai’s refusal of a pattern inline, without the manifest’s path, and keeps the field', async () => {
			backend(shared(true), {
				shared: [
					400,
					{
						error:
							'/home/me/git/sf/system-flow.yaml: cannot add "design/adrs" to claims.shared: it is in the list already; nothing to add'
					}
				]
			});
			await show();
			type('new-shared', 'design/adrs');
			q<HTMLButtonElement>('add-shared')!.click();
			await settle();
			const said = q('said-shared')!;
			expect(said.getAttribute('role')).toBe('alert');
			expect(said.textContent!.trim()).toBe(
				'cannot add "design/adrs" to claims.shared: it is in the list already; nothing to add'
			);
			expect(q<HTMLInputElement>('new-shared')!.value).toBe('design/adrs');
		});

		it('checks a typed path, with settings off, and says which pattern matched or that none did', async () => {
			const sent = backend(shared(false), {
				shared_check: [
					200,
					{
						entries: [
							{
								entry: 'docs/users/flai.md',
								path: 'docs/users/flai.md',
								shared: true,
								pattern: 'docs/users/*.md'
							}
						],
						warnings: []
					}
				]
			});
			await show();
			type('shared-check-text', ' docs/users/flai.md ');
			q<HTMLButtonElement>('shared-check')!.click();
			await settle();
			expect(sent).toEqual([{ kind: 'shared_check', paths: ['docs/users/flai.md'] }]);
			expect(q('checked-docs/users/flai.md')!.textContent!.replace(/\s+/g, ' ')).toContain(
				'docs/users/flai.md: shared, inside docs/users/*.md'
			);

			backend(shared(false), {
				shared_check: [
					200,
					{
						entries: [{ entry: 'flai/cmd/', path: 'flai/cmd', shared: false }],
						warnings: []
					}
				]
			});
			type('shared-check-text', 'flai/cmd/');
			q<HTMLButtonElement>('shared-check')!.click();
			await settle();
			expect(q('checked-flai/cmd/')!.textContent!.replace(/\s+/g, ' ')).toContain(
				'flai/cmd/ (as flai/cmd): not shared: no pattern matches it'
			);
		});

		it('checks a story’s claim entry by entry, and says flai’s refusal of one it cannot find', async () => {
			const sent = backend(shared(true), {
				shared_check: [
					200,
					{
						entries: [
							{
								entry: 'design/adrs',
								path: 'design/adrs',
								shared: true,
								pattern: 'design/adrs',
								story: 'S-0001'
							},
							{ entry: 'flai/cmd', path: 'flai/cmd', shared: false, story: 'S-0001' }
						],
						warnings: []
					}
				]
			});
			await show();
			type('shared-check-text', 'S-0001');
			q<HTMLButtonElement>('shared-check')!.click();
			await settle();
			expect(sent).toEqual([{ kind: 'shared_check', story: 'S-0001' }]);
			expect(q('shared-checked')!.textContent).toContain("S-0001's claim:");
			expect(q('checked-design/adrs')!.textContent!.replace(/\s+/g, ' ')).toContain(
				'shared, inside design/adrs'
			);
			expect(q('checked-flai/cmd')!.textContent).toContain('not shared');

			backend(shared(true), { shared_check: [404, { error: 'no story S-0999' }] });
			type('shared-check-text', 'S-0999');
			q<HTMLButtonElement>('shared-check')!.click();
			await settle();
			expect(q('shared-checked')!.getAttribute('role')).toBe('alert');
			expect(q('shared-checked')!.textContent).toContain('no story S-0999');
		});
	});

	// S-0211, ADR-0084: when flai serve plans again on its own, read-only, from the manifest and the
	// plan host action.
	describe('planning triggers', () => {
		const planning = (p: PlanningTriggers) => {
			const v = view(false, false);
			v.host!.planning = p;
			return v;
		};
		const text = (id: string) => (q(id)!.textContent ?? '').replace(/\s+/g, ' ');

		it('shows whether edits replan, the default policy and what it means, and the schedule with its next run', async () => {
			backend(
				planning({
					plan: true,
					replan: 'deterministic',
					replan_set: false,
					schedule: '0 6 * * 1-5',
					next: '2026-10-05T06:00:00Z'
				})
			);
			await show();
			// S-0229: the keys are edited on the Planner page, not by hand
			expect(text('section-planning')).toContain('edited on the Planner page');
			expect(q('planning-edit-link')!.getAttribute('href')).toBe('/workflow/planner');
			expect(text('planning-edits')).toContain('Edits to a planned story start the planner: on');
			expect(text('planning-edits')).toContain('The plan host action is on for this project');
			expect(text('planning-replan')).toContain('deterministic (default)');
			expect(text('planning-replan')).toContain('forecast deliveries follow, with no agent');
			// S-0329: the expression is read in UTC, and its next run is a time, shown local
			expect(text('planning-schedule')).toContain('in UTC: 0 6 * * 1-5');
			expect(text('planning-schedule')).toContain('0 6 * * 1-5, next run 2026-10-05 02:00 EDT');
			expect(text('planning-schedule')).not.toContain('2026-10-05T06:00:00Z');
			expect(text('planning-schedule').split('next run')[1]).not.toContain('UTC');
		});

		it('says none without a schedule, a set policy without (default), and what flai cannot read', async () => {
			backend(
				planning({
					plan: false,
					replan: 'agent',
					replan_set: true,
					schedule: ''
				})
			);
			await show();
			expect(text('planning-edits')).toContain('start the planner: off');
			expect(text('planning-replan')).toContain('agent');
			expect(text('planning-replan')).not.toContain('(default)');
			expect(text('planning-replan')).toContain(
				'the planner runs for each story whose delivery moved'
			);
			expect(text('planning-schedule')).toContain('in UTC: none');
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
			backend(
				planning({
					plan: true,
					replan: 'sometimes',
					replan_set: true,
					replan_error: 'planning.replan "sometimes" is not a replan policy',
					schedule: '61 * * * *',
					schedule_error: 'planning.schedule "61 * * * *": minute 61 is out of range'
				})
			);
			await show();
			expect(q('planning-replan')!.querySelector('[role="alert"]')!.textContent).toContain(
				'is not a replan policy'
			);
			expect(text('planning-schedule')).not.toContain('next run');
			expect(q('planning-schedule')!.querySelector('[role="alert"]')!.textContent).toContain(
				'minute 61 is out of range'
			);
		});

		it('leaves the section out when flai does not say', async () => {
			backend(view(false, false));
			await show();
			expect(q('section-planning')).toBeNull();
		});
	});
});
