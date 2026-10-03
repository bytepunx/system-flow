import { describe, expect, it } from 'vitest';
import { amount, costOfDelayLines, forecastLines } from './planning';

describe('planning (S-0199)', () => {
	it('says an amount with separators and at most two decimals', () => {
		expect([amount(1500), amount(0), amount(12.345), amount(1234567.5)]).toEqual([
			'1,500',
			'0',
			'12.35',
			'1,234,567.5'
		]);
	});

	it('says a whole cost of delay: value, inputs, and who set it', () => {
		expect(
			costOfDelayLines({
				inputs: { revenue_per_week: 1000, penalty_per_week: 500, time_lost_per_cycle: '4h' },
				value: 1500,
				by: 'alex',
				at: '2026-10-03T09:00:00Z'
			})
		).toEqual([
			'cost of delay: 1,500 per week',
			'revenue 1,000 per week · penalty 500 per week · 4h lost per cycle',
			'set by alex at 2026-10-03T09:00:00Z'
		]);
	});

	it('says a partial cost of delay: a zero amount is given, an absent one is not', () => {
		expect(costOfDelayLines({ inputs: { penalty_per_week: 0 }, by: 'planner' })).toEqual([
			'cost of delay',
			'penalty 0 per week',
			'set by planner'
		]);
		expect(costOfDelayLines({ value: 250 })).toEqual(['cost of delay: 250 per week']);
	});

	it('says nothing of an absent or empty cost of delay', () => {
		expect(costOfDelayLines(undefined)).toEqual([]);
		expect(costOfDelayLines({})).toEqual([]);
		expect(costOfDelayLines({ inputs: {} })).toEqual([]);
	});

	it('says a whole forecast: duration, delivery, basis, and who set it', () => {
		expect(
			forecastLines({
				duration: '16h',
				delivery: '2026-10-07T17:00:00Z',
				basis: 'Three similar stories took two days each.',
				by: 'planner',
				at: '2026-10-03T09:00:00Z'
			})
		).toEqual([
			'forecast: 16h of work, delivery 2026-10-07T17:00:00Z',
			'Three similar stories took two days each.',
			'set by planner at 2026-10-03T09:00:00Z'
		]);
	});

	it('says a partial forecast', () => {
		expect(forecastLines({ delivery: '2026-10-07T17:00:00Z' })).toEqual([
			'forecast: delivery 2026-10-07T17:00:00Z'
		]);
		expect(forecastLines({ basis: 'A guess.', at: '2026-10-03T09:00:00Z' })).toEqual([
			'forecast',
			'A guess.',
			'set at 2026-10-03T09:00:00Z'
		]);
	});

	it('says nothing of an absent or empty forecast', () => {
		expect(forecastLines(undefined)).toEqual([]);
		expect(forecastLines({})).toEqual([]);
	});
});
