import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import NewItemForm from './NewItemForm.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, p?: { id: string }) => route.replace('[id]', p?.id ?? '')
}));

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const ok = (body: unknown) => ({ ok: true, statusText: 'OK', json: async () => body });
const STORY_BODY = '## Goal\n\n## Acceptance criteria\n- [ ]\n\n## Tasks\n\n## Notes\n';
const EPIC_BODY = '## Outcome\n\n## Stories\n\n## Notes\n';

function backend(
	create: (body: Record<string, unknown>) => unknown = () => ok({ item: { id: 'S-0099' } }),
	move: (id: string, to: string) => unknown = (id, to) => ok({ id, status: to })
) {
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		const moving = url.match(/^\/api\/items\/([^/]+)\/move$/);
		if (moving) return move(moving[1], JSON.parse(init!.body!).to);
		if (url.startsWith('/api/items/template'))
			return ok({ body: url.endsWith('epic') ? EPIC_BODY : STORY_BODY });
		if (url.startsWith('/api/items?'))
			return ok([
				{ id: 'E-0001', title: 'Standard', status: 'done', archived: false },
				{ id: 'E-0006', title: 'Workbench', status: 'ready', archived: false },
				{ id: 'E-0004', title: 'Documentation', status: 'backlog', archived: false }
			]);
		if (url === '/api/items' && init?.method === 'POST') return create(JSON.parse(init.body!));
		throw new Error('unexpected ' + url);
	});
}
const q = <T extends HTMLElement>(id: string) =>
	document.querySelector<T>(`[data-testid="${id}"]`)!;
const type = (el: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement, value: string) => {
	el.value = value;
	el.dispatchEvent(
		new Event(el instanceof HTMLSelectElement ? 'change' : 'input', { bubbles: true })
	);
	flushSync();
};
const button = () => document.querySelector<HTMLButtonElement>('form button')!;

describe('NewItemForm', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('starts from the template’s sections, shows no front matter, and offers open epics only', async () => {
		backend();
		c = mount(NewItemForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		expect(q<HTMLTextAreaElement>('body').value).toBe(STORY_BODY);
		expect(document.body.textContent).not.toContain('---');
		expect(document.querySelectorAll('input[name], textarea[name]')).toHaveLength(0);
		const options = [...q<HTMLSelectElement>('parent').options].map((o) => o.value);
		expect(options).toEqual(['', 'E-0006', 'E-0004']);
		expect(q('preview').innerHTML).toContain('<h2');
	});

	it('says what each nature means, release consequence included', async () => {
		backend();
		c = mount(NewItemForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		expect(q('meaning').textContent).toContain('feature: New capability');
		type(q<HTMLSelectElement>('nature'), 'research');
		expect(q('meaning').textContent).toContain('no release');
		expect([...q<HTMLSelectElement>('nature').options].map((o) => o.value)).toEqual([
			'feature',
			'improvement',
			'remediation',
			'research',
			'experiment'
		]);
	});

	it('needs only a title and a body; a story needs no epic (S-0092), an epic needs no parent and gets its own sections', async () => {
		backend();
		c = mount(NewItemForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		expect(button().disabled).toBe(true);
		type(q<HTMLInputElement>('title'), 'Created from the board');
		expect(button().disabled).toBe(false); // no epic chosen, and none is needed
		expect(q<HTMLSelectElement>('parent').value).toBe('');
		type(q<HTMLSelectElement>('parent'), 'E-0006');
		expect(button().disabled).toBe(false);

		document.querySelector<HTMLInputElement>('input[type=radio][value=epic]')!.click();
		await settle();
		expect(document.querySelector('[data-testid="parent"]')).toBeNull();
		expect(q<HTMLTextAreaElement>('body').value).toBe(EPIC_BODY);
		expect(button().textContent).toContain('Create epic');
		expect(button().disabled).toBe(false);
	});

	it('creates a story with no epic when "No epic" stays chosen', async () => {
		let sent: Record<string, unknown> | undefined;
		backend((body) => {
			sent = body;
			return ok({ item: { id: 'S-0099' } });
		});
		const oncreated = vi.fn();
		c = mount(NewItemForm, { target: document.body, props: { oncreated } });
		await settle();
		type(q<HTMLInputElement>('title'), 'Standalone');
		expect(q<HTMLSelectElement>('parent').value).toBe('');
		q<HTMLFormElement>('new-item').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
		expect(sent?.parent).toBe('');
		expect(oncreated).toHaveBeenCalledWith('S-0099');
	});

	it("starts as the item page's New link asks: its type, and a story under its epic (S-0171)", async () => {
		let sent: Record<string, unknown> | undefined;
		backend((body) => {
			sent = body;
			return ok({ item: { id: 'S-0099' } });
		});
		c = mount(NewItemForm, {
			target: document.body,
			props: { oncreated: vi.fn(), initialType: 'story', initialParent: 'E-0004' }
		});
		await settle();
		expect(q<HTMLSelectElement>('parent').value).toBe('E-0004');
		type(q<HTMLInputElement>('title'), 'Another under Documentation');
		q<HTMLFormElement>('new-item').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
		expect(sent?.type).toBe('story');
		expect(sent?.parent).toBe('E-0004');

		unmount(c);
		document.body.innerHTML = '';
		c = mount(NewItemForm, {
			target: document.body,
			props: { oncreated: vi.fn(), initialType: 'epic' }
		});
		await settle();
		expect(document.querySelector<HTMLInputElement>('input[value=epic]')!.checked).toBe(true);
		expect(q<HTMLTextAreaElement>('body').value).toBe(EPIC_BODY);
		expect(button().textContent).toContain('Create epic');
	});

	it('does not start under an epic that has closed since (S-0171)', async () => {
		backend();
		c = mount(NewItemForm, {
			target: document.body,
			props: { oncreated: vi.fn(), initialParent: 'E-0001' }
		});
		await settle();
		expect(q<HTMLSelectElement>('parent').value).toBe('');
	});

	it('does not replace text the designer has written when the type changes', async () => {
		backend();
		c = mount(NewItemForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		type(q<HTMLTextAreaElement>('body'), '## Goal\nMine.\n');
		document.querySelector<HTMLInputElement>('input[type=radio][value=epic]')!.click();
		await settle();
		expect(q<HTMLTextAreaElement>('body').value).toBe('## Goal\nMine.\n');
	});

	it('sends the choices and the markdown, and hands the new ID on', async () => {
		let sent: Record<string, unknown> | undefined;
		backend((body) => {
			sent = body;
			return ok({ item: { id: 'S-0099' } });
		});
		const oncreated = vi.fn();
		c = mount(NewItemForm, { target: document.body, props: { oncreated } });
		await settle();
		type(q<HTMLInputElement>('title'), 'Created from the board');
		type(q<HTMLSelectElement>('parent'), 'E-0006');
		type(q<HTMLSelectElement>('nature'), 'improvement');
		type(q<HTMLInputElement>('tags'), 'dashboard, cli');
		type(q<HTMLInputElement>('touches'), 'flaiover/src');
		type(q<HTMLInputElement>('topics'), 'logging release');
		type(
			q<HTMLTextAreaElement>('body'),
			'## Goal\nG\n\n## Acceptance criteria\n- [ ] one\n\n## Tasks\n\n## Notes\n'
		);
		q<HTMLFormElement>('new-item').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
		expect(sent).toEqual({
			type: 'story',
			title: 'Created from the board',
			nature: 'improvement',
			parent: 'E-0006',
			tags: ['dashboard', 'cli'],
			topics: ['logging', 'release'],
			touches: ['flaiover/src'],
			body: '## Goal\nG\n\n## Acceptance criteria\n- [ ] one\n\n## Tasks\n\n## Notes\n'
		});
		expect(oncreated).toHaveBeenCalledWith('S-0099');
	});

	it('shows a refusal with its findings and keeps the text', async () => {
		backend(() => ({
			ok: false,
			statusText: 'Unprocessable',
			json: async () => ({
				error: 'flai check has 1 finding(s) with this story; nothing was created',
				findings: [
					{
						level: 'warning',
						rule: 'item.heading',
						path: 'wip/x.md',
						line: 1,
						message: 'body is missing the "## Notes" section'
					},
					{
						level: 'warning',
						rule: 'markdown.MD007',
						path: 'wip/x.md',
						line: 23,
						message: 'MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]'
					}
				]
			})
		}));
		const oncreated = vi.fn();
		c = mount(NewItemForm, { target: document.body, props: { oncreated } });
		await settle();
		type(q<HTMLInputElement>('title'), 'Refused');
		type(q<HTMLSelectElement>('parent'), 'E-0006');
		type(q<HTMLTextAreaElement>('body'), '## Goal\nNo notes.\n');
		q<HTMLFormElement>('new-item').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
		const refusal = q('refusal').textContent!;
		expect(refusal).toContain('nothing was created');
		expect(refusal).toContain('item.heading');
		expect(refusal).toContain('missing the "## Notes" section');
		// each finding as the CLI prints it, so a lint finding names its line (S-0240)
		expect(refusal).toContain(
			'wip/x.md:23: warning: markdown.MD007: MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]'
		);
		expect(q<HTMLTextAreaElement>('body').value).toBe('## Goal\nNo notes.\n');
		expect(q<HTMLInputElement>('title').value).toBe('Refused');
		expect(oncreated).not.toHaveBeenCalled();
	});

	// S-0103: a story's agent over the project's default, of which only what is typed is sent
	it("shows the project's default agent and sends only the story's overrides", async () => {
		let sent: Record<string, unknown> | undefined;
		backend((body) => {
			sent = body;
			return ok({ item: { id: 'S-0100' } });
		});
		const inner = api.getMockImplementation()!;
		api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) =>
			url === '/api/manifest'
				? ok({
						agent: { harness: 'claude-code', model: 'claude-opus-5-5', config: { effort: 'high' } }
					})
				: inner(url, init)
		);
		c = mount(NewItemForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		expect(q('agent-default').textContent).toContain(
			'project default: claude-code, claude-opus-5-5, effort=high'
		);
		expect(q<HTMLInputElement>('agent-model').placeholder).toBe('claude-opus-5-5');
		type(q<HTMLInputElement>('title'), 'With its own model');
		type(q<HTMLTextAreaElement>('body'), '## Goal\nG\n');
		type(q<HTMLInputElement>('agent-model'), 'claude-sonnet-5');
		type(q<HTMLTextAreaElement>('agent-config'), 'max_turns=20');
		q<HTMLFormElement>('new-item').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
		expect(sent?.agent).toEqual({ model: 'claude-sonnet-5', config: { max_turns: '20' } });

		// a config line that is not key=value is said, and nothing is sent
		sent = undefined;
		type(q<HTMLTextAreaElement>('agent-config'), 'nonsense');
		q<HTMLFormElement>('new-item').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
		expect(sent).toBeUndefined();
		expect(q('refusal').textContent).toContain('agent config: "nonsense" is not key=value');
	});

	// S-0167: the board's lane menu opens the form on a lane, and the item is moved there once made
	describe('the lane it starts in', () => {
		const submit = async (lane?: 'backlog' | 'ready' | 'in-progress') => {
			const moves: string[] = [];
			backend(undefined, (id, to) => {
				moves.push(`${id}→${to}`);
				return ok({ id, status: to });
			});
			const oncreated = vi.fn();
			c = mount(NewItemForm, {
				target: document.body,
				props: lane ? { oncreated, initialLane: lane } : { oncreated }
			});
			await settle();
			type(q<HTMLInputElement>('title'), 'From a lane');
			q<HTMLFormElement>('new-item').dispatchEvent(
				new Event('submit', { bubbles: true, cancelable: true })
			);
			await settle();
			return { moves, oncreated };
		};

		it('is backlog unless given, and backlog needs no move', async () => {
			const { moves, oncreated } = await submit();
			expect(q<HTMLSelectElement>('lane').value).toBe('backlog');
			expect([...q<HTMLSelectElement>('lane').options].map((o) => o.value)).toEqual([
				'backlog',
				'ready',
				'in-progress'
			]);
			expect(moves).toEqual([]);
			expect(oncreated).toHaveBeenCalledWith('S-0099');
		});

		it('moves a new item to ready, and through ready to in-progress', async () => {
			let run = await submit('ready');
			expect(run.moves).toEqual(['S-0099→ready']);
			expect(run.oncreated).toHaveBeenCalledWith('S-0099');
			unmount(c!);
			document.body.innerHTML = '';
			run = await submit('in-progress');
			expect(run.moves).toEqual(['S-0099→ready', 'S-0099→in-progress']);
			expect(run.oncreated).toHaveBeenCalledWith('S-0099');
		});

		it('says which move was refused, links the item made, and does not make it again', async () => {
			backend(undefined, () => ({
				ok: false,
				statusText: 'Unprocessable',
				json: async () => ({ error: 'rule: a story needs an acceptance criterion' })
			}));
			const oncreated = vi.fn();
			c = mount(NewItemForm, {
				target: document.body,
				props: { oncreated, initialLane: 'in-progress' }
			});
			await settle();
			type(q<HTMLInputElement>('title'), 'Not ready');
			q<HTMLFormElement>('new-item').dispatchEvent(
				new Event('submit', { bubbles: true, cancelable: true })
			);
			await settle();
			expect(q('refusal').textContent).toContain(
				'S-0099 was created, but moving it to ready was refused: rule: a story needs an acceptance criterion'
			);
			expect(q<HTMLAnchorElement>('created').getAttribute('href')).toBe('/items/S-0099');
			expect(oncreated).not.toHaveBeenCalled();
			expect(button().disabled).toBe(true);
		});
	});
});
