import { describe, expect, it, vi } from 'vitest';
import { Repo, type Ask } from './repo';
import { activity } from './activity';
import { inbox } from './inbox';
import { adrs } from './search';
import { stats } from './stats';

// S-0161: a change forgets only the answers read from its kind of file.
const project = {
	name: 'Harbour',
	key: 'harbour',
	layout: { design: 'design', docs: 'docs', wip: 'wip' }
};
const answers: Record<string, unknown> = {
	'project.info': project,
	'board.get': { columns: {} },
	'items.list': [],
	'items.count': { active: 0, archived: 0 },
	'item.get': {
		item: { id: 'S-0001', transitions: null, blocked: null, tags: null },
		children: null
	},
	'threads.list': [],
	'docs.tree': [],
	'doc.get': { path: 'docs/users/guide.md', frontMatter: null, body: '', raw: '' },
	'inbox.designer': { total: 0, counts: {}, entries: null, notes: null },
	'activity.get': { streams: null },
	'stats.get': { data: {} },
	'adrs.list': []
};

/** Every answer the Repo keeps, each once, and a count of the asks per method. */
async function held() {
	const asked: string[] = [];
	const ask = vi.fn(async (method: string) => {
		asked.push(method);
		return answers[method];
	}) as unknown as Ask;
	const r = new Repo('/nowhere', ask);
	const all = () =>
		Promise.all([
			r.manifest(),
			r.boardView(),
			r.items(),
			r.itemCount(),
			r.itemById('S-0001'),
			r.threads(),
			r.threadsFor('S-0001'),
			r.docsTree(),
			r.docFile('docs/users/guide.md'),
			inbox(r),
			activity(r),
			stats(r, { since: '30d' }),
			adrs(r)
		]);
	await all();
	asked.length = 0;
	/** The methods asked again after a change to path, sorted. */
	const after = async (path: string) => {
		r.changed(path);
		await all();
		const again = [...asked].sort();
		asked.length = 0;
		return again;
	};
	return { r, after, all, asked };
}

describe('Repo.changed forgets what the changed file affects (S-0161)', () => {
	it('keeps the board and the items when a narrative changes', async () => {
		const { after } = await held();
		expect(await after('wip/agents/S-0161.md')).toEqual([
			'activity.get',
			'docs.tree',
			'inbox.designer'
		]);
	});

	it('keeps the board when a thread changes; the threads, the inbox, and the tree are asked again', async () => {
		const { after } = await held();
		expect(await after('wip/threads/TH-0001-berths.md')).toEqual([
			'docs.tree',
			'inbox.designer',
			'threads.list',
			'threads.list'
		]);
	});

	it('keeps everything but the tree when a design document changes', async () => {
		const { after } = await held();
		expect(await after('design/system/overview.md')).toEqual(['docs.tree']);
	});

	it('forgets a document it holds when that document changes', async () => {
		const { after } = await held();
		expect(await after('docs/users/guide.md')).toEqual(['doc.get', 'docs.tree']);
	});

	it('forgets the ADR list only for an ADR', async () => {
		const { after } = await held();
		expect(await after('design/adrs/0052-x.md')).toEqual(['adrs.list', 'docs.tree']);
	});

	it('forgets every answer read from the items when a work item changes, and keeps the rest', async () => {
		const { after } = await held();
		expect(await after('wip/kanban/stories/S-0001-berths.md')).toEqual([
			'activity.get',
			'board.get',
			'docs.tree',
			'inbox.designer',
			'item.get',
			'items.count',
			'items.list',
			'stats.get',
			'threads.list',
			'threads.list'
		]);
	});

	it("treats the board's policy and the archive as work items", async () => {
		const { after } = await held();
		const items = await after('wip/kanban/board.md');
		expect(items).toContain('board.get');
		expect(items).not.toContain('adrs.list');
		expect(await after('wip/archive/kanban/stories/S-0002-quay.md')).toEqual(items);
	});

	it('forgets everything when the manifest or a file outside the three folders changes', async () => {
		const { after } = await held();
		const every = Object.keys(answers).concat('threads.list').sort();
		expect(await after('system-flow.yaml')).toEqual(every);
		expect(await after('scripts/build.sh')).toEqual(every);
	});

	it('forgets everything while the layout is not known, and learns it from the manifest', async () => {
		const asked: string[] = [];
		const ask = vi.fn(async (method: string) => {
			asked.push(method);
			return answers[method];
		}) as unknown as Ask;
		const r = new Repo('/nowhere', ask);
		await r.boardView();
		const kinds: string[] = [];
		r.on('change', (_p: string, kind: string) => kinds.push(kind));
		r.changed('wip/agents/S-0161.md');
		await r.boardView();
		expect(asked).toEqual(['board.get', 'project.info', 'board.get']);
		r.changed('wip/agents/S-0161.md');
		await r.boardView();
		expect(asked).toHaveLength(3);
		expect(kinds).toEqual(['other', 'narrative']);
	});

	it('announces the path with its kind', async () => {
		const { r } = await held();
		const seen: [string, string][] = [];
		r.on('change', (p: string, kind: string) => seen.push([p, kind]));
		r.changed('wip/threads/TH-0001-berths.md');
		expect(seen).toEqual([['wip/threads/TH-0001-berths.md', 'thread']]);
	});
});

// S-0225: item.get's expected cost reaches the item page, and its absence adds no key.
describe('itemById', () => {
	it('passes on the expected cost flai priced', async () => {
		const expected_cost = { cost: 24, from: 'estimate', estimated: true };
		const ask = vi.fn(async () => ({
			...(answers['item.get'] as object),
			expected_cost
		})) as unknown as Ask;
		expect((await new Repo('/nowhere', ask).itemById('S-0001')).expected_cost).toEqual(
			expected_cost
		);
		const without = vi.fn(async () => answers['item.get']) as unknown as Ask;
		expect(await new Repo('/nowhere', without).itemById('S-0001')).not.toHaveProperty(
			'expected_cost'
		);
	});
});
