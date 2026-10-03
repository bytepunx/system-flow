import { describe, expect, it } from 'vitest';
import { costFields, costInputs, costPatch, costSummary, type CostFields } from './costofdelay';

const none: CostFields = { revenue_per_week: '', penalty_per_week: '', time_lost_per_cycle: '' };

describe('cost of delay fields (S-0204)', () => {
	it('fills the fields from an item’s inputs, empty where an input is absent', () => {
		expect(costFields(undefined)).toEqual(none);
		expect(costFields({ value: 900 })).toEqual(none);
		expect(
			costFields({
				inputs: { revenue_per_week: 1200.5, penalty_per_week: 0, time_lost_per_cycle: '1h30m' },
				value: 900
			})
		).toEqual({ revenue_per_week: '1200.5', penalty_per_week: '0', time_lost_per_cycle: '1h30m' });
	});

	it('says the inputs given in one line with the currency, and nothing when none is', () => {
		expect(costSummary(none, 'USD')).toBe('');
		expect(costSummary({ ...none, time_lost_per_cycle: '  ' }, 'USD')).toBe('');
		expect(
			costSummary(
				{ revenue_per_week: '1200', penalty_per_week: ' 300 ', time_lost_per_cycle: '6h' },
				'EUR'
			)
		).toBe('revenue 1,200 EUR/week · penalty 300 EUR/week · 6h lost per cycle');
		expect(costSummary({ ...none, penalty_per_week: '0' }, 'USD')).toBe('penalty 0 USD/week');
		// what is not a number is shown as typed, for flai to refuse
		expect(costSummary({ ...none, revenue_per_week: 'lots' }, 'USD')).toBe('revenue lots USD/week');
	});

	it('patches only the inputs that changed, an emptied one as "", and nothing when none did', () => {
		const loaded: CostFields = {
			revenue_per_week: '1200',
			penalty_per_week: '',
			time_lost_per_cycle: '6h'
		};
		expect(costPatch(loaded, { ...loaded })).toBeUndefined();
		expect(costPatch(loaded, { ...loaded, revenue_per_week: ' 1200 ' })).toBeUndefined();
		expect(
			costPatch(loaded, { ...loaded, revenue_per_week: '1500', penalty_per_week: '50' })
		).toEqual({ revenue_per_week: '1500', penalty_per_week: '50' });
		expect(costPatch(loaded, none)).toEqual({ revenue_per_week: '', time_lost_per_cycle: '' });
	});

	it('gives a new item only the inputs given, and nothing when none is', () => {
		expect(costInputs(none)).toBeUndefined();
		expect(costInputs({ ...none, penalty_per_week: ' 300 ', time_lost_per_cycle: '6h' })).toEqual({
			penalty_per_week: '300',
			time_lost_per_cycle: '6h'
		});
	});
});
