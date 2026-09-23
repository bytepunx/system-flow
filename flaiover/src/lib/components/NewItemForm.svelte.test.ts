import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import NewItemForm from './NewItemForm.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const ok = (body: unknown) => ({ ok: true, statusText: 'OK', json: async () => body });
const STORY_BODY = '## Goal\n\n## Acceptance criteria\n- [ ]\n\n## Tasks\n\n## Notes\n';
const EPIC_BODY = '## Outcome\n\n## Stories\n\n## Notes\n';

function backend(
	create: (body: Record<string, unknown>) => unknown = () => ok({ item: { id: 'S-0099' } })
) {
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
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
		expect(q<HTMLTextAreaElement>('body').value).toBe('## Goal\nNo notes.\n');
		expect(q<HTMLInputElement>('title').value).toBe('Refused');
		expect(oncreated).not.toHaveBeenCalled();
	});
});
