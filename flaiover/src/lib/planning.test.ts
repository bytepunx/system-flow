import { describe, expect, it } from 'vitest';
import { amount, costOfDelayLines, costOfDelayStale, forecastLines } from './planning';

describe('planning (S-0199)', () => {
	it('says an amount with separators and at most two decimals', () => {
		expect([amount(1500), amount(0), amount(12.345), amount(1234567.5)]).toEqual([
			'1,500',
			'0',
			'12.35',
			'1,234,567.5'
		]);
	});

	// ADR-0080: the inputs and the value each say who set them and when
	it('says a whole cost of delay: the value with its stamp, then the inputs with theirs', () => {
		expect(
			costOfDelayLines({
				inputs: {
					revenue_per_week: 1200,
					penalty_per_week: 300,
					time_lost_per_cycle: '6h',
					by: 'alex',
					at: '2026-10-03T09:00:00Z'
				},
				value: 1650,
				by: 'planner',
				at: '2026-10-03T09:20:00Z'
			})
		).toEqual([
			'cost of delay: 1,650 per week, set by planner at 2026-10-03 05:20 EDT',
			'inputs: revenue 1,200 per week · penalty 300 per week · 6h lost per cycle, set by alex at 2026-10-03 05:00 EDT'
		]);
	});

	it('says there is no value yet when only inputs are given; a zero amount is given', () => {
		expect(
			costOfDelayLines({
				inputs: { penalty_per_week: 0, by: 'alex', at: '2026-10-03T09:00:00Z' }
			})
		).toEqual([
			'cost of delay: no value yet',
			'inputs: penalty 0 per week, set by alex at 2026-10-03 05:00 EDT'
		]);
	});

	it('says a value with no inputs, such as a share apportioned from an epic, alone', () => {
		expect(costOfDelayLines({ value: 250, by: 'planner', at: '2026-10-03T09:20:00Z' })).toEqual([
			'cost of delay: 250 per week, set by planner at 2026-10-03 05:20 EDT'
		]);
		expect(costOfDelayLines({ value: 250 })).toEqual(['cost of delay: 250 per week']);
	});

	it('says nothing of an absent or empty cost of delay', () => {
		expect(costOfDelayLines(undefined)).toEqual([]);
		expect(costOfDelayLines({})).toEqual([]);
		expect(costOfDelayLines({ inputs: {} })).toEqual([]);
	});

	describe('a stale value (ADR-0080)', () => {
		const value = { value: 1650, by: 'planner', at: '2026-10-03T09:20:00Z' };
		const inputs = (at: string) => ({ revenue_per_week: 1200, by: 'alex', at });

		it('is stale when the inputs were set after the value', () => {
			expect(costOfDelayStale({ ...value, inputs: inputs('2026-10-03T09:20:01Z') })).toBe(true);
		});

		it('is not stale when the inputs were set before the value, or at the same time', () => {
			expect(costOfDelayStale({ ...value, inputs: inputs('2026-10-03T09:00:00Z') })).toBe(false);
			expect(costOfDelayStale({ ...value, inputs: inputs('2026-10-03T09:20:00Z') })).toBe(false);
		});

		it('is not stale without a value or without inputs', () => {
			expect(costOfDelayStale({ inputs: inputs('2026-10-03T10:00:00Z') })).toBe(false);
			expect(costOfDelayStale(value)).toBe(false);
			expect(
				costOfDelayStale({ ...value, inputs: { by: 'alex', at: '2026-10-03T10:00:00Z' } })
			).toBe(false);
			expect(costOfDelayStale(undefined)).toBe(false);
		});

		it('is not stale when either time is not one', () => {
			expect(costOfDelayStale({ ...value, inputs: inputs('soon') })).toBe(false);
			expect(costOfDelayStale({ ...value, inputs: inputs('2026-02-30T10:00:00Z') })).toBe(false);
			expect(costOfDelayStale({ ...value, inputs: { revenue_per_week: 1200 } })).toBe(false);
			expect(
				costOfDelayStale({ ...value, at: 'yesterday', inputs: inputs('2026-10-03T10:00:00Z') })
			).toBe(false);
		});
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
			'forecast: 16h of work, delivery 2026-10-07 13:00 EDT',
			'Three similar stories took two days each.',
			'set by planner at 2026-10-03 05:00 EDT'
		]);
	});

	it('says a partial forecast', () => {
		expect(forecastLines({ delivery: '2026-10-07T17:00:00Z' })).toEqual([
			'forecast: delivery 2026-10-07 13:00 EDT'
		]);
		expect(forecastLines({ basis: 'A guess.', at: '2026-10-03T09:00:00Z' })).toEqual([
			'forecast',
			'A guess.',
			'set at 2026-10-03 05:00 EDT'
		]);
	});

	it('says nothing of an absent or empty forecast', () => {
		expect(forecastLines(undefined)).toEqual([]);
		expect(forecastLines({})).toEqual([]);
	});
});
