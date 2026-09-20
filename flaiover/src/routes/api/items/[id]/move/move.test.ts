import { describe, expect, it } from 'vitest';
import { _moveArgs } from './+server';

describe('move arguments', () => {
	it('passes --yes only for an acceptance the designer chose to include files in', () => {
		expect(_moveArgs('S-0051', { to: 'done' })).toEqual(['move', 'S-0051', 'done']);
		expect(_moveArgs('S-0051', { to: 'done', include_uncommitted: true })).toEqual([
			'move',
			'S-0051',
			'done',
			'--yes'
		]);
	});
	it('never passes --yes for any other move, or for a value that is not true', () => {
		expect(_moveArgs('S-0051', { to: 'review', include_uncommitted: true })).toEqual([
			'move',
			'S-0051',
			'review'
		]);
		expect(
			_moveArgs('S-0051', { to: 'done', include_uncommitted: 'yes' as unknown as boolean })
		).toEqual(['move', 'S-0051', 'done']);
	});
	it('keeps reason and by', () => {
		expect(_moveArgs('S-1', { to: 'cancelled', reason: 'dup', by: 'alex' })).toEqual([
			'move',
			'S-1',
			'cancelled',
			'--reason',
			'dup',
			'--by',
			'alex'
		]);
	});
	it('previews a cancellation with --dry-run, and only a cancellation', () => {
		expect(_moveArgs('E-7', { to: 'cancelled', dry_run: true })).toEqual([
			'move',
			'E-7',
			'cancelled',
			'--reason',
			'preview',
			'--dry-run'
		]);
		expect(_moveArgs('E-7', { to: 'cancelled', dry_run: true, reason: 'why' })).toEqual([
			'move',
			'E-7',
			'cancelled',
			'--reason',
			'why',
			'--dry-run'
		]);
		expect(_moveArgs('S-1', { to: 'review', dry_run: true })).toEqual(['move', 'S-1', 'review']);
		expect(
			_moveArgs('E-7', { to: 'cancelled', reason: 'why', dry_run: 'yes' as unknown as boolean })
		).toEqual(['move', 'E-7', 'cancelled', '--reason', 'why']);
	});
});
