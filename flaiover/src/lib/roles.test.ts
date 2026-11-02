// Behaviour tests for how the role page reads a role as agent.status's roles carry it (S-0420): the
// order its runs are listed in, the controls a role offers, the write each control is, where a
// run's stream is read, which agent events concern a role, and what its action's state says.

import { describe, expect, it } from 'vitest';
import {
	concernsRole,
	roleControls,
	roleRuns,
	roleState,
	roleStream,
	roleWrite,
	triggerLine,
	type RoleState
} from './roles';
import type { RoleRun } from './strategic';

const reviewing: RoleRun = {
	item: 'S-0001',
	role: 'review',
	role_name: 'reviewer',
	agent: 'reviewer',
	command: 'claude',
	pid: 42,
	started: '2026-10-28T10:00:00Z'
};

const role = (name: string, over: Partial<RoleState> = {}): RoleState => ({
	name,
	role: 'review',
	path: `design/workflows/roles/${name}.md`,
	item: [],
	trigger: ['asked'],
	action: { name: 'review', enabled: true },
	run: null,
	running: [],
	runs: [],
	held: false,
	...over
});

const reviewer = role('reviewer', {
	item: ['story'],
	trigger: ['asked', 'schedule'],
	run: reviewing,
	running: [reviewing]
});

describe('roles (S-0420)', () => {
	it('reads a role whose lists or run flai left null or out as empty, and none', () => {
		const raw = { ...role('reviewer'), item: null, trigger: null, running: null, runs: null };
		delete (raw as Partial<typeof raw>).run;
		expect(roleState(raw)).toEqual(role('reviewer', { trigger: [] }));
	});

	it('lists the runs under way before those that ended', () => {
		const ended = { ...reviewing, item: 'S-0002', ended: '2026-10-28T09:00:00Z' };
		expect(roleRuns({ ...reviewer, runs: [ended] })).toEqual([reviewing, ended]);
	});

	it('offers Run, with an item when its run takes one, for a role the operator may ask for', () => {
		expect(roleControls(reviewer)).toEqual({ run: true, item: true, stop: false, start: false });
		expect(roleControls(role('auditor'))).toEqual({
			run: true,
			item: false,
			stop: false,
			start: false
		});
		expect(roleControls(role('watcher', { item: ['story'], trigger: ['event'] }))).toEqual({
			run: false,
			item: false,
			stop: false,
			start: false
		});
	});

	it('offers Stop while the orchestrator is not held and Start while it is, and no other role either', () => {
		// the baseline orchestrator.md names asked among its triggers: it is started, not run
		const orchestrator = role('orchestrator', { role: 'orchestrate', trigger: ['asked', 'event'] });
		expect(roleControls(orchestrator)).toMatchObject({ run: false, stop: true, start: false });
		expect(roleControls({ ...orchestrator, held: true })).toMatchObject({
			stop: false,
			start: true
		});
		expect(roleControls({ ...reviewer, held: true })).toMatchObject({ stop: false, start: false });
	});

	it('runs a defined role through role.run, on its item when given, and a built-in by its own method', () => {
		expect(roleWrite(reviewer, 'run', 'S-0001')).toEqual({
			method: 'role.run',
			params: { name: 'reviewer', item: 'S-0001' }
		});
		expect(roleWrite(role('auditor'), 'run')).toEqual({
			method: 'role.run',
			params: { name: 'auditor' }
		});
		expect(roleWrite(role('planner', { role: 'plan' }), 'run', 'E-0016')).toEqual({
			method: 'plan.run',
			params: { id: 'E-0016' }
		});
		expect(roleWrite(role('planner', { role: 'plan' }), 'run')).toEqual({
			refused: 'the planner plans an epic or a story: give its item'
		});
		expect(roleWrite(role('analyzer', { role: 'analyze' }), 'run')).toEqual({
			method: 'analyze.run',
			params: {}
		});
		expect(roleWrite(role('curator', { role: 'curate' }), 'run')).toEqual({
			method: 'curate.run',
			params: {}
		});
		expect(roleWrite(role('orchestrator', { role: 'orchestrate' }), 'run')).toHaveProperty(
			'refused'
		);
	});

	it("stops and starts the orchestrator's hold, and refuses either for a role with none", () => {
		const orchestrator = role('orchestrator', { role: 'orchestrate' });
		expect(roleWrite(orchestrator, 'stop')).toEqual({ method: 'orchestrate.stop', params: {} });
		expect(roleWrite(orchestrator, 'start')).toEqual({ method: 'orchestrate.start', params: {} });
		expect(roleWrite(reviewer, 'stop')).toEqual({
			refused: 'only the orchestrator has a hold to stop; the reviewer has none'
		});
		expect(roleWrite(reviewer, 'start')).toHaveProperty('refused');
	});

	it("reads a planner run's stream by its item, the orchestrator's and analyzer's by verb, and any other's by name", () => {
		expect(roleStream(role('planner', { role: 'plan' }), { item: 'E-0016' })).toEqual({
			story: 'E-0016',
			plan: true
		});
		expect(roleStream(role('orchestrator', { role: 'orchestrate' }), {})).toEqual({
			story: 'orchestrator',
			role: 'orchestrate'
		});
		expect(roleStream(role('analyzer', { role: 'analyze' }), {})).toEqual({
			story: 'analyzer',
			role: 'analyze'
		});
		expect(roleStream(reviewer, reviewing)).toEqual({ story: 'reviewer', role: 'reviewer' });
		expect(roleStream(role('curator', { role: 'curate' }), {})).toEqual({
			story: 'curator',
			role: 'curator'
		});
	});

	it('is concerned by an agent event naming it, and, when its run takes an item, by one naming an item', () => {
		expect(concernsRole(role('auditor'), 'auditor')).toBe(true);
		expect(concernsRole(role('auditor'), 'S-0001')).toBe(false);
		expect(concernsRole(role('auditor'), 'orchestrator')).toBe(false);
		expect(concernsRole(reviewer, 'S-0001')).toBe(true);
		expect(concernsRole(reviewer, 'E-0016')).toBe(true);
		expect(concernsRole(reviewer, 'analyzer')).toBe(false);
	});

	it('says what its host action does from its triggers', () => {
		expect(triggerLine(reviewer)).toBe(
			'flai serve starts the reviewer when you Run it and when its schedule comes round.'
		);
		expect(triggerLine(role('watcher', { trigger: ['asked', 'event', 'schedule'] }))).toBe(
			'flai serve starts the watcher when you Run it, on the events its file names and when its schedule comes round.'
		);
		expect(triggerLine(role('auditor', { trigger: [] }))).toBe(
			'flai serve starts the auditor as its file says.'
		);
		expect(triggerLine(role('orchestrator', { role: 'orchestrate' }))).toContain(
			'for as long as the action stays on'
		);
	});
});
