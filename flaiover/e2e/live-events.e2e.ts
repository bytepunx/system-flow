// S-0454: tabs that follow a project's live events must not hold the browser's connections to the
// dashboard. Chromium keeps at most six HTTP/1.1 connections to one host, and a server-sent event
// stream holds one for as long as its tab is open: with six tabs following events, every other
// request from every tab waited in the browser's queue, with no timeout. Over a WebSocket
// (ADR-0182) none is held. Seven pages of one browser context follow one project's events here, and
// each, the last among them, still gets an answer to a fetch once its stream is open.
import { expect, test, type Page } from '@playwright/test';
import { createHmac, randomBytes } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import WebSocket from 'ws';

const PAGES = 7;
const PROJECT = 'e2e';
/** How long a fetch from a page may take: /_health answers at once, so this is only the margin. */
const ANSWER_MS = 5000;
/** How long a page may take to load and open its stream. */
const OPEN_MS = 15000;

// what playwright.config.ts made for the run and gave the server
const dir = process.env.FLAIOVER_E2E_DIR ?? '';
const secret = (name: string) => readFileSync(join(dir, name), 'utf8').trim();

/** A flai on the host as far as the dashboard can tell (ADR-0029): it proves the shared credential
 * on /agent for one project and answers what it is asked at once, the project's manifest for
 * project.info and an error for the rest, so that no page's request waits on it. */
function fakeFlai(base: string, key: string): Promise<WebSocket> {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(base.replace(/^http/, 'ws') + '/agent');
		const mine = randomBytes(16).toString('hex');
		let proven = false;
		const send = (m: object) => ws.send(JSON.stringify({ jsonrpc: '2.0', ...m }));
		ws.on('open', () =>
			send({
				id: 'hello',
				method: 'hello',
				params: {
					protocol: 1,
					nonce: mine,
					flai: '0.0.0-e2e',
					project: { key: PROJECT, name: 'E2E' },
					methods: ['project.info']
				}
			})
		);
		ws.on('message', (data) => {
			const m = JSON.parse(data.toString());
			if (!proven) {
				if (!m.result?.nonce) return reject(new Error(`hello refused: ${data.toString()}`));
				proven = true;
				const proof = createHmac('sha256', key)
					.update(`flai|${m.result.nonce}|${mine}`)
					.digest('hex');
				send({ method: 'hello.prove', params: { proof } });
				return resolve(ws);
			}
			if (typeof m.id !== 'number') return;
			if (m.method === 'project.info')
				send({
					id: m.id,
					result: {
						name: 'E2E',
						version: 1,
						layout: { design: 'design', docs: 'docs', wip: 'wip' }
					}
				});
			else
				send({ id: m.id, error: { code: -32601, message: `${m.method}: not in this fake flai` } });
		});
		ws.on('close', (code, reason) => {
			if (!proven) reject(new Error(`/agent closed ${code} ${reason}`));
		});
		ws.on('error', reject);
	});
}

/** Resolves with the transport once the page has the project's event stream open: a WebSocket that
 * has had its `ready`, or, before ADR-0182, an EventSource whose response has begun. */
function streamOpen(page: Page): Promise<{ transport: string; url: string }> {
	const ours = (url: string) => new URL(url).pathname === '/api/events';
	return new Promise((resolve) => {
		page.on('websocket', (ws) => {
			if (!ours(ws.url())) return;
			ws.on('framereceived', ({ payload }) => {
				try {
					if (JSON.parse(String(payload)).type === 'ready')
						resolve({ transport: 'WebSocket', url: ws.url() });
				} catch {
					// not one of ours
				}
			});
		});
		page.on('response', (r) => {
			if (ours(r.url()) && r.status() === 200)
				resolve({ transport: 'server-sent events', url: r.url() });
		});
	});
}

function within<T>(ms: number, p: Promise<T>, what: string): Promise<T> {
	let timer: ReturnType<typeof setTimeout> | undefined;
	const late = new Promise<never>((_, reject) => {
		timer = setTimeout(() => reject(new Error(`${what}: nothing within ${ms} ms`)), ms);
	});
	return Promise.race([p, late]).finally(() => clearTimeout(timer));
}

let flai: WebSocket | undefined;

test.beforeAll(async ({ playwright, baseURL }) => {
	flai = await fakeFlai(baseURL!, secret('agent-key'));
	// connected once the dashboard lists the project as such
	const api = await playwright.request.newContext({
		baseURL,
		extraHTTPHeaders: { authorization: `Bearer ${secret('token')}` }
	});
	await expect
		.poll(async () => {
			const body = await (await api.get('/api/projects')).json();
			return body.projects?.find((p: { key: string }) => p.key === PROJECT)?.connected ?? false;
		})
		.toBe(true);
	await api.dispose();
});

test.afterAll(() => {
	flai?.terminate();
});

test(`${PAGES} pages of one context follow one project's events, and a fetch from each still answers`, async ({
	context,
	baseURL
}) => {
	test.setTimeout(PAGES * (OPEN_MS + ANSWER_MS) + 10000);
	await context.addCookies([{ name: 'flaiover_session', value: secret('token'), url: baseURL! }]);

	for (let n = 1; n <= PAGES; n++) {
		const page = await context.newPage();
		const open = streamOpen(page);
		await page.goto(`/license?project=${PROJECT}`, { waitUntil: 'commit', timeout: OPEN_MS });
		const stream = await within(OPEN_MS, open, `page ${n}'s event stream`);
		expect(new URL(stream.url).searchParams.get('project')).toBe(PROJECT);

		const answer = await page.evaluate(async (ms) => {
			const started = performance.now();
			try {
				const r = await fetch('/_health', { cache: 'no-store', signal: AbortSignal.timeout(ms) });
				return { status: r.status, ms: Math.round(performance.now() - started) };
			} catch (e) {
				return { error: String(e), ms: Math.round(performance.now() - started) };
			}
		}, ANSWER_MS);
		console.log(`page ${n}: ${stream.transport} open; /_health ${JSON.stringify(answer)}`);
		expect(
			answer,
			`a fetch from page ${n} of ${PAGES}, with ${n} event streams open over ${stream.transport}`
		).toMatchObject({ status: 200 });
	}
});
