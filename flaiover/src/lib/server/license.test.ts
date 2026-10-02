// S-0231: the dashboard shows the license it ships with, read from the image beside the server, or
// from the repository root when run from a checkout.
import { describe, expect, it } from 'vitest';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { licenseName, licensePaths, readLicense } from './license';

describe('the shipped license (S-0231)', () => {
	it('looks beside the server, then one directory up, and first where FLAIOVER_LICENSE points', () => {
		expect(licensePaths({}, '/app')).toEqual(['/app/LICENSE.md', '/LICENSE.md']);
		expect(licensePaths({ FLAIOVER_LICENSE: '/etc/flaiover/terms.md' }, '/app')).toEqual([
			'/etc/flaiover/terms.md',
			'/app/LICENSE.md',
			'/LICENSE.md'
		]);
	});

	it('reads the first file that exists and names it by its heading', () => {
		const dir = mkdtempSync(join(tmpdir(), 'flaiover-license-'));
		const path = join(dir, 'LICENSE.md');
		writeFileSync(path, '# Example License 2.0\n\nUse it well.\n');
		const license = readLicense([join(dir, 'missing.md'), path]);
		expect(license).toEqual({
			name: 'Example License 2.0',
			text: '# Example License 2.0\n\nUse it well.\n',
			path
		});
		expect(readLicense([join(dir, 'missing.md')])).toBeNull();
	});

	it('finds the repository root LICENSE.md from a checkout of flaiover', () => {
		const license = readLicense();
		expect(license?.path).toBe(resolve('..', 'LICENSE.md'));
		expect(license?.name).toBe('Bytepunx Shield License 1.0');
		expect(license?.text).toContain('## Prohibited Uses');
	});

	it('takes the name from the first heading, however deep', () => {
		expect(licenseName('## Terms\n\nbody')).toBe('Terms');
		expect(licenseName('')).toBe('');
	});
});
