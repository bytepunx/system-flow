import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, readdir, readFile, rm, writeFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo, RepoError } from './repo';
import { board } from './board';
import { flai, flaiBinary, resetFlaiBinary } from './flai';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');
const haveFlai = existsSync(bin);

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
