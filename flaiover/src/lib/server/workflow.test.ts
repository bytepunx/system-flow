// S-0756: a state added to the workflow file alone, in a project rendered from the template, is
// followed by the dashboard's surfaces, run against the flai built from this tree: board.get lays
// the column out, item.move makes and refuses the moves the file gives it, and stats.get and the
// Charts page's builders report and draw the timer it adds and stack the state. The fixture is the
// one flai/cmd/workflow_e2e_test.go writes over the template's board: testing, of the active kind
// with a WIP limit, between in-progress and review; has-task on the move into it; no move from
// in-progress to review; and the timer testing-to-accept, from the first move into testing to the
// acceptance.
import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { execFile, execFileSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { mkdtemp, readdir, readFile, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { promisify } from 'node:util';
import { Repo } from './repo';
import { board } from './board';
import { stats } from './stats';
import { flaiAsk, flaiBin, haveFlai, shell } from './testing';
import { build, normalise, timerCharts, timerKind, titleOf, type Report } from '$lib/viz/charts';
import { theme } from '$lib/viz/palette';
import type { Flow } from '$lib/lanes';

const template = resolve('../template');
const fixture = resolve('../flai/internal/workflow/testdata/e2e-board.md');
const haveGit = (() => {
	try {
		execFileSync('git', ['--version'], { stdio: 'ignore' });
		return true;
	} catch {
		return false;
	}
})();
const run = promisify(execFile);

/** The tree's flai in a folder, with its output as it prints it, for a command with no --json. */
async function flai(cwd: string, args: string[]): Promise<string> {
	const { stdout } = await run(flaiBin, args, { cwd, env: { ...process.env, FLAI_AGENT: '' } });
	return stdout;
}

function git(cwd: string, ...args: string[]): string {
	return execFileSync('git', args, { cwd, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] });
}

/** A second on: the clock flai stamps a move with counts whole seconds. */
const tick = () => new Promise((done) => setTimeout(done, 1100));

type Bars = { series: { name: string; stack?: string; data: { ids: string[] }[] }[] };
type Areas = { series: { name: string; data: [string, number][] }[] };
type Scatter = { series: { data: { id: string; value: [string, number] }[] }[] };

describe.skipIf(!haveFlai || !haveGit || !existsSync(join(template, 'template.yaml')))(
	"a state added to the template's workflow file, through the dashboard's surfaces",
	() => {
		let base: string;
		let dir: string;
		const saved: Record<string, string | undefined> = {};
		const repo = () => new Repo(dir, flaiAsk(dir));
		/** The item's state as the board reads it, or undefined when no column holds it. */
		const status = async (id: string) =>
			Object.values((await board(repo())).columns)
				.flat()
				.find((c) => c.id === id)?.status;

		beforeAll(async () => {
			base = await mkdtemp(join(tmpdir(), 'flaiover-workflow-'));
			dir = join(base, 'e2e');
			// a configuration of the test's own, and none of the session's story or role
			for (const k of ['FLAI_CONFIG', 'FLAI_STORY', 'FLAI_ROLE', 'FLAI_STARTED_BY']) {
				saved[k] = process.env[k];
				delete process.env[k];
			}
			process.env.FLAI_CONFIG = join(base, 'config.json');
			await flai(base, [
				'new',
				dir,
				'--template',
				template,
				'--defaults',
				'--no-git',
				'--var',
				'description=End to end',
				'--var',
				'owner=olive'
			]);
			await writeFile(join(dir, 'design/workflows/board.md'), await readFile(fixture, 'utf8'));
			await flai(dir, ['config', 'set', 'author', 'olive']);
			await flai(dir, ['agent', 'set', '--harness', 'command']);
			await flai(dir, ['board', 'limit', 'testing', '2']);
			git(dir, 'init', '-q', '-b', 'main');
			git(dir, 'config', 'user.email', 't@t');
			git(dir, 'config', 'user.name', 't');
			git(dir, 'add', '-A');
			git(dir, 'commit', '-q', '-m', 'chore: render the template with testing added');
		}, 60000);
		afterAll(async () => {
			for (const [k, v] of Object.entries(saved)) {
				if (v === undefined) delete process.env[k];
				else process.env[k] = v;
			}
			if (base) await rm(base, { recursive: true, force: true });
		});

		it('board.get lays testing out between in-progress and review, active and with its limit (criterion 3)', async () => {
			const b = await board(repo());
			const names = b.lanes.map((l) => l.name);
			const at = names.indexOf('testing');
			expect(names.slice(at - 1, at + 2)).toEqual(['in-progress', 'testing', 'review']);
			expect(b.lanes[at]).toMatchObject({
				name: 'testing',
				title: 'Testing',
				kind: 'active',
				wip_limit: 2,
				takes_limit: true
			});
			expect(b.wip_limits.testing).toBe(2);
			expect(b.columns.testing).toEqual([]);
			// the moves the board may offer into and out of it are the file's
			expect(b.transitions).toContainEqual(
				expect.objectContaining({ from: ['in-progress'], to: 'testing', operator: true })
			);
			expect(b.transitions).toContainEqual(
				expect.objectContaining({ from: ['testing'], to: 'review', operator: true })
			);
			expect(b.transitions.some((m) => m.from.includes('in-progress') && m.to === 'review')).toBe(
				false
			);
		}, 30000);

		it('item.move takes a story through testing only once it has a task, and refuses the moves the file does not give, changing nothing (criterion 4)', async () => {
			await flai(dir, ['epic', 'new', 'Epic']);
			for (const title of ['Through testing', 'Left in testing'])
				await flai(dir, ['story', 'new', title, '--epic', 'E-0001']);
			// a criterion written, as the move to ready needs
			const stories = join(dir, 'wip/kanban/stories');
			for (const id of ['S-0001', 'S-0002']) {
				const name = (await readdir(stories)).find((f) => f.startsWith(id + '-'))!;
				const body = await readFile(join(stories, name), 'utf8');
				expect(body).toContain('- [ ]');
				await writeFile(join(stories, name), body.replace(/- \[ \][^\n]*/, '- [x] it works'));
			}
			const r = repo();
			const move = async (id: string, to: string) => {
				const { data } = await r.write<{ id: string; status: string }>('item.move', { id, to });
				expect(data).toMatchObject({ id, status: to });
				expect(await status(id)).toBe(to);
			};
			const refused = async (id: string, to: string, why: string) => {
				const before = await status(id);
				await expect(r.write('item.move', { id, to })).rejects.toMatchObject({
					status: 400,
					message: expect.stringContaining(why)
				});
				expect(await status(id)).toBe(before);
			};

			for (const id of ['S-0001', 'S-0002']) {
				await move(id, 'ready');
				await move(id, 'in-progress');
			}
			// its branch and worktree, which the acceptance merges
			await flai(dir, ['stream', 'open', 'S-0001']);
			for (const id of ['S-0001', 'S-0002']) {
				// no transition leads from in-progress to review on this board
				await refused(id, 'review', 'cannot go from in-progress to review');
				// has-task refuses the move into testing until the story has a task
				await refused(id, 'testing', 'needs at least one task');
				await shell(dir, ['task', 'new', '--story', id, 'Step']);
				await move(id, 'testing');
			}
			// both are in the column, within its limit
			expect((await board(repo())).columns.testing.map((c) => c.id).sort()).toEqual([
				'S-0001',
				'S-0002'
			]);
			// a story in testing goes on to review, and not back past in-progress
			await refused('S-0002', 'ready', 'cannot go from testing to ready');
			await tick();
			await move('S-0001', 'review');
			// S-0001's task is done from in-progress, as the acceptance needs
			for (const to of ['ready', 'in-progress', 'done']) await flai(dir, ['move', 'T-0001', to]);
		}, 90000);

		it('stats.get reports testing-to-accept for the accepted story, and the Charts page draws it and stacks testing (criterion 6)', async () => {
			git(dir, 'add', '-A');
			git(dir, 'commit', '-q', '-m', 'chore: the stories through testing');
			const wt = join(dir, '.flai-cache/worktrees/S-0001');
			await writeFile(join(wt, 'design/tested.md'), "the story's work\n");
			git(wt, 'add', '-A');
			git(wt, 'commit', '-q', '-m', "docs: [S-0001] the story's work");
			await tick();
			const accepted = await flai(dir, [
				'accept',
				'S-0001',
				'--by',
				'olive',
				'--trailer',
				'Co-Authored-By: t <t@t>'
			]);
			expect(accepted).toContain('accepted S-0001: done');

			const report = await stats<Report>(repo(), { since: '30d' });
			expect(report.timers).toContainEqual(
				expect.objectContaining({ name: 'testing-to-accept', title: 'Testing to acceptance' })
			);
			const s1 = report.items.find((i) => i.id === 'S-0001')!;
			expect(s1).toBeDefined();
			expect(s1.completed).toBeTruthy();
			// from the first move into testing to the acceptance: its time in testing and in review
			const inState = s1.time_in_state_seconds;
			expect(inState.testing).toBeGreaterThan(0);
			expect(s1.timers?.['testing-to-accept']).toBeGreaterThan(0);
			expect(s1.timers?.['testing-to-accept']).toBe(inState.testing + inState.review);
			expect(report.summary.timers?.['testing-to-accept']).toBeDefined();
			expect(report.summary.time_in_state_share.testing).toBeGreaterThan(0);

			// the Charts page: the board's flow beside the report, as the page asks for both
			const b = await board(repo());
			const flow: Flow = { lanes: b.lanes, transitions: b.transitions };
			const t = theme(false);
			const kind = timerKind('testing-to-accept');
			expect(timerCharts(report.timers ?? [])).toContainEqual({
				kind,
				title: 'Testing to acceptance'
			});
			expect(titleOf(kind, 'day', report.timers)).toBe('Testing to acceptance');
			const scatter = build(kind, report, t, undefined, {}, flow) as Scatter;
			expect(scatter.series.flatMap((s) => s.data.map((d) => d.id))).toEqual(['S-0001']);

			// time in state stacks testing in board order, with S-0001's day among its bars
			const bars = build('time-in-state', report, t, undefined, {}, flow) as Bars;
			const stacked = bars.series.map((s) => s.name);
			expect(stacked.slice(stacked.indexOf('in-progress'), stacked.indexOf('review') + 1)).toEqual([
				'in-progress',
				'testing',
				'review'
			]);
			const testing = bars.series.find((s) => s.name === 'testing')!;
			expect(testing.stack).toBe('state');
			expect(testing.data.flatMap((d) => d.ids)).toContain('S-0001');
			// and from the report's own states, with no board read
			expect((build('time-in-state', report, t) as Bars).series.map((s) => s.name)).toContain(
				'testing'
			);

			// cumulative flow stacks testing, which S-0002 is still in
			const flowChart = build('cfd', normalise(report), t, undefined, {}, flow) as Areas;
			const band = flowChart.series.find((s) => s.name === 'testing');
			expect(band).toBeDefined();
			expect(band!.data.some(([, n]) => n > 0)).toBe(true);
		}, 90000);
	}
);
