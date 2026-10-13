import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import EditorDialog from './EditorDialog.svelte';
import { closeEditor, editorState, openEditor } from '$lib/editor.svelte';

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

const DOC = '# Guide\n\nSome text.\n';
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

function backend() {
	api.mockImplementation(async (url: string, init?: { method?: string }) => {
		if (init?.method === 'PUT') return json({ hash: 'h2', committed: true, commit: 'abc1234' });
		if (url.startsWith('/api/docs/lint')) return json({ findings: [] });
		if (url.startsWith('/api/docs/edit'))
			return json({ path: 'docs/users/guide.md', content: DOC, hash: 'h1', mode: 'full' });
		return json([]);
	});
}
const dialog = () => document.querySelector('[role=dialog]');
const bodyBox = () =>
	document.querySelector('textarea[aria-label="Body (markdown)"]') as HTMLTextAreaElement;
const button = (text: RegExp) =>
	[...document.querySelectorAll('button')].find((b) => text.test(b.textContent ?? ''))!;
function type(el: HTMLTextAreaElement, value: string) {
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
}

describe('EditorDialog', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		closeEditor();
		flushSync();
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		vi.restoreAllMocks();
		document.body.innerHTML = '';
	});

	it('opens the document openEditor names, at its line and on its story, over the page', async () => {
		backend();
		c = mount(EditorDialog, { target: document.body });
		flushSync();
		expect(dialog()).toBeNull();
		openEditor({ path: 'docs/users/guide.md', line: 3, story: 'S-0378' });
		await settle();
		expect(dialog()?.getAttribute('aria-label')).toBe('Edit docs/users/guide.md');
		expect(api).toHaveBeenCalledWith('/api/docs/edit?path=docs%2Fusers%2Fguide.md&story=S-0378');
		expect(dialog()!.textContent).toContain('story/S-0378');
		expect(bodyBox().selectionStart).toBe(DOC.indexOf('Some text.'));
		button(/^Close$/).click();
		flushSync();
		expect(dialog()).toBeNull();
		expect(editorState.target).toBeNull();
	});

	it('asks before closing with unsaved changes, and stays open when told no', async () => {
		backend();
		const ask = vi.spyOn(window, 'confirm').mockReturnValue(false);
		c = mount(EditorDialog, { target: document.body });
		openEditor({ path: 'docs/users/guide.md' });
		await settle();
		type(bodyBox(), DOC + 'more\n');
		button(/^Close$/).click();
		flushSync();
		expect(ask).toHaveBeenCalledWith('Close without saving your changes?');
		expect(dialog()).not.toBeNull();
		window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
		flushSync();
		expect(ask).toHaveBeenCalledTimes(2);
		expect(dialog()).not.toBeNull();
		ask.mockReturnValue(true);
		button(/^Close$/).click();
		flushSync();
		expect(dialog()).toBeNull();
	});

	it('closes without asking when nothing changed', async () => {
		backend();
		const ask = vi.spyOn(window, 'confirm');
		c = mount(EditorDialog, { target: document.body });
		openEditor({ path: 'docs/users/guide.md' });
		await settle();
		window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
		flushSync();
		expect(ask).not.toHaveBeenCalled();
		expect(dialog()).toBeNull();
	});

	it('asks before a link replaces a document with unsaved changes, and keeps it when told no', async () => {
		backend();
		const ask = vi.spyOn(window, 'confirm').mockReturnValue(false);
		c = mount(EditorDialog, { target: document.body });
		openEditor({ path: 'docs/users/guide.md' });
		await settle();
		type(bodyBox(), DOC + 'more\n');
		openEditor({ path: 'docs/users/other.md', line: 2 });
		flushSync();
		expect(ask).toHaveBeenCalledWith('Open another document without saving your changes?');
		expect(editorState.target?.path).toBe('docs/users/guide.md');
		ask.mockReturnValue(true);
		openEditor({ path: 'docs/users/other.md', line: 2 });
		flushSync();
		expect(editorState.target?.path).toBe('docs/users/other.md');
		expect(editorState.dirty).toBe(false);
	});

	it('calls onsaved after flai accepted the save', async () => {
		backend();
		const onsaved = vi.fn();
		c = mount(EditorDialog, { target: document.body });
		openEditor({ path: 'docs/users/guide.md', story: 'S-0378', onsaved });
		await settle();
		type(bodyBox(), DOC + 'fixed\n');
		expect(onsaved).not.toHaveBeenCalled();
		button(/^Save$/).click();
		await settle();
		expect(onsaved).toHaveBeenCalledWith({
			path: 'docs/users/guide.md',
			hash: 'h2',
			commit: 'abc1234',
			committed: true
		});
		expect(dialog()).not.toBeNull();
	});
});
