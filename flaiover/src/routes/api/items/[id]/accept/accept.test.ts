import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { Repo, useRepo } from '$lib/server/repo';
import { flaiAsk } from '$lib/server/testing';
import { chmod, mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { resetFlaiBinary } from '$lib/server/flai';

type Handler = (event: never) => Promise<Response>;

// A fake flai (S-0041): a script standing in for the binary through FLAI_BIN. It records its
// arguments, writes what $FAKE_STDERR holds to stderr a line at a time with a pause between, then
// $FAKE_STDOUT to stdout, and exits with $FAKE_EXIT.
const FAKE = `#!/bin/sh
printf '%s\\n' "$@" > "$FAKE_DIR/args"
if [ -n "$FAKE_STDERR" ]; then
  printf '%s\\n' "$FAKE_STDERR" | while IFS= read -r line; do printf '%s\\n' "$line" >&2; sleep 0.05; done
fi
[ -n "$FAKE_STDOUT" ] && printf '%s\\n' "$FAKE_STDOUT"
exit "\${FAKE_EXIT:-0}"
`;

describe('acceptance and diff endpoints with a fake flai', () => {
	let dir: string;
	let accept: Handler;
	let diff: Handler;
	const args = async () => (await readFile(join(dir, 'args'), 'utf8')).trim().split('\n');
	const lines = async (r: Response) =>
		(await r.text())
			.trim()
			.split('\n')
			.map((l) => JSON.parse(l));
	const post = (id: string, body: unknown) =>
		accept({
			params: { id },
			request: new Request('http://x', { method: 'POST', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-fake-'));
		await writeFile(
			join(dir, 'system-flow.yaml'),
			'version: 1\nname: t\nowner: dana\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n'
		);
		const bin = join(dir, 'flai');
		await writeFile(bin, FAKE);
		await chmod(bin, 0o755);
		process.env.PROJECT_DIR = dir;
		useRepo(new Repo(dir, flaiAsk(dir)));
		process.env.FLAI_BIN = bin;
		process.env.FAKE_DIR = dir;
		resetFlaiBinary();
		accept = (await import('./+server')).POST as unknown as Handler;
		diff = (await import('../diff/+server')).GET as unknown as Handler;
	});
	beforeEach(() => {
		delete process.env.FAKE_STDERR;
		delete process.env.FAKE_STDOUT;
		delete process.env.FAKE_EXIT;
	});
	afterAll(async () => {
		useRepo(null);
		await rm(dir, { recursive: true, force: true });
	});

	it('accepts as the designer and streams each step before the result', async () => {
		process.env.FAKE_STDERR = [
			'{"level":"INFO","msg":"story branch merged","branch":"story/S-0041"}',
			'{"level":"INFO","msg":"acceptance step","item":"S-0041","step":"merged","detail":"story/S-0041 rebased and fast-forwarded into the main branch"}',
			'{"level":"INFO","msg":"acceptance step","item":"S-0041","step":"committed","detail":"chore: [S-0041] accept and archive"}',
			'{"level":"WARN","msg":"accepted locally but not pushed","detail":"run: git push origin HEAD"}',
			'{"level":"INFO","msg":"acceptance step","item":"S-0041","step":"not-pushed","detail":"accepted locally; run: git push origin HEAD"}'
		].join('\n');
		process.env.FAKE_STDOUT =
			'{"id":"S-0041","status":"done","merged":true,"tags":["flaiover/v0.13.0"]}';
		const r = await post('S-0041', {});
		expect(r.status).toBe(200);
		expect(r.headers.get('content-type')).toBe('application/x-ndjson');
		const out = await lines(r);
		expect(out.map((l) => l.event)).toEqual([
			'progress',
			'progress',
			'warning',
			'progress',
			'done'
		]);
		expect(out[0]).toMatchObject({
			step: 'merged',
			msg: expect.stringContaining('fast-forwarded')
		});
		expect(out.at(-1).result).toMatchObject({ status: 'done', tags: ['flaiover/v0.13.0'] });
		expect(await args()).toEqual(['accept', 'S-0041', '--by', 'dana', '--json']);
	});

	it('passes --yes only when the designer chose to include uncommitted files', async () => {
		process.env.FAKE_STDOUT = '{"id":"S-0041","status":"done"}';
		await (await post('S-0041', { include_uncommitted: true })).text();
		expect(await args()).toEqual(['accept', 'S-0041', '--by', 'dana', '--yes', '--json']);
		await (await post('S-0041', { include_uncommitted: 'yes' })).text();
		expect(await args()).not.toContain('--yes');
	});

	it('ends with flai’s message verbatim when the acceptance fails', async () => {
		const message =
			'rebase of story/S-0041 onto main stopped with conflicts in docs/users/flai.md; resolve them in .flai-cache/worktrees/S-0041 and run flai stream sync S-0041';
		process.env.FAKE_STDERR = JSON.stringify({
			level: 'FATAL',
			msg: 'command failed',
			err: message
		});
		process.env.FAKE_EXIT = '1';
		const out = await lines(await post('S-0041', {}));
		expect(out).toHaveLength(1);
		expect(out[0]).toEqual({ event: 'error', status: 500, error: message });
	});

	it('passes the branch diff through', async () => {
		process.env.FAKE_STDOUT =
			'{"story":"S-0041","branch":"story/S-0041","files":[{"path":"a.md","status":"added","additions":1,"deletions":0,"patch":"@@ -0,0 +1 @@\\n+a"}]}';
		const r = await diff({ params: { id: 'S-0041' } } as never);
		expect(r.status).toBe(200);
		expect((await r.json()).files[0]).toMatchObject({ path: 'a.md', status: 'added' });
		expect(await args()).toEqual(['stream', 'diff', 'S-0041', '--json']);
	});

	it('answers a story without a branch with flai’s reason', async () => {
		process.env.FAKE_STDERR = JSON.stringify({
			level: 'FATAL',
			msg: 'command failed',
			err: 'S-0002 has no branch story/S-0002: it was never opened with flai stream open, or it has been merged and removed'
		});
		process.env.FAKE_EXIT = '1';
		const r = await diff({ params: { id: 'S-0002' } } as never);
		expect(r.status).toBe(500);
		expect((await r.json()).error).toContain('has no branch story/S-0002');
	});
});
