import { describe, expect, it } from 'vitest';
import { cardMenu } from './cardmenu';
import type { StoryActivity } from './activity';

// S-0202: a board card's menu offers what the item's page does, to a writer
describe('a card’s menu', () => {
	const story = { id: 'S-0202', type: 'story', status: 'backlog', blocked: false };
	const writer = { writable: true, agentEnabled: true };
	const labels = (...args: Parameters<typeof cardMenu>) => cardMenu(...args).map((e) => e.label);
	const failed: StoryActivity = {
		state: 'failed',
		run: {
			story: 'S-0202',
			command: 'claude',
			agent: 'agent-S-0202',
			started: '2026-10-03T08:00:00Z',
			ended: '2026-10-03T08:30:00Z',
			outcome: 'failed'
		}
	};

	it('offers Finalize on a writable draft story only', () => {
		expect(labels({ ...story, draft: true }, writer)).toContain('Finalize');
		expect(labels({ ...story, draft: false }, writer)).not.toContain('Finalize');
		expect(labels(story, writer)).not.toContain('Finalize');
		expect(labels({ ...story, type: 'task', draft: true }, writer)).not.toContain('Finalize');
		expect(labels({ ...story, type: 'epic', draft: true }, writer)).not.toContain('Finalize');
		expect(labels({ ...story, draft: true, archived: true }, writer)).not.toContain('Finalize');
		expect(labels({ ...story, draft: true }, { ...writer, writable: false })).not.toContain(
			'Finalize'
		);
	});
	it('offers only Open on a read-only board', () => {
		const reader = { writable: false, agentEnabled: true };
		expect(cardMenu({ ...story, status: 'ready', draft: true }, reader)).toEqual([
			{ action: 'open', label: 'Open' }
		]);
		expect(labels({ ...story, status: 'in-progress', blocked: true }, reader)).toEqual(['Open']);
	});
	it('offers Unblock on a blocked item, else Block…', () => {
		expect(labels(story, writer)).toContain('Block…');
		expect(labels(story, writer)).not.toContain('Unblock');
		expect(labels({ ...story, blocked: true }, writer)).toContain('Unblock');
		expect(labels({ ...story, blocked: true }, writer)).not.toContain('Block…');
	});
	it('offers no Block… or Cancel… on a done or cancelled item, nor anything on an archived one', () => {
		for (const status of ['done', 'cancelled']) {
			expect(labels({ ...story, status }, writer)).toEqual(['Open']);
		}
		expect(labels({ ...story, archived: true, blocked: true }, writer)).toEqual(['Open']);
	});
	it('offers the story page’s agent action, with the agent action on', () => {
		expect(cardMenu({ ...story, status: 'ready' }, writer)).toContainEqual({
			action: 'agent',
			label: 'Start agent',
			agent: 'start'
		});
		expect(
			cardMenu({ ...story, status: 'in-progress' }, { ...writer, activity: failed })
		).toContainEqual({ action: 'agent', label: 'Retry', agent: 'restart' });
		expect(
			cardMenu({ ...story, status: 'ready' }, { ...writer, agentEnabled: false }).map(
				(e) => e.action
			)
		).not.toContain('agent');
		expect(
			cardMenu({ ...story, type: 'task', status: 'ready' }, writer).map((e) => e.action)
		).not.toContain('agent');
	});
	it('lists Open, Finalize, the agent, Block… or Unblock, then Cancel…', () => {
		expect(cardMenu({ ...story, status: 'ready', draft: true }, writer)).toEqual([
			{ action: 'open', label: 'Open' },
			{ action: 'finalize', label: 'Finalize' },
			{ action: 'agent', label: 'Start agent', agent: 'start' },
			{ action: 'block', label: 'Block…' },
			{ action: 'cancel', label: 'Cancel…' }
		]);
		expect(labels({ ...story, status: 'in-progress', blocked: true }, writer)).toEqual([
			'Open',
			'Unblock',
			'Cancel…'
		]);
	});
});
