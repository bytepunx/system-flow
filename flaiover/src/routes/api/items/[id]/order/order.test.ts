import { describe, expect, it } from 'vitest';
import { _orderArgs } from './+server';

describe('order arguments', () => {
	it('builds one flai order call per kind of placement', () => {
		expect(_orderArgs('S-0059', { before: 'S-0061' })).toEqual([
			'order',
			'S-0059',
			'--before',
			'S-0061'
		]);
		expect(_orderArgs('S-0059', { after: 'S-0061' })).toEqual([
			'order',
			'S-0059',
			'--after',
			'S-0061'
		]);
		expect(_orderArgs('S-0059', { top: true })).toEqual(['order', 'S-0059', '--top']);
		expect(_orderArgs('S-0059', { bottom: true })).toEqual(['order', 'S-0059', '--bottom']);
	});
	it('refuses none or more than one placement', () => {
		expect(() => _orderArgs('S-0059', {})).toThrow(/exactly one of/);
		expect(() => _orderArgs('S-0059', { top: true, before: 'S-0061' })).toThrow(/exactly one of/);
		expect(() => _orderArgs('S-0059', { top: 'yes' as unknown as boolean })).toThrow(
			/exactly one of/
		);
	});
	it('passes only item IDs to flai, never something it could read as a flag', () => {
		expect(() => _orderArgs('S-0059', { before: '--top' })).toThrow(/not an item ID/);
		expect(() => _orderArgs('--json', { top: true })).toThrow(/not an item ID/);
		expect(() => _orderArgs('S-0059', { after: 'S-0061 --bottom' })).toThrow(/not an item ID/);
	});
});
