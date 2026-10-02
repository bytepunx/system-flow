// S-0231: GET /api/license answers the license the image carries.
import { afterEach, describe, expect, it } from 'vitest';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

describe('GET /api/license', () => {
	const before = process.env.FLAIOVER_LICENSE;
	afterEach(() => {
		if (before === undefined) delete process.env.FLAIOVER_LICENSE;
		else process.env.FLAIOVER_LICENSE = before;
	});

	it('answers the name and text of the license at FLAIOVER_LICENSE first', async () => {
		const dir = mkdtempSync(join(tmpdir(), 'flaiover-license-api-'));
		const path = join(dir, 'terms.md');
		writeFileSync(path, '# Example License 2.0\n\nUse it well.\n');
		process.env.FLAIOVER_LICENSE = path;
		const { GET } = await import('./+server');
		const res = GET();
		expect(res.status).toBe(200);
		expect(await res.json()).toEqual({
			name: 'Example License 2.0',
			text: '# Example License 2.0\n\nUse it well.\n'
		});
	});

	it('answers the repository root LICENSE.md from a checkout', async () => {
		delete process.env.FLAIOVER_LICENSE;
		const { GET } = await import('./+server');
		const res = GET();
		expect(res.status).toBe(200);
		const body = await res.json();
		expect(body.name).toBe('system-flow Shield License 1.0');
		expect(body.text).toContain('## Permitted Purposes');
	});
});
