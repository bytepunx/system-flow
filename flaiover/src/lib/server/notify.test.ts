import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { cp, mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { createServer, type Server } from 'node:http';
import type { AddressInfo } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { Repo } from './repo';
import { resetFlaiBinary } from './flai';
import { resetInboxCache } from './inbox';
import { notifyUrl, startNotifier, type Notifier, type NotifyBody } from './notify';

const fixture = resolve('../flai/internal/metrics/testdata/good');

describe('inbox webhook', () => {
	let dir: string;
	let repo: Repo;
	let server: Server;
	let received: { body: NotifyBody; headers: Record<string, unknown> }[];
	let status: number;
	let notifier: Notifier | null;

	const setManifest = async (extra: string) => {
		const p = join(dir, 'system-flow.yaml');
		const base = (await readFile(p, 'utf8')).replace(/\ndashboard:[\s\S]*$/, '\n');
		await writeFile(p, base + extra);
	};
	const addQuestion = async (q: string) => {
		const p = join(dir, 'wip/agents/S-004.md');
		await writeFile(
			p,
			(await readFile(p, 'utf8')).replace('## Open questions\n', `## Open questions\n- ${q}\n`)
		);
	};
	const waitFor = async (ok: () => boolean, ms = 4000) => {
		const end = Date.now() + ms;
		while (!ok() && Date.now() < end) await new Promise((r) => setTimeout(r, 25));
	};

	beforeEach(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-notify-'));
		await cp(fixture, dir, { recursive: true });
		process.env.PROJECT_DIR = dir;
		// no flai here: overlaps are left out, which is not what this tests
		process.env.FLAI_BIN = join(dir, 'no-such-flai');
		process.env.PATH = '';
		resetFlaiBinary();
		resetInboxCache();
		received = [];
		status = 200;
		server = createServer((req, res) => {
			let raw = '';
			req.on('data', (d) => (raw += d));
			req.on('end', () => {
				received.push({ body: JSON.parse(raw), headers: req.headers });
				res.statusCode = status;
				res.end();
			});
		});
		await new Promise<void>((r) => server.listen(0, '127.0.0.1', r));
		repo = new Repo(dir);
		notifier = null;
	});
	afterEach(async () => {
		notifier?.stop();
		await repo.close();
		await new Promise((r) => server.close(r));
		await rm(dir, { recursive: true, force: true });
	});
	const url = () => `http://127.0.0.1:${(server.address() as AddressInfo).port}/hook`;

	it('does nothing when the manifest sets no notify_url', async () => {
		expect(await notifyUrl(repo)).toBeNull();
		expect(await startNotifier(repo, 20)).toBeNull();
		await setManifest('dashboard:\n  notify_url: "ftp://example.com/x"\n');
		expect(await notifyUrl(new Repo(dir))).toBeNull();
	});

	it('posts an entry that appears after it started, once, and never what was already there', async () => {
		await addQuestion('Already here when the server started?');
		await setManifest(`dashboard:\n  notify_url: "${url()}"\n`);
		notifier = await startNotifier(repo, 20);
		expect(notifier).not.toBeNull();
		await repo.watch();
		await new Promise((r) => setTimeout(r, 400)); // let the watcher settle before writing

		await addQuestion('Which port should it use?');
		await waitFor(() => received.length >= 1);
		await notifier!.idle();
		expect(received).toHaveLength(1);
		const { body, headers } = received[0];
		expect(body.entry).toMatchObject({
			kind: 'question',
			title: 'Which port should it use?',
			href: '/docs/wip/agents/S-004.md'
		});
		expect(Object.keys(body).sort()).toEqual(['entry', 'project']);
		expect(Object.keys(body.entry).sort()).toEqual(['at', 'href', 'key', 'kind', 'title']);
		expect(headers['content-type']).toBe('application/json');
		expect(headers.authorization).toBeUndefined();

		// an unrelated change posts nothing more
		await writeFile(join(dir, 'design/system/extra.md'), '---\ntitle: x\n---\n');
		await new Promise((r) => setTimeout(r, 600));
		await notifier!.idle();
		expect(received).toHaveLength(1);
	});

	it('logs a failed post and carries on', async () => {
		await setManifest(`dashboard:\n  notify_url: "${url()}"\n`);
		notifier = await startNotifier(repo, 20);
		await repo.watch();
		await new Promise((r) => setTimeout(r, 400)); // let the watcher settle before writing
		status = 500;
		await addQuestion('First?');
		await waitFor(() => received.length >= 1);
		await notifier!.idle();
		status = 200;
		await addQuestion('Second?');
		await waitFor(() => received.length >= 2);
		await notifier!.idle();
		// the failed one is not retried; the next new entry still goes out
		expect(received.map((r) => r.body.entry.title)).toEqual(['First?', 'Second?']);
	});
});
