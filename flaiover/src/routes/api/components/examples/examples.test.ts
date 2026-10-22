import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import {
	Repo,
	useRepo,
	type Ask,
	type ComponentAnswer,
	type ComponentExample
} from '$lib/server/repo';

// /api/components/examples adds an example to a component's examples file, or replaces the one of
// its name, through flai's file.get and file.save (S-0443, ADR-0164), against a stubbed host.
type Handler = (event: never) => Promise<Response>;

const CARD = 'src/lib/components/Card.svelte';
const FILE = 'design/components/web/src/lib/components/Card.examples.json';
const HASH = 'a'.repeat(64);
const SAVED = 'b'.repeat(64);
const LONG: ComponentExample = { name: 'Long title', props: { title: 'A long one' } };
// ADR-0187's answers: a fetch's and an EventSource's
const ANSWERS: ComponentAnswer[] = [
	{ path: '/api/threads?item=S-0001', body: { threads: [] } },
	{ path: '/api/events', events: [{ event: 'change', data: { kind: 'edited' } }] }
];

const catalogue = (examples: ComponentExample[] | undefined, path: string | null = FILE) => ({
	projects: [
		{
			name: 'web',
			path: 'web',
			kind: 'sveltekit',
			supported: true,
			components: [
				{
					project: 'web',
					path: CARD,
					name: 'Card',
					props: [{ name: 'title', kind: 'string', required: true }],
					...(examples ? { examples } : {}),
					...(path ? { examples_path: path } : {})
				}
			]
		}
	]
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
// each method's answers, in turn; the last one answers again
let script: Record<string, { data?: unknown; error?: unknown }[]> = {};
// components.list and project.info answer raw, as flai's reads do; file.get and file.save as commands do
const RAW = new Set(['components.list', 'project.info']);
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const answers = script[method];
	if (!answers?.length) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	const s = answers.length > 1 ? answers.shift()! : answers[0];
	if (s.error) throw s.error;
	return RAW.has(method) ? s.data : { data: s.data, warnings: ['saved'] };
}) as Ask;

const conflict = (current: string, hash: string) =>
	new AgentError(502, `${FILE} changed after it was loaded`, -32009, {
		current,
		hash,
		diff: '+mine'
	});

describe('/api/components/examples (S-0443)', () => {
	let PUT: Handler;
	const put = (body: unknown) =>
		PUT({
			request: new Request('http://x/api/components/examples', {
				method: 'PUT',
				body: JSON.stringify(body)
			})
		} as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);
	const saved = (i = 0) => asks('file.save')[i].params;
	const body = (example: unknown = LONG, more: Record<string, unknown> = {}) => ({
		project: 'web',
		component: CARD,
		example,
		...more
	});
	const serve = (c: ReturnType<typeof catalogue>) => (script['components.list'] = [{ data: c }]);

	beforeAll(async () => {
		PUT = (await import('./+server')).PUT as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': [{ data: { name: 'p', key: 'p' } }],
			'file.save': [{ data: { path: FILE, hash: SAVED, committed: true } }]
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('creates the file a first example goes in, with an empty hash, without reading it', async () => {
		serve(catalogue([]));
		const r = await put(body({ ...LONG, theme: 'dark', width: 375 }));
		expect(r.status).toBe(200);
		const example = { ...LONG, theme: 'dark', width: 375 };
		expect(await r.json()).toMatchObject({
			path: FILE,
			hash: SAVED,
			committed: true,
			examples: [example],
			log: ['saved']
		});
		expect(asks('file.get')).toEqual([]);
		expect(saved()).toMatchObject({ path: FILE, hash: '' });
		expect(saved().story).toBeUndefined();
		expect(saved().content).toBe(
			'{\n  "examples": [\n    {\n      "name": "Long title",\n      "props": {\n        "title": "A long one"\n      },\n      "theme": "dark",\n      "width": 375\n    }\n  ]\n}\n'
		);
	});

	it('adds an example to the file, read afresh with its hash, keeping the ones it holds', async () => {
		const short = { name: 'Short', props: { title: 'Hi' }, extra: 'kept as written' };
		serve(catalogue([{ name: 'Short', props: { title: 'Hi' } }]));
		script['file.get'] = [
			{ data: { path: FILE, content: JSON.stringify({ examples: [short] }), hash: HASH } }
		];
		const r = await put(body(LONG, { story: 'S-0443' }));
		expect(r.status).toBe(200);
		expect((await r.json()).examples).toEqual([short, LONG]);
		expect(asks('components.list')[0].params).toEqual({ story: 'S-0443' });
		expect(asks('file.get')[0].params).toEqual({ path: FILE, story: 'S-0443' });
		expect(saved()).toMatchObject({ path: FILE, hash: HASH, story: 'S-0443' });
		expect(JSON.parse(saved().content as string)).toEqual({ examples: [short, LONG] });
	});

	it('replaces the example of the same name in its place', async () => {
		const first = { name: 'First', props: {} };
		serve(catalogue([first, LONG]));
		const content = JSON.stringify({ examples: [first, LONG, { name: 'Last', props: {} }] });
		script['file.get'] = [{ data: { path: FILE, content, hash: HASH } }];
		const changed = { name: 'Long title', props: { title: 'Longer still' }, width: 768 };
		await put(body(changed));
		expect(JSON.parse(saved().content as string).examples).toEqual([
			first,
			changed,
			{ name: 'Last', props: {} }
		]);
	});

	it('keeps the answers a saved example carries, after its other keys (S-0476)', async () => {
		serve(catalogue([]));
		const answered = { name: 'Answered', props: {}, answers: ANSWERS, width: 375 };
		const r = await put(body(answered));
		expect(r.status).toBe(200);
		const example = { name: 'Answered', props: {}, width: 375, answers: ANSWERS };
		expect((await r.json()).examples).toEqual([example]);
		expect(saved().content).toBe(JSON.stringify({ examples: [example] }, null, 2) + '\n');
	});

	it('writes an example without answers, or with none, without the key (S-0476)', async () => {
		serve(catalogue([]));
		await put(body(LONG));
		await put(body({ ...LONG, answers: [] }));
		for (const i of [0, 1]) {
			expect(JSON.parse(saved(i).content as string)).toEqual({ examples: [LONG] });
			expect(saved(i).content).not.toContain('answers');
		}
	});

	it("keeps another example's answers when one is saved (S-0476)", async () => {
		const answered = { name: 'Answered', props: {}, answers: ANSWERS };
		serve(catalogue([answered, LONG]));
		const content = JSON.stringify({ examples: [answered, LONG] });
		script['file.get'] = [{ data: { path: FILE, content, hash: HASH } }];
		const changed = { name: 'Long title', props: { title: 'Longer still' } };
		const r = await put(body(changed));
		expect(r.status).toBe(200);
		expect((await r.json()).examples).toEqual([answered, changed]);
		expect(JSON.parse(saved().content as string).examples).toEqual([answered, changed]);
	});

	it('adds to a file that holds no examples once the create finds it, with its hash', async () => {
		serve(catalogue([]));
		script['file.save'] = [
			{ error: conflict('{\n  "examples": []\n}\n', HASH) },
			{ data: { path: FILE, hash: SAVED, committed: true } }
		];
		const r = await put(body());
		expect(r.status).toBe(200);
		expect(asks('file.save')).toHaveLength(2);
		expect(saved(0).hash).toBe('');
		expect(saved(1).hash).toBe(HASH);
		expect(JSON.parse(saved(1).content as string)).toEqual({ examples: [LONG] });
	});

	it('answers a save that loses a race with 409, the current content, its hash, and a diff', async () => {
		serve(catalogue([LONG]));
		const current = JSON.stringify({ examples: [LONG, { name: 'Theirs', props: {} }] });
		script['file.get'] = [
			{ data: { path: FILE, content: JSON.stringify({ examples: [LONG] }), hash: HASH } }
		];
		script['file.save'] = [{ error: conflict(current, 'c'.repeat(64)) }];
		const r = await put(body({ name: 'Mine', props: {} }));
		expect(r.status).toBe(409);
		expect(await r.json()).toMatchObject({
			error: `${FILE} changed after it was loaded`,
			current,
			hash: 'c'.repeat(64),
			diff: '+mine'
		});
		expect(asks('file.save')).toHaveLength(1);
	});

	it("answers flai's refusal with 422 and its findings", async () => {
		serve(catalogue([]));
		const findings = [
			{ level: 'warning', rule: 'components.examples', path: FILE, line: 6, message: 'no prop x' }
		];
		script['file.save'] = [
			{ error: new AgentError(502, 'flai check refused the save', -32010, { findings }) }
		];
		const r = await put(body({ name: 'Bad', props: { x: 1 } }));
		expect(r.status).toBe(422);
		expect(await r.json()).toMatchObject({ error: 'flai check refused the save', findings });
	});

	it('answers 422 when the file is not an examples file, before saving', async () => {
		serve(catalogue([LONG]));
		script['file.get'] = [{ data: { path: FILE, content: '{ "examples": ', hash: HASH } }];
		const r = await put(body());
		expect(r.status).toBe(422);
		expect((await r.json()).error).toContain(`${FILE} is not JSON`);
		expect(asks('file.save')).toEqual([]);
	});

	it('answers 404 for a component the catalogue lacks, and 501 for a host too old to say where', async () => {
		serve(catalogue([]));
		const r = await put({ ...body(), component: 'src/lib/Nope.svelte' });
		expect(r.status).toBe(404);
		expect((await r.json()).error).toContain('src/lib/Nope.svelte');
		serve(catalogue(undefined, null));
		const old = await put(body());
		expect(old.status).toBe(501);
		expect((await old.json()).error).toContain('upgrade');
		expect(asks('file.save')).toEqual([]);
	});

	it('refuses a malformed body with 400, before asking flai', async () => {
		serve(catalogue([]));
		const bodies = [
			{ component: CARD, example: LONG },
			{ project: 'web', example: LONG },
			body(null),
			body({ name: '  ', props: {} }),
			body({ name: 'x', props: [] }),
			body({ name: 'x', props: {}, theme: 'blue' }),
			body({ name: 'x', props: {}, width: 0 }),
			body({ name: 'x', props: {}, width: 1.5 }),
			body({ name: 'x', props: {}, answers: ANSWERS[0] }),
			body(LONG, { story: 'main' })
		];
		for (const b of bodies) expect((await put(b)).status, JSON.stringify(b)).toBe(400);
		expect(asks('components.list')).toEqual([]);
		expect(asks('file.save')).toEqual([]);
	});
});
