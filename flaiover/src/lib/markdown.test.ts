import { describe, expect, it } from 'vitest';
import { render, rewriteLink, resolveRelative } from './markdown';

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
