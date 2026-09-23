import { describe, expect, it } from 'vitest';
import { activityLine, anyRunning, dotClass, type StoryActivity } from './activity';

const run = {
	story: 'S-0104',
	harness: 'claude-code',
	model: 'claude-haiku-4-5',
	command: 'claude',
	agent: 'agent-S-0104',
	started: '2026-09-23T18:00:00Z'
};

describe('agent activity (S-0104)', () => {
	it('is green while working, yellow while waiting, red when failed, and no dot once finished', () => {
		expect(dotClass('working')).toBe('bg-good');
		expect(dotClass('waiting')).toBe('bg-warn');
		expect(dotClass('failed')).toBe('bg-danger');
		expect(dotClass('worked')).toBeNull();
	});
	it('says what the agent is doing, with which harness and model, and why', () => {
		const a: StoryActivity = {
			state: 'waiting',
			why: 'waiting for an answer to TH-0001: Which port?',
			run
		};
		expect(activityLine(a)).toBe(
			'agent waiting (claude-code, claude-haiku-4-5): waiting for an answer to TH-0001: Which port?'
		);
		expect(
			activityLine({ state: 'working', run: { ...run, harness: undefined, model: undefined } })
		).toBe('agent working (claude)');
		expect(
			activityLine({ state: 'failed', why: 'ended (exit 1) with S-0104 in ready', run })
		).toContain('agent failed (claude-code, claude-haiku-4-5): ended (exit 1)');
	});
	it('knows whether any agent still runs', () => {
		expect(anyRunning(null)).toBe(false);
		expect(
			anyRunning({
				enabled: true,
				state: { command: '', stories: { 'S-0104': { state: 'failed', run } } }
			})
		).toBe(false);
		expect(
			anyRunning({
				enabled: true,
				state: { command: '', stories: { 'S-0104': { state: 'waiting', run } } }
			})
		).toBe(true);
	});
});
