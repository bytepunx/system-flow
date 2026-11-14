import { describe, expect, it } from 'vitest';
import { logHref, splitLogPaths } from './logs';

// S-0785: the log paths in free text from flai, each to be a link to the log page, and the rest
// plain text; and the log page's href for one.
const paths = (text: string) => splitLogPaths(text).flatMap((s) => ('path' in s ? [s.path] : []));

describe('splitLogPaths', () => {
	it('splits a message into its text and the log path it names', () => {
		const text =
			'planner started for S-0785 (pid 42); its output is in /home/op/.flai/serve/agents/planner.log on the host';
		expect(splitLogPaths(text)).toEqual([
			{ text: 'planner started for S-0785 (pid 42); its output is in ' },
			{ path: '/home/op/.flai/serve/agents/planner.log' },
			{ text: ' on the host' }
		]);
	});

	it('finds several paths of each kind', () => {
		const text =
			'see .flai-cache/kluster/logs/s-0595/build-api.log, ~/.flai/serve/flai-serve.log and /tmp/x.log';
		expect(paths(text)).toEqual([
			'.flai-cache/kluster/logs/s-0595/build-api.log',
			'~/.flai/serve/flai-serve.log',
			'/tmp/x.log'
		]);
	});

	it('leaves trailing punctuation out of a path', () => {
		for (const [text, path] of [
			['its output is in /tmp/a.log.', '/tmp/a.log'],
			['in /tmp/a.log, then', '/tmp/a.log'],
			['in /tmp/a.log; then', '/tmp/a.log'],
			['in /tmp/a.log: read it', '/tmp/a.log'],
			['(in /tmp/a.log)', '/tmp/a.log'],
			['(in /tmp/a.log).', '/tmp/a.log']
		]) {
			const segments = splitLogPaths(text);
			expect(paths(text), text).toEqual([path]);
			expect(segments.map((s) => ('path' in s ? s.path : s.text)).join(''), text).toBe(text);
		}
		expect(splitLogPaths('in /tmp/a.log.')).toEqual([
			{ text: 'in ' },
			{ path: '/tmp/a.log' },
			{ text: '.' }
		]);
	});

	it('finds a path inside backticks and quotes', () => {
		expect(splitLogPaths('its output: `/tmp/a.log`')).toEqual([
			{ text: 'its output: `' },
			{ path: '/tmp/a.log' },
			{ text: '`' }
		]);
		expect(paths(`read "~/.flai/x.log" or '/tmp/y.log'`)).toEqual(['~/.flai/x.log', '/tmp/y.log']);
	});

	it('cuts a path at whitespace', () => {
		expect(paths('/tmp/a b.log')).toEqual([]);
		expect(paths('/tmp/a.log\n/tmp/b.log\tdone')).toEqual(['/tmp/a.log', '/tmp/b.log']);
	});

	it('finds nothing in text without a log path', () => {
		expect(splitLogPaths('')).toEqual([]);
		expect(splitLogPaths('nothing to see')).toEqual([{ text: 'nothing to see' }]);
		for (const text of [
			'the catalog is at /usr/share/catalog',
			'see /var/log for more',
			'the logs are in /tmp/run.logs',
			'a backup /tmp/a.log.1',
			'a folder /tmp/a.log/b',
			'a line /tmp/a.log:12',
			'relative wip/a.log and design/x.log',
			'inside a word foo/bar.log and x.flai-cache/a.log',
			'a URL https://example.com/a.log',
			'.flai-cachex/a.log',
			'~user/a.log'
		])
			expect(paths(text), text).toEqual([]);
	});

	it('keeps text with no path as one run', () => {
		expect(splitLogPaths('the catalog')).toEqual([{ text: 'the catalog' }]);
	});
});

describe('logHref', () => {
	it('links the log page with the path, and the project when one is chosen', () => {
		const path = '~/.flai/serve/agents/s-0785 run.log';
		expect(logHref(path)).toBe('/ops/log?path=%7E%2F.flai%2Fserve%2Fagents%2Fs-0785+run.log');
		expect(logHref(path, null)).toBe(logHref(path));
		const href = new URL(logHref(path, 'acme & co'), 'http://x');
		expect(href.pathname).toBe('/ops/log');
		expect(href.searchParams.get('path')).toBe(path);
		expect(href.searchParams.get('project')).toBe('acme & co');
	});
});
