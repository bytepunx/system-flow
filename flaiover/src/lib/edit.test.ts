import { describe, expect, it } from 'vitest';
import { compose, headingAt, headingsOf, splitRaw } from './edit';
import { touching } from './touches';

const DOC = '---\ntitle: Plan\nupdated: 2026-09-01\n---\n\n# Plan\n\n## Shape\ntext\n';

describe('front matter split', () => {
	it('round-trips a document byte for byte', () => {
		const { frontMatter, body } = splitRaw(DOC);
		expect(frontMatter).toBe('title: Plan\nupdated: 2026-09-01\n');
		expect(body).toBe('\n# Plan\n\n## Shape\ntext\n');
		expect(compose(frontMatter, body)).toBe(DOC);
	});
	it('leaves a document without front matter alone', () => {
		expect(splitRaw('# Bare\n')).toEqual({ frontMatter: null, body: '# Bare\n' });
		expect(compose(null, '# Bare\n')).toBe('# Bare\n');
		expect(splitRaw('---\nunterminated\n').frontMatter).toBeNull();
	});
	it('closes front matter the designer left without a final newline', () => {
		expect(compose('title: X', 'b\n')).toBe('---\ntitle: X\n---\nb\n');
	});
});

describe('heading under the caret', () => {
	const body = 'intro\n# One\na\n## Two ##\nb\n```\n# not a heading\n```\nc\n# Three\n';
	it('is the last heading at or before the caret', () => {
		expect(headingAt(body, 0)).toBeNull();
		expect(headingAt(body, body.indexOf('a\n'))).toBe('One');
		expect(headingAt(body, body.indexOf('b\n'))).toBe('Two');
		expect(headingAt(body, body.indexOf('c\n'))).toBe('Two'); // the fenced line is code
		expect(headingAt(body, body.length)).toBe('Three');
	});
	it('lists headings the same way', () => {
		expect(headingsOf(body)).toEqual(['One', 'Two', 'Three']);
	});
});

describe('touches', () => {
	const items = [
		{ id: 'S-1', title: 'a', status: 'in-progress', touches: ['design/system'] },
		{ id: 'S-2', title: 'b', status: 'review', touches: ['docs/users/flai.md'] },
		{ id: 'S-3', title: 'c', status: 'backlog', touches: ['design/system'] },
		{ id: 'S-4', title: 'd', status: 'in-progress', touches: ['design/sys'] }
	];
	it('matches the path or a folder above it, for work in progress or review only', () => {
		expect(touching('design/system/plan.md', items).map((w) => w.id)).toEqual(['S-1']);
		expect(touching('docs/users/flai.md', items).map((w) => w.id)).toEqual(['S-2']);
		expect(touching('docs/users/other.md', items)).toEqual([]);
		expect(touching('', items)).toEqual([]);
	});
});
