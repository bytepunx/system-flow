import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo } from './repo';
import { flaiAsk } from './testing';
import { flai, resetFlaiBinary } from './flai';
import { activity } from './activity';
import { hrefFor, inbox } from './inbox';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');

// How narratives are parsed is flai's now, and tested there (internal/hostapi); the dashboard adds
// the link each entry leads to.
describe('where an inbox entry leads', () => {
	it('sends a story in review to its review page, an item to its page, a document to the docs', () => {
		expect(hrefFor({ kind: 'review', item: 'S-0004' })).toBe('/review/S-0004');
		expect(hrefFor({ kind: 'blocked', item: 'T-0001' })).toBe('/items/T-0001');
		expect(
			hrefFor({ kind: 'thread', item: 'S-0001', path: 'wip/kanban/stories/S-0001-a.md' })
		).toBe('/items/S-0001');
		expect(hrefFor({ kind: 'thread', path: 'design/system/overview.md' })).toBe(
			'/docs/design/system/overview.md'
		);
		expect(hrefFor({ kind: 'question', path: 'wip/agents/S-0001.md' })).toBe(
			'/docs/wip/agents/S-0001.md'
		);
		expect(hrefFor({ kind: 'overlap' })).toBe('/board');
	});
});

describe.skipIf(!existsSync(bin))('activity and inbox on a project', () => {
	let dir: string;
	let repo: Repo;
	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-inbox-'));
		await cp(fixture, dir, { recursive: true });
		const manifest = join(dir, 'system-flow.yaml');
		await writeFile(
			manifest,
			(await readFile(manifest, 'utf8')).replace(/^owner:.*\n/m, '') + 'owner: alex\n'
		);
		const narrative = join(dir, 'wip/agents/S-004.md');
		await writeFile(
			narrative,
			(await readFile(narrative, 'utf8'))
				.replace('## Open questions\n', '## Open questions\n- Which port should it use?\n')
				.replace(
					'Stream opened.\n',
					'Stream opened.\n\n### 2026-08-31T12:00:00Z\nT-003 half done.\n'
				)
		);
		process.env.PROJECT_DIR = dir;
		process.env.FLAI_BIN = bin;
		resetFlaiBinary();
		// a second story in progress that touches what S-004 touches, and a block
		await flai(dir, ['touches', 'S-004', 'flai/cmd']);
		const { data } = await flai<{ id: string }>(dir, [
			'story',
			'new',
			'Overlapping',
			'--epic',
			'E-001',
			'--touches',
			'flai/cmd/move.go'
		]);
		const file = join(dir, 'wip/kanban/stories');
		const { readdir } = await import('node:fs/promises');
		const name = (await readdir(file)).find((f) => f.startsWith(data.id + '-'))!;
		const body = await readFile(join(file, name), 'utf8');
		await writeFile(join(file, name), body.replace('- [ ]\n', '- [x] ok\n'));
		await flai(dir, ['move', data.id, 'ready', '--by', 'alex']);
		await flai(dir, ['move', data.id, 'in-progress', '--by', 'alex']);
		await flai(dir, ['block', 'T-003', '--reason', 'waiting on the designer']);
		repo = new Repo(dir, flaiAsk(dir));
	});
	afterAll(async () => {
		await repo?.close();
		await rm(dir, { recursive: true, force: true });
	});

	it('shows each stream with its agent, task in progress, blocked flag, and last log entry', async () => {
		const { streams } = await activity(repo, new Date('2026-09-01T12:00:00Z'));
		expect(streams).toHaveLength(1);
		expect(streams[0]).toMatchObject({
			stream: 'S-004',
			title: 'Four',
			agent: 'bot',
			status: 'in-progress',
			blocked: true,
			task: { id: 'T-003' },
			last_log: { at: '2026-08-31T12:00:00Z', text: 'T-003 half done.' },
			path: 'wip/agents/S-004.md'
		});
		expect(streams[0].age_seconds).toBe(26 * 3600);
	});

	it('lists what needs a human, each with a stable key and a link', async () => {
		const box = await inbox(repo);
		const byKind = (k: string) => box.entries.filter((e) => e.kind === k);
		expect(byKind('question')).toMatchObject([
			{ title: 'Which port should it use?', href: '/docs/wip/agents/S-004.md' }
		]);
		expect(byKind('blocked')).toMatchObject([
			{
				title: expect.stringContaining('T-003'),
				detail: 'blocked: waiting on the designer',
				href: '/items/T-003'
			}
		]);
		expect(byKind('overlap')).toHaveLength(1);
		expect(byKind('overlap')[0].title).toContain('flai/cmd');
		expect(byKind('overlap')[0].href).toMatch(/^\/items\/S-/);
		expect(box.total).toBe(box.entries.length);
		expect(box.counts.question + box.counts.blocked + box.counts.overlap).toBeGreaterThanOrEqual(3);
		expect(box.notes).toEqual([]);
		// the same repository gives the same keys, so "new" means something
		repo.forget();
		expect((await inbox(repo)).entries.map((e) => e.key)).toEqual(box.entries.map((e) => e.key));
	});

	it('counts a thread only while the last word is not the designer’s', async () => {
		const threads = (await inbox(repo)).entries.filter((e) => e.kind === 'thread');
		const all = await repo.threads();
		const expected = all.filter(
			(t) => t.status !== 'resolved' && t.entries.at(-1) && t.entries.at(-1)!.author !== 'alex'
		);
		expect(threads.map((t) => t.key)).toEqual(expected.map((t) => `thread:${t.id}`));
	});

	it('lists a story in review with a link to its review page', async () => {
		await flai(dir, ['move', 'T-003', 'done', '--by', 'bot']).catch(() => undefined);
		await flai(dir, ['unblock', 'T-003']).catch(() => undefined);
		await flai(dir, ['move', 'S-004', 'review', '--by', 'bot']);
		const fresh = new Repo(dir, flaiAsk(dir));
		const review = (await inbox(fresh)).entries.filter((e) => e.kind === 'review');
		await fresh.close();
		expect(review).toMatchObject([{ key: 'review:S-004', href: '/review/S-004' }]);
	});
});
