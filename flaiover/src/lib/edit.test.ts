import { describe, expect, it } from 'vitest';
import { compose, headingAt, headingsOf, lineAt, locate, offsetOfLine, splitRaw } from './edit';
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

describe('lines and offsets', () => {
	const text = 'one\ntwo\n\nfour';
	it('finds the start of a line, clamped to the first and the last', () => {
		expect(offsetOfLine(text, 1)).toBe(0);
		expect(offsetOfLine(text, 2)).toBe(4);
		expect(offsetOfLine(text, 3)).toBe(8);
		expect(offsetOfLine(text, 4)).toBe(9);
		expect(offsetOfLine(text, 99)).toBe(9);
		expect(offsetOfLine(text, 0)).toBe(0);
		expect(offsetOfLine('', 3)).toBe(0);
	});
	it('names the line that holds an offset, the inverse of offsetOfLine', () => {
		for (const n of [1, 2, 3, 4]) expect(lineAt(text, offsetOfLine(text, n))).toBe(n);
		expect(lineAt(text, 6)).toBe(2);
		expect(lineAt(text, -1)).toBe(1);
		expect(lineAt(text, 1000)).toBe(4);
	});
	it('locates a file line in the front matter or the body', () => {
		const { frontMatter, body } = splitRaw(DOC); // two front matter lines, body from file line 5
		expect(locate(frontMatter, 1)).toEqual({ part: 'frontMatter', line: 1 });
		expect(locate(frontMatter, 3)).toEqual({ part: 'frontMatter', line: 2 });
		expect(locate(frontMatter, 4)).toEqual({ part: 'frontMatter', line: 2 });
		expect(locate(frontMatter, 6)).toEqual({ part: 'body', line: 2 });
		expect(body.split('\n')[locate(frontMatter, 8).line - 1]).toBe('## Shape');
		expect(DOC.split('\n')[8 - 1]).toBe('## Shape');
		expect(locate(null, 3)).toEqual({ part: 'body', line: 3 });
		expect(locate('title: X', 3)).toEqual({ part: 'frontMatter', line: 1 });
		expect(locate('title: X', 4)).toEqual({ part: 'body', line: 1 });
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
