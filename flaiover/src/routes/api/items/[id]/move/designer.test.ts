import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { Repo, useRepo } from '$lib/server/repo';
import { flaiAsk } from '$lib/server/testing';
import { cp, mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { resetFlaiBinary } from '$lib/server/flai';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');
type Handler = (event: never) => Promise<Response>;

// S-0058: a move made on the board is recorded as the designer's, not as
// "flaiover", which is what FLAI_AGENT is for the dashboard's flai calls.
describe.skipIf(!existsSync(bin))('move endpoint attribution', () => {
	let dir: string;
	let POST: Handler;
	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-move-'));
		await cp(fixture, dir, { recursive: true });
		const manifest = join(dir, 'system-flow.yaml');
		const yaml = (await readFile(manifest, 'utf8')).replace(/^owner:.*\n/m, '');
		await writeFile(manifest, yaml + 'owner: dana\n');
		process.env.PROJECT_DIR = dir;
		useRepo(new Repo(dir, flaiAsk(dir)));
		process.env.FLAI_BIN = bin;
		resetFlaiBinary();
		POST = (await import('./+server')).POST as unknown as Handler;
	});
	afterAll(async () => {
		useRepo(null);
		await rm(dir, { recursive: true, force: true });
	});
	const move = (id: string, body: unknown) =>
		POST({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	it('names the manifest owner when the request names nobody', async () => {
		const r = await move('T-003', { to: 'done' });
		expect(r.status).toBe(200);
		const file = await readFile(join(dir, 'wip/kanban/tasks/T-003-t3.md'), 'utf8');
		expect(file).toMatch(/- to: done\n\s+at: [^\n]+\n\s+by: dana/);
		expect(file).not.toContain('by: flaiover');
	});
});
