import { afterAll, beforeAll, describe, expect, it, vi } from 'vitest';
import { execFileSync } from 'node:child_process';
import { cp, mkdtemp, mkdir, readFile, rm, writeFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo, useRepo } from '$lib/server/repo';
import { flaiAsk } from '$lib/server/testing';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');
const DOC =
	'---\ntitle: Plan\nupdated: 2026-09-01\nstatus: active\n---\n\n# Plan\n\n## Shape\ntext\n';

type Handler = (event: never) => Promise<Response>;

describe.skipIf(!existsSync(bin))('document editing endpoints', () => {
	let dir: string;
	let GET: Handler;
	let PUT: Handler;
	const show = async (path: string) =>
		(await GET({ url: new URL(`http://x/api/docs/edit?path=${path}`) } as never)).json();
	const save = (body: unknown) =>
		PUT({
			request: new Request('http://x/api/docs/file', { method: 'PUT', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-edit-'));
		await cp(fixture, dir, { recursive: true });
		await mkdir(join(dir, 'design/system'), { recursive: true });
		await writeFile(join(dir, 'design/system/plan.md'), DOC);
		useRepo(new Repo(dir, flaiAsk(dir)));
		GET = (await import('../edit/+server')).GET as unknown as Handler;
		PUT = (await import('./+server')).PUT as unknown as Handler;
	});
	afterAll(async () => {
		useRepo(null);
		await rm(dir, { recursive: true, force: true });
	});

	it('loads a document with its hash and mode', async () => {
		const doc = await show('design/system/plan.md');
		expect(doc).toMatchObject({ path: 'design/system/plan.md', mode: 'full', content: DOC });
		expect(doc.hash).toMatch(/^[0-9a-f]{64}$/);
		const story = await show('wip/kanban/stories/S-004-four.md');
		expect(story.mode).toBe('body');
		expect(story.reason).toContain('front matter');
	});

	it('saves through flai and returns the new hash', async () => {
		const doc = await show('design/system/plan.md');
		const r = await save({
			path: doc.path,
			hash: doc.hash,
			content: doc.content.replace('text\n', 'better text\n')
		});
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body.hash).not.toBe(doc.hash);
		expect(body.committed).toBe(false); // the fixture is not a git repository
		const file = await readFile(join(dir, 'design/system/plan.md'), 'utf8');
		expect(file).toContain('better text');
	});

	it('answers a stale hash with 409, the current version, and a diff', async () => {
		const doc = await show('design/system/plan.md');
		await writeFile(
			join(dir, 'design/system/plan.md'),
			doc.content.replace('better text', 'their text')
		);
		const r = await save({
			path: doc.path,
			hash: doc.hash,
			content: doc.content.replace('better text', 'my text')
		});
		expect(r.status).toBe(409);
		const body = await r.json();
		expect(body.error).toContain('changed after it was loaded');
		expect(body.current).toContain('their text');
		expect(body.hash).toMatch(/^[0-9a-f]{64}$/);
		expect(body.diff).toContain('+my text');
		expect(await readFile(join(dir, 'design/system/plan.md'), 'utf8')).toContain('their text');
	});

	it('answers content flai refuses with 422 and the findings, and restores the file', async () => {
		const doc = await show('design/system/plan.md');
		const r = await save({
			path: doc.path,
			hash: doc.hash,
			content: doc.content.replace('title: Plan\n', '')
		});
		expect(r.status).toBe(422);
		const body = await r.json();
		expect(body.findings.map((f: { rule: string }) => f.rule)).toContain('doc.title');
		expect(await readFile(join(dir, 'design/system/plan.md'), 'utf8')).toBe(doc.content);
	});

	it('refuses a change to front matter flai owns with 422', async () => {
		const story = await show('wip/kanban/stories/S-004-four.md');
		const r = await save({
			path: story.path,
			hash: story.hash,
			content: story.content.replace('status: in-progress', 'status: done')
		});
		expect(r.status).toBe(422);
		expect((await r.json()).error).toContain('front matter');
	});

	it('requires path, content, and hash', async () => {
		expect((await save({ path: 'design/system/plan.md' })).status).toBe(400);
	});

	it('refuses a story that is not a story ID with 400, before asking flai', async () => {
		for (const story of ['T-0001', 'S-1 --no-commit', '--help', 42]) {
			const r = await save({ path: 'design/system/plan.md', content: DOC, hash: 'h', story });
			expect(r.status, String(story)).toBe(400);
			expect((await r.json()).error).toContain('story must be');
		}
		const shown = await GET({
			url: new URL('http://x/api/docs/edit?path=design/system/plan.md&story=E-001')
		} as never);
		expect(shown.status).toBe(400);
	});
});

// S-0378: with a story, the editor loads and saves the file in the story's worktree, and the save
// is committed on story/<id>, through the flai built from this tree, as the dashboard's link names.
describe.skipIf(!existsSync(bin))("document editing on a story's branch", () => {
	let dir: string;
	let wt: string;
	let GET: Handler;
	let PUT: Handler;
	const PLAN = 'design/system/plan.md';
	const git = (cwd: string, ...args: string[]) =>
		execFileSync('git', args, { cwd, env: process.env }).toString().trim();
	const show = async (story: string) =>
		(
			await GET({
				url: new URL(`http://x/api/docs/edit?path=${PLAN}&story=${story}`)
			} as never)
		).json();
	const save = (body: unknown) =>
		PUT({
			request: new Request('http://x/api/docs/file', { method: 'PUT', body: JSON.stringify(body) })
		} as never);

	beforeAll(async () => {
		for (const v of ['GIT_AUTHOR', 'GIT_COMMITTER']) {
			vi.stubEnv(`${v}_NAME`, 'Dana Designer');
			vi.stubEnv(`${v}_EMAIL`, 'dana@example.com');
		}
		dir = await mkdtemp(join(tmpdir(), 'flaiover-edit-story-'));
		await cp(fixture, dir, { recursive: true });
		await mkdir(join(dir, 'design/system'), { recursive: true });
		await writeFile(join(dir, PLAN), DOC);
		await writeFile(join(dir, '.gitignore'), '.flai-cache/\n');
		git(dir, 'init', '-q', '-b', 'main');
		git(dir, 'add', '-A');
		git(dir, 'commit', '-q', '-m', 'init');
		wt = join(dir, '.flai-cache/worktrees/S-004');
		git(dir, 'worktree', 'add', '-q', '-b', 'story/S-004', wt);
		await writeFile(join(wt, PLAN), DOC.replace('text\n', 'branch text\n'));
		git(wt, 'commit', '-q', '-am', 'docs: [S-004] branch plan');
		useRepo(new Repo(dir, flaiAsk(dir)));
		GET = (await import('../edit/+server')).GET as unknown as Handler;
		PUT = (await import('./+server')).PUT as unknown as Handler;
	});
	afterAll(async () => {
		vi.unstubAllEnvs();
		useRepo(null);
		await rm(dir, { recursive: true, force: true });
	});

	it("loads the story's file and commits the save on its branch, the main checkout left as it was", async () => {
		const mainHead = git(dir, 'rev-parse', 'HEAD');
		const doc = await show('S-004');
		expect(doc).toMatchObject({ path: PLAN, mode: 'full' });
		expect(doc.content).toContain('branch text');
		const r = await save({
			path: PLAN,
			hash: doc.hash,
			content: doc.content.replace('branch text', 'better branch text'),
			message: 'clarify the shape',
			story: 'S-004'
		});
		expect(r.status).toBe(200);
		const body = await r.json();
		expect(body).toMatchObject({ path: PLAN, committed: true });
		expect(git(dir, 'rev-parse', '--short', 'story/S-004')).toBe(body.commit);
		expect(git(wt, 'log', '-1', '--format=%s')).toBe('docs: [S-004] clarify the shape');
		expect(await readFile(join(wt, PLAN), 'utf8')).toContain('better branch text');
		expect(await readFile(join(dir, PLAN), 'utf8')).toBe(DOC);
		expect(git(dir, 'rev-parse', 'HEAD')).toBe(mainHead);
	});

	it("answers a stale hash on the story's branch with 409 and the branch's current version", async () => {
		const doc = await show('S-004');
		await writeFile(join(wt, PLAN), doc.content.replace('better branch text', 'their text'));
		const r = await save({
			path: PLAN,
			hash: doc.hash,
			content: doc.content.replace('better branch text', 'my text'),
			story: 'S-004'
		});
		expect(r.status).toBe(409);
		const body = await r.json();
		expect(body.current).toContain('their text');
		expect(body.diff).toContain('+my text');
		expect(await readFile(join(wt, PLAN), 'utf8')).toContain('their text');
	});

	it("answers content flai refuses on the story's branch with 422 and its findings", async () => {
		const doc = await show('S-004');
		const r = await save({
			path: PLAN,
			hash: doc.hash,
			content: doc.content.replace('title: Plan\n', ''),
			story: 'S-004'
		});
		expect(r.status).toBe(422);
		expect((await r.json()).findings.map((f: { rule: string }) => f.rule)).toContain('doc.title');
		expect(await readFile(join(wt, PLAN), 'utf8')).toBe(doc.content);
	});
});
