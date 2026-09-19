import { describe, expect, it } from 'vitest';
import {
	canReorderOnto,
	dropPlacement,
	endPlacement,
	keyStep,
	reorderable,
	stepPlacement
} from './reorder';

const card = (id: string, status: string, type = 'story') => ({ id, status, type });
const column = ['S-1', 'S-2', 'S-3', 'S-4'];

describe('what can be reordered', () => {
	it('is a story in backlog or ready on a writable board', () => {
		expect(reorderable(card('S-1', 'ready'), true)).toBe(true);
		expect(reorderable(card('S-1', 'backlog'), true)).toBe(true);
		for (const status of ['in-progress', 'review', 'done', 'cancelled'])
			expect(reorderable(card('S-1', status), true), status).toBe(false);
		expect(reorderable(card('E-1', 'backlog', 'epic'), true)).toBe(false);
		expect(reorderable(card('T-1', 'ready', 'task'), true)).toBe(false);
		expect(reorderable(card('S-1', 'ready'), false)).toBe(false);
		expect(reorderable(undefined, true)).toBe(false);
	});
	it('is onto another story of the same column only', () => {
		expect(canReorderOnto(card('S-1', 'ready'), card('S-2', 'ready'), true)).toBe(true);
		expect(canReorderOnto(card('S-1', 'ready'), card('S-1', 'ready'), true)).toBe(false);
		expect(canReorderOnto(card('S-1', 'ready'), card('S-2', 'backlog'), true)).toBe(false);
		expect(canReorderOnto(card('S-1', 'ready'), card('T-2', 'ready', 'task'), true)).toBe(false);
		expect(canReorderOnto(card('S-1', 'review'), card('S-2', 'review'), true)).toBe(false);
		expect(canReorderOnto(card('S-1', 'ready'), card('S-2', 'ready'), false)).toBe(false);
	});
});

describe('a drop', () => {
	it('asks for before or after the card it landed on', () => {
		expect(dropPlacement(column, 'S-4', 'S-1', 'above')).toEqual({ before: 'S-1' });
		expect(dropPlacement(column, 'S-4', 'S-1', 'below')).toEqual({ after: 'S-1' });
		expect(dropPlacement(column, 'S-1', 'S-4', 'below')).toEqual({ after: 'S-4' });
		expect(dropPlacement(column, 'S-1', 'S-3', 'above')).toEqual({ before: 'S-3' });
	});
	it('asks for nothing when the card would stay where it is', () => {
		expect(dropPlacement(column, 'S-2', 'S-3', 'above')).toBeNull();
		expect(dropPlacement(column, 'S-2', 'S-1', 'below')).toBeNull();
		expect(dropPlacement(column, 'S-2', 'S-2', 'above')).toBeNull();
		expect(dropPlacement(column, 'S-9', 'S-2', 'above')).toBeNull();
		expect(dropPlacement(column, 'S-2', 'S-9', 'above')).toBeNull();
	});
	it('in the column’s empty space goes last, unless it already is', () => {
		expect(endPlacement(column, 'S-2')).toEqual({ bottom: true });
		expect(endPlacement(column, 'S-4')).toBeNull();
		expect(endPlacement(column, 'S-9')).toBeNull();
	});
});

describe('a step', () => {
	it('moves past the neighbour', () => {
		expect(stepPlacement(column, 'S-2', 'up')).toEqual({ before: 'S-1' });
		expect(stepPlacement(column, 'S-2', 'down')).toEqual({ after: 'S-3' });
	});
	it('has nowhere to go at the ends', () => {
		expect(stepPlacement(column, 'S-1', 'up')).toBeNull();
		expect(stepPlacement(column, 'S-4', 'down')).toBeNull();
		expect(stepPlacement(['S-1'], 'S-1', 'down')).toBeNull();
		expect(stepPlacement(column, 'S-9', 'up')).toBeNull();
	});
	it('is Alt with an arrow and nothing else', () => {
		const key = (
			key: string,
			mods: Partial<Record<'altKey' | 'ctrlKey' | 'metaKey' | 'shiftKey', boolean>>
		) => keyStep({ key, altKey: false, ctrlKey: false, metaKey: false, shiftKey: false, ...mods });
		expect(key('ArrowUp', { altKey: true })).toBe('up');
		expect(key('ArrowDown', { altKey: true })).toBe('down');
		expect(key('ArrowUp', {})).toBeNull();
		expect(key('ArrowUp', { altKey: true, shiftKey: true })).toBeNull();
		expect(key('ArrowLeft', { altKey: true })).toBeNull();
	});
});
