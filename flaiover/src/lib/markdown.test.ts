// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const { openEditor } = vi.hoisted(() => ({ openEditor: vi.fn() }));
vi.mock('./editor.svelte', () => ({ openEditor }));

import { enhance, markdownInto, render, rewriteLink, resolveRelative } from './markdown';
import { projectState } from './project.svelte';

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

describe('raw HTML (I-0294)', () => {
	function shown(markdown: string) {
		const root = document.createElement('div');
		root.innerHTML = render(markdown, 'design/system/a.md');
		return root;
	}

	it('shows a script block and an inline handler as text, not markup', () => {
		const root = shown(
			'<script>alert(1)</script>\n\nText <img src=x onerror=alert(1)> here.\n\n<div onclick="alert(1)">x</div>\n'
		);
		expect(root.querySelector('script')).toBeNull();
		expect(root.querySelector('img')).toBeNull();
		expect(root.querySelector('[onerror], [onclick]')).toBeNull();
		expect(root.textContent).toContain('<script>alert(1)</script>');
		expect(root.textContent).toContain('<img src=x onerror=alert(1)>');
	});

	it('gives a javascript: link no href that runs it', () => {
		const root = shown('[x](javascript:alert(1)) and <javascript:alert(1)>\n');
		const hrefs = [...root.querySelectorAll('a')].map((a) => a.getAttribute('href') ?? '');
		expect(hrefs.filter((h) => /^\s*javascript:/i.test(h))).toEqual([]);
		expect(root.textContent).toContain('[x](javascript:alert(1))');
	});

	it('still renders criteria checkboxes and heading permalinks', () => {
		const root = shown('## Acceptance criteria\n\n- [ ] open\n- [x] done\n');
		expect(root.querySelectorAll('input[type="checkbox"]')).toHaveLength(2);
		const heading = root.querySelector('h2')!;
		expect(heading.id).toBe('acceptance-criteria');
		expect(heading.querySelector('a')!.getAttribute('href')).toBe('#acceptance-criteria');
	});

	it('rewrites each link once however many times it renders', () => {
		for (let i = 0; i < 50; i++) render('[a](b.md)\n', 'design/system/a.md');
		expect(render('[a](b.md)\n', 'design/x/a.md')).toContain('href="/docs/design/x/b.md"');
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

describe('log paths (S-0785)', () => {
	afterEach(() => {
		projectState.current = null;
	});

	function shown(markdown: string) {
		const root = document.createElement('div');
		root.innerHTML = render(markdown, '');
		return root;
	}
	const logLinks = (root: HTMLElement) =>
		[...root.querySelectorAll<HTMLAnchorElement>('a[data-path]')].map((a) => ({
			path: a.dataset.path,
			href: a.getAttribute('href'),
			code: a.querySelector('code')?.textContent ?? null
		}));
	const href = (path: string) => `/ops/log?${new URLSearchParams({ path })}`;

	it('links an absolute, a ~/, and a .flai-cache/ log path in text to the log page', () => {
		const root = shown(
			'Deploy failed: see /home/op/.flai/serve/kluster/s-0785.log, ~/.flai/serve/flai-serve.log and .flai-cache/kluster/logs/s-0785/build-api.log.\n'
		);
		expect(logLinks(root)).toEqual([
			{
				path: '/home/op/.flai/serve/kluster/s-0785.log',
				href: href('/home/op/.flai/serve/kluster/s-0785.log'),
				code: null
			},
			{
				path: '~/.flai/serve/flai-serve.log',
				href: href('~/.flai/serve/flai-serve.log'),
				code: null
			},
			{
				path: '.flai-cache/kluster/logs/s-0785/build-api.log',
				href: href('.flai-cache/kluster/logs/s-0785/build-api.log'),
				code: null
			}
		]);
		expect(root.textContent!.trim()).toBe(
			'Deploy failed: see /home/op/.flai/serve/kluster/s-0785.log, ~/.flai/serve/flai-serve.log and .flai-cache/kluster/logs/s-0785/build-api.log.'
		);
	});

	it('links a log path in backticks with its code element inside the link', () => {
		const root = shown(
			'deploy of S-0785 failed at build-api; its log is `/home/op/a.log`, `~/b.log`, and `.flai-cache/c.log`\n'
		);
		expect(logLinks(root)).toEqual([
			{ path: '/home/op/a.log', href: href('/home/op/a.log'), code: '/home/op/a.log' },
			{ path: '~/b.log', href: href('~/b.log'), code: '~/b.log' },
			{ path: '.flai-cache/c.log', href: href('.flai-cache/c.log'), code: '.flai-cache/c.log' }
		]);
		expect(root.querySelector('a[data-path] > code')).not.toBeNull();
		expect(root.querySelector('a[data-path]')!.getAttribute('title')).toBe('Show /home/op/a.log');
	});

	it('links only the log path in a code span that says more, the rest staying code', () => {
		const root = shown('Run `tail -f /tmp/x.log` to follow it.\n');
		expect(logLinks(root)).toEqual([
			{ path: '/tmp/x.log', href: href('/tmp/x.log'), code: '/tmp/x.log' }
		]);
		expect([...root.querySelectorAll('code')].map((c) => c.textContent)).toEqual([
			'tail -f ',
			'/tmp/x.log'
		]);
	});

	it('keeps the project this tab views on the link', () => {
		projectState.current = 'shop';
		const [link] = logLinks(shown('see `/tmp/x.log`\n'));
		expect(link.href).toBe(`${href('/tmp/x.log')}&project=shop`);
	});

	it('leaves a fenced block, a path not ending in .log, a URL, and a markdown link as they are', () => {
		const root = shown(
			[
				'```',
				'/home/op/.flai/serve/kluster/s-0785.log',
				'```',
				'',
				'    /tmp/indented.log',
				'',
				'see /home/op/a.log.bak, `~/b.txt`, https://example.com/logs/c.log and [the ADR](design/adrs/0001-x.md)',
				''
			].join('\n')
		);
		expect(logLinks(root)).toEqual([]);
		expect(root.querySelector('pre')!.textContent).toBe(
			'/home/op/.flai/serve/kluster/s-0785.log\n'
		);
		expect([...root.querySelectorAll('a')].map((a) => a.getAttribute('href'))).toEqual([
			'https://example.com/logs/c.log',
			'/docs/design/adrs/0001-x.md'
		]);
	});

	it('leaves a log path that is the text of a link to that link', () => {
		const root = shown(
			'[the log](https://example.com/a) and [/tmp/x.log](https://example.com/x)\n'
		);
		expect(logLinks(root)).toEqual([]);
		expect(root.querySelectorAll('a')).toHaveLength(2);
	});
});
