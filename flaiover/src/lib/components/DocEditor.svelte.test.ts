import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import DocEditor from './DocEditor.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[...path]', params.path ?? '').replace('[id]', params.id ?? '')
}));
vi.mock('$lib/markdown', () => ({
	render: (md: string) => `<p>${md.length} chars</p>`,
	enhance: async () => {}
}));

const DOC = '---\ntitle: Plan\nupdated: 2026-09-01\n---\n\n# Plan\n\n## Shape\ntext\n';
const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const json = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: String(status),
	json: async () => body
});

type Reply = ReturnType<typeof json>;
/** Route the mocked api by URL and method; `put` answers the save. */
function backend(doc: Record<string, unknown>, put?: () => Reply, items: unknown[] = []) {
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		if (init?.method === 'PUT') return put ? put() : json({});
		if (url.startsWith('/api/docs/edit')) return json(doc);
		if (url.startsWith('/api/items')) return json(items);
		if (url.startsWith('/api/threads')) return json([]);
		return json({});
	});
}
const bodyBox = () =>
	document.querySelector('textarea[aria-label="Body (markdown)"]') as HTMLTextAreaElement;
const saveButton = () =>
	[...document.querySelectorAll('button')].find((b) => /^Sav/.test(b.textContent ?? ''))!;
function type(el: HTMLTextAreaElement, value: string) {
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
}
const sent = () =>
	JSON.parse(
		api.mock.calls.filter((c) => (c[1] as { method?: string })?.method === 'PUT').at(-1)![1].body
	);

describe('DocEditor', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('edits body and front matter of a design document and saves with the loaded hash', async () => {
		backend({ path: 'design/system/plan.md', content: DOC, hash: 'h1', mode: 'full' }, () =>
			json({ hash: 'h2', committed: true, commit: 'abc1234' })
		);
		const ondirty = vi.fn();
		c = mount(DocEditor, {
			target: document.body,
			props: { path: 'design/system/plan.md', ondirty }
		});
		await settle();
		expect(document.querySelector('textarea[aria-label="Front matter (YAML)"]')).not.toBeNull();
		expect(saveButton().disabled).toBe(true); // nothing to save yet
		type(bodyBox(), bodyBox().value.replace('text', 'better text'));
		expect(ondirty).toHaveBeenLastCalledWith(true);
		expect(saveButton().disabled).toBe(false);
		saveButton().click();
		await settle();
		expect(sent()).toMatchObject({
			path: 'design/system/plan.md',
			hash: 'h1',
			content: DOC.replace('text', 'better text')
		});
		expect(document.body.textContent).toContain('committed as abc1234');
	});

	it('shows flai-owned front matter read-only and sends it back unchanged', async () => {
		const story = '---\nid: S-0001\nstatus: backlog\n---\n\n# S-0001\n\n## Goal\n';
		backend(
			{
				path: 'wip/kanban/stories/S-0001-x.md',
				content: story,
				hash: 'h1',
				mode: 'body',
				reason: 'flai owns a work item’s front matter'
			},
			() => json({ hash: 'h2', committed: false })
		);
		c = mount(DocEditor, {
			target: document.body,
			props: { path: 'wip/kanban/stories/S-0001-x.md' }
		});
		await settle();
		expect(document.querySelector('textarea[aria-label="Front matter (YAML)"]')).toBeNull();
		expect(document.querySelector('[data-testid="front-matter-readonly"]')?.textContent).toContain(
			'status: backlog'
		);
		expect(document.body.textContent).toContain('flai owns');
		type(bodyBox(), bodyBox().value + 'Edit bodies.\n');
		saveButton().click();
		await settle();
		expect(sent().content).toBe(story + 'Edit bodies.\n');
		expect(document.body.textContent).toContain('Not committed');
	});

	it('offers no editor for a document flai will not let be edited', async () => {
		backend({
			path: 'wip/agents/index.md',
			content: '# Index\n',
			hash: 'h',
			mode: 'none',
			reason: 'index.md is generated from the narratives'
		});
		c = mount(DocEditor, { target: document.body, props: { path: 'wip/agents/index.md' } });
		await settle();
		expect(document.querySelector('textarea')).toBeNull();
		expect(document.body.textContent).toContain('generated from the narratives');
	});

	it('lists the findings when flai refuses, and keeps the text', async () => {
		backend({ path: 'design/system/plan.md', content: DOC, hash: 'h1', mode: 'full' }, () =>
			json(
				{
					error: 'flai check has 1 finding(s) with this content; nothing was saved',
					findings: [
						{
							level: 'warning',
							rule: 'doc.title',
							path: 'design/system/plan.md',
							line: 1,
							message: 'front matter needs a title'
						}
					]
				},
				422
			)
		);
		c = mount(DocEditor, { target: document.body, props: { path: 'design/system/plan.md' } });
		await settle();
		type(bodyBox(), bodyBox().value + 'more\n');
		saveButton().click();
		await settle();
		const alert = document.querySelector('[role=alert]')!.textContent ?? '';
		expect(alert).toContain('nothing was saved');
		expect(alert).toContain('doc.title: front matter needs a title');
		expect(bodyBox().value).toContain('more');
	});

	it('shows a conflict with its diff, and saves over it only with the current hash', async () => {
		let saves = 0;
		backend({ path: 'design/system/plan.md', content: DOC, hash: 'h1', mode: 'full' }, () =>
			++saves === 1
				? json(
						{
							error: 'design/system/plan.md changed after it was loaded',
							current: DOC.replace('text', 'their text'),
							hash: 'h-theirs',
							diff: '-their text\n+my text'
						},
						409
					)
				: json({ hash: 'h3', committed: true, commit: 'def5678' })
		);
		c = mount(DocEditor, { target: document.body, props: { path: 'design/system/plan.md' } });
		await settle();
		type(bodyBox(), bodyBox().value.replace('text', 'my text'));
		saveButton().click();
		await settle();
		expect(document.querySelector('[role=alert]')!.textContent).toContain('+my text');
		const over = [...document.querySelectorAll('button')].find((b) =>
			b.textContent?.includes('Save mine over it')
		)!;
		over.click();
		await settle();
		expect(sent().hash).toBe('h-theirs');
		expect(document.body.textContent).toContain('committed as def5678');
	});

	it('loads the current version on a conflict when asked', async () => {
		backend({ path: 'design/system/plan.md', content: DOC, hash: 'h1', mode: 'full' }, () =>
			json({ error: 'changed', current: DOC.replace('text', 'their text'), hash: 'h-theirs' }, 409)
		);
		c = mount(DocEditor, { target: document.body, props: { path: 'design/system/plan.md' } });
		await settle();
		type(bodyBox(), bodyBox().value.replace('text', 'my text'));
		saveButton().click();
		await settle();
		[...document.querySelectorAll('button')]
			.find((b) => b.textContent?.includes('Load the current version'))!
			.click();
		await settle();
		expect(bodyBox().value).toContain('their text');
		expect(saveButton().disabled).toBe(true);
	});

	it('warns when a story touches the document and waits for the acknowledgement', async () => {
		backend(
			{ path: 'design/system/plan.md', content: DOC, hash: 'h1', mode: 'full' },
			() => json({ hash: 'h2', committed: true, commit: 'abc' }),
			[{ id: 'S-0040', title: 'Editing', status: 'in-progress', touches: ['design/system'] }]
		);
		c = mount(DocEditor, { target: document.body, props: { path: 'design/system/plan.md' } });
		await settle();
		expect(document.body.textContent).toContain('Being worked on by');
		expect(document.body.textContent).toContain('S-0040');
		type(bodyBox(), bodyBox().value + 'x\n');
		expect(saveButton().disabled).toBe(true);
		(document.querySelector('input[type=checkbox]') as HTMLInputElement).click();
		flushSync();
		expect(saveButton().disabled).toBe(false);
	});

	it('opens a thread on the heading under the cursor', async () => {
		backend({ path: 'design/system/plan.md', content: DOC, hash: 'h1', mode: 'full' });
		c = mount(DocEditor, { target: document.body, props: { path: 'design/system/plan.md' } });
		await settle();
		const threadButton = () =>
			[...document.querySelectorAll('button')].find((b) =>
				b.textContent?.includes('Open a thread on')
			)!;
		expect(threadButton().disabled).toBe(true); // no caret yet
		const box = bodyBox();
		const at = box.value.indexOf('text');
		box.setSelectionRange(at, at);
		box.dispatchEvent(new Event('click', { bubbles: true }));
		flushSync();
		expect(threadButton().textContent).toContain('Shape');
		threadButton().click();
		await settle();
		// the Threads composer is open with that heading selected
		const select = document.querySelector('select') as HTMLSelectElement | null;
		expect(select?.value).toBe('Shape');
	});
});
