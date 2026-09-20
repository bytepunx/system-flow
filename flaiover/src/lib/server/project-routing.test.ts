// S-0080: repo() and agent() are called the way they always were, everywhere in the app, but resolve
// to the project current for the request in progress (or the one dashboard's default project, for
// code with no request of its own). This is the plumbing hooks.server.ts relies on; it is tested
// directly here so the routing itself, not any one route, is what is proven.
import { afterEach, describe, expect, it } from 'vitest';
import { AgentError, AgentRegistry, defaultProjectKey, withProject, type AgentHub } from './agent';
import { Repo, RepoError, useRepo, knownRepos, repo, type Ask } from './repo';

describe('project routing (S-0080)', () => {
	afterEach(() => {
		for (const key of [...knownRepos().keys()]) useRepo(null, key);
	});

	function fakeSource(project: string) {
		const asked: string[] = [];
		return {
			asked,
			source: (async (method: string) => {
				asked.push(`${project}:${method}`);
				return {
					name: project,
					version: 1,
					layout: { design: 'design', docs: 'docs', wip: 'wip' }
				};
			}) as Ask
		};
	}

	it('repo() resolves to the project current when it is called, one Repo per key', async () => {
		const a = fakeSource('a');
		const b = fakeSource('b');
		useRepo(new Repo('/a', a.source), 'a');
		useRepo(new Repo('/b', b.source), 'b');

		const seenA = withProject('a', () => repo());
		const seenB = withProject('b', () => repo());
		expect(seenA).not.toBe(seenB);
		expect(withProject('a', () => repo())).toBe(seenA); // the same Repo again, not a new one

		await withProject('a', () => repo().manifest());
		await withProject('b', () => repo().manifest());
		expect(a.asked).toEqual(['a:project.info']);
		expect(b.asked).toEqual(['b:project.info']);
	});

	it('one project’s cache never answers for another’s', async () => {
		const a = fakeSource('a');
		const b = fakeSource('b');
		useRepo(new Repo('/a', a.source), 'a');
		useRepo(new Repo('/b', b.source), 'b');
		const [ma, mb] = await Promise.all([
			withProject('a', () => repo().manifest()),
			withProject('b', () => repo().manifest())
		]);
		expect(ma.name).toBe('a');
		expect(mb.name).toBe('b');
		// asked twice more, once per project: no cached answer crossed over
		await withProject('a', () => repo().manifest());
		await withProject('b', () => repo().manifest());
		expect(a.asked).toEqual(['a:project.info']); // remembered within "a", not re-asked
		expect(b.asked).toEqual(['b:project.info']);
	});

	it('with no project set, repo() and agent() are the one dashboard has always had', () => {
		expect(repo()).toBe(repo());
		const withKey = withProject(defaultProjectKey, () => repo());
		expect(withKey).toBe(repo());
	});

	it('a project key nothing has ever connected for is refused, told apart from one merely not connected', async () => {
		const registry = new AgentRegistry('shared-key');
		registry.hub('known'); // as a real connection attempt would vivify it
		try {
			const known = registry.peek('known') as AgentHub;
			await expect(known.ask('project.info')).rejects.toMatchObject({ status: 503 }); // not connected right now
			expect(registry.peek('never-heard-of')).toBeUndefined();
		} finally {
			registry.close();
		}
	});

	it('RepoError carries the dashboard’s own 404 for a key nothing is known about, the same shape as flai’s', () => {
		const e = new AgentError(404, 'no project "x" is being served here');
		expect(e.status).toBe(404);
		expect(new RepoError(e.status, e.message)).toMatchObject({ status: 404 });
	});
});
