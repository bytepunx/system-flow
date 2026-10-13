// @vitest-environment jsdom
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { openEditor } = vi.hoisted(() => ({ openEditor: vi.fn() }));
vi.mock('./editor.svelte', () => ({ openEditor }));

import { enhance, markdownInto, render, rewriteLink, resolveRelative } from './markdown';

describe('links', () => {
	it('resolves relative paths against the document folder', () => {
		expect(resolveRelative('design/system/overview.md', 'workflow.md')).toBe(
			'design/system/workflow.md'
		);
		expect(resolveRelative('design/system/overview.md', '../adrs/0001-x.md')).toBe(
			'design/adrs/0001-x.md'
		);
		expect(resolveRelative('README.md', 'docs/users/index.md')).toBe('docs/users/index.md');
	});
	it('rewrites relative markdown links to explorer routes and leaves the rest', () => {
		expect(rewriteLink('design/system/a.md', '../adrs/b.md#ctx')).toBe(
			'/docs/design/adrs/b.md#ctx'
		);
		expect(rewriteLink('design/system/a.md', 'https://example.com/x.md')).toBe(
			'https://example.com/x.md'
		);
		expect(rewriteLink('design/system/a.md', '#local')).toBe('#local');
		expect(rewriteLink('design/system/a.md', '/absolute')).toBe('/absolute');
	});
});

describe('render', () => {
	it('turns mermaid fences into pre.mermaid and keeps other fences', () => {
		const html = render(
			'```mermaid\nflowchart LR\n  A --> B\n```\n\n```go\nx := 1\n```\n',
			'design/system/a.md'
		);
		expect(html).toContain('<pre class="mermaid">flowchart LR');
		expect(html).toContain('language-go');
	});
	it('renders task lists, heading anchors, and rewritten links', () => {
		const html = render(
			'# Title\n\n- [x] done\n- [ ] open\n\n[see](../adrs/0001-x.md)\n',
			'design/system/a.md'
		);
		expect(html).toContain('type="checkbox"');
		expect(html).toContain('id="title"');
		expect(html).toContain('href="/docs/design/adrs/0001-x.md"');
	});
});

describe('edit controls', () => {
	beforeEach(() => openEditor.mockClear());

	async function shown(markdown: string, docPath: string, story?: string) {
		const root = document.createElement('div');
		root.innerHTML = render(markdown, docPath);
		await enhance(root, false, { story });
		return root;
	}
	const control = (root: HTMLElement, label: string) =>
		root.querySelector<HTMLButtonElement>(`button[aria-label="${label}"]`);

	it('puts a control beside a relative link to a markdown file; the link still shows it', async () => {
		const root = await shown(
			'See [the ADR](../adrs/0001-x.md#context).\n',
			'design/system/a.md',
			'S-0378'
		);
		const link = root.querySelector('a')!;
		expect(link.getAttribute('href')).toBe('/docs/design/adrs/0001-x.md#context');
		const button = control(root, 'Edit design/adrs/0001-x.md')!;
		expect(link.nextSibling).toBe(button);
		button.click();
		expect(openEditor).toHaveBeenCalledWith({
			path: 'design/adrs/0001-x.md',
			line: undefined,
			story: 'S-0378'
		});
	});

	it('puts a control beside a code span naming a file and line, opening it there', async () => {
		const root = await shown('Fix `docs/users/b.md:12:3` first.\n', 'design/system/a.md');
		const button = control(root, 'Edit docs/users/b.md at line 12')!;
		expect(root.querySelector('code')!.nextSibling).toBe(button);
		button.click();
		expect(openEditor).toHaveBeenCalledWith({
			path: 'docs/users/b.md',
			line: 12,
			story: undefined
		});
	});

	it('gives none to files the editor cannot edit, other links, code blocks, or prose', async () => {
		const root = await shown(
			[
				'[thread](../../wip/threads/TH-0001.md) [site](https://example.com/a.md) [go](main.go)',
				'',
				'`design/issues/I-0001-x.md` `flaiover/README.md` design/system/c.md',
				'',
				'```',
				'design/system/c.md',
				'```',
				''
			].join('\n'),
			'design/system/a.md'
		);
		expect(root.querySelectorAll('button')).toHaveLength(0);
	});

	it('adds each control once, however often it runs, and none when asked', async () => {
		const root = await shown('`design/a.md`\n', 'design/system/a.md');
		await enhance(root, false);
		expect(root.querySelectorAll('button[data-edit]')).toHaveLength(1);
		const plain = document.createElement('div');
		plain.innerHTML = render('`design/a.md`\n', 'design/system/a.md');
		await enhance(plain, false, false);
		expect(plain.querySelectorAll('button')).toHaveLength(0);
	});

	it('renders markdown into an element with its controls, as an attachment', () => {
		const node = document.createElement('div');
		markdownInto('Read [it](design/a.md).\n', '', { story: 'S-0001' })(node);
		expect(node.querySelector('a')!.getAttribute('href')).toBe('/docs/design/a.md');
		control(node, 'Edit design/a.md')!.click();
		expect(openEditor).toHaveBeenCalledWith({
			path: 'design/a.md',
			line: undefined,
			story: 'S-0001'
		});
	});
});
