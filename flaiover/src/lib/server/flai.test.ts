import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, readdir, readFile, rm, writeFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo, RepoError } from './repo';
import { board } from './board';
import { flai, flaiBinary, resetFlaiBinary, withJson } from './flai';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');
const haveFlai = existsSync(bin);

describe('withJson', () => {
	it('asks for JSON before a -- that ends the flags, else at the end', () => {
		expect(withJson(['move', 'S-1', 'ready'])).toEqual(['move', 'S-1', 'ready', '--json']);
		expect(withJson(['story', 'new', '--epic=E-1', '--', '--json is my title'])).toEqual([
			'story',
			'new',
			'--epic=E-1',
			'--json',
			'--',
			'--json is my title'
		]);
	});
});

describe('board reader', () => {
	it('groups active items by state with ages and limits', async () => {
		const b = await board(new Repo(fixture), new Date('2026-09-01T12:00:00Z'));
		expect(b.wip_limits).toEqual({ ready: 5, 'in-progress': 2, review: 3 });
		expect(b.columns['in-progress'].map((c) => c.id)).toEqual(['S-004', 'T-003']);
		const s4 = b.columns['in-progress'][0];
		expect(s4.age_seconds).toBe(93600);
		expect(s4.blocked).toBe(false);
		expect(b.columns.done).toEqual([]); // done items in the fixture are archived
	});
	it('names the parent of every card that has one', async () => {
		const b = await board(new Repo(fixture), new Date('2026-09-01T12:00:00Z'));
		const [story, task] = b.columns['in-progress'];
		expect(story).toMatchObject({ id: 'S-004', parent: 'E-001', parent_title: 'Epic' });
		expect(task).toMatchObject({ id: 'T-003', parent: 'S-004', parent_title: 'Four' });
		const epic = b.columns.backlog.find((c) => c.id === 'E-001')!;
		expect(epic.parent).toBeUndefined();
		expect(epic.parent_title).toBeUndefined();
	});
});

describe.skipIf(!haveFlai)('flai wrapper on a temp project', () => {
	let dir: string;
	beforeAll(async () => {
		process.env.FLAI_BIN = bin;
		resetFlaiBinary();
		dir = await mkdtemp(join(tmpdir(), 'flaiover-'));
		await cp(fixture, dir, { recursive: true });
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('finds the binary', async () => {
		expect(await flaiBinary()).toBe(bin);
	});
	it('performs a valid move and reports the new status', async () => {
		const { data } = await flai<{ id: string; status: string }>(dir, [
			'move',
			'T-003',
			'done',
			'--by',
			'test'
		]);
		expect(data).toMatchObject({ id: 'T-003', status: 'done' });
		const file = await readFile(join(dir, 'wip/kanban/tasks/T-003-t3.md'), 'utf8');
		expect(file).toContain('status: done');
		expect(file).toMatch(/- to: done\n\s+at: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\n\s+by: test/);
	});
	it('refuses an invalid move with the rule text as a 400', async () => {
		await expect(flai(dir, ['move', 'S-004', 'done'])).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('cannot go from in-progress to done')
		});
		await expect(flai(dir, ['move', 'S-004', 'cancelled'])).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('needs --reason')
		});
	});
	it('moves a story with no tasks to ready and in-progress, and is refused at review (ADR-0021)', async () => {
		const { data: made } = await flai<{ id: string }>(dir, [
			'story',
			'new',
			'No tasks yet',
			'--epic',
			'E-001'
		]);
		const stories = join(dir, 'wip/kanban/stories');
		const name = (await readdir(stories)).find((f) => f.startsWith(made.id + '-'))!;
		const path = join(stories, name);
		const body = await readFile(path, 'utf8');
		await writeFile(path, body.replace('- [ ]\n', '- [ ] works\n'));
		for (const to of ['ready', 'in-progress']) {
			const { data } = await flai<{ status: string }>(dir, ['move', made.id, to, '--by', 'test']);
			expect(data.status).toBe(to);
		}
		await expect(flai(dir, ['move', made.id, 'review', '--by', 'test'])).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('needs at least one task before it goes to review')
		});
	});
	it('places a story with flai order, and the board reads the column in that order (S-0057)', async () => {
		const ids: string[] = [];
		for (const title of ['Order one', 'Order two', 'Order three']) {
			const { data } = await flai<{ id: string }>(dir, ['story', 'new', title, '--epic', 'E-001']);
			ids.push(data.id);
		}
		const [one, two, three] = ids;
		const backlog = async () =>
			(await board(new Repo(dir))).columns.backlog
				.filter((c) => ids.includes(c.id))
				.map((c) => c.id);
		expect(await backlog()).toEqual([one, two, three]);
		const { data } = await flai<{ status: string; sequence: string[]; order: string[] }>(dir, [
			'order',
			three,
			'--before',
			one
		]);
		expect(data.status).toBe('backlog');
		expect(data.sequence.filter((id) => ids.includes(id))).toEqual([three, one, two]);
		expect(await backlog()).toEqual([three, one, two]);
		// an epic sharing the column keeps its place
		const column = (await board(new Repo(dir))).columns.backlog;
		expect(column.findIndex((c) => c.id === 'E-001')).toBe(
			column.findIndex((c) => c.type !== 'story')
		);
		await expect(flai(dir, ['order', three, '--before', 'S-004'])).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('within one column')
		});
		await expect(flai(dir, ['order', 'E-001', '--top'])).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('only stories are in the pull order')
		});
	});
	it('asks flai, offline, whether an acceptance is unpushed (S-0063)', async () => {
		const { execFileSync } = await import('node:child_process');
		const base = await mkdtemp(join(tmpdir(), 'flaiover-unpushed-'));
		const git = (cwd: string, ...args: string[]) =>
			execFileSync('git', args, {
				cwd,
				env: {
					...process.env,
					GIT_AUTHOR_NAME: 't',
					GIT_AUTHOR_EMAIL: 't@t',
					GIT_COMMITTER_NAME: 't',
					GIT_COMMITTER_EMAIL: 't@t'
				}
			}).toString();
		try {
			const clone = join(base, 'clone');
			await cp(fixture, clone, { recursive: true });
			git(base, 'init', '-q', '--bare', '-b', 'main', 'origin.git');
			git(clone, 'init', '-q', '-b', 'main');
			git(clone, 'add', '-A');
			git(clone, 'commit', '-q', '-m', 'init');
			git(clone, 'remote', 'add', 'origin', join(base, 'origin.git'));
			git(clone, 'push', '-q', '-u', 'origin', 'main');
			const ask = () =>
				flai<{
					pushed: boolean;
					reason?: string;
					unpushed?: { acceptances: string[]; tags: string[] };
				}>(clone, ['push', '--pending', '--dry-run']);
			expect((await ask()).data).toMatchObject({ pushed: false, reason: 'nothing pending' });
			git(
				clone,
				'commit',
				'-q',
				'--allow-empty',
				'-m',
				'chore: [S-004] accept and archive; release cli 1.1.0'
			);
			git(clone, 'tag', '-a', 'cli/v1.1.0', '-m', 'cli 1.1.0');
			const { data } = await ask();
			expect(data.pushed).toBe(false);
			expect(data.unpushed).toMatchObject({ acceptances: ['S-004'], tags: ['cli/v1.1.0'] });
			// a dry run pushes nothing
			expect(git(join(base, 'origin.git'), 'tag', '--list').trim()).toBe('');
		} finally {
			await rm(base, { recursive: true, force: true });
		}
	});
	it('blocks, unblocks, and logs to a stream', async () => {
		await flai(dir, ['block', 'S-004', '--reason', 'waiting']);
		let file = await readFile(join(dir, 'wip/kanban/stories/S-004-four.md'), 'utf8');
		expect(file).toContain('reason: waiting');
		await flai(dir, ['unblock', 'S-004']);
		file = await readFile(join(dir, 'wip/kanban/stories/S-004-four.md'), 'utf8');
		expect(file).toMatch(/until: \d{4}/);
		await flai(dir, ['stream', 'log', 'S-004', 'from the dashboard']);
		const narrative = await readFile(join(dir, 'wip/agents/S-004.md'), 'utf8');
		expect(narrative).toContain('from the dashboard');
	});
	it('creates a story with a body in one step, and is refused by the check with nothing left (S-0059)', async () => {
		const { data: body } = await flai<{ body: string }>(dir, ['story', 'new', '--print-body']);
		expect(body.body).toContain('## Goal');
		expect(body.body).not.toContain('---');
		const args = (title: string) => [
			'story',
			'new',
			'--nature=improvement',
			'--owner=olive',
			'--epic=E-001',
			'--body-stdin',
			'--autocommit',
			'--',
			title
		];
		const { data } = await flai<{ item: { id: string; owner: string }; path: string }>(
			dir,
			args('--json is a title here'),
			{
				input:
					'## Goal\nFrom the board.\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n',
				exitStatus: { 4: 422 }
			}
		);
		expect(data.item.owner).toBe('olive');
		const file = await readFile(join(dir, data.path), 'utf8');
		expect(file).toContain(`# ${data.item.id} --json is a title here\n\n## Goal\nFrom the board.`);
		const before = await readdir(join(dir, 'wip/kanban/stories'));
		await expect(
			flai(dir, args('Refused'), {
				input:
					'## Goal\nNo notes section, which flai check reports.\n\n## Acceptance criteria\n- [ ] x\n\n## Tasks\n',
				exitStatus: { 4: 422 }
			})
		).rejects.toMatchObject({ status: 422, data: { refused: { findings: expect.any(Array) } } });
		expect(await readdir(join(dir, 'wip/kanban/stories'))).toEqual(before);
	});
	it('reports a missing binary as 503', async () => {
		process.env.FLAI_BIN = '/nonexistent/flai';
		const savedPath = process.env.PATH;
		process.env.PATH = '';
		resetFlaiBinary();
		await expect(flai(dir, ['board'])).rejects.toBeInstanceOf(RepoError);
		await expect(flai(dir, ['board'])).rejects.toMatchObject({ status: 503 });
		process.env.PATH = savedPath;
		process.env.FLAI_BIN = bin;
		resetFlaiBinary();
	});
});
