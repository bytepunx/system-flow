import { describe, expect, it } from 'vitest';
import {
	activityLine,
	anyRunning,
	dotClass,
	holdLine,
	holdWaitsFor,
	reasonParts,
	storyActivity,
	type StoryActivity
} from './activity';

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
		expect(dotClass('working')).toBe('bg-dot-working');
		expect(dotClass('waiting')).toBe('bg-dot-waiting');
		expect(dotClass('failed')).toBe('bg-dot-failed');
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

describe('a held story (S-0129)', () => {
	const overlap = {
		code: 'overlap',
		reason:
			'held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 is accepted, cancelled, or sent back'
	};
	const noTouches = {
		code: 'no-touches',
		reason:
			"held (no-touches): declares no touches, so it may change what S-0128 (in progress) and S-0131 (in review) change; starts when it declares touches that overlap no open story's, or when S-0128 and S-0131 are accepted, cancelled, or sent back"
	};
	// a story that never had an agent: flai gives it a run that names only the story
	const standIn = { story: 'S-0129', command: '', agent: '', started: '' };

	it('waits for the stories its reason says clear it, once each', () => {
		expect(holdWaitsFor(overlap)).toEqual(['S-0128']);
		expect(holdWaitsFor(noTouches)).toEqual(['S-0128', 'S-0131']);
		expect(holdWaitsFor({ code: 'overlap', reason: 'held (overlap)' })).toEqual([]);
	});
	// S-0130: after: names the stories to wait for; a note may name the story itself, and a story
	// held by after: and an overlap has two clauses of what clears it
	it('waits for the stories a story names in after, and for an overlap besides', () => {
		const cancelled = {
			code: 'after',
			reason:
				'held (after): waits for S-0004 (cancelled); starts when S-0004 is done; S-0004 was cancelled, so drop it from after: if S-0009 no longer needs it'
		};
		const both = {
			code: 'after',
			reason:
				'held (after): waits for S-0001 (in backlog); starts when S-0001 is done; also held (overlap): touches flai/cmd, inside flai which S-0002 (in progress) touches; starts when S-0002 is accepted, cancelled, or sent back'
		};
		expect(holdWaitsFor(cancelled)).toEqual(['S-0004']);
		expect(holdWaitsFor(both)).toEqual(['S-0001', 'S-0002']);
		expect(holdLine(both)).toBe('held (after): S-0001, S-0002');
	});
	it('puts the code and the stories it waits for on the card', () => {
		expect(holdLine(overlap)).toBe('held (overlap): S-0128');
		expect(holdLine(noTouches)).toBe('held (no-touches): S-0128, S-0131');
		expect(holdLine({ code: 'overlap', reason: 'held (overlap)' })).toBe('held (overlap)');
	});
	it('cuts the reason at each story ID so each can be a link', () => {
		expect(
			reasonParts('touches x, which S-0128 (in progress) touches; starts when S-0128 is')
		).toEqual([
			{ text: 'touches x, which ' },
			{ text: 'S-0128', id: 'S-0128' },
			{ text: ' (in progress) touches; starts when ' },
			{ text: 'S-0128', id: 'S-0128' },
			{ text: ' is' }
		]);
	});
	it('says the hold, not an agent with no name', () => {
		expect(
			activityLine({ state: 'waiting', why: overlap.reason, run: standIn, hold: overlap })
		).toBe(overlap.reason);
	});
	it('shows held stories while the agent action is off, and every story while it is on', () => {
		const stories: Record<string, StoryActivity> = {
			'S-0104': { state: 'working', run },
			'S-0129': { state: 'waiting', run: standIn, hold: overlap }
		};
		expect(Object.keys(storyActivity({ enabled: true, state: { command: '', stories } }))).toEqual([
			'S-0104',
			'S-0129'
		]);
		expect(Object.keys(storyActivity({ enabled: false, state: { command: '', stories } }))).toEqual(
			['S-0129']
		);
		expect(storyActivity(null)).toEqual({});
		expect(storyActivity({ enabled: false })).toEqual({});
	});
});
