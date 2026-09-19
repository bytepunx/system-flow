import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import NewAdrForm from './NewAdrForm.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const ok = (body: unknown) => ({ ok: true, statusText: 'OK', json: async () => body });
const TEMPLATE = '## Context\n\n## Decision\n\n## Consequences\n\n## Alternatives considered\n';

function backend(
	create: (body: Record<string, unknown>) => unknown = () =>
		ok({ id: 'ADR-0027', path: 'design/adrs/0027-x.md' })
) {
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		if (url === '/api/adrs/template') return ok({ body: TEMPLATE });
		if (url === '/api/docs/adrs')
			return ok([
				{ id: 'ADR-0000', title: 'Template', status: 'template' },
				{ id: 'ADR-0007', title: 'SvelteKit SPA', status: 'accepted' },
				{ id: 'ADR-0018', title: 'Dashboard token', status: 'accepted' }
			]);
		if (url === '/api/adrs' && init?.method === 'POST') return create(JSON.parse(init.body!));
		throw new Error('unexpected ' + url);
	});
}
const q = <T extends HTMLElement>(id: string) =>
	document.querySelector<T>(`[data-testid="${id}"]`)!;
const type = (el: HTMLInputElement | HTMLTextAreaElement, value: string) => {
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
};
const box = (group: string, id: string) =>
	q(group).querySelector<HTMLInputElement>(`input[value="${id}"]`)!;
const button = () => document.querySelector<HTMLButtonElement>('form button')!;
const submit = async () => {
	q<HTMLFormElement>('new-adr').dispatchEvent(
		new Event('submit', { bubbles: true, cancelable: true })
	);
	await settle();
};

describe('NewAdrForm', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('starts from the template’s sections, proposed, with no front matter and no template to pick', async () => {
		backend();
		c = mount(NewAdrForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		expect(q<HTMLTextAreaElement>('body').value).toBe(TEMPLATE);
		expect(document.body.textContent).not.toContain('---');
		expect(button().textContent).toContain('Record as proposed');
		expect(button().disabled).toBe(true); // no title yet
		const ids = [...q('supersedes').querySelectorAll('input')].map(
			(i) => (i as HTMLInputElement).value
		);
		expect(ids).toEqual(['ADR-0007', 'ADR-0018']);
		expect(q('preview').innerHTML).toContain('<h2');
	});

	it('sends the title, the status, the chain, and the markdown, and hands the new path on', async () => {
		let sent: Record<string, unknown> | undefined;
		backend((body) => {
			sent = body;
			return ok({ id: 'ADR-0027', path: 'design/adrs/0027-one-token-per-project.md' });
		});
		const oncreated = vi.fn();
		c = mount(NewAdrForm, { target: document.body, props: { oncreated } });
		await settle();
		type(q<HTMLInputElement>('title'), 'One token per project');
		document.querySelector<HTMLInputElement>('input[type=radio][value=accepted]')!.click();
		box('supersedes', 'ADR-0007').click();
		box('refines', 'ADR-0018').click();
		type(q<HTMLTextAreaElement>('body'), '## Context\nc\n\n## Decision\nd\n');
		flushSync();
		expect(button().textContent).toContain('Record as accepted');
		await submit();
		expect(sent).toEqual({
			title: 'One token per project',
			status: 'accepted',
			supersedes: ['ADR-0007'],
			refines: ['ADR-0018'],
			body: '## Context\nc\n\n## Decision\nd\n'
		});
		expect(oncreated).toHaveBeenCalledWith('design/adrs/0027-one-token-per-project.md', 'ADR-0027');
	});

	it('will not record an ADR that both supersedes and refines the same decision', async () => {
		backend();
		c = mount(NewAdrForm, { target: document.body, props: { oncreated: vi.fn() } });
		await settle();
		type(q<HTMLInputElement>('title'), 'X');
		box('supersedes', 'ADR-0007').click();
		box('refines', 'ADR-0007').click();
		flushSync();
		expect(q('both').textContent).toContain('ADR-0007 cannot be both');
		expect(button().disabled).toBe(true);
	});

	it('shows a refusal with its findings and keeps the text', async () => {
		backend(() => ({
			ok: false,
			statusText: 'Unprocessable',
			json: async () => ({
				error: 'flai check has 1 finding(s) with this ADR; nothing was created',
				findings: [
					{
						level: 'warning',
						rule: 'adr.index',
						path: 'design/adrs/README.md',
						line: 9,
						message: 'the row links to a file that is not there'
					}
				]
			})
		}));
		const oncreated = vi.fn();
		c = mount(NewAdrForm, { target: document.body, props: { oncreated } });
		await settle();
		type(q<HTMLInputElement>('title'), 'Refused');
		type(q<HTMLTextAreaElement>('body'), '## Context\nmine\n');
		await submit();
		expect(q('refusal').textContent).toContain('nothing was created');
		expect(q('refusal').textContent).toContain('adr.index');
		expect(q<HTMLTextAreaElement>('body').value).toBe('## Context\nmine\n');
		expect(oncreated).not.toHaveBeenCalled();
	});
});
