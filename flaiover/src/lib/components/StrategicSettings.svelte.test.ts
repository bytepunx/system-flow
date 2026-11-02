import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import StrategicSettings from './StrategicSettings.svelte';
import type { StrategicBlock, StrategicSettings as Strategic } from '$lib/settings';

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
function type(el: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement, value: string) {
	el.value = value;
	el.dispatchEvent(
		new Event(el instanceof HTMLSelectElement ? 'change' : 'input', { bubbles: true })
	);
	flushSync();
}

function strategic(editable: boolean): Strategic {
	return {
		editable,
		enable: 'flai serve enable settings',
		settings: [
			{
				key: 'orchestration.permissions.promote_to_ready',
				kind: 'boolean',
				values: ['true', 'false'],
				default: 'false',
				meaning: 'Lets the orchestrator move a story to ready.',
				risk: 'Agents may pull, work, and spend on a story you have not chosen to start.',
				value: true,
				set: true
			},
			{
				key: 'orchestration.permissions.answer_threads',
				kind: 'choice',
				values: ['off', 'recommend', 'autonomous'],
				default: 'off',
				meaning: 'How the orchestrator answers threads.',
				risk: 'With autonomous, agents act on answers you did not give.',
				set: false
			},
			{
				key: 'orchestration.policy',
				kind: 'choice',
				values: ['cod', 'wsjf', 'throughput', 'fifo'],
				default: 'fifo',
				meaning: 'How the ready column is ordered.',
				value: 'cod',
				set: true
			},
			{
				key: 'orchestration.release.count',
				kind: 'number',
				whole: true,
				meaning: 'The number of accepted stories at which a release is due.',
				value: 3,
				set: true
			},
			{
				key: 'planning.agent',
				kind: 'agent',
				meaning: "The planner's agent, over the project's.",
				value: { model: 'claude-sonnet-5', config: { effort: 'medium' } },
				set: true
			},
			{
				key: 'planning.schedule',
				kind: 'cron',
				meaning: 'When flai serve runs the planner over the ready column.',
				value: 'daily',
				set: true
			},
			{
				key: 'planning.cycle',
				kind: 'duration',
				default: '168h',
				meaning: "The period a cost of delay's time lost is counted over.",
				set: false
			},
			{
				key: 'analysis.schedule',
				kind: 'cron',
				meaning: 'When flai serve runs the analyzer.',
				set: false
			}
		]
	};
}

/** The strategic block with a role the project defines, reviewer, as settings.get lists it (S-0420). */
function withReviewer(editable: boolean): Strategic {
	return {
		...strategic(editable),
		roles: [
			{
				name: 'reviewer',
				title: 'Reviewer',
				path: 'design/workflows/roles/reviewer.md',
				description: 'Reviews a story against the conventions, for flai serve.',
				action: 'review',
				settings: [
					{
						key: 'roles.reviewer.agent',
						kind: 'agent',
						meaning: 'The agent of the role reviewer, over the project agent.',
						set: false
					},
					{
						key: 'roles.reviewer.schedule',
						kind: 'cron',
						meaning: 'When the role reviewer runs on its schedule.',
						value: 'daily',
						set: true
					},
					{
						key: 'roles.reviewer.permissions.publish',
						kind: 'boolean',
						values: ['true', 'false'],
						default: 'false',
						meaning: 'Lets the role reviewer publish a release.',
						value: true,
						set: true,
						allows: ['release_publish', 'flai release when own-thread']
					}
				]
			}
		]
	};
}

describe('StrategicSettings (S-0229)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	function show(block: StrategicBlock, s: Strategic, onsaved?: () => void) {
		c = mount(StrategicSettings, {
			target: document.body,
			props: { block, strategic: s, onsaved }
		});
		flushSync();
	}
	function posts(status = 200, body: unknown = { set: [], unset: [], warnings: [] }) {
		const sent: Record<string, unknown>[] = [];
		api.mockImplementation(async (_url: string, init?: { body?: string }) => {
			sent.push(JSON.parse(init!.body!));
			return answer(status, body);
		});
		return sent;
	}

	it("shows only its block's keys, each with its meaning, a permission's risk, and an unset key's default", () => {
		show('orchestration', strategic(true));
		expect(q('row-orchestration.policy')).not.toBeNull();
		expect(q('row-planning.agent')).toBeNull();
		expect(q('row-orchestration.permissions.promote_to_ready')?.textContent).toContain(
			'Lets the orchestrator move a story to ready.'
		);
		expect(q('risk-orchestration.permissions.promote_to_ready')?.textContent).toContain(
			'Agents may pull'
		);
		expect(q('risk-orchestration.policy')).toBeNull();
		expect(q('default-orchestration.permissions.answer_threads')?.textContent).toContain('off');
		expect(q('default-orchestration.policy')).toBeNull();
		const toggle = q<HTMLInputElement>('input-orchestration.permissions.promote_to_ready')!;
		expect(toggle.getAttribute('role')).toBe('switch');
		expect(toggle.checked).toBe(true);
		const label = document.querySelector(`label[for="${toggle.id}"]`);
		expect(label?.textContent).toBe('promote_to_ready');
		expect(q<HTMLInputElement>('input-orchestration.release.count')!.step).toBe('1');
	});

	it('saves a changed key and posts only that key, then says so and calls onsaved', async () => {
		const sent = posts(200, { set: [], unset: [], commit: 'abcdef123456', warnings: [] });
		const onsaved = vi.fn();
		show('orchestration', strategic(true), onsaved);
		expect(q<HTMLButtonElement>('strategic-save')!.disabled).toBe(true);
		type(q<HTMLSelectElement>('input-orchestration.policy')!, 'wsjf');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(api).toHaveBeenCalledWith('/api/settings', expect.objectContaining({ method: 'POST' }));
		expect(sent).toEqual([{ kind: 'manifest', set: { 'orchestration.policy': 'wsjf' } }]);
		expect(onsaved).toHaveBeenCalledOnce();
		expect(q('strategic-said')?.textContent).toContain('saved in abcdef1');
		expect(q<HTMLButtonElement>('strategic-save')!.disabled).toBe(true);
	});

	it('posts an emptied field as an unset, a switch as a boolean, and an agent as an object', async () => {
		const sent = posts();
		show('planning', strategic(true));
		type(q<HTMLInputElement>('input-planning.schedule')!, '  ');
		type(q<HTMLInputElement>('input-planning.agent-harness')!, 'claude-code');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(sent).toEqual([
			{
				kind: 'manifest',
				set: {
					'planning.agent': {
						harness: 'claude-code',
						model: 'claude-sonnet-5',
						config: { effort: 'medium' }
					}
				},
				unset: ['planning.schedule']
			}
		]);

		unmount(c!);
		document.body.innerHTML = '';
		show('orchestration', strategic(true));
		const toggle = q<HTMLInputElement>('input-orchestration.permissions.promote_to_ready')!;
		toggle.click();
		flushSync();
		type(q<HTMLInputElement>('input-orchestration.release.count')!, '5');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(sent[1]).toEqual({
			kind: 'manifest',
			set: {
				'orchestration.permissions.promote_to_ready': false,
				'orchestration.release.count': 5
			}
		});
	});

	it("shows a refusal's reason under its field, keeps what was typed, and a field not shown at the top", async () => {
		const sent = posts(422, {
			error: 'orchestration.release.count: must be a whole number; orchestration.release: x',
			refused: [
				{ field: 'orchestration.release.count', reason: 'must be a whole number' },
				{ field: 'orchestration.release', reason: 'theme needs an epic or a tag' }
			]
		});
		const onsaved = vi.fn();
		show('orchestration', strategic(true), onsaved);
		type(q<HTMLInputElement>('input-orchestration.release.count')!, '2.5');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(sent).toEqual([{ kind: 'manifest', set: { 'orchestration.release.count': 2.5 } }]);
		const under = q('refused-orchestration.release.count')!;
		expect(under.textContent).toContain('must be a whole number');
		expect(q('row-orchestration.release.count')!.contains(under)).toBe(true);
		const input = q<HTMLInputElement>('input-orchestration.release.count')!;
		expect(input.value).toBe('2.5');
		expect(input.getAttribute('aria-invalid')).toBe('true');
		expect(input.getAttribute('aria-describedby')).toContain(under.id);
		expect(q('strategic-refused')?.textContent).toContain('theme needs an epic or a tag');
		expect(q('strategic-said')?.getAttribute('role')).toBe('alert');
		expect(onsaved).not.toHaveBeenCalled();
		expect(q<HTMLButtonElement>('strategic-save')!.disabled).toBe(false);
	});

	it('refuses a config line that is not key=value under its field without posting', async () => {
		posts();
		show('planning', strategic(true));
		type(q<HTMLTextAreaElement>('input-planning.agent-config')!, 'effort');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(api).not.toHaveBeenCalled();
		expect(q('refused-planning.agent')?.textContent).toContain('config: "effort" is not key=value');
		expect(q<HTMLTextAreaElement>('input-planning.agent-config')!.value).toBe('effort');
	});

	it('is read-only, with the reason and the command, while settings is off', () => {
		show('analysis', strategic(false));
		const note = q('strategic-readonly')!;
		expect(note.textContent).toContain('settings host action is off');
		expect(note.textContent).toContain('flai serve enable settings');
		const inputs = document.querySelectorAll<HTMLInputElement>('[data-testid^="input-"]');
		expect(inputs.length).toBeGreaterThan(0);
		for (const el of inputs) expect(el.disabled).toBe(true);
		expect(q<HTMLButtonElement>('strategic-save')!.disabled).toBe(true);
	});

	it('says the command when the host answers Disabled', async () => {
		posts(403, {
			error: 'the settings host action is off',
			action: 'settings',
			enable: 'flai serve enable settings'
		});
		show('planning', strategic(true));
		type(q<HTMLInputElement>('input-planning.cycle')!, '72h');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		const said = q('strategic-said')!;
		expect(said.textContent).toContain('the settings host action is off');
		expect(said.textContent).toContain('flai serve enable settings');
	});
});

describe('StrategicSettings for a role (S-0420)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	function show(role: string, s: Strategic, onsaved?: () => void) {
		c = mount(StrategicSettings, { target: document.body, props: { role, strategic: s, onsaved } });
		flushSync();
	}

	it("lists a defined role's agent, schedule, and permissions, each permission with what it allows", () => {
		show('reviewer', withReviewer(true));
		const panel = q('strategic-reviewer')!;
		expect(panel.querySelector('h2')?.textContent).toBe("The reviewer's settings");
		for (const key of [
			'roles.reviewer.agent',
			'roles.reviewer.schedule',
			'roles.reviewer.permissions.publish'
		])
			expect(q(`row-${key}`)).not.toBeNull();
		expect(q('row-planning.agent')).toBeNull();
		expect(q('row-orchestration.policy')).toBeNull();
		expect(q('default-roles.reviewer.agent')?.textContent).toContain("the project's agent");
		expect(q<HTMLInputElement>('input-roles.reviewer.agent-harness')).not.toBeNull();
		expect(q<HTMLInputElement>('input-roles.reviewer.schedule')!.value).toBe('daily');
		const toggle = q<HTMLInputElement>('input-roles.reviewer.permissions.publish')!;
		expect(toggle.getAttribute('role')).toBe('switch');
		expect(toggle.checked).toBe(true);
		expect(document.querySelector(`label[for="${toggle.id}"]`)?.textContent).toBe('publish');
		const allows = q('allows-roles.reviewer.permissions.publish')!;
		expect(allows.textContent).toContain('release_publish');
		expect(allows.textContent).toContain('flai release when own-thread');
		expect(toggle.getAttribute('aria-describedby')).toContain(allows.id);
		expect(q('row-roles.reviewer.permissions.publish')!.contains(allows)).toBe(true);
		expect(panel.querySelector('h3')?.textContent).toBe('permissions');
	});

	it("saves a defined role's changed keys through settings.manifest while settings is on", async () => {
		const sent: Record<string, unknown>[] = [];
		api.mockImplementation(async (_url: string, init?: { body?: string }) => {
			sent.push(JSON.parse(init!.body!));
			return answer(200, { set: [], unset: [], commit: 'abcdef123456', warnings: [] });
		});
		const onsaved = vi.fn();
		show('reviewer', withReviewer(true), onsaved);
		expect(q('strategic-readonly')).toBeNull();
		type(q<HTMLInputElement>('input-roles.reviewer.schedule')!, '0 6 * * 1');
		q<HTMLInputElement>('input-roles.reviewer.permissions.publish')!.click();
		type(q<HTMLInputElement>('input-roles.reviewer.agent-model')!, 'claude-sonnet-5');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(api).toHaveBeenCalledWith('/api/settings', expect.objectContaining({ method: 'POST' }));
		expect(sent).toEqual([
			{
				kind: 'manifest',
				set: {
					'roles.reviewer.agent': { model: 'claude-sonnet-5' },
					'roles.reviewer.schedule': '0 6 * * 1',
					'roles.reviewer.permissions.publish': false
				}
			}
		]);
		expect(onsaved).toHaveBeenCalledOnce();
		expect(q('strategic-said')?.textContent).toContain('saved in abcdef1');
	});

	it("shows a defined role's settings read-only while settings is off", () => {
		show('reviewer', withReviewer(false));
		expect(q('strategic-readonly')?.textContent).toContain('flai serve enable settings');
		const inputs = document.querySelectorAll<HTMLInputElement>('[data-testid^="input-"]');
		expect(inputs.length).toBe(5);
		for (const el of inputs) expect(el.disabled).toBe(true);
		expect(q<HTMLButtonElement>('strategic-save')!.disabled).toBe(true);
		expect(q('allows-roles.reviewer.permissions.publish')).not.toBeNull();
	});

	it("shows a built-in role's block, and a line for a role flai lists no settings for", () => {
		show('planner', withReviewer(true));
		expect(q('strategic-planning')?.querySelector('h2')?.textContent).toBe(
			"The planner's settings"
		);
		expect(q('row-planning.agent')).not.toBeNull();
		expect(q('row-roles.reviewer.agent')).toBeNull();
		unmount(c!);
		document.body.innerHTML = '';

		show('reviewer', strategic(true));
		expect(q('strategic-none')?.textContent).toContain('no settings for this role');
		expect(q('strategic-save')).toBeNull();
		expect(document.querySelectorAll('[data-testid^="input-"]').length).toBe(0);
	});
});
