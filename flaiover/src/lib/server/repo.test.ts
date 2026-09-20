import { describe, expect, it } from 'vitest';
import { resolve } from 'node:path';
import { existsSync } from 'node:fs';
import { Repo, RepoError } from './repo';
import { flaiAsk } from './testing';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const monorepo = resolve('..');

describe('Repo on the metrics fixture', () => {
	const r = new Repo(fixture, flaiAsk(fixture));
	it('reads the manifest', async () => {
		const m = await r.manifest();
		expect(m.name).toBe('good');
		expect(m.layout).toEqual({ design: 'design', docs: 'docs', wip: 'wip' });
	});
	it('lists items across kanban and archive, sorted, with derived fields', async () => {
		const items = await r.items();
		expect(items.map((i) => i.id)).toEqual([
			'E-001',
			'S-001',
			'S-002',
			'S-003',
			'S-004',
			'S-005',
			'T-001',
			'T-002',
			'T-003',
			'T-004'
		]);
		const s1 = items.find((i) => i.id === 'S-001')!;
		expect(s1.archived).toBe(true);
		expect(s1.transitions.at(-1)).toEqual({ to: 'done', at: '2026-08-03T12:00:00Z', by: 'alex' });
		expect(s1.blocked?.[0].until).toBe('2026-08-02T18:00:00Z');
		expect(s1.estimate).toBe('20h');
		for (const it of items) {
			const last = it.transitions.at(-1)?.to ?? 'backlog';
			expect(it.status, it.id).toBe(last);
		}
	});
	it('reads threads and finds them by path or item in any padding', async () => {
		const all = await r.threads();
		expect(all.map((t) => t.id)).toEqual(['TH-0001']);
		expect(all[0].entries).toHaveLength(2);
		expect(all[0].entries[1]).toMatchObject({ at: '2026-08-30T09:30:00Z', author: 'agent' });
		expect(all[0].entries[1].text).toContain('Ticked now');
		expect((await r.threadsFor('S-4')).map((t) => t.id)).toEqual(['TH-0001']);
		expect((await r.threadsFor('wip/kanban/stories/S-004-four.md')).map((t) => t.id)).toEqual([
			'TH-0001'
		]);
		expect(await r.threadsFor('S-0099')).toEqual([]);
	});
	it('returns an item with its children', async () => {
		const { item, children } = await r.itemById('E-001');
		expect(item.type).toBe('epic');
		expect(children.map((c) => c.id)).toEqual(['S-001', 'S-002', 'S-003', 'S-004', 'S-005']);
		await expect(r.itemById('S-999')).rejects.toBeInstanceOf(RepoError);
		expect((await r.itemById('e-1')).item.id).toBe('E-001');
		expect((await r.itemById('E-0001')).item.id).toBe('E-001');
		await expect(r.itemById('X-001')).rejects.toBeInstanceOf(RepoError);
	});
	it('builds the documentation trees', async () => {
		const trees = await r.docsTree();
		expect(trees.map((t) => t.path)).toEqual(['design', 'docs', 'wip']);
		const design = trees[0];
		expect(design.children?.map((c) => c.name)).toContain('conventions');
		const conv = design.children!.find((c) => c.name === 'conventions')!;
		const git = conv.children!.find((c) => c.name === 'git.md')!;
		expect(git.title).toBe('Git');
	});
	it('serves one file and rejects escapes and non-markdown', async () => {
		const f = await r.docFile('design/conventions/git.md');
		expect(f.frontMatter?.order).toBe(20);
		expect(f.body).toContain('Commit at landing');
		await expect(r.docFile('../etc/passwd')).rejects.toMatchObject({ status: 400 });
		await expect(r.docFile('system-flow.yaml')).rejects.toMatchObject({ status: 400 });
		await expect(r.docFile('design/nope.md')).rejects.toMatchObject({ status: 404 });
	});
});

describe.skipIf(!existsSync(resolve(monorepo, 'system-flow.yaml')))('Repo on the monorepo', () => {
	it('reads every item with a consistent status', async () => {
		const r = new Repo(monorepo, flaiAsk(monorepo));
		const items = await r.items();
		expect(items.length).toBeGreaterThan(50);
		for (const it of items) {
			expect(it.status, it.path).toBe(it.transitions.at(-1)?.to ?? 'backlog');
		}
	});
});
