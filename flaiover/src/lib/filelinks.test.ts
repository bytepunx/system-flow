import { beforeEach, describe, expect, it, vi } from 'vitest';

const { openEditor } = vi.hoisted(() => ({ openEditor: vi.fn() }));
vi.mock('./editor.svelte', () => ({ openEditor }));

import { editable, editLabel, fileRef, linkSegments, parseFileRefs } from './filelinks';

describe('parseFileRefs', () => {
	it('finds a path, a path with a line, and a path with a line and a column', () => {
		const text = 'see design/a.md, docs/b.md:12 and wip/agents/S-0001.md:3:7: MD012 blank.';
		expect(parseFileRefs(text)).toEqual([
			{ text: 'see ' },
			{ text: 'design/a.md', ref: { path: 'design/a.md' } },
			{ text: ', ' },
			{ text: 'docs/b.md:12', ref: { path: 'docs/b.md', line: 12 } },
			{ text: ' and ' },
			{
				text: 'wip/agents/S-0001.md:3:7',
				ref: { path: 'wip/agents/S-0001.md', line: 3, column: 7 }
			},
			{ text: ': MD012 blank.' }
		]);
	});

	it('leaves paths outside the folders, inside longer paths, and other files as text', () => {
		for (const text of [
			'flaiover/README.md:4',
			'template/root/design/system/a.md',
			'/home/x/design/a.md',
			'design/a.mdx',
			'design/a.md.bak',
			'designs/a.md',
			'design/../flai/x.md',
			'design/system/a.go:12'
		])
			expect(parseFileRefs(text), text).toEqual([{ text }]);
	});

	it('takes the folders from the layout it is given', () => {
		const layout = { design: 'spec/', docs: 'manual', wip: 'work' };
		expect(fileRef('spec/a.md:2', layout)).toEqual({ path: 'spec/a.md', line: 2 });
		expect(fileRef('design/a.md', layout)).toBeNull();
	});

	it('reads a line of 0 as no line', () => {
		expect(fileRef('design/a.md:0')).toBeNull();
		expect(parseFileRefs('design/a.md:0')[0]).toEqual({
			text: 'design/a.md',
			ref: { path: 'design/a.md' }
		});
	});
});

describe('fileRef', () => {
	it('names the file only when the whole text is one reference', () => {
		expect(fileRef(' design/a.md:9 ')).toEqual({ path: 'design/a.md', line: 9 });
		expect(fileRef('see design/a.md')).toBeNull();
		expect(fileRef('')).toBeNull();
	});
});

describe('editable', () => {
	it('leaves out the files ADR-0023 gives mode none by their path', () => {
		expect(editable('design/system/a.md')).toBe(true);
		expect(editable('wip/kanban/stories/S-0001-x.md')).toBe(true);
		expect(editable('wip/agents/S-0001.md')).toBe(true);
		expect(editable('wip/archive/stories/S-0001-x.md')).toBe(false);
		expect(editable('wip/threads/TH-0001.md')).toBe(false);
		expect(editable('design/issues/I-0001-x.md')).toBe(false);
		expect(editable('design/issues/summary.md')).toBe(false);
		expect(editable('wip/agents/index.md')).toBe(false);
	});
});

describe('linkSegments', () => {
	beforeEach(() => openEditor.mockClear());

	it('opens an editable reference at its line on the story', () => {
		const [ref, sep, threadRef] = linkSegments('design/a.md:12 wip/threads/TH-0001.md', {
			story: 'S-0378'
		});
		expect(ref.text).toBe('design/a.md:12');
		expect(ref.edit?.label).toBe('Edit design/a.md at line 12');
		ref.edit!.open();
		expect(openEditor).toHaveBeenCalledWith({ path: 'design/a.md', line: 12, story: 'S-0378' });
		expect(sep).toEqual({ text: ' ' });
		expect(threadRef).toEqual({ text: 'wip/threads/TH-0001.md' });
	});

	it("opens the files of a refusal for the main checkout's uncommitted files there, not on the story", () => {
		const text =
			"S-0378 cannot be accepted: the main checkout's uncommitted files hold lint the fix of the markdown tier could not fix:\n  markdown: docs/users/flai.md:12: MD009";
		const ref = linkSegments(text, { story: 'S-0378' }).find((s) => s.edit)!;
		ref.edit!.open();
		expect(openEditor).toHaveBeenCalledWith({
			path: 'docs/users/flai.md',
			line: 12,
			story: undefined
		});
	});

	it('opens a wip file on the main checkout, never on the story', () => {
		const [ref] = linkSegments('wip/kanban/stories/S-0378-x.md:3', { story: 'S-0378' });
		ref.edit!.open();
		expect(openEditor).toHaveBeenCalledWith({
			path: 'wip/kanban/stories/S-0378-x.md',
			line: 3,
			story: undefined
		});
	});

	it('labels a reference without a line by its path', () => {
		expect(editLabel({ path: 'docs/a.md' })).toBe('Edit docs/a.md');
	});
});
