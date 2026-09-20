import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, readdir, readFile, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo } from './repo';
import { flaiAsk, haveFlai, shell } from './testing';
import { board } from './board';

// The dashboard's writes are flai's methods on the host (S-0075). Here they run through `flai
// hostapi` from this tree: the same method table flai serve offers, the same validation, the same
// commands, against a copy of the fixture.
const fixture = resolve('../flai/internal/metrics/testdata/good');

describe('board reader', () => {
	it('groups active items by state with ages and limits', async () => {
		const b = await board(new Repo(fixture, flaiAsk(fixture)), new Date('2026-09-01T12:00:00Z'));
		expect(b.wip_limits).toEqual({ ready: 5, 'in-progress': 2, review: 3 });
		expect(b.columns['in-progress'].map((c) => c.id)).toEqual(['S-004', 'T-003']);
		const s4 = b.columns['in-progress'][0];
		expect(s4.age_seconds).toBe(93600);
		expect(s4.blocked).toBe(false);
		expect(b.columns.done).toEqual([]); // done items in the fixture are archived
	});
	it('names the parent of every card that has one', async () => {
		const b = await board(new Repo(fixture, flaiAsk(fixture)), new Date('2026-09-01T12:00:00Z'));
		const [story, task] = b.columns['in-progress'];
		expect(story).toMatchObject({ id: 'S-004', parent: 'E-001', parent_title: 'Epic' });
		expect(task).toMatchObject({ id: 'T-003', parent: 'S-004', parent_title: 'Four' });
		const epic = b.columns.backlog.find((c) => c.id === 'E-001')!;
		expect(epic.parent).toBeUndefined();
		expect(epic.parent_title).toBeUndefined();
	});
});

describe.skipIf(!haveFlai)('writes through flai on a temp project', () => {
	let dir: string;
	let r: Repo;
	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-'));
		await cp(fixture, dir, { recursive: true });
		r = new Repo(dir, flaiAsk(dir));
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('performs a valid move as the designer and reports the new status', async () => {
		const { data } = await r.write<{ id: string; status: string }>('item.move', {
			id: 'T-003',
			to: 'done'
		});
		expect(data).toMatchObject({ id: 'T-003', status: 'done' });
		const file = await readFile(join(dir, 'wip/kanban/tasks/T-003-t3.md'), 'utf8');
		expect(file).toContain('status: done');
		// the fixture names no owner, so flai records the designer; the dashboard names nobody
		expect(file).toMatch(
			/- to: done\n\s+at: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\n\s+by: designer/
		);
	});
	it('refuses an invalid move with the rule text as a 400, and an argument that is not one', async () => {
		await expect(r.write('item.move', { id: 'S-004', to: 'done' })).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('cannot go from in-progress to done')
		});
		await expect(r.write('item.move', { id: 'S-004', to: 'cancelled' })).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('needs --reason')
		});
		await expect(r.write('item.move', { id: '--help', to: 'ready' })).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('not an item ID')
		});
	});
	it('moves a story with no tasks to ready and in-progress, and is refused at review (ADR-0021)', async () => {
		const made = await shell<{ id: string }>(dir, [
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
			const { data } = await r.write<{ status: string }>('item.move', { id: made.id, to });
			expect(data.status).toBe(to);
		}
		await expect(r.write('item.move', { id: made.id, to: 'review' })).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('needs at least one task before it goes to review')
		});
	});
	it('places a story in the pull order, and the board reads the column in that order (S-0057)', async () => {
		const ids: string[] = [];
		for (const title of ['Order one', 'Order two', 'Order three'])
			ids.push((await shell<{ id: string }>(dir, ['story', 'new', title, '--epic', 'E-001'])).id);
		const [one, two, three] = ids;
		const backlog = async () =>
			(await board(new Repo(dir, flaiAsk(dir)))).columns.backlog
				.filter((c) => ids.includes(c.id))
				.map((c) => c.id);
		expect(await backlog()).toEqual([one, two, three]);
		const { data } = await r.write<{ status: string; sequence: string[] }>('item.order', {
			id: three,
			before: one
		});
		expect(data.status).toBe('backlog');
		expect(data.sequence.filter((id) => ids.includes(id))).toEqual([three, one, two]);
		expect(await backlog()).toEqual([three, one, two]);
		// an epic sharing the column keeps its place
		const column = (await board(new Repo(dir, flaiAsk(dir)))).columns.backlog;
		expect(column.findIndex((c) => c.id === 'E-001')).toBe(
			column.findIndex((c) => c.type !== 'story')
		);
		await expect(r.write('item.order', { id: three, before: 'S-004' })).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('within one column')
		});
		await expect(r.write('item.order', { id: 'E-001', top: true })).rejects.toMatchObject({
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
				new Repo(clone, flaiAsk(clone)).run<{
					pushed: boolean;
					reason?: string;
					unpushed?: { acceptances: string[]; tags: string[] };
				}>('push.pending');
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
	it('is refused a push while the operator has not enabled it, and told what enables it (S-0078)', async () => {
		// a configuration of its own: the refusal is journalled beside it, and that is not the operator's
		const was = process.env.FLAI_CONFIG;
		const host = await mkdtemp(join(tmpdir(), 'flaiover-host-'));
		process.env.FLAI_CONFIG = join(host, 'config.json');
		try {
			const info = await r.ask<{ host_actions: Record<string, boolean> }>('project.info');
			expect(info.host_actions).toEqual({ agent: false, push: false });
			await expect(r.write('push.run')).rejects.toMatchObject({
				status: 403,
				message: expect.stringContaining('flai serve enable push'),
				data: { action: 'push', enable: 'flai serve enable push' }
			});
			const journal = await readFile(join(host, 'serve', 'journal.jsonl'), 'utf8');
			expect(JSON.parse(journal.trim())).toMatchObject({
				action: 'push',
				method: 'push.run',
				outcome: 'disabled'
			});
		} finally {
			if (was === undefined) delete process.env.FLAI_CONFIG;
			else process.env.FLAI_CONFIG = was;
			await rm(host, { recursive: true, force: true });
		}
	});
	it('blocks, unblocks, and logs to a stream', async () => {
		await r.write('item.block', { id: 'S-004', reason: 'waiting' });
		let file = await readFile(join(dir, 'wip/kanban/stories/S-004-four.md'), 'utf8');
		expect(file).toContain('reason: waiting');
		await r.write('item.unblock', { id: 'S-004' });
		file = await readFile(join(dir, 'wip/kanban/stories/S-004-four.md'), 'utf8');
		expect(file).toMatch(/until: \d{4}/);
		await r.write('stream.log', { id: 'S-004', entry: '--from the dashboard' });
		const narrative = await readFile(join(dir, 'wip/agents/S-004.md'), 'utf8');
		expect(narrative).toContain('--from the dashboard');
	});
	it('creates a story with a body in one step, and is refused by the check with nothing left (S-0059)', async () => {
		const { data: template } = await r.run<{ body: string }>('item.template', { type: 'story' });
		expect(template.body).toContain('## Goal');
		expect(template.body).not.toContain('---');
		const make = (title: string, body: string) =>
			r.write<{ item: { id: string; owner: string }; path: string }>('item.new', {
				type: 'story',
				title,
				nature: 'improvement',
				parent: 'E-001',
				body
			});
		const { data } = await make(
			'--json is a title here',
			'## Goal\nFrom the board.\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n'
		);
		expect(data.item.owner).toBe('designer');
		const file = await readFile(join(dir, data.path), 'utf8');
		expect(file).toContain(`# ${data.item.id} --json is a title here\n\n## Goal\nFrom the board.`);
		const before = await readdir(join(dir, 'wip/kanban/stories'));
		await expect(
			make(
				'Refused',
				'## Goal\nNo notes section, which flai check reports.\n\n## Acceptance criteria\n- [ ] x\n\n## Tasks\n'
			)
		).rejects.toMatchObject({ status: 422, data: { findings: expect.any(Array) } });
		expect(await readdir(join(dir, 'wip/kanban/stories'))).toEqual(before);
	});
	it('reads what changes nothing without a request ID: a cancellation preview, a diff of nothing, stats', async () => {
		const { data } = await r.run<{ dry_run: boolean; cancelled: unknown[] }>('item.move.preview', {
			id: 'E-001'
		});
		expect(data.dry_run).toBe(true);
		expect(Array.isArray(data.cancelled)).toBe(true);
	});
});

// S-0085: an item's own words are changed by flai on the host, in one checked step.
describe.skipIf(!haveFlai)('editing an item through flai', () => {
	let dir: string;
	let r: Repo;
	type View = { title: string; body: string; hash: string; editable: boolean; parents: unknown[] };
	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-edit-'));
		await cp(fixture, dir, { recursive: true });
		r = new Repo(dir, flaiAsk(dir));
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('shows what an editor needs, retitles everywhere, and refuses a stale hash', async () => {
		const { data: v } = await r.run<View>('item.show', { id: 'S-004' });
		expect(v).toMatchObject({ title: 'Four', editable: true });
		expect(v.body).not.toMatch(/^# /);
		expect(v.hash).toMatch(/^[0-9a-f]{64}$/);
		expect(v.parents.length).toBeGreaterThan(0);

		const { data } = await r.write<{ changed: string[]; renamed_from?: string; path: string }>(
			'item.edit',
			{ id: 'S-004', hash: v.hash, title: 'Four, renamed', tags: ['dashboard'] }
		);
		expect(data.changed).toEqual(['title', 'tags']);
		expect(data.renamed_from).toBe('wip/kanban/stories/S-004-four.md');
		const file = await readFile(join(dir, data.path), 'utf8');
		expect(file).toContain('title: Four, renamed');
		expect(file).toContain('# S-004 Four, renamed');
		expect(file).toContain('tags: [dashboard]');

		// the hash is of the file as it was read: it no longer is
		await expect(
			r.write('item.edit', { id: 'S-004', hash: v.hash, title: 'Too late' })
		).rejects.toMatchObject({
			status: 409,
			data: { hash: expect.stringMatching(/^[0-9a-f]{64}$/) }
		});
	});

	it('is refused by the check with nothing changed, and by flai for what is not a value', async () => {
		const { data: v } = await r.run<View>('item.show', { id: 'S-004' });
		await expect(
			r.write('item.edit', { id: 'S-004', hash: v.hash, body: '## Goal\nNo criteria left.\n' })
		).rejects.toMatchObject({ status: 422, data: { findings: expect.any(Array) } });
		expect((await r.run<View>('item.show', { id: 'S-004' })).data.hash).toBe(v.hash);
		await expect(
			r.write('item.edit', { id: 'S-004', hash: v.hash, tags: ['--by=eve'] })
		).rejects.toMatchObject({ status: 400 });
		await expect(r.write('item.edit', { id: 'S-004', hash: v.hash })).rejects.toMatchObject({
			status: 400
		});
		// what flai itself refuses as not allowed is the caller's mistake too: 400, never 500
		const { data: epic } = await r.run<View>('item.show', { id: 'E-001' });
		await expect(
			r.write('item.edit', { id: 'E-001', hash: epic.hash, touches: ['flai'] })
		).rejects.toMatchObject({
			status: 400,
			message: expect.stringContaining('touches belong to stories and tasks')
		});
		await expect(
			r.write('item.edit', { id: 'S-004', hash: v.hash, parent: 'S-004' })
		).rejects.toMatchObject({ status: 400 });
		await expect(
			r.write('item.edit', { id: 'S-004', hash: v.hash, parent: 'E-999' })
		).rejects.toMatchObject({ status: 404 });
	});
});
