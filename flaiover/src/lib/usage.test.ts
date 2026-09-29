import { describe, expect, it } from 'vitest';
import { count, dollars, modelLine, spent, usageLine, type Usage } from './usage';

const u: Usage = {
	source: 'log',
	seconds: 1083,
	estimated: true,
	models: [
		{
			model: 'claude-opus-5-5',
			input: 256,
			output: 89342,
			cache_read: 19723140,
			cache_write: 327605,
			cost: 8.1258
		}
	]
};

describe('usage', () => {
	it('says what was spent in a line, as flai show does', () => {
		expect(usageLine(u)).toBe(
			"20.1M tokens · $8.13 (estimated) · 18m of agent work · measured from its agents' logs"
		);
		expect(usageLine({ ...u, source: 'sum', estimated: false, seconds: 3725 })).toBe(
			'20.1M tokens · $8.13 · 1h2m of agent work · summed from its children'
		);
		expect(modelLine(u.models[0])).toBe(
			'claude-opus-5-5: 256 in · 89.3K out · 19.7M cache read · 327.6K cache write · $8.13'
		);
	});
	it('counts and prices short', () => {
		expect([count(950), count(12345), count(2.5e9)]).toEqual(['950', '12.3K', '2.5B']);
		expect([dollars(0.0421), dollars(0), dollars(12.5)]).toEqual(['$0.042', '$0.00', '$12.50']);
	});
	it('an empty measurement says nothing', () => {
		expect(spent({ source: 'log', seconds: 0, models: [] })).toBe(false);
		expect(spent(undefined)).toBe(false);
		expect(spent(u)).toBe(true);
	});
});
