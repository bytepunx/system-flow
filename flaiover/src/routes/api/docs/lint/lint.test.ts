import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask } from '$lib/server/repo';
import { flaiAsk, haveFlai } from '$lib/server/testing';

// POST /api/docs/lint answers flai's doc.lint on the host (S-0378): the findings of the content the
// editor holds, through the flai built from this tree, against a copy of the fixture.
const fixture = resolve('../flai/internal/metrics/testdata/good');

type Handler = (event: never) => Promise<Response>;

describe.skipIf(!haveFlai)('document lint endpoint', () => {
	let dir: string;
	let POST: Handler;
	const lint = (body: unknown) =>
		POST({
			request: new Request('http://x/api/docs/lint', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-lint-'));
		await cp(fixture, dir, { recursive: true });
		await writeFile(join(dir, '.markdownlint.yaml'), 'default: true\nMD013: false\nMD041: false\n');
		POST = (await import('./+server')).POST as unknown as Handler;
	});
	afterEach(() => useRepo(null));
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('answers the findings of the content with their lines and rules', async () => {
		useRepo(new Repo(dir, flaiAsk(dir)));
		const r = await lint({
			path: 'design/system/plan.md',
			content: '# Title\n\n\n## Same\n\ntext \n\n## Same\n'
		});
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.path).toBe('design/system/plan.md');
		expect(body.findings).toEqual([
			{
				line: 3,
				rule: 'MD012',
				alias: 'no-multiple-blanks',
				message: 'Multiple consecutive blank lines [Expected: 1; Actual: 2]'
			},
			{
				line: 6,
				rule: 'MD009',
				alias: 'no-trailing-spaces',
				message: 'Trailing spaces [Expected: 0 or 2; Actual: 1]'
			},
			{
				line: 8,
				rule: 'MD024',
				alias: 'no-duplicate-heading',
				message: 'Multiple headings with the same content [Context: "Same"]'
			}
		]);
	});

	it('answers clean content with no findings', async () => {
		useRepo(new Repo(dir, flaiAsk(dir)));
		const r = await lint({ path: 'docs/users/new.md', content: '# Title\n\ntext\n' });
		expect(r.status).toBe(200);
		expect((await r.json()).findings).toEqual([]);
	});

	it('answers 400 for a path or content missing, and a path flai refuses', async () => {
		useRepo(new Repo(dir, flaiAsk(dir)));
		expect((await lint({ content: '# Title\n' })).status).toBe(400);
		expect((await lint({ path: 'design/system/plan.md' })).status).toBe(400);
		const outside = await lint({ path: '../outside.md', content: '# Title\n' });
		expect(outside.status).toBe(400);
	});

	it('answers 501 when the host flai is older than doc.lint', async () => {
		const flai = flaiAsk(dir);
		const older = (async (method: string, params?: Record<string, unknown>) => {
			if (method === 'doc.lint')
				throw new AgentError(502, 'flai offers no method doc.lint', -32601);
			return flai(method, params);
		}) as unknown as Ask;
		useRepo(new Repo(dir, older));
		const r = await lint({ path: 'design/system/plan.md', content: '# Title\n' });
		expect(r.status).toBe(501);
		expect((await r.json()).error).toContain('doc.lint');
	});
});
