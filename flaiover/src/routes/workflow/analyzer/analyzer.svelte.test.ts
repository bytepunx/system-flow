// S-0229: the analyzer's page shows the analysis block of settings.get in its settings panel, and
// saves a key through /api/settings and reads the settings again.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { SettingsView } from '$lib/settings';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$lib/events', () => ({ follow: () => () => {} }));

import AnalyzerPage from './+page.svelte';

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
function type(el: HTMLInputElement, value: string) {
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
}

/** A settings.get answer with the analysis block as flai lists it, and a key of another block. */
function settingsGet(): SettingsView {
	return {
		here: true,
		everywhere: false,
		enable: 'flai serve enable settings',
		enable_everywhere: 'flai serve enable settings --everywhere',
		host: {
			actions: [],
			agent: { name: 'agent', command: [], harnesses: {} },
			checks: { commands: [], timeout_minutes: 10 },
			import_roots: [],
			strategic: {
				editable: true,
				enable: 'flai serve enable settings',
				settings: [
					{
						key: 'orchestration.policy',
						kind: 'choice',
						values: ['cod', 'wsjf', 'throughput', 'fifo'],
						default: 'fifo',
						meaning: 'How the ready column is ordered.',
						set: false
					},
					{
						key: 'analysis.agent',
						kind: 'agent',
						meaning: "The analyzer's agent, over the project's.",
						value: { model: 'claude-sonnet-5' },
						set: true
					},
					{
						key: 'analysis.schedule',
						kind: 'cron',
						meaning: 'When flai serve runs the analyzer.',
						set: false
					}
				]
			}
		}
	};
}
const reads = () =>
	api.mock.calls.filter(([url, init]) => url === '/api/settings' && !init?.method).length;

describe('the analyzer page (S-0229)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});
	const show = async () => {
		c = mount(AnalyzerPage, { target: document.body });
		await settle();
	};

	it("says its status comes later and shows the analysis block's keys", async () => {
		api.mockResolvedValue(answer(200, settingsGet()));
		await show();
		expect(document.querySelector('h1')!.textContent).toBe('Analyzer');
		expect(document.body.textContent).toContain(
			'Its status, activity log, and runs will be shown here.'
		);
		const rows = [...document.querySelectorAll('[data-testid^="row-"]')].map((r) =>
			r.getAttribute('data-testid')!.slice('row-'.length)
		);
		expect(rows).toEqual(['analysis.agent', 'analysis.schedule']);
		expect(q<HTMLInputElement>('input-analysis.agent-model')!.value).toBe('claude-sonnet-5');
	});

	it('saves a changed key, then reads the settings again', async () => {
		const sent: unknown[] = [];
		api.mockImplementation(async (_url: string, init?: { method?: string; body?: string }) => {
			if (init?.method === 'POST') {
				sent.push(JSON.parse(init.body!));
				return answer(200, { set: [], unset: [], warnings: [] });
			}
			return answer(200, settingsGet());
		});
		await show();
		type(q<HTMLInputElement>('input-analysis.schedule')!, 'daily');
		q<HTMLButtonElement>('strategic-save')!.click();
		await settle();
		expect(sent).toEqual([{ kind: 'manifest', set: { 'analysis.schedule': 'daily' } }]);
		expect(reads()).toBe(2);
	});
});
