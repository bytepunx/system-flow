import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { RepoError } from './repo';
import { respond } from './respond';
import { projectIdentity, setIdentityHeaders } from './project';
import { repo } from './repo';

const fixture = resolve('../flai/internal/metrics/testdata/good');

describe('project identity in responses (ADR-0024)', () => {
	let dir: string;
	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-project-'));
		await cp(fixture, dir, { recursive: true });
		process.env.PROJECT_DIR = dir;
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('reads name and key from the manifest', async () => {
		expect(await projectIdentity(repo())).toEqual({ name: 'good', key: 'g' });
	});

	it('adds project to a JSON object body', async () => {
		const r = await respond(async () => ({ total: 3 }));
		expect(await r.json()).toEqual({ total: 3, project: { name: 'good', key: 'g' } });
	});

	it('leaves an array body as it is', async () => {
		const r = await respond(async () => [{ id: 'S-001' }]);
		expect(await r.json()).toEqual([{ id: 'S-001' }]);
	});

	it('adds project to an error body, beside the error and its data', async () => {
		const r = await respond(async () => {
			throw new RepoError(409, 'changed', { hash: 'h' });
		});
		expect(r.status).toBe(409);
		expect(await r.json()).toEqual({
			error: 'changed',
			hash: 'h',
			project: { name: 'good', key: 'g' }
		});
	});

	it('names the project in headers on any response, a stream included, with the name encoded', () => {
		const stream = new Response(new ReadableStream(), {
			headers: { 'content-type': 'application/x-ndjson' }
		});
		const out = setIdentityHeaders(stream, { name: 'Größe & co', key: 'gc' });
		expect(out.headers.get('x-flai-project-key')).toBe('gc');
		expect(decodeURIComponent(out.headers.get('x-flai-project-name')!)).toBe('Größe & co');
		expect(out.headers.get('content-type')).toBe('application/x-ndjson');
	});
});
