import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import SettingsPanel from './SettingsPanel.svelte';
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
				attended_minutes: 0,
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
});
