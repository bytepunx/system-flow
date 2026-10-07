import { describe, expect, it } from 'vitest';
import { localDate, localTime } from './localtime';

describe('localTime', () => {
	it('shows a recorded time in the zone the tests pin, not in UTC', () => {
		expect(Intl.DateTimeFormat().resolvedOptions().timeZone).toBe('America/New_York');
		expect(localTime('2026-10-07T19:35:12Z')).toBe('2026-10-07 15:35 EDT');
	});

	it('crosses midnight into the previous or the next day by zone', () => {
		expect(localTime('2026-10-07T02:10:00Z', 'America/Los_Angeles')).toBe('2026-10-06 19:10 PDT');
		expect(localTime('2026-10-07T23:10:00Z', 'Asia/Tokyo')).toBe('2026-10-08 08:10 GMT+9');
	});

	it('keeps a half-hour offset', () => {
		expect(localTime('2026-10-07T12:00:00Z', 'Asia/Kolkata')).toBe('2026-10-07 17:30 GMT+5:30');
	});

	it('follows summer and winter time', () => {
		expect(localTime('2026-07-01T12:00:00Z')).toBe('2026-07-01 08:00 EDT');
		expect(localTime('2026-01-15T12:00:00Z')).toBe('2026-01-15 07:00 EST');
	});

	it('shows midnight as 00, not 24', () => {
		expect(localTime('2026-10-07T04:00:00Z')).toBe('2026-10-07 00:00 EDT');
	});

	it('returns what it was given when that does not parse', () => {
		expect(localTime('')).toBe('');
		expect(localTime('-')).toBe('-');
		expect(localTime('2026-13-45T99:99:99Z')).toBe('2026-13-45T99:99:99Z');
	});
});

describe('localDate', () => {
	it('gives the calendar date in the zone', () => {
		expect(localDate('2026-10-07T02:10:00Z')).toBe('2026-10-06');
		expect(localDate('2026-10-07T23:10:00Z', 'Asia/Tokyo')).toBe('2026-10-08');
		expect(localDate('2026-10-07T12:00:00Z', 'UTC')).toBe('2026-10-07');
	});

	it('returns what it was given when that does not parse', () => {
		expect(localDate('soon')).toBe('soon');
	});
});
